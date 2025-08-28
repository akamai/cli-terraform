package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	section    = "test_section"
	contractID = "test_contract"

	expectListIncludes = func(client *papi.Mock) {
		listIncludesReq := papi.ListIncludesRequest{
			ContractID: "test_contract",
		}

		includes := papi.ListIncludesResponse{
			Includes: papi.IncludeItems{
				Items: []papi.Include{
					{
						AccountID:         "test_account",
						AssetID:           "test_asset",
						ContractID:        "test_contract",
						GroupID:           "test_group",
						IncludeID:         "inc_123456",
						IncludeName:       "test_include",
						IncludeType:       papi.IncludeTypeMicroServices,
						LatestVersion:     2,
						StagingVersion:    ptr.To(1),
						ProductionVersion: ptr.To(1),
					},
					{
						AccountID:      "test_account_1",
						AssetID:        "test_asset_1",
						ContractID:     "test_contract",
						GroupID:        "test_group",
						IncludeID:      "inc_456789",
						IncludeName:    "test_include_1",
						IncludeType:    papi.IncludeTypeCommonSettings,
						LatestVersion:  1,
						StagingVersion: ptr.To(1),
					},
				},
			},
		}

		client.On("ListIncludes", mock.Anything, listIncludesReq).Return(&includes, nil).Once()
	}

	expectGetIncludeVersion = func(client *papi.Mock, format string) {
		getIncludeVersionReq := papi.GetIncludeVersionRequest{
			ContractID: "test_contract",
			GroupID:    "test_group",
			IncludeID:  "inc_123456",
			Version:    2,
		}

		version := papi.GetIncludeVersionResponse{
			AccountID:   "test_account",
			AssetID:     "test_asset",
			ContractID:  "test_contract",
			GroupID:     "test_group",
			IncludeID:   "inc_123456",
			IncludeName: "test_include",
			IncludeType: papi.IncludeTypeMicroServices,
			IncludeVersions: papi.Versions{
				Items: []papi.IncludeVersion{
					{
						UpdatedByUser:    "test_user",
						UpdatedDate:      "2022-08-22T07:17:48Z",
						ProductID:        "test_product",
						ProductionStatus: papi.VersionStatusInactive,
						Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
						RuleFormat:       format,
						IncludeVersion:   2,
						StagingStatus:    papi.VersionStatusInactive,
					},
				},
			},
			IncludeVersion: papi.IncludeVersion{
				UpdatedByUser:    "test_user",
				UpdatedDate:      "2022-08-22T07:17:48Z",
				ProductID:        "test_product",
				ProductionStatus: papi.VersionStatusInactive,
				Etag:             "1d8ed19bce0833a3fe93e62ae5d5579a38cc2dbe",
				RuleFormat:       format,
				IncludeVersion:   2,
				StagingStatus:    papi.VersionStatusInactive,
			},
		}

		client.On("GetIncludeVersion", mock.Anything, getIncludeVersionReq).Return(&version, nil).Once()
	}

	getIncludeRuleTreeReq = papi.GetIncludeRuleTreeRequest{
		ContractID:     "test_contract",
		GroupID:        "test_group",
		IncludeID:      "inc_123456",
		IncludeVersion: 2,
		RuleFormat:     "v2020-11-02",
	}

	getIncludeRuleTreeReqRulesAsHCL = papi.GetIncludeRuleTreeRequest{
		ContractID:     "test_contract",
		GroupID:        "test_group",
		IncludeID:      "inc_123456",
		IncludeVersion: 2,
		RuleFormat:     "v2023-01-05",
	}

	expectListIncludeActivations = func(client *papi.Mock) {
		listIncludeActivationsReq := papi.ListIncludeActivationsRequest{
			ContractID: "test_contract",
			GroupID:    "test_group",
			IncludeID:  "inc_123456",
		}

		activations := papi.ListIncludeActivationsResponse{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "test_group",
			Activations: papi.IncludeActivationsRes{
				Items: []papi.IncludeActivation{
					{
						ActivationID:       "atv_12344",
						Network:            papi.ActivationNetworkStaging,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-10-27T12:27:54Z",
						UpdateDate:         "2022-10-27T12:28:54Z",
						Note:               "test staging activation",
						NotifyEmails:       []string{"test@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12343",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-10-27T12:27:54Z",
						UpdateDate:         "2022-10-27T12:28:54Z",
						Note:               "test production activation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12342",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeDeactivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-09-27T12:27:54Z",
						UpdateDate:         "2022-09-27T12:28:54Z",
						Note:               "test production deactivation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
					{
						ActivationID:       "atv_12341",
						Network:            papi.ActivationNetworkProduction,
						ActivationType:     papi.ActivationTypeActivate,
						Status:             papi.ActivationStatusActive,
						SubmitDate:         "2022-08-27T12:27:54Z",
						UpdateDate:         "2022-08-27T12:28:54Z",
						Note:               "test production old activation",
						NotifyEmails:       []string{"test@example.com", "test1@example.com"},
						FMAActivationState: "steady",
						FallbackInfo: &papi.ActivationFallbackInfo{
							FastFallbackAttempted:      false,
							FallbackVersion:            1,
							CanFastFallback:            false,
							SteadyStateTime:            1666873734,
							FastFallbackExpirationTime: 1666877334,
						},
						IncludeID:      "inc_123456",
						IncludeName:    "test_include",
						IncludeType:    papi.IncludeTypeMicroServices,
						IncludeVersion: 1,
					},
				},
			},
		}

		client.On("ListIncludeActivations", mock.Anything, listIncludeActivationsReq).Return(&activations, nil).Once()
	}

	expectAllProcessTemplates = func(p *templates.MockProcessor, testData TFData, filterFuncs ...func([]string) ([]string, error)) *mock.Call {
		var call *mock.Call
		if len(filterFuncs) != 0 {
			call = p.On("ProcessTemplates", testData, mock.AnythingOfType("func([]string) ([]string, error)"))
		} else {
			call = p.On("ProcessTemplates", testData)
		}

		return call.Return(nil)
	}
)

func TestCreateInclude(t *testing.T) {
	tests := map[string]struct {
		init                func(*papi.Mock, *templates.MockProcessor, *templates.MockMultiTargetProcessor, string)
		includeName         string
		dir                 string
		snippetFilesToCheck []string
		jsonDir             string
		withError           error
		rulesAsHCL          bool
		splitDepth          *int
	}{
		"include basic": {
			init: func(c *papi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2020-11-02")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				expectAllProcessTemplates(p, getTestData("include basic"))
			},

			includeName: "test_include",
			dir:         "include_basic",
			jsonDir:     "include_basic/property-snippets",
			snippetFilesToCheck: []string{
				"test_include.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"include basic rules as hcl": {
			init: func(c *papi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2023-01-05")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReqRulesAsHCL).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				expectAllProcessTemplates(p, (&tfDataBuilder{}).withData(getTestData("include basic rules as hcl")).
					withIncludeRules(0, flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))).
					build(), useThisOnlyRuleFormat("v2023-01-05"))
				mockAddTemplateTargetRules(p)
				mockAddTemplateTargetIncludesRules(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)

			},

			includeName: "test_include",
			dir:         "include_basic_rules_as_hcl",
			rulesAsHCL:  true,
		},
		"include basic, unsupported schema version": {
			init: func(c *papi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2020-11-02")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				expectAllProcessTemplates(p, (&tfDataBuilder{}).withData(getTestData("include basic")).
					withIncludeRules(0, flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))).
					build())
				mockAddTemplateTargetRules(p)
				mockTemplateExist(p, "rules_v2020-11-02.tmpl", false)
			},

			includeName: "test_include",
			dir:         "include_basic",
			rulesAsHCL:  true,
			withError:   ErrUnsupportedRuleFormat,
		},
		"include with children with split-depth=0": {
			init: func(c *papi.Mock, p *templates.MockProcessor, mm *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2023-01-05")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReqRulesAsHCL).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				tfData := (&tfDataBuilder{}).withData(getTestData("include basic rules as hcl")).
					withSplitDepth(true, "").
					asHCL(true).
					build()
				tfData.Includes[0].RootRule = "test_include_rule_default"
				expectAllProcessTemplates(p, tfData, useThisOnlyRuleFormat("v2023-01-05"))
				mockModuleConfig(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)

				multiTargetData := templates.MultiTargetData{
					"includes_rules.tmpl": templates.DataForTarget{
						"rules/test_include_default.tf": (&tfDataBuilder{}).
							withIncludes([]TFIncludeData{{}}).
							withIncludeRules(0, flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))).
							withSplitDepth(true, "").build(),
					},
				}
				mm.On("ProcessTemplates", multiTargetData, mock.AnythingOfType("func([]string) ([]string, error)")).Return(nil).Once()
			},

			includeName: "test_include",
			dir:         "include_multitarget",
			rulesAsHCL:  true,
			splitDepth:  ptr.To(0),
		},
		"include with children with split-depth=2": {
			init: func(c *papi.Mock, p *templates.MockProcessor, mm *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2023-01-05")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReqRulesAsHCL).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)
				tfData := (&tfDataBuilder{}).withData(getTestData("include basic rules as hcl")).
					withSplitDepth(true, "").
					asHCL(true).
					build()
				tfData.Includes[0].RootRule = "test_include_rule_default"
				expectAllProcessTemplates(p, tfData, useThisOnlyRuleFormat("v2023-01-05"))
				mockModuleConfig(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)

				multiTargetData := templates.MultiTargetData{
					"includes_rules.tmpl": templates.DataForTarget{
						"rules/test_include_default.tf": (&tfDataBuilder{}).
							withIncludes([]TFIncludeData{{}}).
							withIncludeRules(0, []*WrappedRules{flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))[0]}).
							withSplitDepth(true, "").build(),
						"rules/test_include_default_new_rule.tf": (&tfDataBuilder{}).
							withIncludes([]TFIncludeData{{}}).
							withIncludeRules(0, []*WrappedRules{flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))[1]}).
							withSplitDepth(true, "").build(),
						"rules/test_include_default_new_rule_new_rule_1.tf": (&tfDataBuilder{}).
							withIncludes([]TFIncludeData{{}}).
							withIncludeRules(0, []*WrappedRules{flattenRules(wrapAndNameRules("test_include", ruleResponse.Rules))[2]}).
							withSplitDepth(true, "").build(),
					},
				}
				mm.On("ProcessTemplates", multiTargetData, mock.AnythingOfType("func([]string) ([]string, error)")).Return(nil).Once()
			},

			includeName: "test_include",
			dir:         "include_multitarget_flatten",
			rulesAsHCL:  true,
			splitDepth:  ptr.To(2),
		},
		"error include not found": {
			init: func(c *papi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				c.On("ListIncludes", mock.Anything, papi.ListIncludesRequest{ContractID: "test_contract"}).
					Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrIncludeNotFound,
		},
		"error fetching include version": {
			init: func(c *papi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				expectListIncludes(c)
				c.On("GetIncludeVersion", mock.Anything, papi.GetIncludeVersionRequest{
					ContractID: "test_contract",
					GroupID:    "test_group",
					IncludeID:  "inc_123456",
					Version:    2,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrFetchingLatestIncludeVersion,
			includeName: "test_include",
		},
		"error include rules not found": {
			init: func(c *papi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2020-11-02")
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrIncludeRulesNotFound,
			includeName: "test_include",
		},
		"error fetching activations": {
			init: func(c *papi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2020-11-02")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				c.On("ListIncludeActivations", mock.Anything, papi.ListIncludeActivationsRequest{
					ContractID: "test_contract",
					GroupID:    "test_group",
					IncludeID:  "inc_123456",
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError:   ErrFetchingActivations,
			includeName: "test_include",
			dir:         "include_basic",
		},
		"error saving files": {
			init: func(c *papi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				expectListIncludes(c)
				expectGetIncludeVersion(c, "v2020-11-02")

				// Rule Tree
				var ruleResponse papi.GetIncludeRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(rules, &ruleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&ruleResponse, nil).Once()

				expectListIncludeActivations(c)

				p.On("ProcessTemplates", getTestData("include basic")).Return(fmt.Errorf("oops")).Once()
			},
			withError:   ErrSavingFiles,
			includeName: "test_include",
			dir:         "include_basic",
			jsonDir:     "include_basic/property-snippets",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mc := new(papi.Mock)
			mp := new(templates.MockProcessor)
			mm := new(templates.MockMultiTargetProcessor)
			test.init(mc, mp, mm, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			options := includeOptions{
				contractID:  contractID,
				includeName: test.includeName,
				section:     section,
				jsonDir:     fmt.Sprintf("./testdata/res/%s", test.jsonDir),
				tfWorkPath:  "./",
				rulesAsHCL:  test.rulesAsHCL,
				splitDepth:  test.splitDepth,
			}
			err := createInclude(ctx, options, mc, mp, mm)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			for _, f := range test.snippetFilesToCheck {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.jsonDir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.jsonDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
			require.NoError(t, err)
			mc.AssertExpectations(t)
			mp.AssertExpectations(t)
			mm.AssertExpectations(t)
		})
	}
}

func TestProcessIncludeTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
		rulesAsHCL   bool
		filterFuncs  []func([]string) ([]string, error)
	}{
		"include basic": {
			givenData:    getTestData("include basic"),
			dir:          "include_basic",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
		"include basic rules as hcl": {
			givenData:    getTestData("include basic rules as hcl"),
			dir:          "include_basic_rules_as_hcl",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh", "includes_rules.tf"},
			rulesAsHCL:   true,
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"include basic rules using split-depth": {
			givenData:    getTestData("include basic rules using split-depth"),
			dir:          "include_multitarget",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh", "module_config.tf"},
			rulesAsHCL:   true,
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"include single network": {
			givenData:    getTestData("include single network"),
			dir:          "include_single_network",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
		"include no network": {
			givenData:    getTestData("include no network"),
			dir:          "include_no_network",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
		"include basic with multiline note": {
			givenData:    getTestData("include with multiline notes"),
			dir:          "include_basic_multiline_notes",
			filesToCheck: []string{"includes.tf", "variables.tf", "import.sh"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.rulesAsHCL {
				ruleResponse := getRuleTreeResponse(test.dir, t)
				test.givenData.Includes[0].Rules = flattenRules(wrapAndNameRules(test.givenData.Includes[0].IncludeName, ruleResponse.Rules))
				test.givenData.RulesAsHCL = true
			}
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			templateToFile := map[string]string{
				"includes.tmpl":  fmt.Sprintf("./testdata/res/%s/includes.tf", test.dir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
			}
			if test.rulesAsHCL && !test.givenData.UseSplitDepth {
				templateToFile["includes_rules.tmpl"] = fmt.Sprintf("./testdata/res/%s/includes_rules.tf", test.dir)
			}
			if test.givenData.UseSplitDepth {
				templateToFile["rules_module.tmpl"] = fmt.Sprintf("./testdata/res/%s/module_config.tf", test.dir)
			}
			processor := templates.FSTemplateProcessor{
				TemplatesFS:     templateFiles,
				TemplateTargets: templateToFile,
				AdditionalFuncs: additionalFuncs,
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData, test.filterFuncs...))

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

func TestMultiTargetProcessIncludeTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    templates.MultiTargetData
		dir          string
		filesToCheck []string
	}{
		"include with children (split-depth=0)": {
			givenData: templates.MultiTargetData{
				"includes_rules.tmpl": {
					"./testdata/res/include_multitarget/include_default.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
						Includes: []TFIncludeData{
							{
								Rules: []*WrappedRules{
									{
										IsRoot: true,
										Rule: papi.Rules{
											Name: "default",
											Children: []papi.Rules{
												{
													Name: "New Rule",
													Children: []papi.Rules{
														{
															Name: "New Rule 1",
														},
													},
												},
											},
										},
										Children: []*WrappedRules{
											{
												Rule: papi.Rules{
													Name: "New Rule",
													Children: []papi.Rules{
														{
															Name: "New Rule 1",
														},
													},
												},
												FileName:      "include_default_new_rule",
												TerraformName: "include_rule_new_rule",
												Children: []*WrappedRules{
													{
														Rule:          papi.Rules{Name: "New Rule 1"},
														FileName:      "include_default_new_rule_new_rule_1",
														TerraformName: "include_rule_new_rule_1",
													},
												},
											},
										},
										FileName:      "include_default",
										TerraformName: "include_rule_default",
									},
									{
										Rule: papi.Rules{
											Name: "New Rule",
											Children: []papi.Rules{
												{
													Name: "New Rule 1",
												},
											},
										},
										FileName:      "include_default_new_rule",
										TerraformName: "include_rule_new_rule",
										Children: []*WrappedRules{
											{
												Rule:          papi.Rules{Name: "New Rule 1"},
												FileName:      "include_default_new_rule_new_rule_1",
												TerraformName: "include_rule_new_rule_1",
											},
										},
									},
									{
										Rule:          papi.Rules{Name: "New Rule 1"},
										FileName:      "include_default_new_rule_new_rule_1",
										TerraformName: "include_rule_new_rule_1",
									},
								},
							},
						},
					},
				},
			},
			dir:          "include_multitarget",
			filesToCheck: []string{"include_default.tf"},
		},
		"include with children (split-depth=2)": {
			givenData: templates.MultiTargetData{
				"includes_rules.tmpl": {
					"./testdata/res/include_multitarget_flatten/include_default.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
						Includes: []TFIncludeData{
							{
								Rules: []*WrappedRules{
									{
										IsRoot: true,
										Rule: papi.Rules{
											Name: "default",
											Children: []papi.Rules{
												{
													Name: "New Rule",
													Children: []papi.Rules{
														{
															Name: "New Rule 1",
														},
													},
												},
											},
										},
										Children: []*WrappedRules{
											{
												Rule: papi.Rules{
													Name: "New Rule",
													Children: []papi.Rules{
														{
															Name: "New Rule 1",
														},
													},
												},
												FileName:      "include_default_new_rule",
												TerraformName: "include_rule_new_rule",
												Children: []*WrappedRules{
													{
														Rule:          papi.Rules{Name: "New Rule 1"},
														FileName:      "include_default_new_rule_new_rule_1",
														TerraformName: "include_rule_new_rule_1",
													},
												},
											},
										},
										FileName:      "include_default",
										TerraformName: "include_rule_default",
									},
								},
							},
						},
					},
					"./testdata/res/include_multitarget_flatten/include_default_new_rule.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
						Includes: []TFIncludeData{
							{
								Rules: []*WrappedRules{
									{
										Rule: papi.Rules{
											Name: "New Rule",
											Children: []papi.Rules{
												{
													Name: "New Rule 1",
												},
											},
										},
										FileName:      "include_default_new_rule",
										TerraformName: "include_rule_new_rule",
										Children: []*WrappedRules{
											{
												Rule:          papi.Rules{Name: "New Rule 1"},
												FileName:      "include_default_new_rule_new_rule_1",
												TerraformName: "include_rule_new_rule_1",
											},
										},
									},
								},
							},
						},
					},
					"./testdata/res/include_multitarget_flatten/include_default_new_rule_new_rule_1.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
						Includes: []TFIncludeData{
							{
								Rules: []*WrappedRules{
									{
										Rule:          papi.Rules{Name: "New Rule 1"},
										FileName:      "include_default_new_rule_new_rule_1",
										TerraformName: "include_rule_new_rule_1",
									},
								},
							},
						},
					},
				},
			},
			dir:          "include_multitarget_flatten",
			filesToCheck: []string{"include_default.tf", "include_default_new_rule_new_rule_1.tf", "include_default_new_rule.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSMultiTargetProcessor{
				TemplatesFS:     templateFiles,
				AdditionalFuncs: additionalFuncs,
			}
			err := processor.ProcessTemplates(test.givenData, useThisOnlyRuleFormat("v2023-01-05"))
			reportedErrors = []string{}
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

func getTestData(key string) TFData {
	TFDataMap := map[string]TFData{
		"include basic": {
			Section: section,
			Includes: []TFIncludeData{
				{
					StagingInfo: NetworkInfo{
						ActivationNote:          "test staging activation",
						Emails:                  []string{"test@example.com"},
						Version:                 1,
						HasActivation:           true,
						IsActiveOnLatestVersion: false,
					},
					ProductionInfo: NetworkInfo{
						ActivationNote:          "test production activation",
						Emails:                  []string{"test@example.com", "test1@example.com"},
						Version:                 1,
						HasActivation:           true,
						IsActiveOnLatestVersion: false,
					},
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2020-11-02",
				},
			},
		},
		"include basic rules as hcl": {
			Section: section,
			Includes: []TFIncludeData{
				{
					StagingInfo: NetworkInfo{
						ActivationNote: "test staging activation",
						Emails:         []string{"test@example.com"},
						Version:        1,
						HasActivation:  true,
					},
					ProductionInfo: NetworkInfo{
						ActivationNote: "test production activation",
						Emails:         []string{"test@example.com", "test1@example.com"},
						Version:        1,
						HasActivation:  true,
					},
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2023-01-05",
				},
			},
		},
		"include basic rules using split-depth": {
			Section: section,
			Includes: []TFIncludeData{
				{
					StagingInfo: NetworkInfo{
						ActivationNote: "test staging activation",
						Emails:         []string{"test@example.com"},
						Version:        1,
						HasActivation:  true,
					},
					ProductionInfo: NetworkInfo{
						ActivationNote: "test production activation",
						Emails:         []string{"test@example.com", "test1@example.com"},
						Version:        1,
						HasActivation:  true,
					},
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2023-01-05",
					RootRule:    "test_include_default",
				},
			},
			UseSplitDepth: true,
			RulesAsHCL:    true,
		},
		"include single network": {
			Section: section,
			Includes: []TFIncludeData{
				{
					StagingInfo: NetworkInfo{
						Emails:        []string{"test@example.com"},
						Version:       3,
						HasActivation: true,
					},
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2020-11-02",
				},
			},
		},
		"include no network": {
			Section: section,
			Includes: []TFIncludeData{
				{
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2020-11-02",
				},
			},
		},
		"include with multiline notes": {
			Section: section,
			Includes: []TFIncludeData{
				{
					StagingInfo: NetworkInfo{
						ActivationNote:          "first\nsecond\n\nlast",
						Emails:                  []string{"test@example.com"},
						Version:                 1,
						HasActivation:           true,
						IsActiveOnLatestVersion: false,
					},
					ProductionInfo: NetworkInfo{
						ActivationNote:          "first\nsecond\n",
						Emails:                  []string{"test@example.com", "test1@example.com"},
						Version:                 2,
						HasActivation:           true,
						IsActiveOnLatestVersion: true,
					},
					ContractID:  "test_contract",
					GroupID:     "test_group",
					IncludeID:   "inc_123456",
					IncludeName: "test_include",
					IncludeType: string(papi.IncludeTypeMicroServices),
					ProductID:   "test_product",
					RuleFormat:  "v2020-11-02",
				},
			},
		},
		"basic property with multiple children as hcl": {
			Property: TFPropertyData{
				GroupName:            "test_group",
				GroupID:              "grp_12345",
				ContractID:           "test_contract",
				PropertyResourceName: "test-edgesuite-net",
				PropertyName:         "test.edgesuite.net",
				PropertyID:           "prp_12345",
				ProductID:            "prd_HTTP_Content_Del",
				ProductName:          "HTTP_Content_Del",
				RuleFormat:           "latest",
				IsSecure:             "false",
				EdgeHostnames: map[string]EdgeHostname{
					"test-edgesuite-net": {
						EdgeHostname:             "test.edgesuite.net",
						EdgeHostnameID:           "ehn_2867480",
						ContractID:               "test_contract",
						GroupID:                  "grp_12345",
						ID:                       "",
						IPv6:                     "IPV6_COMPLIANCE",
						SecurityType:             "STANDARD-TLS",
						EdgeHostnameResourceName: "test-edgesuite-net",
					},
				},
				Hostnames: map[string]Hostname{
					"test.edgesuite.net": {
						CnameFrom:                "test.edgesuite.net",
						CnameTo:                  "test.edgesuite.net",
						EdgeHostnameResourceName: "test-edgesuite-net",
						CertProvisioningType:     "CPS_MANAGED",
						IsActive:                 true,
					},
				},
				StagingInfo: NetworkInfo{
					Emails:                  []string{"jsmith@akamai.com"},
					HasActivation:           true,
					Version:                 5,
					IsActiveOnLatestVersion: true,
				},
				ReadVersion: "LATEST",
			},
			Section: "test_section",
		},
	}

	return TFDataMap[key]
}
