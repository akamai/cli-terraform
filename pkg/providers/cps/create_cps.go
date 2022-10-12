package cps

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFCPSData represents the data used in CPS templates
	TFCPSData struct {
		Enrollment   cps.Enrollment
		EnrollmentID int
		ContractID   string
		Section      string
	}
)

//go:embed templates/*
var templateFiles embed.FS

var (
	// ErrFetchingEnrollment is returned when fetching enrollment fails
	ErrFetchingEnrollment = errors.New("unable to fetch enrollment with given id")
)

// CmdCreateCPS is an entrypoint to create-cps command
func CmdCreateCPS(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := cps.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	enrollmentPath := filepath.Join(tfWorkPath, "enrollment.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(enrollmentPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"enrollment.tmpl": enrollmentPath,
		"variables.tmpl":  variablesPath,
		"imports.tmpl":    importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
	}

	enrollmentID, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	contractID := c.Args().Get(1)
	section := edgegrid.GetEdgercSection(c)
	if err = createCPS(ctx, contractID, enrollmentID, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting enrollment HCL: %s", err)), 1)
	}
	return nil
}

func createCPS(ctx context.Context, contractID string, enrollmentID int,
	section string, client cps.CPS, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	fmt.Println("Exporting CPS configuration")

	term.Spinner().Start(fmt.Sprintf("Fetching enrollment for the given id %d", enrollmentID))
	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil || enrollment == nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingEnrollment, err)
	}

	term.Spinner().OK()

	tfData := TFCPSData{
		Enrollment:   *enrollment,
		EnrollmentID: enrollmentID,
		ContractID:   contractID,
		Section:      section,
	}

	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for enrollment '%d' was saved successfully\n", enrollmentID)

	return nil
}
