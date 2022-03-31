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

package papi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
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
}

// Hostname represents edge hostname resource
type Hostname struct {
	Hostname                 string
	EdgeHostnameResourceName string
}

// TFData holds template data
type TFData struct {
	GroupName            string
	GroupID              string
	ContractID           string
	PropertyResourceName string
	PropertyName         string
	PropertyID           string
	CPCodeID             string
	ProductID            string
	ProductName          string
	RuleFormat           string
	IsSecure             string
	EdgeHostnames        map[string]EdgeHostname
	Hostnames            map[string]Hostname
	Section              string
	Emails               []string
}

// RulesTemplate represent data used for rules
type RulesTemplate struct {
	AccountID       string        `json:"accountId"`
	ContractID      string        `json:"contractId"`
	GroupID         string        `json:"groupId"`
	PropertyID      string        `json:"propertyId"`
	PropertyVersion int           `json:"propertyVersion"`
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

// CmdCreateProperty is an entrypoint to create-property command
func CmdCreateProperty(c *cli.Context) error {
	ctx := c.Context
	log.SetOutput(ioutil.Discard)
	if c.NArg() != 1 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.Exit(color.RedString("Error displaying help command"), 1)
		}
		return cli.Exit(color.RedString("property name is required"), 1)
	}

	sess := edgegrid.GetSession(c.Context)
	client := papi.Client(sess)
	clientHapi := hapi.Client(sess)
	if c.IsSet("tfworkpath") {
		tools.TFWorkPath = c.String("tfworkpath")
	}
	tools.TFWorkPath = filepath.FromSlash(tools.TFWorkPath)
	if stat, err := os.Stat(tools.TFWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}
	err := tools.CheckFiles(
		tools.CreateTFFilename("property"),
		tools.CreateTFFilename("versions"),
		tools.CreateTFFilename("variables"),
		filepath.Join(tools.TFWorkPath, "rules.json"),
		filepath.Join(tools.TFWorkPath, "import.sh"))
	if err != nil {
		return cli.Exit(err, 1)
	}

	var tfData TFData
	tfData.EdgeHostnames = make(map[string]EdgeHostname)
	tfData.Hostnames = make(map[string]Hostname)
	tfData.Emails = make([]string, 0)

	tfData.Section = c.String("section")

	// Get Property
	propertyName := c.Args().First()
	term := terminal.Get(ctx)
	fmt.Println("Configuring Property")
	term.Spinner().Start("Fetching property ")
	property := findProperty(ctx, client, propertyName)
	if property == nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Property not found "), 1)
	}

	tfData.ContractID = property.ContractID
	tfData.PropertyName = property.PropertyName
	tfData.PropertyID = property.PropertyID
	tfData.PropertyResourceName = strings.Replace(property.PropertyName, ".", "-", -1)
	term.Spinner().OK()

	// Get Property Rules
	term.Spinner().Start("Fetching property rules ")
	rules, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
		RuleFormat:      property.RuleFormat,
	})

	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Property rules not found: ", err), 1)
	}

	tfData.IsSecure = "false"
	if rules.Rules.Options.IsSecure {
		tfData.IsSecure = "true"
	}

	for _, behaviour := range rules.Rules.Behaviors {
		_ = behaviour
		if behaviour.Name == "cpCode" {
			options := behaviour.Options
			value := options["value"].(map[string]interface{})
			tfData.CPCodeID = fmt.Sprintf("%.0f", value["id"].(float64))
		}
	}

	// Get Rule Format
	tfData.RuleFormat = rules.RuleFormat

	// Set up template structure
	var ruletemplate RuleTemplate
	ruletemplate.Name = rules.Rules.Name
	ruletemplate.Criteria = rules.Rules.Criteria
	ruletemplate.Behaviors = rules.Rules.Behaviors
	ruletemplate.Comments = rules.Rules.Comments
	ruletemplate.CriteriaLocked = rules.Rules.CriteriaLocked
	ruletemplate.CriteriaMustSatisfy = rules.Rules.CriteriaMustSatisfy
	ruletemplate.UUID = rules.Rules.UUID
	ruletemplate.Variables = rules.Rules.Variables
	ruletemplate.AdvancedOverride = rules.Rules.AdvancedOverride
	ruletemplate.Children = make([]string, 0)
	ruletemplate.Options = rules.Rules.Options

	var rulestemplate RulesTemplate
	rulestemplate.AccountID = rules.AccountID
	rulestemplate.ContractID = rules.ContractID
	rulestemplate.GroupID = rules.GroupID
	rulestemplate.PropertyID = rules.PropertyID
	rulestemplate.PropertyVersion = rules.PropertyVersion
	rulestemplate.Etag = rules.Etag
	rulestemplate.RuleFormat = rules.RuleFormat
	rulestemplate.Rule = &ruletemplate

	// Save snippets
	snippetspath := filepath.Join(tools.TFWorkPath, "property-snippets")
	os.Mkdir(snippetspath, 0755)

	for _, rule := range rules.Rules.Children {
		jsonBody, err := json.MarshalIndent(rule, "", "  ")
		if err != nil {
			term.Spinner().Fail()
			return cli.Exit(color.RedString("Can't marshall property rule snippets: ", err), 1)
		}
		name := strings.ReplaceAll(rule.Name, " ", "_")
		rulesnamepath := filepath.Join(snippetspath, fmt.Sprintf("%s.json", name))
		err = ioutil.WriteFile(rulesnamepath, jsonBody, 0644)
		if err != nil {
			term.Spinner().Fail()
			return cli.Exit(color.RedString("Can't write property rule snippets: ", err), 1)
		}
		ruletemplate.Children = append(ruletemplate.Children, fmt.Sprintf("#include:%s.json", name))
	}

	jsonBody, err := json.MarshalIndent(rulestemplate, "", "  ")
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Can't marshall rule template: ", err), 1)
	}
	templatepath := filepath.Join(snippetspath, "main.json")
	err = ioutil.WriteFile(templatepath, jsonBody, 0644)
	if err != nil {
		term.Spinner().Fail()

		return cli.Exit(color.RedString("Can't write property rule template: ", err), 1)
	}

	term.Spinner().OK()

	// Get Group
	term.Spinner().Start("Fetching group ")
	group, err := getGroup(ctx, client, property.GroupID)
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Group not found: %s", err), 1)
	}

	tfData.GroupName = group.GroupName
	tfData.GroupID = group.GroupID

	term.Spinner().OK()

	// Get Version
	term.Spinner().Start("Fetching property version ")
	version, err := getVersion(ctx, client, property)
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Version not found: %s", err), 1)
	}

	tfData.ProductID = version.Version.ProductID

	term.Spinner().OK()

	// Get Product
	term.Spinner().Start("Fetching product name ")
	product, err := getProduct(ctx, client, tfData.ProductID, property.ContractID)
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Product not found: %s", err), 1)
	}

	tfData.ProductName = product.ProductName

	term.Spinner().OK()

	// Get Hostnames
	term.Spinner().Start("Fetching hostnames ")
	hostnames, err := getHostnames(ctx, client, property, version)

	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Hostnames not found: %s", err), 1)
	}

	for _, hostname := range hostnames.Items {
		_ = hostname

		if hostname.EdgeHostnameID == "" {
			continue
		}

		// Get slot details
		ehnid, err := strconv.Atoi(strings.Replace(hostname.EdgeHostnameID, "ehn_", "", 1))
		if err != nil {
			term.Spinner().Fail()
			return cli.Exit(color.RedString("Invalid Hostname id: %s", err), 1)
		}

		edgehostname, err := clientHapi.GetEdgeHostname(ctx, ehnid)
		if err != nil {
			term.Spinner().Fail()
			return cli.Exit(color.RedString("Edge Hostname not found: %s", err), 1)
		}

		cnameTo := hostname.CnameTo
		cnameFrom := hostname.CnameFrom
		cnameToResource := strings.Replace(cnameTo, ".", "-", -1)

		var edgeHostnameN EdgeHostname
		edgeHostnameN.EdgeHostname = cnameTo
		edgeHostnameN.EdgeHostnameID = hostname.EdgeHostnameID
		edgeHostnameN.EdgeHostnameResourceName = cnameToResource
		edgeHostnameN.ProductName = product.ProductName
		edgeHostnameN.IPv6 = getIPv6(ctx, client, property, hostname.EdgeHostnameID)
		edgeHostnameN.SlotNumber = edgehostname.SlotNumber
		edgeHostnameN.SecurityType = edgehostname.SecurityType
		edgeHostnameN.ContractID = property.ContractID
		edgeHostnameN.GroupID = group.GroupID
		tfData.EdgeHostnames[cnameToResource] = edgeHostnameN

		var hostnamesN Hostname
		hostnamesN.Hostname = cnameFrom
		hostnamesN.EdgeHostnameResourceName = cnameToResource
		tfData.Hostnames[cnameFrom] = hostnamesN

	}

	term.Spinner().OK()

	// Get contact details
	term.Spinner().Start("Fetching contact details ")
	activationsResponse, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		tfData.Emails = append(tfData.Emails, "")
	} else {
		a, err := getLatestStagingActivation(activationsResponse.Activations, "")
		if err != nil {
			tfData.Emails = append(tfData.Emails, "")
		} else {
			tfData.Emails = a.NotifyEmails
		}
	}

	term.Spinner().OK()

	// Save file
	term.Spinner().Start("Saving TF configurations ")
	err = saveTerraformDefinition(tfData)
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Couldn't save tf file: %s", err), 1)
	}

	term.Spinner().OK()
	fmt.Println("Property configuration completed")

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

func getIPv6(ctx context.Context, client papi.PAPI, property *papi.Property, ehn string) string {
	edgeHostnames, err := client.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
		Options:    nil,
	})
	if err != nil {
		return "false"
	}
	for _, edgehostname := range edgeHostnames.EdgeHostnames.Items {
		_ = edgehostname
		if edgehostname.ID == ehn {
			return edgehostname.IPVersionBehavior
		}
	}
	return ""
}

func findProperty(ctx context.Context, client papi.PAPI, name string) *papi.Property {
	results, err := client.SearchProperties(ctx, papi.SearchRequest{
		Key:   papi.SearchKeyPropertyName,
		Value: name,
	})
	if err != nil {
		return nil
	}

	if err != nil || results == nil || len(results.Versions.Items) == 0 {
		return nil
	}

	response, err := client.GetProperty(ctx, papi.GetPropertyRequest{
		PropertyID: results.Versions.Items[0].PropertyID,
		GroupID:    results.Versions.Items[0].GroupID,
		ContractID: results.Versions.Items[0].ContractID,
	})
	if err != nil {
		return nil
	}

	return response.Property
}

func getVersion(ctx context.Context, client papi.PAPI, property *papi.Property) (*papi.GetPropertyVersionsResponse, error) {
	versions, err := client.GetPropertyVersions(ctx, papi.GetPropertyVersionsRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, err
	}

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

func getGroup(ctx context.Context, client papi.PAPI, groupID string) (*papi.Group, error) {
	groups, err := client.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	group, e := findGroup(groups.Groups, groupID)
	if e != nil {
		return nil, e
	}

	return group, nil
}

// FindGroup finds a specific group by ID
func findGroup(groups papi.GroupItems, id string) (*papi.Group, error) {
	var group *papi.Group
	var groupFound bool

	if id == "" {
		goto err
	}

	for _, group = range groups.Items {
		if group.GroupID == id {
			groupFound = true
			break
		}
	}

err:
	if !groupFound {
		return nil, fmt.Errorf("unable to find group: \"%s\"", id)
	}

	return group, nil
}

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

	product, e := findProduct(products, productID)
	if e != nil {
		return nil, e
	}

	return product, nil
}

// findProduct finds a specific product by ID
func findProduct(products *papi.GetProductsResponse, id string) (*papi.ProductItem, error) {
	var product papi.ProductItem
	var productFound bool
	for _, product = range products.Products.Items {
		if product.ProductID == id {
			productFound = true
			break
		}
	}

	if !productFound {
		return nil, fmt.Errorf("unable to find product: \"%s\"", id)
	}

	return &product, nil
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

func saveTerraformDefinition(data TFData) error {

	tftemplate, err := template.New("tf").Parse(
		"provider \"akamai\" {\n" +
			" edgerc = \"~/.edgerc\"\n" +
			" config_section = \"{{.Section}}\"\n" +
			"}\n" +
			"\n" +
			"data \"akamai_group\" \"group\" {\n" +
			" group_name = \"{{.GroupName}}\"\n" +
			" contract_id = \"{{.ContractID}}\"\n" +
			"}\n" +
			"\n" +
			"data \"akamai_contract\" \"contract\" {\n" +
			"  group_name = data.akamai_group.group.name\n" +
			"}\n" +
			"\n" +
			"data \"akamai_property_rules_template\" \"rules\" {\n" +
			"  template_file = abspath(\"${path.module}/property-snippets/main.json\")\n" +
			"}\n" +
			"\n" +

			// Edge hostname loop
			"{{range .EdgeHostnames}}" +
			"resource \"akamai_edge_hostname\" \"{{.EdgeHostnameResourceName}}\" {\n" +
			" product_id  = \"prd_{{.ProductName}}\"\n" +
			" contract_id = data.akamai_contract.contract.id\n" +
			" group_id = data.akamai_group.group.id\n" +
			" ip_behavior = \"{{.IPv6}}\"\n" +
			" edge_hostname = \"{{.EdgeHostname}}\"\n" +
			"{{if .SlotNumber}}" +
			" certificate = {{.SlotNumber}}\n" +
			"{{end}}" +
			"}\n" +
			"\n" +
			"{{end}}" +
			"resource \"akamai_property\" \"{{.PropertyResourceName}}\" {\n" +
			" name = \"{{.PropertyName}}\"\n" +
			" contract_id = data.akamai_contract.contract.id\n" +
			" group_id = data.akamai_group.group.id\n" +
			" product_id = \"prd_{{.ProductName}}\"\n" +
			" rule_format = \"{{.RuleFormat}}\"\n" +
			"{{range .Hostnames}}" +
			" hostnames {\n" +
			"  cname_from = \"{{.Hostname}}\"\n" +
			"  cname_to = akamai_edge_hostname.{{.EdgeHostnameResourceName}}.edge_hostname\n" +
			"  cert_provisioning_type = \"CPS_MANAGED\"\n" +
			" }\n" +
			"{{end}}" +
			" rules = data.akamai_property_rules_template.rules.json\n" +
			"}\n" +
			"\n" +
			"resource \"akamai_property_activation\" \"{{.PropertyResourceName}}\" {\n" +
			" property_id = akamai_property.{{.PropertyResourceName}}.id\n" +
			" contact = [\"{{range $index, $element := .Emails}}{{if $index}},{{end}}{{$element}}{{end}}\"]\n" +
			" version = akamai_property.{{.PropertyResourceName}}.latest_version\n" +
			" network = upper(var.env)\n" +
			"}\n")

	if err != nil {
		return err
	}
	propertyconfigfilename := tools.CreateTFFilename("property")
	f, err := os.Create(propertyconfigfilename)
	if err != nil {
		return err
	}

	err = tftemplate.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// version
	versionsconfigfilename := tools.CreateTFFilename("versions")
	f, err = os.Create(versionsconfigfilename)
	if err != nil {
		return err
	}

	f.WriteString("terraform {\n")
	f.WriteString("required_providers {\n")
	f.WriteString("akamai = {\n")
	f.WriteString("source = \"akamai/akamai\"\n")
	f.WriteString("}\n")
	f.WriteString("}\n")
	f.WriteString("required_version = \">= 0.13\"\n")
	f.WriteString("}\n")
	f.Close()

	// variables
	variablesconfigfilename := tools.CreateTFFilename("variables")
	f, err = os.Create(variablesconfigfilename)
	if err != nil {
		return err
	}

	f.WriteString("variable \"env\" {\n")
	f.WriteString(" default = \"staging\"\n")
	f.WriteString("}\n")
	f.Close()

	// import script
	importfilename := filepath.Join(tools.TFWorkPath, "import.sh")
	f, err = os.Create(importfilename)
	if err != nil {
		return err
	}

	importtemplate, err := template.New("import").Parse(
		"terraform init\n" +
			"{{range .EdgeHostnames}}" +
			"terraform import akamai_edge_hostname.{{.EdgeHostnameResourceName}} {{.EdgeHostnameID}},{{.ContractID}},{{.GroupID}}\n" +
			"{{end}}" +
			"terraform import akamai_property.{{.PropertyResourceName}} {{.PropertyID}},{{.ContractID}},{{.GroupID}}\n")
	if err != nil {
		return err
	}

	err = importtemplate.Execute(f, data)
	if err != nil {
		return err
	}

	return f.Close()
}
