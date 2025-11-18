package iam

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// CmdCreateIAMClient is an entrypoint to create-iam client command
func CmdCreateIAMClient(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)

	clientPath := filepath.Join(tfWorkPath, "client.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

	err := tools.CheckFiles(clientPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}

	templateToFile := map[string]string{
		"client.tmpl":    clientPath,
		"imports.tmpl":   importPath,
		"variables.tmpl": variablesPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}

	section := edgegrid.GetEdgercSection(c)
	clientID := c.Args().First()
	if err = createIAMAPIClient(ctx, clientID, section, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting HCL for IAM: %s", err), 1)
	}
	return nil
}

// createIAMAPIClient with provided id
func createIAMAPIClient(ctx context.Context, clientID, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	message := "Exporting Identity and Access Management API client configuration"

	if _, err := term.Writeln(message); err != nil {
		return err
	}

	if clientID == "" {
		term.Spinner().Start("Fetching your API client")
	} else {
		term.Spinner().Start("Fetching client by id " + clientID)
	}
	apiClient, err := client.GetAPIClient(ctx, iam.GetAPIClientRequest{
		ClientID:    clientID,
		GroupAccess: true,
		APIAccess:   true,
		IPACL:       true,
		Credentials: true,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("could not get API client with ID '%v': %w", clientID, err)
	}
	if len(apiClient.Credentials) == 0 {
		term.Spinner().Fail()
		return fmt.Errorf("API client with ID '%v' has no credentials. It's impossible to manage API Client with no credential via Terraform", clientID)
	}
	term.Spinner().OK()

	tfAPIClient := getTFClient(apiClient)

	tfData := TFData{
		TFClient:   tfAPIClient,
		Section:    section,
		Subcommand: "client",
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	_, err = term.Writeln(fmt.Sprintf("Terraform configuration for API client with id '%v' was saved successfully",
		apiClient.ClientID))
	if err != nil {
		return nil
	}

	return nil
}
