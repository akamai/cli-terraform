package gtm

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
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
	domainPath = "./testdata/res/testdata_domain.tf"
	if _, err := os.Create(domainPath); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := os.RemoveAll("./testdata/res"); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

var (
	domain = &gtm.Domain{
		Name:                    "1test.name.akadns.net",
		Type:                    "test",
		ModificationComments:    "cli-terraform test domain",
		EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
		DefaultTimeoutPenalty:   10,
		LoadImbalancePercentage: 50,
		DefaultErrorPenalty:     90,
		CnameCoalescingEnabled:  true,
		LoadFeedback:            true,
		EndUserMappingEnabled:   true,
		Datacenters: []*gtm.Datacenter{
			{
				Nickname:     "TEST1",
				DatacenterId: 123,
			},
			{
				Nickname:     "TEST2",
				DatacenterId: 124,
			},
			{
				Nickname:     "DEFAULT",
				DatacenterId: 5400,
			},
		},
		Resources: []*gtm.Resource{
			{
				Name: "test resource1",
			},
			{
				Name: "test resource2",
			},
		},
		Properties: []*gtm.Property{
			{
				Name: "test property1",
			},
			{
				Name: "test property2",
			},
		},
		AsMaps: []*gtm.AsMap{
			{
				Name: "test_asmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterId: 5004,
				},
			},
		},
		GeographicMaps: []*gtm.GeoMap{
			{
				Name: "test_geomap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterId: 5004,
				},
			},
		},
		CidrMaps: []*gtm.CidrMap{
			{
				Name: "test_cidrmap",
				DefaultDatacenter: &gtm.DatacenterBase{
					Nickname:     "default",
					DatacenterId: 5004,
				},
			},
		},
	}

	domainData = TFDomainData{
		Section:                 "test_section",
		Name:                    "1test.name.akadns.net",
		NormalizedName:          "_1test_name",
		Type:                    "test",
		Comment:                 "cli-terraform test domain",
		EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
		DefaultTimeoutPenalty:   10,
		LoadImbalancePercentage: 50,
		DefaultErrorPenalty:     90,
		CnameCoalescingEnabled:  true,
		LoadFeedback:            true,
		EndUserMappingEnabled:   true,
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
		DatacentersImportList: map[int]string{
			123: "TEST1",
			124: "TEST2",
		},
		Resources: map[string][]int{
			"test resource1": {},
			"test resource2": {},
		},
		Properties: map[string][]int{
			"test property1": {},
			"test property2": {},
		},
		Asmaps: map[string][]int{
			"test_asmap": {},
		},
		Cidrmaps: map[string][]int{
			"test_cidrmap": {},
		},
		Geomaps: map[string][]int{
			"test_geomap": {},
		},
	}

	expectGTMProcessTemplates = func(mp *mockProcessor, data TFDomainData, err error) *mock.Call {
		call := mp.On("ProcessTemplates", data)
		if err != nil {
			return call.Return(err)
		}
		return call.Return(nil)
	}

	expectGetDomain = func(mg *mockGTM, domainName string, domain *gtm.Domain, err error) *mock.Call {
		call := mg.On("GetDomain", mock.Anything, domainName)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(domain, nil)
	}

	expectNullFieldMap = func(mg *mockGTM, domain *gtm.Domain) *mock.Call {
		return mg.On("NullFieldMap", mock.Anything, domain).Return(&gtm.NullFieldMapStruct{}, nil)
	}
)

func TestCreateDomain(t *testing.T) {
	section := "test_section"
	domainName := "test.name.net"

	tests := map[string]struct {
		init      func(*mockGTM, *mockProcessor)
		withError error
	}{
		"fetch domain success": {
			init: func(mg *mockGTM, mp *mockProcessor) {
				expectGetDomain(mg, domainName, domain, nil).Once()
				expectNullFieldMap(mg, domain).Once()
				expectGTMProcessTemplates(mp, domainData, nil).Once()
			},
		},
		"error fetching domain": {
			init: func(mg *mockGTM, mp *mockProcessor) {
				expectGetDomain(mg, domainName, domain, fmt.Errorf("oops")).Once()
			},
			withError: ErrFetchingDomain,
		},
		"error processing template": {
			init: func(mg *mockGTM, mp *mockProcessor) {
				expectGetDomain(mg, domainName, domain, nil).Once()
				expectGTMProcessTemplates(mp, domainData, templates.ErrSavingFiles).Once()
			},
			withError: templates.ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mgtm := new(mockGTM)
			mp := new(mockProcessor)
			test.init(mgtm, mp)

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createDomain(ctx, mgtm, domainName, section, mp)
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
	tests := map[string]struct {
		givenData    interface{}
		dir          string
		filesToCheck []string
	}{
		"variables file correct": {
			givenData:    TFDomainData{Section: "test_section"},
			dir:          "only_variables",
			filesToCheck: []string{"variables.tf"},
		},
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
				Resources: map[string][]int{
					"test resource1": {},
					"test resource2": {},
				},
				Properties: map[string][]int{
					"test property1": {},
					"test property2": {},
				},
				Asmaps: map[string][]int{
					"test_asmap": {},
				},
				Cidrmaps: map[string][]int{
					"test_cidrmap": {},
				},
				Geomaps: map[string][]int{
					"test_geomap": {},
				},
			},
			dir:          "import_script",
			filesToCheck: []string{"import.sh"},
		},
		"domain without other resources": {
			givenData: TFDomainData{
				Section:                 "default",
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CnameCoalescingEnabled:  true,
				LoadFeedback:            true,
				EndUserMappingEnabled:   false,
			},
			dir:          "domain_file",
			filesToCheck: []string{"domain.tf", "variables.tf", "import.sh"},
		},
		"simple domain with datacenters": {
			givenData: TFDomainData{
				Section:                 "test_section",
				Name:                    "test.name.akadns.net",
				NormalizedName:          "test_name",
				Type:                    "basic",
				Comment:                 "test",
				EmailNotificationList:   []string{"john@akamai.com", "jdoe@akamai.com"},
				DefaultTimeoutPenalty:   10,
				LoadImbalancePercentage: 50,
				DefaultErrorPenalty:     90,
				CnameCoalescingEnabled:  true,
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
			},
			dir:          "with_datacenters",
			filesToCheck: []string{"domain.tf", "datacenters.tf", "variables.tf", "import.sh"},
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
					"variables.tmpl":   filepath.Join(outDir, "variables.tf"),
				},
				AdditionalFuncs: template.FuncMap{
					"normalize": normalizeResourceName,
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))
			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(filepath.Join("./testdata", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(filepath.Join(outDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
