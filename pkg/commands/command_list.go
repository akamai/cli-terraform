package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func cmdList(c *cli.Context) error {
	color.Yellow("\nCommands:\n\n")

	for _, command := range c.App.Commands {
		bold := color.New(color.FgWhite, color.Bold)
		fmt.Print(bold.Sprintf("  %s", command.Name))
		if len(command.Aliases) > 0 {
			var aliases string

			if len(command.Aliases) == 1 {
				aliases = "alias"
			} else {
				aliases = "aliases"
			}

			fmt.Printf(" (%s: ", aliases)
			for i, alias := range command.Aliases {
				fmt.Print(bold.Sprint(alias))
				if i < len(command.Aliases)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Print(")")
		}

		fmt.Println()

		if command.Description != "" {
			fmt.Printf("    %s\n", command.Description)
		}
	}

	fmt.Printf("\nSee \"%s\" for details.\n", color.BlueString("%s help [command]", c.App.Name))
	return nil
}
