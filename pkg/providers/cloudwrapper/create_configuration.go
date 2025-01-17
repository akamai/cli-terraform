// Package cloudwrapper contains code for exporting CloudWrapper configuration
package cloudwrapper

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudwrapper"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (
	// TFCloudWrapperData represents the data used in CloudWrapper
	TFCloudWrapperData struct {
		Configuration TFCWConfiguration
		Section       string
	}

	// TFCWConfiguration represents the data used for export CloudWrapper configuration
	TFCWConfiguration struct {
		ConfigurationResourceName string
		ID                        int64
		ContractID                string
		Name                      string
		PropertyIDs               []string
		Comments                  string
		Status                    string
		NotificationEmails        []string
		RetainIdleObjects         bool
		CapacityAlertsThreshold   *int
		Locations                 []Location
		IsActive                  bool
	}

	// Location represents CloudWrapper location
	Location struct {
		TrafficTypeID int
		Comments      string
		Capacity      Capacity
	}

	// Capacity represents capacity of location
	Capacity struct {
		Value int64
		Unit  string
	}
)

//go:embed templates/*
var templateFiles embed.FS

var additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})

var (
	// ErrFetchingConfiguration is returned when problem occurred during fetching configuration
	ErrFetchingConfiguration = errors.New("problem with fetching configuration")
	// ErrContainMultiCDNSettings is returned when configuration contains Multi CDN Setting
	ErrContainMultiCDNSettings = errors.New("configuration contains Multi CDN Settings")
	// ErrExportingCloudWrapper is returned when there is an issue related to export of CloudWrapper
	ErrExportingCloudWrapper = errors.New("error exporting CloudWrapper")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving terraform project files")
)

// CmdCreateCloudWrapper is an entrypoint to export-cloudwrapper command
func CmdCreateCloudWrapper(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := cloudwrapper.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	cloudwrapperPath := filepath.Join(tfWorkPath, "cloudwrapper.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(cloudwrapperPath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"cloudwrapper.tmpl": cloudwrapperPath,
		"variables.tmpl":    variablesPath,
		"imports.tmpl":      importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}
	configID, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	section := edgegrid.GetEdgercSection(c)
	if err = createCloudWrapper(ctx, configID, section, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting cloudwraper: %s", err), 1)
	}
	return nil
}

func createCloudWrapper(ctx context.Context, configID int64, section string, client cloudwrapper.CloudWrapper, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	term.Spinner().Start("Fetching configuration " + strconv.Itoa(int(configID)))
	configuration, err := client.GetConfiguration(ctx, cloudwrapper.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingConfiguration, err)
	}
	if configuration.MultiCDNSettings != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%s: %w", ErrExportingCloudWrapper, ErrContainMultiCDNSettings)
	}
	tfCloudWrapperData := populateCloudWrapperData(configID, section, configuration)

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfCloudWrapperData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for CloudWrapper configuration '%d' was saved successfully\n", tfCloudWrapperData.Configuration.ID)

	return nil
}

func populateCloudWrapperData(configID int64, section string, configuration *cloudwrapper.Configuration) TFCloudWrapperData {
	tfCloudWrapperData := TFCloudWrapperData{
		Configuration: TFCWConfiguration{
			ID:                        configID,
			ContractID:                configuration.ContractID,
			PropertyIDs:               configuration.PropertyIDs,
			Name:                      configuration.ConfigName,
			Comments:                  configuration.Comments,
			NotificationEmails:        configuration.NotificationEmails,
			ConfigurationResourceName: strings.ReplaceAll(configuration.ConfigName, "-", "_"),
			RetainIdleObjects:         configuration.RetainIdleObjects,
			CapacityAlertsThreshold:   configuration.CapacityAlertsThreshold,
		},
		Section: section,
	}
	tfCloudWrapperData.Configuration.Status = string(configuration.Status)
	for _, configLocation := range configuration.Locations {
		tfCloudWrapperData.Configuration.Locations = append(tfCloudWrapperData.Configuration.Locations, Location{
			TrafficTypeID: configLocation.TrafficTypeID,
			Comments:      configLocation.Comments,
			Capacity: Capacity{
				Unit:  string(configLocation.Capacity.Unit),
				Value: configLocation.Capacity.Value,
			},
		})
	}
	if tfCloudWrapperData.Configuration.Status == string(cloudwrapper.StatusActive) {
		tfCloudWrapperData.Configuration.IsActive = true
	}
	return tfCloudWrapperData
}
