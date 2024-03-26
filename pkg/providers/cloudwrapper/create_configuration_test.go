package cloudwrapper

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudwrapper"
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

var (
	section                              = "test_section"
	configID                             = 12345
	getConfigurationReq                  = cloudwrapper.GetConfigurationRequest{ConfigID: int64(configID)}
	getConfigurationActiveResponse       = generateCloudWrapperResponseMock(cloudwrapper.StatusActive, false)
	getConfigurationNoActiveResponse     = generateCloudWrapperResponseMock(cloudwrapper.StatusFailed, false)
	getConfigurationResponseWithMultiCDN = generateCloudWrapperResponseMock(cloudwrapper.StatusActive, true)
	processor                            = func(testdir string) templates.FSTemplateProcessor {
		return templates.FSTemplateProcessor{
			TemplatesFS:     templateFiles,
			AdditionalFuncs: additionalFunctions,
			TemplateTargets: map[string]string{
				"cloudwrapper.tmpl": fmt.Sprintf("./testdata/res/%s/cloudwrapper.tf", testdir),
				"variables.tmpl":    fmt.Sprintf("./testdata/res/%s/variables.tf", testdir),
				"imports.tmpl":      fmt.Sprintf("./testdata/res/%s/import.sh", testdir),
			},
		}
	}
)

func TestCreateCloudWrapper(t *testing.T) {
	tests := map[string]struct {
		init      func(*cloudwrapper.Mock, *templates.MockProcessor, string)
		dir       string
		withError error
	}{
		"configuration all fields": {
			init: func(c *cloudwrapper.Mock, p *templates.MockProcessor, dir string) {
				mockGetConfiguration(c, getConfigurationReq, &getConfigurationActiveResponse, nil)
				mockProcessTemplates(p, (&tfCloudWrapperDataBuilder{}).withDefaults().withStatus(cloudwrapper.StatusActive).build(), nil)
			},
			dir: "all_fields_config",
		},
		"configuration not active status": {
			init: func(c *cloudwrapper.Mock, p *templates.MockProcessor, dir string) {
				mockGetConfiguration(c, getConfigurationReq, &getConfigurationNoActiveResponse, nil)
				mockProcessTemplates(p, (&tfCloudWrapperDataBuilder{}).withDefaults().withStatus(cloudwrapper.StatusFailed).build(), nil)
			},
			dir: "not_active_configuration",
		},
		"error problem with fetching configuration": {
			init: func(c *cloudwrapper.Mock, p *templates.MockProcessor, dir string) {
				mockGetConfiguration(c, getConfigurationReq, nil, ErrFetchingConfiguration)
			},
			dir:       "all_fields_config",
			withError: ErrFetchingConfiguration,
		},
		"error configuration contains multi cdn": {
			init: func(c *cloudwrapper.Mock, p *templates.MockProcessor, dir string) {
				mockGetConfiguration(c, getConfigurationReq, &getConfigurationResponseWithMultiCDN, nil)
			},
			dir:       "all_fields_config",
			withError: ErrContainMultiCDNSettings,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			mock := new(cloudwrapper.Mock)
			templateProcessor := new(templates.MockProcessor)
			test.init(mock, templateProcessor, test.dir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCloudWrapper(ctx, int64(configID), section, mock, templateProcessor)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			mock.AssertExpectations(t)
			templateProcessor.AssertExpectations(t)
		})
	}
}

func TestProcessCloudWrapperTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFCloudWrapperData
		dir          string
		filesToCheck []string
	}{
		"basic configuration with required fields": {
			givenData: TFCloudWrapperData{
				Configuration: TFCWConfiguration{
					ID:                        int64(12345),
					Comments:                  "test",
					PropertyIDs:               []string{"123", "456"},
					ContractID:                "1234",
					ConfigurationResourceName: "test_configuration",
					Name:                      "test_configuration",
					Locations: []Location{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
						{
							Comments:      "TestComments",
							TrafficTypeID: 2,
							Capacity: Capacity{
								Unit:  "TB",
								Value: 2,
							},
						},
					},
				},
				Section: section,
			},
			dir:          "basic_req_fields",
			filesToCheck: []string{"cloudwrapper.tf", "import.sh", "variables.tf"},
		},
		"configuration with all fields": {
			givenData: TFCloudWrapperData{
				Configuration: TFCWConfiguration{
					ID:                        int64(12345),
					Comments:                  "test",
					PropertyIDs:               []string{"123", "456"},
					ContractID:                "1234",
					ConfigurationResourceName: "test_configuration",
					Name:                      "test_configuration",
					NotificationEmails:        []string{"testuser@akamai.com"},
					RetainIdleObjects:         false,
					CapacityAlertsThreshold:   tools.IntPtr(75),
					Status:                    string(cloudwrapper.StatusActive),
					IsActive:                  true,
					Locations: []Location{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
						{
							Comments:      "TestComments",
							TrafficTypeID: 2,
							Capacity: Capacity{
								Unit:  "TB",
								Value: 2,
							},
						},
					},
				},
				Section: section,
			},
			dir:          "all_fields_config",
			filesToCheck: []string{"cloudwrapper.tf", "import.sh", "variables.tf"},
		},
		"configuration multiline comments": {
			givenData: TFCloudWrapperData{
				Configuration: TFCWConfiguration{
					ID:                        int64(12345),
					Comments:                  "first\nsecond\n\nlast",
					PropertyIDs:               []string{"123", "456"},
					ContractID:                "1234",
					ConfigurationResourceName: "test_configuration",
					Name:                      "test_configuration",
					Locations: []Location{
						{
							Comments:      "first\nsecond\n",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
						{
							Comments:      "TestComments",
							TrafficTypeID: 2,
							Capacity: Capacity{
								Unit:  "TB",
								Value: 2,
							},
						},
					},
				},
				Section: section,
			},
			dir:          "multiline_comment",
			filesToCheck: []string{"cloudwrapper.tf", "import.sh", "variables.tf"},
		},
		"not active configuration": {
			givenData: TFCloudWrapperData{
				Configuration: TFCWConfiguration{
					ID:                        int64(12345),
					Comments:                  "test",
					PropertyIDs:               []string{"123", "456"},
					ContractID:                "1234",
					ConfigurationResourceName: "test_configuration",
					Name:                      "test_configuration",
					NotificationEmails:        []string{"testuser@akamai.com"},
					RetainIdleObjects:         false,
					CapacityAlertsThreshold:   tools.IntPtr(75),
					Status:                    string(cloudwrapper.StatusFailed),
					IsActive:                  false,
					Locations: []Location{
						{
							Comments:      "TestComments",
							TrafficTypeID: 1,
							Capacity: Capacity{
								Unit:  "GB",
								Value: 1,
							},
						},
						{
							Comments:      "TestComments",
							TrafficTypeID: 2,
							Capacity: Capacity{
								Unit:  "TB",
								Value: 2,
							},
						},
					},
				},
				Section: section,
			},
			dir:          "not_active_configuration",
			filesToCheck: []string{"cloudwrapper.tf", "import.sh", "variables.tf"},
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

func mockGetConfiguration(c *cloudwrapper.Mock, getConfigurationRequest cloudwrapper.GetConfigurationRequest, getConfigurationResponse *cloudwrapper.Configuration, err error) {
	c.On("GetConfiguration", mock.Anything, getConfigurationRequest).
		Return(getConfigurationResponse, err).Once()
}

func mockProcessTemplates(p *templates.MockProcessor, tfCloudWrapperData TFCloudWrapperData, err error) {
	p.On("ProcessTemplates", tfCloudWrapperData).Return(err).Once()
}

type tfCloudWrapperDataBuilder struct {
	tfCloudWrapperData TFCloudWrapperData
}

func (t *tfCloudWrapperDataBuilder) withDefaults() *tfCloudWrapperDataBuilder {
	t.tfCloudWrapperData = TFCloudWrapperData{
		Configuration: TFCWConfiguration{
			ID:                        int64(12345),
			Comments:                  "test",
			PropertyIDs:               []string{"123", "456"},
			ContractID:                "1234",
			ConfigurationResourceName: "testName",
			Name:                      "testName",
			NotificationEmails:        []string{"testuser@akamai.com"},
			RetainIdleObjects:         false,
			CapacityAlertsThreshold:   tools.IntPtr(75),
			Locations: []Location{
				{
					Comments:      "TestComments",
					TrafficTypeID: 1,
					Capacity: Capacity{
						Unit:  "GB",
						Value: 1,
					},
				},
				{
					Comments:      "TestComments",
					TrafficTypeID: 2,
					Capacity: Capacity{
						Unit:  "TB",
						Value: 2,
					},
				},
			},
		},
		Section: "test_section",
	}
	return t
}

func (t *tfCloudWrapperDataBuilder) withStatus(status cloudwrapper.StatusType) *tfCloudWrapperDataBuilder {
	t.tfCloudWrapperData.Configuration.Status = string(status)
	if status == cloudwrapper.StatusActive {
		t.tfCloudWrapperData.Configuration.IsActive = true
	} else {
		t.tfCloudWrapperData.Configuration.IsActive = false
	}
	return t
}

func (t *tfCloudWrapperDataBuilder) build() TFCloudWrapperData {
	return t.tfCloudWrapperData
}

func generateCloudWrapperResponseMock(status cloudwrapper.StatusType, withMultiCDNSettings bool) cloudwrapper.Configuration {
	configurationResponse := cloudwrapper.Configuration{
		Status:                  status,
		Comments:                "test",
		ContractID:              "1234",
		ConfigID:                12345,
		CapacityAlertsThreshold: tools.IntPtr(75),
		ConfigName:              "testName",
		LastActivatedBy:         tools.StringPtr("user"),
		LastActivatedDate:       tools.StringPtr("2018-03-07T23:40:45Z"),
		LastUpdatedDate:         "2018-03-07T23:40:45Z",
		LastUpdatedBy:           "testUser",
		NotificationEmails:      []string{"testuser@akamai.com"},
		PropertyIDs: []string{
			"123",
			"456",
		},
		RetainIdleObjects: false,
		Locations: []cloudwrapper.ConfigLocationResp{
			{
				Comments:      "TestComments",
				TrafficTypeID: 1,
				Capacity: cloudwrapper.Capacity{
					Unit:  "GB",
					Value: 1,
				},
			},
			{
				Comments:      "TestComments",
				TrafficTypeID: 2,
				Capacity: cloudwrapper.Capacity{
					Unit:  "TB",
					Value: 2,
				},
			},
		},
	}
	if withMultiCDNSettings {
		configurationResponse.MultiCDNSettings =
			&cloudwrapper.MultiCDNSettings{
				CDNs: []cloudwrapper.CDN{
					{
						CDNCode: "testCode",
					},
				},
			}
	}

	return configurationResponse
}
