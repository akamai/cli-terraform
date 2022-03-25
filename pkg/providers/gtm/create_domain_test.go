package gtm

import (
	"fmt"
	"io/ioutil"
	"os"
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
			givenData:    map[string]interface{}{},
			dir:          "only_variables",
			filesToCheck: []string{"gtmvars.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/gtmvars.tf", test.dir),
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
