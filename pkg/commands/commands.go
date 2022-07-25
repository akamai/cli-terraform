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
	"github.com/akamai/cli-terraform/pkg/providers/appsec"
	"github.com/akamai/cli-terraform/pkg/providers/cloudlets"
	"github.com/akamai/cli-terraform/pkg/providers/dns"
	"github.com/akamai/cli-terraform/pkg/providers/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/providers/gtm"
	"github.com/akamai/cli-terraform/pkg/providers/iam"
	"github.com/akamai/cli-terraform/pkg/providers/imaging"
	"github.com/akamai/cli-terraform/pkg/providers/papi"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/apphelp"
	"github.com/akamai/cli/pkg/autocomplete"
	"github.com/urfave/cli/v2"
)

// CommandLocator creates and returns a list of subcommands
func CommandLocator() ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:        "export-domain",
		Aliases:     []string{"create-domain"},
		Description: "Generates Terraform configuration for Domain resources",
		Usage:       "export-domain",
		ArgsUsage:   "<domain>",
		Action:      validatedAction(gtm.CmdCreateDomain, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-zone",
		Aliases:     []string{"create-zone"},
		Description: "Generates Terraform configuration for Zone resources",
		Usage:       "export-zone",
		ArgsUsage:   "<zone>",
		Action:      validatedAction(dns.CmdCreateZone, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.BoolFlag{
				Name:  "resources",
				Usage: "Create json formatted resource import list file, <zone>_resources.json. Used as input by createconfig.",
			},
			&cli.BoolFlag{
				Name:  "createconfig",
				Usage: "Create Terraform configuration (<zone>.tf), dnsvars.tf from generated resources file. Saves zone config for import.",
			},
			&cli.BoolFlag{
				Name:  "importscript",
				Usage: "Create import script for generated Terraform configuration script (<zone>_import.script) files",
			},
			&cli.BoolFlag{
				Name:  "segmentconfig",
				Usage: "Directive for createconfig. Group and segment records by name into separate config files.",
			},
			&cli.BoolFlag{
				Name:  "configonly",
				Usage: "Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for importscript. Ignores any existing resource json file.",
			},
			&cli.BoolFlag{
				Name:  "namesonly",
				Usage: "Directive for both resource gathering and config generation. All record set types assumed.",
			},
			&cli.StringSliceFlag{
				Name:  "recordname",
				Usage: "Used in resources gathering or with configonly to filter recordsets. Multiple recordname flags may be specified.",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-appsec",
		Aliases:     []string{"create-appsec"},
		Description: "Generates Terraform configuration for Application Security resources",
		Usage:       "export-appsec",
		ArgsUsage:   "<security configuration name>",
		Action:      validatedAction(appsec.CmdCreateAppsec, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-property",
		Aliases:     []string{"create-property"},
		Description: "Generates Terraform configuration for Property resources",
		Usage:       "export-property",
		ArgsUsage:   "<property name>",
		Action:      validatedAction(papi.CmdCreateProperty, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "Property version to import",
				DefaultText: "LATEST",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-cloudlets-policy",
		Aliases:     []string{"create-cloudlets-policy"},
		Description: "Generates Terraform configuration for Cloudlets Policy resources",
		Usage:       "export-cloudlets-policy",
		ArgsUsage:   "<policy_name>",
		Action:      validatedAction(cloudlets.CmdCreatePolicy, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-edgekv",
		Aliases:     []string{"create-edgekv"},
		Description: "Generates Terraform configuration for EdgeKV resources",
		Usage:       "export-edgekv",
		ArgsUsage:   "<namespace_name> <network>",
		Action:      validatedAction(edgeworkers.CmdCreateEdgeKV, requireValidWorkpath, requireNArguments(2)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-edgeworker",
		Aliases:     []string{"create-edgeworker"},
		Description: "Generates Terraform configuration for EdgeWorker resources",
		Usage:       "export-edgeworker",
		ArgsUsage:   "<edgeworker_id>",
		Action:      validatedAction(edgeworkers.CmdCreateEdgeWorker, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "bundlepath",
				Usage: "Path location for placement of EdgeWorkers tgz code bundle. Default: same value as tfworkpath",
			},
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:            "export-iam",
		Aliases:         []string{"create-iam"},
		Description:     "Generates Terraform configuration for Identity and Access Management resources",
		Usage:           "export-iam",
		HideHelpCommand: true,
		Action:          validatedAction(iam.CmdCreateIAM, requireValidWorkpath, validateSubCommands),
		Subcommands: []*cli.Command{
			{
				Name:        "all",
				Description: "Create all available Terraform Users, Groups and Roles",
				Action:      validatedAction(iam.CmdCreateIAMAll, requireValidWorkpath),
			},
			{
				Name:        "group",
				Description: "Create Terraform Group resource with relevant users and roles resources",
				ArgsUsage:   "<group_id>",
				Action:      validatedAction(iam.CmdCreateIAMGroup, requireValidWorkpath, requireNArguments(1)),
			},
			{
				Name:        "role",
				Description: "Create Terraform Role resource with relevant users and groups resources",
				ArgsUsage:   "<role_id>",
				Action:      validatedAction(iam.CmdCreateIAMRole, requireValidWorkpath, requireNArguments(1)),
			},
			{
				Name:        "user",
				Description: "Create Terraform User resource with relevant groups and roles resources",
				ArgsUsage:   "<user_email>",
				Action:      validatedAction(iam.CmdCreateIAMUser, requireValidWorkpath, requireNArguments(1)),
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-imaging",
		Aliases:     []string{"create-imaging"},
		Description: "Generates Terraform configuration for Image and Video Manager resources",
		Usage:       "export-imaging",
		ArgsUsage:   "<contract_id> <policy_set_id>",
		Action:      validatedAction(imaging.CmdCreateImaging, requireValidWorkpath, requireNArguments(2)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "policy-json-dir",
				Usage: "Path location for placement of policy jsons. Default: same value as tfworkpath",
			},
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.BoolFlag{
				Name:        "schema",
				Usage:       "Generate content of the policy using HCL instead of JSON file",
				Destination: &tools.Schema,
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:               "list",
		Description:        "List commands",
		Action:             cmdList,
		CustomHelpTemplate: apphelp.SimplifiedHelpTemplate,
	})

	return commands, nil
}
