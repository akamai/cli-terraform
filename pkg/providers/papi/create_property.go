// Package papi contains code for exporting properties.
package papi

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// EdgeHostname represents EdgeHostname resource
type EdgeHostname struct {
	EdgeHostname             string
	EdgeHostnameID           string
	ContractID               string
	GroupID                  string
	ID                       string
	IPv6                     string
	EdgeHostnameResourceName string
	SecurityType             string
	UseCases                 string
	CertificateID            int64
	TTL                      int
}

// Hostname represents edge hostname resource
type Hostname struct {
	CnameFrom                string
	CnameTo                  string
	EdgeHostnameResourceName string
	CertProvisioningType     string
	IsActive                 bool
}

// WrappedRules is a wrapper around Rule which simplifies flattening rule tree into list and adjust names of the datasources
type WrappedRules struct {
	Rule          papi.Rules
	TerraformName string
	FileName      string
	Children      []*WrappedRules
	IsRoot        bool
}

// TFData holds template data
type TFData struct {
	Includes      []TFIncludeData
	Property      TFPropertyData
	Section       string
	Rules         []*WrappedRules
	RulesAsHCL    bool
	UseBootstrap  bool
	UseSplitDepth bool
	RootRule      string
}

// TFIncludeData holds template data for include
type TFIncludeData struct {
	ContractID     string
	GroupID        string
	IncludeID      string
	IncludeName    string
	IncludeType    string
	RuleFormat     string
	Rules          []*WrappedRules
	ProductID      string
	ProductionInfo NetworkInfo
	StagingInfo    NetworkInfo
	RootRule       string
}

// TFPropertyData holds template data for property
type TFPropertyData struct {
	GroupName            string
	GroupID              string
	ContractID           string
	PropertyResourceName string
	PropertyName         string
	PropertyID           string
	ProductID            string
	ProductName          string
	RuleFormat           string
	IsSecure             string
	EdgeHostnames        map[string]EdgeHostname
	Hostnames            map[string]Hostname
	UseHostnameBucket    bool
	ReadVersion          string
	ProductionInfo       NetworkInfo
	StagingInfo          NetworkInfo
}

// NetworkInfo holds details for specific network
type NetworkInfo struct {
	Emails                  []string
	ActivationNote          string
	HasActivation           bool
	Version                 int
	IsActiveOnLatestVersion bool
}

// RulesTemplate represent data used for rules
type RulesTemplate struct {
	AccountID       string        `json:"accountId"`
	ContractID      string        `json:"contractId"`
	GroupID         string        `json:"groupId"`
	PropertyID      string        `json:"propertyId,omitempty"`
	IncludeID       string        `json:"includeId,omitempty"`
	PropertyVersion int           `json:"propertyVersion,omitempty"`
	IncludeVersion  int           `json:"includeVersion,omitempty"`
	IncludeType     string        `json:"includeType,omitempty"`
	Etag            string        `json:"etag"`
	RuleFormat      string        `json:"ruleFormat"`
	Comments        string        `json:"comments,omitempty"`
	Rule            *RuleTemplate `json:"rules"`
	Errors          []*papi.Error `json:"errors,omitempty"`
}

// RuleTemplate represent data used for single rule
type RuleTemplate struct {
	Name                string                       `json:"name"`
	Criteria            []papi.RuleBehavior          `json:"criteria,omitempty"`
	Behaviors           []papi.RuleBehavior          `json:"behaviors,omitempty"`
	Children            []string                     `json:"children,omitempty"`
	Comments            string                       `json:"comments,omitempty"`
	CriteriaLocked      bool                         `json:"criteriaLocked,omitempty"`
	CriteriaMustSatisfy papi.RuleCriteriaMustSatisfy `json:"criteriaMustSatisfy,omitempty"`
	UUID                string                       `json:"uuid,omitempty"`
	Variables           []papi.RuleVariable          `json:"variables,omitempty"`
	AdvancedOverride    string                       `json:"advancedOverride,omitempty"`

	Options struct {
		IsSecure bool `json:"is_secure,omitempty"`
	} `json:"options,omitempty"`

	CustomOverride *papi.RuleCustomOverride `json:"customOverride,omitempty"`
}

type propertyOptions struct {
	propertyName  string
	section       string
	tfWorkPath    string
	version       string
	rulesAsHCL    bool
	withBootstrap bool
	splitDepth    *int
}

type splitDepthRuleWrapper func([]*WrappedRules) TFData

//go:embed templates/*
var templateFiles embed.FS

// normalizeRuleNameRegexp is a regexp for finding invalid characters in a path created from the rule name
var normalizeRuleNameRegexp = regexp.MustCompile(`[^\w-.]`)

var (
	// ErrHostnamesNotFound is returned when hostnames couldn't be found
	ErrHostnamesNotFound = errors.New("hostnames not found")
	// ErrPropertyVersionNotFound is returned when property version couldn't be found
	ErrPropertyVersionNotFound = errors.New("property version not found")
	// ErrPropertyVersionNotValid is returned when property version couldn't be found
	ErrPropertyVersionNotValid = errors.New("property version not valid")
	// ErrProductNameNotFound is returned when product couldn't be found
	ErrProductNameNotFound = errors.New("product name not found")
	// ErrFetchingActivationDetails is returned when fetching activation details request failed
	ErrFetchingActivationDetails = errors.New("fetching activations")
	// ErrFetchingHostnameDetails is returned when fetching hostname details request failed
	ErrFetchingHostnameDetails = errors.New("fetching hostnames")
	// ErrSavingSnippets is returned when error appeared while saving property snippet JSON files
	ErrSavingSnippets = errors.New("saving snippets")
	// ErrPropertyRulesNotFound is returned when property rules couldn't be found
	ErrPropertyRulesNotFound = errors.New("property rules not found")
	// ErrGroupNotFound is returned when group couldn't be found
	ErrGroupNotFound = errors.New("group not found")
	// ErrPropertyNotFound is returned when property couldn't be found
	ErrPropertyNotFound = errors.New("property not found")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving terraform project files")
	// ErrUnsupportedRuleFormat is returned when there is no template for provided rule format
	ErrUnsupportedRuleFormat = errors.New("unsupported rule format")
	errCreateRulesDirectory  = errors.New("create rule directory")
	errReadRuleMode          = errors.New("reading export directory mode")
)

var additionalFuncs = tools.DecorateWithMultilineHandlingFunctions(
	map[string]any{
		"TerraformName": tools.TerraformName,
		"AsInt":         AsInt,
		"ReportError":   ReportError,
		"CheckErrors":   CheckErrors,
	})

var propertyRuleWrapper splitDepthRuleWrapper = func(rules []*WrappedRules) TFData {
	return TFData{
		RulesAsHCL:    true,
		UseSplitDepth: true,
		Rules:         rules,
	}
}

var includeRuleWrapper splitDepthRuleWrapper = func(rules []*WrappedRules) TFData {
	return TFData{
		RulesAsHCL:    true,
		UseSplitDepth: true,
		Includes: []TFIncludeData{
			{
				Rules: rules,
			},
		},
	}
}

// CmdCreateProperty is an entrypoint to create-property command
func CmdCreateProperty(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := papi.Client(sess)
	clientHapi := hapi.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	var version string
	if c.IsSet("version") {
		version = c.String("version")
	}

	propertyPath := filepath.Join(tfWorkPath, "property.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(propertyPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"property.tmpl":  propertyPath,
		"variables.tmpl": variablesPath,
		"imports.tmpl":   importPath,
	}

	var rulesAsHCL bool
	if c.IsSet("rules-as-hcl") {
		rulesAsHCL = c.Bool("rules-as-hcl")
	}

	var splitDepth *int
	if c.IsSet("split-depth") {
		splitDepth = ptr.To(c.Int("split-depth"))
	}

	var isBootstrap bool
	if c.IsSet("akamai-property-bootstrap") {
		isBootstrap = c.Bool("akamai-property-bootstrap")
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	var multiTargetProcessor templates.MultiTargetProcessor

	if splitDepth != nil {
		multiTargetProcessor = templates.FSMultiTargetProcessor{
			TemplatesFS:     templateFiles,
			AdditionalFuncs: additionalFuncs,
		}
		err = createSplitRulesDir(tfWorkPath)
		if err != nil {
			return cli.Exit(color.RedString("Error creating directory for rules: %s", err), 1)
		}
	}

	options := propertyOptions{
		propertyName:  c.Args().First(),
		section:       edgegrid.GetEdgercSection(c),
		tfWorkPath:    tfWorkPath,
		version:       version,
		rulesAsHCL:    rulesAsHCL,
		withBootstrap: isBootstrap,
		splitDepth:    splitDepth,
	}
	if err = createProperty(ctx, options, "property-snippets", client, clientHapi, processor, multiTargetProcessor); err != nil {
		return cli.Exit(color.RedString("Error exporting property \"%s\": %s", options.propertyName, err), 1)
	}
	return nil
}

func createSplitRulesDir(tfWorkPath string) error {
	stat, err := os.Stat(tfWorkPath)
	if err != nil {
		return fmt.Errorf("%w: %s", errReadRuleMode, err)
	}
	err = os.Mkdir(filepath.Join(tfWorkPath, "rules"), stat.Mode())
	if err != nil {
		return fmt.Errorf("%w: %s", errCreateRulesDirectory, err)
	}
	return nil
}

func isHostnameBucketProperty(p *papi.Property) bool {
	return p.PropertyType != nil && *p.PropertyType == "HOSTNAME_BUCKET"
}

//nolint:gocyclo
func createProperty(ctx context.Context, options propertyOptions, jsonDir string, client papi.PAPI, clientHapi hapi.HAPI, templateProcessor templates.TemplateProcessor, multiTargetProcessor templates.MultiTargetProcessor) error {
	term := terminal.Get(ctx)

	tfData := TFData{
		Property: TFPropertyData{
			EdgeHostnames: make(map[string]EdgeHostname),
		},
		Section:       options.section,
		RulesAsHCL:    options.rulesAsHCL,
		UseBootstrap:  options.withBootstrap,
		UseSplitDepth: options.splitDepth != nil,
	}

	// Get Property
	term.Spinner().Start("Fetching property " + options.propertyName)
	property, err := findProperty(ctx, client, options.propertyName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrPropertyNotFound, err)
	}

	tfData.Property.ContractID = property.ContractID
	tfData.Property.PropertyName = property.PropertyName
	tfData.Property.PropertyID = property.PropertyID
	tfData.Property.PropertyResourceName = formatResourceName(property.PropertyName)
	if isHostnameBucketProperty(property) {
		tfData.Property.UseHostnameBucket = true
	}

	term.Spinner().OK()

	// Get Group
	term.Spinner().Start("Fetching group ")
	group, err := getGroup(ctx, client, property.GroupID)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrGroupNotFound, err)
	}

	tfData.Property.GroupName = group.GroupName
	tfData.Property.GroupID = group.GroupID

	term.Spinner().OK()

	if options.version == "" {
		options.version = "LATEST"
	}

	// Get Version
	term.Spinner().Start("Fetching property version ")
	version, latestVersion, err := getVersion(ctx, client, property, options.version)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrPropertyVersionNotFound, err)
	}

	tfData.Property.ProductID = version.Version.ProductID
	tfData.Property.ReadVersion = options.version

	term.Spinner().OK()

	multiTargetData := make(templates.MultiTargetData)

	// Get Property Rules
	term.Spinner().Start("Fetching property rules ")
	rules, err := getPropertyRules(ctx, client, version)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrPropertyRulesNotFound, err)
	}

	tfData.Property.IsSecure = "false"
	if rules.Rules.Options.IsSecure {
		tfData.Property.IsSecure = "true"
	}

	// Get Rule Format
	tfData.Property.RuleFormat = rules.RuleFormat

	term.Spinner().OK()

	// Get Product
	term.Spinner().Start("Fetching product name ")
	product, err := getProduct(ctx, client, tfData.Property.ProductID, property.ContractID)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrProductNameNotFound, err)
	}

	tfData.Property.ProductName = product.ProductName

	term.Spinner().OK()

	// Get Hostnames
	if !isHostnameBucketProperty(property) {
		term.Spinner().Start("Fetching hostnames ")
		hostnames, err := getPropertyVersionHostnames(ctx, client, property, version)
		if err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrHostnamesNotFound, err)
		}

		tfData.Property.Hostnames, tfData.Property.EdgeHostnames, err =
			getEdgeHostnameDetail(ctx, client, clientHapi, hostnames, property)
		if err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrFetchingHostnameDetails, err)
		}

		term.Spinner().OK()
	}

	term.Spinner().Start("Fetching activation details ")

	activeStagingActivation, err := fetchActiveActivationForNetwork(ctx, client, property, papi.ActivationNetworkStaging)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingActivationDetails, err)
	}
	if activeStagingActivation != nil {
		tfData.Property.StagingInfo.ActivationNote = activeStagingActivation.Note
		tfData.Property.StagingInfo.Emails = getContactEmails(activeStagingActivation)
		tfData.Property.StagingInfo.Version = activeStagingActivation.PropertyVersion
		tfData.Property.StagingInfo.HasActivation = true
		tfData.Property.StagingInfo.IsActiveOnLatestVersion = activeStagingActivation.PropertyVersion == latestVersion.Version.PropertyVersion
	}
	activeProductionActivation, err := fetchActiveActivationForNetwork(ctx, client, property, papi.ActivationNetworkProduction)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingActivationDetails, err)
	}
	if activeProductionActivation != nil {
		tfData.Property.ProductionInfo.ActivationNote = activeProductionActivation.Note
		tfData.Property.ProductionInfo.Emails = getContactEmails(activeProductionActivation)
		tfData.Property.ProductionInfo.Version = activeProductionActivation.PropertyVersion
		tfData.Property.ProductionInfo.HasActivation = true
		tfData.Property.ProductionInfo.IsActiveOnLatestVersion = activeProductionActivation.PropertyVersion == latestVersion.Version.PropertyVersion
	}

	term.Spinner().OK()

	filterFuncs := make([]func([]string) ([]string, error), 0)
	if options.rulesAsHCL {
		ruleTemplate := fmt.Sprintf("rules_%s.tmpl", rules.RuleFormat)
		if !templateProcessor.TemplateExists(ruleTemplate) {
			return fmt.Errorf("%w: \"%s\"", ErrUnsupportedRuleFormat, rules.RuleFormat)
		}
		filterFuncs = append(filterFuncs, useThisOnlyRuleFormat(rules.RuleFormat))
		wrappedRules := wrapAndNameRules(tfData.Property.PropertyName, rules.Rules)

		if options.splitDepth != nil {
			tfData.RootRule = wrappedRules.TerraformName
			multiTargetData.AddData("split-depth-rules.tmpl", prepareRulesForSplitRule(wrappedRules, *options.splitDepth, options.tfWorkPath, propertyRuleWrapper))
			templateProcessor.AddTemplateTarget("rules_module.tmpl", filepath.Join(options.tfWorkPath, "rules", "module_config.tf"))
		} else {
			tfData.Rules = flattenRules(wrappedRules)
			templateProcessor.AddTemplateTarget(ruleTemplate, filepath.Join(options.tfWorkPath, "rules.tf"))
		}
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData, filterFuncs...); err != nil {
		term.Spinner().Fail()
		if _, err := CheckErrors(); err != nil {
			return fmt.Errorf("%w", err)
		}
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	if !options.rulesAsHCL {
		// Save snippets
		ruleTemplate, rulesTemplate := setPropertyRuleTemplates(rules)
		if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(options.tfWorkPath, jsonDir), "main.json"); err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
		}
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
	term.Printf("Terraform configuration for property \"%s\" was saved successfully\n", property.PropertyName)

	return nil
}

func useThisOnlyRuleFormat(acceptedFormat string) func([]string) ([]string, error) {
	reg := regexp.MustCompile(`rules_(v\d{4}-\d{2}-\d{2}).tmpl`)
	return func(input []string) ([]string, error) {
		res := make([]string, 0)
		formatFound := false
		for _, v := range input {
			if reg.MatchString(v) {
				submatch := reg.FindStringSubmatch(v)
				if submatch[1] == acceptedFormat {
					res = append(res, v)
					formatFound = true
				}
			} else {
				res = append(res, v)
			}
		}

		if !formatFound {
			return nil, fmt.Errorf("did not find %s format among %s", acceptedFormat, input)
		}

		return res, nil
	}
}

func wrapAndNameRules(property string, rule papi.Rules) *WrappedRules {
	wrappedRules := wrapRules(rule)
	wrappedRules.IsRoot = true

	_ = setNamesOnAllRules(property, wrappedRules, map[string]int{}, tools.TerraformName(property), true)
	return wrappedRules
}

func flattenRules(rules *WrappedRules) []*WrappedRules {
	var result []*WrappedRules

	result = append(result, rules)
	result = append(result, flattenWrappedRules(rules)...)

	return result
}

func setNamesOnAllRules(propertyName string, rules *WrappedRules, nameOccurrence map[string]int, parentName string, isRoot bool) map[string]int {
	if isRoot {
		nameOccurrence = setNameOnRule(propertyName, rules, nameOccurrence, parentName)
	}

	for _, child := range rules.Children {
		nameOccurrence = setNameOnRule(propertyName, child, nameOccurrence, rules.FileName)
	}

	for _, child := range rules.Children {
		nameOccurrence = setNamesOnAllRules(propertyName, child, nameOccurrence, rules.FileName, false)
	}
	return nameOccurrence
}

func setNameOnRule(propertyName string, rules *WrappedRules, nameOccurrence map[string]int, parentName string) map[string]int {
	name := tools.TerraformName(rules.Rule.Name)
	nameOccurrence[name]++
	if count := nameOccurrence[name]; count > 1 {
		name = fmt.Sprintf("%s%d", name, count-1)
	}
	rules.TerraformName = fmt.Sprintf("%s_rule_%s", tools.TerraformName(propertyName), name)
	rules.FileName = fmt.Sprintf("%s_%s", parentName, name)
	return nameOccurrence
}

func wrapRules(rule papi.Rules) *WrappedRules {
	var children []*WrappedRules
	for _, child := range rule.Children {
		children = append(children, wrapRules(child))
	}

	return &WrappedRules{
		Rule:          rule,
		TerraformName: rule.Name,
		Children:      children,
	}
}

func flattenWrappedRules(rule *WrappedRules) []*WrappedRules {
	var result = make([]*WrappedRules, 0)

	result = append(result, rule.Children...)

	for _, child := range rule.Children {
		result = append(result, flattenWrappedRules(child)...)
	}
	return result
}

func getPropertyVersionHostnames(ctx context.Context, client papi.PAPI, property *papi.Property, version *papi.GetPropertyVersionsResponse) (*papi.HostnameResponseItems, error) {
	if version == nil {
		var err error
		version, err = client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID:  property.PropertyID,
			ActivatedOn: "",
			ContractID:  property.ContractID,
			GroupID:     property.GroupID,
		})
		if err != nil {
			return nil, err
		}
	}
	response, err := client.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:        property.PropertyID,
		PropertyVersion:   version.Version.PropertyVersion,
		ContractID:        property.ContractID,
		GroupID:           property.GroupID,
		ValidateHostnames: false,
		IncludeCertStatus: false,
	})
	if err != nil {
		return nil, err
	}
	return &response.Hostnames, nil
}

func getEdgeHostnameDetail(ctx context.Context, clientPAPI papi.PAPI, clientHAPI hapi.HAPI, hostnames *papi.HostnameResponseItems, property *papi.Property) (map[string]Hostname, map[string]EdgeHostname, error) {

	edgeHostnamesMap := map[string]EdgeHostname{}
	hostnamesMap := map[string]Hostname{}

	for _, hostname := range hostnames.Items {
		cnameTo := hostname.CnameTo
		cnameFrom := hostname.CnameFrom
		cnameToResource := formatResourceName(cnameTo)

		if hostname.EdgeHostnameID != "" {
			// Get slot details
			edgeHostnameID, err := strconv.Atoi(strings.Replace(hostname.EdgeHostnameID, "ehn_", "", 1))
			if err != nil {
				return nil, nil, fmt.Errorf("invalid Hostname id: %s", err)
			}

			edgeHostname, err := clientHAPI.GetEdgeHostname(ctx, edgeHostnameID)
			if err != nil {
				return nil, nil, fmt.Errorf("edge hostname \"%d\" not found: %s", edgeHostnameID, err)
			}
			papiEdgeHostnames, err := clientPAPI.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
				ContractID: property.ContractID,
				GroupID:    property.GroupID,
				Options:    nil,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("cannot list edge hostnames: %s", err)
			}

			useCases, err := getUseCases(papiEdgeHostnames, hostname.EdgeHostnameID)
			if err != nil {
				return nil, nil, fmt.Errorf("cannot get use cases: %s", err)
			}

			var certificateID int64
			if strings.ToUpper(edgeHostname.SecurityType) == "ENHANCED-TLS" {
				certificate, err := clientHAPI.GetCertificate(ctx, hapi.GetCertificateRequest{
					DNSZone:    edgeHostname.DNSZone,
					RecordName: edgeHostname.RecordName,
				})
				if err != nil {
					if !errors.Is(err, hapi.ErrNotFound) {
						return nil, nil, fmt.Errorf("cannot get certificate details: %s", err)
					}
					certificateID = 0
				} else {
					certificateID, err = strconv.ParseInt(certificate.CertificateID, 10, 64)
					if err != nil {
						return nil, nil, fmt.Errorf("invalid certificate details: %s", err)
					}
				}
			}
			ttl := 0
			if !edgeHostname.UseDefaultTTL {
				ttl = edgeHostname.TTL
			}
			edgeHostnamesMap[cnameToResource] = EdgeHostname{
				EdgeHostname:             cnameTo,
				EdgeHostnameID:           hostname.EdgeHostnameID,
				ContractID:               property.ContractID,
				GroupID:                  property.GroupID,
				TTL:                      ttl,
				IPv6:                     getIPv6(papiEdgeHostnames, hostname.EdgeHostnameID),
				EdgeHostnameResourceName: cnameToResource,
				SecurityType:             edgeHostname.SecurityType,
				UseCases:                 useCases,
				CertificateID:            certificateID,
			}
		}

		certProvisioningType := "CPS_MANAGED"
		if hostname.CertProvisioningType != "" {
			certProvisioningType = hostname.CertProvisioningType
		}
		hostnamesMap[cnameFrom] = Hostname{
			CnameFrom:                cnameFrom,
			CnameTo:                  cnameTo,
			EdgeHostnameResourceName: cnameToResource,
			CertProvisioningType:     certProvisioningType,
			IsActive:                 len(hostname.EdgeHostnameID) > 0,
		}
	}

	return hostnamesMap, edgeHostnamesMap, nil
}

func fetchActiveActivationForNetwork(ctx context.Context, client papi.PAPI, property *papi.Property, network papi.ActivationNetwork) (*papi.Activation, error) {
	activationsResponse, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, err
	}
	return getLatestActiveActivation(activationsResponse.Activations, network), nil
}

// getContactEmails gets list of emails from latest activation
func getContactEmails(activation *papi.Activation) []string {
	if activation == nil || len(activation.NotifyEmails) == 0 {
		return []string{""}
	}
	return activation.NotifyEmails
}

// setPropertyRuleTemplates creates templates based on RuleTemplate and RulesTemplate for given property rule tree response
func setPropertyRuleTemplates(rules *papi.GetRuleTreeResponse) (RuleTemplate, RulesTemplate) {
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
		CustomOverride:      rules.Rules.CustomOverride,
		Children:            make([]string, 0),
		Options:             rules.Rules.Options,
	}

	rulesTemplate := RulesTemplate{
		AccountID:       rules.AccountID,
		ContractID:      rules.ContractID,
		GroupID:         rules.GroupID,
		PropertyID:      rules.PropertyID,
		PropertyVersion: rules.PropertyVersion,
		Etag:            rules.Etag,
		Comments:        rules.Comments,
		RuleFormat:      rules.RuleFormat,
	}

	return ruleTemplate, rulesTemplate
}

// saveSnippets saves given property rules into files under jsonDir directory
func saveSnippets(rules papi.Rules, ruleTemplate RuleTemplate, rulesTemplate RulesTemplate, snippetsPath, templateFileName string) error {
	err := os.MkdirAll(snippetsPath, 0755)
	if err != nil {
		return fmt.Errorf("can't create directory for rule snippets: %s", err)
	}

	nameNormalizer := ruleNameNormalizer()
	for _, rule := range rules.Children {
		jsonBody, err := json.MarshalIndent(rule, "", "  ")
		if err != nil {
			return fmt.Errorf("can't marshall property rule snippets: %s", err)
		}
		name := nameNormalizer(rule.Name)
		rulesNamePath := filepath.Join(snippetsPath, fmt.Sprintf("%s.json", name))
		err = os.WriteFile(rulesNamePath, jsonBody, 0644)
		if err != nil {
			return fmt.Errorf("can't write property rule snippets: %s", err)
		}
		ruleTemplate.Children = append(ruleTemplate.Children, fmt.Sprintf("#include:%s.json", name))
	}

	rulesTemplate.Rule = &ruleTemplate

	jsonBody, err := json.MarshalIndent(rulesTemplate, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshall rule template: %s", err)
	}
	templatePath := filepath.Join(snippetsPath, templateFileName)
	err = os.WriteFile(templatePath, jsonBody, 0644)
	if err != nil {
		return fmt.Errorf("can't write property rule template: %s", err)
	}

	return nil
}

// getUseCases finds UseCases for given edgeHostnameID
func getUseCases(edgeHostnames *papi.GetEdgeHostnamesResponse, edgeHostnameID string) (string, error) {
	for _, edgeHostname := range edgeHostnames.EdgeHostnames.Items {
		if edgeHostname.ID == edgeHostnameID && edgeHostname.UseCases != nil {
			useCasesJSON, err := json.MarshalIndent(edgeHostname.UseCases, "", "  ")
			if err != nil {
				return "", fmt.Errorf("error marshaling Use Cases: %s", err)
			}
			return string(useCasesJSON), nil
		}
	}

	return "", nil
}

// getIPv6 find IPVersionBehavior for given edgeHostnameID
func getIPv6(edgeHostnames *papi.GetEdgeHostnamesResponse, edgeHostnameID string) string {
	for _, edgeHostname := range edgeHostnames.EdgeHostnames.Items {
		if edgeHostname.ID == edgeHostnameID {
			return edgeHostname.IPVersionBehavior
		}
	}
	return ""
}

// findProperty searches for a property with a given name
func findProperty(ctx context.Context, client papi.PAPI, name string) (*papi.Property, error) {
	results, err := client.SearchProperties(ctx, papi.SearchRequest{
		Key:   papi.SearchKeyPropertyName,
		Value: name,
	})
	if err != nil {
		return nil, err
	}

	if results == nil || len(results.Versions.Items) == 0 {
		return nil, fmt.Errorf("unable to find property: \"%s\"", name)
	}

	response, err := client.GetProperty(ctx, papi.GetPropertyRequest{
		PropertyID: results.Versions.Items[0].PropertyID,
		GroupID:    results.Versions.Items[0].GroupID,
		ContractID: results.Versions.Items[0].ContractID,
	})
	if err != nil {
		return nil, err
	}

	return response.Property, nil
}

// getPropertyRules fetches property rules for given property version
func getPropertyRules(ctx context.Context, client papi.PAPI, version *papi.GetPropertyVersionsResponse) (*papi.GetRuleTreeResponse, error) {

	return client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      version.PropertyID,
		PropertyVersion: version.Version.PropertyVersion,
		ContractID:      version.ContractID,
		GroupID:         version.GroupID,
		RuleFormat:      version.Version.RuleFormat,
		ValidateRules:   true,
	})
}

// getVersion gets property version for given property from api
func getVersion(ctx context.Context, client papi.PAPI, property *papi.Property, readVersion string) (*papi.GetPropertyVersionsResponse, *papi.GetPropertyVersionsResponse, error) {
	versions, err := client.GetPropertyVersions(ctx, papi.GetPropertyVersionsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, nil, err
	}

	if readVersion == "LATEST" {
		version, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID:  versions.PropertyID,
			ActivatedOn: "",
			ContractID:  versions.ContractID,
			GroupID:     versions.GroupID,
		})
		if err != nil {
			return nil, nil, err
		}

		return version, version, nil
	}

	v, err := strconv.Atoi(readVersion)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrPropertyVersionNotValid, err)
	}
	// Latest will be the first one
	sort.Slice(versions.Versions.Items, func(i, j int) bool {
		return versions.Versions.Items[i].PropertyVersion > versions.Versions.Items[j].PropertyVersion
	})
	var latestVersion *papi.GetPropertyVersionsResponse
	if len(versions.Versions.Items) == 0 {
		return nil, nil, ErrPropertyVersionNotFound
	}
	latestVersion = getPropertyVersionsResponse(versions, versions.Versions.Items[0])
	for _, item := range versions.Versions.Items {
		if item.PropertyVersion == v {
			return getPropertyVersionsResponse(versions, item), latestVersion, nil
		}
	}
	return nil, nil, ErrPropertyVersionNotFound
}

func getPropertyVersionsResponse(versions *papi.GetPropertyVersionsResponse, item papi.PropertyVersionGetItem) *papi.GetPropertyVersionsResponse {
	return &papi.GetPropertyVersionsResponse{
		PropertyID:   versions.PropertyID,
		PropertyName: versions.PropertyName,
		AccountID:    versions.AccountID,
		ContractID:   versions.ContractID,
		GroupID:      versions.GroupID,
		AssetID:      versions.AssetID,
		Version:      item,
	}
}

// getGroup fetches a group with specific groupID
func getGroup(ctx context.Context, client papi.PAPI, groupID string) (*papi.Group, error) {
	groups, err := client.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	group, err := findGroup(groups.Groups, groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// findGroup finds a specific group by ID
func findGroup(groups papi.GroupItems, id string) (*papi.Group, error) {
	if id == "" {
		return nil, fmt.Errorf("unable to find group: \"%s\"", id)
	}

	for _, group := range groups.Items {
		if group.GroupID == id {
			return group, nil
		}
	}

	return nil, fmt.Errorf("unable to find group: \"%s\"", id)
}

// getProduct finds and returns a productItem with given productID
func getProduct(ctx context.Context, client papi.PAPI, productID string, contractID string) (*papi.ProductItem, error) {
	if contractID == "" {
		return nil, nil
	}

	products, err := client.GetProducts(ctx, papi.GetProductsRequest{
		ContractID: contractID,
	})
	if err != nil {
		return nil, err
	}

	product, err := findProduct(products, productID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// findProduct finds a specific product by ID
func findProduct(products *papi.GetProductsResponse, id string) (*papi.ProductItem, error) {
	for _, product := range products.Products.Items {
		if product.ProductID == id {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("unable to find product: \"%s\"", id)
}

// getLatestActiveActivation retrieves the latest active activation for the specified network.
func getLatestActiveActivation(activationItems papi.ActivationsItems, network papi.ActivationNetwork) *papi.Activation {
	activations := activationItems.Items
	if len(activations) == 0 {
		return nil
	}

	sort.Slice(activations, func(i, j int) bool {
		return activations[i].UpdateDate > activations[j].UpdateDate
	})

	for _, activation := range activations {
		if activation.Status == papi.ActivationStatusActive && activation.Network == network {
			if activation.ActivationType == papi.ActivationTypeActivate {
				return activation
			}
			if activation.ActivationType == papi.ActivationTypeDeactivate {
				return nil
			}
		}
	}

	return nil
}

func ruleNameNormalizer() func(string) string {
	names := map[string]int{}
	return func(name string) string {
		name = normalizeRuleName(name)
		caseInsensitiveName := strings.ToLower(name)
		names[caseInsensitiveName]++
		if count := names[caseInsensitiveName]; count > 1 {
			name = fmt.Sprintf("%s%d", name, count-1)
		}
		return name
	}
}

func normalizeRuleName(name string) string {
	return normalizeRuleNameRegexp.ReplaceAllString(name, "_")
}

// AsInt provides proper conversion of values which are integers in reality
func AsInt(f any) int64 {
	return int64(f.(float64))
}

// as go templates do not support well pointers in receivers and function arguments, global variable seems to be the only
// solution to accumulate all issues
var reportedErrors []string

// ReportError is used to report unknown behaviors or criteria during processing the template
func ReportError(format string, a ...any) string {
	message := fmt.Sprintf(format, a...)
	reportedErrors = append(reportedErrors, message)
	return message
}

// CheckErrors is used to fail the processing of the template in case of any unknown behaviors or criteria
//
// As per template.FuncMap definition, for the error to be treated as error (rather than regular value), it has to be returned as second returned value.
// Therefor always returning "" as first returned value.
func CheckErrors() (string, error) {
	if len(reportedErrors) > 0 {
		return "", fmt.Errorf("there were errors reported: %v", strings.Join(reportedErrors, ", "))
	}
	return "", nil
}

func formatResourceName(name string) string {
	// Replace dots with dashes and spaces with underscores
	formattedName := strings.NewReplacer(".", "-", " ", "_").Replace(name)

	// Prepend underscore if the first character is not a letter
	if !unicode.IsLetter(rune(formattedName[0])) {
		return "_" + formattedName
	}
	return formattedName
}

func prepareRulesForSplitRule(rules *WrappedRules, rulesSplitNestingLeft int, tfWorkPath string, splitRulesWrapper splitDepthRuleWrapper) templates.DataForTarget {
	processedRules := templates.DataForTarget{}

	if rulesSplitNestingLeft == 0 {
		processedRules[filepath.Join(tfWorkPath, "rules", fmt.Sprintf("%s.tf", rules.FileName))] = splitRulesWrapper(flattenRules(rules))
		return processedRules
	}

	processedRules[filepath.Join(tfWorkPath, "rules", fmt.Sprintf("%s.tf", rules.FileName))] = splitRulesWrapper([]*WrappedRules{rules})
	for _, child := range rules.Children {
		processedRules.Join(prepareRulesForSplitRule(child, rulesSplitNestingLeft-1, tfWorkPath, splitRulesWrapper))
	}

	return processedRules
}
