package dns

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessStringNoQuotes(t *testing.T) {

	sourceString := "no quotes"
	expectedString := "no quotes"

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessStringEscapedQuotes(t *testing.T) {

	sourceString := "test \"four\" test"
	expectedString := `test \"four\" test`

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessStringEmbeddedQuotes(t *testing.T) {

	sourceString := `first string" "secondString`
	expectedString := `first string\" \"secondString`

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessStringBackslash(t *testing.T) {

	sourceString := `hello\world`
	expectedString := `hello\\world`

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessRecordset(t *testing.T) {
	tests := map[string]struct {
		mod            bool
		expectRootPath string
		expectModPath  string
	}{
		"modSegment=false": {
			mod:            false,
			expectRootPath: "./testdata/recordset/expected_recordsets_resource.tf",
		},
		"modSegment=true": {
			mod:            true,
			expectRootPath: "./testdata/recordset_mod/expected_recordsets_mod_resource.tf",
			expectModPath:  "./testdata/recordset_mod/expected_recordsets_mod_variables.tf",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new(dns.Mock)

			ctx := context.Background()
			zone := "0007770b-08a8-4b5f-a46b-081b772ba605-test.com"
			metadata := dns.Metadata{}
			recordset := dns.RecordSet{
				Name:  "someName",
				Type:  "someType",
				TTL:   1000,
				Rdata: []string{"INTEL-386", "Unix"},
			}
			recordSets := make([]dns.RecordSet, 0)
			recordSets = append(recordSets, recordset)
			response := dns.RecordSetResponse{Metadata: metadata, RecordSets: recordSets}
			m.On("GetRecordSets", ctx, zone, mock.Anything).Return(&response, nil).Once()
			parsedRData := map[string]interface{}{"hardware": "INTEL-386", "software": "Unix"}
			m.On("ParseRData", ctx, recordset.Type, recordset.Rdata).Return(parsedRData).Once()

			fus := new(fileUtilsMock)
			fus.On("appendRootModuleTF", mock.Anything).Return(nil).Once()
			if test.mod {
				fus.On("createModuleTF", "zoneName_someName_someType", mock.Anything, mock.Anything).Return(nil).Once()
			}
			zoneTypeMap := make(map[string]map[string]bool)
			zoneTypeMap["someName"] = map[string]bool{"someType": true}
			config := configStruct{fetchConfig: fetchConfigStruct{ModSegment: test.mod}}
			processingResult, _ := processRecordSets(ctx, m, zone, "zoneName", zoneTypeMap, fus, config)

			assert.Equal(t, 1, len(processingResult))
			types, nameExist := processingResult[recordset.Name]
			assert.True(t, nameExist)
			assert.Equal(t, 1, len(types))
			assert.Equal(t, recordset.Type, types[0])

			assertFileWithContent(t, test.expectRootPath, fus.appendRootArg)
			if test.mod {
				assertFileWithContent(t, test.expectModPath, fus.createModuleArg)
			}
		})
	}
}
