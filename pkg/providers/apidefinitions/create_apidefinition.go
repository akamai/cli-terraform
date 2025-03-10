// Package apidefinitions contains code for exporting API Definition configuration
package apidefinitions

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFAPIWrapperData represents the data used in API Definitions
	TFAPIWrapperData struct {
		API                  string
		ID                   int64
		Version              int64
		ContractID           string
		GroupID              int64
		ResourceName         string
		IsActiveOnStaging    bool
		IsActiveOnProduction bool
		StagingVersionKey    string
		ProductionVersionKey string
		Section              string
		Operations           string
		IsOperationsEmpty    bool
	}

	outputFormat string
)

var (
	//go:embed templates/*
	templateFiles                    embed.FS
	additionalFunctions                           = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})
	errFetchingAPI                                = errors.New("problem with fetching API")
	errFetchingResourceOperationsAPI              = errors.New("problem with fetching API Operations")
	errSavingFiles                                = errors.New("saving terraform project files")
	openAPIFormat                    outputFormat = "openapi"
	jsonFormat                       outputFormat = "json"
)

// CmdCreateAPIDefinition is an entrypoint to export-apidefinition command
func CmdCreateAPIDefinition(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := apidefinitions.Client(sess)
	clientV0 := v0.Client(sess)

	id, err := strconv.ParseInt(trimPrefixAPI(c.Args().Get(0)), 10, 64)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	section := edgegrid.GetEdgercSection(c)

	var rootPath = "./"
	if c.IsSet("tfworkpath") {
		rootPath = c.String("tfworkpath")
	}

	var versionNumber *int64
	if c.IsSet("version") {
		versionNumber = ptr.To(c.Int64("version"))
	}

	var format = openAPIFormat
	if c.IsSet("format") {
		format = outputFormat(c.String("format"))
	}

	processor, err := createTemplateProcessor(rootPath, format)
	if err != nil {
		return cli.Exit(color.RedString("Error exporting API Definition: %s", err), 1)
	}

	if _, err = createAPIDefinition(ctx, section, format, id, versionNumber, client, clientV0, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting API Definition: %s", err), 1)
	}

	if err = makeFileExecutable(filepath.Join(rootPath, "import.sh")); err != nil {
		return cli.Exit(color.RedString("Error adding execute permission to import.sh: %s", err), 1)
	}

	return nil
}

func createAPIDefinition(ctx context.Context, section string, format outputFormat, id int64, versionNumber *int64, client apidefinitions.APIDefinitions, clientV0 v0.APIDefinitions, templateProcessor templates.TemplateProcessor) (*TFAPIWrapperData, error) {
	term := terminal.Get(ctx)
	term.Spinner().Start("Fetching API Definition details for API ID: " + strconv.Itoa(int(id)))

	API, err := client.GetEndpoint(ctx, apidefinitions.GetEndpointRequest{APIEndpointID: id})
	if err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", errFetchingAPI, err)
	}

	versions, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: id,
		PageSize:      1,
		SortBy:        apidefinitions.VersionNumberSort,
		SortOrder:     apidefinitions.DescSortOrder,
	})
	if err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", errFetchingAPI, err)
	}
	latestVersionNumber := &versions.APIVersions[0].VersionNumber
	if versionNumber == nil {
		versionNumber = latestVersionNumber
	}

	var content *string
	switch format {
	case openAPIFormat:
		response, err := clientV0.ToOpenAPIFile(ctx, v0.ToOpenAPIFileRequest{ID: id, Version: *versionNumber})
		if err != nil {
			term.Spinner().Fail()
			return nil, fmt.Errorf("%w: %s", errFetchingAPI, err)
		}
		content = (*string)(response)
	case jsonFormat:
		version, err := clientV0.GetAPIVersion(ctx, v0.GetAPIVersionRequest{ID: id, Version: *versionNumber})
		if err != nil {
			term.Spinner().Fail()
			return nil, fmt.Errorf("%w: %s", errFetchingAPI, err)
		}
		term.Spinner().OK()
		content, err = serializeIndent(version)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize API : %s", err)
		}
	default:
		return nil, fmt.Errorf("value %s is invalid. Must be: '%s' or '%s'", format, openAPIFormat, jsonFormat)
	}

	var operationsContent *string

	operations, err := clientV0.GetResourceOperation(ctx, v0.GetResourceOperationRequest{APIID: id, VersionNumber: *versionNumber})

	if err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", errFetchingResourceOperationsAPI, err)
	}

	term.Spinner().OK()

	operationsContent, err = serializeResourceOperationResponseIndent(operations)

	if err != nil {
		return nil, fmt.Errorf("unable to serialize API Operations : %s", err)
	}

	isOperationsEmpty := false

	if operations.ResourceOperations == nil || operations.ResourceOperations.Len() == 0 {
		isOperationsEmpty = true
	}

	tfAPIData := populateAPIData(section, *content, id, *versionNumber, *latestVersionNumber, API, isOperationsEmpty, *operationsContent)

	term.Spinner().Start("Saving TF configurations ")

	if err = templateProcessor.ProcessTemplates(tfAPIData); err != nil {
		term.Spinner().Fail()
		return nil, fmt.Errorf("%w: %s", errSavingFiles, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for API Definitions '%d' was saved successfully\n", tfAPIData.ID)

	return &tfAPIData, nil
}

func populateAPIData(section, content string, id, versionNumber, latestVersionNumber int64, api *apidefinitions.GetEndpointResponse, isOperationsEmpty bool, operationsContent string) TFAPIWrapperData {
	return TFAPIWrapperData{
		API:                  content,
		ID:                   id,
		Version:              versionNumber,
		ContractID:           api.ContractID,
		GroupID:              api.GroupID,
		ResourceName:         sanitizeName(api.APIEndpointName),
		Section:              section,
		IsActiveOnStaging:    isActive(api.StagingVersion),
		StagingVersionKey:    versionKey("staging", api.StagingVersion, latestVersionNumber),
		IsActiveOnProduction: isActive(api.ProductionVersion),
		ProductionVersionKey: versionKey("production", api.ProductionVersion, latestVersionNumber),
		Operations:           operationsContent,
		IsOperationsEmpty:    isOperationsEmpty,
	}
}

func serializeIndent(version *v0.GetAPIVersionResponse) (*string, error) {
	jsonBody, err := json.MarshalIndent(version.RegisterAPIRequest, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func serializeResourceOperationResponseIndent(response *v0.GetResourceOperationResponse) (*string, error) {
	jsonBody, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func createTemplateProcessor(rootPath string, format outputFormat) (*templates.FSTemplateProcessor, error) {
	modulesPath := filepath.Join(rootPath, "modules")
	definitionModulePath := filepath.Join(modulesPath, "definition")
	activationModulePath := filepath.Join(modulesPath, "activation")

	paths := []string{modulesPath, definitionModulePath, activationModulePath}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	templateToFile := map[string]string{
		"apidefinitions.tmpl":       filepath.Join(rootPath, "apidefinitions.tf"),
		"variables.tmpl":            filepath.Join(rootPath, "variables.tf"),
		"import.tmpl":               filepath.Join(rootPath, "import.sh"),
		"activation-main.tmpl":      filepath.Join(activationModulePath, "main.tf"),
		"activation-variables.tmpl": filepath.Join(activationModulePath, "variables.tf"),
		"operations-main.tmpl":      filepath.Join(definitionModulePath, "operations.tf"),
		"operations-api.tmpl":       filepath.Join(definitionModulePath, "operations-api.json"),
		"definition-variables.tmpl": filepath.Join(definitionModulePath, "variables.tf"),
	}

	switch format {
	case openAPIFormat:
		templateToFile["definition-openapi-main.tmpl"] = filepath.Join(definitionModulePath, "main.tf")
		templateToFile["definition-api.tmpl"] = filepath.Join(definitionModulePath, "api.yml")
	case jsonFormat:
		templateToFile["definition-json-main.tmpl"] = filepath.Join(definitionModulePath, "main.tf")
		templateToFile["definition-api.tmpl"] = filepath.Join(definitionModulePath, "api.json")
	default:
		return nil, fmt.Errorf("value %s is invalid. Must be: '%s' or '%s'", format, openAPIFormat, jsonFormat)
	}

	for _, file := range templateToFile {
		err := tools.CheckFiles(file)
		if err != nil {
			return nil, err
		}
	}

	return &templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}, nil
}

func sanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(name, "")
}

func versionKey(network string, state apidefinitions.VersionState, latestVersion int64) string {
	if !state.IsActive() || *state.VersionNumber == latestVersion {
		return "api_latest_version"
	}

	return fmt.Sprintf("api_%s_version", network)
}

func isActive(n apidefinitions.VersionState) bool {
	return n.Status != nil && *n.Status == apidefinitions.ActivationStatusActive
}

func makeFileExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.Chmod(path, info.Mode()|0111)
}

func trimPrefixAPI(s string) string {
	return strings.TrimPrefix(s, "API_")
}
