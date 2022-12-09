package papi

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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
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

func TestCreateProperty(t *testing.T) {
	section := "test_section"

	searchPropertiesResponse := papi.SearchResponse{
		Versions: papi.SearchItems{
			Items: []papi.SearchItem{
				{
					AccountID:        "test_account",
					AssetID:          "aid_10541511",
					ContractID:       "test_contract",
					EdgeHostname:     "",
					GroupID:          "grp_12345",
					Hostname:         "",
					ProductionStatus: "ACTIVE",
					PropertyID:       "prp_12345",
					PropertyName:     "test.edgesuite.net",
					PropertyVersion:  2,
					StagingStatus:    "ACTIVE",
					UpdatedByUser:    "jsmith",
					UpdatedDate:      "2018-03-07T23:40:45Z",
				},
				{
					AccountID:        "test_account",
					AssetID:          "aid_10541511",
					ContractID:       "test_contract",
					EdgeHostname:     "",
					GroupID:          "grp_12345",
					Hostname:         "",
					ProductionStatus: "INACTIVE",
					PropertyID:       "prp_12345",
					PropertyName:     "test.edgesuite.net",
					PropertyVersion:  5,
					StagingStatus:    "INACTIVE",
					UpdatedByUser:    "jsmith",
					UpdatedDate:      "2019-04-24T18:54:03Z",
				},
			},
		},
	}
	getPropertyResponse := papi.GetPropertyResponse{
		Properties: papi.PropertiesItems{
			Items: []*papi.Property{
				{
					AccountID:         "test_account",
					AssetID:           "aid_10541511",
					ContractID:        "test_contract",
					GroupID:           "grp_12345",
					LatestVersion:     5,
					Note:              "",
					ProductID:         "prd_HTTP_Content_Del",
					ProductionVersion: nil,
					PropertyID:        "prp_12345",
					PropertyName:      "test.edgesuite.net",
					RuleFormat:        "latest",
					StagingVersion:    nil,
				},
			},
		},
		Property: &papi.Property{
			AccountID:         "test_account",
			AssetID:           "aid_10541511",
			ContractID:        "test_contract",
			GroupID:           "grp_12345",
			LatestVersion:     5,
			Note:              "",
			ProductID:         "prd_HTTP_Content_Del",
			ProductionVersion: nil,
			PropertyID:        "prp_12345",
			PropertyName:      "test.edgesuite.net",
			RuleFormat:        "latest",
			StagingVersion:    nil,
		},
	}
	getGroupsResponse := papi.GetGroupsResponse{
		AccountID:   "test_account",
		AccountName: "Test Account",
		Groups: papi.GroupItems{Items: []*papi.Group{
			{
				GroupID:       "grp_12345",
				GroupName:     "test_group",
				ParentGroupID: "grp_12345",
				ContractIDs:   nil,
			},
		}}}

	getPropertyVersionsResponse := papi.GetPropertyVersionsResponse{
		PropertyID:   "prp_12345",
		PropertyName: "test.edgesuite.net",
		AccountID:    "test_account",
		ContractID:   "test_contract",
		GroupID:      "grp_12345",
		AssetID:      "aid_10541511",
		Versions: papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					Etag:             "4607f363da8bc05b0c0f0f75249",
					Note:             "",
					ProductID:        "prd_HTTP_Content_Del",
					ProductionStatus: "INACTIVE",
					PropertyVersion:  1,
					RuleFormat:       "latest",
					StagingStatus:    "INACTIVE",
					UpdatedByUser:    "jsmith",
					UpdatedDate:      "2019-04-23T18:54:03Z",
				},
				{
					Etag:             "4607f363da8bc05b0c0f0f75249",
					Note:             "",
					ProductID:        "prd_HTTP_Content_Del",
					ProductionStatus: "INACTIVE",
					PropertyVersion:  5,
					RuleFormat:       "latest",
					StagingStatus:    "INACTIVE",
					UpdatedByUser:    "jsmith",
					UpdatedDate:      "2019-04-24T18:54:03Z",
				},
			},
		},
	}

	getLatestVersionResponse := papi.GetPropertyVersionsResponse{
		PropertyID:   "prp_12345",
		PropertyName: "test.edgesuite.net",
		AccountID:    "test_account", ContractID: "test_contract",
		GroupID: "grp_12345",
		AssetID: "aid_10541511",
		Versions: papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					Etag:             "4607f363da8bc05b0c0f0f75249",
					Note:             "",
					ProductID:        "prd_HTTP_Content_Del",
					ProductionStatus: "INACTIVE",
					PropertyVersion:  5,
					RuleFormat:       "latest",
					StagingStatus:    "INACTIVE",
					UpdatedByUser:    "jsmith",
					UpdatedDate:      "2019-04-24T18:54:03Z",
				},
			},
		},
		Version: papi.PropertyVersionGetItem{
			Etag:             "4607f363da8bc05b0c0f0f75249",
			Note:             "",
			ProductID:        "prd_HTTP_Content_Del",
			ProductionStatus: "INACTIVE",
			PropertyVersion:  5,
			RuleFormat:       "latest",
			StagingStatus:    "INACTIVE",
			UpdatedByUser:    "jsmith",
			UpdatedDate:      "2019-04-24T18:54:03Z",
		},
	}

	getListReferencedIncludesResponse := papi.ListReferencedIncludesResponse{
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
					StagingVersion:    tools.IntPtr(1),
					ProductionVersion: tools.IntPtr(1),
				},
			},
		},
	}

	getProductsResponse := papi.GetProductsResponse{
		AccountID:  "test_account",
		ContractID: "test_contract",
		Products: papi.ProductsItems{
			Items: []papi.ProductItem{
				{
					ProductName: "Enterprise",
					ProductID:   "prd_Enterprise",
				},
				{
					ProductName: "Progressive_Media",
					ProductID:   "prd_Progressive_Media",
				},
				{
					ProductName: "Site_Defender",
					ProductID:   "prd_Site_Defender",
				},
				{
					ProductName: "Site_Del",
					ProductID:   "prd_Site_Del",
				},
				{
					ProductName: "Mobile_Accel",
					ProductID:   "prd_Mobile_Accel",
				},
				{
					ProductName: "Web_App_Accel",
					ProductID:   "prd_Web_App_Accel",
				},
				{
					ProductName: "Adaptive_Media_Delivery",
					ProductID:   "prd_Adaptive_Media_Delivery",
				},
				{
					ProductName: "HTTP_Content_Del",
					ProductID:   "prd_HTTP_Content_Del",
				},
			},
		},
	}

	getPropertyVersionHostnamesResponse := papi.GetPropertyVersionHostnamesResponse{
		AccountID:       "test_account",
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
		PropertyID:      "prp_12345",
		PropertyVersion: 5,
		Etag:            "4607f363da8bc05b0c0f0f7524985d2fbc5d864d",
		Hostnames: papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_2867480",
					CnameFrom:            "test.edgesuite.net",
					CnameTo:              "test.edgesuite.net",
					CertProvisioningType: "CPS_MANAGED",
				},
			},
		},
	}

	getPropertyVersion1HostnamesResponse := papi.GetPropertyVersionHostnamesResponse{
		AccountID:       "test_account",
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
		PropertyID:      "prp_12345",
		PropertyVersion: 1,
		Etag:            "4607f363da8bc05b0c0f0f7524985d2fbc5d864d",
		Hostnames: papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_2867480",
					CnameFrom:            "test.edgesuite.net",
					CnameTo:              "test.edgesuite.net",
					CertProvisioningType: "CPS_MANAGED",
				},
			},
		},
	}

	getPropertyVersion2HostnamesResponse := papi.GetPropertyVersionHostnamesResponse{
		AccountID:       "test_account",
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
		PropertyID:      "prp_12345",
		PropertyVersion: 5,
		Etag:            "4607f363da8bc05b0c0f0f7524985d2fbc5d864d",
		Hostnames: papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_2867480",
					CnameFrom:            "test.edgesuite.net",
					CnameTo:              "test.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				},
			},
		},
	}

	getActivationsResponse := papi.GetActivationsResponse{
		Response: papi.Response{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "grp_12345",
		},
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				{
					AccountID:              "",
					ActivationID:           "atv_5594260",
					ActivationType:         "ACTIVATE",
					UseFastFallback:        false,
					FallbackInfo:           nil,
					AcknowledgeWarnings:    nil,
					AcknowledgeAllWarnings: false,
					FastPush:               false,
					FMAActivationState:     "",
					GroupID:                "",
					IgnoreHTTPErrors:       false,
					PropertyName:           "test.edgesuite.net",
					PropertyID:             "prp_12345",
					PropertyVersion:        2,
					Network:                "STAGING",
					Status:                 "ACTIVE",
					NotifyEmails:           []string{"jsmith@akamai.com"},
				},
			},
		},
	}

	getActivations1Response := papi.GetActivationsResponse{
		Response: papi.Response{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "grp_12345",
		},
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				{
					AccountID:              "",
					ActivationID:           "atv_5594260",
					ActivationType:         "ACTIVATE",
					UseFastFallback:        false,
					FallbackInfo:           nil,
					AcknowledgeWarnings:    nil,
					AcknowledgeAllWarnings: false,
					FastPush:               false,
					FMAActivationState:     "",
					GroupID:                "",
					IgnoreHTTPErrors:       false,
					PropertyName:           "test.edgesuite.net",
					PropertyID:             "prp_12345",
					PropertyVersion:        1,
					Network:                "STAGING",
					Status:                 "ACTIVE",
					NotifyEmails:           []string{"jsmith@akamai.com"},
				},
			},
		},
	}

	getActivationsResponseWithNote := papi.GetActivationsResponse{
		Response: papi.Response{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "grp_12345",
		},
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				{
					ActivationID:    "atv_5594260",
					ActivationType:  "ACTIVATE",
					PropertyName:    "test.edgesuite.net",
					PropertyID:      "prp_12345",
					PropertyVersion: 2,
					Network:         "STAGING",
					Status:          "ACTIVE",
					NotifyEmails:    []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
					Note:            "example note",
				},
			},
		},
	}

	getActivationsResponseWithEmptyEmails := papi.GetActivationsResponse{
		Response: papi.Response{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "grp_12345",
		},
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				{
					ActivationID:    "atv_5594260",
					ActivationType:  "ACTIVATE",
					PropertyName:    "test.edgesuite.net",
					PropertyID:      "prp_12345",
					PropertyVersion: 2,
					Network:         "STAGING",
					Status:          "ACTIVE",
					NotifyEmails:    []string{},
					Note:            "example note",
				},
			},
		},
	}

	tests := map[string]struct {
		init                func(*papi.Mock, *hapi.Mock, *mockProcessor, string)
		dir                 string
		snippetFilesToCheck []string
		jsonDir             string
		withError           error
		readVersion         string
		withIncludes        bool
	}{
		"basic property": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponse, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "LATEST",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir:     "basic",
			jsonDir: "basic/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"basic property with includes": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				// Includes
				c.On("ListReferencedIncludes", mock.Anything, papi.ListReferencedIncludesRequest{
					PropertyID:      "prp_12345",
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
					PropertyVersion: 5,
				}).Return(&getListReferencedIncludesResponse, nil).Once()
				expectGetIncludeVersion(c)
				var includeRuleResponse papi.GetIncludeRuleTreeResponse
				includeRules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_include_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(includeRules, &includeRuleResponse))
				c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(&includeRuleResponse, nil).Once()
				expectListIncludeActivations(c)

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponse, nil).Once()

				p.On("ProcessTemplates", TFData{
					Includes: []TFIncludeData{
						{
							ActivationNoteProduction:   "test production activation",
							ActivationNoteStaging:      "test staging activation",
							ContractID:                 "test_contract",
							ActivationEmailsProduction: []string{"test@example.com", "test1@example.com"},
							ActivationEmailsStaging:    []string{"test@example.com"},
							GroupID:                    "test_group",
							IncludeID:                  "inc_123456",
							IncludeName:                "test_include",
							IncludeType:                string(papi.IncludeTypeMicroServices),
							Networks:                   []string{"STAGING", "PRODUCTION"},
							RuleFormat:                 "v2020-11-02",
							VersionProduction:          "1",
							VersionStaging:             "1",
						},
					},
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "LATEST",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir:     "basic_property_with_includes",
			jsonDir: "basic_property_with_includes/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"test_include.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			withIncludes: true,
		},
		"basic property with cert provisioning type": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersion2HostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponse, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "DEFAULT",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "LATEST",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir:     "basic_with_cert_provisioning_type",
			jsonDir: "basic_with_cert_provisioning_type/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"import LATEST property version": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponse, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "LATEST",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir:     "basic",
			jsonDir: "basic/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"import not the latest property version": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 1, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 1,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersion1HostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivations1Response, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "1",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir:     "basic-v1",
			jsonDir: "basic-v1/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			readVersion: "1",
		},
		"property activation with note": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponseWithNote, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:         []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote: "example note",
						Version:        "LATEST",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir: "basic",
		},
		"property activation with empty emails": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponseWithEmptyEmails, nil).Once()

				p.On("ProcessTemplates", TFData{
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
						Version:              "LATEST",
						EdgeHostnames: map[string]EdgeHostname{
							"test-edgesuite-net": {
								EdgeHostname:             "test.edgesuite.net",
								EdgeHostnameID:           "ehn_2867480",
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:         []string{""},
						ActivationNote: "example note",
					},
					Section: "test_section",
				}).Return(nil).Once()
			},
			dir: "basic",
		},
		"error property not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrPropertyNotFound,
		},
		"error group not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&papi.GetRuleTreeResponse{}, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: ErrGroupNotFound,
		},
		"error property rules not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(nil, fmt.Errorf("oops")).Once()

			},
			withError: ErrPropertyRulesNotFound,
		},
		"error property version not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&papi.GetRuleTreeResponse{}, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(nil, fmt.Errorf("oops")).Once()

			},
			withError: ErrPropertyVersionNotFound,
		},
		"error product name not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&papi.GetRuleTreeResponse{}, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(nil, fmt.Errorf("oops")).Once()

			},
			withError: ErrProductNameNotFound,
		},
		"error hostnames not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&papi.GetRuleTreeResponse{}, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(nil, fmt.Errorf("oops")).Once()

			},
			withError: ErrHostnamesNotFound,
		},
		"error hostname details": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&papi.GetRuleTreeResponse{}, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(nil, fmt.Errorf("oops")).Once()

			},
			withError: ErrFetchingHostnameDetails,
		},
		"error saving files": {
			init: func(c *papi.Mock, h *hapi.Mock, p *mockProcessor, dir string) {
				c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
					Return(&searchPropertiesResponse, nil).Once()

				c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
					Return(&getPropertyResponse, nil).Once()

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: 5, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: false, RuleFormat: "latest"}).
					Return(&ruleResponse, nil).Once()

				c.On("GetGroups", mock.Anything).
					Return(&getGroupsResponse, nil).Once()

				c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getPropertyVersionsResponse, nil).Once()

				c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
					PropertyID:  "prp_12345",
					ActivatedOn: "",
					ContractID:  "test_contract",
					GroupID:     "grp_12345",
				}).Return(&getLatestVersionResponse, nil).Once()

				c.On("GetProducts", mock.Anything, papi.GetProductsRequest{
					ContractID: "test_contract",
				}).Return(&getProductsResponse, nil).Once()

				c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:      "prp_12345",
					PropertyVersion: 5,
					ContractID:      "test_contract",
					GroupID:         "grp_12345",
				}).Return(&getPropertyVersionHostnamesResponse, nil).Once()

				h.On("GetEdgeHostname", mock.Anything, 2867480).
					Return(&hapi.GetEdgeHostnameResponse{
						EdgeHostnameID:    2867480,
						RecordName:        "test",
						DNSZone:           "edgesuite.net",
						SecurityType:      "STANDARD-TLS",
						UseDefaultTTL:     false,
						UseDefaultMap:     false,
						IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
						ProductID:         "",
						TTL:               21600,
						Map:               "a;test.akamai.net",
						SerialNumber:      1461,
					}, nil).Once()

				c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&papi.GetEdgeHostnamesResponse{
					EdgeHostnames: papi.EdgeHostnameItems{
						Items: []papi.EdgeHostnameGetItem{
							{
								ID:                "ehn_2867480",
								Domain:            "test.edgesuite.net",
								ProductID:         "",
								DomainPrefix:      "test",
								DomainSuffix:      "edgesuite.net",
								Status:            "CREATED",
								Secure:            false,
								IPVersionBehavior: "IPV6_COMPLIANCE",
								UseCases:          []papi.UseCase(nil),
							},
						},
					},
				}, nil).Once()

				c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
					PropertyID: "prp_12345",
					ContractID: "test_contract",
					GroupID:    "grp_12345",
				}).Return(&getActivationsResponse, nil).Once()

				p.On("ProcessTemplates", TFData{
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
								ProductName:              "HTTP_Content_Del",
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
								Hostname:                 "test.edgesuite.net",
								EdgeHostnameResourceName: "test-edgesuite-net",
								CertProvisioningType:     "CPS_MANAGED",
							},
						},
						Emails:  []string{"jsmith@akamai.com"},
						Version: "LATEST",
					},
					Section: "test_section",
				}).Return(fmt.Errorf("oops")).Once()
			},
			dir:       "basic",
			withError: ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mc := new(papi.Mock)
			mh := new(hapi.Mock)
			mp := new(mockProcessor)
			test.init(mc, mh, mp, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createProperty(ctx, "test.edgesuite.net", test.readVersion, section, fmt.Sprintf("./testdata/res/%s", test.jsonDir), "./", test.withIncludes, mc, mh, mp)
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
		})
	}
}

func TestProcessPolicyTemplates(t *testing.T) {

	useCases := []papi.UseCase{
		{
			Option:  "BACKGROUND",
			Type:    "GLOBAL",
			UseCase: "Download_Mode",
		},
	}
	useCasesJSON, err := json.MarshalIndent(useCases, "", "  ")
	assert.NoError(t, err)

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
		withIncludes bool
	}{
		"property": {
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
					RuleFormat:           "latest",
					IsSecure:             "false",
					Version:              "LATEST",
					EdgeHostnames: map[string]EdgeHostname{
						"test-edgesuite-net": {
							EdgeHostname:             "test.edgesuite.net",
							EdgeHostnameID:           "ehn_2867480",
							ProductName:              "HTTP_Content_Del",
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
							Hostname:                 "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
						},
					},
					Emails: []string{"jsmith@akamai.com"},
				},
				Section: "test_section",
			},
			dir:          "basic",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with include": {
			givenData: TFData{
				Includes: []TFIncludeData{
					{
						ActivationNoteProduction:   "test production activation",
						ActivationNoteStaging:      "test staging activation",
						ContractID:                 "test_contract",
						ActivationEmailsProduction: []string{"test@example.com", "test1@example.com"},
						ActivationEmailsStaging:    []string{"test@example.com"},
						GroupID:                    "test_group",
						IncludeID:                  "inc_123456",
						IncludeName:                "test_include",
						IncludeType:                string(papi.IncludeTypeMicroServices),
						Networks:                   []string{"STAGING", "PRODUCTION"},
						RuleFormat:                 "v2020-11-02",
						VersionProduction:          "1",
						VersionStaging:             "1",
					},
				},
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
					Version:              "LATEST",
					EdgeHostnames: map[string]EdgeHostname{
						"test-edgesuite-net": {
							EdgeHostname:             "test.edgesuite.net",
							EdgeHostnameID:           "ehn_2867480",
							ProductName:              "HTTP_Content_Del",
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
							Hostname:                 "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
						},
					},
					Emails: []string{"jsmith@akamai.com"},
				},
				Section: "test_section",
			},
			dir:          "basic_property_with_includes",
			filesToCheck: []string{"property.tf", "includes.tf", "variables.tf", "import.sh"},
			withIncludes: true,
		},
		"property with use cases": {
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
					RuleFormat:           "latest",
					IsSecure:             "false",
					Version:              "3",
					EdgeHostnames: map[string]EdgeHostname{
						"test-edgesuite-net": {
							EdgeHostname:             "test.edgesuite.net",
							EdgeHostnameID:           "ehn_2867480",
							ProductName:              "HTTP_Content_Del",
							ContractID:               "test_contract",
							GroupID:                  "grp_12345",
							ID:                       "",
							IPv6:                     "IPV6_COMPLIANCE",
							SecurityType:             "STANDARD-TLS",
							EdgeHostnameResourceName: "test-edgesuite-net",
							UseCases:                 string(useCasesJSON),
						},
					},
					Hostnames: map[string]Hostname{
						"test.edgesuite.net": {
							Hostname:                 "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
						},
					},
					Emails: []string{"jsmith@akamai.com"},
				},
				Section: "test_section",
			},
			dir:          "basic_with_use_cases",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with activation note": {
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
					RuleFormat:           "latest",
					IsSecure:             "false",
					Version:              "LATEST",
					EdgeHostnames: map[string]EdgeHostname{
						"test-edgesuite-net": {
							EdgeHostname:             "test.edgesuite.net",
							EdgeHostnameID:           "ehn_2867480",
							ProductName:              "HTTP_Content_Del",
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
							Hostname:                 "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
						},
					},
					Emails:         []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
					ActivationNote: "example note",
				},
				Section: "test_section",
			},
			dir:          "basic_with_activation_note",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			templateToFile := map[string]string{
				"property.tmpl":  fmt.Sprintf("./testdata/res/%s/property.tf", test.dir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
			}

			if test.withIncludes {
				templateToFile["includes.tmpl"] = fmt.Sprintf("./testdata/res/%s/includes.tf", test.dir)
			}
			processor := templates.FSTemplateProcessor{
				TemplatesFS:     templateFiles,
				TemplateTargets: templateToFile,
				AdditionalFuncs: template.FuncMap{
					"ToLower": strings.ToLower,
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

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

func TestNormalizeRuleName(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"one word": {
			given:    "testName1",
			expected: "testName1",
		},
		"with spaces": {
			given:    "this is test name",
			expected: "this_is_test_name",
		},
		"with slashes": {
			given:    "this\\is test/name",
			expected: "this_is_test_name",
		},
		"with dots and dashes": {
			given:    "this-is.test-name.1",
			expected: "this-is.test-name.1",
		},
		"with other symbols": {
			given:    "this#is!te$tN@me",
			expected: "this_is_te_tN_me",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := normalizeRuleName(test.given)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestRuleNameNormalizer(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
		preTest  []string
	}{
		"no duplicates": {
			given:    "this is test name",
			expected: "this_is_test_name",
			preTest:  []string{"testName1", "test name"},
		},
		"one duplicate": {
			given:    "test@name",
			expected: "test_name1",
			preTest:  []string{"testName", "test name"},
		},
		"two duplicates": {
			given:    "test@name",
			expected: "test_name2",
			preTest:  []string{"test%name", "testName", "test name"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			normalizer := ruleNameNormalizer()
			for _, n := range test.preTest {
				normalizer(n)
			}
			assert.Equal(t, test.expected, normalizer(test.given))
		})
	}
}
