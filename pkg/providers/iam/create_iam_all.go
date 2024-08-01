package iam

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrFetchingGroups is returned when fetching groups fails
	ErrFetchingGroups = errors.New("unable to fetch groups under this account")
	// ErrFetchingRoles is returned when fetching roles fails
	ErrFetchingRoles = errors.New("unable to fetch roles under this account")
)

// CmdCreateIAMAll is an entrypoint to create-iam all command
func CmdCreateIAMAll(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)

	groupsPath := filepath.Join(tfWorkPath, "groups.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	rolesPath := filepath.Join(tfWorkPath, "roles.tf")
	usersPath := filepath.Join(tfWorkPath, "users.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

	err := tools.CheckFiles(groupsPath, importPath, rolesPath, usersPath, variablesPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"groups.tmpl":    groupsPath,
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

	if err := createIAMAll(ctx, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting HCL for IAM: %s", err)), 1)
	}
	return nil
}

func createIAMAll(ctx context.Context, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	_, err := term.Writeln("Exporting all accessible Identity and Access Management configuration")
	if err != nil {
		return err
	}

	term.Spinner().Start("Fetching all available users")
	users, err := client.ListUsers(ctx, iam.ListUsersRequest{Actions: true})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingUsers, err)
	}
	tfUsers, err := getTFUsers(ctx, client, filterUsers(users), term)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Fetching all available groups")
	groups, err := client.ListGroups(ctx, iam.ListGroupsRequest{Actions: true})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingGroups, err)
	}
	tfGroups := make([]TFGroup, 0)
	for _, group := range groups {
		for _, innerGroup := range flattenSubgroups(&group) {
			tfGroups = append(tfGroups, getTFGroup(&innerGroup))
		}
	}
	term.Spinner().OK()

	term.Spinner().Start("Fetching all available roles")
	roles, err := client.ListRoles(ctx, iam.ListRolesRequest{
		Actions:       true,
		IgnoreContext: true,
		Users:         true,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingRoles, err)
	}
	tfRoles, err := getTFRoles(ctx, client, roles)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingRoles, err)
	}
	term.Spinner().OK()

	tfData := TFData{
		TFUsers:    tfUsers,
		TFRoles:    tfRoles,
		TFGroups:   tfGroups,
		Section:    section,
		Subcommand: "all",
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	_, err = term.Writeln("Terraform configuration was saved successfully")
	if err != nil {
		return nil
	}

	return nil
}

func flattenSubgroups(group *iam.Group) []iam.Group {
	groups := make([]iam.Group, 0)
	groups = append(groups, *group)
	for _, subGroup := range group.SubGroups {
		groups = append(groups, flattenSubgroups(&subGroup)...)
	}
	return groups
}

func filterUsers(users []iam.UserListItem) []iam.UserListItem {
	res := make([]iam.UserListItem, 0)
	for _, user := range users {
		if user.Actions != nil && user.Actions.EditProfile {
			res = append(res, user)
		}
	}
	return res
}
