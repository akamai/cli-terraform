package mtlskeystore

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type certificateMockData struct {
	name               string
	id                 int64
	contractID         string
	groupID            int64
	geography          string
	keyAlgorithm       string
	notificationEmails []string
	secureNetwork      string
	subject            string
	signer             string
	versions           []mtlskeystore.ClientCertificateVersion
}

func (d *certificateMockData) setVersions(versions []mtlskeystore.ClientCertificateVersion) {
	d.versions = versions
}

func akamaiCertificateMockData() certificateMockData {
	return certificateMockData{
		name:               "test-akamai-cert",
		id:                 12345,
		contractID:         "C-0NTR4CT",
		groupID:            98765,
		geography:          "CORE",
		keyAlgorithm:       "RSA",
		notificationEmails: []string{"test@mail.com"},
		secureNetwork:      "ENHANCED_TLS",
		subject:            "/C=US/O=Akamai Technologies, Inc./OU=12345 C-0NTR4CT 98765/CN=testCommonName/",
		signer:             "AKAMAI",
		versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:      2,
				CreatedDate:  "2024-10-01T12:00:00Z",
				VersionAlias: ptr.To("CURRENT"),
				Status:       mtlskeystore.DeploymentPending,
			},
			{
				Version:      1,
				CreatedDate:  "2023-10-01T12:00:00Z",
				VersionAlias: ptr.To("PREVIOUS"),
				Status:       mtlskeystore.Deployed,
			},
		},
	}
}

func thirdPartyCertificateMockData() certificateMockData {
	return certificateMockData{
		name:               "test-third-party-cert",
		id:                 12345,
		contractID:         "C-0NTR4CT",
		groupID:            98765,
		geography:          "RUSSIA_AND_CORE",
		keyAlgorithm:       "ECDSA",
		notificationEmails: []string{"test@mail.com"},
		secureNetwork:      "STANDARD_TLS",
		subject:            "/C=US/O=Akamai Technologies, Inc./OU=12345 C-0NTR4CT 98765/CN=testCommonName/",
		signer:             "THIRD_PARTY",
		versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     3,
				CreatedDate: "2025-10-01T12:00:00Z",
				Status:      mtlskeystore.AwaitingSigned,
			},
			{
				Version:     2,
				CreatedDate: "2024-10-01T12:00:00Z",
				Status:      mtlskeystore.Deployed,
			},
			{
				Version:     1,
				CreatedDate: "2023-10-01T12:00:00Z",
				Status:      mtlskeystore.Deployed,
			},
		},
	}
}

func TestProcessMTLSKeystoreTemplates(t *testing.T) {
	tests := map[string]struct {
		dir                 string
		init                func(*certificateMockData, *mtlskeystore.Mock)
		mockData            certificateMockData
		configSection       string
		withError           string
		withTemplatingError bool
	}{
		"akamai signer client certificate with 2 versions": {
			dir:      "akamai",
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockAll(m)
			},
		},
		"akamai signer client certificate with 1 version": {
			dir:      "akamai_one_version",
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
				})
				d.mockAll(m)
			},
		},
		"akamai signer client certificate with 1 deployed version and 1 pending deletion": {
			dir:      "akamai_one_version",
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.DeletePending,
					},
					{
						Version:     2,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
				})
				d.mockAll(m)
			},
		},
		"third party signer client certificate with multiple versions": {
			dir:      "third_party",
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockAll(m)
			},
		},
		"third party signer client certificate with one version": {
			dir:      "third_party_one_version",
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
				})
				d.mockAll(m)
			},
		},
		"third party signer client certificate with multiple versions and non-nil aliases": {
			dir:      "third_party",
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:      3,
						CreatedDate:  "2025-10-01T13:00:00Z",
						VersionAlias: ptr.To("CURRENT"),
						Status:       mtlskeystore.Deployed,
					},
					{
						Version:     3,
						CreatedDate: "2025-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
					{
						Version:      2,
						CreatedDate:  "2024-10-01T13:00:00Z",
						VersionAlias: ptr.To("PREVIOUS"),
						Status:       mtlskeystore.Deployed,
					},
					{
						Version:     2,
						CreatedDate: "2024-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.Deployed,
					},
				})
				d.mockGetClientCertificate(m, nil)
				d.mockListClientCertificateVersions(m, nil)
			},
		},
		"custom config section": {
			dir:      "custom_config_section",
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockAll(m)
			},
			configSection: "other-section",
		},
		"expect error: akamai signer certificate with pending deletion version": {
			dir:      "akamai",
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.DeletePending,
					},
				})
				d.mockAll(m)
			},
			withError: "certificate has no versions or the versions are pending deletion",
		},
		"expect error: third party signer certificate with no versions": {
			dir:      "third_party",
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{})
				d.mockAll(m)
			},
			withError: "error populating terraform data: certificate with ID '12345' has no versions or the versions are pending deletion",
		},
		"expect error: third party signer certificate with 2 versions pending deletion": {
			dir:      "third_party",
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.setVersions([]mtlskeystore.ClientCertificateVersion{
					{
						Version:     1,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.DeletePending,
					},
					{
						Version:     2,
						CreatedDate: "2023-10-01T12:00:00Z",
						Status:      mtlskeystore.DeletePending,
					},
				})
				d.mockAll(m)
			},
			withError: "error populating terraform data: certificate with ID '12345' has no versions or the versions are pending deletion",
		},
		"expect error: cannot get client certificate": {
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockGetClientCertificate(m, fmt.Errorf("get client certificate error"))
			},
			withError: "error fetching client certificate: get client certificate error",
		},
		"expect error: cannot get client certificate versions": {
			mockData: thirdPartyCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockGetClientCertificate(m, nil)
				d.mockListClientCertificateVersions(m, fmt.Errorf("get client certificate versions error"))
			},
			withError: "error fetching client certificate versions: get client certificate versions error",
		},
		"expect error: subject does not contain group and contract": {
			mockData: akamaiCertificateMockData(),
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.subject = "some custom subject value"
				d.mockAll(m)
			},
			withError: "error populating terraform data: unable to extract group and contract from certificate subject: unexpected format: 'some custom subject value'.\nRe-run with following arguments: <certificate_id>  <group_id> <contract_id>",
		},
		"templating error": {
			mockData: akamaiCertificateMockData(),
			dir:      "templating_error",
			init: func(d *certificateMockData, m *mtlskeystore.Mock) {
				d.mockAll(m)
			},
			withTemplatingError: true,
			withError:           "error saving terraform project files: no template file: nosuchtemplate.tmpl",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mtlskeystore.Mock{}
			test.init(&test.mockData, m)

			if test.configSection == "" {
				test.configSection = "default"
			}

			params := createCertificateParams{
				id:            test.mockData.id,
				configSection: test.configSection,
				client:        m,
			}

			var tempDir string
			if test.dir != "" {
				tempDir = t.TempDir()
				var targets map[string]string
				if !test.withTemplatingError {
					targets = map[string]string{
						"mtlskeystore.tmpl": fmt.Sprintf("/%s/mtlskeystore.tf", tempDir),
						"variables.tmpl":    fmt.Sprintf("%s/variables.tf", tempDir),
						"imports.tmpl":      fmt.Sprintf("%s/import.sh", tempDir),
					}
				} else {
					targets = map[string]string{
						"nosuchtemplate.tmpl": fmt.Sprintf("%s/nosuchtemplate.tf", tempDir),
					}
				}

				params.templateProcessor = templates.FSTemplateProcessor{
					TemplatesFS:     templateFiles,
					AdditionalFuncs: additionalFunctions,
					TemplateTargets: targets,
				}
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCertificate(ctx, params)

			m.AssertExpectations(t)

			if test.withError != "" {
				assert.ErrorContains(t, err, test.withError)
				return
			}
			require.NoError(t, err)

			for _, f := range []string{"mtlskeystore.tf", "import.sh", "variables.tf"} {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("%s/%s", tempDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}

func (d *certificateMockData) mockGetClientCertificate(m *mtlskeystore.Mock, err error) {
	req := mtlskeystore.GetClientCertificateRequest{
		CertificateID: d.id,
	}

	res := &mtlskeystore.GetClientCertificateResponse{
		CertificateID:      d.id,
		CertificateName:    d.name,
		Geography:          mtlskeystore.Geography(d.geography),
		KeyAlgorithm:       mtlskeystore.CryptographicAlgorithm(d.keyAlgorithm),
		NotificationEmails: d.notificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(d.secureNetwork),
		Signer:             mtlskeystore.Signer(d.signer),
		Subject:            d.subject,
	}

	m.On("GetClientCertificate", mock.Anything, req).Return(res, err).Once()
}

func (d *certificateMockData) mockListClientCertificateVersions(m *mtlskeystore.Mock, err error) {
	req := mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: d.id,
	}

	res := &mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{},
	}

	for _, v := range d.versions {
		res.Versions = append(res.Versions, mtlskeystore.ClientCertificateVersion{
			Version:      v.Version,
			CreatedDate:  v.CreatedDate,
			VersionAlias: v.VersionAlias,
			Status:       v.Status,
		})
	}

	m.On("ListClientCertificateVersions", mock.Anything, req).Return(res, err).Once()
}

func (d *certificateMockData) mockAll(m *mtlskeystore.Mock) {
	d.mockGetClientCertificate(m, nil)
	d.mockListClientCertificateVersions(m, nil)
}

func TestExtractContractAndGroup(t *testing.T) {
	tests := map[string]struct {
		subject        string
		expectContract string
		expectGroup    string
		expectError    string
	}{
		"standard subject with contract and group": {
			subject:        "/C=US/O=Akamai Technologies, Inc./OU=12345 C-0NTR4CT 98765/CN=testCommonName/",
			expectContract: "C-0NTR4CT",
			expectGroup:    "98765",
		},
		"subject with extra spaces": {
			subject:        "/C=US/O=Org/OU=  12345   C-0NTR4CT   98765  /CN=foo/",
			expectContract: "C-0NTR4CT",
			expectGroup:    "98765",
		},
		"subject with multiple fields in OU": {
			subject:        "/C=US/O=Org/OU=foo bar C-0NTR4CT 98765/CN=foo/",
			expectContract: "C-0NTR4CT",
			expectGroup:    "98765",
		},
		"subject with OU but not enough fields": {
			subject:     "/C=US/O=Org/OU=onlyone/CN=foo/",
			expectError: "no group or contract",
		},
		"subject with OU but only one field": {
			subject:     "/C=US/O=Org/OU=12345/CN=foo/",
			expectError: "no group or contract",
		},
		"subject with dashes in contract": {
			subject:        "/C=US/O=Org/OU=foo-bar C-123-456 78910/CN=foo/",
			expectContract: "C-123-456",
			expectGroup:    "78910",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctr, grp, err := extractContractAndGroup(tt.subject)
			if tt.expectError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectContract, ctr)
				assert.Equal(t, tt.expectGroup, grp)
			}
		})
	}
}
