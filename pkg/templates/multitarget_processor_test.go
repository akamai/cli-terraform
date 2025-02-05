package templates

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessMultitargetTemplates(t *testing.T) {
	tests := map[string]struct {
		templateDir string
		data        MultiTargetData
		withError   error
		expected    map[string]string
	}{
		"process simple templates": {
			templateDir: "./testdata",
			data: MultiTargetData{
				"1.tmpl": DataForTarget{
					"./testdata/res/1.txt": TestData{A: "Hello"},
					"./testdata/res/2.txt": TestData{A: "World"},
				},
				"2.tmpl": DataForTarget{
					"./testdata/res/3.txt": TestData{B: "Other hello"},
					"./testdata/res/4.txt": TestData{B: "from different planet"},
				},
			},
			expected: map[string]string{
				"./testdata/res/1.txt": "Hello",
				"./testdata/res/2.txt": "World",
				"./testdata/res/3.txt": "Other hello",
				"./testdata/res/4.txt": "from different planet",
			},
		},
		"do not save empty file": {
			templateDir: "./testdata",
			data: MultiTargetData{
				"1.tmpl": DataForTarget{
					"./testdata/res/1.txt": TestData{A: "Hello"},
				},
				"empty.tmpl": DataForTarget{
					"./testdata/res/not_existing.txt": TestData{},
				},
			},
			expected: map[string]string{
				"./testdata/res/1.txt":            "Hello",
				"./testdata/res/not_existing.txt": "",
			},
		},
		"nested template": {
			templateDir: "./testdata",
			data: MultiTargetData{
				"with_nesting.tmpl": DataForTarget{
					"./testdata/res/res.txt": TestData{A: "Hello"},
				},
			},
			expected: map[string]string{
				"./testdata/res/res.txt": "This nests template 1: Hello",
			},
		},
		"error executing template": {
			templateDir: "./testdata",
			data: MultiTargetData{
				"invalid_field.tmpl": DataForTarget{
					"./testdata/res/res.txt": TestData{A: "Hello"},
				},
			},
			withError: ErrTemplateExecution,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			templateFS := os.DirFS(test.templateDir)
			processor := FSMultiTargetProcessor{
				TemplatesFS: templateFS,
			}
			err := processor.ProcessTemplates(test.data)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			for path, val := range test.expected {
				if val == "" {
					_, err = os.Stat(path)
					assert.True(t, errors.Is(err, os.ErrNotExist), "expected no file but found '%s'", path)
					return
				}
				res, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, val, string(res))
			}
		})
	}
}
