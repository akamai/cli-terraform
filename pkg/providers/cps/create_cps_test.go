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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/jinzhu/copier"
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

	certECDSAForTests       = "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE ECDSA REQUEST-----"
	certRSAForTests         = "-----BEGIN CERTIFICATE RSA REQUEST-----\n...\n-----END CERTIFICATE RSA REQUEST-----"
	trustChainRSAForTests   = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST-----"
	trustChainECDSAForTests = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----"
	RSA                     = "RSA"
	ECDSA                   = "ECDSA"

	enrollmentDV = cps.GetEnrollmentResponse{
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
			C:                   "US",
			CN:                  "test.akamai.com",
			L:                   "Cambridge",
			O:                   "Akamai",
			OU:                  "WebEx",
			PreferredTrustChain: "intermediate-a",
			SANS:                []string{"test.akamai.com"},
			ST:                  "MA",
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

	enrollmentDVMin = cps.GetEnrollmentResponse{
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
		PendingChanges: []cps.PendingChange{
			{
				Location:   "change",
				ChangeType: "new-certificate",
			},
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
		ThirdParty:     nil,
		ValidationType: "dv",
	}

	enrollmentDVAll = cps.GetEnrollmentResponse{
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
			C:                   "US",
			CN:                  "test.akamai.com",
			L:                   "Cambridge",
			O:                   "Akamai",
			OU:                  "WebEx",
			PreferredTrustChain: "intermediate-a",
			SANS:                []string{"test.akamai.com", "san.test.akamai.com"},
			ST:                  "MA",
		},
		EnableMultiStackedCertificates: true,
		Location:                       "loc",
		MaxAllowedSanNames:             10,
		MaxAllowedWildcardSanNames:     20,
		NetworkConfiguration: &cps.NetworkConfiguration{
			ClientMutualAuthentication: &cps.ClientMutualAuthentication{
				AuthenticationOptions: &cps.AuthenticationOptions{
					OCSP: &cps.OCSP{
						Enabled: ptr.To(true),
					},
					SendCAListToClient: ptr.To(true),
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
		PendingChanges: []cps.PendingChange{
			{
				Location:   "change",
				ChangeType: "new-certificate",
			},
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
		ThirdParty:     nil,
		ValidationType: "dv",
	}

	enrollmentThirdParty = func(sans []string) cps.GetEnrollmentResponse {
		return cps.GetEnrollmentResponse{
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
				SANS: sans,
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
							Enabled: ptr.To(true),
						},
						SendCAListToClient: ptr.To(true),
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
			PendingChanges: []cps.PendingChange{
				{
					Location:   "change",
					ChangeType: "new-certificate",
				},
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
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: true,
			},
			ValidationType: "third-party",
		}
	}

	enrollmentThirdPartyAll = enrollmentThirdParty([]string{"test.akamai.com", "san.test.akamai.com"})

	enrollmentThirdPartyAllNoSANs = enrollmentThirdParty([]string{"test.akamai.com"})

	enrollmentOV = cps.GetEnrollmentResponse{
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
		RA:                 "symantec",
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
		ValidationType: "ov",
	}

	expectGetEnrollment = func(m *cps.Mock, enrollmentID int, enrollment cps.GetEnrollmentResponse, err error) *mock.Call {
		call := m.On(
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

	expectGetChangeHistory = func(m *cps.Mock, enrollmentID int, response cps.GetChangeHistoryResponse, err error) *mock.Call {
		call := m.On(
			"GetChangeHistory",
			mock.Anything,
			cps.GetChangeHistoryRequest{
				EnrollmentID: enrollmentID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(&response, nil)
	}
)

func TestCreateCPS(t *testing.T) {
	section := "test_section"
	tests := map[string]struct {
		init         func(*cps.Mock)
		enrollmentID int
		contractID   string
		filesToCheck []string
		dataDir      string
		jsonDir      string
		withError    error
	}{
		"export DV enrollment with minimum fields": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentDVMin, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "dv_enrollment_min",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export DV enrollment without DNSNameSettings": {
			init: func(m *cps.Mock) {
				var enrollment cps.GetEnrollmentResponse
				_ = copier.CopyWithOption(&enrollment, enrollmentDVMin, copier.Option{DeepCopy: true})
				enrollment.NetworkConfiguration.DNSNameSettings = nil
				expectGetEnrollment(m, 1, enrollment, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "dv_enrollment_min",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export DV enrollment": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentDVAll, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "dv_enrollment_all_fields",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment ecdsa": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAll, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							PrimaryCertificate: cps.CertificateChangeHistory{
								Certificate:  certECDSAForTests,
								KeyAlgorithm: ECDSA,
								TrustChain:   trustChainECDSAForTests,
							},
							Status: "inactive",
						},
						{
							PrimaryCertificate: cps.CertificateChangeHistory{
								Certificate:  certRSAForTests,
								KeyAlgorithm: RSA,
								TrustChain:   trustChainRSAForTests,
							},
							Status: "active",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_all_fields_ecdsa",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment rsa": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAll, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							PrimaryCertificate: cps.CertificateChangeHistory{
								Certificate:  certRSAForTests,
								KeyAlgorithm: RSA,
								TrustChain:   trustChainRSAForTests,
							},
							Status: "active",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_all_fields_rsa",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment ecdsa+rsa": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAll, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							PrimaryCertificate: cps.CertificateChangeHistory{
								Certificate:  certECDSAForTests,
								KeyAlgorithm: ECDSA,
								TrustChain:   trustChainECDSAForTests,
							},
							MultiStackedCertificates: []cps.CertificateChangeHistory{
								{
									Certificate:  certRSAForTests,
									KeyAlgorithm: RSA,
									TrustChain:   trustChainRSAForTests,
								},
							},
							Status: "active",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_all_fields_ecdsa_rsa",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment renewal": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAll, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							Action:            "renew",
							ActionDescription: "Renew Certificate",
							Status:            "incomplete",
							RA:                "third-party",
						},
						{
							PrimaryCertificate: cps.CertificateChangeHistory{
								Certificate:  certECDSAForTests,
								KeyAlgorithm: ECDSA,
								TrustChain:   trustChainECDSAForTests,
							},
							MultiStackedCertificates: []cps.CertificateChangeHistory{
								{
									Certificate:  certRSAForTests,
									KeyAlgorithm: RSA,
									TrustChain:   trustChainRSAForTests,
								},
							},
							Status: "active",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_all_fields_renewal",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment new certificate": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAll, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							Action:            "new-certificate",
							ActionDescription: "Create New Certificate",
							Status:            "incomplete",
							RA:                "third-party",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_all_fields_new_certificate",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"export third party enrollment new certificate no SANs": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 1, enrollmentThirdPartyAllNoSANs, nil).Once()
				response := cps.GetChangeHistoryResponse{
					Changes: []cps.ChangeHistory{
						{
							Action:            "new-certificate",
							ActionDescription: "Create New Certificate",
							Status:            "incomplete",
							RA:                "third-party",
						},
					},
				}
				expectGetChangeHistory(m, 1, response, nil).Once()
			},
			enrollmentID: 1,
			contractID:   "ctr_1",
			dataDir:      "third_party_enrollment_new_certificate_no_sans",
			filesToCheck: []string{"enrollment.tf", "import.sh", "variables.tf"},
		},
		"error fetching enrollment": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 2, enrollmentDV, fmt.Errorf("oops")).Once()
			},
			enrollmentID: 2,
			withError:    ErrFetchingEnrollment,
		},
		"provided ov enrollment": {
			init: func(m *cps.Mock) {
				expectGetEnrollment(m, 3, enrollmentOV, nil).Once()
			},
			enrollmentID: 3,
			withError:    ErrUnsupportedEnrollmentType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("testdata/res/%s/%s", test.dataDir, test.jsonDir), 0755))
			mi := new(cps.Mock)
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
		"dv enrollment": {
			givenData: TFCPSData{
				Enrollment:   enrollmentDV,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "dv_enrollment",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
		},
		"dv enrollment with all fields set": {
			givenData: TFCPSData{
				Enrollment:   enrollmentDVAll,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "dv_enrollment_all_fields",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
		},
		"dv enrollment with required only fields": {
			givenData: TFCPSData{
				Enrollment:   enrollmentDVMin,
				EnrollmentID: 1,
				ContractID:   "ctr_1",
				Section:      "test_section",
			},
			dir:          "dv_enrollment_min",
			filesToCheck: []string{"enrollment.tf", "variables.tf", "import.sh"},
		},
		"third party enrollment with all fields set": {
			givenData: TFCPSData{
				Enrollment:       enrollmentThirdPartyAll,
				EnrollmentID:     1,
				ContractID:       "ctr_1",
				Section:          "test_section",
				CertificateECDSA: "-----BEGIN CERTIFICATE ECDSA REQUEST-----\\n...\\n-----END CERTIFICATE ECDSA REQUEST-----",
				CertificateRSA:   "-----BEGIN CERTIFICATE RSA REQUEST-----\\n...\\n-----END CERTIFICATE RSA REQUEST-----",
				TrustChainECDSA:  "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\\n...\\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----",
				TrustChainRSA:    "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\\n...\\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST-----",
			},
			dir:          "third_party_enrollment_all_fields_ecdsa_rsa",
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
