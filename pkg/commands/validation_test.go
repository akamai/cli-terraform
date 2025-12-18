package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestValidatedAction(t *testing.T) {
	t.Run("empty validation list", func(t *testing.T) {
		action := func(_ *cli.Context) error {
			return nil
		}
		actionFunc := validatedAction(action)

		err := actionFunc(&cli.Context{})
		assert.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		action := func(_ *cli.Context) error {
			return fmt.Errorf("action error")
		}
		actionFunc := validatedAction(action)

		err := actionFunc(&cli.Context{})
		assert.ErrorContains(t, err, "action error")
	})

	t.Run("action with validation", func(t *testing.T) {
		action := func(_ *cli.Context) error {
			return nil
		}
		validation := func(_ *cli.Context) error {
			return nil
		}
		actionFunc := validatedAction(action, validation)

		err := actionFunc(&cli.Context{})
		assert.NoError(t, err)
	})

	t.Run("assert function call order", func(t *testing.T) {
		var callOrder []string
		action := func(_ *cli.Context) error {
			callOrder = append(callOrder, "action")
			return nil
		}
		firstValidation := func(_ *cli.Context) error {
			callOrder = append(callOrder, "first_validation")
			return nil
		}
		secondValidation := func(_ *cli.Context) error {
			callOrder = append(callOrder, "second_validation")
			return nil
		}
		actionFunc := validatedAction(action, firstValidation, secondValidation)

		err := actionFunc(&cli.Context{})
		assert.NoError(t, err)

		expectedCallOrder := []string{"first_validation", "second_validation", "action"}
		assert.Equal(t, expectedCallOrder, callOrder)
	})

	t.Run("validation error", func(t *testing.T) {
		action := func(_ *cli.Context) error {
			return fmt.Errorf("action error")
		}
		validation := func(_ *cli.Context) error {
			return fmt.Errorf("validation error")
		}
		actionFunc := validatedAction(action, validation)

		err := actionFunc(&cli.Context{})
		assert.ErrorContains(t, err, "validation error")
	})
}

func TestRequireValidWorkpath(t *testing.T) {

	t.Run("tfworkpath not set", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		flagset.String("tfworkpath", "", "") // flag definition

		ctx := cli.NewContext(app, flagset, nil)

		err := requireValidWorkpath(ctx)
		assert.NoError(t, err)
	})

	t.Run("tfworkpath is set to unknown path", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		flagset.String("tfworkpath", "", "") // flag definition
		err := flagset.Set("tfworkpath", "some/not/existing/path")
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)

		err = requireValidWorkpath(ctx)
		assert.ErrorContains(t, err, "Destination work path is not accessible")
	})

	t.Run("tfworkpath is set to known path", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		flagset.String("tfworkpath", "", "") // flag definition
		err := flagset.Set("tfworkpath", "./")
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)

		err = requireValidWorkpath(ctx)
		assert.NoError(t, err)
	})
}

func TestRequireNArguments(t *testing.T) {
	t.Run("require no arguments", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)

		ctx := cli.NewContext(app, flagset, nil)

		validateNArgsFunc := requireNArguments(0)

		err := validateNArgsFunc(ctx)
		assert.NoError(t, err)
	})

	t.Run("require 3 arguments", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagset.Parse([]string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)

		validateNArgsFunc := requireNArguments(3)

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{"arg1"}) // passing one argument
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.ArgsUsage = "<example usage>"

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		validateNArgsFunc := requireNArguments(3) // expecting 3 arguments

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <example usage>")
	})
}

func TestRequireNOptionalArguments(t *testing.T) {
	t.Run("require no arguments", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)

		ctx := cli.NewContext(app, flagset, nil)

		validateNArgsFunc := requiredAndOptionalArguments(0, 0)

		err := validateNArgsFunc(ctx)
		assert.NoError(t, err)
	})

	t.Run("require 1 arguments + 2 optional", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagset.Parse([]string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)

		validateNArgsFunc := requiredAndOptionalArguments(1, 2)

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
	})

	t.Run("require 3 arguments", func(t *testing.T) {
		app := cli.NewApp()
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagset.Parse([]string{"arg1", "arg2", "arg3"})
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)

		validateNArgsFunc := requiredAndOptionalArguments(3, 0)

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{"arg1"}) // passing one argument
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.ArgsUsage = "<example usage>"

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		validateNArgsFunc := requiredAndOptionalArguments(3, 0) // expecting 3 arguments

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <example usage>")
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{"arg1", "arg2"}) // passing two argument
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.ArgsUsage = "<arg1> <arg2> <arg3>"

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		validateNArgsFunc := requiredAndOptionalArguments(1, 2) // expecting 3 arguments

		err = validateNArgsFunc(ctx)
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <arg1> <arg2> <arg3>")
	})
}

func TestShowHelpCommandWithErr(t *testing.T) {
	cmdName := "create-command"

	t.Run("no error", func(t *testing.T) {
		customErr := fmt.Errorf("custom error for test")
		ctx := cli.Context{
			Command: &cli.Command{
				Name: cmdName,
			},
			App: &cli.App{
				CommandNotFound: func(_ *cli.Context, _ string) {},
				ErrWriter:       os.Stderr,
			},
		}
		err := showHelpCommandWithErr(&ctx, customErr.Error())
		require.NoError(t, err)
	})

	t.Run("help command is not found", func(t *testing.T) {
		customErr := fmt.Errorf("custom error for test")
		ctx := cli.Context{
			Command: &cli.Command{
				Name: cmdName,
			},
			App: &cli.App{
				CommandNotFound: nil,
				ErrWriter:       os.Stderr,
			},
		}
		err := showHelpCommandWithErr(&ctx, customErr.Error())
		assert.Equal(t, cli.Exit(color.RedString("No help topic for '%s'", cmdName), 3), err)
	})
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

func TestValidateSubCommands(t *testing.T) {
	t.Run("missed sub-command", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{}) // missed sub-command
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.Name = "create-command"
		ctx.App.Commands = []*cli.Command{
			{Name: "user"},
			{Name: "group"},
			{Name: "role"},
		}

		_ = validateSubCommands(ctx)
		assert.Contains(t, errBuffer.String(), "One of the subcommands is required : [user group role]")
	})

	t.Run("invalid sub-command", func(t *testing.T) {
		invalidSubcommand := "invalid_subcommand"
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{invalidSubcommand}) // invalid sub-command
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.Name = "create-command"
		ctx.App.Commands = []*cli.Command{
			{Name: "user"},
			{Name: "group"},
			{Name: "role"},
		}

		_ = validateSubCommands(ctx)
		assert.Contains(t, errBuffer.String(), fmt.Sprintf("Subcommand '%s' is invalid. Use one of valid subcommands: [user group role]", invalidSubcommand))
	})

	t.Run("valid sub-command", func(t *testing.T) {
		validSubcommand := "user"
		app := cli.NewApp()

		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{validSubcommand}) // valid sub-command
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.Name = "create-command"
		ctx.App.Commands = []*cli.Command{
			{Name: validSubcommand},
			{Name: "group"},
			{Name: "role"},
		}

		err = validateSubCommands(ctx)
		require.NoError(t, err)
		assert.Contains(t, errBuffer.String(), "")
	})

	t.Run("sub-commands must not be there", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		errBuffer := &bytes.Buffer{}
		app.ErrWriter = errBuffer

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagSet.Parse([]string{"some_subcommand"})
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagSet, nil)
		ctx.Command.Name = "create-command"
		ctx.App.Commands = []*cli.Command{} // expect no sub-commands for given command

		_ = validateSubCommands(ctx)
		assert.Contains(t, errBuffer.String(), fmt.Sprintf("Subcommands are not expected for '%s' command", ctx.Command.Name))
	})
}

func TestValidateRuleFormat(t *testing.T) {
	t.Run("rule-format flag not set - should skip validation", func(t *testing.T) {
		app := cli.NewApp()
		app.Writer = io.Discard
		app.ErrWriter = io.Discard

		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		flagSet.String("rule-format", "", "")

		cliCtx := cli.NewContext(app, flagSet, nil)

		validateFunc := validateRuleFormat(false)
		err := validateFunc(cliCtx)
		assert.NoError(t, err)
	})
}

func TestValidateRuleFormatWithPAPIClient(t *testing.T) {
	tests := map[string]struct {
		ruleFormatFlagValue     string
		rulesAsHclSet           bool
		isProperty              bool
		mockRuleFormatsResponse []string
		mockAPIError            error
		expectedError           string
	}{
		"valid rule format": {
			ruleFormatFlagValue:     "v2023-01-05",
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
		},
		"invalid rule format - correct format but non existent for property": {
			ruleFormatFlagValue:     "v2020-01-01",
			isProperty:              true,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: v2020-01-01. Must be vYYYY-MM-DD (with a leading \"v\"). The 'latest' value is allowed only for properties using the JSON rule format.",
		},
		"invalid rule format - correct format but non existent for include": {
			ruleFormatFlagValue:     "v2020-01-01",
			isProperty:              false,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: v2020-01-01. Must be vYYYY-MM-DD (with a leading \"v\").",
		},
		"invalid rule format - incorrect format for property": {
			ruleFormatFlagValue:     "2023-01-05",
			isProperty:              true,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: 2023-01-05. Must be vYYYY-MM-DD (with a leading \"v\"). The 'latest' value is allowed only for properties using the JSON rule format.",
		},
		"invalid rule format - incorrect format for include": {
			ruleFormatFlagValue:     "2023-01-05",
			isProperty:              false,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: 2023-01-05. Must be vYYYY-MM-DD (with a leading \"v\").",
		},
		"latest is valid for property and rules-as-hcl not set": {
			ruleFormatFlagValue:     "latest",
			isProperty:              true,
			rulesAsHclSet:           false,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
		},
		"latest is invalid for property and rules-as-hcl set": {
			ruleFormatFlagValue:     "latest",
			isProperty:              true,
			rulesAsHclSet:           true,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: latest. Must be vYYYY-MM-DD (with a leading \"v\"). The 'latest' value is allowed only for properties using the JSON rule format.",
		},
		"latest is invalid for include": {
			ruleFormatFlagValue:     "latest",
			isProperty:              false,
			mockRuleFormatsResponse: []string{"v2023-01-05", "v2023-05-30", "v2024-03-15", "latest"},
			expectedError:           "Invalid rule format version: latest. Must be vYYYY-MM-DD (with a leading \"v\").",
		},
		"API error when fetching rule formats versions": {
			ruleFormatFlagValue: "v2023-01-05",
			mockAPIError:        fmt.Errorf("API connection failed"),
			expectedError:       "Failed to fetch rule formats",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			app := cli.NewApp()
			app.Writer = io.Discard
			errBuffer := &bytes.Buffer{}
			app.ErrWriter = errBuffer

			flagSet := flag.NewFlagSet("test", flag.PanicOnError)
			flagSet.String("rule-format", "", "")
			flagSet.Bool("rules-as-hcl", false, "")
			err := flagSet.Set("rule-format", test.ruleFormatFlagValue)
			require.NoError(t, err)
			if test.rulesAsHclSet {
				err = flagSet.Set("rules-as-hcl", "true")
				require.NoError(t, err)
			}

			cliCtx := cli.NewContext(app, flagSet, nil)

			mockClient := new(papi.Mock)
			if test.mockAPIError != nil {
				mockClient.On("GetRuleFormats", mock.Anything).Return(nil, test.mockAPIError).Once()
			} else {
				response := &papi.GetRuleFormatsResponse{
					RuleFormats: papi.RuleFormatItems{
						Items: test.mockRuleFormatsResponse,
					},
				}
				mockClient.On("GetRuleFormats", mock.Anything).Return(response, nil).Once()
			}

			err = checkRuleFormat(cliCtx, mockClient, test.isProperty)

			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
