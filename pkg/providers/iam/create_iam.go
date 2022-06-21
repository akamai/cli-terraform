package iam

import (
	"embed"
	"encoding/json"
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
