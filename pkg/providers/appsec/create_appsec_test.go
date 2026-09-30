package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"path/filepath"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/cli-terraform/v3/pkg/templates"
	"github.com/akamai/cli-terraform/v3/pkg/tools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetConfigDescription(t *testing.T) {
	mocks := func(c *appsec.Mock) {
		c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12345}).Return(&appsec.GetConfigurationResponse{Description: "description"}, nil)
		c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12346}).Return(&appsec.GetConfigurationResponse{Description: ""}, nil)
	}

	ma := new(appsec.Mock)
	mocks(ma)

	client = ma

	description, err := getConfigDescription(12345)
	assert.NoError(t, err)
	assert.Equal(t, "description", description)

	description, err = getConfigDescription(12346)
	assert.NoError(t, err)
	assert.Equal(t, "Created by Terraform", description)
}

func TestGetWAFMode(t *testing.T) {
	mocks := func(c *appsec.Mock) {
		c.On("GetWAFMode", mock.Anything, appsec.GetWAFModeRequest{ConfigID: 12345, Version: 1, PolicyID: "ASE1_156138"}).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
	}

	ma := new(appsec.Mock)
	mocks(ma)

	client = ma

	wafMode, err := getWAFMode(12345, 1, "ASE1_156138")
	assert.NoError(t, err)
	assert.Equal(t, "KRS", wafMode)
}

func TestExportCustomDenyList(t *testing.T) {

	testdata := `{
    "name": "Deny Message",
    "ID": 12345,
    "parameters": [
        {
            "name": "response_status_code",
            "value": "433"
        }
    ]
}`

	expected := `{
    "name": "Deny Message",
    "parameters": [
        {
            "name": "response_status_code",
            "value": "433"
        }
    ]
}`

	var i map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(testdata), &i))

	actual, err := exportJSON(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExportReputationProfile(t *testing.T) {

	testdata := `{
    "ID": 12345,
    "context": "WEBATCK",
    "name": "Web Attackers (High Threat)",
    "sharedIpHandling": "NON_SHARED",
    "threshold": 9
}`

	expected := `{
    "context": "WEBATCK",
    "name": "Web Attackers (High Threat)",
    "sharedIpHandling": "NON_SHARED",
    "threshold": 9
}`

	var i map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(testdata), &i))

	actual, err := exportJSON(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExportRatePolicy(t *testing.T) {

	testdata := `{
    "ID": 12345,
    "additionalMatchOptions": null,
    "averageThreshold": 100,
    "burstThreshold": 500,
    "clientIdentifiers": ["ip"],
    "matchType": "path",
    "name": "High Rate",
    "pathMatchType": "Custom",
    "pathUriPositiveMatch": true,
    "requestType": "ClientRequest",
    "sameActionOnIpv6": true,
    "type": "WAF",
    "useXForwardForHeaders": false
}`

	expected := `{
    "additionalMatchOptions": null,
    "averageThreshold": 100,
    "burstThreshold": 500,
    "clientIdentifiers": [
        "ip"
    ],
    "matchType": "path",
    "name": "High Rate",
    "pathMatchType": "Custom",
    "pathUriPositiveMatch": true,
    "requestType": "ClientRequest",
    "sameActionOnIpv6": true,
    "type": "WAF",
    "useXForwardForHeaders": false
}`

	var i map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(testdata), &i))

	actual, err := exportJSON(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExportCustomRule(t *testing.T) {

	testdata := `{
    "conditions": [
        {
            "type": "argsPostMatch",
            "positiveMatch": true,
            "name": "email",
            "value": [
                "me@email.com"
            ]
        }
    ],
    "name": "Custom Rule 1",
    "ID": 12345,
    "tag": [
        "Login"
    ]
}`

	expected := `{
    "conditions": [
        {
            "name": "email",
            "positiveMatch": true,
            "type": "argsPostMatch",
            "value": [
                "me@email.com"
            ]
        }
    ],
    "name": "Custom Rule 1",
    "tag": [
        "Login"
    ]
}`

	var i map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(testdata), &i))

	actual, err := exportJSON(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetRepNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getRepNameByID(getExportConfigurationResponse, 3017089)
	assert.NoError(t, err)
	assert.Equal(t, "dos_attackers_high_threat", desc)
}

func TestGetPolicyNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getPolicyNameByID(getExportConfigurationResponse, "ASE1_156138")
	assert.NoError(t, err)
	assert.Equal(t, "default_policy", desc)
}

func TestGetRateNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getRateNameByID(getExportConfigurationResponse, 177906)
	assert.NoError(t, err)
	assert.Equal(t, "page_view_requests", desc)
}

func TestGetMalwareNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getMalwareNameByID(getExportConfigurationResponse, 1187)
	assert.NoError(t, err)
	assert.Equal(t, "fms_configuration_1", desc)
}

func TestGetCustomRuleNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getCustomRuleNameByID(getExportConfigurationResponse, 60088542)
	assert.NoError(t, err)
	assert.Equal(t, "custom_rule_1", desc)
}

func TestGetRuleNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getRuleNameByID(getExportConfigurationResponse, 3000080)
	assert.NoError(t, err)
	assert.Equal(t, "aseweb_attackxss", desc)
}

func TestGetRuleDescByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase")
	desc, err := getRuleDescByID(getExportConfigurationResponse, 3000080)
	assert.NoError(t, err)
	assert.Equal(t, "Cross-site Scripting (XSS) Attack (Attribute Injection 1)", desc)
}

func getExportConfigurationResponse(filename string) *appsec.GetExportConfigurationResponse {

	jsonFile, err := os.Open(fmt.Sprintf("./testdata/%s.json", filename))
	if err != nil {
		log.Fatal(err)
	}

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var getExportConfigurationResponse appsec.GetExportConfigurationResponse
	err = json.Unmarshal(byteValue, &getExportConfigurationResponse)
	if err != nil {
		log.Fatal(err)
	}

	return &getExportConfigurationResponse
}

func TestProcessPolicyTemplates(t *testing.T) {

	// Our input json for each test case
	configs := []string{"ase", "tcwest"}

	// Mocked API calls
	mocks := func(c *appsec.Mock, _ *templates.MockProcessor) {
		//c.On("GetWAFMode", mock.Anything, appsec.GetWAFModeRequest{ConfigID: 79947, Version: 1, PolicyID: "ASE1_156138"}).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		//c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 79947}).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	}

	// Additional functions for the template processor
	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})

	// Template to path mappings
	security := filepath.Join("modules", "security")
	activateSecurity := filepath.Join("modules", "activate-security")

	tests := map[string]string{
		"appsec.tmpl":                         "appsec.tf",
		"imports.tmpl":                        "appsec-import.sh",
		"main.tmpl":                           "appsec-main.tf",
		"variables.tmpl":                      "appsec-variables.tf",
		"versions.tmpl":                       "appsec-versions.tf",
		"modules-activate-security-main.tmpl": filepath.Join(activateSecurity, "main.tf"),
		"modules-activate-security-variables.tmpl":     filepath.Join(activateSecurity, "variables.tf"),
		"modules-activate-security-versions.tmpl":      filepath.Join(activateSecurity, "versions.tf"),
		"modules-security-advanced.tmpl":               filepath.Join(security, "advanced.tf"),
		"modules-security-api.tmpl":                    filepath.Join(security, "api.tf"),
		"modules-security-custom-rules.tmpl":           filepath.Join(security, "custom-rules.tf"),
		"modules-security-custom-deny.tmpl":            filepath.Join(security, "custom-deny.tf"),
		"modules-security-firewall.tmpl":               filepath.Join(security, "firewall.tf"),
		"modules-security-main.tmpl":                   filepath.Join(security, "main.tf"),
		"modules-security-malware-policies.tmpl":       filepath.Join(security, "malware-policies.tf"),
		"modules-security-malware-policy-actions.tmpl": filepath.Join(security, "malware-policy-actions.tf"),
		"modules-security-eval-penalty-box.tmpl":       filepath.Join(security, "eval-penalty-box.tf"),
		"modules-security-penalty-box.tmpl":            filepath.Join(security, "penalty-box.tf"),
		"modules-security-policies.tmpl":               filepath.Join(security, "policies.tf"),
		"modules-security-protections.tmpl":            filepath.Join(security, "protections.tf"),
		"modules-security-rate-policies.tmpl":          filepath.Join(security, "rate-policies.tf"),
		"modules-security-rate-policy-actions.tmpl":    filepath.Join(security, "rate-policy-actions.tf"),
		"modules-security-reputation.tmpl":             filepath.Join(security, "reputation.tf"),
		"modules-security-reputation-profiles.tmpl":    filepath.Join(security, "reputation-profiles.tf"),
		"modules-security-siem.tmpl":                   filepath.Join(security, "siem.tf"),
		"modules-security-slow-post.tmpl":              filepath.Join(security, "slow-post.tf"),
		"modules-security-variables.tmpl":              filepath.Join(security, "variables.tf"),
		"modules-security-versions.tmpl":               filepath.Join(security, "versions.tf"),
		"modules-security-waf.tmpl":                    filepath.Join(security, "waf.tf"),
		"modules-security-waf-ruleset.tmpl":            filepath.Join(security, "waf-ruleset.tf"),
		"modules-aap-selected-hostnames.tmpl":          filepath.Join(security, "aap-selected-hostnames.tf"),
	}

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	// Let's run our tests
	for _, config := range configs {
		for name, output := range tests {
			t.Run(name, func(t *testing.T) {

				// Create mock client
				ma := new(appsec.Mock)
				mp := new(templates.MockProcessor)
				mocks(ma, mp)
				client = ma

				// Ensure test directory exists
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/activate-security", config), 0755))

				// Run the template
				processor := templates.FSTemplateProcessor{
					TemplatesFS: templateFiles,
					TemplateTargets: map[string]string{
						name: fmt.Sprintf("./testdata/res/%s/%s", config, output),
					},
					AdditionalFuncs: additionalFuncs,
				}

				getExportConfigurationResponse := getExportConfigurationResponse(config)
				require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

				// Validate output
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, output))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, output))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			})
		}
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}
func TestProcessPolicyTemplatesWithBotman(t *testing.T) {

	// Our input json for each test case
	configs := []string{"ase-botman"}

	// Mocked API calls
	mocks := func(c *appsec.Mock, _ *templates.MockProcessor) {
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for\ndemo\n"}, nil)
	}

	botmanMocks := func(c *botman.Mock, _ *templates.MockProcessor) {
		c.On("GetAkamaiBotCategoryList", mock.Anything, mock.Anything).Return(&botman.GetAkamaiBotCategoryListResponse{Categories: []map[string]interface{}{
			{"categoryId": "0b116152-1d20-4715-8fa7-dcacb1c697e2", "categoryName": "Akamai Bot Category A"},
			{"categoryId": "da0596ba-2379-4657-9b84-79b460d66070", "categoryName": "Akamai Bot Category B"},
		}}, nil)
		c.On("GetAkamaiDefinedBotList", mock.Anything, mock.Anything).Return(&botman.GetAkamaiDefinedBotListResponse{Bots: []map[string]interface{}{
			{"botId": "eceac3f9-871b-4c57-9a24-c25b0237949a", "botName": "Akamai Defined Bot A"},
			{"botId": "c590d2e5-a041-4f05-8fda-71608f42d720", "botName": "Akamai Defined Bot B"},
		}}, nil)
		c.On("GetBotDetectionList", mock.Anything, mock.Anything).Return(&botman.GetBotDetectionListResponse{Detections: []map[string]interface{}{
			{"detectionId": "179e6bd6-5077-4f22-9a5b-3b09ee731eca", "detectionName": "Bot Detection A"},
			{"detectionId": "c4d20de1-af7a-476f-911d-73aedd97e294", "detectionName": "Bot Detection B"},
		}}, nil)
	}

	// Additional functions for the template processor
	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})

	// Template to path mappings
	security := filepath.Join("modules", "security")
	activateSecurity := filepath.Join("modules", "activate-security")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	tests := map[string]string{
		"appsec.tmpl":                         "appsec.tf",
		"imports.tmpl":                        "appsec-import.sh",
		"main.tmpl":                           "appsec-main.tf",
		"variables.tmpl":                      "appsec-variables.tf",
		"versions.tmpl":                       "appsec-versions.tf",
		"modules-activate-security-main.tmpl": filepath.Join(activateSecurity, "main.tf"),
		"modules-activate-security-variables.tmpl":      filepath.Join(activateSecurity, "variables.tf"),
		"modules-activate-security-versions.tmpl":       filepath.Join(activateSecurity, "versions.tf"),
		"modules-security-advanced.tmpl":                filepath.Join(security, "advanced.tf"),
		"modules-security-api.tmpl":                     filepath.Join(security, "api.tf"),
		"modules-security-custom-rules.tmpl":            filepath.Join(security, "custom-rules.tf"),
		"modules-security-custom-deny.tmpl":             filepath.Join(security, "custom-deny.tf"),
		"modules-security-firewall.tmpl":                filepath.Join(security, "firewall.tf"),
		"modules-security-main.tmpl":                    filepath.Join(security, "main.tf"),
		"modules-security-malware-policies.tmpl":        filepath.Join(security, "malware-policies.tf"),
		"modules-security-malware-policy-actions.tmpl":  filepath.Join(security, "malware-policy-actions.tf"),
		"modules-security-match-targets.tmpl":           filepath.Join(security, "match-targets.tf"),
		"modules-security-penalty-box.tmpl":             filepath.Join(security, "penalty-box.tf"),
		"modules-security-policies.tmpl":                filepath.Join(security, "policies.tf"),
		"modules-security-protections.tmpl":             filepath.Join(security, "protections.tf"),
		"modules-security-rate-policies.tmpl":           filepath.Join(security, "rate-policies.tf"),
		"modules-security-rate-policy-actions.tmpl":     filepath.Join(security, "rate-policy-actions.tf"),
		"modules-security-reputation.tmpl":              filepath.Join(security, "reputation.tf"),
		"modules-security-reputation-profiles.tmpl":     filepath.Join(security, "reputation-profiles.tf"),
		"modules-security-siem.tmpl":                    filepath.Join(security, "siem.tf"),
		"modules-security-slow-post.tmpl":               filepath.Join(security, "slow-post.tf"),
		"modules-security-variables.tmpl":               filepath.Join(security, "variables.tf"),
		"modules-security-versions.tmpl":                filepath.Join(security, "versions.tf"),
		"modules-security-waf.tmpl":                     filepath.Join(security, "waf.tf"),
		"modules-security-waf-ruleset.tmpl":             filepath.Join(security, "waf-ruleset.tf"),
		"modules-security-bot-directory.tmpl":           filepath.Join(security, "bot-directory.tf"),
		"modules-security-bot-directory-actions.tmpl":   filepath.Join(security, "bot-directory-actions.tf"),
		"modules-security-custom-client.tmpl":           filepath.Join(security, "custom-client.tf"),
		"modules-security-response-actions.tmpl":        filepath.Join(security, "response-actions.tf"),
		"modules-security-advanced-settings.tmpl":       filepath.Join(security, "advanced-settings.tf"),
		"modules-security-javascript-injection.tmpl":    filepath.Join(security, "javascript-injection.tf"),
		"modules-security-transactional-endpoints.tmpl": filepath.Join(security, "transactional-endpoints.tf"),
		"modules-security-content-protection.tmpl":      filepath.Join(security, "content-protection.tf"),
	}

	// Let's run our tests
	for _, config := range configs {
		for name, output := range tests {
			t.Run(name, func(t *testing.T) {

				// Create mock client
				ma := new(appsec.Mock)
				mp := new(templates.MockProcessor)
				mocks(ma, mp)
				client = ma
				mb := new(botman.Mock)
				botmanMocks(mb, new(templates.MockProcessor))
				botmanClient = mb

				// Ensure test directory exists
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/activate-security", config), 0755))

				// Run the template
				processor := templates.FSTemplateProcessor{
					TemplatesFS: templateFiles,
					TemplateTargets: map[string]string{
						name: fmt.Sprintf("./testdata/res/%s/%s", config, output),
					},
					AdditionalFuncs: additionalFuncs,
				}

				getExportConfigurationResponse := getExportConfigurationResponse(config)
				require.NoError(t, addBotmanCommonResources(context.Background(), getExportConfigurationResponse))
				require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

				// Validate output
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, output))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, output))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			})
		}
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

func TestExportJSONWithoutKeys(t *testing.T) {

	testdata := `{
    "arrayKey": [
		"arrayValue1",
		"arrayValue2"
	],
    "objectKey": {
        "innerKey": "innerValue"
    },
    "primitiveKey": "primitiveValue",
    "removeArrayKey": [
        "arrayValue1",
        "arrayValue2"
    ],
    "removeObjectKey": {
        "innerKey": "innerValue"
    },
    "removePrimitiveKey": "primitiveValue"
}`

	expected := `{
    "arrayKey": [
        "arrayValue1",
        "arrayValue2"
    ],
    "objectKey": {
        "innerKey": "innerValue"
    },
    "primitiveKey": "primitiveValue"
}`

	var i map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(testdata), &i))

	actual, err := exportJSONWithoutKeys(i, "removePrimitiveKey", "removeArrayKey", "removeObjectKey")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetCustomBotCategoryNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfigurationResponse("ase-botman")
	name, err := getCustomBotCategoryNameByID(getExportConfigurationResponse.CustomBotCategories, "dae597b8-b552-4c95-ab8b-066a3fef2f75")
	assert.NoError(t, err)
	assert.Equal(t, "category_a", name)
}

func TestExportUrlProtectionAction(t *testing.T) {
	// This test validates the url-protection-action template mapping and output
	configs := []string{"ase"}
	security := filepath.Join("modules", "security")
	templateName := "modules-security-url-protection-action.tmpl"
	outputFile := filepath.Join(security, "url-protection-action.tf")

	mocks := func(c *appsec.Mock, _ *templates.MockProcessor) {
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	}

	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for _, config := range configs {
		t.Run(templateName, func(t *testing.T) {
			ma := new(appsec.Mock)
			mp := new(templates.MockProcessor)
			mocks(ma, mp)
			client = ma

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName: fmt.Sprintf("./testdata/res/%s/%s", config, outputFile),
				},
				AdditionalFuncs: additionalFuncs,
			}

			getExportConfigurationResponse := getExportConfigurationResponse(config)
			require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, outputFile))
			require.NoError(t, err)
			result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, outputFile))
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(result))
		})
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

func TestExportUrlProtectionPolicy(t *testing.T) {

	// This test validates the url-protection-policy template mapping and output
	configs := []string{"ase"}
	security := filepath.Join("modules", "security")
	templateName := "modules-security-url-protection-policy.tmpl"
	outputFile := filepath.Join(security, "url-protection-policy.tf")

	mocks := func(c *appsec.Mock, _ *templates.MockProcessor) {
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	}

	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for _, config := range configs {
		t.Run(templateName, func(t *testing.T) {
			ma := new(appsec.Mock)
			mp := new(templates.MockProcessor)
			mocks(ma, mp)
			client = ma

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName: fmt.Sprintf("./testdata/res/%s/%s", config, outputFile),
				},
				AdditionalFuncs: additionalFuncs,
			}

			getExportConfigurationResponse := getExportConfigurationResponse(config)
			require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, outputFile))
			require.NoError(t, err)
			result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, outputFile))
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(result))
		})
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

func TestExportWAFAIRules(t *testing.T) {
	configs := []string{"ase"}
	security := filepath.Join("modules", "security")
	templateName := "modules-security-waf-ai-rules.tmpl"
	outputFile := filepath.Join(security, "waf-ai-rules.tf")

	mocks := func(c *appsec.Mock, _ *templates.MockProcessor) {
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	}

	additionalFuncs := tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for _, config := range configs {
		t.Run(templateName, func(t *testing.T) {
			ma := new(appsec.Mock)
			mp := new(templates.MockProcessor)
			mocks(ma, mp)
			client = ma

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName: fmt.Sprintf("./testdata/res/%s/%s", config, outputFile),
				},
				AdditionalFuncs: additionalFuncs,
			}

			getExportConfigurationResponse := getExportConfigurationResponse(config)
			require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, outputFile))
			require.NoError(t, err)
			result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, outputFile))
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(result))
		})
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

// wafRulesetAdditionalFuncs returns the template function map shared by all
func wafRulesetAdditionalFuncs() map[string]any {
	return tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getCustomRuleNameByID":                      getCustomRuleNameByID,
		"getRepNameByID":                             getRepNameByID,
		"getRuleNameByID":                            getRuleNameByID,
		"getRuleDescByID":                            getRuleDescByID,
		"getRateNameByID":                            getRateNameByID,
		"getMalwareNameByID":                         getMalwareNameByID,
		"getPolicyNameByID":                          getPolicyNameByID,
		"getWAFMode":                                 getWAFMode,
		"getConfigDescription":                       getConfigDescription,
		"getPrefixFromID":                            getPrefixFromID,
		"getEdgercPath":                              getEdgercPath,
		"getSection":                                 getSection,
		"isStructuredRule":                           isStructuredRule,
		"exportJSON":                                 exportJSON,
		"exportJSONWithoutKeys":                      exportJSONWithoutKeys,
		"getCustomBotCategoryNameByID":               getCustomBotCategoryNameByID,
		"getCustomBotCategoryResourceNamesByIDs":     getCustomBotCategoryResourceNamesByIDs,
		"getCustomClientResourceNamesByIDs":          getCustomClientResourceNamesByIDs,
		"getContentProtectionRuleResourceNamesByIDs": getContentProtectionRuleResourceNamesByIDs,
		"getProtectedHostsByID":                      getProtectedHostsByID,
		"getEvaluatedHostsByID":                      getEvaluatedHostsByID,
		"exportJSONForCustomDefBotsWithoutKeys":      exportJSONForCustomDefBotsWithoutKeys,
		"buildCategoryMap":                           buildCategoryMap,
		"getRapidRulesByPolicyID":                    getRapidRulesByPolicyID,
		"exportRapidRulesJSON":                       exportRapidRulesJSON,
	})
}

// setupWAFRulesetMocks registers the API mock expectations shared by all
func setupWAFRulesetMocks(c *appsec.Mock, _ *templates.MockProcessor) {
	c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
	c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
}

func TestWAFRulesetTemplate(t *testing.T) {
	// This test validates the waf-ruleset template mapping and output
	configs := []string{"ase", "tcwest", "ase-botman", "ase-apr"}
	security := filepath.Join("modules", "security")
	templateName := "modules-security-waf-ruleset.tmpl"
	outputFile := filepath.Join(security, "waf-ruleset.tf")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for _, config := range configs {
		t.Run(templateName+"-"+config, func(t *testing.T) {
			ma := new(appsec.Mock)
			mp := new(templates.MockProcessor)
			setupWAFRulesetMocks(ma, mp)
			client = ma

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName: fmt.Sprintf("./testdata/res/%s/%s", config, outputFile),
				},
				AdditionalFuncs: wafRulesetAdditionalFuncs(),
			}

			resp := getExportConfigurationResponse(config)
			require.NoError(t, processor.ProcessTemplates(resp))

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, outputFile))
			require.NoError(t, err)
			result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, outputFile))
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(result))
		})
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

func TestWAFRulesetTemplateWithEmptyRulesAndAttackGroups(t *testing.T) {
	// This test validates the waf-ruleset template is NOT generated when both rules and attack_groups are empty
	security := filepath.Join("modules", "security")
	templateName := "modules-security-waf-ruleset.tmpl"
	outFile := filepath.Join(t.TempDir(), "ase", security, "waf-ruleset.tf")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	ma := new(appsec.Mock)
	mp := new(templates.MockProcessor)
	setupWAFRulesetMocks(ma, mp)
	client = ma

	require.NoError(t, os.MkdirAll(filepath.Dir(outFile), 0755))

	processor := templates.FSTemplateProcessor{
		TemplatesFS: templateFiles,
		TemplateTargets: map[string]string{
			templateName: outFile,
		},
		AdditionalFuncs: wafRulesetAdditionalFuncs(),
	}

	config := getExportConfigurationResponse("ase")
	require.NotEmpty(t, config.SecurityPolicies)

	// Force template fallback branches - set both RuleActions and AttackGroupActions to empty
	for i := range config.SecurityPolicies {
		require.NotNil(t, config.SecurityPolicies[i].WebApplicationFirewall)
		config.SecurityPolicies[i].WebApplicationFirewall.RuleActions = nil
		config.SecurityPolicies[i].WebApplicationFirewall.AttackGroupActions = nil
	}

	require.NoError(t, processor.ProcessTemplates(config))

	// When both RuleActions and AttackGroupActions are empty, the resource should not be generated
	// The file may be empty or may not exist depending on the processor behavior
	result, err := os.ReadFile(outFile)
	if err != nil {
		// File doesn't exist or can't be read - that's acceptable for empty templates
		return
	}

	out := string(result)
	// When both RuleActions and AttackGroupActions are empty, the resource should not be generated at all
	assert.NotContains(t, out, "akamai_appsec_waf_ruleset")
	assert.NotContains(t, out, "rules")
	assert.NotContains(t, out, "attack_groups")
	// Output should be empty or only contain whitespace
	assert.True(t, strings.TrimSpace(out) == "", "expected empty output when both RuleActions and AttackGroupActions are empty")
}

func TestWAFRulesetTemplateWithEmptyRulesOrAttackGroups(t *testing.T) {
	// This test validates the waf-ruleset template output when rules or attack_groups are empty
	security := filepath.Join("modules", "security")
	templateName := "modules-security-waf-ruleset.tmpl"
	outFile := filepath.Join(t.TempDir(), "ase", security, "waf-ruleset.tf")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	ma := new(appsec.Mock)
	mp := new(templates.MockProcessor)
	setupWAFRulesetMocks(ma, mp)
	client = ma

	require.NoError(t, os.MkdirAll(filepath.Dir(outFile), 0755))

	processor := templates.FSTemplateProcessor{
		TemplatesFS: templateFiles,
		TemplateTargets: map[string]string{
			templateName: outFile,
		},
		AdditionalFuncs: wafRulesetAdditionalFuncs(),
	}

	config := getExportConfigurationResponse("ase")
	require.NotEmpty(t, config.SecurityPolicies)

	// Force template fallback branches that should render empty Terraform lists.
	for i := range config.SecurityPolicies {
		require.NotNil(t, config.SecurityPolicies[i].WebApplicationFirewall)
		config.SecurityPolicies[i].WebApplicationFirewall.AttackGroupActions = nil
	}

	require.NoError(t, processor.ProcessTemplates(config))

	result, err := os.ReadFile(outFile)
	require.NoError(t, err)

	out := string(result)
	// When only AttackGroupActions is empty, rules should be populated but attack_groups should be empty
	assert.Contains(t, out, "rules = [")
	assert.Regexp(t, `(?m)^\s*attack_groups\s*=\s*\[\]$`, out)
	assert.NotContains(t, out, "attack_groups = [{")
}

// rapidRulesConditionException is the condition exception attached to the second mocked rapid rule.
// It exercises the branch in convertRapidRulesToRuleDefinitions which copies a non-empty
// conditionException through to the emitted rule definition.
func rapidRulesConditionException() *appsec.RuleConditionException {
	return &appsec.RuleConditionException{
		Exception: &appsec.RuleException{
			HeaderCookieOrParamValues: []string{"exception-value"},
			SpecificHeaderCookieOrParamPrefix: &appsec.SpecificHeaderCookieOrParamPrefixPtr{
				Prefix:   "exception-prefix",
				Selector: "REQUEST_HEADERS",
			},
			SpecificHeaderCookieParamXMLOrJSONNames: &appsec.SpecificHeaderCookieParamXMLOrJSONNames{
				{
					Names:    []string{"Auth"},
					Selector: "REQUEST_HEADERS",
					Wildcard: true,
				},
			},
		},
	}
}

func rapidRulesAdvancedConditionException() *appsec.RuleConditionException {
	return &appsec.RuleConditionException{
		AdvancedExceptionsList: &appsec.AdvancedExceptions{
			ConditionOperator: "AND",
			SpecificHeaderCookieParamXMLOrJSONNames: &appsec.AttackGroupSpecificHeaderCookieParamXMLOrJSONNamesAdvanced{
				{
					Names:    []string{"X-Trace-Id"},
					Selector: "REQUEST_HEADERS",
					Wildcard: false,
				},
			},
		},
	}
}

// setupRapidRulesMocks registers the API mock expectations needed to render the rapid rules
// template for the ase config, whose single policy (ASE1_156138) has rapid rules enabled.
func setupRapidRulesMocks(c *appsec.Mock) {
	c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
	c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	c.On("GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{
		ConfigID: 79947,
		Version:  1,
		PolicyID: "ASE1_156138",
	}).Return(&appsec.GetRapidRulesResponse{
		Rules: []appsec.PolicyRapidRule{
			{
				ID:              3000101,
				Action:          "deny",
				Lock:            false,
				Name:            "Rapid Rule One",
				Version:         1,
				RiskScoreGroups: []string{"SQL"},
			},
			{
				ID:                 3000102,
				Action:             "alert",
				Lock:               true,
				Name:               "Rapid Rule Two",
				Version:            2,
				RiskScoreGroups:    []string{"XSS"},
				ConditionException: rapidRulesConditionException(),
			},
			{
				ID:                 3000103,
				Action:             "deny",
				Lock:               false,
				Name:               "Rapid Rule Three",
				Version:            1,
				RiskScoreGroups:    []string{"CMD"},
				ConditionException: rapidRulesAdvancedConditionException(),
			},
		},
	}, nil)
	c.On("GetRapidRulesDefaultAction", mock.Anything, appsec.GetRapidRulesDefaultActionRequest{
		ConfigID: 79947,
		Version:  1,
		PolicyID: "ASE1_156138",
	}).Return(&appsec.GetRapidRulesDefaultActionResponse{Action: "alert"}, nil)
}

func TestExportRapidRules(t *testing.T) {
	// This test validates the rapid-rules template mapping and output for a policy which has
	// rapid rules enabled in the exported security configuration.
	configs := []string{"ase"}
	security := filepath.Join("modules", "security")
	templateName := "modules-security-rapid-rules.tmpl"
	outputFile := filepath.Join(security, "rapid-rules.tf")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for _, config := range configs {
		t.Run(templateName+"-"+config, func(t *testing.T) {
			ma := new(appsec.Mock)
			setupRapidRulesMocks(ma)
			client = ma

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s/modules/security", config), 0755))

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName: fmt.Sprintf("./testdata/res/%s/%s", config, outputFile),
				},
				AdditionalFuncs: wafRulesetAdditionalFuncs(),
			}

			getExportConfigurationResponse := getExportConfigurationResponse(config)
			require.NoError(t, addRapidRulesResources(context.Background(), getExportConfigurationResponse))
			require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, outputFile))
			require.NoError(t, err)
			result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, outputFile))
			require.NoError(t, err)
			assert.Equal(t, string(expected), string(result))

			ma.AssertCalled(t, "GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{
				ConfigID: 79947,
				Version:  1,
				PolicyID: "ASE1_156138",
			})
			ma.AssertCalled(t, "GetRapidRulesDefaultAction", mock.Anything, appsec.GetRapidRulesDefaultActionRequest{
				ConfigID: 79947,
				Version:  1,
				PolicyID: "ASE1_156138",
			})
		})
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}

func TestAddRapidRulesResources(t *testing.T) {
	// This test validates the enriched side map built by addRapidRulesResources, and in particular
	// the exact json shape of the rule definitions which end up in the rule_definitions attribute.
	// The akamai_appsec_rapid_rules resource deserializes rule_definitions with
	// DisallowUnknownFields, so only id, action, lock and conditionException may be emitted.
	ma := new(appsec.Mock)
	setupRapidRulesMocks(ma)
	client = ma

	config := getExportConfigurationResponse("ase")
	require.NoError(t, addRapidRulesResources(context.Background(), config))

	require.Len(t, rapidRulesByPolicyID, 1)
	data := rapidRulesByPolicyID["ASE1_156138"]
	require.NotNil(t, data)
	assert.Equal(t, "alert", data.DefaultAction)
	require.Len(t, data.RuleDefinitions, 3)

	// First rule: no condition exception on the source rapid rule, so none must be emitted.
	assert.Equal(t, int64(3000101), *data.RuleDefinitions[0].ID)
	assert.Equal(t, "deny", *data.RuleDefinitions[0].Action)
	assert.False(t, *data.RuleDefinitions[0].Lock)
	assert.Nil(t, data.RuleDefinitions[0].ConditionException)

	// Second rule: non-empty condition exception must be carried through unchanged.
	assert.Equal(t, int64(3000102), *data.RuleDefinitions[1].ID)
	assert.Equal(t, "alert", *data.RuleDefinitions[1].Action)
	assert.True(t, *data.RuleDefinitions[1].Lock)
	assert.Equal(t, rapidRulesConditionException(), data.RuleDefinitions[1].ConditionException)

	assert.Equal(t, int64(3000103), *data.RuleDefinitions[2].ID)
	assert.Equal(t, "deny", *data.RuleDefinitions[2].Action)
	assert.False(t, *data.RuleDefinitions[2].Lock)
	require.NotNil(t, data.RuleDefinitions[2].ConditionException)
	assert.Nil(t, data.RuleDefinitions[2].ConditionException.Exception)
	assert.Equal(t, rapidRulesAdvancedConditionException(), data.RuleDefinitions[2].ConditionException)

	// The serialized shape must contain exactly the keys the resource accepts.
	serialized, err := exportRapidRulesJSON(data.RuleDefinitions)
	require.NoError(t, err)

	var definitions []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(serialized), &definitions))
	require.Len(t, definitions, 3)

	assert.ElementsMatch(t, []string{"id", "action", "lock"}, keysOf(definitions[0]))
	assert.ElementsMatch(t, []string{"id", "action", "lock", "conditionException"}, keysOf(definitions[1]))
	assert.ElementsMatch(t, []string{"id", "action", "lock", "conditionException"}, keysOf(definitions[2]))

	assert.Equal(t, float64(3000101), definitions[0]["id"])
	assert.Equal(t, "deny", definitions[0]["action"])
	assert.Equal(t, false, definitions[0]["lock"])
	assert.Equal(t, float64(3000102), definitions[1]["id"])
	assert.Equal(t, "alert", definitions[1]["action"])
	assert.Equal(t, true, definitions[1]["lock"])
	assert.Equal(t, float64(3000103), definitions[2]["id"])
	assert.Equal(t, "deny", definitions[2]["action"])
	assert.Equal(t, false, definitions[2]["lock"])

	thirdCE, ok := definitions[2]["conditionException"].(map[string]interface{})
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"advancedExceptions"}, keysOf(thirdCE))

	ma.AssertNumberOfCalls(t, "GetRapidRules", 1)
	ma.AssertNumberOfCalls(t, "GetRapidRulesDefaultAction", 1)
}

// keysOf returns the sorted keys of a json object, used to assert the exact serialized shape.
func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestAddRapidRulesResourcesWithConditionExceptionEmpty(t *testing.T) {
	// An all-nil condition exception object must not be emitted; the resource would otherwise
	// receive a conditionException with no content.
	ma := new(appsec.Mock)
	ma.On("GetRapidRules", mock.Anything, mock.Anything).Return(&appsec.GetRapidRulesResponse{
		Rules: []appsec.PolicyRapidRule{
			{
				ID:                 3000101,
				Action:             "deny",
				ConditionException: &appsec.RuleConditionException{},
			},
		},
	}, nil)
	ma.On("GetRapidRulesDefaultAction", mock.Anything, mock.Anything).Return(&appsec.GetRapidRulesDefaultActionResponse{Action: "deny"}, nil)
	client = ma

	config := getExportConfigurationResponse("ase")
	require.NoError(t, addRapidRulesResources(context.Background(), config))

	data := rapidRulesByPolicyID["ASE1_156138"]
	require.NotNil(t, data)
	require.Len(t, data.RuleDefinitions, 1)
	assert.Nil(t, data.RuleDefinitions[0].ConditionException)

	serialized, err := exportRapidRulesJSON(data.RuleDefinitions)
	require.NoError(t, err)
	assert.NotContains(t, serialized, "conditionException")
}

func TestExportRapidRulesJSONPreservesID(t *testing.T) {
	// exportRapidRulesJSON must preserve the id field. The generic exportJSON helper calls
	// removeID, which strips it; the akamai_appsec_rapid_rules resource requires the id of every
	// rule definition, so swapping exportRapidRulesJSON back to exportJSON must fail this test.
	ruleID := int64(3000101)
	action := "deny"
	lock := true
	definitions := []appsec.RuleDefinition{
		{
			ID:     &ruleID,
			Action: &action,
			Lock:   &lock,
		},
	}

	expected := `[
    {
        "id": 3000101,
        "action": "deny",
        "lock": true
    }
]`

	actual, err := exportRapidRulesJSON(definitions)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.Contains(t, actual, `"id": 3000101`)

	// Guard rail: demonstrate that the generic helper drops the id, which is why this helper exists.
	stripped, err := exportJSON(definitions[0])
	require.NoError(t, err)
	assert.NotContains(t, stripped, `"id"`)
	assert.Contains(t, stripped, `"action": "deny"`)
}

func TestExportRapidRulesDisabledPolicy(t *testing.T) {
	// Negative case: neither the akamai_appsec_rapid_rules resource nor the import line may be
	// emitted for policies where rapidRules.enabled is false, or where rapidRules is absent.
	tests := map[string]struct {
		config string
		mutate func(*appsec.GetExportConfigurationResponse)
	}{
		// The tcwest fixture has one policy with rapidRules present and disabled, and two policies
		// with no rapidRules block at all.
		"rapid rules present but disabled": {
			config: "tcwest",
			mutate: func(c *appsec.GetExportConfigurationResponse) {
				require.NotNil(t, c.SecurityPolicies[0].RapidRules)
				require.False(t, c.SecurityPolicies[0].RapidRules.Enabled)
			},
		},
		"rapid rules absent entirely": {
			config: "tcwest",
			mutate: func(c *appsec.GetExportConfigurationResponse) {
				for i := range c.SecurityPolicies {
					c.SecurityPolicies[i].RapidRules = nil
				}
			},
		},
	}

	security := filepath.Join("modules", "security")
	templateName := "modules-security-rapid-rules.tmpl"
	outputFile := filepath.Join(security, "rapid-rules.tf")

	edgercPath = "/non/default/path/to/edgerc"
	section = "non-default-section"

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ma := new(appsec.Mock)
			ma.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
			ma.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
			client = ma

			baseDir := t.TempDir()
			require.NoError(t, os.MkdirAll(filepath.Join(baseDir, security), 0755))

			resourcePath := filepath.Join(baseDir, outputFile)
			importPath := filepath.Join(baseDir, "appsec-import.sh")

			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					templateName:   resourcePath,
					"imports.tmpl": importPath,
				},
				AdditionalFuncs: wafRulesetAdditionalFuncs(),
			}

			config := getExportConfigurationResponse(test.config)
			test.mutate(config)
			require.NoError(t, addRapidRulesResources(context.Background(), config))

			// No rapid rules api calls may be made for policies without rapid rules enabled.
			assert.Empty(t, rapidRulesByPolicyID)
			ma.AssertNotCalled(t, "GetRapidRules", mock.Anything, mock.Anything)
			ma.AssertNotCalled(t, "GetRapidRulesDefaultAction", mock.Anything, mock.Anything)

			require.NoError(t, processor.ProcessTemplates(config))

			// The template produces empty output, so no rapid-rules.tf is written at all.
			if resource, err := os.ReadFile(resourcePath); err == nil {
				assert.NotContains(t, string(resource), "akamai_appsec_rapid_rules")
				assert.Empty(t, strings.TrimSpace(string(resource)))
			}

			// The import script must not contain a rapid rules import line.
			importScript, err := os.ReadFile(importPath)
			require.NoError(t, err)
			assert.NotEmpty(t, strings.TrimSpace(string(importScript)))
			assert.NotContains(t, string(importScript), "akamai_appsec_rapid_rules")
		})
	}
}
