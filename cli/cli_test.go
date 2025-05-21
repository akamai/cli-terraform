package cli

import (
	"bytes"
	"context"
	"flag"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/akamai/cli/v2/pkg/log"
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

	ctxLogger := log.FromContext(cliCtx.Context)
	buffer.Reset()
	ctxLogger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")

	sess, err := session.New()
	assert.NoError(t, err)
	sessLogger := sess.Log(cliCtx.Context)
	buffer.Reset()
	sessLogger.Info("oops")
	assert.Contains(t, buffer.String(), "oops")
}
