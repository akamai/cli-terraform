package imaging

import (
	"context"
	"encoding/json"
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
		RolloutDuration: 3600,
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
		Transformations: []imaging.TransformationType{
			&imaging.RegionOfInterestCrop{
				Transformation: "RegionOfInterestCrop",
				Style:          &imaging.RegionOfInterestCropStyleVariableInline{Value: imaging.RegionOfInterestCropStylePtr("fill")},
				Gravity:        &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
				Width:          &imaging.IntegerVariableInline{Value: tools.IntPtr(7)},
				Height:         &imaging.IntegerVariableInline{Value: tools.IntPtr(8)},
				RegionOfInterest: &imaging.RectangleShapeType{
					Anchor: &imaging.PointShapeType{
						X: &imaging.NumberVariableInline{Value: tools.Float64Ptr(4)},
						Y: &imaging.NumberVariableInline{Value: tools.Float64Ptr(5)},
					},
					Width:  &imaging.NumberVariableInline{Value: tools.Float64Ptr(8)},
					Height: &imaging.NumberVariableInline{Value: tools.Float64Ptr(9)},
				},
			},
			&imaging.Append{
				Transformation:         "Append",
				Gravity:                &imaging.GravityVariableInline{Value: imaging.GravityPtr("Center")},
				GravityPriority:        &imaging.AppendGravityPriorityVariableInline{Value: imaging.AppendGravityPriorityPtr("horizontal")},
				PreserveMinorDimension: &imaging.BooleanVariableInline{Value: tools.BoolPtr(true)},
				Image: &imaging.TextImageType{
					Type:       "Text",
					Fill:       &imaging.StringVariableInline{Value: tools.StringPtr("#000000")},
					Size:       &imaging.NumberVariableInline{Value: tools.Float64Ptr(72)},
					Stroke:     &imaging.StringVariableInline{Value: tools.StringPtr("#FFFFFF")},
					StrokeSize: &imaging.NumberVariableInline{Value: tools.Float64Ptr(0)},
					Text:       &imaging.StringVariableInline{Value: tools.StringPtr("test")},
					Transformation: &imaging.Compound{
						Transformation: "Compound",
					},
				},
			},
			&imaging.Trim{
				Transformation: "Trim",
				Fuzz: &imaging.NumberVariableInline{
					Value: tools.Float64Ptr(0.08),
				},
				Padding: &imaging.IntegerVariableInline{
					Value: tools.IntPtr(0),
				},
			},
			&imaging.IfDimension{
				Transformation: "IfDimension",
				Dimension: &imaging.IfDimensionDimensionVariableInline{
					Value: imaging.IfDimensionDimensionPtr("width"),
				},
				Value: &imaging.IntegerVariableInline{
					Name: tools.StringPtr("MaxDimOld"),
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
								Name: tools.StringPtr("MinDim"),
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
											Name: tools.StringPtr("ResizeDimWithBorder"),
										},
										Height: &imaging.IntegerVariableInline{
											Name: tools.StringPtr("ResizeDimWithBorder"),
										},
									},
									&imaging.Crop{
										Transformation: "Crop",
										XPosition: &imaging.IntegerVariableInline{
											Value: tools.IntPtr(0),
										},
										YPosition: &imaging.IntegerVariableInline{
											Value: tools.IntPtr(0),
										},
										Gravity: &imaging.GravityVariableInline{
											Value: imaging.GravityPtr("Center"),
										},
										AllowExpansion: &imaging.BooleanVariableInline{
											Value: tools.BoolPtr(true),
										},
										Width: &imaging.IntegerVariableInline{
											Name: tools.StringPtr("ResizeDim"),
										},
										Height: &imaging.IntegerVariableInline{
											Name: tools.StringPtr("ResizeDim"),
										},
									},
									&imaging.BackgroundColor{
										Transformation: "BackgroundColor",
										Color: &imaging.StringVariableInline{
											Value: tools.StringPtr("#ffffff"),
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
											Name: tools.StringPtr("MinDim"),
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
														Name: tools.StringPtr("ResizeDimWithBorder"),
													},
													Height: &imaging.IntegerVariableInline{
														Name: tools.StringPtr("ResizeDimWithBorder"),
													},
												},
												&imaging.Crop{
													Transformation: "Crop",
													XPosition: &imaging.IntegerVariableInline{
														Value: tools.IntPtr(0),
													},
													YPosition: &imaging.IntegerVariableInline{
														Value: tools.IntPtr(0),
													},
													Gravity: &imaging.GravityVariableInline{
														Value: imaging.GravityPtr("Center"),
													},
													AllowExpansion: &imaging.BooleanVariableInline{
														Value: tools.BoolPtr(true),
													},
													Width: &imaging.IntegerVariableInline{
														Name: tools.StringPtr("ResizeDim"),
													},
													Height: &imaging.IntegerVariableInline{
														Name: tools.StringPtr("ResizeDim"),
													},
												},
												&imaging.BackgroundColor{
													Transformation: "BackgroundColor",
													Color: &imaging.StringVariableInline{
														Value: tools.StringPtr("#ffffff"),
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
														Name: tools.StringPtr("MaxDimOld"),
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
																	Name: tools.StringPtr("ResizeDimWithBorder"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDimWithBorder"),
																},
															},
															&imaging.Crop{
																Transformation: "Crop",
																XPosition: &imaging.IntegerVariableInline{
																	Value: tools.IntPtr(0),
																},
																YPosition: &imaging.IntegerVariableInline{
																	Value: tools.IntPtr(0),
																},
																Gravity: &imaging.GravityVariableInline{
																	Value: imaging.GravityPtr("Center"),
																},
																AllowExpansion: &imaging.BooleanVariableInline{
																	Value: tools.BoolPtr(true),
																},
																Width: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDim"),
																},
															},
															&imaging.BackgroundColor{
																Transformation: "BackgroundColor",
																Color: &imaging.StringVariableInline{
																	Value: tools.StringPtr("#ffffff"),
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
																	Name: tools.StringPtr("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDim"),
																},
															},
															&imaging.Crop{
																Transformation: "Crop",
																XPosition: &imaging.IntegerVariableInline{
																	Value: tools.IntPtr(0),
																},
																YPosition: &imaging.IntegerVariableInline{
																	Value: tools.IntPtr(0),
																},
																Gravity: &imaging.GravityVariableInline{
																	Value: imaging.GravityPtr("Center"),
																},
																AllowExpansion: &imaging.BooleanVariableInline{
																	Value: tools.BoolPtr(true),
																},
																Width: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDim"),
																},
																Height: &imaging.IntegerVariableInline{
																	Name: tools.StringPtr("ResizeDim"),
																},
															},
															&imaging.BackgroundColor{
																Transformation: "BackgroundColor",
																Color: &imaging.StringVariableInline{
																	Value: tools.StringPtr("#ffffff"),
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
		PostBreakpointTransformations: []imaging.TransformationType{
			&imaging.BackgroundColor{
				Transformation: "BackgroundColor",
				Color: &imaging.StringVariableInline{
					Value: tools.StringPtr("#ffffff"),
				},
			},
			&imaging.IfDimension{
				Transformation: "IfDimension",
				Dimension: &imaging.IfDimensionDimensionVariableInline{
					Value: imaging.IfDimensionDimensionPtr("height"),
				},
				Value: &imaging.IntegerVariableInline{
					Name: tools.StringPtr("MaxDimOld"),
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
								Name: tools.StringPtr("ResizeDimWithBorder"),
							},
							Height: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDimWithBorder"),
							},
						},
						&imaging.Crop{
							Transformation: "Crop",
							XPosition: &imaging.IntegerVariableInline{
								Value: tools.IntPtr(0),
							},
							YPosition: &imaging.IntegerVariableInline{
								Value: tools.IntPtr(0),
							},
							Gravity: &imaging.GravityVariableInline{
								Value: imaging.GravityPtr("Center"),
							},
							AllowExpansion: &imaging.BooleanVariableInline{
								Value: tools.BoolPtr(true),
							},
							Width: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDim"),
							},
							Height: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDim"),
							},
						},
						&imaging.BackgroundColor{
							Transformation: "BackgroundColor",
							Color: &imaging.StringVariableInline{
								Value: tools.StringPtr("#ffffff"),
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
								Name: tools.StringPtr("ResizeDim"),
							},
							Height: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDim"),
							},
						},
						&imaging.Crop{
							Transformation: "Crop",
							XPosition: &imaging.IntegerVariableInline{
								Value: tools.IntPtr(0),
							},
							YPosition: &imaging.IntegerVariableInline{
								Value: tools.IntPtr(0),
							},
							Gravity: &imaging.GravityVariableInline{
								Value: imaging.GravityPtr("Center"),
							},
							AllowExpansion: &imaging.BooleanVariableInline{
								Value: tools.BoolPtr(true),
							},
							Width: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDim"),
							},
							Height: &imaging.IntegerVariableInline{
								Name: tools.StringPtr("ResizeDim"),
							},
						},
						&imaging.BackgroundColor{
							Transformation: "BackgroundColor",
							Color: &imaging.StringVariableInline{
								Value: tools.StringPtr("#ffffff"),
							},
						},
					},
				},
			},
		},
		Breakpoints: &imaging.Breakpoints{
			Widths: []int{280, 1080},
		},
		Output: &imaging.OutputImage{
			PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
				Value: imaging.OutputImagePerceptualQualityPtr("mediumHigh"),
			},
			AdaptiveQuality: 50,
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
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   tools.BoolPtr(false),
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
			Video:   tools.BoolPtr(false),
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
			Video:   tools.BoolPtr(false),
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
			Video:   tools.BoolPtr(true),
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
			Video:   tools.BoolPtr(true),
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
			Video:   tools.BoolPtr(true),
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
		init         func(*mockimaging)
		filesToCheck []string
		jsonDir      string
		withError    error
		schema       bool
	}{
		"fetch policy set with given id and contract and no policies": {
			init: func(i *mockimaging) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies returns zero policies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					nil, 0, nil).Once()
			},
			jsonDir:      "json/image_no_policies",
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with image policies same on production": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/image_policies",
			filesToCheck: []string{"_auto.json", "test_policy_image.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with image policies same on production as schema": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/image_policies_schema",
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
			schema:       true,
		},
		"fetch policy set with image policies same on production as schema, too many levels": {
			init: func(i *mockimaging) {
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
			withError: ErrFetchingPolicy,
			schema:    true,
		},
		"fetch policy set with image policies different on production": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/image_policies_diff_prod",
			filesToCheck: []string{"_auto.json", "test_policy_image.json", "imaging.tf", "import.sh", "variables.tf"},
		},

		"fetch policy set with video policies same on production": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/video_policies",
			filesToCheck: []string{"_auto.json", "test_policy_video.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with video policies same on production as schema": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/video_policies_schema",
			schema:       true,
			filesToCheck: []string{"imaging.tf", "import.sh", "variables.tf"},
		},
		"fetch policy set with video policies different on production": {
			init: func(i *mockimaging) {
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
			jsonDir:      "json/video_policies_diff_prod",
			filesToCheck: []string{"_auto.json", "test_policy_video.json", "imaging.tf", "import.sh", "variables.tf"},
		},
		"error fetching policy set": {
			init: func(i *mockimaging) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "VIDEO", "EMEA", fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicySet,
		},
		"error fetching policies": {
			init: func(i *mockimaging) {
				expectGetPolicySet(i, "test_policyset_id", "ctr_123", "some policy set", "IMAGE", "EMEA", nil).Once()
				// getPolicies
				expectListPolicies(i, "test_policyset_id", "ctr_123", "POLICY", imaging.PolicyNetworkStaging,
					imagePoliciesOutputs, 0, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingPolicy,
		},
		"error fetching policy": {
			init: func(i *mockimaging) {
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("testdata/res/%s", test.jsonDir), 0755))
			mi := new(mockimaging)
			mp := processor(test.jsonDir)
			test.init(mi)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createImaging(ctx, "ctr_123", "test_policyset_id", fmt.Sprintf("testdata/res/%s", test.jsonDir), section, mi, mp, test.schema)
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
		"policy set with image policies schema": {
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
				Section: "test_section",
			},
			dir:          "with_image_policies_schema",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies schema empty": {
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
				Section: "test_section",
			},
			dir:          "with_image_policies_schema_empty",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
		"policy set with image policies schema with image type": {
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
																	PreserveMinorDimension: &imaging.BooleanVariableInline{Value: tools.BoolPtr(false)},
																	Transformation:         "Append",
																},
															},
														},
														Type: "Box",
													},
													PreserveMinorDimension: &imaging.BooleanVariableInline{Value: tools.BoolPtr(false)},
													Transformation:         "Append",
												},
											},
										},
										Type: "Box",
									},
									PreserveMinorDimension: &imaging.BooleanVariableInline{Value: tools.BoolPtr(false)},
									Transformation:         "Append",
								},
							},
						},
					},
				},
				Section: "test_section",
			},
			dir:          "with_image_policies_schema_with_imagetype",
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
		"policy set with video policies schema": {
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
							RolloutDuration: 3600,
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
									Value: tools.StringPtr("some"),
								},
							},
						},
					},
				},
				Section: "test_section",
			},
			dir:          "with_video_policies_schema",
			filesToCheck: []string{"imaging.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			require.NoError(t, processor(test.dir).ProcessTemplates(test.givenData))

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
