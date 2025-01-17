package iam

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// CmdCreateIAMAllowlist is an entrypoint to create-iam allowlist command
func CmdCreateIAMAllowlist(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)

	importPath := filepath.Join(tfWorkPath, "import.sh")
	allowlistPath := filepath.Join(tfWorkPath, "allowlist.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

	err := tools.CheckFiles(allowlistPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"imports.tmpl":   importPath,
		"allowlist.tmpl": allowlistPath,
		"variables.tmpl": variablesPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}

	section := edgegrid.GetEdgercSection(c)
	if err = createIAMAllowlist(ctx, section, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting HCL for IAM: %s", err), 1)
	}
	return nil
}

func createIAMAllowlist(ctx context.Context, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	_, err := term.Writeln("Exporting Identity and Access Management allowlist configuration")
	if err != nil {
		return err
	}
	term.Spinner().Start("Fetching IP allowlist and CIDR blocks")

	status, err := client.GetIPAllowlistStatus(ctx)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingIPAllowlistStatus, err)
	}

	tfCIDRBlocks, err := getTFCIDRBlocks(ctx, client)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingCIDRBlocks, err)
	}

	tfData := TFData{
		TFAllowlist: TFAllowlist{
			Enabled:    status.Enabled,
			CIDRBlocks: tfCIDRBlocks,
		},
		Section:    section,
		Subcommand: "allowlist",
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	_, err = term.Writeln("Terraform configuration for allowlist and CIDR blocks was saved successfully")
	if err != nil {
		return nil
	}

	return nil
}
