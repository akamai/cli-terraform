package cli

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
	akacli "github.com/akamai/cli/v2/pkg/app"
	"github.com/akamai/cli/v2/pkg/log"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

type beforeContextKey struct{}

func newTemplateApp() *cli.Command {
	noop := func(context.Context, *cli.Command) error { return nil }
	return &cli.Command{
		Name:      "akamai-terraform",
		Writer:    io.Discard,
		ErrWriter: io.Discard,
		Commands: []*cli.Command{
			{
				Name:    "some-command",
				Aliases: []string{"other-command"},
				Action:  noop,
				Commands: []*cli.Command{
					{Name: "subcommand", Action: noop},
				},
			},
			{Name: "help", Action: noop},
			{Name: "list", Action: noop},
		},
	}
}

func runSessionRequired(cmd *cli.Command, args []string) (bool, error) {
	var required bool
	cmd.Before = func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		required = sessionRequired(cmd)
		return ctx, nil
	}
	cmd.Action = func(context.Context, *cli.Command) error { return nil }
	cmd.CommandNotFound = func(context.Context, *cli.Command, string) {}
	err := cmd.Run(context.Background(), append([]string{cmd.Name}, args...))
	return required, err
}

func Test_sessionRequired(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected bool
	}{
		"no command": {
			args:     nil,
			expected: false,
		},
		"help": {
			args:     []string{"help"},
			expected: false,
		},
		"help for command": {
			args:     []string{"help", "create-something"},
			expected: false,
		},
		"help for subcommand": {
			args:     []string{"help", "create-something", "subcommand"},
			expected: false,
		},
		"--help for command": {
			args:     []string{"some-command", "--help"},
			expected: false,
		},
		"--help for subcommand": {
			args:     []string{"some-command", "subcommand", "--help"},
			expected: false,
		},
		"list": {
			args:     []string{"list"},
			expected: false,
		},
		"unknown command": {
			args:     []string{"unknown"},
			expected: false,
		},
		"some command which requires auth": {
			args:     []string{"some-command"},
			expected: true,
		},
		"use alias": {
			args:     []string{"other-command"},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := newTemplateApp()
			actual, err := runSessionRequired(app, test.args)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestEnsureBeforePassesUpdatedContext(t *testing.T) {
	first := func(ctx context.Context, _ *cli.Command) (context.Context, error) {
		return context.WithValue(ctx, beforeContextKey{}, "first"), nil
	}
	second := func(ctx context.Context, _ *cli.Command) (context.Context, error) {
		assert.Equal(t, "first", ctx.Value(beforeContextKey{}))
		return context.WithValue(ctx, beforeContextKey{}, "second"), nil
	}

	ctx, err := ensureBefore(first, second)(context.Background(), &cli.Command{})
	assert.NoError(t, err)
	assert.Equal(t, "second", ctx.Value(beforeContextKey{}))
}

func TestPutLoggerInContext(t *testing.T) {
	t.Setenv("AKAMAI_LOG", "debug")
	buffer := &bytes.Buffer{}
	ctx := context.Background()

	cmd := &cli.Command{
		Writer: buffer,
	}
	ctx, err := putLoggerInContext(ctx, cmd)
	assert.NoError(t, err)

	ctxLogger := log.FromContext(ctx)
	buffer.Reset()
	ctxLogger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")

	sess, err := session.New()
	assert.NoError(t, err)
	sessLogger := sess.Log(ctx)
	buffer.Reset()
	sessLogger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")
}

func TestHelpForLeafCommands(t *testing.T) {
	tests := map[string]struct {
		commandName         string
		expectedDescription string
	}{
		"simple leaf command": {
			commandName:         "leaf",
			expectedDescription: "Test leaf command.",
		},
		"hyphenated leaf command": {
			commandName:         "test-leaf",
			expectedDescription: "Test hyphenated leaf command.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))

			app := akacli.CreateAppTemplate(
				ctx,
				"terraform",
				"A CLI Plugin for Akamai Terraform Provider",
				"Export selected resources for faster adoption in Terraform.",
				Version,
			)

			output := &bytes.Buffer{}
			app.Writer = output

			actionCalled := false
			command := &cli.Command{
				Name:        test.commandName,
				Description: test.expectedDescription,
				Action: func(context.Context, *cli.Command) error {
					actionCalled = true
					return nil
				},
			}
			app.Commands = append([]*cli.Command{command}, app.Commands...)

			err := app.Run(ctx, []string{"akamai-terraform", "help", test.commandName})
			outputText := output.String()

			assert.NoError(t, err)
			assert.Contains(t, outputText, "Name:\n  akamai-terraform "+test.commandName+"\n")
			assert.Contains(t, outputText, "Description:\n  "+test.expectedDescription+"\n")
			assert.NotContains(t, outputText, "No help topic")
			assert.False(t, actionCalled, "help should not execute the command action")
		})
	}
}
