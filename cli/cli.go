// Package cli contains code for managing the command line interface functionality.
package cli

import (
	"context"
	"errors"
	"os"

	sesslog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/cli-terraform/pkg/commands"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	akacli "github.com/akamai/cli/pkg/app"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/log"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

var (
	// Version holds current version of cli-terraform
	Version = "2.0.0"
)

// Run initializes the cli and runs it
func Run() error {
	term := terminal.Color()
	ctx := context.Background()
	ctx = terminal.Context(ctx, term)

	app := akacli.CreateAppTemplate(ctx, "terraform",
		"A CLI Plugin for Akamai Terraform Provider",
		"Export selected resources for faster adoption in Terraform.",
		Version)

	cmds, err := commands.CommandLocator()
	if err != nil {
		return errors.New(color.RedString("An error occurred initializing commands: %s", err))
	}
	if len(cmds) > 0 {
		app.Commands = append(cmds, app.Commands...)
	}

	app.Before = ensureBefore(putLoggerInContext, putSessionInContext)
	return app.RunContext(ctx, os.Args)
}

func ensureBefore(bfs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(c *cli.Context) error {
		for _, bf := range bfs {
			err := bf(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func sessionRequired(c *cli.Context) bool {
	command := c.Args().First()

	for _, cmd := range []string{"help", "list", ""} {
		if cmd == command {
			return false
		}
	}

	tail := c.Args().Tail()
	if len(tail) > 0 && tail[len(tail)-1] == "--help" {
		return false
	}

	for _, cmd := range c.App.Commands {
		if cmd.Name == command || sliceContains(cmd.Aliases, command) {
			return true
		}
	}

	return false
}

func sliceContains(slc []string, c string) bool {
	for _, s := range slc {
		if s == c {
			return true
		}
	}
	return false
}

func putSessionInContext(c *cli.Context) error {
	if !sessionRequired(c) {
		return nil
	}
	s, err := edgegrid.InitializeSession(c)
	if err != nil {
		return err
	}
	c.Context = edgegrid.WithSession(c.Context, s)

	return nil
}

func putLoggerInContext(c *cli.Context) error {
	c.Context = log.SetupContext(c.Context, c.App.Writer)

	handler := log.FromContext(c.Context).Handler()
	sessionLogger := sesslog.NewSlogAdapter(handler)

	c.Context = session.ContextWithOptions(c.Context, session.WithContextLog(sessionLogger))

	return nil
}
