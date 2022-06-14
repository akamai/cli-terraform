package templates

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
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
		ProcessTemplates(interface{}) error
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
	// ErrTemplateExecution is returned when template.Execute method fails
	ErrTemplateExecution = errors.New("executing template")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving processed terraform file")
)

// ProcessTemplates parses templates located in fs.FS and executes them using the provided data
// result of each template execution is persisted in location provided in FSTemplateProcessor.TemplateTargets
func (t FSTemplateProcessor) ProcessTemplates(data interface{}) error {
	funcs := template.FuncMap{
		"escape":     tools.EscapeQuotedStringLit,
		"formatIntList": formatIntList,
		"toJSON":     tools.ToJSON,
		"escapeName": tools.EscapeName,
		"toList":     tools.ToList,
	}
	tmpl := template.Must(template.New("templates").Funcs(funcs).Funcs(t.AdditionalFuncs).ParseFS(t.TemplatesFS, "**/*.tmpl"))

	for templateName, targetPath := range t.TemplateTargets {
		buf := bytes.Buffer{}

		if err := tmpl.Lookup(templateName).Execute(&buf, data); err != nil {
			return fmt.Errorf("%w: %s: %s", ErrTemplateExecution, templateName, err)
		}
		out := buf.Bytes()
		if len(bytes.TrimSpace(out)) == 0 {
			continue
		}
		if filepath.Ext(targetPath) == ".tf" {
			out = hclwrite.Format(out)
		}
		if err := os.WriteFile(targetPath, out, 0644); err != nil {
			return fmt.Errorf("%w: '%s': %s", ErrSavingFiles, targetPath, err)
		}
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