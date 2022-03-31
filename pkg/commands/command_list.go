/*
 * Copyright 2018. Akamai Technologies, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
