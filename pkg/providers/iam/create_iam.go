package iam

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFData represents the iam data used in templates
	TFData struct {
		TFUsers    []*TFUser
		TFRoles    []TFRole
		TFGroups   []TFGroup
		Section    string
		Subcommand string
	}

	// TFUser represents the user data used in templates
	TFUser struct {
		TFUserBasicInfo
		IsLocked   bool
		AuthGrants string
	}

	// TFUserBasicInfo represents user basic info data used in templates
	TFUserBasicInfo struct {
		ID                string
		FirstName         string
		LastName          string
		Email             string
		Country           string
		Phone             string
		TFAEnabled        bool
		ContactType       string
		JobTitle          string
		TimeZone          string
		SecondaryEmail    string
		MobilePhone       string
		Address           string
		City              string
		State             string
		ZipCode           string
		PreferredLanguage string
		SessionTimeOut    *int
	}

	// TFRole represents a role used in templates
	TFRole struct {
		RoleID          int64
		RoleName        string
		RoleDescription string
		GrantedRoles    []int
	}

	// TFGroup represents a group used in templates
	TFGroup struct {
		GroupID       int
		ParentGroupID int
		GroupName     string
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	// ErrFetchingUsers is returned when fetching users fails
	ErrFetchingUsers = errors.New("unable to fetch users under this account")
	// ErrFetchingUserByID is returned when fetching user by id fails
	ErrFetchingUserByID = errors.New("unable to fetch user by id")
)

// CmdCreateIAM is an entrypoint to create-iam command
func CmdCreateIAM(c *cli.Context) error {
	if c.NArg() == 0 {
		return showHelpCommandWithErr(c, fmt.Sprintf("One of the subcommands is required : %s", getSubcommandsNames(c)))
	}
	if !isSubcommandValid(c, c.Args().First()) {
		return showHelpCommandWithErr(c, fmt.Sprintf("Subcommand '%v' is invalid. Use one of valid subcommand: %s", c.Args().First(), getSubcommandsNames(c)))
	}
	return nil
}

func isSubcommandValid(ctx *cli.Context, subcommand string) bool {
	commands := ctx.App.Commands
	if len(commands) == 0 {
		return false
	}
	for _, c := range commands {
		if c.Name == subcommand {
			return true
		}
	}
	return false
}

func getSubcommandsNames(ctx *cli.Context) []string {
	var names []string
	for _, c := range ctx.App.Commands {
		names = append(names, c.Name)
	}
	return names
}

func showHelpCommandWithErr(c *cli.Context, stringErr string) error {
	if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
		return cli.Exit(color.RedString("Error displaying help command"), 1)
	}
	return cli.Exit(color.RedString(stringErr), 1)
}

func getTFUsers(ctx context.Context, client iam.IAM, users []iam.UserListItem) ([]*TFUser, error) {
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

func getTFUser(user *iam.User) (*TFUser, error) {
	authGrants, err := getUserAuthGrants(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %s for user with email %s", ErrMarshalUserAuthGrants, err, user.Email)
	}

	return &TFUser{
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
	}, nil
}

func getTFGroup(group *iam.Group) TFGroup {
	return TFGroup{
		GroupID:       int(group.GroupID),
		ParentGroupID: int(group.ParentGroupID),
		GroupName:     group.GroupName,
	}
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

func getTFRoles(roles []iam.Role) []TFRole {
	tfRoles := make([]TFRole, 0)
	for _, r := range roles {
		tfRoles = append(tfRoles, TFRole{
			RoleID:          r.RoleID,
			RoleName:        r.RoleName,
			RoleDescription: r.RoleDescription,
			GrantedRoles:    getGrantedRolesID(r.GrantedRoles),
		})
	}
	return tfRoles
}

func getGrantedRolesID(grantedRoles []iam.RoleGrantedRole) []int {
	rolesIDs := make([]int, 0, len(grantedRoles))
	for _, v := range grantedRoles {
		rolesIDs = append(rolesIDs, int(v.RoleID))
	}
	return rolesIDs
}
