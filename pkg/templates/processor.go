package templates

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

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
		"escape": escapeQuotedStringLit,
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

// escapeQuotedStringLit returns escaped terraform string literal
// https://www.terraform.io/docs/language/expressions/strings.html#escape-sequences
// This function is based on https://github.com/hashicorp/hcl/blob/c7ee8b78101c33b4dfed2641d78cf5e9651eabb8/hclwrite/generate.go#L207-L246
func escapeQuotedStringLit(s string) string {
	if len(s) == 0 {
		return ""
	}
	buf := strings.Builder{}
	for i, r := range s {
		switch r {
		case '\n':
			buf.Write([]byte{'\\', 'n'})
		case '\r':
			buf.Write([]byte{'\\', 'r'})
		case '\t':
			buf.Write([]byte{'\\', 't'})
		case '"':
			buf.Write([]byte{'\\', '"'})
		case '\\':
			buf.Write([]byte{'\\', '\\'})
		case '$', '%':
			buf.WriteRune(r)
			remain := s[i+1:]
			if len(remain) > 0 && remain[0] == '{' {
				// Double up our template introducer symbol to escape it.
				buf.WriteRune(r)
			}
		default:
			if !unicode.IsPrint(r) {
				var fmted string
				if r < 65536 {
					fmted = fmt.Sprintf("\\u%04x", r)
				} else {
					fmted = fmt.Sprintf("\\U%08x", r)
				}
				buf.WriteString(fmted)
			} else {
				buf.WriteRune(r)
			}
		}
	}
	return buf.String()
}
