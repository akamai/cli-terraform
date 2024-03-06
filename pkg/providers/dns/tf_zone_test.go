package dns

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessZone(t *testing.T) {
	defaultZoneResponse := dns.ZoneResponse{
		Zone:               "0007770b-08a8-4b5f-a46b-081b772ba605-test.com",
		Type:               "PRIMARY",
		Masters:            []string{},
		ContractID:         "test_contract",
		ActivationState:    "NEW",
		LastModifiedBy:     "jreed",
		LastActivationDate: "2021-03-16T17:16:59.208264Z",
		VersionID:          "fd858f59-6014-4ce4-8372-c08389d809e8",
		TSIGKey:            &dns.TSIGKey{Name: "some-name", Algorithm: "some-algorithm", Secret: "some-secret"},
	}
	tests := map[string]struct {
		filePath       string
		modSegment     bool
		modName        string
		modContentPath string
		zoneResponse   dns.ZoneResponse
	}{
		"modSegment=false": {
			filePath:     "./testdata/zone/expected_zone.tf",
			modSegment:   false,
			zoneResponse: defaultZoneResponse,
		},
		"modSegment=true": {
			filePath:       "./testdata/zone_mod/expected_zone_mod.tf",
			modSegment:     true,
			modName:        "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			modContentPath: "./testdata/zone_mod/mod/expected_zone_mod_res.tf",
			zoneResponse:   defaultZoneResponse,
		},
		"modSegment=false, multiline comment": {
			filePath:   "./testdata/zone/expected_zone_multiline.tf",
			modSegment: false,
			zoneResponse: dns.ZoneResponse{
				Zone:               "0007770b-08a8-4b5f-a46b-081b772ba605-test.com",
				Type:               "PRIMARY",
				Masters:            []string{},
				ContractID:         "test_contract",
				ActivationState:    "NEW",
				LastModifiedBy:     "jreed",
				LastActivationDate: "2021-03-16T17:16:59.208264Z",
				VersionID:          "fd858f59-6014-4ce4-8372-c08389d809e8",
				Comment:            "first\nsecond\n\nlast",
				TSIGKey:            &dns.TSIGKey{Name: "some-name", Algorithm: "some-algorithm", Secret: "some-secret"},
			},
		},
		"modSegment=false, multiline comment ending newline": {
			filePath:   "./testdata/zone/expected_zone_multiline2.tf",
			modSegment: false,
			zoneResponse: dns.ZoneResponse{
				Zone:               "0007770b-08a8-4b5f-a46b-081b772ba605-test.com",
				Type:               "PRIMARY",
				Masters:            []string{},
				ContractID:         "test_contract",
				ActivationState:    "NEW",
				LastModifiedBy:     "jreed",
				LastActivationDate: "2021-03-16T17:16:59.208264Z",
				VersionID:          "fd858f59-6014-4ce4-8372-c08389d809e8",
				Comment:            "first\nsecond\n",
				TSIGKey:            &dns.TSIGKey{Name: "some-name", Algorithm: "some-algorithm", Secret: "some-secret"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := new(fileUtilsMock)
			if test.modSegment {
				m.On("createModuleTF", test.modName, mock.Anything, mock.Anything).Return(nil).Once()
			}
			zone, err := processZone(context.Background(), &test.zoneResponse, "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com", test.modSegment, m, "./")
			require.NoError(t, err)
			m.AssertExpectations(t)

			if test.modSegment {
				assertFileWithContent(t, test.modContentPath, m.createModuleArg)
			}
			assertFileWithContent(t, test.filePath, zone)
		})
	}
}
