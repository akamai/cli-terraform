package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var osExiter = os.Exit

type actionValidator func(*cli.Context) error

func validatedAction(action cli.ActionFunc, validators ...actionValidator) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		for _, validator := range validators {
			if err := validator(ctx); err != nil {
				return err
			}
		}
		return action(ctx)
	}
}

func requireValidWorkpath(ctx *cli.Context) error {
	if !ctx.IsSet("tfworkpath") {
		return nil
	}
	tfWorkPath := ctx.String("tfworkpath")
	tfWorkPath = filepath.FromSlash(tfWorkPath)
	if stat, err := os.Stat(tfWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}
	return nil
}

func requireNArguments(n int) actionValidator {
	return func(ctx *cli.Context) error {
		if ctx.NArg() != n {
			if err := showHelpCommandWithErr(ctx, fmt.Sprintf("Invalid arguments usage, next arguments are required: %s", ctx.Command.ArgsUsage)); err != nil {
				return err
			}
			osExiter(1)
		}
		return nil
	}
}

func validateSubCommands(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return showHelpCommandWithErr(ctx, fmt.Sprintf("One of the subcommands is required : %s", getSubcommandsNames(ctx)))
	}

	subcommand := ctx.Args().First()
	commands := ctx.App.Commands
	if len(commands) == 0 {
		return showHelpCommandWithErr(ctx, fmt.Sprintf("Subcommands are not expected for '%s' command", ctx.Command.Name))
	}
	for _, c := range commands {
		if c.Name == subcommand {
			return nil
		}
	}
	return showHelpCommandWithErr(ctx, fmt.Sprintf("Subcommand '%v' is invalid. Use one of valid subcommands: %s", ctx.Args().First(), getSubcommandsNames(ctx)))
}

func getSubcommandsNames(ctx *cli.Context) []string {
	var names []string
	for _, c := range ctx.App.Commands {
		names = append(names, c.Name)
	}
	return names
}

func showHelpCommandWithErr(c *cli.Context, stringErr string) error {
	_, err := fmt.Fprintf(c.App.ErrWriter, "%s\n\n", color.RedString(stringErr))
	if err != nil {
		return err
	}
	return cli.ShowCommandHelp(c, c.Command.Name)
}
