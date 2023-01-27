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
	"strconv"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
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
	ProductName              string
	ContractID               string
	GroupID                  string
	ID                       string
	IPv6                     string
	EdgeHostnameResourceName string
	SlotNumber               int
	SecurityType             string
	UseCases                 string
}

// Hostname represents edge hostname resource
type Hostname struct {
	Hostname                 string
	EdgeHostnameResourceName string
	CertProvisioningType     string
}

// TFData holds template data
type TFData struct {
	Includes []TFIncludeData
	Property TFPropertyData
	Section  string
}

// TFIncludeData holds template data for an include
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
	Emails               []string
	ActivationNote       string
	Version              string
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
)

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

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: template.FuncMap{
			"ToLower": strings.ToLower,
		},
	}

	propertyName := c.Args().First()
	section := edgegrid.GetEdgercSection(c)
	if err = createProperty(ctx, propertyName, version, section, "property-snippets", tfWorkPath, withIncludes, client, clientHapi, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting property: %s", err)), 1)
	}
	return nil
}

func createProperty(ctx context.Context, propertyName, readVersion, section, jsonDir, tfWorkPath string, withIncludes bool, client papi.PAPI, clientHapi hapi.HAPI, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	tfData := TFData{
		Property: TFPropertyData{
			EdgeHostnames: make(map[string]EdgeHostname),
			Emails:        make([]string, 0),
		},
		Section: section,
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
	tfData.Property.Version = readVersion

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
			includeData, err := getIncludeData(ctx, &include, jsonDir, tfWorkPath, client)
			if err != nil {
				return err
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
	hostnames, err := getHostnames(ctx, client, property, version)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrHostnamesNotFound, err)
	}

	tfData.Property.Hostnames, tfData.Property.EdgeHostnames, err =
		getEdgeHostnameDetail(ctx, client, clientHapi, hostnames, product.ProductName, property)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingHostnameDetails, err)
	}

	term.Spinner().OK()

	term.Spinner().Start("Fetching activation details ")
	latestActivation, err := fetchLatestActivation(ctx, client, property)
	if err == nil {
		tfData.Property.ActivationNote = latestActivation.Note
		tfData.Property.Emails = getContactEmails(latestActivation)
	}
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	// Save snippets
	ruleTemplate, rulesTemplate := setPropertyRuleTemplates(rules)
	if err = saveSnippets(rules.Rules, ruleTemplate, rulesTemplate, filepath.Join(tfWorkPath, jsonDir), "main.json"); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingSnippets, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for property '%s' was saved successfully\n", property.PropertyName)

	return nil
}

func getHostnames(ctx context.Context, client papi.PAPI, property *papi.Property, version *papi.GetPropertyVersionsResponse) (*papi.HostnameResponseItems, error) {
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

func getEdgeHostnameDetail(ctx context.Context, clientPAPI papi.PAPI, clientHAPI hapi.HAPI, hostnames *papi.HostnameResponseItems,
	productName string, property *papi.Property) (map[string]Hostname, map[string]EdgeHostname, error) {

	edgeHostnamesMap := map[string]EdgeHostname{}
	hostnamesMap := map[string]Hostname{}

	for _, hostname := range hostnames.Items {
		if hostname.EdgeHostnameID == "" {
			continue
		}

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

		cnameTo := hostname.CnameTo
		cnameFrom := hostname.CnameFrom
		cnameToResource := strings.Replace(cnameTo, ".", "-", -1)

		edgeHostnamesMap[cnameToResource] = EdgeHostname{
			EdgeHostname:             cnameTo,
			EdgeHostnameID:           hostname.EdgeHostnameID,
			ProductName:              productName,
			ContractID:               property.ContractID,
			GroupID:                  property.GroupID,
			IPv6:                     getIPv6(papiEdgeHostnames, hostname.EdgeHostnameID),
			EdgeHostnameResourceName: cnameToResource,
			SlotNumber:               edgeHostname.SlotNumber,
			SecurityType:             edgeHostname.SecurityType,
			UseCases:                 useCases,
		}

		certProvisioningType := "CPS_MANAGED"
		if hostname.CertProvisioningType != "" {
			certProvisioningType = hostname.CertProvisioningType
		}
		hostnamesMap[cnameFrom] = Hostname{
			Hostname:                 cnameFrom,
			EdgeHostnameResourceName: cnameToResource,
			CertProvisioningType:     certProvisioningType,
		}
	}

	return hostnamesMap, edgeHostnamesMap, nil
}

func fetchLatestActivation(ctx context.Context, client papi.PAPI, property *papi.Property) (*papi.Activation, error) {
	activationsResponse, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, err
	}

	latestActivation, err := getLatestStagingActivation(activationsResponse.Activations, "")
	if err != nil {
		return nil, err
	}

	return latestActivation, nil
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

// getLatestStagingActivation retrieves the latest activation for the staging network
//
// Pass in a status to check for, defaults to StatusActive
func getLatestStagingActivation(activations papi.ActivationsItems, status papi.ActivationStatus) (*papi.Activation, error) {
	return getLatestActivation(activations, papi.ActivationNetworkStaging, status)
}

// getLatestActivation gets the latest activation for the specified network
//
// Defaults to NetworkProduction. Pass in a status to check for, defaults to StatusActive
//
// This can return an activation OR a deactivation. Check activation.ActivationType and activation.Status for what you're looking for
func getLatestActivation(activations papi.ActivationsItems, network papi.ActivationNetwork, status papi.ActivationStatus) (*papi.Activation, error) {
	if network == "" {
		network = papi.ActivationNetworkProduction
	}

	if status == "" {
		status = papi.ActivationStatusActive
	}

	var latest *papi.Activation
	for _, activation := range activations.Items {
		if activation.Network == network && activation.Status == status && (latest == nil || activation.PropertyVersion > latest.PropertyVersion) {
			latest = activation
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no activation found (network: %s, status: %s)", network, status)
	}

	return latest, nil
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
