package iam

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// CmdCreateIAMGroup is an entrypoint to create-iam group command
func CmdCreateIAMGroup(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)

	groupOnly := c.Bool("only")

	groupPath := filepath.Join(tfWorkPath, "group.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	rolesPath := filepath.Join(tfWorkPath, "roles.tf")
	usersPath := filepath.Join(tfWorkPath, "users.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

	err := tools.CheckFiles(groupPath, usersPath, rolesPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"groups.tmpl":    groupPath,
		"imports.tmpl":   importPath,
		"roles.tmpl":     rolesPath,
		"users.tmpl":     usersPath,
		"variables.tmpl": variablesPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}

	section := edgegrid.GetEdgercSection(c)
	groupID, err := strconv.ParseInt(c.Args().First(), 10, 64)
	if err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Wrong format of group id %v must be a number: %s", groupID, err)), 1)
	}
	if err = createIAMGroupByID(ctx, groupID, section, client, processor, groupOnly); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting HCL for IAM: %s", err)), 1)
	}
	return nil
}

func createIAMGroupByID(ctx context.Context, groupID int64, section string, client iam.IAM, templateProcessor templates.TemplateProcessor, groupOnly bool) error {
	term := terminal.Get(ctx)

	message := "Exporting Identity and Access Management group configuration"
	if !groupOnly {
		message += " with related roles and users"
	}

	if _, err := term.Writeln(message); err != nil {
		return err
	}

	term.Spinner().Start("Fetching group by id " + strconv.FormatInt(groupID, 10))
	group, err := client.GetGroup(ctx, iam.GetGroupRequest{
		GroupID: groupID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("could not get group with ID '%v': %w", groupID, err)
	}
	term.Spinner().OK()

	tfGroup := getTFGroup(group)

	tfData := TFData{
		TFGroups: []TFGroup{
			tfGroup,
		},
		Section:    section,
		Subcommand: "group",
	}

	if !groupOnly {
		term.Spinner().Start("Fetching users within group with id " + strconv.FormatInt(groupID, 10))
		tfData.TFUsers, err = getUsersWithinGroup(ctx, client, groupID, term)
		if err != nil {
			term.Spinner().Fail()
			return err
		}
		term.Spinner().OK()

		term.Spinner().Start("Fetching user's relative roles within group " + strconv.FormatInt(groupID, 10))
		tfData.TFRoles, err = getRolesWithinGroup(ctx, client, groupID)
		if err != nil {
			term.Spinner().Fail()
			return err
		}
		term.Spinner().OK()
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	_, err = term.Writeln(fmt.Sprintf("Terraform configuration for group with id '%v' was saved successfully", groupID))
	if err != nil {
		return nil
	}

	return nil
}

func getUsersWithinGroup(ctx context.Context, client iam.IAM, groupID int64, term terminal.Terminal) ([]*TFUser, error) {
	users, err := client.ListUsers(ctx, iam.ListUsersRequest{
		Actions: true,
		GroupID: tools.Int64Ptr(groupID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v with error %s", ErrFetchingUsersWithinGroup, groupID, err)
	}

	return getTFUsers(ctx, client, filterUsers(users), term)
}

func getRolesWithinGroup(ctx context.Context, client iam.IAM, groupID int64) ([]TFRole, error) {
	roles, err := client.ListRoles(ctx, iam.ListRolesRequest{
		GroupID: tools.Int64Ptr(groupID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v with error %s", ErrFetchingRolesWithinGroup, groupID, err)
	}

	return getTFRoles(ctx, client, roles)
}
