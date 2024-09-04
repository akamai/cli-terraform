package iam

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrFetchingRole is returned when fetching role fails
	ErrFetchingRole = errors.New("unable to fetch role by role_id")
)

// CmdCreateIAMRole is an entrypoint to create-iam role command
func CmdCreateIAMRole(c *cli.Context) error {
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
	rolesPath := filepath.Join(tfWorkPath, "role.tf")
	userPath := filepath.Join(tfWorkPath, "users.tf")
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
	roleID, err := strconv.ParseInt(c.Args().First(), 10, 64)
	if err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Wrong format of role id %v must be a number: %s", roleID, err)), 1)
	}

	if err = createIAMRoleByID(ctx, roleID, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting HCL for IAM: %s", err)), 1)
	}
	return nil
}

func createIAMRoleByID(ctx context.Context, roleID int64, section string, client iam.IAM, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	_, err := term.Writeln("Exporting Identity and Access Management role configuration with related users and groups")
	if err != nil {
		return err
	}
	term.Spinner().Start(fmt.Sprintf("Fetching role by role_id %d", roleID))

	role, err := client.GetRole(ctx, iam.GetRoleRequest{
		ID:           roleID,
		GrantedRoles: true,
		Users:        true,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: could not fetch role with roleID '%v': %s", ErrFetchingRole, roleID, err)
	}
	term.Spinner().OK()

	tfRole := TFRole{
		RoleID:          role.RoleID,
		RoleName:        role.RoleName,
		RoleDescription: role.RoleDescription,
		GrantedRoles:    getGrantedRolesID(role.GrantedRoles),
	}

	term.Spinner().Start(fmt.Sprintf("Fetching users with the given role %d", roleID))
	users, err := getUsersByRole(ctx, term, role.Users, client)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	tfUsers := make([]*TFUser, 0)
	tfGroups := make([]TFGroup, 0)

	term.Spinner().Start(fmt.Sprintf("Fetching groups for users related within role %d", roleID))
	for _, user := range users {
		userData, err := getTFUser(user)
		if err != nil {
			term.Spinner().Fail()
			return err
		}

		tfUsers = append(tfUsers, userData)
		authGrantsList := user.AuthGrants

		if len(authGrantsList) > 0 {
			groupsData, err := getTFUserGroups(ctx, client, authGrantsList)
			if err != nil {
				term.Spinner().Fail()
				return err
			}

			tfGroups = appendUniqueGroups(tfGroups, groupsData)
		}
	}
	term.Spinner().OK()

	tfData := TFData{
		TFUsers:  tfUsers,
		TFGroups: tfGroups,
		TFRoles: []TFRole{
			tfRole,
		},

		Section:    section,
		Subcommand: "role",
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	_, err = term.Writeln(fmt.Sprintf("Terraform configuration for role with id '%d' was saved successfully", tfRole.RoleID))
	if err != nil {
		return nil
	}

	return nil
}

func appendUniqueGroups(tfGroups []TFGroup, groupsData []TFGroup) []TFGroup {
	for _, groupData := range groupsData {
		neverSeenGroup := true
		for _, tfGroup := range tfGroups {
			if groupData.GroupID == tfGroup.GroupID {
				neverSeenGroup = false
				break
			}
		}
		if neverSeenGroup {
			tfGroups = append(tfGroups, groupData)
		}
	}

	return tfGroups
}

func getUsersByRole(ctx context.Context, term terminal.Terminal, roleUsers []iam.RoleUser, client iam.IAM) ([]*iam.User, error) {
	users := make([]*iam.User, 0)

	for _, roleUser := range roleUsers {
		user, err := client.GetUser(ctx, iam.GetUserRequest{
			IdentityID:    roleUser.UIIdentityID,
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		})
		if err != nil {
			_, err := term.Writeln(fmt.Sprintf("[WARN] Unable to fetch user of ID '%s' - skipping:\n%s", roleUser.UIIdentityID, err))
			if err != nil {
				return nil, err
			}
			continue
		}

		users = append(users, user)
	}

	return users, nil
}
