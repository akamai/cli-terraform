package appsec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/cli-terraform/pkg/templates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessTemplates(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

func TestGetConfigDescription(t *testing.T) {
	mocks := func(c *mockAppsec) {
		c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12345}).Return(&appsec.GetConfigurationResponse{Description: "description"}, nil)
		c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12346}).Return(&appsec.GetConfigurationResponse{Description: ""}, nil)
	}

	ma := new(mockAppsec)
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
	mocks := func(c *mockAppsec) {
		c.On("GetWAFMode", mock.Anything, appsec.GetWAFModeRequest{ConfigID: 12345, Version: 1, PolicyID: "ASE1_156138"}).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
	}

	ma := new(mockAppsec)
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
	assert.NoError(t, json.Unmarshal([]byte(string(testdata)), &i))

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
	assert.NoError(t, json.Unmarshal([]byte(string(testdata)), &i))

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
    "clientIdentifier": "ip",
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
    "clientIdentifier": "ip",
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
	assert.NoError(t, json.Unmarshal([]byte(string(testdata)), &i))

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
	assert.NoError(t, json.Unmarshal([]byte(string(testdata)), &i))

	actual, err := exportJSON(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetRepNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getRepNameByID(getExportConfigurationResponse, 3017089)
	assert.NoError(t, err)
	assert.Equal(t, "dos_attackers_high_threat", desc)
}

func TestGetPolicyNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getPolicyNameByID(getExportConfigurationResponse, "ASE1_156138")
	assert.NoError(t, err)
	assert.Equal(t, "default_policy", desc)
}

func TestGetRateNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getRateNameByID(getExportConfigurationResponse, 177906)
	assert.NoError(t, err)
	assert.Equal(t, "page_view_requests", desc)
}

func TestGetCustomRuleNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getCustomRuleNameByID(getExportConfigurationResponse, 60088542)
	assert.NoError(t, err)
	assert.Equal(t, "custom_rule_1", desc)
}

func TestGetRuleNameByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getRuleNameByID(getExportConfigurationResponse, 3000080)
	assert.NoError(t, err)
	assert.Equal(t, "aseweb_attackxss", desc)
}

func TestGetRuleDescByID(t *testing.T) {
	getExportConfigurationResponse := getExportConfiguratonResponse("ase")
	desc, err := getRuleDescByID(getExportConfigurationResponse, 3000080)
	assert.NoError(t, err)
	assert.Equal(t, "Cross-site Scripting (XSS) Attack (Attribute Injection 1)", desc)
}

func getExportConfiguratonResponse(filename string) *appsec.GetExportConfigurationResponse {

	jsonFile, err := os.Open(fmt.Sprintf("./testdata/%s.json", filename))
	if err != nil {
		log.Fatal(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var getExportConfigurationResponse appsec.GetExportConfigurationResponse
	err = json.Unmarshal(byteValue, &getExportConfigurationResponse)
	return &getExportConfigurationResponse
}

func TestProcessPolicyTemplates(t *testing.T) {

	// Our input json for each test case
	configs := []string{"ase", "tcwest"}

	// Mocked API calls
	mocks := func(c *mockAppsec, p *mockProcessor) {
		//c.On("GetWAFMode", mock.Anything, appsec.GetWAFModeRequest{ConfigID: 79947, Version: 1, PolicyID: "ASE1_156138"}).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		c.On("GetWAFMode", mock.Anything, mock.Anything).Return(&appsec.GetWAFModeResponse{Mode: "KRS"}, nil)
		//c.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 79947}).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
		c.On("GetConfiguration", mock.Anything, mock.Anything).Return(&appsec.GetConfigurationResponse{Description: "A security config for demo"}, nil)
	}

	// Additional functions for the template processor
	additionalFuncs := template.FuncMap{
		"getCustomRuleNameByID": getCustomRuleNameByID,
		"getRepNameByID":        getRepNameByID,
		"getRuleNameByID":       getRuleNameByID,
		"getRuleDescByID":       getRuleDescByID,
		"getRateNameByID":       getRateNameByID,
		"getPolicyNameByID":     getPolicyNameByID,
		"getWAFMode":            getWAFMode,
		"getConfigDescription":  getConfigDescription,
		"getPrefixFromID":       getPrefixFromID,
		"getSection":            getSection,
		"isStructuredRule":      isStructuredRule,
		"exportJSON":            exportJSON,
	}

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
		"modules-activate-security-variables.tmpl":  filepath.Join(activateSecurity, "variables.tf"),
		"modules-activate-security-versions.tmpl":   filepath.Join(activateSecurity, "versions.tf"),
		"modules-security-advanced.tmpl":            filepath.Join(security, "advanced.tf"),
		"modules-security-api.tmpl":                 filepath.Join(security, "api.tf"),
		"modules-security-custom-rules.tmpl":        filepath.Join(security, "custom-rules.tf"),
		"modules-security-custom-deny.tmpl":         filepath.Join(security, "custom-deny.tf"),
		"modules-security-firewall.tmpl":            filepath.Join(security, "firewall.tf"),
		"modules-security-main.tmpl":                filepath.Join(security, "main.tf"),
		"modules-security-match-targets.tmpl":       filepath.Join(security, "match-targets.tf"),
		"modules-security-penalty-box.tmpl":         filepath.Join(security, "penalty-box.tf"),
		"modules-security-policies.tmpl":            filepath.Join(security, "policies.tf"),
		"modules-security-protections.tmpl":         filepath.Join(security, "protections.tf"),
		"modules-security-rate-policies.tmpl":       filepath.Join(security, "rate-policies.tf"),
		"modules-security-rate-policy-actions.tmpl": filepath.Join(security, "rate-policy-actions.tf"),
		"modules-security-reputation.tmpl":          filepath.Join(security, "reputation.tf"),
		"modules-security-reputation-profiles.tmpl": filepath.Join(security, "reputation-profiles.tf"),
		"modules-security-selected-hostnames.tmpl":  filepath.Join(security, "selected-hostnames.tf"),
		"modules-security-siem.tmpl":                filepath.Join(security, "siem.tf"),
		"modules-security-slow-post.tmpl":           filepath.Join(security, "slow-post.tf"),
		"modules-security-variables.tmpl":           filepath.Join(security, "variables.tf"),
		"modules-security-versions.tmpl":            filepath.Join(security, "versions.tf"),
		"modules-security-waf.tmpl":                 filepath.Join(security, "waf.tf"),
	}

	// Let's run our tests
	for _, config := range configs {
		for name, output := range tests {
			t.Run(name, func(t *testing.T) {

				// Create mock client
				ma := new(mockAppsec)
				mp := new(mockProcessor)
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

				getExportConfigurationResponse := getExportConfiguratonResponse(config)
				require.NoError(t, processor.ProcessTemplates(getExportConfigurationResponse))

				// Validate output
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", config, output))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", config, output))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			})
		}
	}
	require.NoError(t, os.RemoveAll("./testdata/res"))
}
