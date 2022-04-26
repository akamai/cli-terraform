package gtm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestProcessDomainTemplates(t *testing.T) {
	tests := map[string]struct {
		givenData    interface{}
		dir          string
		filesToCheck []string
	}{
		"variables file correct": {
			givenData:    TFDomainData{Section: "test_section"},
			dir:          "only_variables",
			filesToCheck: []string{"variables.tf"},
		},
		"domain file header correct": {
			givenData:    TFDomainData{Section: "default"},
			dir:          "domain_file",
			filesToCheck: []string{"domain.tf", "variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			outDir := filepath.Join("./testdata/res", test.dir)
			require.NoError(t, os.MkdirAll(outDir, 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"variables.tmpl": filepath.Join(outDir, "variables.tf"),
					"domain.tmpl":    filepath.Join(outDir, "domain.tf"),
				},
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))
			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(filepath.Join("./testdata", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(filepath.Join(outDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
