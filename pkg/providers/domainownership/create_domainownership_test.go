package domainownership

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type domMockData struct {
	domainsRQ []domainownership.Domain
	domainsRS []domainownership.SearchDomainItem
}

func TestProcessDomainOwnershipTemplates(t *testing.T) {
	basicDomainsRQ := []domainownership.Domain{
		{
			DomainName:      "example.com",
			ValidationScope: "DOMAIN",
		},
		{
			DomainName:      "sub.example.com",
			ValidationScope: "HOST",
		},
	}

	basicDomainsRS := []domainownership.SearchDomainItem{
		{
			DomainName:       "example.com",
			ValidationScope:  "DOMAIN",
			ValidationMethod: ptr.To("DNS_CNAME"),
			ValidationLevel:  "FQDN",
			DomainStatus:     "VALIDATED",
		},
		{
			DomainName:       "sub.example.com",
			ValidationScope:  "HOST",
			ValidationMethod: ptr.To("DNS_TXT"),
			ValidationLevel:  "FQDN",
			DomainStatus:     "VALIDATED",
		},
	}

	domainsNoValidationScopeRQ := []domainownership.Domain{
		{
			DomainName:      "example.com",
			ValidationScope: "DOMAIN",
		},
		{
			DomainName:      "sub.example.com",
			ValidationScope: "DOMAIN",
		},
		{
			DomainName:      "sub.example.com",
			ValidationScope: "HOST",
		},
		{
			DomainName:      "sub.example.com",
			ValidationScope: "WILDCARD",
		},
	}

	tests := map[string]struct {
		dir                 string
		init                func(*domMockData, *domainownership.Mock)
		domains             string
		configSection       string
		edgercPath          string
		withError           string
		withTemplatingError bool
	}{
		"basic export": {
			dir:     "basic",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = basicDomainsRS
				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, same domain with two validation scopes": {
			dir:     "basic_multiple_validation_scopes",
			domains: "example.com:DOMAIN,example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = []domainownership.Domain{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
					},
					{
						DomainName:      "example.com",
						ValidationScope: "HOST",
					},
				}

				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:       "example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: ptr.To("DNS_CNAME"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATED",
					},
					{
						DomainName:       "example.com",
						ValidationScope:  "HOST",
						ValidationMethod: ptr.To("DNS_TXT"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATED",
					},
				}

				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, same domain without validation scope returned twice, once domain only with FQDN validation level": {
			dir:     "basic_single_domain",
			domains: "example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = []domainownership.Domain{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
					},
					{
						DomainName:      "example.com",
						ValidationScope: "HOST",
					},
					{
						DomainName:      "example.com",
						ValidationScope: "WILDCARD",
					},
				}

				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:       "example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: ptr.To("DNS_CNAME"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATED",
					},
					{
						DomainName:       "example.com",
						ValidationScope:  "HOST",
						ValidationMethod: ptr.To("DNS_TXT"),
						ValidationLevel:  "ROOT/WILDCARD",
						DomainStatus:     "VALIDATED",
					},
				}

				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, no validation scope": {
			dir:     "basic",
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = basicDomainsRS
				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, one not validated": {
			dir:     "basic_one_not_validated",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:       "example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: ptr.To("DNS_CNAME"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATED",
					},
					{
						DomainName:       "sub.example.com",
						ValidationScope:  "HOST",
						ValidationMethod: ptr.To("DNS_TXT"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATION_IN_PROGRESS",
					},
				}
				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, all not validated": {
			dir:     "basic_not_validated",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:       "example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: ptr.To("DNS_CNAME"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATION_IN_PROGRESS",
					},
					{
						DomainName:       "sub.example.com",
						ValidationScope:  "HOST",
						ValidationMethod: ptr.To("DNS_TXT"),
						ValidationLevel:  "FQDN",
						DomainStatus:     "VALIDATION_IN_PROGRESS",
					},
				}
				d.mockSearchDomains(m, nil)
			},
		},
		"basic export, without validation level of FQDN": {
			dir:     "basic",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "sub.example.com",
						ValidationScope: "HOST",
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'example.com' with validation scope 'DOMAIN' does not exist or the validation level was not 'FQDN'",
		},
		"basic export, no validation scope, without validation level of FQDN": {
			dir:     "basic",
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "sub.example.com",
						ValidationScope: "HOST",
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'sub.example.com' was not found or the validation level was not 'FQDN'",
		},
		"custom config section": {
			dir:     "custom_config_section",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = basicDomainsRS
				d.mockSearchDomains(m, nil)
			},
			configSection: "other-section",
		},
		"custom edgerc path": {
			dir:     "custom_edgerc_path",
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.domainsRS = basicDomainsRS
				d.mockSearchDomains(m, nil)
			},
			edgercPath: "/path/to/edgerc",
		},
		"error cannot list domains": {
			domains: "example.com:DOMAIN,sub.example.com:HOST",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = basicDomainsRQ
				d.mockSearchDomains(m, fmt.Errorf("listing error"))
			},
			withError: "error searching domains: listing error",
		},
		"error multiple validation scopes": {
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "sub.example.com",
						ValidationScope: "HOST",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "sub.example.com",
						ValidationScope: "WILDCARD",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'sub.example.com' specified without validation scope has multiple matches: [sub.example.com:HOST sub.example.com:WILDCARD]",
		},
		"error cannot specify domain with and without validation scope at the same time": {
			domains:   "example.com:DOMAIN,sub.example.com,sub.example.com:HOST",
			withError: "error parsing domains: domain 'sub.example.com' specified multiple times with and without validation scope",
		},
		"error cannot specify domain with the same validation scope twice": {
			domains:   "example.com:DOMAIN,example.com:DOMAIN,sub.example.com:HOST",
			withError: "error parsing domains: domain 'example.com' with validation scope 'DOMAIN' specified multiple times",
		},
		"error cannot specify domain without validation scope twice": {
			domains:   "example.com,example.com,sub.example.com:HOST",
			withError: "error parsing domains: domain 'example.com' specified multiple times without validation scope",
		},
		"error invalid format for domain": {
			domains:   "example.com:DOMAIN:foo,sub.example.com,sub.example.com:HOST",
			withError: "error parsing domains: invalid domain format for 'example.com:DOMAIN:foo', it should be 'domain' or 'domain:validation_scope' separated with comma",
		},
		"error invalid validation scope": {
			domains:   "example.com:foo,sub.example.com,sub.example.com:HOST",
			withError: "error parsing domains: invalid validation scope 'foo' for domain 'example.com', it should be 'HOST', 'WILDCARD' or 'DOMAIN'",
		},
		"error domain with validation scope not found": {
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "sub.example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'example.com' with validation scope 'DOMAIN' does not exist or the validation level was not 'FQDN'",
		},
		"error domain without validation scope not found": {
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'sub.example.com' was not found or the validation level was not 'FQDN'",
		},
		"error domain with validation scope not found with validation level FQDN": {
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "sub.example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'example.com' with validation scope 'DOMAIN' does not exist or the validation level was not 'FQDN'",
		},
		"error domain without validation scope not found with validation level FQDN": {
			domains: "example.com:DOMAIN,sub.example.com",
			init: func(d *domMockData, m *domainownership.Mock) {
				d.domainsRQ = domainsNoValidationScopeRQ
				d.domainsRS = []domainownership.SearchDomainItem{
					{
						DomainName:      "example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					}, {
						DomainName:      "sub.example.com",
						ValidationScope: "DOMAIN",
						ValidationLevel: "FQDN",
						DomainStatus:    "VALIDATED",
					},
				}
				d.mockSearchDomains(m, nil)
			},
			withError: "error parsing domains: domain 'example.com' with validation scope 'DOMAIN' does not exist or the validation level was not 'FQDN'",
		},
		"error more than 1000 domains": {
			domains: func() string {
				s := make([]string, 1001)
				for i := 0; i < 1001; i++ {
					s[i] = fmt.Sprintf("domain%d.com", i)
				}
				return strings.Join(s, ",")
			}(),
			withError: "error parsing domains: maximum of 1000 domains can be processed at once",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &domainownership.Mock{}
			d := domMockData{}
			if test.init != nil {
				test.init(&d, m)
			}

			if test.configSection == "" {
				test.configSection = "default"
			}

			if test.edgercPath == "" {
				test.edgercPath = "~/.edgerc"
			}

			params := createDomainOwnershipParams{
				domains:       test.domains,
				configSection: test.configSection,
				edgercPath:    test.edgercPath,
			}

			if test.dir != "" {
				targets := map[string]string{
					"domainownership.tmpl": fmt.Sprintf("./testdata/res/%s/domainownership.tf", test.dir),
					"variables.tmpl":       fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"import.tmpl":          fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				}

				params.templateProcessor = templates.FSTemplateProcessor{
					TemplatesFS:     templateFiles,
					AdditionalFuncs: additionalFunctions,
					TemplateTargets: targets,
				}
				require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createDomainOwnership(ctx, params, m)

			m.AssertExpectations(t)

			if test.withError != "" {
				assert.ErrorContains(t, err, test.withError)
				return
			}
			require.NoError(t, err)

			for _, f := range []string{"domainownership.tf", "import.sh", "variables.tf"} {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}

func (d *domMockData) mockSearchDomains(m *domainownership.Mock, err error) {
	req := domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: d.domainsRQ,
		},
	}

	res := &domainownership.SearchDomainsResponse{
		Domains: d.domainsRS,
	}

	m.On("SearchDomains", mock.Anything, req).Return(res, err).Once()
}
