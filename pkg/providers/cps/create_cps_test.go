package cps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/tools"
	"github.com/akamai/cli-terraform/pkg/templates"
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
	processor = func(testdir string) templates.FSTemplateProcessor {
		return templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"enrollment.tmpl": fmt.Sprintf("./testdata/res/%s/enrollment.tf", testdir),
				"variables.tmpl":  fmt.Sprintf("./testdata/res/%s/variables.tf", testdir),
				"imports.tmpl":    fmt.Sprintf("./testdata/res/%s/import.sh", testdir),
			},
			AdditionalFuncs: template.FuncMap{
				"ToLower": func(val string) string {
					return strings.ToLower(val)
				},
			},
		}
	}

	enrollment = cps.Enrollment{
		AdminContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r1d1@akamai.com",
			FirstName:        "R1",
			LastName:         "D1",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		CertificateChainType: "default",
		CertificateType:      "san",
		ChangeManagement:     false,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "test.akamai.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"test.akamai.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: false,
		NetworkConfiguration: &cps.NetworkConfiguration{
			DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
			DNSNameSettings: &cps.DNSNameSettings{
				CloneDNSNames: false,
				DNSNames:      []string{},
			},
			Geography:        "core",
			MustHaveCiphers:  "ak-akamai-default",
			OCSPStapling:     "on",
			PreferredCiphers: "ak-akamai-default",
			QuicEnabled:      false,
			SecureNetwork:    "enhanced-tls",
			SNIOnly:          true,
		},
		Org: &cps.Org{
			AddressLineOne: "150 Broadway",
			City:           "Cambridge",
			Country:        "US",
			Name:           "Akamai",
			Phone:          "321321321",
			PostalCode:     "12345",
			Region:         "MA",
		},
		RA:                 "lets-encrypt",
		SignatureAlgorithm: "SHA-256",
		TechContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r2d2@akamai.com",
			FirstName:        "R2",
			LastName:         "D2",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		ValidationType: "dv",
	}

	enrollmentMin = cps.Enrollment{
		AdminContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r1d1@akamai.com",
			FirstName:        "R1",
			LastName:         "D1",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		AutoRenewalStartTime: "2022-10-03",
		CertificateType:      "san",
		CSR: &cps.CSR{
			C:  "US",
			CN: "test.akamai.com",
			L:  "Cambridge",
			O:  "Akamai",
			OU: "WebEx",
			ST: "MA",
		},
		Location:                   "loc",
		MaxAllowedSanNames:         10,
		MaxAllowedWildcardSanNames: 20,
		NetworkConfiguration: &cps.NetworkConfiguration{
			DNSNameSettings: &cps.DNSNameSettings{
				DNSNames: []string{"san.test.akamai.com"},
			},
			Geography:     "core",
			SecureNetwork: "enhanced-tls",
			SNIOnly:       true,
		},
		Org: &cps.Org{
			AddressLineOne: "150 Broadway",
			City:           "Cambridge",
			Country:        "US",
			Name:           "Akamai",
			Phone:          "321321321",
			PostalCode:     "12345",
			Region:         "MA",
		},
		PendingChanges:     []string{"change"},
		RA:                 "lets-encrypt",
		SignatureAlgorithm: "SHA-256",
		TechContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r2d2@akamai.com",
			FirstName:        "R2",
			LastName:         "D2",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		ThirdParty:     nil,
		ValidationType: "dv",
	}

	enrollmentAll = cps.Enrollment{
		AdminContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			AddressLineTwo:   "Aka",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r1d1@akamai.com",
			FirstName:        "R1",
			LastName:         "D1",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
			Title:            "title",
		},
		AutoRenewalStartTime: "2022-10-03",
		CertificateChainType: "default",
		CertificateType:      "san",
		ChangeManagement:     true,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "test.akamai.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"san.test.akamai.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: true,
		Location:                       "loc",
		MaxAllowedSanNames:             10,
		MaxAllowedWildcardSanNames:     20,
		NetworkConfiguration: &cps.NetworkConfiguration{
			ClientMutualAuthentication: &cps.ClientMutualAuthentication{
				AuthenticationOptions: &cps.AuthenticationOptions{
					OCSP: &cps.OCSP{
						Enabled: tools.BoolPtr(true),
					},
					SendCAListToClient: tools.BoolPtr(true),
				},
				SetID: "2",
			},
			DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
			DNSNameSettings: &cps.DNSNameSettings{
				CloneDNSNames: true,
				DNSNames:      []string{"san.test.akamai.com"},
			},
			Geography:        "core",
			MustHaveCiphers:  "ak-akamai-default",
			OCSPStapling:     "on",
			PreferredCiphers: "ak-akamai-default",
			QuicEnabled:      true,
			SecureNetwork:    "enhanced-tls",
			SNIOnly:          true,
		},
		Org: &cps.Org{
			AddressLineOne: "150 Broadway",
			AddressLineTwo: "Aka",
			City:           "Cambridge",
			Country:        "US",
			Name:           "Akamai",
			Phone:          "321321321",
			PostalCode:     "12345",
			Region:         "MA",
		},
		PendingChanges:     []string{"change"},
		RA:                 "lets-encrypt",
		SignatureAlgorithm: "SHA-256",
		TechContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r2d2@akamai.com",
			FirstName:        "R2",
			LastName:         "D2",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		ThirdParty:     nil,
		ValidationType: "dv",
	}

	expectGetEnrollment = func(i *mockcps, enrollmentID int, enrollment cps.Enrollment, err error) *mock.Call {
		call := i.On(
			"GetEnrollment",
			mock.Anything,
			cps.GetEnrollmentRequest{
				EnrollmentID: enrollmentID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(&enrollment, nil)
	}
)

func TestCreateCPS(t *testing.T) {
	section := "test_section"
	tests := map[string]struct {
		init                               func(*mockcps)
		enrollmentID                       int
		contractID                         string
		acknowledgePreVerificationWarnings bool
		allowDuplicateCommonName           bool
		filesToCheck                       []string
		dataDir                            string
		jsonDir                            string
		withError                          error
		schema                             bool
	}{
		"export enrollment with minimum fields": {
			init: func(i *mockcps) {
				expectGetEnrollment(i, 1, enrollmentMin, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "enrollment_min",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export enrollment": {
			init: func(i *mockcps) {
				expectGetEnrollment(i, 1, enrollmentAll, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "enrollment_all_fields",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"error fetching enrollment": {
			init: func(i *mockcps) {
				expectGetEnrollment(i, 2, enrollment, fmt.Errorf("oops")).Once()
			},
			enrollmentID: 2,
			withError:    ErrFetchingEnrollment,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("testdata/res/%s/%s", test.dataDir, test.jsonDir), 0755))
			mi := new(mockcps)
			mp := processor(test.dataDir)
			test.init(mi)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCPS(ctx, test.contractID, test.enrollmentID, section, mi, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

			if test.filesToCheck != nil {
				for _, f := range test.filesToCheck {
					expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dataDir, f))
					require.NoError(t, err)
					result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dataDir, f))
					require.NoError(t, err)
					assert.Equal(t, string(expected), string(result))
				}
			}
			mi.AssertExpectations(t)
		})
	}
}

func TestProcessEnrollmentTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFCPSData
		dir          string
		filesToCheck []string
	}{
		"enrollment": {
			givenData: TFCPSData{
				Enrollment:   enrollment,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "enrollment",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
		},
		"enrollment with all fields set": {
			givenData: TFCPSData{
				Enrollment:   enrollmentAll,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "enrollment_all_fields",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
		},
		"enrollment with required only fields": {
			givenData: TFCPSData{
				Enrollment:   enrollmentMin,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "enrollment_min",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
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
