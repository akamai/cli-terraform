package papi

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrFetchingActivations is returned when fetching activations of include request failed
	ErrFetchingActivations = errors.New("fetching include activations")
	// ErrFetchingLatestIncludeVersion is returned when fetching latest version of include request failed
	ErrFetchingLatestIncludeVersion = errors.New("fetching latest include version")
	// ErrIncludeNotFound is returned when include couldn't be found
	ErrIncludeNotFound = errors.New("include name not found")
	// ErrIncludeRulesNotFound is returned when include rules couldn't be found
	ErrIncludeRulesNotFound = errors.New("include rules not found")
)

// CmdCreateInclude is an entrypoint to export-include include sub-command
func CmdCreateInclude(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := papi.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	includePath := filepath.Join(tfWorkPath, "includes.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(includePath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"includes.tmpl":  includePath,
		"variables.tmpl": variablesPath,
		"imports.tmpl":   importPath,
	}

	var schema bool
	if c.IsSet("schema") {
		schema = c.Bool("schema")
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	contractID := c.Args().First()
	includeName := c.Args().Get(1)
	section := edgegrid.GetEdgercSection(c)

	if err = createInclude(ctx, contractID, includeName, section, "property-snippets", tfWorkPath, schema, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting include: %s", err)), 1)
	}

	return nil
}

func createInclude(ctx context.Context, contractID, includeName, section, jsonDir, tfWorkPath string, schema bool, client papi.PAPI, processor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	tfData := TFData{
		Includes:      make([]TFIncludeData, 0),
		Section:       section,
		RulesAsSchema: schema,
	}

	// Get Include
	term.Spinner().Start("Fetching include " + includeName)
	include, err := findIncludeByName(ctx, client, contractID, includeName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrIncludeNotFound, err)
	}
	term.Spinner().OK()

	includeData, rules, err := getIncludeData(ctx, include, client)
	if err != nil {
		return err
	}

	// Save snippets
	if !schema {
		term.Spinner().Start("Saving snippets ")
		ruleTemplate, rulesTemplate := setIncludeRuleTemplates(rules)
		if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(tfWorkPath, jsonDir), fmt.Sprintf("%s.json", include.IncludeName)); err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
		}
		term.Spinner().OK()
	} else {
		includeData.Rules = flattenRules(includeData.IncludeName, rules.Rules)
	}

	tfData.Includes = append(tfData.Includes, *includeData)
	filterFuncs := make([]func([]string) ([]string, error), 0)
	if schema {
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

func getIncludeData(ctx context.Context, include *papi.Include, client papi.PAPI) (*TFIncludeData, *papi.GetIncludeRuleTreeResponse, error) {
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
		return nil, nil, fmt.Errorf("%w: %s", ErrFetchingLatestIncludeVersion, err)
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
		return nil, nil, fmt.Errorf("%w: %s", ErrIncludeRulesNotFound, err)
	}
	term.Spinner().OK()

	// Get activations
	term.Spinner().Start("Fetching include activations ")
	results, err := client.ListIncludeActivations(ctx, papi.ListIncludeActivationsRequest{
		IncludeID:  include.IncludeID,
		ContractID: include.ContractID,
		GroupID:    include.GroupID,
	})
	if err != nil {
		term.Spinner().Fail()
		return nil, nil, fmt.Errorf("%w: %s", ErrFetchingActivations, err)
	}
	activations := results.Activations.Items
	term.Spinner().OK()

	var stagingActivations, prodActivations []papi.IncludeActivation
	var latestStagingActivation, latestProdActivation *papi.IncludeActivation

	stagingActivations = filterIncludeActivationsByNetwork(activations, papi.ActivationNetworkStaging)
	latestStagingActivation = findLatestIncludeActivation(stagingActivations)

	prodActivations = filterIncludeActivationsByNetwork(activations, papi.ActivationNetworkProduction)
	latestProdActivation = findLatestIncludeActivation(prodActivations)

	// Populate TFIncludeData
	includeData := TFIncludeData{
		ContractID:  include.ContractID,
		GroupID:     include.GroupID,
		IncludeID:   include.IncludeID,
		IncludeName: include.IncludeName,
		IncludeType: string(include.IncludeType),
		Networks:    getActivatedNetworks(include),
		RuleFormat:  latestVersion.IncludeVersion.RuleFormat,
	}

	if latestStagingActivation != nil {
		includeData.ActivationNoteStaging = latestStagingActivation.Note
		includeData.ActivationEmailsStaging = latestStagingActivation.NotifyEmails
		includeData.VersionStaging = strconv.Itoa(latestStagingActivation.IncludeVersion)
	}

	if latestProdActivation != nil {
		includeData.ActivationNoteProduction = latestProdActivation.Note
		includeData.ActivationEmailsProduction = latestProdActivation.NotifyEmails
		includeData.VersionProduction = strconv.Itoa(latestProdActivation.IncludeVersion)
	}

	term.Spinner().OK()

	return &includeData, rules, nil
}

// findLatestIncludeActivation finds the latest activation of type `ACTIVATE` with status `ACTIVE` or `PENDING`.
// If it encounters activation of type `DEACTIVATE` with status `ACTIVE` first or does not find any activation of type
// `ACTIVATE` with `ACTIVE` status, it returns nil
func findLatestIncludeActivation(activations []papi.IncludeActivation) *papi.IncludeActivation {
	if len(activations) == 0 {
		return nil
	}

	sort.Slice(activations, func(i, j int) bool {
		return activations[i].UpdateDate > activations[j].UpdateDate
	})

	for _, activation := range activations {
		if activation.ActivationType == papi.ActivationTypeActivate &&
			(activation.Status == papi.ActivationStatusActive || activation.Status == papi.ActivationStatusPending) {
			return &activation
		}
		if activation.ActivationType == papi.ActivationTypeDeactivate &&
			(activation.Status == papi.ActivationStatusActive || activation.Status == papi.ActivationStatusPending) {
			return nil
		}
	}

	return nil
}

// filterIncludeActivationsByNetwork filters list of activations based on given network
func filterIncludeActivationsByNetwork(activations []papi.IncludeActivation, network papi.ActivationNetwork) []papi.IncludeActivation {
	var filteredActivations []papi.IncludeActivation
	for _, activation := range activations {
		if activation.Network == network {
			filteredActivations = append(filteredActivations, activation)
		}
	}

	return filteredActivations
}

// setIncludeRuleTemplates creates templates based on RuleTemplate and RulesTemplate for given include rule tree response
func setIncludeRuleTemplates(rules *papi.GetIncludeRuleTreeResponse) (RuleTemplate, RulesTemplate) {
	// Set up template structure
	ruleTemplate := RuleTemplate{
		Name:                rules.Rules.Name,
		Criteria:            rules.Rules.Criteria,
		Behaviors:           rules.Rules.Behaviors,
		Comments:            rules.Rules.Comments,
		CriteriaLocked:      rules.Rules.CriteriaLocked,
		CriteriaMustSatisfy: rules.Rules.CriteriaMustSatisfy,
		UUID:                rules.Rules.UUID,
		Variables:           rules.Rules.Variables,
		AdvancedOverride:    rules.Rules.AdvancedOverride,
		Children:            make([]string, 0),
		Options:             rules.Rules.Options,
	}

	rulesTemplate := RulesTemplate{
		AccountID:      rules.AccountID,
		ContractID:     rules.ContractID,
		GroupID:        rules.GroupID,
		IncludeID:      rules.IncludeID,
		IncludeVersion: rules.IncludeVersion,
		IncludeType:    string(rules.IncludeType),
		Etag:           rules.Etag,
		RuleFormat:     rules.RuleFormat,
	}

	return ruleTemplate, rulesTemplate

}

// findIncludeByName searches for an include with a given name and contract_id
func findIncludeByName(ctx context.Context, client papi.PAPI, contractID, includeName string) (*papi.Include, error) {
	result, err := client.ListIncludes(ctx, papi.ListIncludesRequest{
		ContractID: contractID,
	})
	if err != nil {
		return nil, err
	}

	var include *papi.Include
	for _, i := range result.Includes.Items {
		if i.IncludeName == includeName {
			include = &i
			return include, nil
		}
	}

	return nil, fmt.Errorf("unable to find include: \"%s\"", includeName)
}

// getActivatedNetworks returns a list of networks on which the given include is activated
func getActivatedNetworks(include *papi.Include) []string {
	var result []string

	if include.StagingVersion != nil {
		result = append(result, string(papi.ActivationNetworkStaging))
	}

	if include.ProductionVersion != nil {
		result = append(result, string(papi.ActivationNetworkProduction))
	}

	return result
}
