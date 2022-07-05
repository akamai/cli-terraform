package commands

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

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
			if err := cli.ShowCommandHelp(ctx, ctx.Command.Name); err != nil {
				return cli.Exit(color.RedString("Error displaying help command"), 1)
			}
			return cli.Exit(color.RedString("Invalid arguments usage, should be: %s", ctx.Command.ArgsUsage), 1)
		}
		return nil
	}
}
