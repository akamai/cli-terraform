package cloudaccess

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudaccess"
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

var (
	processor = func(testDir string) templates.FSTemplateProcessor {
		return templates.FSTemplateProcessor{
			TemplatesFS:     templateFiles,
			AdditionalFuncs: additionalFunctions,
			TemplateTargets: map[string]string{
				"cloudaccess.tmpl": fmt.Sprintf("./testdata/res/%s/cloudaccess.tf", testDir),
				"variables.tmpl":   fmt.Sprintf("./testdata/res/%s/variables.tf", testDir),
				"imports.tmpl":     fmt.Sprintf("./testdata/res/%s/import.sh", testDir),
			},
		}
	}
	chinaCDN     = cloudaccess.ChinaCDN
	accessKeyUID = int64(1)
	section      = "test_section"
)

func TestCreateCloudAccess(t *testing.T) {
	tests := map[string]struct {
		init      func(*cloudaccess.Mock, *templates.MockProcessor, string, testDataForCloudAccess)
		dir       string
		testData  testDataForCloudAccess
		withError error
	}{
		"access key with two versions and all data": {
			init: func(c *cloudaccess.Mock, p *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
				mockProcessTemplates(p, data)
			},
			dir: "two-versions",
			testData: testDataForCloudAccess{
				section:              section,
				accessKeyUID:         1,
				accessKeyName:        "name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"ctr_111"},
						GroupID:     11,
						GroupName:   tools.StringPtr("group11"),
					},
				},
				accessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key2"),
					},
				},
			},
		},
		"access key with one version": {
			init: func(c *cloudaccess.Mock, p *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
				mockProcessTemplates(p, data)
			},
			dir: "one-version",
			testData: testDataForCloudAccess{
				section:              section,
				accessKeyUID:         1,
				accessKeyName:        "Test name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"ctr_111"},
						GroupID:     11,
						GroupName:   tools.StringPtr("group11"),
					},
				},
				accessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
				},
			},
		},
		"access key with no versions": {
			init: func(c *cloudaccess.Mock, p *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
				mockProcessTemplates(p, data)
			},
			dir: "no-versions",
			testData: testDataForCloudAccess{
				section:              section,
				accessKeyUID:         1,
				accessKeyName:        "name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"ctr_111"},
						GroupID:     11,
						GroupName:   tools.StringPtr("group11"),
					},
				},
			},
		},
		"error getting access key": {
			init: func(c *cloudaccess.Mock, _ *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, ErrFetchingKey)
			},
			withError: ErrFetchingKey,
			testData: testDataForCloudAccess{
				accessKeyUID: 1,
			},
		},
		"error listing access key versions": {
			init: func(c *cloudaccess.Mock, _ *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, ErrListingKeyVersions)
			},
			withError: ErrListingKeyVersions,
			testData: testDataForCloudAccess{
				section:              section,
				accessKeyUID:         1,
				accessKeyName:        "name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"ctr_111"},
						GroupID:     11,
						GroupName:   tools.StringPtr("group11"),
					},
				},
			},
		},
		"error non unique cloud access key id": {
			init: func(c *cloudaccess.Mock, _ *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
			},
			withError: ErrNonUniqueCloudAccessKeyID,
			testData: testDataForCloudAccess{
				section:              section,
				accessKeyUID:         1,
				accessKeyName:        "name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"ctr_111"},
						GroupID:     11,
						GroupName:   tools.StringPtr("group11"),
					},
				},
				accessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
				},
			},
		},
		"error key has no group and contract safeguard": {
			init: func(c *cloudaccess.Mock, _ *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
			},
			withError: ErrNoGroup,
			testData: testDataForCloudAccess{
				accessKeyUID: 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			mc := new(cloudaccess.Mock)
			templateProcessor := new(templates.MockProcessor)
			test.init(mc, templateProcessor, test.dir, test.testData)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCloudAccess(ctx, accessKeyUID, 0, "", section, mc, templateProcessor)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			mc.AssertExpectations(t)
			templateProcessor.AssertExpectations(t)
		})
	}
}

func TestCreateCloudAccess_GroupIDAndConractIDSupplied(t *testing.T) {
	tests := map[string]struct {
		init      func(*cloudaccess.Mock, *templates.MockProcessor, string, testDataForCloudAccess)
		dir       string
		testData  testDataForCloudAccess
		withError error
	}{
		"access key with valid groupID and contractID": {
			init: func(c *cloudaccess.Mock, p *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
				mockProcessTemplates(p, data)
			},
			dir: "groudId_contractId",
			testData: testDataForCloudAccess{
				section:              section,
				flag:                 true,
				accessKeyUID:         1,
				accessKeyName:        "Test name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"C-Contract123"},
						GroupID:     1234,
						GroupName:   tools.StringPtr("group11"),
					},
				},
				accessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
				},
			},
		},
		"error - access key with invalid groupID and contractID combination": {
			init: func(c *cloudaccess.Mock, p *templates.MockProcessor, _ string, data testDataForCloudAccess) {
				mockGetAccessKey(c, data, nil)
				mockListAccessKeyVersions(c, data, nil)
				mockProcessTemplates(p, data)
			},
			withError: errors.New("error populating cloud access data: invalid combination of groupId (1234) and contractId (C-Contract123) for this access key"),
			testData: testDataForCloudAccess{
				section:              section,
				flag:                 true,
				accessKeyUID:         1,
				accessKeyName:        "Test name1",
				authenticationMethod: "AWS4_HMAC_SHA256",
				networkConfiguration: &cloudaccess.SecureNetwork{
					AdditionalCDN:   &chinaCDN,
					SecurityNetwork: "ENHANCED_TLS",
				},
				groups: []cloudaccess.Group{
					{
						ContractIDs: []string{"C-Contract12"},
						GroupID:     1234,
						GroupName:   tools.StringPtr("group11"),
					},
				},
				accessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     1,
						CloudAccessKeyID: tools.StringPtr("key1"),
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			mc := new(cloudaccess.Mock)
			templateProcessor := new(templates.MockProcessor)
			test.init(mc, templateProcessor, test.dir, test.testData)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCloudAccess(ctx, accessKeyUID, 1234, "C-Contract123", section, mc, templateProcessor)
			if test.withError != nil {
				assert.Equal(t, test.withError.Error(), err.Error(), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			mc.AssertExpectations(t)
			templateProcessor.AssertExpectations(t)
		})
	}
}
func TestProcessCloudAccessTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFCloudAccessData
		dir          string
		filesToCheck []string
	}{
		"access key with all data and two versions": {
			givenData: TFCloudAccessData{
				Key: TFCloudAccessKey{
					KeyResourceName:      "TestKeyName",
					AccessKeyName:        "TestKeyName",
					AuthenticationMethod: "AWS4_HMAC_SHA256",
					GroupID:              1234,
					ContractID:           "C-Contract123",
					AccessKeyUID:         1,
					CredentialA:          &Credential{CloudAccessKeyID: "testAccessKey1"},
					CredentialB:          &Credential{CloudAccessKeyID: "testAccessKey2"},
					NetworkConfiguration: &NetworkConfiguration{
						AdditionalCDN:   tools.StringPtr("CHINA_CDN"),
						SecurityNetwork: "ENHANCED_TLS",
					},
				},
				Section: "test_section",
			},
			dir:          "two_versions",
			filesToCheck: []string{"cloudaccess.tf", "import.sh", "variables.tf"},
		},
		"access key with all data and one version": {
			givenData: TFCloudAccessData{
				Key: TFCloudAccessKey{
					KeyResourceName:      "TestKeyName",
					AccessKeyName:        "TestKeyName",
					AuthenticationMethod: "AWS4_HMAC_SHA256",
					GroupID:              1234,
					ContractID:           "C-Contract123",
					AccessKeyUID:         1,
					CredentialA:          &Credential{CloudAccessKeyID: "testAccessKey1"},
					NetworkConfiguration: &NetworkConfiguration{
						AdditionalCDN:   tools.StringPtr("CHINA_CDN"),
						SecurityNetwork: "ENHANCED_TLS",
					},
				},
				Section: "test_section",
			},
			dir:          "one_version",
			filesToCheck: []string{"cloudaccess.tf", "import.sh", "variables.tf"},
		},
		"access key with no versions": {
			givenData: TFCloudAccessData{
				Key: TFCloudAccessKey{
					KeyResourceName:      "TestKeyName",
					AccessKeyName:        "TestKeyName",
					AuthenticationMethod: "AWS4_HMAC_SHA256",
					GroupID:              1234,
					ContractID:           "C-Contract123",
					AccessKeyUID:         1,
					NetworkConfiguration: &NetworkConfiguration{
						AdditionalCDN:   tools.StringPtr("CHINA_CDN"),
						SecurityNetwork: "ENHANCED_TLS",
					},
				},
				Section: "test_section",
			},
			dir:          "no_versions",
			filesToCheck: []string{"cloudaccess.tf", "import.sh", "variables.tf"},
		},
		"access key with no versions and no additional cdn": {
			givenData: TFCloudAccessData{
				Key: TFCloudAccessKey{
					KeyResourceName:      "TestKeyName",
					AccessKeyName:        "TestKeyName",
					AuthenticationMethod: "AWS4_HMAC_SHA256",
					GroupID:              1234,
					ContractID:           "C-Contract123",
					AccessKeyUID:         1,
					NetworkConfiguration: &NetworkConfiguration{
						SecurityNetwork: "ENHANCED_TLS",
					},
				},
				Section: "test_section",
			},
			dir:          "no_additional_cdn",
			filesToCheck: []string{"cloudaccess.tf", "import.sh", "variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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

type testDataForCloudAccess struct {
	section              string
	accessKeyUID         int64
	accessKeyName        string
	authenticationMethod string
	flag                 bool
	networkConfiguration *cloudaccess.SecureNetwork
	groups               []cloudaccess.Group
	accessKeyVersions    []cloudaccess.AccessKeyVersion
}

func mockGetAccessKey(c *cloudaccess.Mock, data testDataForCloudAccess, err error) {
	req := cloudaccess.AccessKeyRequest{
		AccessKeyUID: data.accessKeyUID,
	}
	resp := &cloudaccess.GetAccessKeyResponse{
		AccessKeyUID:         data.accessKeyUID,
		AccessKeyName:        data.accessKeyName,
		AuthenticationMethod: data.authenticationMethod,
		NetworkConfiguration: data.networkConfiguration,
		Groups:               data.groups,
	}
	c.On("GetAccessKey", mock.Anything, req).Return(resp, err).Once()
}

func mockListAccessKeyVersions(c *cloudaccess.Mock, data testDataForCloudAccess, err error) {
	req := cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: data.accessKeyUID,
	}
	resp := &cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: data.accessKeyVersions,
	}
	c.On("ListAccessKeyVersions", mock.Anything, req).Return(resp, err).Once()
}

func mockProcessTemplates(p *templates.MockProcessor, data testDataForCloudAccess) {
	groupID := data.groups[0].GroupID
	contractID := data.groups[0].ContractIDs[0]
	var credA, credB *Credential
	if len(data.accessKeyVersions) == 1 {
		credA = &Credential{CloudAccessKeyID: *data.accessKeyVersions[0].CloudAccessKeyID}
	}
	if len(data.accessKeyVersions) == 2 {
		credA = &Credential{CloudAccessKeyID: *data.accessKeyVersions[0].CloudAccessKeyID}
		credB = &Credential{CloudAccessKeyID: *data.accessKeyVersions[1].CloudAccessKeyID}
	}
	var netConf *NetworkConfiguration
	if data.networkConfiguration != nil {
		netConf = &NetworkConfiguration{
			SecurityNetwork: string(data.networkConfiguration.SecurityNetwork),
		}
		if data.networkConfiguration.AdditionalCDN != nil {
			netConf.AdditionalCDN = tools.StringPtr(string(*data.networkConfiguration.AdditionalCDN))
		}
	}
	p.On("ProcessTemplates", TFCloudAccessData{
		Key: TFCloudAccessKey{
			KeyResourceName:      strings.ReplaceAll(data.accessKeyName, "-", "_"),
			AccessKeyName:        data.accessKeyName,
			AuthenticationMethod: data.authenticationMethod,
			GroupID:              groupID,
			ContractID:           contractID,
			AccessKeyUID:         data.accessKeyUID,
			CredentialA:          credA,
			CredentialB:          credB,
			NetworkConfiguration: netConf,
		},
		Section: data.section,
		Flag:    data.flag,
	}).Return(nil).Once()
}
