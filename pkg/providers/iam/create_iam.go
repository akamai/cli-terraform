package iam

import (
	"embed"
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFUserData represents the user data used in templates
	TFUserData struct {
		TFUserBasicInfo
		IsLocked   bool
		AuthGrants string
		Section    string
		Roles      []TFRole
		Groups     []TFGroup
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
