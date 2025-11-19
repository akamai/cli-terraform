package gtm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
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
	domain = &gtm.GetDomainResponse{
		Name:                    "1test.name.akadns.net",
		Type:                    "test",
		ModificationComments:    "cli-terraform test domain",
		EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
		DefaultTimeoutPenalty:   10,
		LoadImbalancePercentage: 50,
		DefaultErrorPenalty:     90,
		CNameCoalescingEnabled:  true,
		LoadFeedback:            true,
		EndUserMappingEnabled:   true,
		SignAndServe:            true,
		SignAndServeAlgorithm:   tools.StringPtr("RSA-SHA1"),
		Datacenters: []gtm.Datacenter{
			{
				Nickname:     "TEST1",
				DatacenterID: 123,
			},
			{
				Nickname:     "TEST2",
				DatacenterID: 124,
			},
			{
				Nickname:     "DEFAULT",
				DatacenterID: 5400,
			},
		},
		Resources: []gtm.Resource{
			{
				Name: "test resource1",
			},
			{
				Name: "test resource2",
			},
		},
		Properties: []gtm.Property{
			{
				Name:                 "test property1",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
					},
				},
			},
			{
				Name:                 "test property2",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				StaticRRSets: []gtm.StaticRRSet{
					{
						Type:  "test type",
						Rdata: []string{"rdata1", "rdata2"},
					},
				},
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
					{
						DatacenterID: 124,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"7.6.5.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
						HTTPHeaders: []gtm.HTTPHeader{
							{
								Name:  "header1",
								Value: "header1Value",
							},
							{
								Name:  "header2",
								Value: "header2Value",
							},
						},
					},
				},
			},
		},
		ASMaps: []gtm.ASMap{
			{
				Name: "test_asmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		GeographicMaps: []gtm.GeoMap{
			{
				Name: "test_geomap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		CIDRMaps: []gtm.CIDRMap{
			{
				Name: "test_cidrmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
	}

	domainData = TFDomainData{
		EdgercPath:              "~/.edgerc",
		Section:                 "test_section",
		Name:                    "1test.name.akadns.net",
		NormalizedName:          "_1test_name",
		Type:                    "test",
		Comment:                 "cli-terraform test domain",
		EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
		DefaultTimeoutPenalty:   10,
		LoadImbalancePercentage: 50,
		DefaultErrorPenalty:     90,
		CNameCoalescingEnabled:  true,
		LoadFeedback:            true,
		EndUserMappingEnabled:   true,
		SignAndServe:            true,
		SignAndServeAlgorithm:   "RSA-SHA1",
		DefaultDatacenters: []TFDatacenterData{
			{
				Nickname: "DEFAULT",
				ID:       5400,
			},
		},
		Datacenters: []TFDatacenterData{
			{
				Nickname: "TEST1",
				ID:       123,
			},
			{
				Nickname: "TEST2",
				ID:       124,
			},
		},
		ASMaps: []gtm.ASMap{
			{
				Name: "test_asmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		GeoMaps: []gtm.GeoMap{
			{
				Name: "test_geomap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		CIDRMaps: []gtm.CIDRMap{
			{
				Name: "test_cidrmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		Resources: []gtm.Resource{
			{
				Name: "test resource1",
			},
			{
				Name: "test resource2",
			},
		},
		Properties: []gtm.Property{
			{
				Name:                 "test property1",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
					},
				},
			},
			{
				Name:                 "test property2",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				StaticRRSets: []gtm.StaticRRSet{
					{
						Type:  "test type",
						Rdata: []string{"rdata1", "rdata2"},
					},
				},
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
					{
						DatacenterID: 124,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"7.6.5.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
						HTTPHeaders: []gtm.HTTPHeader{
							{
								Name:  "header1",
								Value: "header1Value",
							},
							{
								Name:  "header2",
								Value: "header2Value",
							},
						},
					},
				},
			},
		},
	}

	domainDataWithNonDefaultEdgercPathAndSection = TFDomainData{
		EdgercPath:              "/non/default/path/to/edgerc",
		Section:                 "non_default_section",
		Name:                    "1test.name.akadns.net",
		NormalizedName:          "_1test_name",
		Type:                    "test",
		Comment:                 "cli-terraform test domain",
		EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
		DefaultTimeoutPenalty:   10,
		LoadImbalancePercentage: 50,
		DefaultErrorPenalty:     90,
		CNameCoalescingEnabled:  true,
		LoadFeedback:            true,
		EndUserMappingEnabled:   true,
		SignAndServe:            true,
		SignAndServeAlgorithm:   "RSA-SHA1",
		DefaultDatacenters: []TFDatacenterData{
			{
				Nickname: "DEFAULT",
				ID:       5400,
			},
		},
		Datacenters: []TFDatacenterData{
			{
				Nickname: "TEST1",
				ID:       123,
			},
			{
				Nickname: "TEST2",
				ID:       124,
			},
		},
		ASMaps: []gtm.ASMap{
			{
				Name: "test_asmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		GeoMaps: []gtm.GeoMap{
			{
				Name: "test_geomap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		CIDRMaps: []gtm.CIDRMap{
			{
				Name: "test_cidrmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterID: 5004,
				},
			},
		},
		Resources: []gtm.Resource{
			{
				Name: "test resource1",
			},
			{
				Name: "test resource2",
			},
		},
		Properties: []gtm.Property{
			{
				Name:                 "test property1",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
					},
				},
			},
			{
				Name:                 "test property2",
				Type:                 "performance",
				ScoreAggregationType: "worst",
				DynamicTTL:           60,
				HandoutLimit:         8,
				HandoutMode:          "normal",
				StaticRRSets: []gtm.StaticRRSet{
					{
						Type:  "test type",
						Rdata: []string{"rdata1", "rdata2"},
					},
				},
				TrafficTargets: []gtm.TrafficTarget{
					{
						DatacenterID: 123,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"1.2.3.4"},
					},
					{
						DatacenterID: 124,
						Enabled:      true,
						Weight:       1,
						Servers:      []string{"7.6.5.4"},
					},
				},
				LivenessTests: []gtm.LivenessTest{
					{
						Name:               "HTTP",
						TestInterval:       60,
						TestObject:         "/",
						HTTPError3xx:       true,
						HTTPError4xx:       true,
						HTTPError5xx:       true,
						TestObjectProtocol: "HTTP",
						TestObjectPort:     80,
						TestTimeout:        10,
						HTTPHeaders: []gtm.HTTPHeader{
							{
								Name:  "header1",
								Value: "header1Value",
							},
							{
								Name:  "header2",
								Value: "header2Value",
							},
						},
					},
				},
			},
		},
	}

	expectGTMProcessTemplates = func(mp *templates.MockProcessor, data TFDomainData, err error) *mock.Call {
		call := mp.On("ProcessTemplates", data)
		if err != nil {
			return call.Return(err)
		}
		return call.Return(nil)
	}

	expectGetDomain = func(mg *gtm.Mock, _ gtm.GetDomainRequest, err error) *mock.Call {
		call := mg.On("GetDomain", mock.Anything, mock.AnythingOfType("gtm.GetDomainRequest"))
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(domain, nil)
	}
)

func TestCreateDomain(t *testing.T) {
	domainName := "test.name.net"

	tests := map[string]struct {
		edgercPath string
		section    string
		init       func(*gtm.Mock, *templates.MockProcessor)
		withError  error
	}{
		"fetch domain success": {
			init: func(mg *gtm.Mock, mp *templates.MockProcessor) {
				expectGetDomain(mg, gtm.GetDomainRequest{DomainName: domainName}, nil).Once()
				expectGTMProcessTemplates(mp, domainData, nil).Once()
			},
		},
		"error fetching domain": {
			init: func(mg *gtm.Mock, _ *templates.MockProcessor) {
				expectGetDomain(mg, gtm.GetDomainRequest{DomainName: domain.Name}, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingDomain,
		},
		"error processing template": {
			init: func(mg *gtm.Mock, mp *templates.MockProcessor) {
				expectGetDomain(mg, gtm.GetDomainRequest{DomainName: domain.Name}, nil).Once()
				expectGTMProcessTemplates(mp, domainData, templates.ErrSavingFiles).Once()
			},
			withError: templates.ErrSavingFiles,
		},
		"non default edgerc path and section": {
			edgercPath: "/non/default/path/to/edgerc",
			section:    "non_default_section",
			init: func(mg *gtm.Mock, mp *templates.MockProcessor) {
				expectGetDomain(mg, gtm.GetDomainRequest{DomainName: domainName}, nil).Once()
				expectGTMProcessTemplates(mp, domainDataWithNonDefaultEdgercPathAndSection, nil).Once()
			},
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
			mgtm := new(gtm.Mock)
			mp := new(templates.MockProcessor)
			test.init(mgtm, mp)

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createDomain(ctx, mgtm, domainName, test.edgercPath, test.section, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

			mgtm.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessDomainTemplates(t *testing.T) {
	defaultEdgercPath := "~/.edgerc"
	defaultSection := "test_section"

	tests := map[string]struct {
		givenData    interface{}
		dir          string
		filesToCheck []string
	}{
		"import script correct": {
			givenData: TFDomainData{
				Name:           "test.name.akadns.net",
				NormalizedName: "test_name",
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname: "TEST1",
						ID:       123,
					},
					{
						Nickname: "TEST2",
						ID:       124,
					},
					{
						Nickname: "TEST3",
						ID:       125,
					},
				},
				Resources: []gtm.Resource{
					{
						Name: "test resource1",
					},
					{
						Name: "test resource2",
					},
				},
				Properties: []gtm.Property{
					{
						Name: "test property1",
					},
					{
						Name: "test property2",
					},
				},
				ASMaps: []gtm.ASMap{
					{
						Name: "test_asmap",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 123,
						},
					},
				},
				GeoMaps: []gtm.GeoMap{
					{
						Name: "test_geomap",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 124,
						},
					},
				},
				CIDRMaps: []gtm.CIDRMap{
					{
						Name: "test_cidrmap",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 125,
						},
					},
				},
			},
			dir:          "import_script",
			filesToCheck: []string{"import.sh"},
		},
		"domain without other resources": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				EndUserMappingEnabled:   false,
			},
			dir:          "domain_file",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh"},
		},
		"domain with sign_and_serve_algorithm": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				EndUserMappingEnabled:   false,
				SignAndServe:            true,
				SignAndServeAlgorithm:   "RSA-SHA1",
			},
			dir:          "domain_with_sign_and_serve",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh"},
		},
		"simple domain with datacenters": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "qtr",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 5400,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
					},
				},
			},
			dir:          "with_datacenters",
			filesToCheck: []string{"domain.tf", "datacenters.tf", "properties.tf", "variables.tf", "import.sh"},
		},
		"simple domain with maps": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				ASMaps: []gtm.ASMap{
					{
						Name: "test_asmap",
						Assignments: []gtm.ASAssignment{
							{
								DatacenterBase: gtm.DatacenterBase{
									Nickname:     "TEST1",
									DatacenterID: 123,
								},
								ASNumbers: []int64{1, 2, 3},
							},
						},
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 123,
						},
					},
				},
				GeoMaps: []gtm.GeoMap{
					{
						Name: "test_geomap",
						Assignments: []gtm.GeoAssignment{
							{
								DatacenterBase: gtm.DatacenterBase{
									Nickname:     "TEST1",
									DatacenterID: 123,
								},
								Countries: []string{"US"},
							},
						},
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 124,
						},
					},
				},
				CIDRMaps: []gtm.CIDRMap{
					{
						Name: "test_cidrmap",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "default",
							DatacenterID: 124,
						},
					},
				},
			},
			dir:          "with_maps",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh", "maps.tf"},
		},
		"simple domain with resources": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Resources: []gtm.Resource{
					{
						Type:                "XML load object via HTTP",
						HostHeader:          "header",
						LeastSquaresDecay:   30,
						Description:         "some description",
						LeaderString:        "leader",
						ConstrainedProperty: "**",
						ResourceInstances: []gtm.ResourceInstance{
							{
								DatacenterID:         123,
								UseDefaultLoadObject: false,
								LoadObject: gtm.LoadObject{
									LoadObject:     "load",
									LoadObjectPort: 80,
									LoadServers:    []string{"server"},
								},
							},
						},
						AggregationType:             "latest",
						LoadImbalancePercentage:     51,
						UpperBound:                  20,
						Name:                        "test resource1",
						MaxUMultiplicativeIncrement: 10,
						DecayRate:                   5,
					},
					{
						Name: "test resource2",
					},
				},
			},
			dir:          "with_resources",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh", "resources.tf"},
		},
		"simple domain with properties": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "static",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						Comments:             "some comment",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:                    "HTTP",
								TestInterval:            60,
								TestObject:              "/",
								HTTPError3xx:            true,
								HTTPError4xx:            true,
								HTTPError5xx:            true,
								TestObjectProtocol:      "HTTP",
								TestObjectPort:          80,
								TestTimeout:             10,
								HTTPMethod:              tools.StringPtr("GET"),
								HTTPRequestBody:         tools.StringPtr("Body"),
								AlternateCACertificates: []string{"test1"},
								Pre2023SecurityPosture:  true,
							},
						},
					},
					{
						Name:                 "test property2",
						Type:                 "performance",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						StaticRRSets: []gtm.StaticRRSet{
							{
								Type:  "test type",
								Rdata: []string{"rdata1", "rdata2", "\"properlyescaped\""},
							},
						},
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
							{
								DatacenterID: 124,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"7.6.5.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
								HTTPHeaders: []gtm.HTTPHeader{
									{
										Name:  "header1",
										Value: "header1Value",
									},
									{
										Name:  "header2",
										Value: "header2Value",
									},
								},
							},
						},
					},
					{
						Name:                 "test property3",
						Type:                 "asmapping",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 5400,
								Enabled:      true,
								Weight:       0,
								Servers:      []string{},
							},
							{
								DatacenterID: 124,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{},
							},
						},
					},
				},
			},
			dir:          "with_properties",
			filesToCheck: []string{"domain.tf", "datacenters.tf", "properties.tf", "variables.tf", "import.sh"},
		},
		"simple domain with ranked_failover properties": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "ranked-failover",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						Comments:             "some comment",
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
					},
					{
						Name:                 "test property2",
						Type:                 "ranked-failover",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						StaticRRSets: []gtm.StaticRRSet{
							{
								Type:  "test type",
								Rdata: []string{"rdata1", "rdata2", "\"properlyescaped\""},
							},
						},
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
								Precedence:   tools.IntPtr(10),
							},
							{
								DatacenterID: 124,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"7.6.5.4"},
								Precedence:   tools.IntPtr(200),
							},
							{
								DatacenterID: 5400,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"7.6.5.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
								HTTPHeaders: []gtm.HTTPHeader{
									{
										Name:  "header1",
										Value: "header1Value",
									},
									{
										Name:  "header2",
										Value: "header2Value",
									},
								},
							},
						},
					},
					{
						Name:                 "test property3",
						Type:                 "ranked-failover",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 5400,
								Enabled:      true,
								Weight:       0,
								Servers:      []string{},
								Precedence:   tools.IntPtr(100),
							},
							{
								DatacenterID: 124,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{},
							},
						},
					},
				},
			},
			dir:          "with_ranked_failover_properties",
			filesToCheck: []string{"domain.tf", "datacenters.tf", "properties.tf", "variables.tf", "import.sh"},
		},
		"simple domain with property of type 'qtr'": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
				},
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT_5401",
						ID:       5401,
					},
					{
						Nickname: "DEFAULT_5402",
						ID:       5402,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "qtr",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 5401,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
					},
					{
						Name:                 "test property2",
						Type:                 "qtr",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
							{
								DatacenterID: 5402,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"7.6.5.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
					},
				},
			},
			dir:          "with_qtr_properties",
			filesToCheck: []string{"domain.tf", "datacenters.tf", "properties.tf", "variables.tf", "import.sh"},
		},
		"simple domain with resources and properties with multilines": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "first\nsecond\n\nlast",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Resources: []gtm.Resource{
					{
						Type:                "XML load object via HTTP",
						HostHeader:          "header",
						LeastSquaresDecay:   30,
						Description:         "first\nsecond\n\nlast",
						LeaderString:        "leader",
						ConstrainedProperty: "**",
						ResourceInstances: []gtm.ResourceInstance{
							{
								DatacenterID:         123,
								UseDefaultLoadObject: false,
								LoadObject: gtm.LoadObject{
									LoadObject:     "load",
									LoadObjectPort: 80,
									LoadServers:    []string{"server"},
								},
							},
						},
						AggregationType:             "latest",
						LoadImbalancePercentage:     51,
						UpperBound:                  20,
						Name:                        "test resource1",
						MaxUMultiplicativeIncrement: 10,
						DecayRate:                   5,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "static",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
						Comments: "first\nsecond\n\nlast",
					},
				},
			},
			dir:          "with_multiline",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh", "resources.tf", "properties.tf"},
		},
		"simple domain with resources and properties with multilines - empty line at the end": {
			givenData: TFDomainData{
				EdgercPath:              defaultEdgercPath,
				Section:                 defaultSection,
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "first\nsecond\n",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				DefaultDatacenters: []TFDatacenterData{
					{
						Nickname: "DEFAULT",
						ID:       5400,
					},
				},
				Datacenters: []TFDatacenterData{
					{
						Nickname:        "TEST1",
						ID:              123,
						City:            "New York",
						StateOrProvince: "NY",
						Country:         "US",
						Latitude:        40.71305,
						Longitude:       -74.00723,
						DefaultLoadObject: &gtm.LoadObject{
							LoadObject:     "test load object",
							LoadObjectPort: 111,
							LoadServers:    []string{"loadServer1", "loadServer2", "loadServer3"},
						},
					},
					{
						Nickname:        "TEST2",
						ID:              124,
						City:            "Chicago",
						StateOrProvince: "IL",
						Country:         "US",
						Latitude:        41.88323,
						Longitude:       -87.6324,
					},
				},
				Resources: []gtm.Resource{
					{
						Type:                "XML load object via HTTP",
						HostHeader:          "header",
						LeastSquaresDecay:   30,
						Description:         "first\nsecond\n",
						LeaderString:        "leader",
						ConstrainedProperty: "**",
						ResourceInstances: []gtm.ResourceInstance{
							{
								DatacenterID:         123,
								UseDefaultLoadObject: false,
								LoadObject: gtm.LoadObject{
									LoadObject:     "load",
									LoadObjectPort: 80,
									LoadServers:    []string{"server"},
								},
							},
						},
						AggregationType:             "latest",
						LoadImbalancePercentage:     51,
						UpperBound:                  20,
						Name:                        "test resource1",
						MaxUMultiplicativeIncrement: 10,
						DecayRate:                   5,
					},
				},
				Properties: []gtm.Property{
					{
						Name:                 "test property1",
						Type:                 "static",
						ScoreAggregationType: "worst",
						DynamicTTL:           60,
						HandoutLimit:         8,
						HandoutMode:          "normal",
						TrafficTargets: []gtm.TrafficTarget{
							{
								DatacenterID: 123,
								Enabled:      true,
								Weight:       1,
								Servers:      []string{"1.2.3.4"},
							},
						},
						LivenessTests: []gtm.LivenessTest{
							{
								Name:               "HTTP",
								TestInterval:       60,
								TestObject:         "/",
								HTTPError3xx:       true,
								HTTPError4xx:       true,
								HTTPError5xx:       true,
								TestObjectProtocol: "HTTP",
								TestObjectPort:     80,
								TestTimeout:        10,
							},
						},
						Comments: "first\nsecond\n",
					},
				},
			},
			dir:          "with_multiline2",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh", "resources.tf", "properties.tf"},
		},
		"non default edgerc path and section": {
			givenData: TFDomainData{
				EdgercPath:              "/non/default/path/to/edgerc",
				Section:                 "non_default_section",
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CNameCoalescingEnabled:  true,
				LoadFeedback:            true,
				EndUserMappingEnabled:   false,
			},
			dir:          "non_default_edgerc_path_and_section",
			filesToCheck: []string{"variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			outDir := filepath.Join("./testdata/res", test.dir)
			require.NoError(t, os.MkdirAll(outDir, 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"datacenters.tmpl": filepath.Join(outDir, "datacenters.tf"),
					"domain.tmpl":      filepath.Join(outDir, "domain.tf"),
					"imports.tmpl":     filepath.Join(outDir, "import.sh"),
					"maps.tmpl":        filepath.Join(outDir, "maps.tf"),
					"resources.tmpl":   filepath.Join(outDir, "resources.tf"),
					"properties.tmpl":  filepath.Join(outDir, "properties.tf"),
					"variables.tmpl":   filepath.Join(outDir, "variables.tf"),
				},
				AdditionalFuncs: additionalFunctions,
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))
			for _, f := range test.filesToCheck {
				expected, err := os.ReadFile(filepath.Join("./testdata", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(filepath.Join(outDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
