package iam

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrFetchingUserByID is returned when fetching user by id fails
	ErrFetchingUserByID = errors.New("unable to fetch user by id")
	// ErrFetchingUsersWithinGroup is returned when fetching users within group fails
	ErrFetchingUsersWithinGroup = errors.New("unable to fetch users within group")
	// ErrFetchingRolesWithinGroup is returned when fetching roles within group fails
	ErrFetchingRolesWithinGroup = errors.New("unable to fetch roles within group")
)

// CmdCreateIAMGroup is an entrypoint to create-iam group command
func CmdCreateIAMGroup(c *cli.Context) error {
	ctx := c.Context
	if c.NArg() != 1 {
		return showHelpCommandWithErr(c, "Group id is required")
	}
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	tfWorkPath := "." // default is current dir
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)
	if stat, err := os.Stat(tfWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}

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
	}

	section := edgegrid.GetEdgercSection(c)
	groupID, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Wrong format of group id %v must be a number: %s", groupID, err)), 1)
	}
	if err = createIAMGroupByID(ctx, groupID, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting HCL for IAM: %s", err)), 1)
	}
	return nil
}

func createIAMGroupByID(ctx context.Context, groupID int, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	_, err := term.Writeln("Exporting Identity and Access Management group configuration with related users and groups")
	if err != nil {
		return err
	}

	term.Spinner().Start("Fetching group by id " + strconv.Itoa(groupID))
	group, err := client.GetGroup(ctx, iam.GetGroupRequest{
		GroupID: int64(groupID),
	})
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	tfGroup := getTFGroup(group)

	term.Spinner().Start("Fetching users within group with id " + strconv.Itoa(groupID))
	tfUsers, err := getUsersWithinGroup(ctx, client, groupID)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Fetching user's relative roles within group " + strconv.Itoa(groupID))
	tfRoles, err := getRolesWithinGroup(ctx, client, groupID)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	tfData := TFData{
		TFUsers: tfUsers,
		TFRoles: tfRoles,
		TFGroups: []TFGroup{
			tfGroup,
		},
		Section:    section,
		Subcommand: "group",
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

func getTFGroup(group *iam.Group) TFGroup {
	return TFGroup{
		GroupID:       int(group.GroupID),
		ParentGroupID: int(group.ParentGroupID),
		GroupName:     group.GroupName,
	}
}

func getUsersWithinGroup(ctx context.Context, client iam.IAM, groupID int) ([]*TFUser, error) {
	users, err := client.ListUsers(ctx, iam.ListUsersRequest{
		GroupID: tools.IntPtr(groupID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v with error %s", ErrFetchingUsersWithinGroup, groupID, err)
	}

	res := make([]*TFUser, 0)
	for _, v := range users {
		user, err := client.GetUser(ctx, iam.GetUserRequest{
			IdentityID:    v.IdentityID,
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %s with error %s", ErrFetchingUserByID, v.IdentityID, err)
		}
		tfUser, err := getTFUser(user)
		if err != nil {
			return nil, err
		}
		res = append(res, tfUser)
	}

	return res, nil
}

func getRolesWithinGroup(ctx context.Context, client iam.IAM, groupID int) ([]TFRole, error) {
	roles, err := client.ListRoles(ctx, iam.ListRolesRequest{
		GroupID: tools.Int64Ptr(int64(groupID)),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v with error %s", ErrFetchingRolesWithinGroup, groupID, err)
	}

	tfRoles := make([]TFRole, 0)
	for _, r := range roles {
		tfRoles = append(tfRoles, TFRole{
			RoleID:          r.RoleID,
			RoleName:        r.RoleName,
			RoleDescription: r.RoleDescription,
			GrantedRoles:    getGrantedRolesID(r.GrantedRoles),
		})
	}
	return tfRoles, nil
}
