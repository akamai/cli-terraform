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
				Name:  "tfworkpath",
				Usage: "file path location for placement of created and/or modified artifacts. Default: current directory",
			},
			cli.BoolTFlag{
				Name:  "resources",
				Usage: "Create json formatted resource import list file, <domain>_resources.json. Used as input by createconfig.",
			},
			cli.BoolTFlag{
				Name:  "createconfig",
				Usage: "Create Terraform configuration (<domain>.tf), gtmvars.tf, and import command script (<domain>_import.script) files using resources json",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-zone",
		Description: "Create Terraform Zone Resources",
		ArgsUsage:   "[zone]",
		Action:      cmdCreateZone,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "file path location for placement of created and/or modified artifacts. Default: current directory",
			},
			cli.BoolTFlag{
				Name:  "resources",
				Usage: "Create json formatted resource import list file, <zone>_resources.json. Used as input by createconfig.",
			},
			cli.BoolTFlag{
				Name:  "createconfig",
				Usage: "Create Terraform configuration (<zone>.tf), dnsvars.tf from generated resources file. Saves zone config for import.",
			},
			cli.BoolTFlag{
				Name:  "importscript",
				Usage: "Create import script for generated Terraform configuration script (<zone>_import.script) files",
			},
			cli.BoolTFlag{
				Name:  "segmentconfig",
				Usage: "Directive for createconfig. Group and segment records by name into separate config files.",
			},
			cli.BoolTFlag{
				Name:  "configonly",
				Usage: "Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for importscript. Ignores any existing resource json file.",
			},
			cli.BoolTFlag{
				Name:  "namesonly",
				Usage: "Directive for both resource gathering and config generation. All record set types assumed.",
			},
			cli.StringSliceFlag{
				Name:  "recordname",
				Usage: "Used in resources gathering or with configonly to filter recordsets. Multiple recordname flags may be specified.",
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
