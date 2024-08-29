package papi

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrIncludeRuleNotFound is returned when an include rule couldn't be found
	ErrIncludeRuleNotFound = errors.New("include rule not found")
)

// CmdCreateIncludeRule is an entrypoint to export-include-rule sub-command
func CmdCreateIncludeRule(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := papi.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	includePath := filepath.Join(tfWorkPath, "includes.tf")
	err := tools.CheckFiles(includePath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"includes.tmpl": includePath,
	}

	var rulesAsHCL bool
	if c.IsSet("rules-as-hcl") {
		rulesAsHCL = c.Bool("rules-as-hcl")
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	contractID := c.Args().First()
	includeName := c.Args().Get(1)
	ruleName := c.Args().Get(2)
	section := edgegrid.GetEdgercSection(c)

	if err = createIncludeRule(ctx, contractID, includeName, ruleName, section, "property-snippets", tfWorkPath, rulesAsHCL, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting include: %s", err)), 1)
	}

	return nil
}

func createIncludeRule(ctx context.Context, contractID, includeName, ruleName, section, jsonDir, tfWorkPath string, rulesAsHCL bool, client papi.PAPI, processor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	var includeData TFIncludeData
	tfData := TFData{
		Includes:   make([]TFIncludeData, 0),
		Section:    section,
		RulesAsHCL: rulesAsHCL,
	}

	// Get Include
	term.Spinner().Start("Fetching include " + includeName)
	include, err := findIncludeByName(ctx, client, contractID, includeName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrIncludeNotFound, err)
	}
	term.Spinner().OK()

	rules, err := getIncludeRuleData(ctx, include, ruleName, client)
	if err != nil {
		return err
	}

	// Save snippets
	if !rulesAsHCL {
		term.Spinner().Start("Saving snippets ")
		ruleTemplate, rulesTemplate := setIncludeRuleTemplates(rules)
		if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(tfWorkPath, jsonDir), fmt.Sprintf("%s.json", include.IncludeName)); err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
		}
		term.Spinner().OK()
	} else {
		dummyDefaultRule := papi.Rules{
			Children: []papi.Rules{
				rules.Rules,
			},
		}
		includeData.Rules = flattenRules(ruleName, dummyDefaultRule)
	}

	tfData.Includes = append(tfData.Includes, includeData)
	filterFuncs := make([]func([]string) ([]string, error), 0)
	if rulesAsHCL {
		ruleTemplate := fmt.Sprintf("rules_%s.tmpl", rules.RuleFormat)
		if !processor.TemplateExists(ruleTemplate) {
			return fmt.Errorf("%w: %s", ErrUnsupportedRuleFormat, rules.RuleFormat)
		}
		processor.AddTemplateTarget(ruleTemplate, filepath.Join(tfWorkPath, "rules.tf"))
		processor.AddTemplateTarget("includes_rules.tmpl", filepath.Join(tfWorkPath, "includes_rules.tf"))
		filterFuncs = append(filterFuncs, useThisOnlyRuleFormat(rules.RuleFormat))
	}
	term.Spinner().Start("Saving TF configurations ")
	if err = processor.ProcessTemplates(tfData, filterFuncs...); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for include '%s' was saved successfully\n", includeData.IncludeName)

	return nil
}

func getIncludeRuleData(ctx context.Context, include *papi.Include, ruleName string, client papi.PAPI) (*papi.GetIncludeRuleTreeResponse, error) {
	term := terminal.Get(ctx)

	// Get the latest version of include
	term.Spinner().Start("Fetching the latest version of include ")
	latestVersion, err := client.GetIncludeVersion(ctx, papi.GetIncludeVersionRequest{
		ContractID: include.ContractID,
		GroupID:    include.GroupID,
		IncludeID:  include.IncludeID,
		Version:    include.LatestVersion,
	})
	if err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", ErrFetchingLatestIncludeVersion, err)
	}

	// Get include rules
	term.Spinner().Start("Fetching include rules ")
	rules, err := client.GetIncludeRuleTree(ctx, papi.GetIncludeRuleTreeRequest{
		ContractID:     include.ContractID,
		GroupID:        include.GroupID,
		IncludeID:      include.IncludeID,
		IncludeVersion: include.LatestVersion,
		RuleFormat:     latestVersion.IncludeVersion.RuleFormat,
	})
	if err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", ErrIncludeRuleNotFound, err)
	}
	term.Spinner().OK()

	singleRule, err := findSingleRule(ctx, ruleName, rules.Rules)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrIncludeRuleNotFound, err)
	}
	rules.Rules = singleRule

	term.Spinner().OK()

	return rules, nil
}

func findSingleRule(ctx context.Context, ruleName string, rules papi.Rules) (papi.Rules, error) {

	if rules.Name == ruleName {
		return rules, nil
	}
	for _, childRule := range rules.Children {
		rules, err := findSingleRule(ctx, ruleName, childRule)
		if err == nil {
			return rules, nil
		}
	}

	return rules, ErrIncludeRulesNotFound
}
