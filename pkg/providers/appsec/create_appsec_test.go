package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"path/filepath"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"

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
		"modules-aap-selected-hostnames.tmpl":          filepath.Join(security, "aap-selected-hostnames.tf"),
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
