// Package appsec contains code for exporting application security configuration
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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/botman"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templateFiles embed.FS
var client appsec.APPSEC
var botmanClient botman.BotMan

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
	sess := edgegrid.GetSession(ctx)
	client = appsec.Client(sess)
	botmanClient = botman.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	// Directory Paths
	modulesPath := filepath.Join(tfWorkPath, "modules")
	securityModulePath := filepath.Join(modulesPath, "security")
	activateSecurityModulePath := filepath.Join(modulesPath, "activate-security")
	paths := []string{modulesPath, securityModulePath, activateSecurityModulePath}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return cli.Exit(color.RedString(err.Error()), 1)
		}
	}

	// File Paths
	appsecPath := filepath.Join(tfWorkPath, "appsec.tf")

	err := tools.CheckFiles(appsecPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	// Save our section for use later
	section = edgegrid.GetEdgercSection(c)

	// Template to path mappings
	templateToFile := map[string]string{
		"appsec.tmpl":                         appsecPath,
		"imports.tmpl":                        filepath.Join(tfWorkPath, "appsec-import.sh"),
		"main.tmpl":                           filepath.Join(tfWorkPath, "appsec-main.tf"),
		"modules-activate-security-main.tmpl": filepath.Join(activateSecurityModulePath, "main.tf"),
		"modules-activate-security-variables.tmpl":      filepath.Join(activateSecurityModulePath, "variables.tf"),
		"modules-activate-security-versions.tmpl":       filepath.Join(activateSecurityModulePath, "versions.tf"),
		"modules-security-advanced.tmpl":                filepath.Join(securityModulePath, "advanced.tf"),
		"modules-security-api.tmpl":                     filepath.Join(securityModulePath, "api.tf"),
		"modules-security-custom-deny.tmpl":             filepath.Join(securityModulePath, "custom-deny.tf"),
		"modules-security-custom-rules.tmpl":            filepath.Join(securityModulePath, "custom-rules.tf"),
		"modules-security-firewall.tmpl":                filepath.Join(securityModulePath, "firewall.tf"),
		"modules-security-main.tmpl":                    filepath.Join(securityModulePath, "main.tf"),
		"modules-security-malware-policies.tmpl":        filepath.Join(securityModulePath, "malware-policies.tf"),
		"modules-security-malware-policy-actions.tmpl":  filepath.Join(securityModulePath, "malware-policy-actions.tf"),
		"modules-security-match-targets.tmpl":           filepath.Join(securityModulePath, "match-targets.tf"),
		"modules-security-penalty-box.tmpl":             filepath.Join(securityModulePath, "penalty-box.tf"),
		"modules-security-eval-penalty-box.tmpl":        filepath.Join(securityModulePath, "eval-penalty-box.tf"),
		"modules-security-policies.tmpl":                filepath.Join(securityModulePath, "policies.tf"),
		"modules-security-protections.tmpl":             filepath.Join(securityModulePath, "protections.tf"),
		"modules-security-rate-policies.tmpl":           filepath.Join(securityModulePath, "rate-policies.tf"),
		"modules-security-rate-policy-actions.tmpl":     filepath.Join(securityModulePath, "rate-policy-actions.tf"),
		"modules-security-reputation-profiles.tmpl":     filepath.Join(securityModulePath, "reputation-profiles.tf"),
		"modules-security-reputation.tmpl":              filepath.Join(securityModulePath, "reputation.tf"),
		"modules-security-siem.tmpl":                    filepath.Join(securityModulePath, "siem.tf"),
		"modules-security-slow-post.tmpl":               filepath.Join(securityModulePath, "slow-post.tf"),
		"modules-security-variables.tmpl":               filepath.Join(securityModulePath, "variables.tf"),
		"modules-security-versions.tmpl":                filepath.Join(securityModulePath, "versions.tf"),
		"modules-security-waf.tmpl":                     filepath.Join(securityModulePath, "waf.tf"),
		"modules-security-bot-directory.tmpl":           filepath.Join(securityModulePath, "bot-directory.tf"),
		"modules-security-bot-directory-actions.tmpl":   filepath.Join(securityModulePath, "bot-directory-actions.tf"),
		"modules-security-custom-client.tmpl":           filepath.Join(securityModulePath, "custom-client.tf"),
		"modules-security-response-actions.tmpl":        filepath.Join(securityModulePath, "response-actions.tf"),
		"modules-security-advanced-settings.tmpl":       filepath.Join(securityModulePath, "advanced-settings.tf"),
		"modules-security-javascript-injection.tmpl":    filepath.Join(securityModulePath, "javascript-injection.tf"),
		"modules-security-transactional-endpoints.tmpl": filepath.Join(securityModulePath, "transactional-endpoints.tf"),
		"modules-security-content-protection.tmpl":      filepath.Join(securityModulePath, "content-protection.tf"),
		"modules-aap-selected-hostnames.tmpl":           filepath.Join(securityModulePath, "aap-selected-hostnames.tf"),
		"variables.tmpl":                                filepath.Join(tfWorkPath, "appsec-variables.tf"),
		"versions.tmpl":                                 filepath.Join(tfWorkPath, "appsec-versions.tf"),
	}

	// Provide custom helper functions to get data that does not exist in the security config export
	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"exportJSON":                                 exportJSON,
		"getConfigDescription":                       getConfigDescription,
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getPrefixFromID":                            getPrefixFromID,
		"getRateNameByID":                            getRateNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getSection":                                 getSection,
		"getWAFMode":                                 getWAFMode,
		"isStructuredRule":                           isStructuredRule,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
	})

	// The template processor
	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFuncs,
	}

	appsecName := c.Args().First()
	if err = createAppsec(ctx, appsecName, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting appsec config HCL: %s", err), 1)
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

	if err := addBotmanCommonResources(ctx, configuration); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("error fetching botman common values: %s", err)
	}

	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err := templateProcessor.ProcessTemplates(configuration); err != nil {
		term.Spinner().Fail()
		fmt.Printf("------------------FAILED ----------------------")
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
		Source:   "TF",
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
	err = json.Unmarshal(js, &dest)
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
func getWAFMode(configID int, version int, policyID string) (string, error) {

	getWAFModeResponse, err := client.GetWAFMode(context.Background(), appsec.GetWAFModeRequest{
		ConfigID: configID,
		PolicyID: policyID,
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

// Get the malware policy name by id
func getMalwareNameByID(configuration *appsec.GetExportConfigurationResponse, id int) (string, error) {
	for _, element := range configuration.MalwarePolicies {
		if element.MalwarePolicyID == id {
			return tools.EscapeName(element.Name)
		}
	}

	return "", errors.New("Can't find malware policy name")
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

// addBotmanCommonResources makes api call to get akamaiBotCategories, botDetection and akamaiDefinedBots to fetch names to be used in terraform resource names
func addBotmanCommonResources(ctx context.Context, configuration *appsec.GetExportConfigurationResponse) error {
	hasAkamaiBotCategoryAction := false
	for _, policy := range configuration.SecurityPolicies {
		if policy.BotManagement != nil && len(policy.BotManagement.AkamaiBotCategoryActions) > 0 {
			hasAkamaiBotCategoryAction = true
			break
		}
	}
	if hasAkamaiBotCategoryAction {
		var akamaiBotCategoryList, err = botmanClient.GetAkamaiBotCategoryList(ctx, botman.GetAkamaiBotCategoryListRequest{})
		if err != nil {
			return err
		}
		akamaiBotCategoryMap := make(map[string]string)
		for _, akamaiBotCategory := range akamaiBotCategoryList.Categories {
			akamaiBotCategoryMap[akamaiBotCategory["categoryId"].(string)] = akamaiBotCategory["categoryName"].(string)
		}
		for _, policy := range configuration.SecurityPolicies {
			if policy.BotManagement == nil {
				continue
			}
			for _, akamaiBotCategoryAction := range policy.BotManagement.AkamaiBotCategoryActions {
				akamaiBotCategoryAction["categoryName"] = akamaiBotCategoryMap[akamaiBotCategoryAction["categoryId"].(string)]
			}
		}
	}

	hasBotDetectionAction := false
	for _, policy := range configuration.SecurityPolicies {
		if policy.BotManagement != nil && len(policy.BotManagement.BotDetectionActions) > 0 {
			hasBotDetectionAction = true
			break
		}
	}
	if hasBotDetectionAction {
		var botDetectionList, err = botmanClient.GetBotDetectionList(ctx, botman.GetBotDetectionListRequest{})
		if err != nil {
			return err
		}
		botDetectionMap := make(map[string]string)
		for _, botDetection := range botDetectionList.Detections {
			botDetectionMap[botDetection["detectionId"].(string)] = botDetection["detectionName"].(string)
		}
		for _, policy := range configuration.SecurityPolicies {
			if policy.BotManagement == nil {
				continue
			}
			for _, botDetectionAction := range policy.BotManagement.BotDetectionActions {
				botDetectionAction["detectionName"] = botDetectionMap[botDetectionAction["detectionId"].(string)]
			}
		}
	}
	var akamaiBotIDMap map[string]string
	for _, category := range configuration.CustomBotCategories {
		if category["metadata"] == nil {
			continue
		}
		metadata := category["metadata"].(map[string]interface{})
		if metadata["akamaiDefinedBotIds"] == nil {
			continue
		}
		akamaiDefinedBotIDs := metadata["akamaiDefinedBotIds"].([]interface{})
		if len(akamaiDefinedBotIDs) == 0 {
			continue
		}
		if akamaiBotIDMap == nil {
			akamaiDefinedBotList, err := botmanClient.GetAkamaiDefinedBotList(ctx, botman.GetAkamaiDefinedBotListRequest{})
			if err != nil {
				return err
			}
			akamaiBotIDMap = make(map[string]string)
			for _, akamaiBot := range akamaiDefinedBotList.Bots {
				akamaiBotIDMap[akamaiBot["botId"].(string)] = akamaiBot["botName"].(string)
			}
		}
		recategorizedAkamaiDefinedBot := make([]map[string]string, len(akamaiDefinedBotIDs))
		for i, botID := range akamaiDefinedBotIDs {
			akamaiBot := make(map[string]string)
			akamaiBot["botId"] = botID.(string)
			akamaiBot["botName"] = akamaiBotIDMap[botID.(string)]
			recategorizedAkamaiDefinedBot[i] = akamaiBot
		}
		metadata["akamaiDefinedBots"] = recategorizedAkamaiDefinedBot
	}

	return nil
}

// exportJSONWithoutKeys returns json string without specified keys
func exportJSONWithoutKeys(source map[string]interface{}, keys ...string) (string, error) {
	// deep copy source by converting to json
	js, err := json.Marshal(source)
	if err != nil {
		return "", err
	}
	dest := make(map[string]interface{})
	err = json.Unmarshal(js, &dest)
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		delete(dest, key)
	}

	js, err = json.MarshalIndent(dest, "", "    ")
	if err != nil {
		return "", err
	}

	return string(js), nil
}

func getCustomBotCategoryNameByID(customBotCategories []map[string]interface{}, categoryID string) (string, error) {
	for _, category := range customBotCategories {
		if category["categoryId"].(string) == categoryID {
			return tools.EscapeName(category["categoryName"].(string))
		}
	}
	return "", fmt.Errorf("cannot find custom bot category name for id %s", categoryID)
}

// getCustomBotCategoryResourceNamesByIDs returns comma separated custom bot category resource names in the same order as the provided categoryIDs
func getCustomBotCategoryResourceNamesByIDs(customBotCategories []map[string]interface{}, categoryIDs []string) (string, error) {
	customBotCategoryMap := make(map[string]string)
	for _, category := range customBotCategories {
		categoryName, ok := category["categoryName"].(string)
		if !ok {
			return "", errors.New("cannot convert categoryName to string")
		}
		name, err := tools.EscapeName(categoryName)
		if err != nil {
			return "", err
		}
		categoryID, ok := category["categoryId"].(string)
		if !ok {
			return "", errors.New("cannot convert categoryId to string")
		}
		customBotCategoryMap[categoryID] = name
	}
	categoryResourceNames := make([]string, len(categoryIDs))
	for i, categoryID := range categoryIDs {
		categoryName, ok := customBotCategoryMap[categoryID]
		if !ok {
			return "", fmt.Errorf("cannot find custom bot category name for id %s", categoryID)
		}
		categoryResourceNames[i] = fmt.Sprintf("akamai_botman_custom_bot_category.%s_%s.category_id", categoryName, categoryID)
	}
	return strings.Join(categoryResourceNames, ","), nil
}

// getCustomClientResourceNamesByIDs returns comma separated custom-client resource names in the same order as the provided customClientIDs
func getCustomClientResourceNamesByIDs(customClients []map[string]interface{}, customClientIDs []string) (string, error) {
	customClientMap := make(map[string]string)
	for _, customClient := range customClients {
		customClientName, ok := customClient["customClientName"].(string)
		if !ok {
			return "", fmt.Errorf("cannot convert custom client name %s to string", customClient["customClientName"])
		}
		name, err := tools.EscapeName(customClientName)
		if err != nil {
			return "", err
		}
		customClientID, ok := customClient["customClientId"].(string)
		if !ok {
			return "", errors.New("cannot convert customClientId to string")
		}
		customClientMap[customClientID] = name
	}
	customClientResourceNames := make([]string, len(customClientIDs))
	for i, customClientID := range customClientIDs {
		customClientName, ok := customClientMap[customClientID]
		if !ok {
			return "", fmt.Errorf("cannot find custom client name for id %s", customClientID)
		}
		customClientResourceNames[i] = fmt.Sprintf("akamai_botman_custom_client.%s_%s.custom_client_id", customClientName, customClientID)
	}
	return strings.Join(customClientResourceNames, ",\n"), nil
}

// getContentProtectionRuleResourceNamesByIDs returns comma separated content-protection-rule resource names in the same order as the provided contentProtectionRuleIDs
func getContentProtectionRuleResourceNamesByIDs(policyID string, contentProtectionRules []map[string]interface{}, contentProtectionRuleIDs []string) (string, error) {
	contentProtectionRuleMap := make(map[string]string)
	for _, contentProtectionRule := range contentProtectionRules {
		contentProtectionRuleName, ok := contentProtectionRule["contentProtectionRuleName"].(string)
		if !ok {
			return "", fmt.Errorf("cannot convert content protection rule name %s to string", contentProtectionRule["contentProtectionRuleName"])
		}
		name, err := tools.EscapeName(contentProtectionRuleName)
		if err != nil {
			return "", err
		}
		contentProtectionRuleID, ok := contentProtectionRule["contentProtectionRuleId"].(string)
		if !ok {
			return "", errors.New("cannot convert contentProtectionRuleId to string")
		}
		contentProtectionRuleMap[contentProtectionRuleID] = name
	}
	contentProtectionRuleResourceNames := make([]string, len(contentProtectionRuleIDs))
	for i, contentProtectionRuleID := range contentProtectionRuleIDs {
		contentProtectionRuleName, ok := contentProtectionRuleMap[contentProtectionRuleID]
		if !ok {
			return "", fmt.Errorf("cannot find content protection rule for id %s", contentProtectionRuleID)
		}
		contentProtectionRuleResourceNames[i] = fmt.Sprintf("akamai_botman_content_protection_rule.%s_%s_%s.content_protection_rule_id", policyID, contentProtectionRuleName, contentProtectionRuleID)
	}
	return strings.Join(contentProtectionRuleResourceNames, ",\n"), nil
}

// Get the protected hosts by policy id
func getProtectedHostsByID(configuration *appsec.GetExportConfigurationResponse, id string) ([]string, error) {
	protectedHosts := make([]string, 0)
	for _, websiteTarget := range configuration.MatchTargets.WebsiteTargets {
		if websiteTarget.SecurityPolicy.PolicyID == id {
			return websiteTarget.Hostnames, nil
		}
	}
	return protectedHosts, nil
}

// Get the protected hosts by policy id
func getEvaluatedHostsByID(configuration *appsec.GetExportConfigurationResponse, id string) ([]string, error) {
	evaluatingHostnames := make([]string, 0)
	for _, policy := range configuration.Evaluating.SecurityPolicies {
		if policy.SecurityPolicyID == id {
			return policy.Hostnames, nil
		}
	}
	return evaluatingHostnames, nil
}
