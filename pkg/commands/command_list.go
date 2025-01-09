package commands

import (
	"fmt"

	"github.com/akamai/cli/pkg/color"
	"github.com/urfave/cli/v2"
)

func cmdList(c *cli.Context) error {
	fmt.Print(color.YellowString("\nCommands:\n\n"))

	for _, command := range c.App.Commands {
		fmt.Print(color.BoldString("  %s", command.Name))
		if len(command.Aliases) > 0 {
			var aliases string

			if len(command.Aliases) == 1 {
				aliases = "alias"
			} else {
				aliases = "aliases"
			}

			fmt.Printf(" (%s: ", aliases)
			for i, alias := range command.Aliases {
				fmt.Print(color.BoldString(alias))
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
