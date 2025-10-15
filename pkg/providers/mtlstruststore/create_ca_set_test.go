package mtlstruststore

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type caSetMockData struct {
	name                    string
	id                      string
	latestVersion           *int64
	stagingVersion          *int64
	productionVersion       *int64
	userVersion             int64
	requestedVersion        int64
	caSets                  []string
	caSetIDs                []string
	certificates            []string
	description             *string
	versionDescription      *string
	certificateDescriptions []*string
}

var cert1 = `-----BEGIN CERTIFICATE-----
FAKECERTSTARTSEQ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKL
MNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN
OPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOO
PQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNEND
SEQ==
-----END CERTIFICATE-----`

var cert2 = `-----BEGIN CERTIFICATE-----
ZYXWVUTSRQPONMLKJIHGFEDCBA9876543210zyxwvutsrqponmlkjihgfedcba988
76543210ZYXWVUTSRQPONMLKJIHGFEDCBA9876543210zyxwvutsrqponmlkjihgf
edcba9876543210ZYXWVUTSRQPONMLKJIHGFEDCBA9876543210zyxwvutsrqponm
lkjihgfedcba9876543210ZYXWVUTSRQPONMLKJIHGFEDCBA9876543210zyxwvut
ENDSEQ==
-----END CERTIFICATE-----`

func defaultCaSetMockData() caSetMockData {
	return caSetMockData{
		name:                    "test-ca-set-name",
		id:                      "12345",
		latestVersion:           ptr.To(int64(5)),
		stagingVersion:          ptr.To(int64(4)),
		productionVersion:       ptr.To(int64(3)),
		caSets:                  []string{"test-ca-set-name", "other-ca-set-name"},
		caSetIDs:                []string{"12345", "67890"},
		certificates:            []string{cert1},
		certificateDescriptions: []*string{nil},
	}
}

func TestProcessMtlsTruststoreTemplates(t *testing.T) {
	tests := map[string]struct {
		dir                 string
		init                func(*caSetMockData, *mtlstruststore.Mock)
		configSection       string
		withError           string
		withTemplatingError bool
	}{
		"basic export generation": {
			dir: "basic",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.mockAll(m)
			},
		},
		"basic export generation with descriptions": {
			dir: "basic_descriptions",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.description = ptr.To("test CA set description")
				d.versionDescription = ptr.To("test CA set version description")
				d.certificateDescriptions = []*string{ptr.To("test certificate description")}
				d.noUserVersion()
				d.mockAll(m)
			},
		},
		"basic export generation with descriptions containing multilines": {
			dir: "with_multiline",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.description = ptr.To("test CA set\ndescription")
				d.versionDescription = ptr.To("test CA set\n\nversion description")
				d.certificateDescriptions = []*string{ptr.To("test certificate\ndescription")}
				d.noUserVersion()
				d.mockAll(m)
			},
		},
		"basic export generation with descriptions containing multilines - empty line at the end": {
			dir: "with_multiline2",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.description = ptr.To("test CA set\ndescription\n")
				d.versionDescription = ptr.To("test CA set\n\nversion description\n")
				d.certificateDescriptions = []*string{ptr.To("test certificate\ndescription\n")}
				d.noUserVersion()
				d.mockAll(m)
			},
		},
		"export two certificates": {
			dir: "two_certificates",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.description = ptr.To("test CA set description")
				d.versionDescription = ptr.To("test CA set version description")
				d.noUserVersion()
				d.certificates = []string{cert1, cert2}
				d.certificateDescriptions = []*string{
					ptr.To("test certificate description 1"),
					ptr.To("test certificate description 2"),
				}
				d.mockAll(m)
			},
		},
		"only staging activation": {
			dir: "staging_activation",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.productionVersion = nil
				d.mockAll(m)
			},
		},
		"only production activation": {
			dir: "production_activation",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.stagingVersion = nil
				d.mockAll(m)
			},
		},
		"latest version on staging": {
			dir: "latest_version_staging",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.stagingVersion = ptr.To(*d.latestVersion)
				d.mockAll(m)
			},
		},
		"latest version on production": {
			dir: "latest_version_production",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.productionVersion = ptr.To(*d.latestVersion)
				d.mockAll(m)
			},
		},
		"no activations": {
			dir: "no_activations",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.stagingVersion = nil
				d.productionVersion = nil
				d.mockAll(m)
			},
		},
		"requested specific version": {
			dir: "specific_version",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.userVersion = 2
				d.requestedVersion = 2
				d.mockAll(m)
			},
		},
		"custom config section": {
			dir: "custom_config_section",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.mockAll(m)
			},
			configSection: "other-section",
		},
		"funny CA set name handled properly": {
			dir: "funny_ca_set_name",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.name = "a_funny-set.name"
				d.caSets = []string{"a_funny-set.name", "regular-name"}
				d.mockAll(m)
			},
		},
		"error cannot list CA sets": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.mockListCASets(m, fmt.Errorf("listing error"))
			},
			withError: "error finding CA set: failed to list CA sets: listing error",
		},
		"error no such CA set": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.name = "non-existing-ca-set"
				d.mockListCASets(m, nil)
			},
			withError: "error finding CA set: no CA set found with name 'non-existing-ca-set'",
		},
		"error multiple CA sets": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.caSets = []string{"test-ca-set-name", "test-ca-set-name"}
				d.mockListCASets(m, nil)
			},
			withError: "error finding CA set: multiple CA sets found with name 'test-ca-set-name'",
		},
		"error fetching single CA set": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.mockListCASets(m, nil)
				d.mockGetCASet(m, fmt.Errorf("404 no set for given ID"))
			},
			withError: "error fetching CA set: 404 no set for given ID",
		},
		"error no version specified and CA set has no versions": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.latestVersion = nil
				d.mockListCASets(m, nil)
				d.mockGetCASet(m, nil)
			},
			withError: "error fetching CA set: CA set 'test-ca-set-name' has no versions",
		},
		"error fetching CA set version": {
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.userVersion = 2
				d.requestedVersion = 2
				d.mockListCASets(m, nil)
				d.mockGetCASet(m, nil)
				d.mockGetCASetVersion(m, fmt.Errorf("CA set version is not found"))
			},
			withError: "error fetching CA set version: CA set version is not found",
		},
		"templating error": {
			dir: "templating_error",
			init: func(d *caSetMockData, m *mtlstruststore.Mock) {
				d.noUserVersion()
				d.mockAll(m)
			},
			withTemplatingError: true,
			withError:           "error saving terraform project files: no template file: nosuchtemplate.tmpl",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mtlstruststore.Mock{}
			d := defaultCaSetMockData()
			test.init(&d, m)

			if test.configSection == "" {
				test.configSection = "default"
			}

			params := createCASetParams{
				name:          d.name,
				userVersion:   d.userVersion,
				configSection: test.configSection,
				client:        m,
			}

			if test.dir != "" {
				var targets map[string]string
				if !test.withTemplatingError {
					targets = map[string]string{
						"mtlstruststore.tmpl": fmt.Sprintf("./testdata/res/%s/mtlstruststore.tf", test.dir),
						"variables.tmpl":      fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
						"imports.tmpl":        fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
					}
				} else {
					targets = map[string]string{
						"nosuchtemplate.tmpl": fmt.Sprintf("./testdata/res/%s/nosuchtemplate.tf", test.dir),
					}
				}

				params.templateProcessor = templates.FSTemplateProcessor{
					TemplatesFS:     templateFiles,
					AdditionalFuncs: additionalFunctions,
					TemplateTargets: targets,
				}
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createCASet(ctx, params)

			m.AssertExpectations(t)

			if test.withError != "" {
				assert.ErrorContains(t, err, test.withError)
				return
			}
			require.NoError(t, err)

			for _, f := range []string{"mtlstruststore.tf", "import.sh", "variables.tf"} {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}

// If the user did not specify a version, we expect the latest version will be fetched.
func (d *caSetMockData) noUserVersion() {
	d.userVersion = 0
	d.requestedVersion = *d.latestVersion
}

func (d *caSetMockData) mockListCASets(m *mtlstruststore.Mock, err error) {
	req := mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: d.name,
	}

	res := &mtlstruststore.ListCASetsResponse{}
	for i, caSet := range d.caSets {
		res.CASets = append(res.CASets, mtlstruststore.CASetResponse{
			CASetName:   caSet,
			CASetStatus: "NOT_DELETED",
			CASetID:     d.caSetIDs[i],
		})
	}

	m.On("ListCASets", mock.Anything, req).Return(res, err).Once()
}

func (d *caSetMockData) mockGetCASet(m *mtlstruststore.Mock, err error) {
	req := mtlstruststore.GetCASetRequest{
		CASetID: d.id,
	}

	res := &mtlstruststore.GetCASetResponse{
		CASetName:         d.name,
		CASetID:           d.id,
		LatestVersion:     d.latestVersion,
		StagingVersion:    d.stagingVersion,
		ProductionVersion: d.productionVersion,
		Description:       d.description,
	}

	m.On("GetCASet", mock.Anything, req).Return(res, err).Once()
}

func (d *caSetMockData) mockGetCASetVersion(m *mtlstruststore.Mock, err error) {
	req := mtlstruststore.GetCASetVersionRequest{
		CASetID: d.id,
		Version: d.requestedVersion,
	}

	var certs []mtlstruststore.CertificateResponse
	for i, cert := range d.certificates {
		certs = append(certs, mtlstruststore.CertificateResponse{
			CertificatePEM: cert,
			Description:    d.certificateDescriptions[i],
		})
	}

	res := &mtlstruststore.GetCASetVersionResponse{
		Certificates: certs,
		Description:  d.versionDescription,
	}

	m.On("GetCASetVersion", mock.Anything, req).Return(res, err).Once()
}

func (d *caSetMockData) mockAll(m *mtlstruststore.Mock) {
	d.mockListCASets(m, nil)
	d.mockGetCASet(m, nil)
	d.mockGetCASetVersion(m, nil)
}
