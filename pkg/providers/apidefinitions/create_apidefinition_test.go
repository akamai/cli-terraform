package apidefinitions

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wk8/go-ordered-map/v2"
)

var api = "{\n  \"name\": \"Pet Store\",\n  \"hostnames\": null,\n  \"contractId\": \"\",\n  \"groupId\": 0\n}"

var apiOperations = "{\n  \"operations\": {\n    \"/base\": {\n      \"test login\": {\n        \"method\": \"POST\",\n        \"purpose\": \"SEARCH\"\n      }\n    }\n  }\n}"

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
	processor = func(testdir string, format outputFormat) templates.FSTemplateProcessor {
		processor, err := createTemplateProcessor(fmt.Sprintf("./testdata/res/%s/", testdir), format)
		if err != nil {
			panic(err)
		}
		return *processor
	}
)

func TestCreateAPIDefinition(t *testing.T) {
	tests := map[string]struct {
		format         outputFormat
		init           func(*apidefinitions.Mock, *v0.Mock, *templates.MockProcessor)
		withError      error
		expectedResult TFAPIWrapperData
	}{
		"inactive - openAPI format": {
			format: openAPIFormat,
			init: func(client *apidefinitions.Mock, clientV0 *v0.Mock, p *templates.MockProcessor) {
				mockGetAPI(client, nil, nil)
				mockToOpenAPIFile(clientV0)
				mockGetAPIVersions(client, 1)
				mockProcessTemplates(p, nil)
				mockGetResourceOperation(clientV0)
			},
			expectedResult: inActive,
		},
		"inactive - json format": {
			format: jsonFormat,
			init: func(client *apidefinitions.Mock, clientV0 *v0.Mock, p *templates.MockProcessor) {
				mockGetAPI(client, nil, nil)
				mockGetAPIVersion(clientV0)
				mockGetAPIVersions(client, 1)
				mockProcessTemplates(p, nil)
				mockGetResourceOperation(clientV0)
			},
			expectedResult: inActive,
		},
		"active - openAPI format": {
			format: openAPIFormat,
			init: func(client *apidefinitions.Mock, clientV0 *v0.Mock, p *templates.MockProcessor) {
				mockGetAPI(client, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive))
				mockToOpenAPIFile(clientV0)
				mockGetAPIVersions(client, 2)
				mockProcessTemplates(p, nil)
				mockGetResourceOperation(clientV0)
			},
			expectedResult: active,
		},
		"active - json format": {
			format: jsonFormat,
			init: func(client *apidefinitions.Mock, clientV0 *v0.Mock, p *templates.MockProcessor) {
				mockGetAPI(client, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive))
				mockGetAPIVersion(clientV0)
				mockGetAPIVersions(client, 2)
				mockProcessTemplates(p, nil)
				mockGetResourceOperation(clientV0)
			},
			expectedResult: active,
		},
		"latest_active": {
			format: openAPIFormat,
			init: func(client *apidefinitions.Mock, clientV0 *v0.Mock, p *templates.MockProcessor) {
				mockGetAPI(client, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive))
				mockToOpenAPIFile(clientV0)
				mockGetAPIVersions(client, 1)
				mockProcessTemplates(p, nil)
				mockGetResourceOperation(clientV0)
			},
			expectedResult: latestActive,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mock := new(apidefinitions.Mock)
			mockV0 := new(v0.Mock)
			templateProcessor := new(templates.MockProcessor)
			test.init(mock, mockV0, templateProcessor)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			data, err := createAPIDefinition(ctx, "test_section", test.format, int64(1), nil, mock, mockV0, templateProcessor)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}

			assert.Equal(t, test.expectedResult, *data)
			require.NoError(t, err)
			mock.AssertExpectations(t)
			templateProcessor.AssertExpectations(t)
		})
	}
}

func TestProcessAPIDefinitionTemplates(t *testing.T) {
	tests := map[string]struct {
		format       outputFormat
		givenData    TFAPIWrapperData
		dir          string
		filesToCheck []string
	}{
		"openapi - active on Staging and Production": {
			format:    openAPIFormat,
			givenData: active,
			dir:       "openapi_active",
			filesToCheck: []string{"apidefinitions.tf", "import.sh", "variables.tf",
				"modules/definition/main.tf", "modules/definition/api.yml", "modules/definition/variables.tf",
				"modules/operations/operations.tf",
			},
		},
		"json - active on Staging and Production": {
			format:    jsonFormat,
			givenData: active,
			dir:       "json_active",
			filesToCheck: []string{"apidefinitions.tf", "import.sh", "variables.tf",
				"modules/activation/main.tf", "modules/activation/variables.tf",
				"modules/definition/main.tf", "modules/definition/api.json",
				"modules/operations/operations.tf",
			},
		},
		"inactive on Staging and Production": {
			format:       openAPIFormat,
			givenData:    inActive,
			dir:          "inactive",
			filesToCheck: []string{"apidefinitions.tf", "import.sh"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			require.NoError(t, processor(test.dir, test.format).ProcessTemplates(test.givenData))
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

func TestTrimPrefixAPI(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1234", "1234"},
		{"API_1234", "1234"},
		{"API_", ""},
		{"", ""},
	}

	for _, test := range tests {
		result := trimPrefixAPI(test.input)
		if result != test.expected {
			t.Errorf("TrimPrefixAPI(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
func mockGetAPI(c *apidefinitions.Mock, version *int64, status *apidefinitions.ActivationStatus) {
	c.On("GetEndpoint", mock.Anything, mock.Anything).
		Return(&apidefinitions.GetEndpointResponse{
			Endpoint: apidefinitions.Endpoint{
				APIEndpointName: "Pet Store",
				StagingVersion: apidefinitions.VersionState{
					VersionNumber: version,
					Status:        status,
				},
				ProductionVersion: apidefinitions.VersionState{
					VersionNumber: version,
					Status:        status,
				},
				Locked:     false,
				ContractID: "Contract-1",
				GroupID:    1,
			},
		}, nil).
		Once()
}

func mockGetAPIVersion(c *v0.Mock) {
	c.On("GetAPIVersion", mock.Anything, mock.Anything).
		Return(&v0.GetAPIVersionResponse{
			RegisterAPIRequest: v0.RegisterAPIRequest{
				APIAttributes: v0.APIAttributes{
					Name: "Pet Store",
				},
			},
		}, nil).
		Once()
}

func mockToOpenAPIFile(c *v0.Mock) {
	c.On("ToOpenAPIFile", mock.Anything, mock.Anything).
		Return(ptr.To(v0.ToOpenAPIFileResponse(api)), nil).
		Once()
}

func mockGetAPIVersions(c *apidefinitions.Mock, latestVersion int64) {
	c.On("ListEndpointVersions", mock.Anything, mock.Anything).
		Return(&apidefinitions.ListEndpointVersionsResponse{
			APIVersions: []apidefinitions.APIVersion{{
				VersionNumber: latestVersion,
			}},
		}, nil).
		Once()
}

func mockGetResourceOperation(c *v0.Mock) {
	operations := orderedmap.New[string, *orderedmap.OrderedMap[string, v0.Operation]]()
	baseOperations := orderedmap.New[string, v0.Operation]()
	baseOperations.Set("test login", v0.Operation{
		Method:  ptr.To("POST"),
		Purpose: ptr.To("SEARCH"),
	})
	operations.Set("/base", baseOperations)

	c.On("GetResourceOperation", mock.Anything, mock.Anything).
		Return(&v0.GetResourceOperationResponse{
			ResourceOperations: operations,
		}, nil).Once()
}

func mockProcessTemplates(p *templates.MockProcessor, err error) {
	p.On("ProcessTemplates", mock.Anything).Return(err).Once()
}

var (
	inActive = TFAPIWrapperData{
		API:                  api,
		ID:                   1,
		Version:              1,
		ContractID:           "Contract-1",
		GroupID:              1,
		ResourceName:         "pet_store",
		IsActiveOnStaging:    false,
		IsActiveOnProduction: false,
		StagingVersionKey:    "api_latest_version",
		ProductionVersionKey: "api_latest_version",
		Section:              "test_section",
		Operations:           apiOperations,
	}

	active = TFAPIWrapperData{
		API:                  api,
		ID:                   1,
		Version:              2,
		ContractID:           "Contract-1",
		GroupID:              1,
		ResourceName:         "pet_store",
		IsActiveOnStaging:    true,
		IsActiveOnProduction: true,
		StagingVersionKey:    "api_staging_version",
		ProductionVersionKey: "api_production_version",
		Section:              "test_section",
		Operations:           apiOperations,
	}

	latestActive = TFAPIWrapperData{
		API:                  api,
		ID:                   1,
		Version:              1,
		ContractID:           "Contract-1",
		GroupID:              1,
		ResourceName:         "pet_store",
		IsActiveOnStaging:    true,
		IsActiveOnProduction: true,
		StagingVersionKey:    "api_latest_version",
		ProductionVersionKey: "api_latest_version",
		Section:              "test_section",
		Operations:           apiOperations,
	}
)
