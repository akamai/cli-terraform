package appsec

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templateFiles embed.FS
var client appsec.APPSEC

var (
	// ErrFetchingPolicy is returned when fetching policy fails
	ErrFetchingPolicy = errors.New("unable to fetch policy with given name")
	// ErrFetchingVersion is returned when fetching policy version fails
	ErrFetchingVersion = errors.New("unable to fetch latest policy version")
	// ErrCloudletTypeNotSupported is returned when a provided cloudlet type is not yet supported
	ErrCloudletTypeNotSupported = errors.New("cloudlet type not supported")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving terraform project files")

	section string
)

// CmdCreateAppsec is an entrypoint to create-appsec command
func CmdCreateAppsec(c *cli.Context) error {
	ctx := c.Context
	if c.NArg() == 0 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.NewExitError(color.RedString("Error displaying help command"), 1)
		}
		return cli.NewExitError(color.RedString("Appsec configuration name is required"), 1)
	}

	sess := edgegrid.GetSession(ctx)
	client = appsec.Client(sess)

	tfWorkPath := "." // default is current directory

	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)
	if stat, err := os.Stat(tfWorkPath); err != nil || !stat.IsDir() {
		return cli.NewExitError(color.RedString("Destination work path is not accessible"), 1)
	}

	// Directory Paths
	modulesPath := filepath.Join(tfWorkPath, "modules")
	securityModulePath := filepath.Join(modulesPath, "security")
	activateSecurityModulePath := filepath.Join(modulesPath, "activate-security")
	paths := []string{modulesPath, securityModulePath, activateSecurityModulePath}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return cli.NewExitError(color.RedString(err.Error()), 1)
		}
	}

	// File Paths
	appsecPath := filepath.Join(tfWorkPath, "appsec.tf")

	err := tools.CheckFiles(appsecPath)
	if err != nil {
		return cli.NewExitError(color.RedString(err.Error()), 1)
	}

	// Save our section for use later
	section = edgegrid.GetEdgercSection(c)

	// Template to path mappings
	templateToFile := map[string]string{
		"appsec.tmpl":                         appsecPath,
		"imports.tmpl":                        filepath.Join(tfWorkPath, "appsec-import.sh"),
		"main.tmpl":                           filepath.Join(tfWorkPath, "appsec-main.tf"),
		"modules-activate-security-main.tmpl": filepath.Join(activateSecurityModulePath, "main.tf"),
		"modules-activate-security-variables.tmpl":  filepath.Join(activateSecurityModulePath, "variables.tf"),
		"modules-activate-security-versions.tmpl":   filepath.Join(activateSecurityModulePath, "versions.tf"),
		"modules-security-advanced.tmpl":            filepath.Join(securityModulePath, "advanced.tf"),
		"modules-security-api.tmpl":                 filepath.Join(securityModulePath, "api.tf"),
		"modules-security-custom-deny.tmpl":         filepath.Join(securityModulePath, "custom-deny.tf"),
		"modules-security-custom-rules.tmpl":        filepath.Join(securityModulePath, "custom-rules.tf"),
		"modules-security-firewall.tmpl":            filepath.Join(securityModulePath, "firewall.tf"),
		"modules-security-main.tmpl":                filepath.Join(securityModulePath, "main.tf"),
		"modules-security-match-targets.tmpl":       filepath.Join(securityModulePath, "match-targets.tf"),
		"modules-security-penalty-box.tmpl":         filepath.Join(securityModulePath, "penalty-box.tf"),
		"modules-security-policies.tmpl":            filepath.Join(securityModulePath, "policies.tf"),
		"modules-security-protections.tmpl":         filepath.Join(securityModulePath, "protections.tf"),
		"modules-security-rate-policies.tmpl":       filepath.Join(securityModulePath, "rate-policies.tf"),
		"modules-security-rate-policy-actions.tmpl": filepath.Join(securityModulePath, "rate-policy-actions.tf"),
		"modules-security-reputation-profiles.tmpl": filepath.Join(securityModulePath, "reputation-profiles.tf"),
		"modules-security-reputation.tmpl":          filepath.Join(securityModulePath, "reputation.tf"),
		"modules-security-siem.tmpl":                filepath.Join(securityModulePath, "siem.tf"),
		"modules-security-slow-post.tmpl":           filepath.Join(securityModulePath, "slow-post.tf"),
		"modules-security-variables.tmpl":           filepath.Join(securityModulePath, "variables.tf"),
		"modules-security-versions.tmpl":            filepath.Join(securityModulePath, "versions.tf"),
		"modules-security-waf.tmpl":                 filepath.Join(securityModulePath, "waf.tf"),
		"variables.tmpl":                            filepath.Join(tfWorkPath, "appsec-variables.tf"),
		"versions.tmpl":                             filepath.Join(tfWorkPath, "appsec-versions.tf"),
	}

	// Provide custom helper functions to get data that does not exist in the security config export
	additionalFuncs := template.FuncMap{
		"exportJSON":            exportJSON,
		"getConfigDescription":  getConfigDescription,
		"getCustomRuleNameByID": getCustomRuleNameByID,
		"getPolicyNameByID":     getPolicyNameByID,
		"getPrefixFromID":       getPrefixFromID,
		"getRateNameByID":       getRateNameByID,
		"getRepNameByID":        getRepNameByID,
		"getRuleDescByID":       getRuleDescByID,
		"getRuleNameByID":       getRuleNameByID,
		"getSection":            getSection,
		"getWAFMode":            getWAFMode,
		"isStructuredRule":      isStructuredRule,
	}

	// The template processor
	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	appsecName := c.Args().First()
	if err = createAppsec(ctx, appsecName, client, processor); err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Error exporting appsec config HCL: %s", err)), 1)
	}
	return nil
}

func createAppsec(ctx context.Context, configName string, client appsec.APPSEC, templateProcessor templates.TemplateProcessor) error {

	term := terminal.Get(ctx)

	fmt.Println("Configuring Appsec")
	term.Spinner().Start("Finding appsec configuration " + configName)

	id, version, err := findConfigurationIDByName(ctx, configName, client)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicy, err)
	}

	term.Spinner().OK()

	term.Spinner().Start("Fetching appsec configuration " + configName)

	configuration, err := exportConfiguration(ctx, id, version, client)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicy, err)
	}

	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err := templateProcessor.ProcessTemplates(configuration); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for configuration '%s' was saved successfully\n", configName)

	return nil
}

// Find the id of a security configuration if we know its name
func findConfigurationIDByName(ctx context.Context, name string, client appsec.APPSEC) (int, int, error) {
	getConfigurationsResponse, err := client.GetConfigurations(ctx, appsec.GetConfigurationsRequest{
		ConfigID: 0,
	})
	if err != nil {
		return 0, 0, err
	}

	for _, configuration := range getConfigurationsResponse.Configurations {
		if configuration.Name == name {
			return configuration.ID, configuration.LatestVersion, nil
		}
	}

	return 0, 0, fmt.Errorf("configuration '%s' does not exist", name)
}

// Export the security config
func exportConfiguration(ctx context.Context, id int, version int, client appsec.APPSEC) (*appsec.GetExportConfigurationResponse, error) {
	getExportConfigurationResponse, err := client.GetExportConfiguration(ctx, appsec.GetExportConfigurationRequest{
		ConfigID: id,
		Version:  version,
	})
	if err != nil {
		return nil, err
	}

	return getExportConfigurationResponse, nil
}

// Get the description for the given security configuration id
func getConfigDescription(configid int) (string, error) {

	getConfigurationResponse, err := client.GetConfiguration(context.Background(), appsec.GetConfigurationRequest{
		ConfigID: configid,
	})
	if err != nil {
		return "", err
	}

	description := getConfigurationResponse.Description
	if description == "" {
		description = "Created by Terraform"
	}

	return description, nil
}

// Recursively remove ID field from structure
func removeID(dest *map[string]interface{}) {

	for _, v := range *dest {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			m := v.(map[string]interface{})
			removeID(&m)
		}
	}

	fieldNames := []string{"id", "ID"}
	for _, fieldName := range fieldNames {
		delete(*dest, fieldName)
	}
}

// Export the given struct in a JSON format appropriate for Terraform
func exportJSON(source interface{}) (string, error) {

	// Marshal Edgegrid-golang object back to Json
	js, err := json.Marshal(source)
	if err != nil {
		return "", err
	}

	// Unmarshal into our anonymous struct
	dest := map[string]interface{}{}
	err = json.Unmarshal([]byte(string(js)), &dest)
	if err != nil {
		return "", err
	}

	// Remove any ID fields
	removeID(&dest)

	// Marshal anonymous struct back into JSON
	js, err = json.MarshalIndent(dest, "", "    ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

// Get the WAF mode for the given security policy
func getWAFMode(configid int, version int, policyid string) (string, error) {

	getWAFModeResponse, err := client.GetWAFMode(context.Background(), appsec.GetWAFModeRequest{
		ConfigID: configid,
		PolicyID: policyid,
		Version:  version,
	})
	if err != nil {
		return "", err
	}

	return getWAFModeResponse.Mode, nil
}

// Get the prefix from the ID string
func getPrefixFromID(s string) string {
	split := strings.Split(s, "_")
	if len(split) > 0 {
		return split[0]
	}
	return ""
}

// Get our config section
func getSection() string {
	return section
}

// Get the reputation profile name by id
func getRepNameByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, element := range configuration.ReputationProfiles {
		if element.ID == id {
			return tools.EscapeName(element.Name)
		}
	}

	return "", errors.New("Can't find reputation profile name")
}

// Get the security policy name by id
func getPolicyNameByID(configuration *appsec.GetExportConfigurationResponse, id string) (string, error) {
	for _, element := range configuration.SecurityPolicies {
		if element.ID == id {
			return tools.EscapeName(element.Name)
		}
	}

	return "", errors.New("Can't find security policy name")
}

// Get the rate name by id
func getRateNameByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, element := range configuration.RatePolicies {
		if element.ID == id {
			return tools.EscapeName(element.Name)
		}
	}

	return "", errors.New("Can't find rate control name")
}

// Get the custom rule name by id
func getCustomRuleNameByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, rule := range configuration.CustomRules {
		if rule.ID == id {
			return tools.EscapeName(rule.Name)
		}
	}

	return "", errors.New("Can't find custom rule name")
}

// Get the rule name by id
func getRuleNameByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, element := range configuration.Rulesets {
		for _, rule := range *element.Rules {
			if rule.ID == id {
				return tools.EscapeName(rule.Tag)
			}
		}
	}

	return "", errors.New("can't find rule name")
}

// Get the rule description by id
func getRuleDescByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, element := range configuration.Rulesets {
		for _, rule := range *element.Rules {
			if rule.ID == id {
				return rule.Title, nil
			}
		}
	}

	return "", errors.New("can't find rule description")
}

// Is this a structured rule?
// Note: A "structured" rule is one that the customer has created themselves via the GUI. An "unstructured" rule is one that Akamai have created using XML. Our TF provider doesn't support unstructured rules.
func isStructuredRule(configuration *appsec.GetExportConfigurationResponse, id int) bool {
	for _, rule := range configuration.CustomRules {
		if rule.ID == id {
			if len(rule.Tag) > 0 {
				return true
			}
			break
		}
	}
	return false
}
