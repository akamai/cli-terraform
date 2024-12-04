// Package templates allows processing multiple templates which use common data
package templates

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type (
	// TemplateProcessor allows processing multiple templates which use common data
	TemplateProcessor interface {
		// ProcessTemplates is used to parse given template/templates using the given data as input
		// If template execution fails, ProcessTemplates should return ErrTemplateExecution
		ProcessTemplates(interface{}, ...func([]string) ([]string, error)) error
		// AddTemplateTarget provides ability to specify additional template target after the processor was created
		AddTemplateTarget(string, string)
		// TemplateExists returns information if given template exists
		TemplateExists(string) bool
	}

	// FSTemplateProcessor allows working with templates stored as fs.FS
	// it contains the fs.FS as source of templates
	// as well as a map which stores template names with target files to which the result should be written
	// All templates within TemplatesFS should have .tmpl extension
	// AdditionalFuncs can be used to add custom template functions
	FSTemplateProcessor struct {
		TemplatesFS     fs.FS
		TemplateTargets map[string]string
		AdditionalFuncs template.FuncMap
	}
)

var (
	errEmptyProcessingOutput = errors.New("empty processing output")
	errTemplateCreation      = errors.New("creating template")
	// ErrTemplateExecution is returned when template.Execute method fails
	ErrTemplateExecution = errors.New("executing template")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving processed terraform file")
	// ErrNoFile is returned when there is no template file
	ErrNoFile = errors.New("no template file")
)

// ProcessTemplates parses templates located in fs.FS and executes them using the provided data
// result of each template execution is persisted in location provided in FSTemplateProcessor.TemplateTargets
func (t FSTemplateProcessor) ProcessTemplates(data interface{}, filterFuncs ...func([]string) ([]string, error)) error {
	tmpl, err := getTemplate(t.TemplatesFS, t.AdditionalFuncs, filterFuncs)
	if err != nil {
		return fmt.Errorf("%w: %s", errTemplateCreation, err)
	}

	for templateName, targetPath := range t.TemplateTargets {
		if err = processTemplateToFile(tmpl, templateName, targetPath, data); err != nil && !errors.Is(err, errEmptyProcessingOutput) {
			return err
		}
	}
	return nil
}

func getTemplate(templatesFS fs.FS, additionalFuncs template.FuncMap, filterFuncs []func([]string) ([]string, error)) (*template.Template, error) {
	funcs := template.FuncMap{
		"escape":        tools.EscapeQuotedStringLit,
		"formatIntList": formatIntList,
		"toJSON":        tools.ToJSON,
		"escapeName":    tools.EscapeName,
		"toList":        tools.ToList,
	}
	files, err := findTemplateFiles(templatesFS)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", "error filtering template files", err)
	}

	for _, f := range filterFuncs {
		files, err = f(files)
		if err != nil {
			return nil, fmt.Errorf("%w: error filtering template files", err)
		}
	}

	tmpl := template.Must(template.New("templates").Funcs(funcs).Funcs(additionalFuncs).
		ParseFS(templatesFS, files...))
	return tmpl, nil
}

// AddTemplateTarget provides ability to specify additional template target after the processor was created
func (t FSTemplateProcessor) AddTemplateTarget(templateName, targetPath string) {
	t.TemplateTargets[templateName] = targetPath
}

// TemplateExists returns information if given template exists
func (t FSTemplateProcessor) TemplateExists(fileName string) bool {
	return templateExistInFS(fileName, t.TemplatesFS)
}

func templateExistInFS(fileName string, fs fs.FS) bool {
	files, err := findTemplateFiles(fs)
	if err != nil {
		return false
	}

	for _, file := range files {
		if _, name := path.Split(file); name == fileName {
			return true
		}
	}
	return false
}

func processTemplateToFile(tmpl *template.Template, templateName, targetPath string, data interface{}) error {
	buf := bytes.Buffer{}

	if err := tmpl.Lookup(templateName).Execute(&buf, data); err != nil {
		return fmt.Errorf("%w: %s: %s", ErrTemplateExecution, templateName, err)
	}
	out := buf.Bytes()
	if len(bytes.TrimSpace(out)) == 0 {
		return errEmptyProcessingOutput
	}
	if filepath.Ext(targetPath) == ".tf" {
		out = hclwrite.Format(out)
	}
	if err := os.WriteFile(targetPath, out, 0644); err != nil {
		return fmt.Errorf("%w: '%s': %s", ErrSavingFiles, targetPath, err)
	}
	return nil
}

func formatIntList(items []int) string {
	if len(items) == 0 {
		return "[]"
	}
	var list []string
	for _, v := range items {
		list = append(list, strconv.Itoa(v))
	}
	output := strings.Join(list, ", ")
	return "[" + output + "]"
}

func findTemplateFiles(dirFS fs.FS) ([]string, error) {
	var files []string

	err := fs.WalkDir(dirFS, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path.Ext(filePath) == ".tmpl" {
			files = append(files, filePath)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return files, nil
}
