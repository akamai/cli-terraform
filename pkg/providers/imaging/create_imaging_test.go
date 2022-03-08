package imaging

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/tools"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
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

var (
	policiesRequests = []imaging.GetPolicyRequest{
		{
			PolicyID:    ".auto",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "ctr_123",
			PolicySetID: "test_policyset_id",
		},
		{
			PolicyID:    "test_policy_image",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "ctr_123",
			PolicySetID: "test_policyset_id",
		},
		{
			PolicyID:    "test_policy_video",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "ctr_123",
			PolicySetID: "test_policyset_id",
		},
	}

	imagePoliciesData = []TFPolicy{
		{
			PolicyID:             ".auto",
			JSON:                 "testdata/res/json/image_policies/_auto.json",
			ActivateOnProduction: true,
		},
		{
			PolicyID:             "test_policy_image",
			JSON:                 "testdata/res/json/image_policies/test_policy_image.json",
			ActivateOnProduction: true,
		},
		{
			// In test cases where staging is not same as prod
			PolicyID:             "test_policy_image",
			JSON:                 "testdata/res/json/image_policies/test_policy_image.json",
			ActivateOnProduction: false,
		},
	}

	imagePoliciesOutputs = imaging.PolicyOutputs{
		&imaging.PolicyOutputImage{
			ID: ".auto",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   false,
		},
		&imaging.PolicyOutputImage{
			ID: "test_policy_image",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{420, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 2,
			Video:   false,
		},
		&imaging.PolicyOutputImage{
			// Return object for GetPolicy when older version is in production
			ID: "test_policy_image",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{420, 740, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   false,
		},
	}

	videoPoliciesData = []TFPolicy{
		{
			PolicyID:             ".auto",
			JSON:                 "testdata/res/json/video_policies/_auto.json",
			ActivateOnProduction: true,
		},
		{
			PolicyID:             "test_policy_video",
			JSON:                 "testdata/res/json/video_policies/test_policy_video.json",
			ActivateOnProduction: true,
		},
		{
			// In test cases where staging is not same as prod
			PolicyID:             "test_policy_video",
			JSON:                 "testdata/res/json/video_policies/test_policy_video.json",
			ActivateOnProduction: false,
		},
	}

	videoPoliciesOutputs = imaging.PolicyOutputs{
		&imaging.PolicyOutputVideo{
			ID: ".auto",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputVideo{
				PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
					Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMediumHigh),
				},
			},
			Version: 1,
			Video:   true,
		},
		&imaging.PolicyOutputVideo{
			ID: "test_policy_video",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{420, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputVideo{
				PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
					Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMediumHigh),
				},
			},
			Version: 2,
			Video:   true,
		},
		&imaging.PolicyOutputVideo{
			// Return object for GetPolicy when older version is in production
			ID: "test_policy_video",
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{420, 740, 1024, 2048, 5000},
			},
			Output: &imaging.OutputVideo{
				PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
					Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMediumHigh),
				},
			},
			Version: 1,
			Video:   true,
		},
	}

	expectImagingProcessTemplates = func(p *mockProcessor, policySetID, contractID, name, policyType, region, section string,
		policies []TFPolicy, err error) *mock.Call {
		call := p.On(
			"ProcessTemplates",
			TFImagingData{
				PolicySet: TFPolicySet{
					ID:         policySetID,
					ContractID: contractID,
					Name:       name,
					Region:     region,
					Type:       policyType,
				},
				Policies: policies,
				Section:  section,
			},
		)
		if err != nil {
			return call.Return(err)
		}
		return call.Return(nil)
	}

	expectGetPolicySet = func(i *mockimaging, policySetID, contractID, name, policyType string, region imaging.Region,
		err error) *mock.Call {
		call := i.On(
			"GetPolicySet",
			mock.Anything,
			imaging.GetPolicySetRequest{
				PolicySetID: policySetID,
				ContractID:  contractID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(
			&imaging.PolicySet{
				ID:     policySetID,
				Name:   name,
				Region: region,
				Type:   policyType,
			}, nil)
	}

	expectGetPolicy = func(i *mockimaging, policyRequest imaging.GetPolicyRequest, policyOutput imaging.PolicyOutput, err error) *mock.Call {
		call := i.On(
			"GetPolicy",
			mock.Anything,
			policyRequest,
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(policyOutput, nil)
	}

	expectListPolicies = func(i *mockimaging, policySetID, contractID, itemKind string, network imaging.PolicyNetwork,
		items imaging.PolicyOutputs, totalItems int, err error) *mock.Call {
		call := i.On(
			"ListPolicies",
			mock.Anything,
			imaging.ListPoliciesRequest{
				Network:     network,
				PolicySetID: policySetID,
				ContractID:  contractID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(
			&imaging.ListPoliciesResponse{
				ItemKind:   itemKind,
				Items:      items,
				TotalItems: totalItems,
			}, nil)
	}
)

func TestCreateImaging(t *testing.T) {
	section := "test_section"
	tests := map[string]struct {
		init         func(*mockimaging, *mockProcessor)
		filesToCheck []string
		jsonDir      string
		withError    error
	}{
		"fetch policy set with given id and contract and no policies": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns zero policies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					nil, 0, nil).Once()

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"IMAGE", "EMEA", "test_section", nil, nil).Once()
			},
		},
		"fetch policy set with image policies same on production": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[1], nil)

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"IMAGE", "EMEA", "test_section", imagePoliciesData[:2], nil).Once()
			},
			jsonDir:      "json/image_policies",
			filesToCheck: []string{"_auto.json", "test_policy_image.json"},
		},
		"fetch policy set with image policies different on production": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - different on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[2], nil)

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"IMAGE", "EMEA", "test_section", []TFPolicy{imagePoliciesData[0], imagePoliciesData[2]}, nil).Once()
			},
			jsonDir:      "json/image_policies",
			filesToCheck: []string{"_auto.json", "test_policy_image.json"},
		},

		"fetch policy set with video policies same on production": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[1], nil)

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"VIDEO", "EMEA", "test_section", videoPoliciesData[:2], nil).Once()
			},
			jsonDir:      "json/video_policies",
			filesToCheck: []string{"_auto.json", "test_policy_video.json"},
		},
		"fetch policy set with video policies different on production": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - different on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[2], nil)

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"VIDEO", "EMEA", "test_section", []TFPolicy{videoPoliciesData[0], videoPoliciesData[2]}, nil).Once()
			},
			jsonDir:      "json/video_policies",
			filesToCheck: []string{"_auto.json", "test_policy_video.json"},
		},
		"error fetching policy set": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicySet,
		},
		"error fetching policies": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs, 0, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicy,
		},
		"error fetching policy": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], fmt.Errorf("oops"))
			},
			withError: ErrFetchingPolicy,
		},
		"error processing template": {
			init: func(i *mockimaging, p *mockProcessor) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[1], nil)

				expectImagingProcessTemplates(p, "test_policyset_id", "ctr_123", "some policy set",
					"IMAGE", "EMEA", "test_section", imagePoliciesData[:2], templates.ErrSavingFiles).Once()
			},
			jsonDir:   "json/image_policies",
			withError: templates.ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("testdata/res/%s", test.jsonDir), 0755))
			mi := new(mockimaging)
			mp := new(mockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createImaging(ctx, "ctr_123", "test_policyset_id", fmt.Sprintf("testdata/res/%s", test.jsonDir), section, mi, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

			if test.filesToCheck != nil {
				for _, f := range test.filesToCheck {
					expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.jsonDir, f))
					require.NoError(t, err)
					result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.jsonDir, f))
					require.NoError(t, err)
					assert.Equal(t, string(expected), string(result))
				}
			}
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessPolicyTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFImagingData
		dir          string
		filesToCheck []string
	}{
		"policy set with no policies": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Section: "test_section",
			},
			dir:          "only_policy_set",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Policies: imagePoliciesData[:2],
				Section:  "test_section",
			},
			dir:          "with_image_policies",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with video policies": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "VIDEO",
				},
				Policies: []TFPolicy{videoPoliciesData[0], videoPoliciesData[2]},
				Section:  "test_section",
			},
			dir:          "with_video_policies",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"imaging.tmpl":   fmt.Sprintf("./testdata/res/%s/imaging.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				},
				AdditionalFuncs: template.FuncMap{
					"ToLower": func(val string) string {
						return strings.ToLower(val)
					},
					"RemoveSymbols": func(val string) string {
						return RemoveSymbols.ReplaceAllString(val, "_")
					},
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
