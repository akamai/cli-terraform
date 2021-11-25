package templates

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"text/template"
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
	FSTemplateProcessor struct {
		TemplatesFS     fs.FS
		TemplateTargets map[string]string
	}
)

var (
	// ErrTemplateExecution is returned when template.Execute method fails
	ErrTemplateExecution = errors.New("executing template")
)

// ProcessTemplates parses templates located in fs.FS and executes them using the provided data
// result of each template execution is persisted in location provided in FSTemplateProcessor.TemplateTargets
func (t FSTemplateProcessor) ProcessTemplates(data interface{}) error {
	tmpl := template.Must(template.ParseFS(t.TemplatesFS, "**/*.tmpl"))
	for templateName, targetPath := range t.TemplateTargets {
		buf := bytes.Buffer{}
		if err := tmpl.Lookup(templateName).Execute(&buf, data); err != nil {
			return fmt.Errorf("%w: %s: %s", ErrTemplateExecution, templateName, err)
		}
		data := buf.Bytes()
		if len(data) == 0 || len(bytes.TrimSpace(data)) == 0 {
			continue
		}
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("creating '%s': %s", targetPath, err)
		}
	}
	return nil
}
