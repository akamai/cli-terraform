package imaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
	veryDeepPolicy = imaging.PolicyInputImage{
		RolloutDuration: ptr.To(3600),
		Hosts:           []string{"host1", "host2"},
		Variables: []imaging.Variable{
			{
				Name:         "ResizeDim",
				Type:         "number",
				DefaultValue: "280",
			},
			{
				Name:         "ResizeDimWithBorder",
				Type:         "number",
				DefaultValue: "260",
			},
			{
				Name:         "VariableWithoutDefaultValue",
				Type:         "string",
				DefaultValue: "",
			},
			{
				Name:         "MinDim",
				Type:         "number",
				DefaultValue: "1000",
				EnumOptions: []*imaging.EnumOptions{
					{
						ID:    "1",
						Value: "value1",
					},
					{
						ID:    "2",
						Value: "value2",
					},
				},
			},
			{
				Name:         "MinDimNew",
				Type:         "number",
				DefaultValue: "1450",
			},
			{
				Name:         "MaxDimOld",
				Type:         "number",
				DefaultValue: "1500",
			},
		},
		Transformations: []imaging.TransformationType{
			&imaging.RegionOfInterestCrop{
				Transformation: "RegionOfInterestCrop",
				Style:          &imaging.RegionOfInterestCropStyleVariableInline{Value: imaging.RegionOfInterestCropStylePtr("fill")},
				Gravity:        &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
				Width:          &imaging.IntegerVariableInline{Value: ptr.To(7)},
				Height:         &imaging.IntegerVariableInline{Value: ptr.To(8)},
				RegionOfInterest: &imaging.RectangleShapeType{
					Anchor: &imaging.PointShapeType{
						X: &imaging.NumberVariableInline{Value: ptr.To(float64(4))},
						Y: &imaging.NumberVariableInline{Value: ptr.To(float64(5))},
					},
					Width:  &imaging.NumberVariableInline{Value: ptr.To(float64(8))},
					Height: &imaging.NumberVariableInline{Value: ptr.To(float64(9))},
				},
			},
			&imaging.Append{
				Transformation:         "Append",
				Gravity:                &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
				GravityPriority:        &imaging.AppendGravityPriorityVariableInline{Value: imaging.AppendGravityPriorityPtr("horizontal")},
				PreserveMinorDimension: &imaging.BooleanVariableInline{Value: ptr.To(true)},
				Image: &imaging.TextImageType{
					Type:       "Text",
					Fill:       &imaging.StringVariableInline{Value: ptr.To("#000000")},
					Size:       &imaging.NumberVariableInline{Value: ptr.To(float64(72))},
					Stroke:     &imaging.StringVariableInline{Value: ptr.To("#FFFFFF")},
					StrokeSize: &imaging.NumberVariableInline{Value: ptr.To(float64(0))},
					Text:       &imaging.StringVariableInline{Value: ptr.To("test")},
					Transformation: &imaging.Compound{
						Transformation: "Compound",
					},
				},
			},
			&imaging.Trim{
				Transformation: "Trim",
				Fuzz: &imaging.NumberVariableInline{
					Value: ptr.To(0.08),
				},
				Padding: &imaging.IntegerVariableInline{
					Value: ptr.To(0),
				},
			},
			&imaging.IfDimension{
				Transformation: "IfDimension",
				Dimension: &imaging.IfDimensionDimensionVariableInline{
					Value: imaging.IfDimensionDimensionPtr("width"),
				},
				Value: &imaging.IntegerVariableInline{
					Name: ptr.To("MaxDimOld"),
				},
				Default: &imaging.Compound{
					Transformation: "Compound",
					Transformations: []imaging.TransformationType{
						&imaging.IfDimension{
							Transformation: "IfDimension",
							Dimension: &imaging.IfDimensionDimensionVariableInline{
								Value: imaging.IfDimensionDimensionPtr("width"),
							},
							Value: &imaging.IntegerVariableInline{
								Name: ptr.To("MinDim"),
							},
							LessThan: &imaging.Compound{
								Transformation: "Compound",
								Transformations: []imaging.TransformationType{
									&imaging.Resize{
										Transformation: "Resize",
										Aspect: &imaging.ResizeAspectVariableInline{
											Value: imaging.ResizeAspectPtr("fit"),
										},
										Type: &imaging.ResizeTypeVariableInline{
											Value: imaging.ResizeTypePtr("normal"),
										},
										Width: &imaging.IntegerVariableInline{
											Name: ptr.To("ResizeDimWithBorder"),
										},
										Height: &imaging.IntegerVariableInline{
											Name: ptr.To("ResizeDimWithBorder"),
										},
									},
									&imaging.Crop{
										Transformation: "Crop",
										XPosition: &imaging.IntegerVariableInline{
											Value: ptr.To(0),
										},
										YPosition: &imaging.IntegerVariableInline{
											Value: ptr.To(0),
										},
										Gravity: &imaging.GravityVariableInline{
											Value: imaging.GravityPtr("Center"),
										},
										AllowExpansion: &imaging.BooleanVariableInline{
											Value: ptr.To(true),
										},
										Width: &imaging.IntegerVariableInline{
											Name: ptr.To("ResizeDim"),
										},
										Height: &imaging.IntegerVariableInline{
											Name: ptr.To("ResizeDim"),
										},
									},
									&imaging.BackgroundColor{
										Transformation: "BackgroundColor",
										Color: &imaging.StringVariableInline{
											Value: ptr.To("#ffffff"),
										},
									},
								},
							},
							Default: &imaging.Compound{
								Transformation: "Compound",
								Transformations: []imaging.TransformationType{
									&imaging.IfDimension{
										Transformation: "IfDimension",
										Dimension: &imaging.IfDimensionDimensionVariableInline{
											Value: imaging.IfDimensionDimensionPtr("height"),
										},
										Value: &imaging.IntegerVariableInline{
											Name: ptr.To("MinDim"),
										},
										LessThan: &imaging.Compound{
											Transformation: "Compound",
											Transformations: []imaging.TransformationType{
												&imaging.Resize{
													Transformation: "Resize",
													Aspect: &imaging.ResizeAspectVariableInline{
														Value: imaging.ResizeAspectPtr("fit"),
													},
													Type: &imaging.ResizeTypeVariableInline{
														Value: imaging.ResizeTypePtr("normal"),
													},
													Width: &imaging.IntegerVariableInline{
														Name: ptr.To("ResizeDimWithBorder"),
													},
													Height: &imaging.IntegerVariableInline{
														Name: ptr.To("ResizeDimWithBorder"),
													},
												},
												&imaging.Crop{
													Transformation: "Crop",
													XPosition: &imaging.IntegerVariableInline{
														Value: ptr.To(0),
													},
													YPosition: &imaging.IntegerVariableInline{
														Value: ptr.To(0),
													},
													Gravity: &imaging.GravityVariableInline{
														Value: imaging.GravityPtr("Center"),
													},
													AllowExpansion: &imaging.BooleanVariableInline{
														Value: ptr.To(true),
													},
													Width: &imaging.IntegerVariableInline{
														Name: ptr.To("ResizeDim"),
													},
													Height: &imaging.IntegerVariableInline{
														Name: ptr.To("ResizeDim"),
													},
												},
												&imaging.BackgroundColor{
													Transformation: "BackgroundColor",
													Color: &imaging.StringVariableInline{
														Value: ptr.To("#ffffff"),
													},
												},
											},
										},
										Default: &imaging.Compound{
											Transformation: "Compound",
											Transformations: []imaging.TransformationType{
												&imaging.IfDimension{
													Transformation: "IfDimension",
													Dimension: &imaging.IfDimensionDimensionVariableInline{
														Value: imaging.IfDimensionDimensionPtr("height"),
													},
													Value: &imaging.IntegerVariableInline{
														Name: ptr.To("MaxDimOld"),
													},
													GreaterThan: &imaging.Compound{
														Transformation: "Compound",
														Transformations: []imaging.TransformationType{
															&imaging.Resize{
																Transformation: "Resize",
																Aspect: &imaging.ResizeAspectVariableInline{
																	Value: imaging.ResizeAspectPtr("fit"),
																},
																Type: &imaging.ResizeTypeVariableInline{
																	Value: imaging.ResizeTypePtr("normal"),
																},

																Width: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDimWithBorder"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDimWithBorder"),
																},
															},
															&imaging.Crop{
																Transformation: "Crop",
																XPosition: &imaging.IntegerVariableInline{
																	Value: ptr.To(0),
																},
																YPosition: &imaging.IntegerVariableInline{
																	Value: ptr.To(0),
																},
																Gravity: &imaging.GravityVariableInline{
																	Value: imaging.GravityPtr("Center"),
																},
																AllowExpansion: &imaging.BooleanVariableInline{
																	Value: ptr.To(true),
																},
																Width: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
															},
															&imaging.BackgroundColor{
																Transformation: "BackgroundColor",
																Color: &imaging.StringVariableInline{
																	Value: ptr.To("#ffffff"),
																},
															},
														},
													},
													Default: &imaging.Compound{
														Transformation: "Compound",
														Transformations: []imaging.TransformationType{
															&imaging.Resize{
																Transformation: "Resize",
																Aspect: &imaging.ResizeAspectVariableInline{
																	Value: imaging.ResizeAspectPtr("fit"),
																},
																Type: &imaging.ResizeTypeVariableInline{
																	Value: imaging.ResizeTypePtr("normal"),
																},
																Width: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
															},
															&imaging.Crop{
																Transformation: "Crop",
																XPosition: &imaging.IntegerVariableInline{
																	Value: ptr.To(0),
																},
																YPosition: &imaging.IntegerVariableInline{
																	Value: ptr.To(0),
																},
																Gravity: &imaging.GravityVariableInline{
																	Value: imaging.GravityPtr("Center"),
																},
																AllowExpansion: &imaging.BooleanVariableInline{
																	Value: ptr.To(true),
																},
																Width: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: ptr.To("ResizeDim"),
																},
															},
															&imaging.BackgroundColor{
																Transformation: "BackgroundColor",
																Color: &imaging.StringVariableInline{
																	Value: ptr.To("#ffffff"),
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		PostBreakpointTransformations: []imaging.TransformationTypePost{
			&imaging.BackgroundColor{
				Transformation: "BackgroundColor",
				Color: &imaging.StringVariableInline{
					Value: ptr.To("#ffffff"),
				},
			},
			&imaging.IfDimensionPost{
				Transformation: "IfDimension",
				Dimension: &imaging.IfDimensionPostDimensionVariableInline{
					Value: imaging.IfDimensionPostDimensionPtr("height"),
				},
				Value: &imaging.IntegerVariableInline{
					Name: ptr.To("MaxDimOld"),
				},
				GreaterThan: &imaging.CompoundPost{
					Transformation: "Compound",
					Transformations: []imaging.TransformationTypePost{
						&imaging.BackgroundColor{
							Transformation: "BackgroundColor",
							Color: &imaging.StringVariableInline{
								Value: ptr.To("#ffffff"),
							},
						},
					},
				},
				Default: &imaging.CompoundPost{
					Transformation: "Compound",
					Transformations: []imaging.TransformationTypePost{
						&imaging.BackgroundColor{
							Transformation: "BackgroundColor",
							Color: &imaging.StringVariableInline{
								Value: ptr.To("#ffffff"),
							},
						},
					},
				},
			},
			&imaging.CompositePost{
				Gravity: &imaging.GravityPostVariableInline{Value: imaging.GravityPostPtr("NorthWest")},
				Image: &imaging.TextImageTypePost{
					Fill:       &imaging.StringVariableInline{Value: ptr.To("#000000")},
					Size:       &imaging.NumberVariableInline{Value: ptr.To(float64(72))},
					Stroke:     &imaging.StringVariableInline{Value: ptr.To("#FFFFFF")},
					StrokeSize: &imaging.NumberVariableInline{Value: ptr.To(float64(0))},
					Text:       &imaging.StringVariableInline{Value: ptr.To("test")},
					Type:       imaging.TextImageTypePostTypeText,
					Transformation: &imaging.CompoundPost{
						Transformation: imaging.CompoundPostTransformationCompound,
					},
				},
				Placement:      &imaging.CompositePostPlacementVariableInline{Value: imaging.CompositePostPlacementPtr(imaging.CompositePostPlacementOver)},
				Transformation: imaging.CompositePostTransformationComposite,
				XPosition:      &imaging.IntegerVariableInline{Value: ptr.To(0)},
				YPosition:      &imaging.IntegerVariableInline{Value: ptr.To(0)},
			},
		},
		Breakpoints: &imaging.Breakpoints{
			Widths: []int{280, 1080},
		},
		Output: &imaging.OutputImage{
			PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
				Value: imaging.OutputImagePerceptualQualityPtr("mediumHigh"),
			},
			AdaptiveQuality: ptr.To(50),
		},
	}

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
						Value: ptr.To(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   ptr.To(false),
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
						Value: ptr.To(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 2,
			Video:   ptr.To(false),
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
						Value: ptr.To(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   ptr.To(false),
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
			Video:   ptr.To(true),
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
			Video:   ptr.To(true),
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
			Video:   ptr.To(true),
		},
	}

	processor = func(testdir string) templates.FSTemplateProcessor {
		return templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"imaging.tmpl":   fmt.Sprintf("./testdata/res/%s/imaging.tf", testdir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", testdir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", testdir),
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
	}

	expectGetPolicySet = func(i *imaging.Mock, policySetID, contractID, name, policyType string, region imaging.Region,
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

	expectGetPolicy = func(i *imaging.Mock, policyRequest imaging.GetPolicyRequest, policyOutput imaging.PolicyOutput, err error) *mock.Call {
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

	expectListPolicies = func(i *imaging.Mock, policySetID, contractID, itemKind string, network imaging.PolicyNetwork,
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
	tests := map[string]struct {
		edgercPath   string
		section      string
		init         func(*imaging.Mock)
		filesToCheck []string
		dataDir      string
		jsonDir      string
		withError    error
		policyAsHCL  bool
	}{
		"fetch policy set with given id and contract and no policies": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns zero policies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					nil, 0, nil).Once()
			},
			dataDir:      "json/image_no_policies",
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with image policies same on production": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[1], nil)
			},
			dataDir:      "json/image_policies",
			filesToCheck: []string{"_auto.json", "test_policy_image.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with image policies same on production with jsondir": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[1], nil)
			},
			dataDir:      "json/image_policies_jsondir",
			jsonDir:      "jsondir",
			filesToCheck: []string{"jsondir/_auto.json", "jsondir/test_policy_image.json", "imaging.tf", "import.sh", "variables.tf"},
		},

		"fetch policy set with image policies same on production as as hcl": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[1], nil)
			},
			dataDir:      "json/image_policies_as_hcl",
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
			policyAsHCL:  true,
		},
		"fetch policy set with image policies same on production as as hcl, too many levels": {
			init: func(i *imaging.Mock) {
				policy, err := convertPolicyInputImage(&veryDeepPolicy)
				require.NoError(t, err)
				policy.ID = ".auto"

				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					[]imaging.PolicyOutput{policy}, 1, nil).Once()
				expectGetPolicy(i, imaging.GetPolicyRequest{
					PolicyID:    ".auto",
					Network:     imaging.PolicyNetworkProduction,
					ContractID:  "ctr_123",
					PolicySetID: "test_policyset_id",
				}, policy, nil)
			},
			withError:   ErrFetchingPolicy,
			policyAsHCL: true,
		},
		"fetch policy set with image policies different on production": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesImageData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], imagePoliciesOutputs[0], nil)
				// policyID: test_policy_image - different on production
				expectGetPolicy(i, policiesRequests[1], imagePoliciesOutputs[2], nil)
			},
			dataDir:      "json/image_policies_diff_prod",
			filesToCheck: []string{"_auto.json", "test_policy_image.json", "imaging.tf", "import.sh", "variables.tf"},
		},

		"fetch policy set with video policies same on production": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[1], nil)
			},
			dataDir:      "json/video_policies",
			filesToCheck: []string{"_auto.json", "test_policy_video.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with video policies same on production with jsondir": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[1], nil)
			},
			dataDir:      "json/video_policies_jsondir",
			jsonDir:      "jsondir",
			filesToCheck: []string{"jsondir/_auto.json", "jsondir/test_policy_video.json", "imaging.tf", "import.sh", "variables.tf"},
		},

		"fetch policy set with video policies same on production as as hcl": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - same on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[1], nil)
			},
			dataDir:      "json/video_policies_as_hcl",
			policyAsHCL:  true,
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with video policies different on production": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", nil).Once()
				// getPolicies returns two policy outputs from staging
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					videoPoliciesOutputs[:2], 2, nil).Once()
				// getPoliciesVideoData
				// policyID: .auto - same on production
				expectGetPolicy(i, policiesRequests[0], videoPoliciesOutputs[0], nil)
				// policyID: test_policy_image - different on production
				expectGetPolicy(i, policiesRequests[2], videoPoliciesOutputs[2], nil)
			},
			dataDir:      "json/video_policies_diff_prod",
			filesToCheck: []string{"_auto.json", "test_policy_video.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"error fetching policy set": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicySet,
		},
		"error fetching policies": {
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs, 0, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicy,
		},
		"error fetching policy": {
			init: func(i *imaging.Mock) {
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
		"non default edgerc path and section": {
			edgercPath: "/non/default/path/to/edgerc",
			section:    "non_default_section",
			init: func(i *imaging.Mock) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns zero policies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					nil, 0, nil).Once()
			},
			dataDir:      "non_default_edgerc_path_and_section",
			filesToCheck: []string{"variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.edgercPath == "" {
				test.edgercPath = "~/.edgerc"
			}
			if test.section == "" {
				test.section = "test_section"
			}
			require.NoError(t, os.MkdirAll(fmt.Sprintf("testdata/res/%s/%s", test.dataDir, test.jsonDir), 0755))
			tfWorkPath := fmt.Sprintf("testdata/res/%s", test.dataDir)
			mi := new(imaging.Mock)
			mp := processor(test.dataDir)
			test.init(mi)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createImaging(ctx, "ctr_123", "test_policyset_id", tfWorkPath, test.jsonDir, test.edgercPath, test.section, mi, mp, test.policyAsHCL)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

			if test.filesToCheck != nil {
				for _, f := range test.filesToCheck {
					expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dataDir, f))
					require.NoError(t, err)
					result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dataDir, f))
					require.NoError(t, err)
					assert.Equal(t, string(expected), string(result))
				}
			}
			mi.AssertExpectations(t)
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
			},
			dir:          "with_image_policies",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies as hcl": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Policies: []TFPolicy{
					{
						PolicyID:             "test_policy_image",
						ActivateOnProduction: true,
						Policy:               &veryDeepPolicy,
					},
				},
			},
			dir:          "with_image_policies_as_hcl",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies as hcl empty": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Policies: []TFPolicy{
					{
						PolicyID:             "test_policy_image",
						ActivateOnProduction: true,
						Policy:               &imaging.PolicyInputImage{},
					},
				},
			},
			dir:          "with_image_policies_as_hcl_empty",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies as hcl with image type": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				Policies: []TFPolicy{
					{
						PolicyID:             "test_policy_image",
						ActivateOnProduction: true,
						Policy: &imaging.PolicyInputImage{
							Output: &imaging.OutputImage{
								AllowPristineOnDownsize: ptr.To(true),
								PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
									Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
								},
								PreferModernFormats: ptr.To(false),
							},
							Transformations: []imaging.TransformationType{
								&imaging.Append{
									Gravity:         &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
									GravityPriority: &imaging.AppendGravityPriorityVariableInline{Value: imaging.AppendGravityPriorityPtr("horizontal")},
									Image: &imaging.BoxImageType{
										Transformation: &imaging.Compound{
											Transformation: "Compound",
											Transformations: []imaging.TransformationType{
												&imaging.Append{
													Gravity:         &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
													GravityPriority: &imaging.AppendGravityPriorityVariableInline{Value: imaging.AppendGravityPriorityPtr("horizontal")},
													Image: &imaging.BoxImageType{
														Transformation: &imaging.Compound{
															Transformation: "Compound",
															Transformations: []imaging.TransformationType{
																&imaging.Append{
																	Gravity:         &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
																	GravityPriority: &imaging.AppendGravityPriorityVariableInline{Value: imaging.AppendGravityPriorityPtr("horizontal")},
																	Image: &imaging.BoxImageType{
																		Transformation: &imaging.Compound{},
																		Type:           "Box",
																	},
																	PreserveMinorDimension: &imaging.BooleanVariableInline{Value: ptr.To(false)},
																	Transformation:         "Append",
																},
															},
														},
														Type: "Box",
													},
													PreserveMinorDimension: &imaging.BooleanVariableInline{Value: ptr.To(false)},
													Transformation:         "Append",
												},
											},
										},
										Type: "Box",
									},
									PreserveMinorDimension: &imaging.BooleanVariableInline{Value: ptr.To(false)},
									Transformation:         "Append",
								},
							},
							ServeStaleDuration: ptr.To(3600),
						},
					},
				},
			},
			dir:          "with_image_policies_as_hcl_with_imagetype",
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
			},
			dir:          "with_video_policies",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with video policies as hcl": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "VIDEO",
				},
				Policies: []TFPolicy{
					{
						PolicyID:             "test_policy_video",
						ActivateOnProduction: true,
						Policy: &imaging.PolicyInputVideo{
							RolloutDuration: ptr.To(3600),
							Hosts:           []string{"host1", "host2"},
							Variables: []imaging.Variable{
								{
									Name:         "ResizeDim",
									Type:         "number",
									DefaultValue: "280",
								},
								{
									Name:         "ResizeDimWithBorder",
									Type:         "number",
									DefaultValue: "260",
								},
								{
									Name:         "MinDim",
									Type:         "number",
									DefaultValue: "1000",
									EnumOptions: []*imaging.EnumOptions{
										{
											ID:    "1",
											Value: "value1",
										},
										{
											ID:    "2",
											Value: "value2",
										},
									},
								},
								{
									Name:         "MinDimNew",
									Type:         "number",
									DefaultValue: "1450",
								},
								{
									Name:         "MaxDimOld",
									Type:         "number",
									DefaultValue: "1500",
								},
							},
							Breakpoints: &imaging.Breakpoints{
								Widths: []int{280, 1080},
							},
							Output: &imaging.OutputVideo{
								PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
									Value: imaging.OutputVideoPerceptualQualityPtr("mediumHigh"),
								},
								PlaceholderVideoURL: &imaging.StringVariableInline{
									Value: ptr.To("some"),
								},
							},
						},
					},
				},
			},
			dir:          "with_video_policies_as_hcl",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"non default edgerc path and section": {
			givenData: TFImagingData{
				PolicySet: TFPolicySet{
					ID:         "test_policyset_id",
					ContractID: "ctr_123",
					Name:       "some policy set",
					Region:     "EMEA",
					Type:       "IMAGE",
				},
				EdgercPath: "/non/default/path/to/edgerc",
				Section:    "non_default_section",
			},
			dir:          "non_default_edgerc_path_and_section",
			filesToCheck: []string{"variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.givenData.EdgercPath == "" {
				test.givenData.EdgercPath = "~/.edgerc"
			}
			if test.givenData.Section == "" {
				test.givenData.Section = "test_section"
			}
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			require.NoError(t, processor(test.dir).ProcessTemplates(test.givenData))

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

func TestGetDepth(t *testing.T) {
	tests := map[string]struct {
		args  interface{}
		depth int
	}{
		"2 depth": {
			args:  imagePoliciesOutputs[0].(*imaging.PolicyOutputImage),
			depth: 2,
		},
		"14 depth": {
			args:  veryDeepPolicy,
			depth: 14,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.depth, getDepth(tt.args, 0), "getDepth(%v, %v)", tt.args, tt.depth)
		})
	}
}

func convertPolicyInputImage(policy imaging.PolicyInput) (*imaging.PolicyOutputImage, error) {
	var policyJSON []byte
	var err error
	if policyJSON, err = json.Marshal(policy); err != nil {
		return nil, err
	}

	var policyInput imaging.PolicyOutputImage
	if err := json.Unmarshal(policyJSON, &policyInput); err != nil {
		return nil, err
	}
	return &policyInput, nil
}

func TestEnsureDirExists(t *testing.T) {
	t.Run("no json dir specified", func(t *testing.T) {
		tfDir, err := os.MkdirTemp("", "tfworkpath")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, os.RemoveAll(tfDir)) }()
		jsonDirPath := path.Join(tfDir, ".")

		err = ensureDirExists(jsonDirPath)

		assert.NoError(t, err)
		assert.DirExists(t, jsonDirPath)
	})

	t.Run("json dir already exists", func(t *testing.T) {
		tfDir, err := os.MkdirTemp("", "tfworkpath")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, os.RemoveAll(tfDir)) }()
		jsonDirPath := path.Join(tfDir, "jsondir")

		// create a dir in path where json dir is expected
		err = os.MkdirAll(jsonDirPath, 0755)
		assert.NoError(t, err)

		err = ensureDirExists(jsonDirPath)

		assert.NoError(t, err)
		assert.DirExists(t, jsonDirPath)
	})

	t.Run("json dir does not exists", func(t *testing.T) {
		tfDir, err := os.MkdirTemp("", "tfworkpath")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, os.RemoveAll(tfDir)) }()
		jsonDirPath := path.Join(tfDir, "jsondir")

		err = ensureDirExists(jsonDirPath)

		assert.NoError(t, err)
		assert.DirExists(t, jsonDirPath)
	})

	t.Run("path exists but is not a dir", func(t *testing.T) {
		tfDir, err := os.MkdirTemp("", "tfworkpath")
		assert.NoError(t, err)
		defer func() { assert.NoError(t, os.RemoveAll(tfDir)) }()
		jsonDirPath := path.Join(tfDir, "jsondir")

		// create a file in path where json dir is expected
		err = os.MkdirAll(tfDir, 0755)
		assert.NoError(t, err)
		f, err := os.Create(jsonDirPath)
		assert.NoError(t, err)
		assert.NoError(t, f.Close())

		err = ensureDirExists(jsonDirPath)

		assert.ErrorIs(t, err, ErrCreateDir)
		assert.NoDirExists(t, jsonDirPath)
	})
}
