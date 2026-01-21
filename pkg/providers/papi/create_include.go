package papi

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
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

type includeOptions struct {
	contractID                string
	includeName               string
	edgercPath                string
	section                   string
	jsonDir                   string
	tfWorkPath                string
	rulesAsHCL                bool
	splitDepth                *int
	ruleFormatVersionOverride string
}

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
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}
	templateToFile := map[string]string{
		"includes.tmpl":  includePath,
		"variables.tmpl": variablesPath,
		"imports.tmpl":   importPath,
	}

	options := includeOptions{
		contractID:  c.Args().First(),
		includeName: c.Args().Get(1),
		section:     edgegrid.GetEdgercSection(c),
		jsonDir:     "property-snippets",
		tfWorkPath:  tfWorkPath,
	}

	if c.IsSet("rules-as-hcl") {
		options.rulesAsHCL = c.Bool("rules-as-hcl")
	}

	if c.IsSet("split-depth") {
		options.splitDepth = ptr.To(c.Int("split-depth"))
	}

	if c.IsSet("rule-format") {
		options.ruleFormatVersionOverride = c.String("rule-format")
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	var multiTargetProcessor templates.MultiTargetProcessor
	if options.splitDepth != nil {
		multiTargetProcessor = templates.FSMultiTargetProcessor{
			TemplatesFS:     templateFiles,
			AdditionalFuncs: additionalFuncs,
		}
		err = createSplitRulesDir(tfWorkPath)
		if err != nil {
			return cli.Exit(color.RedString("Error creating directory for include rules: %s", err), 1)
		}
	}

	if err = createInclude(ctx, options, client, processor, multiTargetProcessor); err != nil {
		return cli.Exit(color.RedString("Error exporting include: %s", err), 1)
	}

	return nil
}

func createInclude(ctx context.Context, options includeOptions, client papi.PAPI, processor templates.TemplateProcessor, multiTargetProcessor templates.MultiTargetProcessor) error {
	term := terminal.Get(ctx)

	multiTargetData := make(templates.MultiTargetData)

	tfData := TFData{
		Includes:   make([]TFIncludeData, 0),
		EdgercPath: options.edgercPath,
		Section:    options.section,
		RulesAsHCL: options.rulesAsHCL,
	}

	// Get Include
	term.Spinner().Start("Fetching include " + options.includeName)
	include, err := findIncludeByName(ctx, client, options.contractID, options.includeName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrIncludeNotFound, err)
	}
	term.Spinner().OK()

	includeData, rules, err := getIncludeData(ctx, include, client, options.ruleFormatVersionOverride)
	if err != nil {
		return err
	}

	// Save snippets
	if !options.rulesAsHCL {
		term.Spinner().Start("Saving snippets ")
		ruleTemplate, rulesTemplate := setIncludeRuleTemplates(rules)
		if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(options.tfWorkPath, options.jsonDir), fmt.Sprintf("%s.json", include.IncludeName)); err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
		}
		term.Spinner().OK()
	} else {
		wrappedRules := wrapAndNameRules(includeData.IncludeName, rules.Rules)
		if options.splitDepth != nil {
			includeData.RootRule = wrappedRules.TerraformName
			multiTargetData.AddData("includes_rules.tmpl", prepareRulesForSplitRule(wrappedRules, *options.splitDepth, options.tfWorkPath, includeRuleWrapper))
		} else {
			includeData.Rules = flattenRules(wrappedRules)
		}

	}

	tfData.Includes = append(tfData.Includes, *includeData)
	filterFuncs := make([]func([]string) ([]string, error), 0)
	if options.rulesAsHCL {
		ruleTemplate := fmt.Sprintf("rules_%s.tmpl", rules.RuleFormat)
		if !processor.TemplateExists(ruleTemplate) {
			return fmt.Errorf("%w: %s", ErrUnsupportedRuleFormat, rules.RuleFormat)
		}
		if options.splitDepth == nil {
			processor.AddTemplateTarget(ruleTemplate, filepath.Join(options.tfWorkPath, "rules.tf"))
			processor.AddTemplateTarget("includes_rules.tmpl", filepath.Join(options.tfWorkPath, "includes_rules.tf"))
		} else {
			processor.AddTemplateTarget("rules_module.tmpl", filepath.Join(options.tfWorkPath, "rules", "module_config.tf"))
			tfData.UseSplitDepth = true
		}
		filterFuncs = append(filterFuncs, useThisOnlyRuleFormat(rules.RuleFormat))
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = processor.ProcessTemplates(tfData, filterFuncs...); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	if options.splitDepth != nil {
		if err = multiTargetProcessor.ProcessTemplates(multiTargetData, filterFuncs...); err != nil {
			term.Spinner().Fail()
			if _, err := CheckErrors(); err != nil {
				return fmt.Errorf("%w", err)
			}
			return fmt.Errorf("%w: %s", ErrSavingFiles, err)
		}
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for include '%s' was saved successfully\n", includeData.IncludeName)

	return nil
}

func getIncludeData(ctx context.Context, include *papi.Include, client papi.PAPI, ruleFormatVersionOverride string) (*TFIncludeData, *papi.GetIncludeRuleTreeResponse, error) {
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
	rules, err := getIncludeRules(ctx, client, include, latestVersion, ruleFormatVersionOverride)
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
		ProductID:   latestVersion.IncludeVersion.ProductID,
		RuleFormat:  rules.RuleFormat,
	}

	if latestStagingActivation != nil {
		includeData.StagingInfo.ActivationNote = latestStagingActivation.Note
		includeData.StagingInfo.Emails = latestStagingActivation.NotifyEmails
		includeData.StagingInfo.Version = latestStagingActivation.IncludeVersion
		includeData.StagingInfo.HasActivation = true
		includeData.StagingInfo.IsActiveOnLatestVersion = latestStagingActivation.IncludeVersion == include.LatestVersion
	}

	if latestProdActivation != nil {
		includeData.ProductionInfo.ActivationNote = latestProdActivation.Note
		includeData.ProductionInfo.Emails = latestProdActivation.NotifyEmails
		includeData.ProductionInfo.Version = latestProdActivation.IncludeVersion
		includeData.ProductionInfo.HasActivation = true
		includeData.ProductionInfo.IsActiveOnLatestVersion = latestProdActivation.IncludeVersion == include.LatestVersion
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
		Comments:       rules.Comments,
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

func getIncludeRules(ctx context.Context, client papi.PAPI, include *papi.Include, latestVersion *papi.GetIncludeVersionResponse, ruleFormatVersionOverride string) (*papi.GetIncludeRuleTreeResponse, error) {
	ruleFormatVersion := latestVersion.IncludeVersion.RuleFormat
	if ruleFormatVersionOverride != "" {
		ruleFormatVersion = ruleFormatVersionOverride
	}

	return client.GetIncludeRuleTree(ctx, papi.GetIncludeRuleTreeRequest{
		ContractID:     include.ContractID,
		GroupID:        include.GroupID,
		IncludeID:      include.IncludeID,
		IncludeVersion: include.LatestVersion,
		RuleFormat:     ruleFormatVersion,
	})
}
