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

package commands

import (
	cloudlets "github.com/akamai/cli-terraform/pkg/providers/cloudlets"
	"github.com/akamai/cli-terraform/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/providers/dns"
	"github.com/akamai/cli-terraform/pkg/providers/gtm"
	"github.com/akamai/cli-terraform/pkg/providers/papi"
	akacli "github.com/akamai/cli/pkg/app"
	"github.com/urfave/cli/v2"
)

func CommandLocator() ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:        "create-domain",
		Description: "Create Terraform Domain Resources",
		ArgsUsage:   "<domain>",
		Action:      gtm.CmdCreateDomain,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created and/or modified artifacts. Default: current directory",
			},
			&cli.BoolFlag{
				Name:  "resources",
				Value: true,
				Usage: "Create json formatted resource import list file, <domain>_resources.json. Used as input by createconfig.",
			},
			&cli.BoolFlag{
				Name:  "createconfig",
				Value: true,
				Usage: "Create Terraform configuration (<domain>.tf), gtmvars.tf, and import command script (<domain>_import.script) files using resources json",
			},
		},
		BashComplete: akacli.DefaultAutoComplete,
	})

	commands = append(commands, &cli.Command{
		Name:        "create-zone",
		Description: "Create Terraform Zone Resources",
		ArgsUsage:   "<zone>",
		Action:      dns.CmdCreateZone,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created and/or modified artifacts. Default: current directory",
			},
			&cli.BoolFlag{
				Name:  "resources",
				Value: true,
				Usage: "Create json formatted resource import list file, <zone>_resources.json. Used as input by createconfig.",
			},
			&cli.BoolFlag{
				Name:  "createconfig",
				Value: true,
				Usage: "Create Terraform configuration (<zone>.tf), dnsvars.tf from generated resources file. Saves zone config for import.",
			},
			&cli.BoolFlag{
				Name:  "importscript",
				Value: true,
				Usage: "Create import script for generated Terraform configuration script (<zone>_import.script) files",
			},
			&cli.BoolFlag{
				Name:  "segmentconfig",
				Value: true,
				Usage: "Directive for createconfig. Group and segment records by name into separate config files.",
			},
			&cli.BoolFlag{
				Name:  "configonly",
				Value: true,
				Usage: "Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for importscript. Ignores any existing resource json file.",
			},
			&cli.BoolFlag{
				Name:  "namesonly",
				Value: true,
				Usage: "Directive for both resource gathering and config generation. All record set types assumed.",
			},
			&cli.StringSliceFlag{
				Name:  "recordname",
				Usage: "Used in resources gathering or with configonly to filter recordsets. Multiple recordname flags may be specified.",
			},
		},
		BashComplete: akacli.DefaultAutoComplete,
	})

	commands = append(commands, &cli.Command{
		Name:        "create-property",
		Description: "Create Terraform Property Resource",
		Usage:       "create-property",
		ArgsUsage:   "<property name>",
		Action:      papi.CmdCreateProperty,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created artifacts. Default: current directory",
			},
		},
		BashComplete: akacli.DefaultAutoComplete,
	})

	commands = append(commands, &cli.Command{
		Name:        "create-cloudlets-policy",
		Description: "Create Terraform Cloudlets Policy Resource",
		Usage:       "create-cloudlets-policy",
		ArgsUsage:   "<policy_name>",
		Action:      cloudlets.CmdCreatePolicy,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created artifacts. Default: current directory",
			},
		},
		BashComplete: akacli.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-edgekv",
		Description: "Create Terraform EdgeKV Resource",
		Usage:       "create-edgekv",
		ArgsUsage:   "<namespace_name>, <network>",
		Action:      edgeworkers.CmdCreateEdgeKV,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created artifacts. Default: current directory",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-edgeworker",
		Description: "Create Terraform EdgeWorker Resource",
		Usage:       "create-edgeworker",
		ArgsUsage:   "<edgeworker_id>",
		Action:      edgeworkers.CmdCreateEdgeWorker,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "bundlepath",
				Usage: "Path location for placement of EdgeWorkers tgz code bundle. Default: same value as tfworkpath",
			},
			cli.StringFlag{
				Name:  "tfworkpath",
				Usage: "Path location for placement of created artifacts. Default: current directory",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands,
		&cli.Command{
			Name:        "list",
			Description: "List commands",
			Action:      cmdList,
		},
		&cli.Command{
			Name:         "help",
			Description:  "Displays help information",
			ArgsUsage:    "<command> <sub-command>",
			Action:       cmdHelp,
			BashComplete: akacli.DefaultAutoComplete,
		},
	)

	return commands, nil
}
