package dns

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-terraform/v2/pkg/tools/tests"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_createNamedModulePath(t *testing.T) {
	tests := map[string]struct {
		tfWorkPath   string
		modName      string
		expectedPath string
	}{
		"tfWorkPath = ./": {
			tfWorkPath:   "./",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = test_path": {
			tfWorkPath:   "test_path",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "test_path/" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = ../": {
			tfWorkPath:   "../",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "../" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"blank tfWorkPath": {
			tfWorkPath:   "",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, test.expectedPath, createNamedModulePath(test.modName, test.tfWorkPath), "createNamedModulePath(%v, %v)", test.modName, test.tfWorkPath)
		})
	}
}

func Test_createDNSVarsConfig(t *testing.T) {
	tests := map[string]struct {
		edgercPath    string
		edgercSection string
		contractID    string
		expectedFile  string
	}{
		"default edgerc path and section": {
			edgercPath:    "~/.edgerc",
			edgercSection: "default",
			contractID:    "ctr_1-23456",
			expectedFile:  "testdata/dnsvars/dnsvars_default.tf",
		},
		"non default edgerc path and section": {
			edgercPath:    "/non/default/path/to/edgerc",
			edgercSection: "non_default_section",
			contractID:    "ctr_A-B-C",
			expectedFile:  "testdata/dnsvars/dnsvars_non_default.tf",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			contractID = test.contractID

			f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
			}()

			ctx := terminal.Context(context.Background(), terminal.New(f, nil, f))
			term := terminal.Get(ctx)
			err = createDNSVarsConfig(term, tmpDir, test.edgercPath, test.edgercSection)
			require.NoError(t, err)

			dnsVarsContent, err := os.ReadFile(filepath.Join(tmpDir, "dnsvars.tf"))
			require.NoError(t, err)

			expectedContent, err := os.ReadFile(test.expectedFile)
			require.NoError(t, err)
			assert.Equal(t, string(expectedContent), string(dnsVarsContent))
		})
	}
}

func TestCreateZone(t *testing.T) {
	tests := map[string]struct {
		init        func(*dns.Mock)
		setup       func(*testing.T, string)
		config      configStruct
		expectedErr error
		dir         string
	}{
		"successful export with --resources flag": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
				mockGetZoneNameTypes(m, "example.com", "www.example.com", []string{"A"})
				mockGetZoneNameTypes(m, "example.com", "example.com", []string{"NS", "SOA"})
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
			},
			dir: "create_zone_resources_flag",
		},
		"successful export with --createconfig flag without --resources flag": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			setup: func(t *testing.T, resDir string) {
				content := `{"Zone":"example.com","RecordSets":{"example.com":["NS","SOA"]}}`
				require.NoError(t, os.WriteFile(filepath.Join(resDir, "example_com_resources.json"), []byte(content), 0644))
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
			},
			dir: "create_zone_createconfig_from_file",
		},
		"successful export with --importscript flag without --createconfig flag": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
			},
			setup: func(t *testing.T, resDir string) {
				content := `{"example.com":["NS","SOA"]}`
				require.NoError(t, os.WriteFile(filepath.Join(resDir, "example_com_zoneconfig.json"), []byte(content), 0644))
			},
			config: configStruct{
				zoneName:     "example.com",
				importScript: true,
			},
			dir: "create_zone_importscript_from_file",
		},
		"successful export with --resources, --namesonly and --createconfig flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				fetchConfig:            fetchConfigStruct{NamesOnly: true},
			},
			dir: "create_zone_namesonly_and_config",
		},
		"successful export with --resources and --namesonly flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				fetchConfig:            fetchConfigStruct{NamesOnly: true},
			},
			dir: "create_zone_resources_namesonly_flag",
		},
		"successful export with --resources and --createconfig flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
				mockGetZoneNameTypes(m, "example.com", "www.example.com", []string{"A"})
				mockGetZoneNameTypes(m, "example.com", "example.com", []string{"NS", "SOA"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
			},
			dir: "create_zone_resources_and_config",
		},
		"successful export with --recordname and --configonly flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "abc.example.com", Type: "TXT", TTL: 300, Rdata: []string{"\"dummy text abc\""}},
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
				})
				mockParseRData(m, "TXT", []string{"\"dummy text abc\""}, map[string]interface{}{"target": []string{"\"dummy text abc\""}})
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
				fetchConfig:  fetchConfigStruct{ConfigOnly: true},
				recordNames:  []string{"abc.example.com"},
			},
			dir: "create_zone_recordname_and_configonly",
		},
		"successful export with --createconfig and --configonly flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
				fetchConfig:  fetchConfigStruct{ConfigOnly: true},
			},
			dir: "create_zone_configonly",
		},
		"successful export with --resources, --createconfig and --importscript flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
				mockGetZoneNameTypes(m, "example.com", "www.example.com", []string{"A"})
				mockGetZoneNameTypes(m, "example.com", "example.com", []string{"NS", "SOA"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				importScript:           true,
			},
			dir: "create_zone_resources_config_and_importscript",
		},
		"successful export with --resources, --createconfig, --importscript and --segmentconfig flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"www.example.com", "example.com"})
				mockGetZoneNameTypes(m, "example.com", "www.example.com", []string{"A"})
				mockGetZoneNameTypes(m, "example.com", "example.com", []string{"NS", "SOA"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				importScript:           true,
				fetchConfig:            fetchConfigStruct{ModSegment: true},
			},
			dir: "create_zone_segmentconfig",
		},
		"successful export with --recordname and --segmentconfig flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNameTypes(m, "example.com", "abc.example.com", []string{"TXT"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "abc.example.com", Type: "TXT", TTL: 300, Rdata: []string{"\"dummy text abc\""}},
				})
				mockParseRData(m, "TXT", []string{"\"dummy text abc\""}, map[string]interface{}{"target": []string{"\"dummy text abc\""}})
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				recordNames:            []string{"abc.example.com"},
				fetchConfig:            fetchConfigStruct{ModSegment: true},
			},
			dir: "create_zone_recordname_and_segmentconfig",
		},
		"successful export with single --recordname flag": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNameTypes(m, "example.com", "abc.example.com", []string{"TXT"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "abc.example.com", Type: "TXT", TTL: 300, Rdata: []string{"\"dummy text abc\""}},
				})
				mockParseRData(m, "TXT", []string{"\"dummy text abc\""}, map[string]interface{}{"target": []string{"\"dummy text abc\""}})
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				recordNames:            []string{"abc.example.com"},
			},
			dir: "create_zone_single_recordname",
		},
		"successful export with multiple --recordname flags": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNameTypes(m, "example.com", "abc.example.com", []string{"TXT"})
				mockGetZoneNameTypes(m, "example.com", "def.example.com", []string{"TXT"})
				mockGetRecordSets(m, "example.com", []dns.RecordSet{
					{Name: "abc.example.com", Type: "TXT", TTL: 300, Rdata: []string{"\"dummy text abc\""}},
					{Name: "def.example.com", Type: "TXT", TTL: 300, Rdata: []string{"\"dummy text def\""}},
				})
				mockParseRData(m, "TXT", []string{"\"dummy text abc\""}, map[string]interface{}{"target": []string{"\"dummy text abc\""}})
				mockParseRData(m, "TXT", []string{"\"dummy text def\""}, map[string]interface{}{"target": []string{"\"dummy text def\""}})
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
				createConfig:           true,
				recordNames:            []string{"abc.example.com", "def.example.com"},
			},
			dir: "create_zone_multiple_recordnames",
		},
		"successful export with --createconfig flag and paginated GetRecordSets": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetRecordSetsPage(m, "example.com", 1, 2, []dns.RecordSet{
					{Name: "example.com", Type: "NS", TTL: 3600, Rdata: []string{"ns1.example.com."}},
				})
				mockGetRecordSetsPage(m, "example.com", 2, 2, []dns.RecordSet{
					{Name: "example.com", Type: "SOA", TTL: 3600, Rdata: []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}},
				})
				mockParseRDataNS(m)
				mockParseRDataSOA(m)
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
				fetchConfig:  fetchConfigStruct{ConfigOnly: true},
			},
			dir: "create_zone_paginated_recordsets",
		},
		"zone type ALIAS is not supported": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "ALIAS", "test_contract")
			},
			config:      configStruct{zoneName: "example.com"},
			expectedErr: ErrAliasZoneNotSupported,
		},
		"GetZone API error": {
			init: func(m *dns.Mock) {
				m.On("GetZone", mock.Anything, dns.GetZoneRequest{Zone: "example.com"}).
					Return(nil, errors.New("zone not found")).Once()
			},
			config:      configStruct{zoneName: "example.com"},
			expectedErr: ErrZoneRetrievalFailed,
		},
		"resources file already exists on disk": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"example.com"})
				mockGetZoneNameTypes(m, "example.com", "example.com", []string{"NS", "SOA"})
			},
			setup: func(t *testing.T, resDir string) {
				require.NoError(t, os.WriteFile(filepath.Join(resDir, "example_com_resources.json"), []byte("{}"), 0644))
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
			},
			dir:         "create_zone_resources_file_exists",
			expectedErr: ErrResourceListFileExists,
		},
		"GetZoneNames API error": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				m.On("GetZoneNames", mock.Anything, dns.GetZoneNamesRequest{Zone: "example.com"}).
					Return(nil, errors.New("API failure")).Once()
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
			},
			expectedErr: ErrZoneNamesRetrievalFailed,
		},
		"GetZoneNameTypes API error": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				mockGetZoneNames(m, "example.com", []string{"example.com"})
				m.On("GetZoneNameTypes", mock.Anything, dns.GetZoneNameTypesRequest{Zone: "example.com", ZoneName: "example.com"}).
					Return(nil, errors.New("API failure")).Once()
			},
			config: configStruct{
				zoneName:               "example.com",
				shouldCreateImportList: true,
			},
			expectedErr: ErrZoneNameTypesRetrievalFailed,
		},
		"GetRecordSets API error": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				m.On("GetRecordSets", mock.Anything, matchGetRecordSetsRequest("example.com")).
					Return(nil, errors.New("API failure")).Once()
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
				fetchConfig:  fetchConfigStruct{ConfigOnly: true},
			},
			dir:         "create_zone_recordsets_api_error",
			expectedErr: ErrRecordSetsRetrievalFailed,
		},
		"GetRecordSets API error on subsequent page": {
			init: func(m *dns.Mock) {
				mockGetZone(m, "example.com", "PRIMARY", "test_contract")
				m.On("GetRecordSets", mock.Anything, matchGetRecordSetsRequestPage("example.com", 1)).
					Return(&dns.GetRecordSetsResponse{
						Metadata:   dns.Metadata{Page: 1, LastPage: 2},
						RecordSets: []dns.RecordSet{},
					}, nil).Once()
				m.On("GetRecordSets", mock.Anything, matchGetRecordSetsRequestPage("example.com", 2)).
					Return(nil, errors.New("API failure")).Once()
			},
			config: configStruct{
				zoneName:     "example.com",
				createConfig: true,
				fetchConfig:  fetchConfigStruct{ConfigOnly: true},
			},
			dir:         "create_zone_recordsets_api_error_subsequent_page",
			expectedErr: ErrRecordSetsRetrievalFailed,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new(dns.Mock)
			test.init(m)

			cfg := test.config

			if test.dir != "" {
				resDir := filepath.Join("testdata", "res", test.dir)
				require.NoError(t, os.RemoveAll(resDir))
				require.NoError(t, os.MkdirAll(resDir, 0755))
				cfg.tfWorkPath = resDir
				if test.setup != nil {
					test.setup(t, resDir)
				}
			}

			// set package-level zoneName used inside inventorZone
			zoneName = cfg.zoneName

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))

			err := createZone(ctx, "~/.edgerc", "default", m, cfg)

			if test.expectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.expectedErr))
			} else {
				require.NoError(t, err)
				resDir := filepath.Join("testdata", "res", test.dir)
				expectedDir := filepath.Join("testdata", "end_to_end", test.dir)
				assertOutputDir(t, expectedDir, resDir)
			}

			m.AssertExpectations(t)
		})
	}
}

// assertOutputDir compares every file in expectedDir against the file with the same name
// in actualDir, recursing into subdirectories.
func assertOutputDir(t *testing.T, expectedDir, actualDir string) {
	files, err := os.ReadDir(expectedDir)
	require.NoError(t, err, "failed to read expected output dir %s", expectedDir)
	for _, file := range files {
		if file.IsDir() {
			assertOutputDir(t,
				filepath.Join(expectedDir, file.Name()),
				filepath.Join(actualDir, file.Name()),
			)
			continue
		}
		expectedContent, err := os.ReadFile(filepath.Join(expectedDir, file.Name()))
		require.NoError(t, err)
		actualContent, err := os.ReadFile(filepath.Join(actualDir, file.Name()))
		require.NoError(t, err, "error opening generated file %s", file.Name())

		normalize := selectContentNormalizer(file.Name())
		assert.Equal(t,
			normalize(string(expectedContent)),
			normalize(string(actualContent)),
			"file %s content mismatch", file.Name(),
		)
	}
}

// mockGetZone sets up a GetZone mock expectation returning the given zone type and contractID.
func mockGetZone(m *dns.Mock, zone, zoneType, contractID string) {
	m.On("GetZone", mock.Anything, dns.GetZoneRequest{Zone: zone}).
		Return(&dns.GetZoneResponse{
			Zone:       zone,
			Type:       zoneType,
			ContractID: contractID,
		}, nil).Once()
}

// mockGetZoneNameTypes sets up a GetZoneNameTypes mock expectation returning the provided types for the given zone and record name.
func mockGetZoneNameTypes(m *dns.Mock, zone, zoneName string, types []string) {
	m.On("GetZoneNameTypes", mock.Anything, dns.GetZoneNameTypesRequest{Zone: zone, ZoneName: zoneName}).
		Return(&dns.GetZoneNameTypesResponse{Types: types}, nil).Once()
}

// mockGetZoneNames sets up a GetZoneNames mock expectation returning the provided names for the given zone.
func mockGetZoneNames(m *dns.Mock, zone string, names []string) {
	m.On("GetZoneNames", mock.Anything, dns.GetZoneNamesRequest{Zone: zone}).
		Return(&dns.GetZoneNamesResponse{Names: names}, nil).Once()
}

// mockGetRecordSets sets up a GetRecordSets mock expectation returning the provided record sets for the given zone.
func mockGetRecordSets(m *dns.Mock, zone string, recordSets []dns.RecordSet) {
	mockGetRecordSetsPage(m, zone, 1, 1, recordSets)
}

// mockGetRecordSetsPage sets up a GetRecordSets mock expectation for a specific page returning the provided record sets.
func mockGetRecordSetsPage(m *dns.Mock, zone string, page, lastPage int, recordSets []dns.RecordSet) {
	m.On("GetRecordSets", mock.Anything, matchGetRecordSetsRequestPage(zone, page)).
		Return(&dns.GetRecordSetsResponse{
			Metadata:   dns.Metadata{Page: page, LastPage: lastPage},
			RecordSets: recordSets,
		}, nil).Once()
}

// matchGetRecordSetsRequest returns a matcher for GetRecordSetsRequest with the given zone.
// Using exact request is complicated due to getQueryArguments dependent on the environment. To be refactored later.
func matchGetRecordSetsRequest(zone string) any {
	return mock.MatchedBy(func(req dns.GetRecordSetsRequest) bool {
		return req.Zone == zone
	})
}

// matchGetRecordSetsRequestPage returns a matcher for GetRecordSetsRequest with the given zone and page.
// Using exact request is complicated due to getQueryArguments dependent on the environment. To be refactored later.
func matchGetRecordSetsRequestPage(zone string, page int) any {
	return mock.MatchedBy(func(req dns.GetRecordSetsRequest) bool {
		return req.Zone == zone && req.QueryArgs.Page == page
	})
}

// mockParseRData sets up a ParseRData mock expectation returning the provided result for the given type and rdata.
func mockParseRData(m *dns.Mock, recordType string, rdata []string, result map[string]interface{}) {
	m.On("ParseRData", mock.Anything, recordType, rdata).Return(result).Once()
}

// mockParseRDataNS sets up a ParseRData mock expectation for the standard NS record.
func mockParseRDataNS(m *dns.Mock) {
	mockParseRData(m, "NS", []string{"ns1.example.com."}, map[string]interface{}{"target": []string{"ns1.example.com."}})
}

// mockParseRDataSOA sets up a ParseRData mock expectation for the standard SOA record.
func mockParseRDataSOA(m *dns.Mock) {
	mockParseRData(m, "SOA", []string{"ns1.example.com. admin.example.com. 2024010101 3600 600 604800 300"}, map[string]interface{}{
		"originserver": "ns1.example.com.",
		"contact":      "admin.example.com.",
		"serial":       2024010101,
		"refresh":      3600,
		"retry":        600,
		"expiry":       604800,
		"minimum":      300,
	})
}

func selectContentNormalizer(filename string) func(string) string {
	switch filepath.Ext(filename) {
	case ".tf":
		return func(content string) string {
			// hclwrite.Format can be dropped once the export-zone start using FSTemplateProcessor
			contentBytes := hclwrite.Format([]byte(content))
			content = tests.NormalizeFieldsInBlock(string(contentBytes), `resource "akamai_dns_record"`) //ToDo: tej częsci nie jestem pewny czy faktycznie tutaj musi być normalizacja dwa razy. To jest chyba tylko dlatego jakby w jednym pliku modułowym był resource i moduł (mimo że nie dotyczą tego samego)
			return tests.NormalizeBlocksInFile(content, `module "`)
		}
	case ".script":
		return tests.NormalizeWholeFile
	default:
		return func(content string) string { return content }
	}
}
