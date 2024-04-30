// Package cli contains code for managing the command line interface functionality.
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	sesslog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
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
	Version = "1.19.0"
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
		return fmt.Errorf(color.RedString("An error occurred initializing commands: %s", err))
	}
	if len(cmds) > 0 {
		app.Commands = append(cmds, app.Commands...)
	}

	app.Before = ensureBefore(putLoggerInContext, putSessionInContext, deprecationInfoForCreateCommands, deprecationInfoForSchemaFlags)
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

func deprecationInfoForCreateCommands(c *cli.Context) error {
	if !c.Args().Present() {
		return nil
	}
	command := c.Args().First()
	if command == "help" {
		command = c.Args().Get(1)
	}
	if strings.HasPrefix(command, "create-") {
		fmt.Fprintln(c.App.Writer, color.HiYellowString("Warning:"), "create command names are now deprecated, use export commands instead.")
		fmt.Fprintln(c.App.Writer)
	}
	return nil
}

func deprecationInfoForSchemaFlags(c *cli.Context) error {
	if !c.Args().Present() {
		return nil
	}
	command := c.Args().First()
	newFlagName := ""
	if strings.HasSuffix(command, "property") {
		newFlagName = "--rules-as-hcl"
	}
	if strings.HasSuffix(command, "imaging") {
		newFlagName = "--policy-as-hcl"
	}
	if newFlagName == "" {
		return nil
	}

	hasSchemaFlag := false
	for _, f := range c.Args().Slice() {
		if f == "--schema" {
			hasSchemaFlag = true
			break
		}
	}

	if hasSchemaFlag {
		fmt.Fprint(c.App.Writer, color.HiYellowString("Warning: "))
		fmt.Fprintf(c.App.Writer, "flag --schema is now deprecated, use %s flag instead.\n\n", newFlagName)
	}
	return nil
}
