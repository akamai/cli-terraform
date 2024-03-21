package cli

import (
	"bytes"
	"context"
	"flag"
	"io"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/cli/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func newContextFromStringSlice(ss []string, app *cli.App) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse(ss)
	return cli.NewContext(app, set, nil)
}

func newTemplateApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []*cli.Command{{Name: "some-command", Aliases: []string{"other-command"}}, {Name: "help"}, {Name: "list"}}
	return app
}

func Test_sessionRequired(t *testing.T) {
	tests := map[string]struct {
		c        func() *cli.Context
		expected bool
	}{
		"no command": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{""}, cli.NewApp())
			},
			expected: false,
		},
		"help": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"help"}, newTemplateApp())
			},
			expected: false,
		},
		"help for command": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"help", "create-something"}, newTemplateApp())
			},
			expected: false,
		},
		"help for subcommand": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"help", "create-something", "subcommand"}, newTemplateApp())
			},
			expected: false,
		},
		"--help for command": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"create-something", "--help"}, newTemplateApp())
			},
			expected: false,
		},
		"--help for subcommand": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"create-something", "subcommand", "--help"}, newTemplateApp())
			},
			expected: false,
		},
		"list": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"list"}, newTemplateApp())
			},
			expected: false,
		},
		"unknown command": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"unknown"}, newTemplateApp())
			},
			expected: false,
		},
		"some command which requires auth": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"some-command"}, newTemplateApp())
			},
			expected: true,
		},
		"use alias": {
			c: func() *cli.Context {
				return newContextFromStringSlice([]string{"other-command"}, newTemplateApp())
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, sessionRequired(test.c()))
		})
	}
}

func TestPutLoggerInContext(t *testing.T) {
	t.Setenv("AKAMAI_LOG", "debug")
	app := cli.NewApp()
	buffer := &bytes.Buffer{}
	app.Writer = buffer

	ctx := context.Background()

	cliCtx := &cli.Context{
		Context: ctx,
		App:     app,
	}
	err := putLoggerInContext(cliCtx)
	assert.NoError(t, err)

	logger := log.FromContext(cliCtx.Context)
	buffer.Reset()
	logger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")

	sess, err := session.New()
	assert.NoError(t, err)
	logger = sess.Log(cliCtx.Context)
	buffer.Reset()
	logger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")
}

func TestDeprecationInfo(t *testing.T) {
	app := cli.NewApp()
	app.Commands = []*cli.Command{{Name: "export-command", Aliases: []string{"create-command"}}, {Name: "help"}, {Name: "list"}}

	buf := &bytes.Buffer{}
	app.Writer = buf
	app.ErrWriter = io.Discard
	app.Before = ensureBefore(deprecationInfoForCreateCommands)

	tests := map[string]struct {
		args          []string
		expectWarning bool
	}{
		"create": {
			args:          []string{"cmd", "create-command"},
			expectWarning: true,
		},
		"export": {
			args:          []string{"cmd", "export-command"},
			expectWarning: false,
		},
		"help create": {
			args:          []string{"cmd", "help", "create-command"},
			expectWarning: true,
		},
		"help export": {
			args:          []string{"cmd", "help", "export-command"},
			expectWarning: false,
		},
	}

	for name, test := range tests {
		buf.Reset()
		t.Run(name, func(t *testing.T) {
			err := app.Run(test.args)
			assert.NoError(t, err)

			if test.expectWarning {
				assert.Contains(t, buf.String(), "Warning: create command names are now deprecated, use export commands instead")
			} else {
				assert.NotContains(t, buf.String(), "Warning")
			}
		})
	}
}
func TestDeprecationInfoForSchemaFlag(t *testing.T) {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name: "export-property", Aliases: []string{"create-property"}, Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "rules-as-hcl",
					Aliases: []string{"schema"},
				},
			},
		},
		{
			Name: "export-imaging", Aliases: []string{"create-imaging"}, Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "policy-as-hcl",
					Aliases: []string{"schema"},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	app.Writer = buf
	app.ErrWriter = io.Discard
	app.Before = ensureBefore(deprecationInfoForSchemaFlags)

	tests := map[string]struct {
		args          []string
		expectWarning bool
	}{
		"property cmd when --schema flag passed": {
			args:          []string{"cmd", "create-property", "--schema"},
			expectWarning: true,
		},
		"property cmd when no --schema flag passed": {
			args:          []string{"cmd", "export-property"},
			expectWarning: false,
		},
		"imaging cmd when --schema flag passed": {
			args:          []string{"cmd", "create-imaging", "--schema"},
			expectWarning: true,
		},
		"imaging cmd when no --schema flag passed": {
			args:          []string{"cmd", "export-imaging"},
			expectWarning: false,
		},
	}

	for name, test := range tests {
		buf.Reset()
		t.Run(name, func(t *testing.T) {
			err := app.Run(test.args)
			assert.NoError(t, err)

			if test.expectWarning {
				assert.Contains(t, buf.String(), "Warning: flag --schema is now deprecated")
			} else {
				assert.NotContains(t, buf.String(), "Warning")
			}
		})
	}
}
