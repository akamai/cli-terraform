package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	// ErrFetchingUsers is returned when fetching users fails
	ErrFetchingUsers = errors.New("unable to fetch users under this account")
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
	if c.NArg() != 1 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.Exit(color.RedString("Error displaying help command"), 1)
		}
		return cli.Exit(color.RedString("User's email is required"), 1)
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
		"user.tmpl":      userPath,
		"variables.tmpl": variablesPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
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
	fmt.Println("Exporting Identity and Access Management user configuration with related role and groups")
	term.Spinner().Start("Fetching user by email " + userEmail)

	user, err := getUserByEmail(ctx, client, userEmail)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	tfData, err := getTFUserData(user, section)
	if err != nil {
		term.Spinner().Fail()
		return err
	}

	authGrantsList := user.AuthGrants

	if len(authGrantsList) > 0 {
		tfData.Roles, err = getTFUserRoles(ctx, client, authGrantsList)
		if err != nil {
			term.Spinner().Fail()
			return err
		}

		tfData.Groups, err = getTFUserGroups(ctx, client, authGrantsList)
		if err != nil {
			term.Spinner().Fail()
			return err
		}
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for user with email '%s' was saved successfully\n", tfData.Email)

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
				return nil, err
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
	groups := make([]TFGroup, 0)
	for i := range authGrantsList {
		groupID := authGrantsList[i].GroupID
		if groupID > 0 {
			group, err := client.GetGroup(ctx, iam.GetGroupRequest{
				GroupID: int64(groupID),
			})
			if err != nil {
				return nil, err
			}
			groups = append(groups, TFGroup{
				GroupID:       int(group.GroupID),
				ParentGroupID: int(group.ParentGroupID),
				GroupName:     group.GroupName,
			})
		}
	}
	return groups, nil
}

func getTFUserData(user *iam.User, section string) (*TFUserData, error) {
	authGrants, err := getUserAuthGrants(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %s for user with email %s", ErrMarshalUserAuthGrants, err, user.Email)
	}

	return &TFUserData{
		TFUserBasicInfo: TFUserBasicInfo{
			ID:                user.IdentityID,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			Country:           user.Country,
			Phone:             user.Phone,
			TFAEnabled:        user.TFAEnabled,
			ContactType:       user.ContactType,
			JobTitle:          user.JobTitle,
			TimeZone:          user.TimeZone,
			SecondaryEmail:    user.SecondaryEmail,
			MobilePhone:       user.MobilePhone,
			Address:           user.Address,
			City:              user.City,
			State:             user.State,
			ZipCode:           user.ZipCode,
			PreferredLanguage: user.PreferredLanguage,
			SessionTimeOut:    user.SessionTimeOut,
		},
		IsLocked:   user.IsLocked,
		AuthGrants: authGrants,
		Section:    section,
	}, nil
}

func getUserAuthGrants(user *iam.User) (string, error) {
	var authGrantsJSON []byte
	var err error
	if len(user.AuthGrants) > 0 {
		authGrantsJSON, err = json.Marshal(user.AuthGrants)
		if err != nil {
			return "", err
		}
	}
	return string(authGrantsJSON), nil
}

func getGrantedRolesID(grantedRoles []iam.RoleGrantedRole) []int {
	rolesIDs := make([]int, 0, len(grantedRoles))
	for _, v := range grantedRoles {
		rolesIDs = append(rolesIDs, int(v.RoleID))
	}
	return rolesIDs
}
