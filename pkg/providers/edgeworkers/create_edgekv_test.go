package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) ProcessTemplates(i interface{}) error {
	args := m.Called(i)
	return args.Error(0)
}

var (
	intPtr = func(i int) *int {
		return &i
	}

	expectGetEdgeKVNamespace = func(e *edgeworkers.Mock, network edgeworkers.NamespaceNetwork, name string, geoLocation string,
		retention *int, groupID *int, err error) *mock.Call {
		call := e.On(
			"GetEdgeKVNamespace",
			mock.Anything,
			edgeworkers.GetEdgeKVNamespaceRequest{
				Network: network,
				Name:    name,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(
			&edgeworkers.Namespace{
				Name:        name,
				GeoLocation: geoLocation,
				Retention:   retention,
				GroupID:     groupID,
			}, nil)
	}

	expectProcessTemplates = func(p *mockProcessor, network edgeworkers.NamespaceNetwork, name string, geoLocation string,
		retention int, groupID *int, section string, err error) *mock.Call {
		var tfData TFEdgeKVData
		tfData = TFEdgeKVData{
			Name:        name,
			Network:     network,
			Retention:   retention,
			GeoLocation: geoLocation,
			Section:     section,
		}
		if groupID != nil {
			tfData.GroupID = *groupID
		}

		call := p.On(
			"ProcessTemplates",
			tfData,
		)
		if err != nil {
			return call.Return(err)
		}
		return call.Return(nil)
	}
)

func TestCreateEdgeKV(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init      func(*edgeworkers.Mock, *mockProcessor)
		withError error
	}{
		"fetch edgekv based on namespace and network": {
			init: func(e *edgeworkers.Mock, p *mockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, intPtr(123), section, nil).Once()
			},
		},
		"fetch edgekv based on namespace and network with no group_id returned": {
			init: func(e *edgeworkers.Mock, p *mockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), nil, nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, nil, section, nil).Once()
			},
		},
		"error fetching edgekv": {
			init: func(e *edgeworkers.Mock, p *mockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), fmt.Errorf("error")).Once()
			},
			withError: ErrFetchingEdgeKV,
		},
		"error processing template": {
			init: func(e *edgeworkers.Mock, p *mockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, intPtr(123), section, fmt.Errorf("error")).Once()
			},
			withError: templates.ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			me := new(edgeworkers.Mock)
			mp := new(mockProcessor)
			test.init(me, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createEdgeKV(ctx, "test_namespace", edgeworkers.NamespaceStagingNetwork, section, me, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			me.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessEdgeKVTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFEdgeKVData
		dir          string
		filesToCheck []string
	}{
		"edgekv with staging network": {
			givenData: TFEdgeKVData{
				Name:        "test_namespace",
				Network:     edgeworkers.NamespaceStagingNetwork,
				GroupID:     123,
				Retention:   0,
				GeoLocation: "EU",
				Section:     "test_section",
			},
			dir:          "edgekv_with_staging_network",
			filesToCheck: []string{"edgekv.tf", "variables.tf", "import.sh"},
		},
		"edgekv with staging network no group_id": {
			givenData: TFEdgeKVData{
				Name:        "test_namespace",
				Network:     edgeworkers.NamespaceStagingNetwork,
				Retention:   0,
				GeoLocation: "EU",
				Section:     "test_section",
			},
			dir:          "edgekv_with_staging_network_no_group_id",
			filesToCheck: []string{"edgekv.tf", "variables.tf", "import.sh"},
		},
		"edgekv with production network": {
			givenData: TFEdgeKVData{
				Name:        "test_namespace",
				Network:     edgeworkers.NamespaceProductionNetwork,
				GroupID:     123,
				Retention:   0,
				GeoLocation: "EU",
				Section:     "test_section",
			},
			dir:          "edgekv_with_production_network",
			filesToCheck: []string{"edgekv.tf", "variables.tf", "import.sh"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"edgekv.tmpl":           fmt.Sprintf("./testdata/res/%s/edgekv.tf", test.dir),
					"edgekv-variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"edgekv-imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
				},
				AdditionalFuncs: template.FuncMap{
					"ToLower": func(network edgeworkers.ActivationNetwork) string {
						return strings.ToLower(string(network))
					},
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
