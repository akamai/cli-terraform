// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	//"errors"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/urfave/cli"
	//"strconv"
	//"strings"
)

var commandLocator akamai.CommandLocator = func() ([]cli.Command, error) {
	var commands []cli.Command

	commands = append(commands, cli.Command{
		Name:        "create-domain",
		Description: "Create Terraform Domain Resources",
		ArgsUsage:   "[domain]",
		Action:      cmdCreateDomain,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "statepath",
				Usage: "Full path to terraform.tfstate file. Default: current directory",
			},
                        cli.StringFlag{
                                Name:  "configpath",
                                Usage: "Path to create config file, <domain>.tf. Default: current directory",
                        },
                        cli.StringFlag{
                                Name:  "importlistpath",
                                Usage: "Path to output list of resoures to be imported, <domain>_import.list. Default: current directory",
                        },
                        cli.StringFlag{
                                Name:  "tfworkpath",
                                Usage: "Path to output list of resoures to be imported, <domain>_import.list. Default: current directory",
                        },
                        cli.StringSliceFlag{
                                Name:  "property",
                                Usage: "Specific property and dependent resources to import. Multiple proeprty flags may be specified",
                        },
			cli.BoolTFlag{
				Name:  "importlist",
				Usage: "Create import list",
			},
			cli.BoolFlag{
				Name:  "createconfig",
				Usage: "Create terrform configuration and import commnad script",
			},
                        cli.BoolFlag{
                                Name:  "importconfig",
                                Usage: "Import terrform configuration",
                        },
/*
			cli.BoolFlag{
				Name:  "verbose",
				Usage: "Display verbose result status.",
			},
			cli.BoolFlag{
				Name:  "json",
				Usage: "Return status in JSON format.",
			},
*/
			cli.BoolFlag{
				Name:  "complete",
				Usage: "Wait up to 5 minutes for change completion",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands,
		cli.Command{
			Name:        "list",
			Description: "List commands",
			Action:      akamai.CmdList,
		},
		cli.Command{
			Name:         "help",
			Description:  "Displays help information",
			ArgsUsage:    "[command] [sub-command]",
			Action:       akamai.CmdHelp,
			BashComplete: akamai.DefaultAutoComplete,
		},
	)

	return commands, nil
}
