package commands

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestValidatedAction(t *testing.T) {
	t.Run("empty validation list", func(t *testing.T) {
		action := func(ctx *cli.Context) error {
			return nil
		}
		actionFunc := validatedAction(action)

		err := actionFunc(&cli.Context{})
		assert.NoError(t, err)
	})

	t.Run("action error", func(t *testing.T) {
		action := func(ctx *cli.Context) error {
			return fmt.Errorf("action error")
		}
		actionFunc := validatedAction(action)

		err := actionFunc(&cli.Context{})
		assert.ErrorContains(t, err, "action error")
	})

	t.Run("action with validation", func(t *testing.T) {
		action := func(ctx *cli.Context) error {
			return nil
		}
		validation := func(ctx *cli.Context) error {
			return nil
		}
		actionFunc := validatedAction(action, validation)

		err := actionFunc(&cli.Context{})
		assert.NoError(t, err)
	})

	t.Run("assert function call order", func(t *testing.T) {
		var callOrder []string
		action := func(ctx *cli.Context) error {
			callOrder = append(callOrder, "action")
			return nil
		}
		firstValidation := func(ctx *cli.Context) error {
			callOrder = append(callOrder, "first_validation")
			return nil
		}
		secondValidation := func(ctx *cli.Context) error {
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
		action := func(ctx *cli.Context) error {
			return fmt.Errorf("action error")
		}
		validation := func(ctx *cli.Context) error {
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
		flagset := flag.NewFlagSet("test", flag.PanicOnError)
		err := flagset.Parse([]string{"arg1"}) // passing one argument
		assert.NoError(t, err)

		ctx := cli.NewContext(app, flagset, nil)
		ctx.Command.ArgsUsage = "<example usage>"

		validateNArgsFunc := requireNArguments(3) // expecting 3 arguments

		err = validateNArgsFunc(ctx)
		assert.ErrorContains(t, err, "Invalid arguments usage, should be: <example usage>")
	})
}
