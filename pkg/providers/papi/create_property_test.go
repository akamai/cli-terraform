package papi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
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

	getListReferencedMultipleIncludesResponse := papi.ListReferencedIncludesResponse{
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
				{
					AccountID:      "test_account",
					AssetID:        "test_asset",
					ContractID:     "test_contract",
					GroupID:        "test_group",
					IncludeID:      "inc_78910",
					IncludeName:    "test_include_1",
					IncludeType:    papi.IncludeTypeMicroServices,
					LatestVersion:  2,
					StagingVersion: tools.IntPtr(1),
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

	hapiGetEdgeHostnameResponse := hapi.GetEdgeHostnameResponse{
		EdgeHostnameID:    2867480,
		RecordName:        "test",
		DNSZone:           "edgesuite.net",
		SecurityType:      "STANDARD-TLS",
		UseDefaultTTL:     false,
		UseDefaultMap:     false,
		IPVersionBehavior: "IPV6_IPV4_DUALSTACK",
		TTL:               21600,
		Map:               "a;test.akamai.net",
		SerialNumber:      1461,
	}

	tfIncludeData := TFIncludeData{
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
	}

	tfIncludeData1 := TFIncludeData{
		ActivationNoteStaging:   "test staging activation",
		ContractID:              "test_contract",
		ActivationEmailsStaging: []string{"test@example.com"},
		GroupID:                 "test_group",
		IncludeID:               "inc_78910",
		IncludeName:             "test_include_1",
		IncludeType:             string(papi.IncludeTypeMicroServices),
		Networks:                []string{"STAGING"},
		RuleFormat:              "v2020-11-02",
		VersionStaging:          "1",
	}

	tests := map[string]struct {
		init                func(*papi.Mock, *hapi.Mock, *templates.MockProcessor, string)
		dir                 string
		snippetFilesToCheck []string
		jsonDir             string
		withError           error
		readVersion         string
		withIncludes        bool
		schema              bool
	}{
		"basic property": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().build(), nil)
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
		"basic property with empty hostname id": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

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
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withEdgeHostname(map[string]EdgeHostname{}).
					withIsActive(false).build(), nil)
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
		"basic property with rules as datasource": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockAddTemplateTarget(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withRuleFormat("v2023-01-05").
					withRules(flattenRules("test.edgesuite.net", ruleResponse.Rules)).build(), nil)
			},
			dir:    "basic-rules-datasource",
			schema: true,
		},
		"basic property with rules as datasource with unsupported rule format": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockAddTemplateTarget(p)
				mockTemplateExist(p, "rules_latest.tmpl", false)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withIsActive(false).
					withRules(flattenRules("test.edgesuite.net", ruleResponse.Rules)).build(), nil)
			},
			withError: ErrUnsupportedRuleFormat,
			dir:       "basic",
			schema:    true,
		},
		"basic property with rules as datasource with unsupported behaviors and criteria": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockAddTemplateTarget(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withRuleFormat("v2023-01-05").
					withRules(flattenRules("test.edgesuite.net", ruleResponse.Rules)).build(), ErrSavingFiles)
			},
			withError: ErrSavingFiles,
			dir:       "basic-rules-datasource-unknown",
			schema:    true,
		},
		"basic property with include": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)

				// Includes
				mockListReferencedIncludes(c, &getListReferencedIncludesResponse)
				expectGetIncludeVersion(c)
				var includeRuleResponse papi.GetIncludeRuleTreeResponse
				includeRules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_include_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(includeRules, &includeRuleResponse))
				mockGetIncludeRuleTree(c, getIncludeRuleTreeReq, &includeRuleResponse)
				expectListIncludeActivations(c)

				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withIncludes([]TFIncludeData{tfIncludeData}).build(), nil)
			},
			dir:     "basic_property_with_include",
			jsonDir: "basic_property_with_include/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"test_include.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			withIncludes: true,
		},
		"basic property with multiple includes": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)

				// Includes
				mockListReferencedIncludes(c, &getListReferencedMultipleIncludesResponse)
				expectGetIncludeVersion(c)

				var includeRuleResponse papi.GetIncludeRuleTreeResponse
				includeRules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_include_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(includeRules, &includeRuleResponse))
				mockGetIncludeRuleTree(c, getIncludeRuleTreeReq, &includeRuleResponse)
				expectListIncludeActivations(c)

				var secondIncludeRuleResponse papi.GetIncludeRuleTreeResponse
				expectGetSeconIncludeVersion(c)
				secondIncludeRules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_second_include_rules.json"))
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(secondIncludeRules, &secondIncludeRuleResponse))

				mockGetIncludeRuleTree(c, papi.GetIncludeRuleTreeRequest{
					ContractID:     "test_contract",
					GroupID:        "test_group",
					IncludeID:      "inc_78910",
					IncludeVersion: 2,
					RuleFormat:     "v2020-11-02",
				}, &secondIncludeRuleResponse)
				expectListSecondIncludeActivations(c)

				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withIncludes([]TFIncludeData{tfIncludeData, tfIncludeData1}).build(), nil)
			},
			dir:     "basic_property_with_multiple_includes",
			jsonDir: "basic_property_with_multiple_includes/property-snippets",
			snippetFilesToCheck: []string{
				"main.json",
				"test_include.json",
				"test_include_1.json",
				"Content_Compression.json",
				"Static_Content.json",
				"Dynamic_Content.json",
			},
			withIncludes: true,
		},
		"basic property with cert provisioning type": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersion2HostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withCertProvisioningType("DEFAULT").build(), nil)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().build(), nil)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 1, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 1, &getPropertyVersion1HostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivations1Response)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withVersion("1").build(), nil)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponseWithNote)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withActivationNote("example note").
					withEmails([]string{"jsmith@akamai.com", "rjohnson@akamai.com"}).build(), nil)
			},
			dir: "basic",
		},
		"property activation with empty emails": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponseWithEmptyEmails)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().withActivationNote("example note").
					withEmails([]string{""}).build(), nil)
			},
			dir: "basic",
		},
		"error property not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, nil, fmt.Errorf("oops"))
			},
			withError: ErrPropertyNotFound,
		},
		"error group not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, nil, fmt.Errorf("oops"))
			},
			withError: ErrGroupNotFound,
		},
		"error property rules not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)

				mockGetRuleTree(c, 5, nil, fmt.Errorf("oops"))

			},
			withError: ErrPropertyRulesNotFound,
		},
		"error property version not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, nil, fmt.Errorf("oops"))

			},
			withError: ErrPropertyVersionNotFound,
		},
		"error product name not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
				mockGetRuleTree(c, 5, &papi.GetRuleTreeResponse{}, nil)
				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, nil, fmt.Errorf("oops"))

			},
			withError: ErrProductNameNotFound,
		},
		"error hostnames not found": {
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)
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
			init: func(c *papi.Mock, h *hapi.Mock, p *templates.MockProcessor, dir string) {
				mockSearchProperties(c, &searchPropertiesResponse, nil)
				mockGetProperty(c, &getPropertyResponse)

				ruleResponse := getRuleTreeResponse(dir, t)
				mockGetRuleTree(c, 5, &ruleResponse, nil)

				mockGetGroups(c, &getGroupsResponse, nil)
				mockGetPropertyVersions(c, &getPropertyVersionsResponse, nil)
				mockGetLatestVersion(c, &getLatestVersionResponse)
				mockGetProducts(c, &getProductsResponse, nil)
				mockGetPropertyVersionHostnames(c, 5, &getPropertyVersionHostnamesResponse, nil)
				mockGetEdgeHostname(h, &hapiGetEdgeHostnameResponse, nil)
				mockGetEdgeHostnames(c)
				mockGetActivations(c, &getActivationsResponse)
				mockAddTemplateTarget(p)
				mockTemplateExist(p, "rules_v2023-01-05.tmpl", true)
				mockProcessTemplates(p, (&tfDataBuilder{}).withDefaults().build(), fmt.Errorf("oops"))
			},
			dir:       "basic",
			withError: ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mc := new(papi.Mock)
			mh := new(hapi.Mock)
			mp := new(templates.MockProcessor)
			test.init(mc, mh, mp, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createProperty(ctx, "test.edgesuite.net", test.readVersion, section, fmt.Sprintf("./testdata/res/%s", test.jsonDir), "./", test.withIncludes, test.schema, mc, mh, mp)
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

func mockSearchProperties(c *papi.Mock, searchPropertiesResponse *papi.SearchResponse, err error) {
	c.On("SearchProperties", mock.Anything, papi.SearchRequest{Key: "propertyName", Value: "test.edgesuite.net"}).
		Return(searchPropertiesResponse, err).Once()
}

func mockGetProperty(c *papi.Mock, getPropertyResponse *papi.GetPropertyResponse) {
	c.On("GetProperty", mock.Anything, papi.GetPropertyRequest{ContractID: "test_contract", GroupID: "grp_12345", PropertyID: "prp_12345"}).
		Return(getPropertyResponse, nil).Once()
}

func mockGetRuleTree(c *papi.Mock, propertyVersion int, ruleResponse *papi.GetRuleTreeResponse, err error) {
	c.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{PropertyID: "prp_12345", PropertyVersion: propertyVersion, ContractID: "test_contract", GroupID: "grp_12345", ValidateMode: "", ValidateRules: true, RuleFormat: "latest"}).
		Return(ruleResponse, err).Once()
}

func mockGetGroups(c *papi.Mock, getGroupsResponse *papi.GetGroupsResponse, err error) {
	c.On("GetGroups", mock.Anything).
		Return(getGroupsResponse, err).Once()
}

func mockGetPropertyVersions(c *papi.Mock, getPropertyVersionsResponse *papi.GetPropertyVersionsResponse, err error) {
	c.On("GetPropertyVersions", mock.Anything, papi.GetPropertyVersionsRequest{
		PropertyID: "prp_12345",
		ContractID: "test_contract",
		GroupID:    "grp_12345",
	}).Return(getPropertyVersionsResponse, err).Once()
}

func mockGetLatestVersion(c *papi.Mock, getLatestVersionResponse *papi.GetPropertyVersionsResponse) {
	c.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
		PropertyID:  "prp_12345",
		ActivatedOn: "",
		ContractID:  "test_contract",
		GroupID:     "grp_12345",
	}).Return(getLatestVersionResponse, nil).Once()
}

func mockGetProducts(c *papi.Mock, getProductsResponse *papi.GetProductsResponse, err error) {
	c.On("GetProducts", mock.Anything, papi.GetProductsRequest{ContractID: "test_contract"}).
		Return(getProductsResponse, err).Once()
}

func mockGetPropertyVersionHostnames(c *papi.Mock, propertyVersion int, getPropertyVersionHostnamesResponse *papi.GetPropertyVersionHostnamesResponse, err error) {
	c.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      "prp_12345",
		PropertyVersion: propertyVersion,
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
	}).Return(getPropertyVersionHostnamesResponse, err).Once()
}

func mockGetEdgeHostnames(c *papi.Mock) *mock.Call {
	return c.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
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

func mockGetActivations(c *papi.Mock, getActivationsResponse *papi.GetActivationsResponse) {
	c.On("GetActivations", mock.Anything, papi.GetActivationsRequest{
		PropertyID: "prp_12345",
		ContractID: "test_contract",
		GroupID:    "grp_12345",
	}).Return(getActivationsResponse, nil).Once()
}

func mockAddTemplateTarget(p *templates.MockProcessor) {
	p.On("AddTemplateTarget", "rules_v2023-01-05.tmpl", "rules.tf")
}

func mockTemplateExist(p *templates.MockProcessor, templatePath string, templateExist bool) *mock.Call {
	return p.On("TemplateExists", templatePath).Return(templateExist).Once()
}

func mockListReferencedIncludes(c *papi.Mock, getListReferencedIncludesResponse *papi.ListReferencedIncludesResponse) {
	c.On("ListReferencedIncludes", mock.Anything, papi.ListReferencedIncludesRequest{
		PropertyID:      "prp_12345",
		ContractID:      "test_contract",
		GroupID:         "grp_12345",
		PropertyVersion: 5,
	}).Return(getListReferencedIncludesResponse, nil).Once()
}

func mockGetIncludeRuleTree(c *papi.Mock, getIncludeRuleTreeReq papi.GetIncludeRuleTreeRequest, includeRuleResponse *papi.GetIncludeRuleTreeResponse) {
	c.On("GetIncludeRuleTree", mock.Anything, getIncludeRuleTreeReq).Return(includeRuleResponse, nil).Once()
}

func mockProcessTemplates(p *templates.MockProcessor, tfData TFData, err error) {
	p.On("ProcessTemplates", tfData).Return(err).Once()
}

func getRuleTreeResponse(dir string, t *testing.T) papi.GetRuleTreeResponse {
	var ruleResponse papi.GetRuleTreeResponse
	rules, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", dir, "mock_rules.json"))
	assert.NoError(t, err)
	err = json.Unmarshal(rules, &ruleResponse)
	assert.NoError(t, err)
	return ruleResponse
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
		schema       bool
		ruleResponse papi.GetRuleTreeResponse
		withError    string
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
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
					Version:              "LATEST",
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic-rules-datasource",
			schema:       true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
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
					Version:              "LATEST",
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic-rules-datasource-unknown",
			schema:       true,
			filesToCheck: []string{"property.tf", "rules.tf", "variables.tf", "import.sh"},
			withError:    "there were errors reported: Unknown behavior 'caching-unknown', Unknown behavior 'allowPost-unknown'",
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic_property_with_include",
			filesToCheck: []string{"property.tf", "includes.tf", "variables.tf", "import.sh"},
			withIncludes: true,
		},
		"property with multiple includes": {
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
					{
						ActivationNoteStaging:   "test staging activation",
						ContractID:              "test_contract",
						ActivationEmailsStaging: []string{"test@example.com"},
						GroupID:                 "test_group",
						IncludeID:               "inc_78910",
						IncludeName:             "test_include_1",
						IncludeType:             string(papi.IncludeTypeMicroServices),
						Networks:                []string{"STAGING"},
						RuleFormat:              "v2020-11-02",
						VersionStaging:          "1",
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic_property_with_multiple_includes",
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
					Emails:               []string{"jsmith@akamai.com"},
					HasStagingActivation: true,
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
					Emails:               []string{"jsmith@akamai.com", "rjohnson@akamai.com"},
					ActivationNote:       "example note",
					HasStagingActivation: true,
				},
				Section: "test_section",
			},
			dir:          "basic_with_activation_note",
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
					Version:              "LATEST",
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
				Section: "test_section",
			},
			dir:          "basic_without_activation",
			filesToCheck: []string{"property.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.schema {
				ruleResponse := getRuleTreeResponse(test.dir, t)
				test.givenData.Rules = flattenRules("test.edgesuite.net", ruleResponse.Rules)
				test.givenData.RulesAsSchema = true
			}
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			templateToFile := map[string]string{
				"property.tmpl":  fmt.Sprintf("./testdata/res/%s/property.tf", test.dir),
				"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
			}

			if test.withIncludes {
				templateToFile["includes.tmpl"] = fmt.Sprintf("./testdata/res/%s/includes.tf", test.dir)
			}
			if test.schema {
				templateToFile["rules_v2023-01-05.tmpl"] = fmt.Sprintf("./testdata/res/%s/rules.tf", test.dir)
			}
			processor := templates.FSTemplateProcessor{
				TemplatesFS:     templateFiles,
				TemplateTargets: templateToFile,
				AdditionalFuncs: additionalFuncs,
			}
			err := processor.ProcessTemplates(test.givenData)
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
			assert.Equal(t, test.expected, TerraformName(test.given))
		})
	}
}

func TestIsMultiline(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected bool
	}{
		"has new lines": {
			given:    "this\nis test\n",
			expected: true,
		},
		"empty": {
			given:    "",
			expected: false,
		},
		"no new lines": {
			given:    "no new lines",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsMultiline(test.given))
		})
	}
}

func TestNoNewlineAtTheEnd(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected bool
	}{
		"has new line at the end": {
			given:    "this\nis test\n",
			expected: false,
		},
		"has new line in the middle": {
			given:    "this\nis test",
			expected: true,
		},
		"empty": {
			given:    "",
			expected: true,
		},
		"no new lines": {
			given:    "no new lines",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NoNewlineAtTheEnd(test.given))
		})
	}
}

func TestRemoveLastNewline(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"has new line at the end": {
			given:    "this\nis test\n",
			expected: "this\nis test",
		},
		"has new line in the middle": {
			given:    "this\nis test",
			expected: "this\nis test",
		},
		"empty": {
			given:    "",
			expected: "",
		},
		"no new lines": {
			given:    "no new lines",
			expected: "no new lines",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, RemoveLastNewline(test.given))
		})
	}
}

func TestGetEOT(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"has new line": {
			given:    "this\nis test\n",
			expected: "EOT",
		},
		"has EOT inside": {
			given:    "this\nEOT",
			expected: "EOTA",
		},
		"empty": {
			given:    "",
			expected: "EOT",
		},
		"has two EOTs": {
			given:    "some\nEOT\nEOTA\ntext",
			expected: "EOTAA",
		},
		"has EOT": {
			given:    "comment\nnewline\nand\nEOT\ninside\n",
			expected: "EOTA",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, GetEOT(test.given))
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
			Emails:               []string{"jsmith@akamai.com"},
			Version:              "LATEST",
			HasStagingActivation: true,
		},
		Section: "test_section",
	}
	return t
}

func (t *tfDataBuilder) withRuleFormat(ruleFormat string) *tfDataBuilder {
	t.tfData.Property.RuleFormat = ruleFormat
	return t
}

func (t *tfDataBuilder) withEdgeHostname(edgeHostname map[string]EdgeHostname) *tfDataBuilder {
	t.tfData.Property.EdgeHostnames = edgeHostname
	return t
}

func (t *tfDataBuilder) withIsActive(isActive bool) *tfDataBuilder {
	hostname := t.tfData.Property.Hostnames["test.edgesuite.net"]
	hostname.IsActive = isActive
	t.tfData.Property.Hostnames["test.edgesuite.net"] = hostname
	return t
}

func (t *tfDataBuilder) withEmails(emails []string) *tfDataBuilder {
	t.tfData.Property.Emails = emails
	return t
}

func (t *tfDataBuilder) withVersion(version string) *tfDataBuilder {
	t.tfData.Property.Version = version
	return t
}

func (t *tfDataBuilder) withRules(rules []*WrappedRules) *tfDataBuilder {
	t.tfData.Rules = rules
	t.tfData.RulesAsSchema = true
	return t
}

func (t *tfDataBuilder) withCertProvisioningType(certProvisioningType string) *tfDataBuilder {
	hostname := t.tfData.Property.Hostnames["test.edgesuite.net"]
	hostname.CertProvisioningType = certProvisioningType
	t.tfData.Property.Hostnames["test.edgesuite.net"] = hostname
	return t
}

func (t *tfDataBuilder) withActivationNote(activationNote string) *tfDataBuilder {
	t.tfData.Property.ActivationNote = activationNote
	return t
}

func (t *tfDataBuilder) withIncludes(includes []TFIncludeData) *tfDataBuilder {
	t.tfData.Includes = includes
	return t
}

func (t *tfDataBuilder) build() TFData {
	return t.tfData
}
