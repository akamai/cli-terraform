package edgeworkers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/edgeworkers"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/terminal"
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
	expectEdgeWorkerProcessTemplates = func(p *templates.MockProcessor, edgeWorkerID int, name string, groupID int64, resourceTierID int,
		localBundle, edgercPath, section, note string, err error) *mock.Call {
		tfData := TFEdgeWorkerData{
			EdgeWorkerID:   edgeWorkerID,
			Name:           name,
			GroupID:        groupID,
			ResourceTierID: resourceTierID,
			LocalBundle:    localBundle,
			EdgercPath:     edgercPath,
			Section:        section,
			Note:           note,
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

	expectGetEdgeWorkerID = func(e *edgeworkers.Mock, edgeWorkerID int, name string, groupID int64, resourceTierID int, err error) *mock.Call {
		call := e.On(
			"GetEdgeWorkerID",
			mock.Anything,
			edgeworkers.GetEdgeWorkerIDRequest{
				EdgeWorkerID: edgeWorkerID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(
			&edgeworkers.EdgeWorkerID{
				EdgeWorkerID:   edgeWorkerID,
				Name:           name,
				GroupID:        groupID,
				ResourceTierID: resourceTierID,
			}, nil)
	}

	expectGetEdgeWorkerVersionContent = func(e *edgeworkers.Mock, edgeWorkerID int, version string, versionContent *bytes.Buffer, err error) *mock.Call {
		call := e.On(
			"GetEdgeWorkerVersionContent",
			mock.Anything,
			edgeworkers.GetEdgeWorkerVersionContentRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		return call.Return(
			&edgeworkers.Bundle{
				Reader: versionContent,
			}, nil)
	}

	expectListEdgeWorkerVersions = func(e *edgeworkers.Mock, edgeWorkerID int, empty bool, err error) *mock.Call {
		var versions []edgeworkers.EdgeWorkerVersion
		call := e.On(
			"ListEdgeWorkerVersions",
			mock.Anything,
			edgeworkers.ListEdgeWorkerVersionsRequest{
				EdgeWorkerID: edgeWorkerID,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}
		if !empty {
			versions = []edgeworkers.EdgeWorkerVersion{
				{
					EdgeWorkerID:   edgeWorkerID,
					Version:        "1.23",
					AccountID:      "B-3-WNKA6P",
					Checksum:       "868f28f16c26f46d418d83e24973520534d9ea4e4dbfd8a69ab00c1c37f28ca4",
					SequenceNumber: 3,
					CreatedBy:      "jsmith",
					CreatedTime:    "2021-12-20T08:28:37Z",
				},
				{
					EdgeWorkerID:   edgeWorkerID,
					Version:        "1.24.5",
					AccountID:      "B-3-WNKA6P",
					Checksum:       "ad9c18a7f2ed5d7bbcd31c55b94a0a00ae1771c6a15fd9265aeae08f5ef41e1f",
					SequenceNumber: 4,
					CreatedBy:      "jsmith",
					CreatedTime:    "2021-12-20T09:39:48Z",
				},
			}
		}
		return call.Return(
			&edgeworkers.ListEdgeWorkerVersionsResponse{
				EdgeWorkerVersions: versions,
			}, nil)
	}

	expectListActivations = func(e *edgeworkers.Mock, edgeWorkerID int, version string, err error) *mock.Call {
		call := e.On(
			"ListActivations",
			mock.Anything,
			edgeworkers.ListActivationsRequest{
				EdgeWorkerID: edgeWorkerID,
				Version:      version,
			},
		)
		if err != nil {
			return call.Return(nil, err)
		}

		var activations = []edgeworkers.Activation{
			{
				Note:    "note",
				Network: "STAGING",
				Version: version,
				Status:  "COMPLETE",
			},
		}

		return call.Return(
			&edgeworkers.ListActivationsResponse{
				Activations: activations,
			}, nil)
	}
)

func TestCreateEdgeWorker(t *testing.T) {
	defaultEdgercPath := "~/.edgerc"
	defaultSection := "test_section"
	localBundlePath := "testdata/res/bundle"
	localBundle := fmt.Sprintf("%s/1.24.5.tgz", localBundlePath)
	bundleBytes, err := os.ReadFile("./testdata/bundle/sampleBundle.tgz")
	if err != nil {
		require.NoError(t, err)
	}
	versionContent := bytes.NewBuffer(bundleBytes)

	tests := map[string]struct {
		givenData  TFEdgeWorkerData
		init       func(*edgeworkers.Mock, *templates.MockProcessor)
		withError  error
		withBundle bool
	}{
		"fetch edgeworker with no version": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, true, nil).Once()
				expectEdgeWorkerProcessTemplates(p, 123, "test_edgeworker", 1, 2, "", defaultEdgercPath, defaultSection, "", nil).Once()
			},
		},
		"fetch edgeworker with version": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, false, nil).Once()
				expectListActivations(e, 123, "1.24.5", nil).Once()
				expectGetEdgeWorkerVersionContent(e, 123, "1.24.5", versionContent, nil).Once()
				expectEdgeWorkerProcessTemplates(p, 123, "test_edgeworker", 1, 2, localBundle, defaultEdgercPath, defaultSection, "note", nil).Once()
			},
			withBundle: true,
		},
		"error fetching edgeworker": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, _ *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, fmt.Errorf("error")).Once()
			},
			withError: ErrFetchingEdgeWorker,
		},
		"error fetching edgeworker versions": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, _ *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, false, fmt.Errorf("error")).Once()
			},
			withError: ErrFetchingEdgeWorker,
		},
		"error fetching edgeworker version content": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, _ *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, false, nil).Once()
				expectListActivations(e, 123, "1.24.5", nil).Once()
				expectGetEdgeWorkerVersionContent(e, 123, "1.24.5", versionContent, fmt.Errorf("error")).Once()
			},
			withError: ErrFetchingEdgeWorker,
		},
		"error processing template": {
			givenData: TFEdgeWorkerData{
				EdgercPath: defaultEdgercPath,
				Section:    defaultSection,
			},
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, true, nil).Once()
				expectEdgeWorkerProcessTemplates(p, 123, "test_edgeworker", 1, 2, "", defaultEdgercPath, defaultSection, "", fmt.Errorf("error")).Once()
			},
			withError: templates.ErrSavingFiles,
		},
		"non default edgerc path and section": {
			givenData: TFEdgeWorkerData{
				EdgercPath: "/non/default/path/to/edgerc",
				Section:    "non_default_section",
			},
			init: func(e *edgeworkers.Mock, p *templates.MockProcessor) {
				expectGetEdgeWorkerID(e, 123, "test_edgeworker", 1, 2, nil).Once()
				expectListEdgeWorkerVersions(e, 123, true, nil).Once()
				expectEdgeWorkerProcessTemplates(p, 123, "test_edgeworker", 1, 2, "", "/non/default/path/to/edgerc", "non_default_section", "", nil).Once()
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(localBundlePath, 0755))

			me := new(edgeworkers.Mock)
			mp := new(templates.MockProcessor)
			test.init(me, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createEdgeWorker(ctx, 123, localBundlePath, test.givenData.EdgercPath, test.givenData.Section, me, mp)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)

			if test.withBundle {
				_, err := os.Stat(localBundle)
				require.NoError(t, err)
				require.NoError(t, os.RemoveAll(localBundle))
			}

			me.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessEdgeWorkerTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    TFEdgeWorkerData
		dir          string
		filesToCheck []string
	}{
		"edgeworker with no local bundle": {
			givenData: TFEdgeWorkerData{
				EdgeWorkerID:   123,
				Name:           "test_edgeworker",
				GroupID:        1,
				ResourceTierID: 2,
				EdgercPath:     "~/.edgerc",
				Section:        "test_section",
			},
			dir:          "edgeworker_with_no_local_bundle",
			filesToCheck: []string{"edgeworker.tf", "variables.tf", "import.sh"},
		},
		"edgeworker with local bundle": {
			givenData: TFEdgeWorkerData{
				EdgeWorkerID:   123,
				Name:           "test_edgeworker",
				GroupID:        1,
				ResourceTierID: 2,
				LocalBundle:    "testdata/bundle/sampleBundle.tgz",
				EdgercPath:     "~/.edgerc",
				Section:        "test_section",
				Note:           "note",
			},
			dir:          "edgeworker_with_local_bundle",
			filesToCheck: []string{"edgeworker.tf", "variables.tf", "import.sh"},
		},
		"edgeworker with non default edgerc path and section": {
			givenData: TFEdgeWorkerData{
				EdgeWorkerID:   123,
				Name:           "test_edgeworker",
				GroupID:        1,
				ResourceTierID: 2,
				EdgercPath:     "/non/default/path/to/edgerc",
				Section:        "non_default_section",
			},
			dir:          "edgeworker_with_non_default_edgerc_path_and_section",
			filesToCheck: []string{"variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"edgeworker.tmpl":           fmt.Sprintf("./testdata/res/%s/edgeworker.tf", test.dir),
					"edgeworker-variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"edgeworker-imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
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
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
