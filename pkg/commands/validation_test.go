package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func runValidationAction(cmd *cli.Command, args []string, validator actionValidator) error {
	cmd.CommandNotFound = func(context.Context, *cli.Command, string) {}
	cmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		return validator(ctx, cmd)
	}
	return cmd.Run(context.Background(), append([]string{cmd.Name}, args...))
}

func TestValidatedAction(t *testing.T) {
	t.Run("empty validation list", func(t *testing.T) {
		action := func(_ context.Context, _ *cli.Command) error {
			return nil
		}
		actionFunc := validatedAction(action)

		err := actionFunc(context.Background(), &cli.Command{})
		assert.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		action := func(_ context.Context, _ *cli.Command) error {
			return fmt.Errorf("action error")
		}
		actionFunc := validatedAction(action)

		err := actionFunc(context.Background(), &cli.Command{})
		assert.ErrorContains(t, err, "action error")
	})

	t.Run("action with validation", func(t *testing.T) {
		action := func(_ context.Context, _ *cli.Command) error {
			return nil
		}
		validation := func(_ context.Context, _ *cli.Command) error {
			return nil
		}
		actionFunc := validatedAction(action, validation)

		err := actionFunc(context.Background(), &cli.Command{})
		assert.NoError(t, err)
	})

	t.Run("assert function call order", func(t *testing.T) {
		var callOrder []string
		action := func(_ context.Context, _ *cli.Command) error {
			callOrder = append(callOrder, "action")
			return nil
		}
		firstValidation := func(_ context.Context, _ *cli.Command) error {
			callOrder = append(callOrder, "first_validation")
			return nil
		}
		secondValidation := func(_ context.Context, _ *cli.Command) error {
			callOrder = append(callOrder, "second_validation")
			return nil
		}
		actionFunc := validatedAction(action, firstValidation, secondValidation)

		err := actionFunc(context.Background(), &cli.Command{})
		assert.NoError(t, err)

		expectedCallOrder := []string{"first_validation", "second_validation", "action"}
		assert.Equal(t, expectedCallOrder, callOrder)
	})

	t.Run("validation error", func(t *testing.T) {
		action := func(_ context.Context, _ *cli.Command) error {
			return fmt.Errorf("action error")
		}
		validation := func(_ context.Context, _ *cli.Command) error {
			return fmt.Errorf("validation error")
		}
		actionFunc := validatedAction(action, validation)

		err := actionFunc(context.Background(), &cli.Command{})
		assert.ErrorContains(t, err, "validation error")
	})
}

func TestRequireValidWorkpath(t *testing.T) {
	t.Run("tfworkpath not set", func(t *testing.T) {
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "tfworkpath"},
			},
		}
		err := requireValidWorkpath(context.Background(), cmd)
		assert.NoError(t, err)
	})

	t.Run("tfworkpath is set to unknown path", func(t *testing.T) {
		cmd := &cli.Command{Flags: []cli.Flag{&cli.StringFlag{Name: "tfworkpath"}}}
		require.NoError(t, cmd.Set("tfworkpath", "some/not/existing/path"))
		err := requireValidWorkpath(context.Background(), cmd)
		assert.ErrorContains(t, err, "Destination work path is not accessible")
	})

	t.Run("tfworkpath is set to known path", func(t *testing.T) {
		cmd := &cli.Command{Flags: []cli.Flag{&cli.StringFlag{Name: "tfworkpath"}}}
		require.NoError(t, cmd.Set("tfworkpath", "./"))
		err := requireValidWorkpath(context.Background(), cmd)
		assert.NoError(t, err)
	})
}

func TestRequireNArguments(t *testing.T) {
	t.Run("require no arguments", func(t *testing.T) {
		cmd := &cli.Command{Name: "test-command"}
		err := runValidationAction(cmd, nil, requireNArguments(0))
		assert.NoError(t, err)
	})

	t.Run("require 3 arguments", func(t *testing.T) {
		cmd := &cli.Command{Name: "test-command"}
		err := runValidationAction(cmd, []string{"arg1", "arg2", "arg3"}, requireNArguments(3))
		assert.NoError(t, err)
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		cmd := &cli.Command{
			Name:      "test-command",
			ArgsUsage: "<example usage>",
			Writer:    io.Discard,
			ErrWriter: errBuffer,
		}

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		err := runValidationAction(cmd, []string{"arg1"}, requireNArguments(3))
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <example usage>")
	})
}

func TestRequireNOptionalArguments(t *testing.T) {
	t.Run("require no arguments", func(t *testing.T) {
		cmd := &cli.Command{Name: "test-command"}
		err := runValidationAction(cmd, nil, requiredAndOptionalArguments(0, 0))
		assert.NoError(t, err)
	})

	t.Run("require 1 arguments + 2 optional", func(t *testing.T) {
		cmd := &cli.Command{Name: "test-command"}
		err := runValidationAction(cmd, []string{"arg1", "arg2", "arg3"}, requiredAndOptionalArguments(1, 2))
		assert.NoError(t, err)
	})

	t.Run("require 3 arguments", func(t *testing.T) {
		cmd := &cli.Command{Name: "test-command"}
		err := runValidationAction(cmd, []string{"arg1", "arg2", "arg3"}, requiredAndOptionalArguments(3, 0))
		assert.NoError(t, err)
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		cmd := &cli.Command{
			Name:      "test-command",
			ArgsUsage: "<example usage>",
			Writer:    io.Discard,
			ErrWriter: errBuffer,
		}

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		err := runValidationAction(cmd, []string{"arg1"}, requiredAndOptionalArguments(3, 0))
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <example usage>")
	})

	t.Run("error not enough arguments", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		cmd := &cli.Command{
			Name:      "test-command",
			ArgsUsage: "<arg1> <arg2> <arg3>",
			Writer:    io.Discard,
			ErrWriter: errBuffer,
		}

		exitOsCalled := false
		// patch osExiter
		defer func(restore func(_ int)) {
			osExiter = restore
		}(osExiter)
		osExiter = func(_ int) {
			exitOsCalled = true
		}

		err := runValidationAction(cmd, []string{"arg1", "arg2"}, requiredAndOptionalArguments(1, 2))
		assert.NoError(t, err)
		assert.True(t, exitOsCalled)
		assert.Contains(t, errBuffer.String(), "Invalid arguments usage, next arguments are required: <arg1> <arg2> <arg3>")
	})
}

func TestShowHelpCommandWithErr(t *testing.T) {
	cmdName := "create-command"
	t.Run("no error", func(t *testing.T) {
		cmd := &cli.Command{
			Name:            cmdName,
			Writer:          io.Discard,
			ErrWriter:       io.Discard,
			CommandNotFound: func(context.Context, *cli.Command, string) {},
		}
		err := showHelpCommandWithErr(context.Background(), cmd, "custom error for test")
		require.NoError(t, err)
	})
	t.Run("help command is not found", func(t *testing.T) {
		cmd := &cli.Command{Name: cmdName, Writer: io.Discard, ErrWriter: io.Discard}
		err := showHelpCommandWithErr(context.Background(), cmd, "custom error for test")
		assert.Equal(t, cli.Exit(color.RedString("No help topic for '%s'", cmdName), 3), err)
	})
}

func TestRequireNArgumentsShowsCommandHelp(t *testing.T) {
	previousOSExiter := osExiter
	t.Cleanup(func() {
		osExiter = previousOSExiter
	})

	tests := []struct {
		name        string
		commandPath []string
		argsUsage   string
	}{
		{
			name:        "top-level command",
			commandPath: []string{"export-property"},
			argsUsage:   "<property name>",
		},
		{
			name:        "nested command",
			commandPath: []string{"export-iam", "group"},
			argsUsage:   "<group_id>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			var errOutput bytes.Buffer
			exitCode := 0
			osExiter = func(code int) {
				exitCode = code
			}

			app := &cli.Command{
				Name:      "akamai-terraform",
				Writer:    &output,
				ErrWriter: &errOutput,
			}
			parent := app
			for _, name := range test.commandPath {
				command := &cli.Command{
					Name:      name,
					Writer:    &output,
					ErrWriter: &errOutput,
				}
				parent.Commands = append(parent.Commands, command)
				parent = command
			}
			parent.ArgsUsage = test.argsUsage
			parent.Action = func(ctx context.Context, cmd *cli.Command) error {
				return requireNArguments(1)(ctx, cmd)
			}

			args := append([]string{app.Name}, test.commandPath...)
			err := app.Run(context.Background(), args)

			require.NoError(t, err)
			assert.Equal(t, 1, exitCode)
			assert.Contains(t, errOutput.String(), "Invalid arguments usage, next arguments are required: "+test.argsUsage)
			assert.Contains(t, output.String(), "USAGE:")
			assert.NotContains(t, errOutput.String(), "No help topic")
		})
	}
}

func TestGetSubcommandsNames(t *testing.T) {
	tests := map[string]struct {
		cmd      *cli.Command
		expected []string
	}{
		"subcommands exist": {
			cmd: &cli.Command{Commands: []*cli.Command{
				{Name: "user"},
				{Name: "group"},
				{Name: "role"},
			}},
			expected: []string{"user", "group", "role"},
		},
		"subcommands don't exist": {
			cmd: &cli.Command{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, getSubcommandsNames(test.cmd))
		})
	}
}

func TestValidateSubCommands(t *testing.T) {
	runValidation := func(args []string, subcommands []*cli.Command, errWriter io.Writer) error {
		cmd := &cli.Command{
			Name:      "create-command",
			ArgsUsage: "<subcommand>",
			Writer:    io.Discard,
			ErrWriter: errWriter,
		}
		validator := func(ctx context.Context, cmd *cli.Command) error {
			cmd.Commands = subcommands
			return validateSubCommands(ctx, cmd)
		}
		return runValidationAction(cmd, args, validator)
	}

	t.Run("missed sub-command", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		subcommands := []*cli.Command{{Name: "user"}, {Name: "group"}, {Name: "role"}}
		_ = runValidation(nil, subcommands, errBuffer)
		assert.Contains(t, errBuffer.String(), "One of the subcommands is required : [user group role]")
	})
	t.Run("invalid sub-command", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		_ = runValidation([]string{"invalid_subcommand"}, []*cli.Command{{Name: "user"}, {Name: "group"}, {Name: "role"}}, errBuffer)
		assert.Contains(t, errBuffer.String(), "Subcommand 'invalid_subcommand' is invalid. Use one of valid subcommands: [user group role]")
	})
	t.Run("valid sub-command", func(t *testing.T) {
		err := runValidation([]string{"user"}, []*cli.Command{{Name: "user"}, {Name: "group"}, {Name: "role"}}, io.Discard)
		require.NoError(t, err)
	})
	t.Run("sub-commands must not be there", func(t *testing.T) {
		errBuffer := &bytes.Buffer{}
		_ = runValidation([]string{"some_subcommand"}, nil, errBuffer)
		assert.Contains(t, errBuffer.String(), "Subcommands are not expected for 'create-command' command")
	})
}

func TestValidateRuleFormat(t *testing.T) {
	t.Run("rule-format flag not set - should skip validation", func(t *testing.T) {
		cmd := &cli.Command{
			Flags: []cli.Flag{&cli.StringFlag{Name: "rule-format"}},
		}
		validateFunc := validateRuleFormat(false)
		err := validateFunc(context.Background(), cmd)
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
			cmd := &cli.Command{
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "rule-format"},
					&cli.BoolFlag{Name: "rules-as-hcl"},
				},
			}
			require.NoError(t, cmd.Set("rule-format", test.ruleFormatFlagValue))
			if test.rulesAsHclSet {
				require.NoError(t, cmd.Set("rules-as-hcl", "true"))
			}

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

			err := checkRuleFormat(context.Background(), cmd, mockClient, test.isProperty)

			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
