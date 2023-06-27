package dns

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessZone(t *testing.T) {
	tests := map[string]struct {
		filePath       string
		modSegment     bool
		modName        string
		modContentPath string
	}{
		"modSegment=false": {
			filePath:   "./testdata/zone/expected_zone.tf",
			modSegment: false,
		},
		"modSegment=true": {
			filePath:       "./testdata/zone_mod/expected_zone_mod.tf",
			modSegment:     true,
			modName:        "_0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin_com",
			modContentPath: "./testdata/zone_mod/mod/expected_zone_mod_res.tf",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new(fileutilsmock)
			if test.modSegment {
				m.On("createModuleTF", test.modName, mock.Anything, mock.Anything).Return(nil).Once()
			}
			zoneResponse := dns.ZoneResponse{
				Zone:               "0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin.com",
				Type:               "PRIMARY",
				Masters:            []string{},
				ContractID:         "1-3CV382",
				ActivationState:    "NEW",
				LastModifiedBy:     "jreed",
				LastActivationDate: "2021-03-16T17:16:59.208264Z",
				VersionId:          "fd858f59-6014-4ce4-8372-c08389d809e8",
				TsigKey:            &dns.TSIGKey{Name: "some-name", Algorithm: "some-algorithm", Secret: "some-secret"},
			}
			zone, err := processZone(context.Background(), &zoneResponse, "_0007770b-08a8-4b5f-a46b-081b772ba605-sbodden-calvin_com", test.modSegment, m, "./")
			require.NoError(t, err)
			m.AssertExpectations(t)

			if test.modSegment {
				assertFileWithContent(t, test.modContentPath, m.createModuleArg)
			}
			assertFileWithContent(t, test.filePath, zone)
		})
	}
}
