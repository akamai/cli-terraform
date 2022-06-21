package iam

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestIsSubcommandValid(t *testing.T) {
	tests := map[string]struct {
		ctx        *cli.Context
		subcommand string
		isValid    bool
	}{
		"valid subcommand": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{
						{
							Name: "user",
						},
						{
							Name: "group",
						},
						{
							Name: "role",
						},
					},
				},
			},
			subcommand: "user",
			isValid:    true,
		},
		"invalid subcommand": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{
						{
							Name: "user",
						},
						{
							Name: "group",
						},
						{
							Name: "role",
						},
					},
				},
			},
			subcommand: "asdf",
			isValid:    false,
		},
		"no subcommands on context": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{},
				},
			},
			subcommand: "asdf",
			isValid:    false,
		},
		"empty subcommand": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{{Name: "user"}},
				},
			},
			subcommand: "",
			isValid:    false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			isValid := isSubcommandValid(test.ctx, test.subcommand)
			assert.Equal(t, test.isValid, isValid)
		})
	}
}

func TestGetSubcommandsNames(t *testing.T) {
	tests := map[string]struct {
		ctx                *cli.Context
		expectedSubcommand []string
	}{
		"subcommands exist": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{
						{
							Name: "user",
						},
						{
							Name: "group",
						},
						{
							Name: "role",
						},
					},
				},
			},
			expectedSubcommand: []string{"user", "group", "role"},
		},
		"subcommands don't exist": {
			ctx: &cli.Context{
				App: &cli.App{
					Commands: []*cli.Command{},
				},
			},
			expectedSubcommand: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subcommands := getSubcommandsNames(test.ctx)
			assert.Equal(t, test.expectedSubcommand, subcommands)
		})
	}
}

func TestShowHelpCommandWithErr(t *testing.T) {
	tests := map[string]struct {
		ctx         *cli.Context
		stringErr   string
		expectedErr error
	}{
		"return proper error": {
			ctx: &cli.Context{
				Command: &cli.Command{
					Name: "create-iam",
				},
				App: &cli.App{
					Commands: []*cli.Command{
						{
							Name: "user",
						},
					},
					CommandNotFound: func(c *cli.Context, command string) {},
				},
			},
			stringErr:   "valid error",
			expectedErr: cli.Exit(color.RedString("valid error"), 1),
		},
		"help command is nil": {
			ctx: &cli.Context{
				Command: &cli.Command{
					Name: "create-iam",
				},
				App: &cli.App{
					Commands: []*cli.Command{
						{
							Name: "user",
						},
					},
					CommandNotFound: nil,
				},
			},
			stringErr:   "valid error",
			expectedErr: cli.Exit(color.RedString("Error displaying help command"), 1),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := showHelpCommandWithErr(test.ctx, test.stringErr)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestGetGrantedRolesID(t *testing.T) {
	tests := map[string]struct {
		grantedRoles []iam.RoleGrantedRole
		expectedIDs  []int
	}{
		"granted roles": {
			grantedRoles: []iam.RoleGrantedRole{
				{
					RoleID: 123,
				},
				{
					RoleID: 321,
				},
				{
					RoleID: 456,
				},
			},
			expectedIDs: []int{123, 321, 456},
		},
		"empty granted roles": {
			grantedRoles: []iam.RoleGrantedRole{},
			expectedIDs:  []int{},
		},
		"nil granted roles": {
			grantedRoles: nil,
			expectedIDs:  []int{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			grantedRolesIDs := getGrantedRolesID(test.grantedRoles)
			assert.Equal(t, test.expectedIDs, grantedRolesIDs)
		})
	}
}
