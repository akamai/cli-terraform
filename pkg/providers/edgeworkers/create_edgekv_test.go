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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

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

	expectListGroupsWithinNamespace = func(e *edgeworkers.Mock, network edgeworkers.NamespaceNetwork, name string, groups []string, err error) *mock.Call {
		call := e.On(
			"ListGroupsWithinNamespace",
			mock.Anything,
			edgeworkers.ListGroupsWithinNamespaceRequest{
				Network:     network,
				NamespaceID: name,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(groups, nil)
	}

	expectListItems = func(e *edgeworkers.Mock, network edgeworkers.NamespaceNetwork, name string, group string, items *edgeworkers.ListItemsResponse, err error) *mock.Call {
		call := e.On(
			"ListItems",
			mock.Anything,
			edgeworkers.ListItemsRequest{
				ItemsRequestParams: edgeworkers.ItemsRequestParams{
					Network:     edgeworkers.ItemNetwork(network),
					NamespaceID: name,
					GroupID:     group,
				},
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(items, nil)
	}

	expectGetItem = func(e *edgeworkers.Mock, network edgeworkers.NamespaceNetwork, name, group, itemID, item string, err error) *mock.Call {
		call := e.On(
			"GetItem",
			mock.Anything,
			edgeworkers.GetItemRequest{
				ItemID: itemID,
				ItemsRequestParams: edgeworkers.ItemsRequestParams{
					Network:     edgeworkers.ItemNetwork(network),
					NamespaceID: name,
					GroupID:     group,
				},
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		it := edgeworkers.Item(item)
		return call.Return(&it, nil)
	}

	expectProcessTemplates = func(p *templates.MockProcessor, network edgeworkers.NamespaceNetwork, name string, geoLocation string,
		retention int, groupID *int, section string, items map[string]map[string]edgeworkers.Item, err error) *mock.Call {
		var tfData TFEdgeKVData
		tfData = TFEdgeKVData{
			Name:        name,
			Network:     network,
			Retention:   retention,
			GeoLocation: geoLocation,
			Section:     section,
			GroupItems:  items,
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

	emptyItems = map[string]map[string]edgeworkers.Item{}

	items = map[string]map[string]edgeworkers.Item{
		"group1": {
			"item1.1": edgeworkers.Item("value1.1"),
			"item1.2": edgeworkers.Item("value1.2"),
		},
		"group2": {
			"item2.1": edgeworkers.Item("value2.1"),
			"item2.2": edgeworkers.Item("value\n2.2"),
		},
	}
)

func TestCreateEdgeKV(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init      func(*edgeworkers.Mock, *templates.MockProcessor)
		withError error
	}{
		"fetch edgekv based on namespace and network": {
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), nil).Once()
				expectListGroupsWithinNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", []string{}, nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, intPtr(123), section, emptyItems, nil).Once()
			},
		},
		"fetch edgekv based on namespace and network with group items": {
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), nil).Once()
				expectListGroupsWithinNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", []string{"group1", "group2"}, nil).Once()
				expectListItems(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group1", &edgeworkers.ListItemsResponse{"item1.1", "item1.2"}, nil).Once()
				expectListItems(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group2", &edgeworkers.ListItemsResponse{"item2.1", "item2.2"}, nil).Once()
				expectGetItem(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group1", "item1.1", "value1.1", nil).Once()
				expectGetItem(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group1", "item1.2", "value1.2", nil).Once()
				expectGetItem(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group2", "item2.1", "value2.1", nil).Once()
				expectGetItem(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "group2", "item2.2", "value\n2.2", nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, intPtr(123), section, items, nil).Once()
			},
		},
		"fetch edgekv based on namespace and network with no group_id returned": {
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), nil, nil).Once()
				expectListGroupsWithinNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", []string{}, nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, nil, section, emptyItems, nil).Once()
			},
		},
		"error fetching edgekv": {
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), fmt.Errorf("error")).Once()
			},
			withError: ErrFetchingEdgeKV,
		},
		"error processing template": {
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeKVNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", intPtr(0), intPtr(123), nil).Once()
				expectListGroupsWithinNamespace(e, edgeworkers.NamespaceStagingNetwork, "test_namespace", []string{}, nil).Once()
				expectProcessTemplates(p, edgeworkers.NamespaceStagingNetwork, "test_namespace", "EU", 0, intPtr(123), section, emptyItems, fmt.Errorf("error")).Once()
			},
			withError: templates.ErrSavingFiles,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			me := new(edgeworkers.Mock)
			mp := new(templates.MockProcessor)
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
		"edgekv with staging network and items": {
			givenData: TFEdgeKVData{
				Name:        "test_namespace",
				Network:     edgeworkers.NamespaceStagingNetwork,
				GroupID:     123,
				Retention:   0,
				GeoLocation: "EU",
				Section:     "test_section",
				GroupItems:  items,
			},
			dir:          "edgekv_with_staging_network_and_items",
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
					"Escape": tools.Escape,
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
