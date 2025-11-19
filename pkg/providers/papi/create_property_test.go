package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
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

func TestCreateProperty(t *testing.T) {
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
	getPropertyResponse := func() *papi.GetPropertyResponse {
		return &papi.GetPropertyResponse{
			Properties: papi.PropertiesItems{
				Items: []*papi.Property{
					{
						AccountID:         "test_account",
						AssetID:           "aid_10541511",
						ContractID:        "test_contract",
						GroupID:           "grp_12345",
						LatestVersion:     5,
						Note:              "",
						ProductionVersion: nil,
						PropertyID:        "prp_12345",
						PropertyName:      "test.edgesuite.net",
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
				ProductionVersion: nil,
				PropertyID:        "prp_12345",
				PropertyName:      "test.edgesuite.net",
				StagingVersion:    nil,
			},
		}
	}

	getPropertyEnhancementTLSResponse := papi.GetPropertyResponse{
		Properties: papi.PropertiesItems{
			Items: []*papi.Property{
				{
					AccountID:         "test_account",
					AssetID:           "aid_10541511",
					ContractID:        "test_contract",
					GroupID:           "grp_12345",
					LatestVersion:     5,
					Note:              "",
					ProductionVersion: nil,
					PropertyID:        "prp_12345",
					PropertyName:      "test.edgekey.net",
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
			ProductionVersion: nil,
			PropertyID:        "prp_12345",
			PropertyName:      "test.edgekey.net",
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

	getPropertyVersionWithDigitsAndSpacesInHostnamesResponse := papi.GetPropertyVersionHostnamesResponse{
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
					CnameTo:              "1 test.edgesuite.net",
					CertProvisioningType: "CPS_MANAGED",
				},
			},
		},
	}

	getPropertyVersionWithCCMHostnamesResponse := papi.GetPropertyVersionHostnamesResponse{
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
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_34343434",
					CnameFrom:            "foo.com",
					CnameTo:              "foo.com.edgekey.net",
					CertProvisioningType: "CCM",
					CCMCertificates: &papi.CCMCertificates{
						RSACertID:     "2226",
						RSACertLink:   "/ccm/v1/certificates/2226",
						ECDSACertID:   "7890",
						ECDSACertLink: "/ccm/v1/certificates/7890",
					},
				},
			},
		},
	}

	getPropertyVersionHostnamesEnhancementTLSResponse := papi.GetPropertyVersionHostnamesResponse{
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
					CnameFrom:            "test.edgekey.net",
					CnameTo:              "test.edgekey.net",
					CertProvisioningType: "CPS_MANAGED",
				},
			},
		},
	}

	getPropertyVersionEmptyHostnameIDResponse := papi.GetPropertyVersionHostnamesResponse{
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
					EdgeHostnameID:       "",
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
					PropertyVersion:        5,
					Network:                "STAGING",
					Status:                 "ACTIVE",
					NotifyEmails:           []string{"jsmith@akamai.com"},
				},
			},
		},
	}

	getProductionActivationsResponse := papi.GetActivationsResponse{
		Response: papi.Response{
			AccountID:  "test_account",
			ContractID: "test_contract",
			GroupID:    "grp_12345",
		},
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				{
					ActivationID:           "atv_5594260",
					ActivationType:         "ACTIVATE",
					PropertyName:           "test.edgesuite.net",
					PropertyID:             "prp_12345",
					PropertyVersion:        5,
					Network:                "PRODUCTION",
					AcknowledgeAllWarnings: false,
					Status:                 "ACTIVE",
					NotifyEmails:           []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
					Note:                   "example production note",
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
					ActivationID:           "atv_5594260",
					ActivationType:         "ACTIVATE",
					PropertyName:           "test.edgesuite.net",
					PropertyID:             "prp_12345",
					PropertyVersion:        5,
					AcknowledgeAllWarnings: false,
					Network:                "STAGING",
					Status:                 "ACTIVE",
					NotifyEmails:           []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
					Note:                   "example staging note",
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
					ActivationID:           "atv_5594260",
					ActivationType:         "ACTIVATE",
					PropertyName:           "test.edgesuite.net",
					PropertyID:             "prp_12345",
					PropertyVersion:        5,
					AcknowledgeAllWarnings: false,
					Network:                "STAGING",
					Status:                 "ACTIVE",
					NotifyEmails:           []string{},
					Note:                   "example note",
				},
			},
		},
	}

	hapiGetEdgeHostnameResponse := hapi.GetEdgeHostnameResponse{
		EdgeHostnameID:    2867480,
		RecordName:        "test",
		DNSZone:           "edgesuite.net",
		SecurityType:      "STANDARD-TLS",
		UseDefaultTTL:     true,
		UseDefaultMap:     false,
		IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
		TTL:               21600,
		Map:               "a;test.akamai.net",
		SerialNumber:      1461,
	}

	hapiGetEdgeHostnameResponseNonDefaultTTL := hapi.GetEdgeHostnameResponse{
		EdgeHostnameID:    2867480,
		RecordName:        "test",
		DNSZone:           "edgesuite.net",
		SecurityType:      "STANDARD-TLS",
		UseDefaultTTL:     false,
		UseDefaultMap:     false,
		IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
		TTL:               600,
		Map:               "a;test.akamai.net",
		SerialNumber:      1461,
	}

	edgeHostnameWithTTL := map[string]EdgeHostname{
		"test-edgesuite-net": {
			EdgeHostname:             "test.edgesuite.net",
			EdgeHostnameID:           "ehn_2867480",
			ContractID:               "test_contract",
			GroupID:                  "grp_12345",
			ID:                       "",
			IPv6:                     "IPV6_COMPLIANCE",
			SecurityType:             "STANDARD-TLS",
			EdgeHostnameResourceName: "test-edgesuite-net",
			TTL:                      600,
		},
	}

	hapiGetEdgeHostnameResponseEnhancementTLS := hapi.GetEdgeHostnameResponse{
		EdgeHostnameID:    2867480,
		RecordName:        "test",
		DNSZone:           "edgekey.net",
		SecurityType:      "ENHANCED-TLS",
		UseDefaultTTL:     true,
		UseDefaultMap:     false,
		IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
		TTL:               21600,
		Map:               "a;test.akamai.net",
		SerialNumber:      1461,
	}

	var noFilters []func([]string) ([]string, error)
	otherRuleFormatFilter := []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")}

	tests := map[string]struct {
		init                func(*papi.Mock, *hapi.Mock, *templates.MockProcessor, *templates.MockMultiTargetProcessor, string)
		dir                 string
		snippetFilesToCheck []string
		jsonDir             string
		withError           error
		readVersion         string
		rulesAsHCL          bool
		withBootstrap       bool
		splitDepth          *int
		edgercPath          string
		edgercSection       string
	}{
		"basic property (with hostname's cnameTo starting with a digit)": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionWithDigitsAndSpacesInHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaultsHavingDigitsAndSpacesInHostnameDetails().build(), noFilters, nil)
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
		"basic property (with ccm hostname)": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionWithCCMHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				hostnames := map[string]Hostname{
					"test.edgesuite.net": {
						CnameFrom:                "test.edgesuite.net",
						CnameTo:                  "test.edgesuite.net",
						EdgeHostnameResourceName: "test-edgesuite-net",
						CertProvisioningType:     "CPS_MANAGED",
						IsActive:                 true,
					},
					"foo.com": {
						CnameFrom:                "foo.com",
						CnameTo:                  "foo.com.edgekey.net",
						EdgeHostnameResourceName: "",
						CertProvisioningType:     "CCM",
						IsActive:                 true,
						CCMCertificates: &CCMCertificates{
							RSACertID:   "2226",
							ECDSACertID: "7890",
						},
					},
				}
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().
					withHostnames(hostnames).build(), noFilters, nil)
			},
			dir:     "basic-ccm-hostnames",
			jsonDir: "basic/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},

		"basic property with edgehostname with non default ttl": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponseNonDefaultTTL, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withEdgeHostname(edgeHostnameWithTTL).build(), noFilters, nil)
			},
			dir:     "basic-non-default-ttl",
			jsonDir: "basic-non-default-ttl/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"property with enhancement tls": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyEnhancementTLSResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				getPropertyVersionsResponse := getPropertyVersionsResponse
				getPropertyVersionsResponse.PropertyName = "test.edgekey.net"
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				getLatestVersionResponse := getLatestVersionResponse
				getLatestVersionResponse.PropertyName = "test.edgekey.net"
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesEnhancementTLSResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponseEnhancementTLS, nil)
				mockGetEdgeHostnames(c)
				cert := hapi.GetCertificateResponse{

					AvailableDomains: []string{"*.dev-exp-terraform-automation-test.com"},
					CertificateID:    "123456",
					CertificateType:  "THIRD_PARTY",
					CommonName:       "*.dev-exp-terraform-automation-test.com",
					ExpirationDate:   *newTimeFromString(t, "2025-05-21T12:24:21.000+00:00"),
					SerialNumber:     "fa:ke:76:82:a9:3f:14:ba:6b:93:01:57:43:10:0c:34",
					SlotNumber:       30278,
					Status:           "DEPLOYED",
					ValidationType:   "THIRD_PARTY",
				}
				mockGetCertificate(h, "edgekey.net", "test", &cert, nil)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				data := TFData{
					Property: TFPropertyData{
						GroupName:            "test_group",
						GroupID:              "grp_12345",
						ContractID:           "test_contract",
						PropertyResourceName: "test-edgekey-net",
						PropertyName:         "test.edgekey.net",
						PropertyID:           "prp_12345",
						ProductID:            "prd_HTTP_Content_Del",
						ProductName:          "HTTP_Content_Del",
						RuleFormat:           "latest",
						IsSecure:             "false",
						EdgeHostnames: map[string]EdgeHostname{
							"test-edgekey-net": {
								EdgeHostname:             "test.edgekey.net",
								EdgeHostnameID:           "ehn_2867480",
								ContractID:               "test_contract",
								GroupID:                  "grp_12345",
								ID:                       "",
								IPv6:                     "IPV6_COMPLIANCE",
								SecurityType:             "ENHANCED-TLS",
								EdgeHostnameResourceName: "test-edgekey-net",
								CertificateID:            123456,
							},
						},
						Hostnames: map[string]Hostname{
							"test.edgekey.net": {
								CnameFrom:                "test.edgekey.net",
								CnameTo:                  "test.edgekey.net",
								EdgeHostnameResourceName: "test-edgekey-net",
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
					EdgercPath: defaultEdgercPath,
					Section:    defaultSection,
				}

				mockProcessTemplates(p, (&tfDataBuilder{}).withData(data).build(), noFilters, nil)
			},
			dir:     "enhancement-tls",
			jsonDir: "enhancement-tls/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"property with enhancement tls but without certificate": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyEnhancementTLSResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				getPropertyVersionsResponse := getPropertyVersionsResponse
				getPropertyVersionsResponse.PropertyName = "test.edgekey.net"
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				getLatestVersionResponse := getLatestVersionResponse
				getLatestVersionResponse.PropertyName = "test.edgekey.net"
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesEnhancementTLSResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponseEnhancementTLS, nil)
				mockGetEdgeHostnames(c)
				resp := hapi.Error{
					Type:            "CERTIFICATE_NOT_FOUND",
					Title:           "Certificate Not Found",
					Status:          404,
					Detail:          "Details are not available for this certificate; the certificate is missing or access is denied",
					Instance:        "/hapi/error-instances/a30f67cc-df20-4e02-bbc3-cf7c204a4aab",
					RequestInstance: "http://origin.pulsar.akamai.com/hapi/open/v1/dns-zones/edgekey.net/edge-hostnames/example.com/certificate?depth=ALL&accountSwitchKey=F-AC-1937217#d7aa7348",
					Method:          "GET",
					RequestTime:     "2022-11-30T18:51:43.482982Z",
				}
				err := fmt.Errorf("%s: %s: %w", hapi.ErrGetCertificate, hapi.ErrNotFound, &resp)
				mockGetCertificate(h, "edgekey.net", "test", nil, err)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				data := TFData{
					Property: TFPropertyData{
						GroupName:            "test_group",
						GroupID:              "grp_12345",
						ContractID:           "test_contract",
						PropertyResourceName: "test-edgekey-net",
						PropertyName:         "test.edgekey.net",
						PropertyID:           "prp_12345",
						ProductID:            "prd_HTTP_Content_Del",
						ProductName:          "HTTP_Content_Del",
						RuleFormat:           "latest",
						IsSecure:             "false",
						EdgeHostnames: map[string]EdgeHostname{
							"test-edgekey-net": {
								EdgeHostname:             "test.edgekey.net",
								EdgeHostnameID:           "ehn_2867480",
								ContractID:               "test_contract",
								GroupID:                  "grp_12345",
								ID:                       "",
								IPv6:                     "IPV6_COMPLIANCE",
								SecurityType:             "ENHANCED-TLS",
								EdgeHostnameResourceName: "test-edgekey-net",
								CertificateID:            0,
							},
						},
						Hostnames: map[string]Hostname{
							"test.edgekey.net": {
								CnameFrom:                "test.edgekey.net",
								CnameTo:                  "test.edgekey.net",
								EdgeHostnameResourceName: "test-edgekey-net",
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
					EdgercPath: defaultEdgercPath,
					Section:    defaultSection,
				}

				mockProcessTemplates(p, (&tfDataBuilder{}).withData(data).build(), noFilters, nil)
			},
			dir:     "enhancement-tls",
			jsonDir: "enhancement-tls/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"basic property not active the latest": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivations1Response, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withStagingVersion(1, false).build(), noFilters, nil)
			},
			dir:     "basic_not_latest",
			jsonDir: "basic_not_latest/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"basic property with empty hostname id": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				var ruleResponse papi.GetRuleTreeResponse
				rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
				assert.NoError(t, err)
				err = json.Unmarshal(rules, &ruleResponse)
				assert.NoError(t, err)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionEmptyHostnameIDResponse, nil)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withEdgeHostname(map[string]EdgeHostname{}).
					withIsActive(false).build(), noFilters, nil)
			},
			dir:     "basic_property_with_empty_hostname_id",
			jsonDir: "basic_property_with_empty_hostname_id/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"basic property with hostname bucket and no hostname activations": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockEmptyListPropertyHostnameActivations(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(nil, nil, "", "", nil, nil).
					build(), noFilters, nil)
			},
			dir:     "hostname-bucket/no-hostnames",
			jsonDir: "hostname-bucket/no-hostnames/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
		},
		"basic property with hostname bucket and only STAGING hostnames": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, true, false)
				mockListActivePropertyHostnames(c, 10, "STAGING", "CPS_MANAGED")
				mockAddTemplateTargetHostnameBucket(p)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(generateHostnames(10, "CPS_MANAGED", "STAGING", "ehn_12345"),
						nil, "staging note", "", []string{"test@mail.com"}, nil).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/staging",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/staging/property-snippets",
		},
		"basic property with hostname bucket and only STAGING hostnames with paging": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, true, false)
				mockListActivePropertyHostnames(c, 2500, "STAGING", "CPS_MANAGED")
				mockAddTemplateTargetHostnameBucket(p)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(generateHostnames(2500, "CPS_MANAGED", "STAGING", "ehn_12345"),
						nil, "staging note", "", []string{"test@mail.com"}, nil).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/staging",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/staging/property-snippets",
		},
		"basic property with hostname bucket and only PRODUCTION hostnames": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, false, true)
				mockListActivePropertyHostnames(c, 10, "PRODUCTION", "DEFAULT")
				mockAddTemplateTargetHostnameBucket(p)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(nil, generateHostnames(10, "DEFAULT", "PRODUCTION", "ehn_12345"),
						"", "production note", nil, []string{"test@mail.com"}).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/production",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/production/property-snippets",
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, true, true)
				mockListActivePropertyHostnames(c, 10, "STAGING", "DEFAULT")
				mockListActivePropertyHostnames(c, 10, "PRODUCTION", "DEFAULT")
				mockAddTemplateTargetHostnameBucket(p)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(generateHostnames(10, "DEFAULT", "STAGING", "ehn_12345"),
						generateHostnames(10, "DEFAULT", "PRODUCTION", "ehn_12345"), "staging note", "production note", []string{"test@mail.com"}, []string{"test@mail.com"}).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/staging-production",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/staging-production/property-snippets",
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames with differences": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, true, true)
				mockListActivePropertyHostnames(c, 5, "STAGING", "DEFAULT")
				mockListActivePropertyHostnames(c, 10, "PRODUCTION", "DEFAULT")
				mockAddTemplateTargetHostnameBucket(p)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(generateHostnames(5, "DEFAULT", "STAGING", "ehn_12345"),
						generateHostnames(10, "DEFAULT", "PRODUCTION", "ehn_12345"), "staging note", "production note", []string{"test@mail.com"}, []string{"test@mail.com"}).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/staging-production-diff",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/staging-production-diff/property-snippets",
		},
		"basic property with hostname bucket and only STAGING activation but no hostnames": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockListPropertyHostnameActivations(c, true, false)
				mockListActivePropertyHostnames(c, 0, "STAGING", "CPS_MANAGED")
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(generateHostnames(0, "CPS_MANAGED", "STAGING", "ehn_12345"),
						nil, "staging note", "", []string{"test@mail.com"}, nil).
					build(), noFilters, nil)
			},
			dir: "hostname-bucket/staging",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			jsonDir: "hostname-bucket/staging/property-snippets",
		},
		"basic property with rules as datasource": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockAddTemplateTargetRules(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withRuleFormat("v2023-01-05").
					withRules(flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))).build(), otherRuleFormatFilter, nil)
			},
			dir:        "ruleformats/basic-rules-datasource",
			rulesAsHCL: true,
		},
		"basic property with rules as datasource with unsupported rule format": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockAddTemplateTargetRules(p)
				mockTemplateExist(p, "rules_latest.tmpl", false)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withIsActive(false).
					withRules(flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))).build(), noFilters, nil)
			},
			withError:  ErrUnsupportedRuleFormat,
			dir:        "basic",
			rulesAsHCL: true,
		},
		"basic property with rules as datasource with unsupported behaviors and criteria": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockAddTemplateTargetRules(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withRuleFormat("v2023-01-05").
					withRules(flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))).build(), otherRuleFormatFilter, ErrSavingFiles)
			},
			withError:  ErrSavingFiles,
			dir:        "ruleformats/basic-rules-datasource-unknown",
			rulesAsHCL: true,
		},
		"basic property with children with split-depth=0": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, mm *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)

				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)

				tfdata := (&tfDataBuilder{}).
					withData(getTestData("basic property with multiple children as hcl")).
					asHCL(true).
					withSplitDepth(true, "test-edgesuite-net_rule_default").
					withRuleFormat("v2023-01-05").
					build()
				mockProcessTemplates(p, tfdata, otherRuleFormatFilter, nil)

				mockModuleConfig(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)

				multiTargetData := templates.MultiTargetData{
					"split-depth-rules.tmpl": templates.DataForTarget{
						"rules/test-edgesuite-net_default.tf": (&tfDataBuilder{}).
							withRules(flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))).
							withSplitDepth(true, "").build(),
					},
				}
				mm.On("ProcessTemplates", multiTargetData, mock.AnythingOfType("func([]string) ([]string, error)")).Return(nil).Once()
			},
			dir:        "multitarget-property-with-children",
			rulesAsHCL: true,
			splitDepth: ptr.To(0),
		},
		"basic property with children with split-depth=2": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, mm *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)

				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)

				tfdata := (&tfDataBuilder{}).
					withData(getTestData("basic property with multiple children as hcl")).
					asHCL(true).
					withSplitDepth(true, "test-edgesuite-net_rule_default").
					withRuleFormat("v2023-01-05").
					build()
				mockProcessTemplates(p, tfdata, otherRuleFormatFilter, nil)

				mockModuleConfig(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)

				multiTargetData := templates.MultiTargetData{
					"split-depth-rules.tmpl": templates.DataForTarget{
						"rules/test-edgesuite-net_default.tf": (&tfDataBuilder{}).
							withRules([]*WrappedRules{flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))[0]}).
							withSplitDepth(true, "").build(),
						"rules/test-edgesuite-net_default_new_rule.tf": (&tfDataBuilder{}).
							withRules([]*WrappedRules{flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))[1]}).
							withSplitDepth(true, "").build(),
						"rules/test-edgesuite-net_default_new_rule_new_rule_1.tf": (&tfDataBuilder{}).
							withRules([]*WrappedRules{flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))[2]}).
							withSplitDepth(true, "").build(),
					},
				}
				mm.On("ProcessTemplates", multiTargetData, mock.AnythingOfType("func([]string) ([]string, error)")).Return(nil).Once()
			},
			dir:        "multitarget-property-with-flatten-children",
			rulesAsHCL: true,
			splitDepth: ptr.To(2),
		},
		"basic property with cert provisioning type": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersion2HostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withCertProvisioningType("DEFAULT").build(), noFilters, nil)
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
		"basic property with bootstrap": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withBootstrap(true).build(), noFilters, nil)
			},
			dir:     "basic-bootstrap",
			jsonDir: "basic-bootstrap/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			withBootstrap: true,
		},
		"basic property with hostname bucket and bootstrap and no hostname activations": {
			init: func(c *papi.Mock, _ *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				propertyResponse := getPropertyResponse()
				propertyResponse.Property.PropertyType = ptr.To("HOSTNAME_BUCKET")
				propertyResponse.Properties.Items[0].PropertyType = ptr.To("HOSTNAME_BUCKET")
				mockGetProperty(c, propertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockEmptyListPropertyHostnameActivations(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).
					withDefaults().
					withEdgeHostname(map[string]EdgeHostname{}).
					withHostnames(nil).
					withHostnameBucket(nil, nil, "", "", nil, nil).
					withBootstrap(true).
					build(), noFilters, nil)
			},
			dir:     "hostname-bucket/bootstrap/no-hostnames",
			jsonDir: "hostname-bucket/bootstrap/no-hostnames/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			withBootstrap: true,
		},
		"import LATEST property version": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().build(), noFilters, nil)
			},
			dir:     "basic",
			jsonDir: "basic/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
				"IPCUID_Invalidation.json",
				"ipcuid_invalidation1.json",
			},
		},
		"import not the latest property version": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 1, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 1, &getPropertyVersion1HostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivations1Response, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withVersion("1").withStagingVersion(1, false).build(), noFilters, nil)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponseWithNote, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withActivationNote("example staging note").
					withEmails([]string{"jsmith@akamai.com", "rjohnson@akamai.com"}).build(), noFilters, nil)
			},
			dir: "basic",
		},
		"property with production activation": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockGetActivations(c, &getProductionActivationsResponse, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withOnlyProductionActivation([]string{"jsmith@akamai.com", "rjohnson@akamai.com"}, "example production note", 5).build(), noFilters, nil)
			},
			dir: "basic",
		},
		"property with both activations": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponseWithNote, nil)
				mockGetActivations(c, &getProductionActivationsResponse, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withOnlyProductionActivation([]string{"jsmith@akamai.com", "rjohnson@akamai.com"}, "example production note", 5).withActivationNote("example staging note").
					withEmails([]string{"jsmith@akamai.com", "rjohnson@akamai.com"}).withStagingVersion(5, true).build(), noFilters, nil)
			},
			dir: "basic",
		},
		"property activation with empty emails": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponseWithEmptyEmails, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withActivationNote("example note").
					withEmails([]string{""}).build(), noFilters, nil)
			},
			dir: "basic",
		},
		"error property not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, nil, fmt.Errorf("oops"))
			},
			withError: ErrPropertyNotFound,
		},
		"error group not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, nil, fmt.Errorf("oops"))
			},
			withError: ErrGroupNotFound,
		},
		"error property rules not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)

				mockGetRuleTree(c, 5, nil, fmt.Errorf("oops"))

			},
			withError: ErrPropertyRulesNotFound,
		},
		"error property version not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, nil, fmt.Errorf("oops"))

			},
			withError: ErrPropertyVersionNotFound,
		},
		"error fetching property activation": {
			init: func(c *papi.Mock, h *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, nil, fmt.Errorf("oops"))
			},
			withError: ErrFetchingActivationDetails,
		},
		"error product name not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, nil, fmt.Errorf("oops"))

			},
			withError: ErrProductNameNotFound,
		},
		"error hostnames not found": {
			init: func(c *papi.Mock, _ *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, nil, fmt.Errorf("oops"))
			},
			withError: ErrHostnamesNotFound,
		},
		"error hostname details": {
			init: func(c *papi.Mock, h *hapi.Mock, _ *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, _ string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, nil, fmt.Errorf("oops"))
			},
			withError: ErrFetchingHostnameDetails,
		},
		"error saving files": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)

				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockAddTemplateTargetRules(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().build(), noFilters, fmt.Errorf("oops"))
			},
			dir:       "basic",
			withError: ErrSavingFiles,
		},
		"non default edgerc path and section": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, _ *templates.MockMultiTargetProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, getPropertyResponse())

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponseNonDefaultTTL, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse, nil)
				mockGetActivations(c, &papi.GetActivationsResponse{}, nil)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withEdgeHostname(edgeHostnameWithTTL).withEdgercPathAndSection("/non/default/path/to/edgerc", "non_default_section").build(), noFilters, nil)
			},
			dir:           "basic-non-default-ttl",
			jsonDir:       "basic-non-default-ttl/property-snippets",
			edgercPath:    "/non/default/path/to/edgerc",
			edgercSection: "non_default_section",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.edgercPath == "" {
				test.edgercPath = defaultEdgercPath
			}
			if test.edgercSection == "" {
				test.edgercSection = defaultSection
			}
			mc := new(papi.Mock)
			mh := new(hapi.Mock)
			mp := new(templates.MockProcessor)
			mm := new(templates.MockMultiTargetProcessor)
			test.init(mc, mh, mp, mm, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			options := propertyOptions{
				propertyName:  "test.edgesuite.net",
				edgercPath:    test.edgercPath,
				section:       test.edgercSection,
				tfWorkPath:    "./",
				version:       test.readVersion,
				rulesAsHCL:    test.rulesAsHCL,
				withBootstrap: test.withBootstrap,
				splitDepth:    test.splitDepth,
			}
			err := createProperty(ctx, options, fmt.Sprintf("./testdata/res/%s", test.jsonDir), mc, mh, mp, mm)
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

func mockSearchProperties(p *papi.Mock, searchPropertiesResponse *papi.SearchResponse, err error) {
	p.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
		Return(searchPropertiesResponse, err).Once()
}

func mockGetProperty(p *papi.Mock, getPropertyResponse *papi.GetPropertyResponse) {
	p.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
		Return(getPropertyResponse, nil).Once()
}

func mockGetRuleTree(p *papi.Mock, propertyVersion int, ruleResponse *papi.GetRuleTreeResponse, err error) {
	p.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: propertyVersion, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: true, RuleFormat: "latest"}).
		Return(ruleResponse, err).Once()
}

func mockGetGroups(p *papi.Mock, getGroupsResponse *papi.GetGroupsResponse, err error) {
	p.On("GetGroups", mock.Anything).
		Return(getGroupsResponse, err).Once()
}

func mockGetPropertyVersions(p *papi.Mock, getPropertyVersionsResponse *papi.GetPropertyVersionsResponse, err error) {
	p.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
		PropertyID: "prp_12345",
		ContractID: "test_contract",
		GroupID:    "grp_12345",
	}).Return(getPropertyVersionsResponse, err).Once()
}

func mockGetLatestVersion(p *papi.Mock, getLatestVersionResponse *papi.GetPropertyVersionsResponse) {
	p.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
		PropertyID:  "prp_12345",
		ActivatedOn: "",
		ContractID:  "test_contract",
		GroupID:     "grp_12345",
	}).Return(getLatestVersionResponse, nil).Once()
}

func mockGetProducts(p *papi.Mock, getProductsResponse *papi.GetProductsResponse, err error) {
	p.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "test_contract"}).
		Return(getProductsResponse, err).Once()
}

func mockEmptyListPropertyHostnameActivations(p *papi.Mock) {
	p.On("ListPropertyHostnameActivations", mock.Anything, papi.ListPropertyHostnameActivationsRequest{
		PropertyID: "prp_12345",
		Limit:      999,
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
	}).Return(&papi.ListPropertyHostnameActivationsResponse{
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
		HostnameActivations: papi.HostnameActivationsList{
			Items: []papi.HostnameActivationListItem{},
		},
	}, nil).Once()
}

func mockListPropertyHostnameActivations(p *papi.Mock, staging, production bool) {
	resp := papi.ListPropertyHostnameActivationsResponse{
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
	}
	if staging {
		resp.HostnameActivations.Items = append(resp.HostnameActivations.Items, papi.HostnameActivationListItem{
			ActivationType:       "ACTIVATE",
			HostnameActivationID: "atv_0",
			PropertyID:           "prp_12345",
			Network:              "STAGING",
			Status:               "ACTIVE",
			Note:                 "staging note",
			NotifyEmails:         []string{"test@mail.com"},
		})
	}
	if production {
		resp.HostnameActivations.Items = append(resp.HostnameActivations.Items, papi.HostnameActivationListItem{
			ActivationType:       "ACTIVATE",
			HostnameActivationID: "atv_1",
			PropertyID:           "prp_12345",
			Network:              "PRODUCTION",
			Status:               "ACTIVE",
			Note:                 "production note",
			NotifyEmails:         []string{"test@mail.com"},
		})
	}

	p.On("ListPropertyHostnameActivations", mock.Anything, papi.ListPropertyHostnameActivationsRequest{
		PropertyID: "prp_12345",
		Limit:      999,
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
	}).Return(&resp, nil).Once()
}

func mockListActivePropertyHostnames(p *papi.Mock, n int, network, certType string) {
	offset := 0
	hostnameItems := generateHostnames(n, certType, network, "ehn_12345")

	for len(hostnameItems) > 999 {
		req := papi.ListActivePropertyHostnamesRequest{
			PropertyID: "prp_12345",
			Limit:      999,
			Offset:     offset,
			Network:    papi.ActivationNetwork(network),
			ContractID: "ctr_test_contract",
			GroupID:    "grp_12345",
		}
		offset += 999
		resp := papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test_contract",
			GroupID:    "grp_12345",
			PropertyID: "prp_12345",
			Hostnames: papi.HostnamesResponseItems{
				Items:            hostnameItems[:999],
				CurrentItemCount: 999,
				TotalItems:       n,
			},
		}
		hostnameItems = hostnameItems[999:]
		p.On("ListActivePropertyHostnames", mock.Anything, req).Return(&resp, nil).Once()
	}

	req := papi.ListActivePropertyHostnamesRequest{
		PropertyID: "prp_12345",
		Limit:      999,
		Offset:     offset,
		Network:    papi.ActivationNetwork(network),
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
	}
	resp := papi.ListActivePropertyHostnamesResponse{
		ContractID: "ctr_test_contract",
		GroupID:    "grp_12345",
		PropertyID: "prp_12345",
		Hostnames: papi.HostnamesResponseItems{
			Items:            hostnameItems,
			CurrentItemCount: len(hostnameItems),
			TotalItems:       n,
		},
	}
	p.On("ListActivePropertyHostnames", mock.Anything, req).Return(&resp, nil).Once()
}

func mockAddTemplateTargetHostnameBucket(t *templates.MockProcessor) {
	t.On("AddTemplateTarget", "hostname_bucket.tmpl", "hostname_bucket.tf").Once()
}

func mockGetPropertyVersionHostnames(p *papi.Mock, propertyVersion int, getPropertyVersionHostnamesResponse *papi.GetPropertyVersionHostnamesResponse, err error) {
	p.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      "prp_12345",
		PropertyVersion: propertyVersion,
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
	}).Return(getPropertyVersionHostnamesResponse, err).Once()
}

func mockGetCertificate(h *hapi.Mock, dnsZone, recordName string, getCertificateResponse *hapi.GetCertificateResponse, err error) {
	h.On("GetCertificate", mock.Anything, hapi.GetCertificateRequest{
		DNSZone:    dnsZone,
		RecordName: recordName,
	}).Return(getCertificateResponse, err).Once()
}

func mockGetEdgeHostnames(p *papi.Mock) *mock.Call {
	return p.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
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
}

func mockGetEdgeHostname(h *hapi.Mock, hapiGetEdgeHostnameResponse *hapi.GetEdgeHostnameResponse, err error) {
	h.On("GetEdgeHostname", mock.Anything, 2867480).
		Return(hapiGetEdgeHostnameResponse, err).Once()
}

func mockGetActivations(p *papi.Mock, getActivationsResponse *papi.GetActivationsResponse, err error) {
	p.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
		PropertyID: "prp_12345",
		ContractID: "test_contract",
		GroupID:    "grp_12345",
	}).Return(getActivationsResponse, err).Once()
}

func mockAddTemplateTargetRules(t *templates.MockProcessor) {
	t.On("AddTemplateTarget", "rules_v2023-01-05.tmpl", "rules.tf")
}

func mockAddTemplateTargetIncludesRules(t *templates.MockProcessor) {
	t.On("AddTemplateTarget", "includes_rules.tmpl", "includes_rules.tf")
}

func mockTemplateExist(t *templates.MockProcessor, templatePath string, templateExist bool) *mock.Call {
	return t.On("TemplateExists", templatePath).Return(templateExist).Once()
}

func mockProcessTemplates(t *templates.MockProcessor, tfData TFData, filterFuncs []func([]string) ([]string, error), err error) {
	if len(filterFuncs) != 0 {
		t.On("ProcessTemplates", tfData, mock.AnythingOfType("func([]string) ([]string, error)")).Return(err).Once()
	} else {
		t.On("ProcessTemplates", tfData).Return(err).Once()
	}
}

func mockModuleConfig(p *templates.MockProcessor) {
	p.On("AddTemplateTarget", "rules_module.tmpl", "rules/module_config.tf")
}

func generateHostnames(n int, certType, network, ehnID string) []papi.HostnameItem {
	result := make([]papi.HostnameItem, n)
	for i := range n {
		var item papi.HostnameItem
		if network == "STAGING" {
			item = papi.HostnameItem{
				CnameFrom:             fmt.Sprintf("www.test.cname_from.%d.com", i),
				CnameType:             "EDGE_HOSTNAME",
				StagingCertType:       papi.CertType(certType),
				StagingCnameTo:        fmt.Sprintf("www.test.cname_to.%d.com", i),
				StagingEdgeHostnameID: ehnID,
			}
		} else {
			item = papi.HostnameItem{
				CnameFrom:                fmt.Sprintf("www.test.cname_from.%d.com", i),
				CnameType:                "EDGE_HOSTNAME",
				ProductionCertType:       papi.CertType(certType),
				ProductionCnameTo:        fmt.Sprintf("www.test.cname_to.%d.com", i),
				ProductionEdgeHostnameID: ehnID,
			}
		}
		result[i] = item
	}

	return result
}

type activationItemData struct {
	actType    string
	actID      string
	network    string
	status     string
	version    int
	updateDate string
}

type activationItemsData []activationItemData

func mockActivationsItems(items activationItemsData) papi.ActivationsItems {
	var activations []*papi.Activation
	for _, activationItem := range items {
		activation := papi.Activation{
			ActivationID:           activationItem.actID,
			ActivationType:         papi.ActivationType(activationItem.actType),
			PropertyName:           "test.edgesuite.net",
			PropertyID:             "prp_12345",
			PropertyVersion:        activationItem.version,
			AcknowledgeAllWarnings: false,
			Network:                papi.ActivationNetwork(activationItem.network),
			Status:                 papi.ActivationStatus(activationItem.status),
			NotifyEmails:           []string{"aa@aa.com"},
			Note:                   "example note",
			UpdateDate:             activationItem.updateDate,
		}
		activations = append(activations, &activation)
	}
	return papi.ActivationsItems{Items: activations}
}

func getRuleTreeResponse(dir string, t *testing.T) papi.GetRuleTreeResponse {
	var ruleResponse papi.GetRuleTreeResponse
	rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
	assert.NoError(t, err)
	err = json.Unmarshal(rules, &ruleResponse)
	assert.NoError(t, err)
	return ruleResponse
}

func TestProcessPropertyTemplates(t *testing.T) {

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
		rulesAsHCL   bool
		ruleResponse papi.GetRuleTreeResponse
		withError    string
		filterFuncs  []func([]string) ([]string, error)
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with CCM hostnames": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
						"foo.edgesuite.net": {
							CnameFrom:                "foo.edgesuite.net",
							CnameTo:                  "foo",
							EdgeHostnameResourceName: "",
							CertProvisioningType:     "CCM",
							IsActive:                 true,
							CCMCertificates: &CCMCertificates{
								RSACertID:   "123456",
								ECDSACertID: "343434",
							},
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				Section: "test_section",
			},
			dir:          "basic-ccm-hostnames",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with edgehostname with non default ttl": {
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
					ReadVersion:          "LATEST",
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
							TTL:                      600,
						},
					},
					Hostnames: map[string]Hostname{
						"test.edgesuite.net": {
							CnameFrom:                "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic-non-default-ttl",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property enhancement tls": {
			givenData: TFData{
				Property: TFPropertyData{
					GroupName:            "test_group",
					GroupID:              "grp_12345",
					ContractID:           "test_contract",
					PropertyResourceName: "test-edgekey-net",
					PropertyName:         "test.edgekey.net",
					PropertyID:           "prp_12345",
					ProductID:            "prd_HTTP_Content_Del",
					ProductName:          "HTTP_Content_Del",
					RuleFormat:           "latest",
					IsSecure:             "false",
					ReadVersion:          "LATEST",
					EdgeHostnames: map[string]EdgeHostname{
						"test-edgekey-net": {
							EdgeHostname:             "test.edgekey.net",
							EdgeHostnameID:           "ehn_2867480",
							ContractID:               "test_contract",
							GroupID:                  "grp_12345",
							ID:                       "",
							IPv6:                     "IPV6_COMPLIANCE",
							SecurityType:             "ENHANCED-TLS",
							EdgeHostnameResourceName: "test-edgekey-net",
							CertificateID:            123456,
						},
					},
					Hostnames: map[string]Hostname{
						"test.edgekey.net": {
							CnameFrom:                "test.edgekey.net",
							EdgeHostnameResourceName: "test-edgekey-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "enhancement-tls",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"basic property with hostname bucket and no hostname activations": {
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
					HostnameBucket:       &HostnameBucketDetails{},
					ReadVersion:          "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/no-hostnames",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"basic property with hostname bucket and only STAGING hostnames": {
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
					HostnameBucket: &HostnameBucketDetails{
						StagingNotifyEmails:  []string{"test@mail.com"},
						StagingNote:          "staging note",
						HasStagingActivation: true,
						Hostnames:            createTFHostnameItems(generateHostnames(10, "CPS_MANAGED", "STAGING", "ehn_12345"), nil),
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/staging",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and only PRODUCTION hostnames and no emails or note": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						Hostnames:               createTFHostnameItems(nil, generateHostnames(1500, "DEFAULT", "PRODUCTION", "ehn_12345")),
					},
					ReadVersion: "LATEST",
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/production",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and PRODUCTION hostnames with STAGING activation that has no hostnames": {
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
					HostnameBucket: &HostnameBucketDetails{
						ProductionNotifyEmails:  []string{"test@mail.com"},
						ProductionNote:          "production note",
						HasStagingActivation:    false,
						HasProductionActivation: true,
						Hostnames:               createTFHostnameItems(nil, generateHostnames(1, "CPS_MANAGED", "PRODUCTION", "ehn_12345")),
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/production-staging-no-hostnames",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						HasStagingActivation:    true,
						StagingNotifyEmails:     []string{"test@mail.com"},
						ProductionNotifyEmails:  []string{"test@mail.com"},
						StagingNote:             "staging note",
						ProductionNote:          "production note",
						Hostnames: createTFHostnameItems(generateHostnames(20, "DEFAULT", "STAGING", "ehn_12345"),
							generateHostnames(20, "DEFAULT", "PRODUCTION", "ehn_12345")),
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/staging-production",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames with diff": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						HasStagingActivation:    true,
						StagingNotifyEmails:     []string{"test@mail.com"},
						ProductionNotifyEmails:  []string{"test@mail.com"},
						StagingNote:             "staging note",
						ProductionNote:          "production note",
						Hostnames: createTFHostnameItems(generateHostnames(20, "DEFAULT", "STAGING", "ehn_12345"),
							generateHostnames(10, "DEFAULT", "PRODUCTION", "ehn_12345")),
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/staging-production-diff",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames with diff edge hostname id": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						HasStagingActivation:    true,
						StagingNotifyEmails:     []string{"test@mail.com"},
						ProductionNotifyEmails:  []string{"test@mail.com"},
						StagingNote:             "staging note",
						ProductionNote:          "production note",
						Hostnames: map[string]HostnameItem{
							"www.test.cname_from.0.com": {
								CnameTo:                        "www.test.cname_to.0.com",
								CertProvisioningType:           "CPS_MANAGED",
								EdgeHostnameID:                 "ehn_12345",
								Staging:                        true,
								ProductionCertProvisioningType: ptr.To("CPS_MANAGED"),
								ProductionEdgeHostnameID:       ptr.To("ehn_54321"),
							},
						},
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/staging-production-diff-edge-hostname",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"basic property with hostname bucket and both STAGING and PRODUCTION hostnames with diff cert provisioning type": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						HasStagingActivation:    true,
						StagingNotifyEmails:     []string{"test@mail.com"},
						ProductionNotifyEmails:  []string{"test@mail.com"},
						StagingNote:             "staging note",
						ProductionNote:          "production note",
						Hostnames: map[string]HostnameItem{
							"www.test.cname_from.0.com": {
								CnameTo:                        "www.test.cname_to.0.com",
								CertProvisioningType:           "CPS_MANAGED",
								EdgeHostnameID:                 "ehn_12345",
								Staging:                        true,
								ProductionCertProvisioningType: ptr.To("DEFAULT"),
								ProductionEdgeHostnameID:       ptr.To("ehn_12345"),
							},
						},
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "hostname-bucket/staging-production-diff-cert-provisioning-type",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"property with rules as datasource": {
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "ruleformats/basic-rules-datasource",
			rulesAsHCL:   true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"property with rules as datasource with serial": {
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "ruleformats/basic-rules-datasource-serial",
			rulesAsHCL:   true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"property with rules as datasource with unknown behaviors and criteria": {
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
							CnameTo:                  "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation: true,
						Emails:        []string{"jsmith@akamai.com"},
					},
				},
			},
			dir:          "ruleformats/basic-rules-datasource-unknown",
			rulesAsHCL:   true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			withError:    "there were errors reported: Unknown behavior 'caching-unknown', Unknown behavior 'allowPost-unknown'",
		},
		"property with rules as datasource with empty options": {
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "ruleformats/basic-rules-datasource-empty-options",
			rulesAsHCL:   true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2024-01-09")},
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
					ReadVersion:          "3",
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
							UseCases:                 string(useCasesJSON),
						},
					},
					Hostnames: map[string]Hostname{
						"test.edgesuite.net": {
							CnameFrom:                "test.edgesuite.net",
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation: true,
						Emails:        []string{"jsmith@akamai.com"},
					},
				},
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "example staging note",
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic_with_activation_note",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with multiline activation note": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "first\nsecond\n\nlast",
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "first\nsecond\n",
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic_with_multiline_activation_note",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with production activation": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						Emails:         nil,
						ActivationNote: "",
						HasActivation:  false,
					},
					ProductionInfo: NetworkInfo{
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "example production note",
						HasActivation:           true,
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic_with_production_activation",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with both activations - staging and production": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "example staging note",
						HasActivation:           true,
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						Emails:                  []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
						ActivationNote:          "example production note",
						HasActivation:           true,
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "basic_with_both_activations",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property without activation - activation resource commented": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
				},
			},
			dir:          "basic_without_activation",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"property with bootstrap": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation: true,
						Emails:        []string{"jsmith@akamai.com"},
					},
				},
				UseBootstrap: true,
			},
			dir:          "basic-bootstrap",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"basic property with hostname bucket and bootstrap and no hostname activations": {
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
					HostnameBucket:       &HostnameBucketDetails{},
					ReadVersion:          "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				UseBootstrap: true,
			},
			dir:          "hostname-bucket/bootstrap/no-hostnames",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"basic property with hostname bucket and bootstrap and PRODUCTION activation but no hostnames": {
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
					HostnameBucket: &HostnameBucketDetails{
						ProductionNote:          "production note",
						ProductionNotifyEmails:  []string{"test@mail.com"},
						HasProductionActivation: true,
						Hostnames:               nil,
					},
					ReadVersion: "LATEST",
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				UseBootstrap: true,
			},
			dir:          "hostname-bucket/bootstrap/production-no-hostnames",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
		"basic property with hostname bucket and bootstrap and hostnames in STAGING and PRODUCTION": {
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
					HostnameBucket: &HostnameBucketDetails{
						HasProductionActivation: true,
						HasStagingActivation:    true,
						StagingNotifyEmails:     []string{"test@mail.com"},
						ProductionNotifyEmails:  []string{"test@mail.com"},
						StagingNote:             "staging note",
						ProductionNote:          "production note",
						Hostnames: createTFHostnameItems(generateHostnames(20, "DEFAULT", "STAGING", "ehn_12345"),
							generateHostnames(20, "DEFAULT", "PRODUCTION", "ehn_12345")),
					},
					ReadVersion: "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
					ProductionInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				UseBootstrap: true,
			},
			dir:          "hostname-bucket/bootstrap/staging-production",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "hostname_bucket.tf"},
		},
		"property using split-depth": {
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
					ReadVersion:          "LATEST",
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				UseSplitDepth: true,
				RulesAsHCL:    true,
				RootRule:      "test-edgesuite-net",
			},
			dir:          "multitarget-property",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh", "module_config.tf"},
			rulesAsHCL:   true,
		},
		"property with rules having rdns details": {
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
							EdgeHostnameResourceName: "test-edgesuite-net",
							CertProvisioningType:     "CPS_MANAGED",
							IsActive:                 true,
						},
					},
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
			},
			dir:          "ruleformats/rules-with-rdn-details-datasource",
			rulesAsHCL:   true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			filterFuncs:  []func([]string) ([]string, error){useThisOnlyRuleFormat("v2023-01-05")},
		},
		"non default edgerc path and section": {
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
					HostnameBucket:       &HostnameBucketDetails{},
					ReadVersion:          "LATEST",
					StagingInfo: NetworkInfo{
						HasActivation:           true,
						Emails:                  []string{"jsmith@akamai.com"},
						IsActiveOnLatestVersion: true,
					},
				},
				EdgercPath: "/non/default/path/to/edgerc",
				Section:    "non_default_section",
			},
			dir:          "property_non_default_edgerc_path_and_section",
			filesToCheck: []string{"variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.givenData.EdgercPath == "" {
				test.givenData.EdgercPath = defaultEdgercPath
			}
			if test.givenData.Section == "" {
				test.givenData.Section = defaultSection
			}
			if test.rulesAsHCL {
				ruleResponse := getRuleTreeResponse(test.dir, t)
				test.givenData.Rules = flattenRules(wrapAndNameRules("test.edgesuite.net", ruleResponse.Rules))
				test.givenData.RulesAsHCL = true
			}
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			templateToFile := map[string]string{
				"property.tmpl":  fmt.Sprintf("./testdata/res/%s/property.tf", test.dir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
			}
			if test.givenData.Property.HostnameBucket != nil && len(test.givenData.Property.HostnameBucket.Hostnames) != 0 {
				templateToFile["hostname_bucket.tmpl"] = fmt.Sprintf("./testdata/res/%s/hostname_bucket.tf", test.dir)
			}

			if test.rulesAsHCL {
				rulesVersion := test.givenData.Property.RuleFormat
				if !test.givenData.UseSplitDepth {
					templateToFile[fmt.Sprintf("rules_%s.tmpl", rulesVersion)] = fmt.Sprintf("./testdata/res/%s/rules.tf", test.dir)
				} else {
					templateToFile["rules_module.tmpl"] = fmt.Sprintf("./testdata/res/%s/module_config.tf", test.dir)
				}
			}

			processor := templates.FSTemplateProcessor{
				TemplatesFS:     templateFiles,
				TemplateTargets: templateToFile,
				AdditionalFuncs: additionalFuncs,
			}
			err := processor.ProcessTemplates(test.givenData, test.filterFuncs...)
			reportedErrors = []string{}
			if test.withError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.withError)
				return
			}
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

func TestMultiTargetProcessPropertyTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    templates.MultiTargetData
		dir          string
		filesToCheck []string
		withError    string
	}{
		"property with children (split-depth=0)": {
			givenData: templates.MultiTargetData{
				"split-depth-rules.tmpl": {
					"./testdata/res/multitarget-property-with-children/test-edgesuite-net_default.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
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
										FileName:      "test-edgesuite-net_default_new_rule",
										TerraformName: "test-edgesuite-net_rule_new_rule",
										Children: []*WrappedRules{
											{
												Rule:          papi.Rules{Name: "New Rule 1"},
												FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
												TerraformName: "test-edgesuite-net_rule_new_rule_1",
											},
										},
									},
								},
								FileName:      "test-edgesuite-net_default",
								TerraformName: "test-edgesuite-net_rule_default",
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
								Children: []*WrappedRules{
									{
										Rule:          papi.Rules{Name: "New Rule 1"},
										FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
										TerraformName: "test-edgesuite-net_rule_new_rule_1",
									},
								},
								FileName:      "test-edgesuite-net_default_new_rule",
								TerraformName: "test-edgesuite-net_rule_new_rule",
							},
							{
								Rule:          papi.Rules{Name: "New Rule 1"},
								FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
								TerraformName: "test-edgesuite-net_rule_new_rule_1",
							},
						},
					},
				},
			},
			dir:          "multitarget-property-with-children",
			filesToCheck: []string{"test-edgesuite-net_default.tf"},
		},
		"property with children (split-depth=2)": {
			givenData: templates.MultiTargetData{
				"split-depth-rules.tmpl": {
					"./testdata/res/multitarget-property-with-flatten-children/test-edgesuite-net_default.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
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
										FileName:      "test-edgesuite-net_default_new_rule",
										TerraformName: "test-edgesuite-net_rule_new_rule",
										Children: []*WrappedRules{
											{
												Rule:          papi.Rules{Name: "New Rule 1"},
												FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
												TerraformName: "test-edgesuite-net_rule_new_rule_1",
											},
										},
									},
								},
								FileName:      "test-edgesuite-net_default",
								TerraformName: "test-edgesuite-net_rule_default",
							},
						},
					},
					"./testdata/res/multitarget-property-with-flatten-children/test-edgesuite-net_default_new_rule.tf": TFData{
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
								Children: []*WrappedRules{
									{
										Rule:          papi.Rules{Name: "New Rule 1"},
										FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
										TerraformName: "test-edgesuite-net_rule_new_rule_1",
									},
								},
								FileName:      "test-edgesuite-net_default_new_rule",
								TerraformName: "test-edgesuite-net_rule_new_rule",
							},
						},
					},
					"./testdata/res/multitarget-property-with-flatten-children/test-edgesuite-net_default_new_rule_new_rule_1.tf": TFData{
						Rules: []*WrappedRules{
							{
								Rule:          papi.Rules{Name: "New Rule 1"},
								FileName:      "test-edgesuite-net_default_new_rule_new_rule_1",
								TerraformName: "test-edgesuite-net_rule_new_rule_1",
							},
						},
					},
				},
			},
			dir:          "multitarget-property-with-flatten-children",
			filesToCheck: []string{"test-edgesuite-net_default.tf", "test-edgesuite-net_default_new_rule.tf", "test-edgesuite-net_default_new_rule_new_rule_1.tf"},
		},
		"property with unknown behaviors": {
			givenData: templates.MultiTargetData{
				"split-depth-rules.tmpl": {
					"./testdata/res/basic-rules-datasource-unknown/test-edgesuite-net_default.tf": TFData{
						RulesAsHCL:    true,
						UseSplitDepth: true,
						Rules: []*WrappedRules{
							{
								Rule: papi.Rules{
									Name: "default",
									Behaviors: []papi.RuleBehavior{
										{
											Name: "caching-unknown",
										},
										{
											Name: "allowPost-unknown",
										},
										{
											Name: "report",
										},
									},
								},
								FileName:      "test-edgesuite-net_default",
								TerraformName: "test-edgesuite-net_rule_default",
							},
						},
					},
				},
			},
			dir:       "basic-rules-datasource-unknown",
			withError: "there were errors reported: Unknown behavior 'caching-unknown', Unknown behavior 'allowPost-unknown'",
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
			if test.withError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.withError)
				return
			}
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
		"two duplicates different case": {
			given:    "ipcuid invalidation",
			expected: "ipcuid_invalidation1",
			preTest:  []string{"IPCUID Invalidation"},
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

func TestGetLatestActiveActivation(t *testing.T) {
	tests := map[string]struct {
		activationItems       activationItemsData
		network               papi.ActivationNetwork
		expectedActiveVersion int
	}{
		"one staging activation": {
			activationItems: activationItemsData{
				activationItemData{
					actType:    "ACTIVATE",
					actID:      "atv_1",
					status:     "ACTIVE",
					network:    "STAGING",
					version:    1,
					updateDate: "2018-03-07T23:40:45Z"},
			},
			network:               papi.ActivationNetworkStaging,
			expectedActiveVersion: 1,
		},
		"two staging activations": {
			activationItems: activationItemsData{
				activationItemData{
					actType:    "ACTIVATE",
					actID:      "atv_2",
					status:     "ACTIVE",
					network:    "STAGING",
					version:    2,
					updateDate: "2018-03-12T15:40:45Z"},
				activationItemData{
					actType:    "DEACTIVATE",
					actID:      "atv_1",
					status:     "INACTIVE",
					network:    "STAGING",
					version:    1,
					updateDate: "2018-03-07T10:11:45Z"},
			},
			network:               papi.ActivationNetworkStaging,
			expectedActiveVersion: 2,
		},
		"three staging activation - 2nd active": {
			activationItems: activationItemsData{
				activationItemData{
					actType:    "DEACTIVATE",
					actID:      "atv_3",
					status:     "INACTIVE",
					network:    "STAGING",
					version:    3,
					updateDate: "2018-04-01T14:00:45Z"},
				activationItemData{
					actType:    "ACTIVATE",
					actID:      "atv_2",
					status:     "ACTIVE",
					network:    "STAGING",
					version:    2,
					updateDate: "2018-04-20T10:00:45Z"},
				activationItemData{
					actType:    "DEACTIVATE",
					actID:      "atv_1",
					status:     "INACTIVE",
					network:    "STAGING",
					version:    1,
					updateDate: "2018-03-07T08:40:45Z"},
			},
			network:               papi.ActivationNetworkStaging,
			expectedActiveVersion: 2,
		},
		"two production activations": {
			activationItems: activationItemsData{
				activationItemData{
					actType:    "ACTIVATE",
					actID:      "atv_2",
					status:     "ACTIVE",
					network:    "PRODUCTION",
					version:    2,
					updateDate: "2018-03-08T12:00:45Z"},
				activationItemData{
					actType:    "DEACTIVATE",
					actID:      "atv_1",
					status:     "INACTIVE",
					network:    "PRODUCTION",
					version:    1,
					updateDate: "2018-03-07T08:00:45Z"},
			},
			network:               papi.ActivationNetworkProduction,
			expectedActiveVersion: 2,
		},
		"no activations": {
			activationItems: activationItemsData{},
			network:         papi.ActivationNetworkProduction,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			items := mockActivationsItems(test.activationItems)
			var activation = getLatestActiveActivation(items, test.network)
			if activation != nil {
				assert.Equal(t, test.expectedActiveVersion, activation.PropertyVersion)
			} else {
				assert.Empty(t, test.expectedActiveVersion)
			}
		})
	}
}

func TestTerraformName(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"with spaces": {
			given:    "this is test name",
			expected: "this_is_test_name",
		},
		"strange characters": {
			given:    "test@name1ą",
			expected: "test-name1ą",
		},
		"Deny by Location": {
			given:    "Deny by Location",
			expected: "deny_by_location",
		},
		"redirect to language specific section": {
			given:    "redirect to language specific section",
			expected: "redirect_to_language_specific_section",
		},
		"mPulse": {
			given:    "mPulse",
			expected: "m_pulse",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, tools.TerraformName(test.given))
		})
	}
}

type tfDataBuilder struct {
	tfData TFData
}

func (t *tfDataBuilder) withDefaults() *tfDataBuilder {
	t.tfData = TFData{
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
		EdgercPath: defaultEdgercPath,
		Section:    defaultSection,
	}
	return t
}

func (t *tfDataBuilder) withDefaultsHavingDigitsAndSpacesInHostnameDetails() *tfDataBuilder {
	t.tfData = TFData{
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
				"_1_test-edgesuite-net": {
					EdgeHostname:             "1 test.edgesuite.net",
					EdgeHostnameID:           "ehn_2867480",
					ContractID:               "test_contract",
					GroupID:                  "grp_12345",
					ID:                       "",
					IPv6:                     "IPV6_COMPLIANCE",
					SecurityType:             "STANDARD-TLS",
					EdgeHostnameResourceName: "_1_test-edgesuite-net",
				},
			},
			Hostnames: map[string]Hostname{
				"test.edgesuite.net": {
					CnameFrom:                "test.edgesuite.net",
					CnameTo:                  "1 test.edgesuite.net",
					EdgeHostnameResourceName: "_1_test-edgesuite-net",
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
		EdgercPath: defaultEdgercPath,
		Section:    defaultSection,
	}
	return t
}

func (t *tfDataBuilder) withData(data TFData) *tfDataBuilder {
	t.tfData = data
	return t
}

func (t *tfDataBuilder) withRuleFormat(ruleFormat string) *tfDataBuilder {
	t.tfData.Property.RuleFormat = ruleFormat
	return t
}

func (t *tfDataBuilder) withHostnames(hostnames map[string]Hostname) *tfDataBuilder {
	t.tfData.Property.Hostnames = hostnames
	return t
}

func (t *tfDataBuilder) withEdgeHostname(edgeHostname map[string]EdgeHostname) *tfDataBuilder {
	t.tfData.Property.EdgeHostnames = edgeHostname
	return t
}

func (t *tfDataBuilder) withHostnameBucket(staging, production []papi.HostnameItem, stagingNote, productionNote string, stagingMails, productionMails []string) *tfDataBuilder {
	if len(staging) == 0 && len(production) == 0 {
		t.tfData.Property.HostnameBucket = &HostnameBucketDetails{}
		return t
	}
	hb := &HostnameBucketDetails{
		StagingNotifyEmails:     stagingMails,
		ProductionNotifyEmails:  productionMails,
		StagingNote:             stagingNote,
		ProductionNote:          productionNote,
		HasStagingActivation:    len(staging) != 0,
		HasProductionActivation: len(production) != 0,
		Hostnames:               createTFHostnameItems(staging, production),
	}
	t.tfData.Property.HostnameBucket = hb

	return t
}

func (t *tfDataBuilder) withIsActive(isActive bool) *tfDataBuilder {
	hostname := t.tfData.Property.Hostnames["test.edgesuite.net"]
	hostname.IsActive = isActive
	t.tfData.Property.Hostnames["test.edgesuite.net"] = hostname
	return t
}

func (t *tfDataBuilder) withEmails(emails []string) *tfDataBuilder {
	t.tfData.Property.StagingInfo.Emails = emails
	return t
}

func (t *tfDataBuilder) withOnlyProductionActivation(emails []string, activationNote string, version int) *tfDataBuilder {

	t.withProductionVersion(version, true)
	t.tfData.Property.ProductionInfo.Emails = emails
	t.tfData.Property.ProductionInfo.ActivationNote = activationNote
	t.tfData.Property.StagingInfo.HasActivation = false
	t.tfData.Property.StagingInfo.ActivationNote = ""
	t.tfData.Property.StagingInfo.Version = 0
	t.tfData.Property.StagingInfo.Emails = nil
	t.tfData.Property.StagingInfo.IsActiveOnLatestVersion = false
	return t
}

func (t *tfDataBuilder) withVersion(version string) *tfDataBuilder {
	t.tfData.Property.ReadVersion = version
	return t
}

func (t *tfDataBuilder) withStagingVersion(version int, isActiveOnLatestVersion bool) *tfDataBuilder {
	t.tfData.Property.StagingInfo.HasActivation = true
	t.tfData.Property.StagingInfo.Version = version
	t.tfData.Property.StagingInfo.IsActiveOnLatestVersion = isActiveOnLatestVersion
	return t
}

func (t *tfDataBuilder) withProductionVersion(version int, isActiveOnLatestVersion bool) *tfDataBuilder {
	t.tfData.Property.ProductionInfo.HasActivation = true
	t.tfData.Property.ProductionInfo.Version = version
	t.tfData.Property.ProductionInfo.IsActiveOnLatestVersion = isActiveOnLatestVersion
	return t
}

func (t *tfDataBuilder) withRules(rules []*WrappedRules) *tfDataBuilder {
	t.tfData.Rules = rules
	t.tfData.RulesAsHCL = true
	return t
}

func (t *tfDataBuilder) withIncludeRules(index int, rules []*WrappedRules) *tfDataBuilder {
	t.tfData.Includes[index].Rules = rules
	t.tfData.RulesAsHCL = true
	return t
}

func (t *tfDataBuilder) withSplitDepth(useSplitDepth bool, rootRule string) *tfDataBuilder {
	t.tfData.UseSplitDepth = useSplitDepth
	t.tfData.RootRule = rootRule
	return t
}

func (t *tfDataBuilder) withCertProvisioningType(certProvisioningType string) *tfDataBuilder {
	hostname := t.tfData.Property.Hostnames["test.edgesuite.net"]
	hostname.CertProvisioningType = certProvisioningType
	t.tfData.Property.Hostnames["test.edgesuite.net"] = hostname
	return t
}

func (t *tfDataBuilder) withActivationNote(activationNote string) *tfDataBuilder {
	t.tfData.Property.StagingInfo.ActivationNote = activationNote
	return t
}

func (t *tfDataBuilder) withIncludes(includes []TFIncludeData) *tfDataBuilder {
	t.tfData.Includes = includes
	return t
}

func (t *tfDataBuilder) withBootstrap(useBootstrap bool) *tfDataBuilder {
	t.tfData.UseBootstrap = useBootstrap
	return t
}

func (t *tfDataBuilder) asHCL(useHCL bool) *tfDataBuilder {
	t.tfData.RulesAsHCL = useHCL
	return t
}

func (t *tfDataBuilder) withEdgercPathAndSection(edgercPath string, edgercSection string) *tfDataBuilder {
	t.tfData.EdgercPath = edgercPath
	t.tfData.Section = edgercSection
	return t
}

func (t *tfDataBuilder) build() TFData {
	return t.tfData
}

func TestUseThisOnlyRuleFormat(t *testing.T) {
	tests := map[string]struct {
		acceptedFormat string
		input          []string
		expected       []string
		err            error
	}{
		"happy case": {
			acceptedFormat: "v2023-05-30",
			input:          []string{"import.sh", "templates/rules_v2023-01-05.tmpl", "templates/rules_v2023-05-30.tmpl"},
			expected:       []string{"import.sh", "templates/rules_v2023-05-30.tmpl"},
		},
		"missing template": {
			acceptedFormat: "v2023-05-32",
			input:          []string{"import.sh", "templates/rules_v2023-01-05.tmpl", "templates/rules_v2023-05-30.tmpl"},
			err:            fmt.Errorf("did not find v2023-05-32 format among [import.sh templates/rules_v2023-01-05.tmpl templates/rules_v2023-05-30.tmpl]"),
		},
		"no formats at all": {
			acceptedFormat: "v2023-05-32",
			input:          []string{"import.sh", "READ.md"},
			err:            fmt.Errorf("did not find v2023-05-32 format among [import.sh READ.md]"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := useThisOnlyRuleFormat(test.acceptedFormat)(test.input)
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, output)
		})
	}
}

func newTimeFromString(t *testing.T, s string) *time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	require.NoError(t, err)
	return &parsedTime
}

func TestCreateTFHostnameItems(t *testing.T) {
	tests := map[string]struct {
		staging    []papi.HostnameItem
		production []papi.HostnameItem
		expected   map[string]HostnameItem
	}{
		"only staging - CPS_MANAGED": {
			staging: generateHostnames(2, "CPS_MANAGED", "STAGING", "ehn_12345"),
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:              "www.test.cname_to.0.com",
					CertProvisioningType: "CPS_MANAGED",
					EdgeHostnameID:       "ehn_12345",
					Staging:              true,
					Production:           false,
				},
				"www.test.cname_from.1.com": {
					CnameTo:              "www.test.cname_to.1.com",
					CertProvisioningType: "CPS_MANAGED",
					EdgeHostnameID:       "ehn_12345",
					Staging:              true,
					Production:           false,
				},
			},
		},
		"only production - DEFAULT": {
			production: generateHostnames(2, "DEFAULT", "PRODUCTION", "ehn_12345"),
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:              "www.test.cname_to.0.com",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_12345",
					Staging:              false,
					Production:           true,
				},
				"www.test.cname_from.1.com": {
					CnameTo:              "www.test.cname_to.1.com",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_12345",
					Staging:              false,
					Production:           true,
				},
			},
		},
		"staging and production - no diff between networks - DEFAULT and ehn_12345 - no additional production fields filled": {
			staging:    generateHostnames(2, "DEFAULT", "STAGING", "ehn_12345"),
			production: generateHostnames(2, "DEFAULT", "PRODUCTION", "ehn_12345"),
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:              "www.test.cname_to.0.com",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_12345",
					Staging:              true,
					Production:           true,
				},
				"www.test.cname_from.1.com": {
					CnameTo:              "www.test.cname_to.1.com",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_12345",
					Staging:              true,
					Production:           true,
				},
			},
		},
		"staging and production - diff in cert_provisioning_type between networks - expect additional production fields to be filled": {
			staging:    generateHostnames(2, "DEFAULT", "STAGING", "ehn_12345"),
			production: generateHostnames(2, "CPS_MANAGED", "PRODUCTION", "ehn_12345"),
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:                        "www.test.cname_to.0.com",
					CertProvisioningType:           "DEFAULT",
					EdgeHostnameID:                 "ehn_12345",
					Staging:                        true,
					Production:                     false,
					ProductionCertProvisioningType: ptr.To("CPS_MANAGED"),
					ProductionEdgeHostnameID:       ptr.To("ehn_12345"),
				},
				"www.test.cname_from.1.com": {
					CnameTo:                        "www.test.cname_to.1.com",
					CertProvisioningType:           "DEFAULT",
					EdgeHostnameID:                 "ehn_12345",
					Staging:                        true,
					Production:                     false,
					ProductionCertProvisioningType: ptr.To("CPS_MANAGED"),
					ProductionEdgeHostnameID:       ptr.To("ehn_12345"),
				},
			},
		},
		"staging and production - diff in edge_hostname_id between networks - expect additional production fields to be filled": {
			staging:    generateHostnames(2, "DEFAULT", "STAGING", "ehn_67890"),
			production: generateHostnames(2, "DEFAULT", "PRODUCTION", "ehn_12345"),
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:                        "www.test.cname_to.0.com",
					CertProvisioningType:           "DEFAULT",
					EdgeHostnameID:                 "ehn_67890",
					Staging:                        true,
					Production:                     false,
					ProductionCertProvisioningType: ptr.To("DEFAULT"),
					ProductionEdgeHostnameID:       ptr.To("ehn_12345"),
				},
				"www.test.cname_from.1.com": {
					CnameTo:                        "www.test.cname_to.1.com",
					CertProvisioningType:           "DEFAULT",
					EdgeHostnameID:                 "ehn_67890",
					Staging:                        true,
					Production:                     false,
					ProductionCertProvisioningType: ptr.To("DEFAULT"),
					ProductionEdgeHostnameID:       ptr.To("ehn_12345"),
				},
			},
		},
		"staging and production - first hostname is the same for both networks - no additional production fields filled, " +
			"diff in edge_hostname_id between networks for the second hostname - expect additional production fields to be filled": {
			staging: []papi.HostnameItem{
				{
					CnameFrom:             "www.test.cname_from.0.com",
					StagingCertType:       "CPS_MANAGED",
					StagingCnameTo:        "www.test.cname_to.0.com",
					StagingEdgeHostnameID: "ehn_12345",
				},
				{
					CnameFrom:             "www.test.cname_from.1.com",
					StagingCertType:       "DEFAULT",
					StagingCnameTo:        "www.test.cname_to.1.com",
					StagingEdgeHostnameID: "ehn_12345",
				},
			},
			production: []papi.HostnameItem{
				{
					CnameFrom:                "www.test.cname_from.0.com",
					ProductionCertType:       "CPS_MANAGED",
					ProductionCnameTo:        "www.test.cname_to.0.com",
					ProductionEdgeHostnameID: "ehn_12345",
				},
				{
					CnameFrom:                "www.test.cname_from.1.com",
					ProductionCertType:       "DEFAULT",
					ProductionCnameTo:        "www.test.cname_to.1.com",
					ProductionEdgeHostnameID: "ehn_67890",
				},
			},
			expected: map[string]HostnameItem{
				"www.test.cname_from.0.com": {
					CnameTo:              "www.test.cname_to.0.com",
					CertProvisioningType: "CPS_MANAGED",
					EdgeHostnameID:       "ehn_12345",
					Staging:              true,
					Production:           true,
				},
				"www.test.cname_from.1.com": {
					CnameTo:                        "www.test.cname_to.1.com",
					CertProvisioningType:           "DEFAULT",
					EdgeHostnameID:                 "ehn_12345",
					Staging:                        true,
					Production:                     false,
					ProductionCertProvisioningType: ptr.To("DEFAULT"),
					ProductionEdgeHostnameID:       ptr.To("ehn_67890"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := createTFHostnameItems(tc.staging, tc.production)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
