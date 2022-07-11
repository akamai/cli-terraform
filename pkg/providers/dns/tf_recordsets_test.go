package dns

import (
	"context"
	"testing"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
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

func TestProcessRecordset(t *testing.T) {
	tests := map[string]struct {
		expectRootPath string
		expectModPath  string
	}{
		"basic case": {
			expectRootPath: "./testdata/recordset_mod/expected_recordsets_mod_resource.tf",
			expectModPath:  "./testdata/recordset_mod/expected_recordsets_mod_variables.tf",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new(mockdns)

			ctx := context.Background()
			zone := "0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin.com"
			metadata := dns.MetadataH{}
			recordset := dns.Recordset{
				Name:  "someName",
				Type:  "someType",
				TTL:   1000,
				Rdata: []string{"INTEL-386", "Unix"},
			}
			recordsets := make([]dns.Recordset, 0)
			recordsets = append(recordsets, recordset)
			response := dns.RecordSetResponse{Metadata: metadata, Recordsets: recordsets}
			m.On("GetRecordsets", ctx, zone, mock.Anything).Return(&response, nil).Once()
			parsedRData := map[string]interface{}{"hardware": "INTEL-386", "software": "Unix"}
			m.On("ParseRData", ctx, recordset.Type, recordset.Rdata).Return(parsedRData).Once()

			fus := new(fileutilsmock)
			fus.On("appendRootModuleTF", mock.Anything).Return(nil).Once()
			fus.On("createModuleTF", "zoneName_someName_someType", mock.Anything).Return(nil).Once()
			zoneTypeMap := make(map[string]map[string]bool)
			zoneTypeMap["someName"] = map[string]bool{"someType": true}
			processingResult, _ := processRecordsets(ctx, m, zone,
				"zoneName", zoneTypeMap, fetchConfigStruct{}, fus)

			assert.Equal(t, 1, len(processingResult))
			types, nameExist := processingResult[recordset.Name]
			assert.True(t, nameExist)
			assert.Equal(t, 1, len(types))
			assert.Equal(t, recordset.Type, types[0])

			assertFileWithContent(t, test.expectRootPath, fus.appendRootArg)
			assertFileWithContent(t, test.expectModPath, fus.createModuleArg)

		})
	}
}
