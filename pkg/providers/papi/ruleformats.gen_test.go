package papi

import (
	"fmt"
	"os"
	"testing"

	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Generated file, do not modify tests manually.
func TestProcessPolicyTemplatesGenerated(t *testing.T) {
	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
		filterFuncs  []func([]string) ([]string, error)
	}{
		// property with rules as datasource - hcl rules version 1 and 2 is a pair of tests that confirms that hcl rules 1 does not use any hcl rules's 2 template definitions and vice versa
		// the behaviour was chosen in a way, so it's easily identifiable which template inner definition was picked (e.g. there was change in field type)
		"property with rules as datasource - hcl rules version v2023-01-05": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2023-01-05",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2023_01_05/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"property with rules as datasource - hcl rules version v2023-05-30": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2023-05-30",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2023_05_30/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-05-30")},
		},
		"property with rules as datasource - hcl rules version v2023-09-20": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2023-09-20",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2023_09_20/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-09-20")},
		},
		"property with rules as datasource - hcl rules version v2023-10-30": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2023-10-30",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2023_10_30/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-10-30")},
		},
		"property with rules as datasource - hcl rules version v2024-01-09": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2024-01-09",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2024_01_09/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-01-09")},
		},
		"property with rules as datasource - hcl rules version v2024-02-12": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2024-02-12",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2024_02_12/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-02-12")},
		},
		"property with rules as datasource - hcl rules version v2024-05-31": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2024-05-31",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2024_05_31/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-05-31")},
		},
		"property with rules as datasource - hcl rules version v2024-08-13": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2024-08-13",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2024_08_13/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-08-13")},
		},
		"property with rules as datasource - hcl rules version v2024-10-21": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgesuite-net",
					PropertyName:         "test.edgesuite.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "v2024-10-21",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
				},
				Section: "test_section",
			},
			dir:          "ruleformats/v2024_10_21/rules_datasource",
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-10-21")},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ruleResponse := getRuleTreeResponse(test.dir, t)
			test.givenData.Rules = flattenRules("test.edgesuite.net", ruleResponse.Rules)
			test.givenData.RulesAsHCL = true

			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			templateToFile := map[string]string{
				"property.tmpl":  fmt.Sprintf("./testdata/res/%s/property.tf", test.dir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
			}

			rulesVersion := ruleResponse.RuleFormat
			templateToFile[fmt.Sprintf("rules_%s.tmpl", rulesVersion)] = fmt.Sprintf("./testdata/res/%s/rules.tf", test.dir)

			processor := templates.FSTemplateProcessor{
				TemplatesFS:     templateFiles,
				TemplateTargets: templateToFile,
				AdditionalFuncs: additionalFuncs,
			}
			err := processor.ProcessTemplates(test.givenData, test.filterFuncs...)
			require.NoError(t, err)

			for _, f := range test.filesToCheck {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
