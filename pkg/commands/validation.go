package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/papi"
	"github.com/akamai/cli-terraform/v3/pkg/edgegrid"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/urfave/cli/v3"
)

var osExiter = os.Exit

type actionValidator func(context.Context, *cli.Command) error

func validatedAction(action cli.ActionFunc, validators ...actionValidator) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		for _, validator := range validators {
			if err := validator(ctx, cmd); err != nil {
				return err
			}
		}
		return action(ctx, cmd)
	}
}

func requireValidWorkpath(_ context.Context, cmd *cli.Command) error {
	if !cmd.IsSet("tfworkpath") {
		return nil
	}
	tfWorkPath := cmd.String("tfworkpath")
	tfWorkPath = filepath.FromSlash(tfWorkPath)
	if stat, err := os.Stat(tfWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}
	return nil
}

func requireNArguments(n int) actionValidator {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != n {
			if err := showHelpCommandWithErr(ctx, cmd, fmt.Sprintf("Invalid arguments usage, next arguments are required: %s", cmd.ArgsUsage)); err != nil {
				return err
			}
			osExiter(1)
		}
		return nil
	}
}

func requiredAndOptionalArguments(required, optional int) actionValidator {
	return func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != required && cmd.NArg() != optional+required {
			if err := showHelpCommandWithErr(ctx, cmd, fmt.Sprintf("Invalid arguments usage, next arguments are required: %s", cmd.ArgsUsage)); err != nil {
				return err
			}
			osExiter(1)
		}
		return nil
	}
}

func validateSplitDepth(_ context.Context, cmd *cli.Command) error {
	if cmd.IsSet("split-depth") && !cmd.IsSet("rules-as-hcl") {
		return cli.Exit(color.RedString(`"split-depth" option must be used along with "rules-as-hcl"`), 1)
	}
	if cmd.IsSet("split-depth") && cmd.Int("split-depth") < 0 {
		return cli.Exit(color.RedString(`"split-depth" cannot have a negative value`), 1)
	}
	return nil
}

func validateSubCommands(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() == 0 {
		return showHelpCommandWithErr(ctx, cmd, fmt.Sprintf("One of the subcommands is required : %s", getSubcommandsNames(cmd)))
	}

	subcommand := cmd.Args().First()
	commands := cmd.Commands
	if len(commands) == 0 {
		return showHelpCommandWithErr(ctx, cmd, fmt.Sprintf("Subcommands are not expected for '%s' command", cmd.Name))
	}
	for _, c := range commands {
		if c.Name == subcommand {
			return nil
		}
	}
	return showHelpCommandWithErr(ctx, cmd, fmt.Sprintf("Subcommand '%v' is invalid. Use one of valid subcommands: %s", cmd.Args().First(), getSubcommandsNames(cmd)))
}

func getSubcommandsNames(cmd *cli.Command) []string {
	var names []string
	for _, c := range cmd.Commands {
		names = append(names, c.Name)
	}
	return names
}

func showHelpCommandWithErr(ctx context.Context, cmd *cli.Command, stringErr string) error {
	_, err := fmt.Fprintf(cmd.ErrWriter, "%s\n\n", color.RedString("%s", stringErr))
	if err != nil {
		return err
	}

	parentCommand := cmd
	if lineage := cmd.Lineage(); len(lineage) > 1 {
		parentCommand = lineage[1]
	}
	return cli.ShowCommandHelp(ctx, parentCommand, cmd.Name)
}

func validateVersion(_ context.Context, cmd *cli.Command) error {
	if !cmd.IsSet("version") {
		return nil
	}
	if cmd.Int64("version") <= 0 {
		return cli.Exit(color.RedString("CA set version should be a positive integer"), 1)
	}
	return nil
}

func validateRuleFormat(isProperty bool) actionValidator {
	return func(ctx context.Context, cmd *cli.Command) error {
		if !cmd.IsSet("rule-format") {
			return nil
		}

		sess := edgegrid.GetSession(ctx)
		client := papi.Client(sess)

		return checkRuleFormat(ctx, cmd, client, isProperty)
	}
}

func checkRuleFormat(ctx context.Context, cmd *cli.Command, client papi.PAPI, isProperty bool) error {
	ruleFormats, err := client.GetRuleFormats(ctx)
	if err != nil {
		return cli.Exit(color.RedString("Failed to fetch rule formats: %s", err), 1)
	}

	ruleFormat := cmd.String("rule-format")

	// "latest" is only allowed for properties using json format for rules
	latestAllowed := isProperty && !cmd.IsSet("rules-as-hcl")
	if !slices.Contains(ruleFormats.RuleFormats.Items, ruleFormat) || (ruleFormat == "latest" && !latestAllowed) {
		if isProperty {
			return cli.Exit(color.RedString("Invalid rule format version: %s. Must be vYYYY-MM-DD (with a leading \"v\"). The 'latest' value is allowed only for properties using the JSON rule format.", ruleFormat), 1)
		}
		return cli.Exit(color.RedString("Invalid rule format version: %s. Must be vYYYY-MM-DD (with a leading \"v\").", ruleFormat), 1)
	}

	return nil
}
