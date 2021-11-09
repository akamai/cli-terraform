package templates

import (
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
		file, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		if err := tmpl.Lookup(templateName).Execute(file, data); err != nil {
			return fmt.Errorf("%w: %s: %s", ErrTemplateExecution, templateName, err)
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}
