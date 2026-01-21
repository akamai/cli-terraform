package dns

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createNamedModulePath(t *testing.T) {
	tests := map[string]struct {
		tfWorkPath   string
		modName      string
		expectedPath string
	}{
		"tfWorkPath = ./": {
			tfWorkPath:   "./",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = test_path": {
			tfWorkPath:   "test_path",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "test_path/" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = ../": {
			tfWorkPath:   "../",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "../" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"blank tfWorkPath": {
			tfWorkPath:   "",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, test.expectedPath, createNamedModulePath(test.modName, test.tfWorkPath), "createNamedModulePath(%v, %v)", test.modName, test.tfWorkPath)
		})
	}
}

func Test_createDNSVarsConfig(t *testing.T) {
	tests := map[string]struct {
		edgercPath    string
		edgercSection string
		contractID    string
		expectedFile  string
	}{
		"default edgerc path and section": {
			edgercPath:    "~/.edgerc",
			edgercSection: "default",
			contractID:    "ctr_1-23456",
			expectedFile:  "testdata/dnsvars/dnsvars_default.tf",
		},
		"non default edgerc path and section": {
			edgercPath:    "/non/default/path/to/edgerc",
			edgercSection: "non_default_section",
			contractID:    "ctr_A-B-C",
			expectedFile:  "testdata/dnsvars/dnsvars_non_default.tf",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			contractID = test.contractID

			f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
			}()

			ctx := terminal.Context(context.Background(), terminal.New(f, nil, f))
			term := terminal.Get(ctx)
			err = createDNSVarsConfig(term, tmpDir, test.edgercPath, test.edgercSection)
			require.NoError(t, err)

			dnsVarsContent, err := os.ReadFile(filepath.Join(tmpDir, "dnsvars.tf"))
			require.NoError(t, err)

			expectedContent, err := os.ReadFile(test.expectedFile)
			require.NoError(t, err)
			assert.Equal(t, string(expectedContent), string(dnsVarsContent))
		})
	}
}
