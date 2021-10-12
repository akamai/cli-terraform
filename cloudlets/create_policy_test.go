package cloudlets

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	common "github.com/akamai/cli-common-golang"
	"github.com/akamai/cli-terraform/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessTemplates(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	if err := os.MkdirAll("./testdata/res", 0755); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := os.RemoveAll("./testdata/res"); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestCreatePolicy(t *testing.T) {
	// TODO this is a workaround to prevent common.StartSpinner and common.StopSpinner from panicking
	// This should be removed once a dependency on "github.com/akamai/cli-common-golang" is removed
	common.App = &cli.App{ErrWriter: ioutil.Discard}

	tests := map[string]struct {
		init      func(*mockCloudlets, *mockProcessor)
		withError error
	}{
		"fetch latest version of policy and produce output": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
					{
						PolicyID:     2,
						GroupID:      234,
						Name:         "test_policy",
						Description:  "test_policy description",
						CloudletID:   0,
						CloudletCode: "ER",
					},
				}, nil).Once()
				c.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{PolicyID: 2}).Return([]cloudlets.PolicyVersion{
					{
						PolicyID: 2,
						Version:  1,
					},
					{
						PolicyID:        2,
						Version:         2,
						Description:     "version 2 description",
						MatchRuleFormat: "1.0",
					},
				}, nil).Once()
				c.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
					PolicyID: 2,
					Version:  2,
				}).Return(&cloudlets.PolicyVersion{
					PolicyID:    2,
					Version:     2,
					Description: "version 2 description",
					MatchRules: cloudlets.MatchRules{
						&cloudlets.MatchRuleER{
							Name:  "some rule",
							Type:  "ER",
							Start: 1,
							End:   2,
							ID:    1234,
						},
					},
					MatchRuleFormat: "1.0",
				}, nil).Once()
				p.On("ProcessTemplates", TFPolicyData{
					Name:            "test_policy",
					CloudletCode:    "ER",
					Description:     "version 2 description",
					GroupID:         234,
					MatchRuleFormat: "1.0",
					MatchRules: cloudlets.MatchRules{
						&cloudlets.MatchRuleER{
							Name:  "some rule",
							Type:  "ER",
							Start: 1,
							End:   2,
							ID:    1234,
						},
					},
				}).Return(nil).Once()
			},
		},
		"error fetching policy": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicy,
		},
		"error policy not found": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
				}, nil).Once()
			},
			withError: ErrFetchingPolicy,
		},
		"unsupported cloudlet type": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
					{
						PolicyID:     2,
						GroupID:      234,
						Name:         "test_policy",
						Description:  "test_policy description",
						CloudletID:   0,
						CloudletCode: "ABC",
					},
				}, nil).Once()
			},
			withError: ErrCloudletTypeNotSupported,
		},
		"error listing versions": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
					{
						PolicyID:     2,
						GroupID:      234,
						Name:         "test_policy",
						Description:  "test_policy description",
						CloudletID:   0,
						CloudletCode: "ER",
					},
				}, nil).Once()
				c.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{PolicyID: 2}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingVersion,
		},
		"error fetching latest version": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
					{
						PolicyID:     2,
						GroupID:      234,
						Name:         "test_policy",
						Description:  "test_policy description",
						CloudletID:   0,
						CloudletCode: "ER",
					},
				}, nil).Once()
				c.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{PolicyID: 2}).Return([]cloudlets.PolicyVersion{
					{
						PolicyID: 2,
						Version:  1,
					},
					{
						PolicyID:        2,
						Version:         2,
						Description:     "version 2 description",
						MatchRuleFormat: "1.0",
					},
				}, nil).Once()
				c.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
					PolicyID: 2,
					Version:  2,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingVersion,
		},
		"error processing template": {
			init: func(c *mockCloudlets, p *mockProcessor) {
				c.On("ListPolicies", mock.Anything, cloudlets.ListPoliciesRequest{}).Return([]cloudlets.Policy{
					{
						PolicyID:     1,
						GroupID:      123,
						Name:         "some policy",
						CloudletID:   0,
						CloudletCode: "ER",
					},
					{
						PolicyID:     2,
						GroupID:      234,
						Name:         "test_policy",
						Description:  "test_policy description",
						CloudletID:   0,
						CloudletCode: "ER",
					},
				}, nil).Once()
				c.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{PolicyID: 2}).Return([]cloudlets.PolicyVersion{
					{
						PolicyID: 2,
						Version:  1,
					},
					{
						PolicyID:        2,
						Version:         2,
						Description:     "version 2 description",
						MatchRuleFormat: "1.0",
					},
				}, nil).Once()
				c.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
					PolicyID: 2,
					Version:  2,
				}).Return(&cloudlets.PolicyVersion{
					PolicyID:    2,
					Version:     2,
					Description: "version 2 description",
					MatchRules: cloudlets.MatchRules{
						&cloudlets.MatchRuleER{
							Name:  "some rule",
							Type:  "ER",
							Start: 1,
							End:   2,
							ID:    1234,
						},
					},
					MatchRuleFormat: "1.0",
				}, nil).Once()
				p.On("ProcessTemplates", TFPolicyData{
					Name:            "test_policy",
					CloudletCode:    "ER",
					Description:     "version 2 description",
					GroupID:         234,
					MatchRuleFormat: "1.0",
					MatchRules: cloudlets.MatchRules{
						&cloudlets.MatchRuleER{
							Name:  "some rule",
							Type:  "ER",
							Start: 1,
							End:   2,
							ID:    1234,
						},
					},
				}).Return(fmt.Errorf("oops")).Once()
			},
			withError: ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mc := new(mockCloudlets)
			mp := new(mockProcessor)
			test.init(mc, mp)
			err := createPolicy("test_policy", mc, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			mc.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessPolicyTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData TFPolicyData
		dir       string
	}{
		"policy with match rules": {
			givenData: TFPolicyData{
				Name:            "test_policy_export",
				CloudletCode:    "ER",
				Description:     "Testing exported policy",
				GroupID:         12345,
				MatchRuleFormat: "1.0",
				MatchRules: cloudlets.MatchRules{
					cloudlets.MatchRuleER{
						Name:  "r1",
						Start: 1,
						End:   2,
						Matches: []cloudlets.MatchCriteriaER{
							{
								MatchType:     "extension",
								MatchValue:    "txt",
								MatchOperator: "equals",
							},
							{
								MatchType:     "cookie",
								MatchValue:    "cookie=cookievalue",
								MatchOperator: "equals",
								CaseSensitive: true,
							},
							{
								MatchType:     "hostname",
								MatchValue:    "3333.dom",
								MatchOperator: "equals",
								CaseSensitive: true,
								Negate:        true,
							},
						},
						UseRelativeURL:           "copy_scheme_hostname",
						StatusCode:               307,
						RedirectURL:              "/abc/sss",
						MatchURL:                 "test.url",
						UseIncomingSchemeAndHost: true,
					},
					cloudlets.MatchRuleER{
						Name:                     "r2",
						UseRelativeURL:           "copy_scheme_hostname",
						StatusCode:               301,
						RedirectURL:              "/ddd",
						MatchURL:                 "abc.com",
						UseIncomingSchemeAndHost: true,
					},
				},
			},
			dir: "with_match_rules",
		},
		"policy without match rules": {
			givenData: TFPolicyData{
				Name:            "test_policy_export",
				CloudletCode:    "ER",
				Description:     "Testing exported policy",
				GroupID:         12345,
				MatchRuleFormat: "1.0",
			},
			dir: "no_match_rules",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"policy.tmpl":    fmt.Sprintf("./testdata/res/%s/policy.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			filesToCheck := []string{"policy.tf", "variables.tf", "import.sh"}
			for _, f := range filesToCheck {
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}

}
