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
	// ErrFetchingUser is returned when fetching user fails
	ErrFetchingUser = errors.New("unable to fetch user by email")
	// ErrUserNotExist is returned when user does not exist
	ErrUserNotExist = errors.New("user does not exist with given email")
	// ErrMarshalUserAuthGrants is returned when marshal user auth grants failed
	ErrMarshalUserAuthGrants = errors.New("unable to marshal AuthGrants ")
)

// CmdCreateIAMUser is an entrypoint to create-iam user command
func CmdCreateIAMUser(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := iam.Client(sess)
	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)

	groupPath := filepath.Join(tfWorkPath, "groups.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	rolesPath := filepath.Join(tfWorkPath, "roles.tf")
	userPath := filepath.Join(tfWorkPath, "user.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

	err := tools.CheckFiles(userPath, groupPath, rolesPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"groups.tmpl":    groupPath,
		"imports.tmpl":   importPath,
		"roles.tmpl":     rolesPath,
		"users.tmpl":     userPath,
		"variables.tmpl": variablesPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}

	section := edgegrid.GetEdgercSection(c)
	email := c.Args().First()
	if err = createIAMUserByEmail(ctx, email, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting HCL for IAM: %s", err)), 1)
	}
	return nil
}

func createIAMUserByEmail(ctx context.Context, userEmail, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	_, err := term.Writeln("Exporting Identity and Access Management user configuration with relevant roles and groups")
	if err != nil {
		return err
	}
	term.Spinner().Start("Fetching user by email " + userEmail)

	user, err := getUserByEmail(ctx, client, userEmail)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	tfUserData, err := getTFUser(user)
	if err != nil {
		term.Spinner().Fail()
		return err
	}

	authGrantsList := user.AuthGrants

	tfData := TFData{
		TFUsers: []*TFUser{
			tfUserData,
		},
		Section:    section,
		Subcommand: "user",
	}

	if len(authGrantsList) > 0 {
		term.Spinner().Start("Fetching roles for user " + userEmail)
		tfData.TFRoles, err = getTFUserRoles(ctx, client, authGrantsList)
		if err != nil {
			term.Spinner().Fail()
			return err
		}
		term.Spinner().OK()

		term.Spinner().Start("Fetching groups for user " + userEmail)
		tfData.TFGroups, err = getTFUserGroups(ctx, client, authGrantsList)
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
	_, err = term.Writeln(fmt.Sprintf("Terraform configuration for user with email '%s' was saved successfully", tfUserData.Email))
	if err != nil {
		return nil
	}

	return nil
}

func getUserByEmail(ctx context.Context, client iam.IAM, email string) (*iam.User, error) {
	users, err := client.ListUsers(ctx, iam.ListUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFetchingUsers, err)
	}
	for _, v := range users {
		if v.Email == email {
			user, err := client.GetUser(ctx, iam.GetUserRequest{
				IdentityID:    v.IdentityID,
				Actions:       true,
				AuthGrants:    true,
				Notifications: true,
			})
			if err != nil {
				return nil, fmt.Errorf("%w: %s with error %s", ErrFetchingUser, email, err)
			}
			return user, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrUserNotExist, email)
}

func getTFUserRoles(ctx context.Context, client iam.IAM, authGrantsList []iam.AuthGrant) ([]TFRole, error) {
	roles := make([]TFRole, 0)
	for i := range authGrantsList {
		roleID := authGrantsList[i].RoleID
		if roleID != nil {
			role, err := client.GetRole(ctx, iam.GetRoleRequest{
				ID:           int64(*roleID),
				GrantedRoles: true,
			})
			if err != nil {
				return nil, fmt.Errorf("could not get role with roleID '%v': %w", roleID, err)
			}
			roles = append(roles, TFRole{
				RoleID:          role.RoleID,
				RoleName:        role.RoleName,
				RoleDescription: role.RoleDescription,
				GrantedRoles:    getGrantedRolesID(role.GrantedRoles),
			})
		}
	}
	return roles, nil
}

func getTFUserGroups(ctx context.Context, client iam.IAM, authGrantsList []iam.AuthGrant) ([]TFGroup, error) {
	allGroups := make([]TFGroup, 0)
	for i := range authGrantsList {
		groupID := authGrantsList[i].GroupID
		if groupID > 0 {
			groups, err := getGroupsInSubtree(ctx, client, groupID)
			if err != nil {
				return nil, err
			}
			allGroups = append(allGroups, groups...)
		}
	}
	return allGroups, nil
}

func getGroupsInSubtree(ctx context.Context, client iam.IAM, groupID int64) ([]TFGroup, error) {
	groups := make([]TFGroup, 0)
	group, err := client.GetGroup(ctx, iam.GetGroupRequest{
		GroupID: groupID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get group with groupID '%v': %w", groupID, err)
	}
	groups = append(groups, TFGroup{
		GroupID:       int(group.GroupID),
		ParentGroupID: int(group.ParentGroupID),
		GroupName:     group.GroupName,
	})
	for _, subGroup := range group.SubGroups {
		subGroups, err := getGroupsInSubtree(ctx, client, subGroup.GroupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, subGroups...)
	}
	return groups, nil
}
