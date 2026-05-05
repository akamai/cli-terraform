// Package reportinggroups contains code for exporting Reporting Groups.
package reportinggroups

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (
	// TFData represents the data used in Reporting Groups templates.
	TFData struct {
		Group      TFGroup
		EdgercPath string
		Section    string
	}

	// TFGroup represents the data used for exporting a Reporting Group.
	TFGroup struct {
		ResourceName       string
		ReportingGroupID   int64
		ReportingGroupName string
		AccessContractID   string
		ContractID         string
		CPCodes            []string
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	// ErrFindingReportingGroup is returned when the named reporting group cannot be found.
	ErrFindingReportingGroup = errors.New("error finding reporting group")
	// ErrFetchingReportingGroup is returned when a reporting group could not be fetched.
	ErrFetchingReportingGroup = errors.New("error fetching reporting group")
	// ErrInvalidReportingGroup is returned when the reporting group data is not in the expected shape.
	ErrInvalidReportingGroup = errors.New("invalid reporting group")
	// ErrSavingFiles is returned when an issue with processing templates occurs.
	ErrSavingFiles = errors.New("error saving terraform project files")
)

// CmdCreateReportingGroup is an entrypoint to the export-reportinggroup command.
func CmdCreateReportingGroup(c *cli.Context) error {
	// tfWorkPath is a target directory for generated terraform resources.
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	reportingGroupPath := filepath.Join(tfWorkPath, "reportinggroups.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(reportingGroupPath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}

	params := createReportingGroupParams{
		name:          c.Args().Get(0),
		edgercPath:    edgegrid.GetEdgercPath(c),
		configSection: edgegrid.GetEdgercSection(c),
		client:        reportinggroups.Client(edgegrid.GetSession(c.Context)),
		templateProcessor: templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"reportinggroups.tmpl": reportingGroupPath,
				"variables.tmpl":       variablesPath,
				"import.tmpl":          importPath,
			},
		},
	}

	if err := createReportingGroup(c.Context, params); err != nil {
		return cli.Exit(color.RedString("Error exporting reporting group: %s", err), 1)
	}
	return nil
}

type createReportingGroupParams struct {
	name              string
	edgercPath        string
	configSection     string
	client            reportinggroups.ReportingGroups
	templateProcessor templates.TemplateProcessor
}

func createReportingGroup(ctx context.Context, params createReportingGroupParams) (e error) {
	term := terminal.Get(ctx)
	term.Spinner().Start("Fetching reporting group '" + params.name + "'")
	defer func() {
		if e != nil {
			term.Spinner().Fail()
		}
	}()

	group, err := findReportingGroup(ctx, params.client, params.name)
	if err != nil {
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Extracting data")
	tfData, err := buildTFData(params.edgercPath, params.configSection, group)
	if err != nil {
		return err
	}
	term.Spinner().OK()
	term.Spinner().Start("Saving TF configurations")
	if err = params.templateProcessor.ProcessTemplates(tfData); err != nil {
		return fmt.Errorf("%w: %w", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for reporting group '%s' was saved successfully\n", params.name)
	return nil
}

func findReportingGroup(ctx context.Context, client reportinggroups.ReportingGroups, name string) (*reportinggroups.ReportingGroup, error) {
	list, err := client.ListReportingGroups(ctx, reportinggroups.ListReportingGroupsRequest{
		ReportingGroupName: name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFetchingReportingGroup, err)
	}

	var matches []reportinggroups.ReportingGroup
	for _, g := range list.Groups {
		if g.ReportingGroupName == name {
			matches = append(matches, g)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("%w: no reporting group found with name %q", ErrFindingReportingGroup, name)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("%w: multiple reporting groups found with name %q", ErrFindingReportingGroup, name)
	}
}

func buildTFData(edgercPath, configSection string, group *reportinggroups.ReportingGroup) (TFData, error) {
	if len(group.Contracts) != 1 {
		return TFData{}, fmt.Errorf("%w: expected exactly one contract, got %d",
			ErrInvalidReportingGroup, len(group.Contracts))
	}

	contract := group.Contracts[0]
	cpCodes := make([]string, 0, len(contract.CPCodes))
	for _, cp := range contract.CPCodes {
		cpCodes = append(cpCodes, strconv.FormatInt(cp.CPCodeID, 10))
	}

	return TFData{
		EdgercPath: edgercPath,
		Section:    configSection,
		Group: TFGroup{
			ResourceName:       tools.SanitizeResourceName(group.ReportingGroupName),
			ReportingGroupID:   group.ReportingGroupID,
			ReportingGroupName: group.ReportingGroupName,
			AccessContractID:   group.AccessGroup.ContractID,
			ContractID:         contract.ContractID,
			CPCodes:            cpCodes,
		},
	}, nil
}
