package dns

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFiles embed.FS

type (
	// Data represents a struct passed to template
	Data struct {
		Zone           string
		ResourceFields map[string]string
		BlockName      string
	}

	// ImportData represents a struct passed to import script template
	ImportData struct {
		Zone          string
		ZoneConfigMap map[string]Types
		ResourceName  string
	}
)

var funcs = template.FuncMap{
	"namedModulePath":           createNamedModulePath,
	"checkForResource":          checkForResource,
	"createUniqueRecordsetName": createUniqueRecordsetName,
}
var tmpl = template.Must(template.New("template").Funcs(funcs).ParseFS(templateFiles, "**/*.tmpl"))

func useTemplate(data interface{}, templateName string, trimBeginning bool) string {
	buf := bytes.Buffer{}

	if err := tmpl.Lookup(templateName).Execute(&buf, data); err != nil {
		return ""
	}

	res := string(buf.Bytes())

	if trimBeginning {
		res = strings.TrimLeft(res, "\n")
	}
	return res
}

// check if resource present in state
func checkForResource(rtype string, name string) bool {

	if tfState == nil {
		if err := readTfState(); err != nil {
			// not differentiating between not exists and file error
			return false
		}
	}
	for _, r := range tfState.Resources {
		if r.Type == rtype && r.Name == name {
			return true
		}
	}

	return false
}
