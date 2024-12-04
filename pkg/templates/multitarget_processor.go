package templates

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"text/template"
)

type (
	// DataForTarget holds data about (template) output file to input data that this template needs
	DataForTarget map[string]any

	// MultiTargetData holds relationship between template to be executed and map of entries for which this template is to be executed
	MultiTargetData map[string]DataForTarget

	// MultiTargetProcessor allows processing multiple templates that each template and each target uses different data
	MultiTargetProcessor interface {
		// ProcessTemplates processes template
		ProcessTemplates(MultiTargetData, ...func([]string) ([]string, error)) error
	}

	// FSMultiTargetProcessor allows to work on templates that are stored as fs.FS as implementation of MultiTargetProcessor interface
	FSMultiTargetProcessor struct {
		TemplatesFS     fs.FS
		AdditionalFuncs template.FuncMap
	}
)

// ProcessTemplates searches for templates inside fs.FS. Later it executes for each template, for each target (and its data) that is provided in input.
func (t FSMultiTargetProcessor) ProcessTemplates(data MultiTargetData, filterFuncs ...func([]string) ([]string, error)) error {
	tmpl, err := getTemplate(t.TemplatesFS, t.AdditionalFuncs, filterFuncs)
	if err != nil {
		return fmt.Errorf("%w: %s", errTemplateCreation, err)
	}

	for templateName, processingData := range data {
		for templateTarget, templateData := range processingData {
			if err = processTemplateToFile(tmpl, templateName, templateTarget, templateData); err != nil && !errors.Is(err, errEmptyProcessingOutput) {
				return err
			}
		}
	}
	return nil
}

// Join merges two DataForTargets
func (d DataForTarget) Join(toAdd DataForTarget) DataForTarget {
	maps.Copy(d, toAdd)
	return d
}

// AddData adds data for processing for a template. If the template already exists, it merges data.
func (m MultiTargetData) AddData(templateTarget string, toAdd DataForTarget) MultiTargetData {
	if added, exist := m[templateTarget]; exist {
		m[templateTarget] = added.Join(toAdd)
	} else {
		m[templateTarget] = toAdd
	}
	return m
}
