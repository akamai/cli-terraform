package cloudcertificates

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type (
	certificateMockData struct {
		id                   string
		name                 string
		contractID           string
		baseName             string
		keyType              string
		keySize              string
		secureNetwork        string
		sans                 []string
		certificateStatus    string
		subject              *cloudcertificates.Subject
		signedCertificatePEM *string
		trustChainPEM        *string
	}
)

var (
	signedCertificatePEM = `-----BEGIN CERTIFICATE-----
testsignedcertificate
-----END CERTIFICATE-----`

	trustChainPEM = `-----BEGIN CERTIFICATE-----
testtrustchaincertificate1
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
testtrustchaincertificate2
-----END CERTIFICATE-----`

	signedCertificatePEMWithNewline       = signedCertificatePEM + "\n"
	trustChainPEMWithNewline              = trustChainPEM + "\n"
	signedCertificatePEMWithDoubleNewline = signedCertificatePEMWithNewline + "\n"
	trustChainPEMWithDoubleNewline        = trustChainPEMWithNewline + "\n"

	emptySubject = cloudcertificates.Subject{}
)

func TestProcessCloudCertificateTemplates(t *testing.T) {
	tests := map[string]struct {
		dir                 string
		init                func(*certificateMockData, *cloudcertificates.Mock)
		mockData            certificateMockData
		edgercPath          string
		configSection       string
		withError           string
		withTemplatingError bool
	}{
		"certificate status is CSR_READY": {
			dir:      "csr_ready",
			mockData: getCertificateMockDataNoPEMs(),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is READY_FOR_USE": {
			dir: "ready_for_use",
			mockData: getCertificateMockData("READY_FOR_USE", "test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is READY_FOR_USE no trust_chain": {
			dir:      "ready_for_use_no_trust_chain",
			mockData: getCertificateMockDataNoTrustChain(),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is ACTIVE": {
			dir: "active",
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is ACTIVE - with paging": {
			dir: "active",
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				firstCertBatch := []cloudcertificates.Certificate{}
				for i := range 100 {
					firstCertBatch = append(firstCertBatch, cloudcertificates.Certificate{
						CertificateID: fmt.Sprintf("cert-id-%d", i),
					})
				}
				// Mock that the first page returns 100 certificates and a next link.
				m.On("ListCertificates", mock.Anything, cloudcertificates.ListCertificatesRequest{
					PageSize: maxPageSize,
					Page:     1,
				}).Return(&cloudcertificates.ListCertificatesResponse{
					Certificates: firstCertBatch,
					Links: cloudcertificates.Links{
						Next: ptr.To("next-page-link"),
					},
				}, nil).Once()
				// Mock that the second page returns the target certificate.
				m.On("ListCertificates", mock.Anything, cloudcertificates.ListCertificatesRequest{
					PageSize: maxPageSize,
					Page:     2,
				}).Return(&cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{
						{
							CertificateID:     d.id,
							ContractID:        d.contractID,
							CertificateName:   d.name,
							KeyType:           cloudcertificates.CryptographicAlgorithm(d.keyType),
							KeySize:           cloudcertificates.KeySize(d.keySize),
							SecureNetwork:     d.secureNetwork,
							SANs:              d.sans,
							CertificateStatus: d.certificateStatus,
							Subject:           d.subject,
						},
					},
				}, nil).Once()

				d.mockGetCertificate(m)
			},
		},
		"certificate status is ACTIVE, PEMs end in newline - no trimsuffix needed": {
			dir: "active_pems_end_in_newline",
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEMWithNewline, &trustChainPEMWithNewline),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is ACTIVE, PEMs end in double newline - no trimsuffix needed": {
			dir: "active_pems_end_in_double_newline",
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEMWithDoubleNewline, &trustChainPEMWithDoubleNewline),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate status is ACTIVE, custom edgerc and config section": {
			dir:           "active_custom_edgerc_section",
			edgercPath:    "custom_edgerc_path",
			configSection: "custom-section",
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate name begins with a number": {
			dir: "name_begins_with_number",
			mockData: getCertificateMockData("ACTIVE", "123test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"certificate name begins with a dot": {
			dir: "name_begins_with_dot",
			mockData: getCertificateMockData("ACTIVE", ".test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"empty certificate subject": {
			dir:      "empty_subject",
			mockData: getCertificateMockDataCustomSubject(&emptySubject),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"nil certificate subject": {
			dir:      "empty_subject",
			mockData: getCertificateMockDataCustomSubject(nil),
			init: func(d *certificateMockData, m *cloudcertificates.Mock) {
				d.mockListCertificates(m)
				d.mockGetCertificate(m)
			},
		},
		"API error": {
			mockData: getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
				&signedCertificatePEM, &trustChainPEM),
			init: func(_ *certificateMockData, m *cloudcertificates.Mock) {
				m.On("ListCertificates", mock.Anything, cloudcertificates.ListCertificatesRequest{
					PageSize: maxPageSize,
					Page:     1,
				}).Return(nil, cloudcertificates.ErrListCertificates).Once()
			},
			withError: ErrListingCloudCertificates.Error(),
		},
		"certificate does not exist": {
			dir: "active",
			mockData: getCertificateMockData("ACTIVE", "certificate_does_not_exist",
				&signedCertificatePEM, &trustChainPEM),
			init: func(_ *certificateMockData, m *cloudcertificates.Mock) {
				m.On("ListCertificates", mock.Anything, cloudcertificates.ListCertificatesRequest{
					PageSize: maxPageSize,
					Page:     1,
				}).Return(&cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{},
				}, nil).Once()
			},
			withError: "failed to fetch certificate: no certificate found with the name \"certificate_does_not_exist\"",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &cloudcertificates.Mock{}
			test.init(&test.mockData, m)

			if test.configSection == "" {
				test.configSection = "default"
			}

			if test.edgercPath == "" {
				test.edgercPath = "~/.edgerc"
			}

			params := createCloudCertificateParams{
				name:          test.mockData.name,
				edgercPath:    test.edgercPath,
				configSection: test.configSection,
				client:        m,
			}

			var tempDir string
			if test.dir != "" {
				tempDir = t.TempDir()
				var targets map[string]string
				if !test.withTemplatingError {
					targets = map[string]string{
						"cloudcertificate.tmpl": fmt.Sprintf("/%s/cloudcertificate.tf", tempDir),
						"variables.tmpl":        fmt.Sprintf("%s/variables.tf", tempDir),
						"import.tmpl":           fmt.Sprintf("%s/import.sh", tempDir),
					}
				} else {
					targets = map[string]string{
						"nosuchtemplate.tmpl": fmt.Sprintf("%s/nosuchtemplate.tf", tempDir),
					}
				}

				params.templateProcessor = templates.FSTemplateProcessor{
					TemplatesFS:     templateFiles,
					TemplateTargets: targets,
					AdditionalFuncs: additionalFunctions,
				}
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCloudCertificate(ctx, params)

			m.AssertExpectations(t)

			if test.withError != "" {
				assert.ErrorContains(t, err, test.withError)
				return
			}
			require.NoError(t, err)

			for _, f := range []string{"cloudcertificate.tf", "import.sh", "variables.tf"} {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("%s/%s", tempDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}

func TestExtractBaseName(t *testing.T) {

	tests := []struct {
		label       string
		name        string
		expBaseName string
	}{
		{"empty", "", ""},
		{"casual name", "foo", "foo"},
		{"renewed name", "foo.renewed.2025-05-01", "foo"},
		{"bad suffix", "foo.rotated.2025-05-01", "foo.rotated.2025-05-01"},
		{"non-existing date", "foo.renewed.2025-99-01", "foo.renewed.2025-99-01"},
		{"no basename", ".renewed.2025-05-01", ".renewed.2025-05-01"},
		{"no basename 2", "renewed.2025-05-01", "renewed.2025-05-01"},
		{"no date", "foo.renewed.", "foo.renewed."},
		{"no date 2", "foo.renewed", "foo.renewed"},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			res := extractBaseName(tc.name)
			assert.Equal(t, tc.expBaseName, res)
		})
	}
}

func getCertificateMockData(certStat, name string, signedCertificate, trustChain *string) certificateMockData {
	return certificateMockData{
		id:                "12345",
		name:              name,
		contractID:        "test_contract",
		baseName:          extractBaseName(name),
		keyType:           "RSA",
		keySize:           "2048",
		secureNetwork:     "ENHANCED_TLS",
		sans:              []string{"test.example.com", "test.example2.com"},
		certificateStatus: certStat,
		subject: &cloudcertificates.Subject{
			CommonName:   "test.example.com",
			Country:      "US",
			Organization: "Test Org",
			State:        "CA",
			Locality:     "Test City",
		},
		signedCertificatePEM: signedCertificate,
		trustChainPEM:        trustChain,
	}
}

func (d *certificateMockData) mockListCertificates(m *cloudcertificates.Mock) {
	m.On("ListCertificates", mock.Anything, cloudcertificates.ListCertificatesRequest{
		PageSize: maxPageSize,
		Page:     1,
	}).Return(&cloudcertificates.ListCertificatesResponse{
		Certificates: []cloudcertificates.Certificate{
			{
				CertificateID:     d.id,
				ContractID:        d.contractID,
				CertificateName:   d.name,
				KeyType:           cloudcertificates.CryptographicAlgorithm(d.keyType),
				KeySize:           cloudcertificates.KeySize(d.keySize),
				SecureNetwork:     d.secureNetwork,
				SANs:              d.sans,
				CertificateStatus: d.certificateStatus,
				Subject:           d.subject,
			},
		},
	}, nil).Once()
}

func (d *certificateMockData) mockGetCertificate(m *cloudcertificates.Mock) {
	m.On("GetCertificate", mock.Anything, cloudcertificates.GetCertificateRequest{
		CertificateID: d.id,
	}).Return(&cloudcertificates.GetCertificateResponse{
		Certificate: cloudcertificates.Certificate{
			CertificateID:        d.id,
			ContractID:           d.contractID,
			CertificateName:      d.name,
			KeyType:              cloudcertificates.CryptographicAlgorithm(d.keyType),
			KeySize:              cloudcertificates.KeySize(d.keySize),
			SecureNetwork:        d.secureNetwork,
			SANs:                 d.sans,
			CertificateStatus:    d.certificateStatus,
			Subject:              d.subject,
			SignedCertificatePEM: d.signedCertificatePEM,
			TrustChainPEM:        d.trustChainPEM,
		},
	}, nil).Once()
}

func getCertificateMockDataCustomSubject(subject *cloudcertificates.Subject) certificateMockData {
	cert := getCertificateMockData("ACTIVE", "test-name.example.com1234567890",
		&signedCertificatePEM, &trustChainPEM)
	cert.subject = subject
	return cert
}

func getCertificateMockDataNoTrustChain() certificateMockData {
	cert := getCertificateMockData("READY_FOR_USE", "test-name.example.com1234567890", &signedCertificatePEM, nil)
	cert.trustChainPEM = nil
	return cert
}

func getCertificateMockDataNoPEMs() certificateMockData {
	cert := getCertificateMockData("CSR_READY", "test-name.example.com1234567890", nil, nil)
	cert.trustChainPEM = nil
	cert.signedCertificatePEM = nil
	return cert
}
