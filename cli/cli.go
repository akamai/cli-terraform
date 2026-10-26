// Package cli contains code for managing the command line interface functionality.
package cli

import (
	"context"
	"os"

	sesslog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/akamai/cli-terraform/v3/pkg/commands"
	"github.com/akamai/cli-terraform/v3/pkg/edgegrid"
	akacli "github.com/akamai/cli/v2/pkg/app"
	"github.com/akamai/cli/v2/pkg/log"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v3"
)

var (
	// Version holds current version of cli-terraform
	Version = "3.1.0"
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

	cmds := commands.CommandLocator()
	if len(cmds) > 0 {
		app.Commands = append(cmds, app.Commands...)
	}

	app.Before = ensureBefore(putLoggerInContext, putSessionInContext)
	return app.Run(ctx, os.Args)
}

func ensureBefore(bfs ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		for _, bf := range bfs {
			nextCtx, err := bf(ctx, cmd)
			if err != nil {
				return ctx, err
			}
			if nextCtx != nil {
				ctx = nextCtx
			}
		}
		return ctx, nil
	}
}

func sessionRequired(cmd *cli.Command) bool {
	command := cmd.Args().First()

	for _, cmd := range []string{"help", "list", ""} {
		if cmd == command {
			return false
		}
	}

	tail := cmd.Args().Tail()
	if len(tail) > 0 && tail[len(tail)-1] == "--help" {
		return false
	}

	for _, c := range cmd.Commands {
		if c.Name == command || sliceContains(c.Aliases, command) {
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

func putSessionInContext(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	if !sessionRequired(cmd) {
		return ctx, nil
	}
	s, err := edgegrid.InitializeSession(ctx, cmd)
	if err != nil {
		return ctx, err
	}

	return edgegrid.WithSession(ctx, s), nil
}

func putLoggerInContext(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	ctx = log.SetupContext(ctx, cmd.Writer)

	handler := log.FromContext(ctx).Handler()
	sessionLogger := sesslog.NewSlogAdapter(handler)

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(sessionLogger))
	return ctx, nil
}
