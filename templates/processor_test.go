package templates

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestData struct {
	A, B string
}

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

func TestProcessTemplates(t *testing.T) {
	tests := map[string]struct {
		templateDir     string
		templateTargets map[string]string
		data            TestData
		withError       error
		expected        map[string]string
	}{
		"process simple templates": {
			templateDir: "./testdata",
			templateTargets: map[string]string{
				"1.tmpl": "./testdata/res/1.txt",
				"2.tmpl": "./testdata/res/2.txt",
			},
			data: TestData{
				A: "Hello",
				B: "World",
			},
			expected: map[string]string{
				"./testdata/res/1.txt": "Hello",
				"./testdata/res/2.txt": "World",
			},
		},
		"do not save empty file": {
			templateDir: "./testdata",
			templateTargets: map[string]string{
				"1.tmpl": "./testdata/res/1.txt",
				"2.tmpl": "./testdata/res/not_existing.txt",
			},
			data: TestData{
				A: "Hello",
			},
			expected: map[string]string{
				"./testdata/res/1.txt":            "Hello",
				"./testdata/res/not_existing.txt": "",
			},
		},
		"nested template": {
			templateDir: "./testdata",
			templateTargets: map[string]string{
				"with_nesting.tmpl": "./testdata/res/res.txt",
			},
			data: TestData{
				A: "Hello",
			},
			expected: map[string]string{
				"./testdata/res/res.txt": "This nests template 1: Hello",
			},
		},
		"error executing template": {
			templateDir: "./testdata",
			templateTargets: map[string]string{
				"invalid_field.tmpl": "./testdata/res/res.txt",
			},
			data: TestData{
				A: "Hello",
			},
			withError: ErrTemplateExecution,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			templateFS := os.DirFS(test.templateDir)
			processor := FSTemplateProcessor{
				TemplatesFS:     templateFS,
				TemplateTargets: test.templateTargets,
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
				res, err := ioutil.ReadFile(path)
				require.NoError(t, err)
				assert.Equal(t, val, string(res))
			}
		})
	}
}
