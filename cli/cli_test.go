package cli

import (
	"flag"
	"testing"

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
	app.Commands = []*cli.Command{{Name: "some-command"}, {Name: "help"}, {Name: "list"}}
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, sessionRequired(test.c()))
		})
	}
}
