/*
 * Copyright 2018-2020. Akamai Technologies, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package papi contains code for exporting properties
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
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
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
	Children      []*WrappedRules
}

// TFData holds template data
type TFData struct {
	Includes      []TFIncludeData
	Property      TFPropertyData
	Section       string
	Rules         []*WrappedRules
	RulesAsSchema bool
}

// TFIncludeData holds template data for include
type TFIncludeData struct {
	ActivationNoteProduction   string
	ActivationNoteStaging      string
	ContractID                 string
	ActivationEmailsProduction []string
	ActivationEmailsStaging    []string
	GroupID                    string
	IncludeID                  string
	IncludeName                string
	IncludeType                string
	Networks                   []string
	RuleFormat                 string
	VersionProduction          string
	VersionStaging             string
	Rules                      []*WrappedRules
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
	ReadVersion          string
	ProductionInfo       NetworkInfo
	StagingInfo          NetworkInfo
}

// NetworkInfo holds details for specific network
type NetworkInfo struct {
	Emails         []string
	ActivationNote string
	HasActivation  bool
	Version        int
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
	// ErrFetchingReferencedIncludes is returned when fetching referenced includes request failed
	ErrFetchingReferencedIncludes = errors.New("fetching referenced includes")
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
)

var additionalFuncs = template.FuncMap{
	"ToLower":           strings.ToLower,
	"TerraformName":     TerraformName,
	"AsInt":             AsInt,
	"Escape":            tools.Escape,
	"ReportError":       ReportError,
	"CheckErrors":       CheckErrors,
	"IsMultiline":       IsMultiline,
	"NoNewlineAtTheEnd": NoNewlineAtTheEnd,
	"RemoveLastNewline": RemoveLastNewline,
	"GetEOT":            GetEOT,
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

	var withIncludes bool
	if c.IsSet("with-includes") {
		withIncludes = c.Bool("with-includes")
		if withIncludes {
			templateToFile["includes.tmpl"] = filepath.Join(tfWorkPath, "includes.tf")
		}
	}

	var schema bool
	if c.IsSet("schema") {
		schema = c.Bool("schema")
	}

	if withIncludes && schema {
		templateToFile["includes_rules.tmpl"] = filepath.Join(tfWorkPath, "includes_rules.tf")
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	propertyName := c.Args().First()
	section := edgegrid.GetEdgercSection(c)
	if err = createProperty(ctx, propertyName, version, section, "property-snippets", tfWorkPath, withIncludes, schema, client, clientHapi, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting property: %s", err)), 1)
	}
	return nil
}

func createProperty(ctx context.Context, propertyName, readVersion, section, jsonDir, tfWorkPath string, withIncludes, schema bool, client papi.PAPI, clientHapi hapi.HAPI, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	tfData := TFData{
		Property: TFPropertyData{
			EdgeHostnames: make(map[string]EdgeHostname),
		},
		Section:       section,
		RulesAsSchema: schema,
	}

	// Get Property
	term.Spinner().Start("Fetching property " + propertyName)
	property, err := findProperty(ctx, client, propertyName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrPropertyNotFound, err)
	}

	tfData.Property.ContractID = property.ContractID
	tfData.Property.PropertyName = property.PropertyName
	tfData.Property.PropertyID = property.PropertyID
	tfData.Property.PropertyResourceName = strings.Replace(property.PropertyName, ".", "-", -1)

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

	if readVersion == "" {
		readVersion = "LATEST"
	}

	// Get Version
	term.Spinner().Start("Fetching property version ")
	version, err := getVersion(ctx, client, property, readVersion)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrPropertyVersionNotFound, err)
	}

	tfData.Property.ProductID = version.Version.ProductID
	tfData.Property.ReadVersion = readVersion

	term.Spinner().OK()

	// Get Includes if withIncludes is set
	if withIncludes {
		term.Spinner().Start("Fetching referenced includes with property " + propertyName)
		includes, err := client.ListReferencedIncludes(ctx, papi.ListReferencedIncludesRequest{
			PropertyID:      property.PropertyID,
			ContractID:      property.ContractID,
			GroupID:         property.GroupID,
			PropertyVersion: version.Version.PropertyVersion,
		})
		if err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrFetchingReferencedIncludes, err)
		}
		term.Spinner().OK()

		tfData.Includes = make([]TFIncludeData, 0)
		for _, include := range includes.Includes.Items {
			includeData, rules, err := getIncludeData(ctx, &include, client)
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
		}
	}

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
	}

	term.Spinner().OK()

	filterFuncs := make([]func([]string) ([]string, error), 0)
	if schema {
		ruleTemplate := fmt.Sprintf("rules_%s.tmpl", rules.RuleFormat)
		if !templateProcessor.TemplateExists(ruleTemplate) {
			return fmt.Errorf("%w: %s", ErrUnsupportedRuleFormat, rules.RuleFormat)
		}
		templateProcessor.AddTemplateTarget(ruleTemplate, filepath.Join(tfWorkPath, "rules.tf"))
		tfData.Rules = flattenRules(tfData.Property.PropertyName, rules.Rules)
		filterFuncs = append(filterFuncs, useThisOnlyRuleFormat(rules.RuleFormat))
	}
	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData, filterFuncs...); err != nil {
		term.Spinner().Fail()
		if _, err := CheckErrors(); err != nil {
			return fmt.Errorf("%w", err)
		}
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	if !schema {
		// Save snippets
		ruleTemplate, rulesTemplate := setPropertyRuleTemplates(rules)
		if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(tfWorkPath, jsonDir), "main.json"); err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
		}
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for property '%s' was saved successfully\n", property.PropertyName)

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

func flattenRules(property string, rule papi.Rules) []*WrappedRules {
	var result []*WrappedRules
	wrappedRules := wrapRules(rule)
	result = append(result, wrappedRules)
	result = append(result, flattenWrappedRules(wrappedRules)...)
	var names = map[string]int{}
	for _, wrappedRules := range result {
		name := TerraformName(wrappedRules.Rule.Name)
		names[name]++
		if count := names[name]; count > 1 {
			name = fmt.Sprintf("%s%d", name, count-1)
		}
		wrappedRules.TerraformName = fmt.Sprintf("%s_rule_%s", TerraformName(property), name)
	}
	return result
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
		cnameToResource := strings.Replace(cnameTo, ".", "-", -1)

		if hostname.EdgeHostnameID != "" {
			// Get slot details
			edgeHostnameID, err := strconv.Atoi(strings.Replace(hostname.EdgeHostnameID, "ehn_", "", 1))
			if err != nil {
				return nil, nil, fmt.Errorf("invalid Hostname id: %s", err)
			}

			edgeHostname, err := clientHAPI.GetEdgeHostname(ctx, edgeHostnameID)
			if err != nil {
				return nil, nil, fmt.Errorf("edge hostname %d not found: %s", edgeHostnameID, err)
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

			edgeHostnamesMap[cnameToResource] = EdgeHostname{
				EdgeHostname:             cnameTo,
				EdgeHostnameID:           hostname.EdgeHostnameID,
				ContractID:               property.ContractID,
				GroupID:                  property.GroupID,
				IPv6:                     getIPv6(papiEdgeHostnames, hostname.EdgeHostnameID),
				EdgeHostnameResourceName: cnameToResource,
				SecurityType:             edgeHostname.SecurityType,
				UseCases:                 useCases,
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
func getVersion(ctx context.Context, client papi.PAPI, property *papi.Property, readVersion string) (*papi.GetPropertyVersionsResponse, error) {
	versions, err := client.GetPropertyVersions(ctx, papi.GetPropertyVersionsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, err
	}

	if readVersion == "LATEST" {
		version, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID:  versions.PropertyID,
			ActivatedOn: "",
			ContractID:  versions.ContractID,
			GroupID:     versions.GroupID,
		})
		if err != nil {
			return nil, err
		}

		return version, nil
	}

	v, err := strconv.Atoi(readVersion)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrPropertyVersionNotValid, err)
	}
	for _, item := range versions.Versions.Items {
		if item.PropertyVersion == v {
			return &papi.GetPropertyVersionsResponse{
				PropertyID:   versions.PropertyID,
				PropertyName: versions.PropertyName,
				AccountID:    versions.AccountID,
				ContractID:   versions.ContractID,
				GroupID:      versions.GroupID,
				AssetID:      versions.AssetID,
				Version:      item,
			}, nil
		}
	}

	return nil, ErrPropertyVersionNotFound
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
		names[name]++
		if count := names[name]; count > 1 {
			name = fmt.Sprintf("%s%d", name, count-1)
		}
		return name
	}
}

func normalizeRuleName(name string) string {
	return normalizeRuleNameRegexp.ReplaceAllString(name, "_")
}

var matchFirstCap = regexp.MustCompile("([^ _])([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase returns name using snake case notation - SomeName -> some_name
func ToSnakeCase(str string) string {
	snake := strings.Replace(str, " ", "_", -1)
	snake = matchFirstCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var nameRegexp = regexp.MustCompile(`[^\p{L}\p{Nl}\p{Mn}\p{Mc}\p{Nd}\p{Pc}\d\-_ ]`)

// TerraformName is used to convert rule name into valid name of the exported data source
// Current implementation is not covering all the cases defined in the terraform specification
// https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#identifiers and http://unicode.org/reports/tr31/ ,
// but only a reasonable subset.
func TerraformName(str string) string {
	str = nameRegexp.ReplaceAllString(str, "-")
	return ToSnakeCase(str)
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
func CheckErrors() (string, error) {
	if len(reportedErrors) > 0 {
		return "", fmt.Errorf("there were errors reported: %v", strings.Join(reportedErrors, ", "))
	}
	return "", nil
}

// IsMultiline returns true if the input string contains at least one new line character
func IsMultiline(str string) bool {
	return strings.LastIndex(str, "\n") >= 0
}

// NoNewlineAtTheEnd returns true if there is no new line character at the end of the string
func NoNewlineAtTheEnd(str string) bool {
	if str == "" {
		return true
	}
	return str[len(str)-1:] != "\n"
}

// RemoveLastNewline removes the new line character if this is the last character in the string
func RemoveLastNewline(str string) string {
	if len(str) > 0 && str[len(str)-1:] == "\n" {
		return str[:len(str)-1]
	}
	return str
}

// GetEOT generates unique delimiter word for heredoc, by default it is EOT
func GetEOT(str string) string {
	eot := "EOT"
	for strings.LastIndex(str, eot) >= 0 {
		eot += "A"
	}
	return eot
}
