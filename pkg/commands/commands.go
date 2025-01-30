// Package commands contains code defining the export commands.
package commands

import (
	"github.com/akamai/cli-terraform/pkg/providers/appsec"
	"github.com/akamai/cli-terraform/pkg/providers/clientlists"
	"github.com/akamai/cli-terraform/pkg/providers/cloudaccess"
	"github.com/akamai/cli-terraform/pkg/providers/cloudlets"
	"github.com/akamai/cli-terraform/pkg/providers/cloudwrapper"
	"github.com/akamai/cli-terraform/pkg/providers/cps"
	"github.com/akamai/cli-terraform/pkg/providers/dns"
	"github.com/akamai/cli-terraform/pkg/providers/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/providers/gtm"
	"github.com/akamai/cli-terraform/pkg/providers/iam"
	"github.com/akamai/cli-terraform/pkg/providers/imaging"
	"github.com/akamai/cli-terraform/pkg/providers/papi"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/v2/pkg/apphelp"
	"github.com/akamai/cli/v2/pkg/autocomplete"
	"github.com/urfave/cli/v2"
)

// CommandLocator creates and returns a list of subcommands
func CommandLocator() ([]*cli.Command, error) {
	var commands []*cli.Command

	commands = append(commands, &cli.Command{
		Name:        "export-domain",
		Description: "Generates Terraform configuration for Domain resources.",
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
		Description: "Generates Terraform configuration for Zone resources.",
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
				Usage: "Creates a JSON-formatted resource import list file, '<zone>_resources.json'. Used as input by the '--createconfig' flag.",
			},
			&cli.BoolFlag{
				Name:  "createconfig",
				Usage: "Creates these Terraform configuration files based on the values in '<zone>_resources.json': '<zone>.tf' and 'dnsvars.tf'. Saves the zone config for import.",
			},
			&cli.BoolFlag{
				Name:  "importscript",
				Usage: "Creates an import script for the generated Terraform configuration script files ('<zone>_import.script').",
			},
			&cli.BoolFlag{
				Name:  "segmentconfig",
				Usage: "Directive for the '--createconfig' flag. Groups and segments records by name into separate config files.",
			},
			&cli.BoolFlag{
				Name:  "configonly",
				Usage: "Directive for the '--createconfig' flag. Creates the entire Terraform zone and recordsets configuration ('<zone>.tf') and 'dnsvars.tf'. Saves the zone config for the '--importscript' flag. Ignores any existing resource json file.",
			},
			&cli.BoolFlag{
				Name:  "namesonly",
				Usage: "Directive for both gathering resources and generating a config file. All record set types are assumed.",
			},
			&cli.StringSliceFlag{
				Name:  "recordname",
				Usage: "Used when gathering resources or with the '--configonly' flag to filter recordsets. You can provide the '--recordname' flag multiple times.",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-appsec",
		Description: "Generates Terraform configuration for Application Security resources.",
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
		Name:        "export-clientlist",
		Description: "Generates Terraform configuration for Client List resources.",
		Usage:       "export-clientlist",
		ArgsUsage:   "<list_id>",
		Action:      validatedAction(clientlists.CmdCreateClientList, requireValidWorkpath, requireNArguments(1)),
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
		Description: "Generates Terraform configuration for Property resources.",
		Usage:       "export-property",
		ArgsUsage:   "<property name>",
		Action:      validatedAction(papi.CmdCreateProperty, requireValidWorkpath, requireNArguments(1), validateSplitDepth),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "Property version to import.",
				DefaultText: "LATEST",
			},
			&cli.BoolFlag{
				Name:  "rules-as-hcl",
				Usage: "Exports the referenced rules as the 'akamai_property_rules_builder' data source.",
			},
			&cli.BoolFlag{
				Name:  "akamai-property-bootstrap",
				Usage: "Exports the referenced property using a combination of the 'akamai-property-bootstrap' and 'akamai-property' resources.",
			},
			&cli.IntFlag{
				Name:  "split-depth",
				Usage: "Exports the rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level.",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-property-include",
		Description: "Generates Terraform configuration for Include resources.",
		Usage:       "export-property-include",
		ArgsUsage:   "<contract_id> <include_name>",
		Action:      validatedAction(papi.CmdCreateInclude, requireValidWorkpath, requireNArguments(2), validateSplitDepth),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.BoolFlag{
				Name:  "rules-as-hcl",
				Usage: "Exports the referenced rules as the 'akamai_property_rules_builder' data source.",
			},
			&cli.IntFlag{
				Name:  "split-depth",
				Usage: "Exports the rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level.",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-cloudwrapper",
		Description: "Generates Terraform configuration for CloudWrapper resources.",
		Usage:       "export-cloudwrapper",
		ArgsUsage:   "<config_id>",
		Action:      validatedAction(cloudwrapper.CmdCreateCloudWrapper, requireValidWorkpath, requireNArguments(1)),
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
		Name:        "export-cloudlets-policy",
		Description: "Generates Terraform configuration for Cloudlets Policy resources.",
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
		Description: "Generates Terraform configuration for EdgeKV resources.",
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
		Description: "Generates Terraform configuration for EdgeWorker resources.",
		Usage:       "export-edgeworker",
		ArgsUsage:   "<edgeworker_id>",
		Action:      validatedAction(edgeworkers.CmdCreateEdgeWorker, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "bundlepath",
				Usage: "Path location of the EdgeWorkers 'tgz' code bundle. Its default value is the same as for the '--tfworkpath' flag.",
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
		Description:     "Generates Terraform configuration for Identity and Access Management resources.",
		Usage:           "export-iam",
		HideHelpCommand: true,
		Action:          validatedAction(iam.CmdCreateIAM, requireValidWorkpath, validateSubCommands),
		Subcommands: []*cli.Command{
			{
				Name:        "all",
				Description: "Exports all available Terraform users, groups, roles, and allowlist details.",
				Action:      validatedAction(iam.CmdCreateIAMAll, requireValidWorkpath),
			},
			{
				Name:        "allowlist",
				Description: "Exports Terraform IP Allowlist and CIDR block resources.",
				Action:      validatedAction(iam.CmdCreateIAMAllowlist, requireValidWorkpath),
			},
			{
				Name:        "group",
				Description: "Exports the Terraform group resource with relevant user and role resources.",
				ArgsUsage:   "<group_id>",
				Action:      validatedAction(iam.CmdCreateIAMGroup, requireValidWorkpath, requireNArguments(1)),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "only",
						Usage: "Exports only the Terraform group resource; excludes the role and user resources when specified.",
					},
				},
			},
			{
				Name:        "role",
				Description: "Exports the Terraform role resource with relevant user and group resources.",
				ArgsUsage:   "<role_id>",
				Action:      validatedAction(iam.CmdCreateIAMRole, requireValidWorkpath, requireNArguments(1)),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "only",
						Usage: "Exports only the Terraform role resource; excludes the user and group resources when specified.",
					},
				},
			},
			{
				Name:        "user",
				Description: "Exports the Terraform user resource with relevant group and role resources.",
				ArgsUsage:   "<user_email>",
				Action:      validatedAction(iam.CmdCreateIAMUser, requireValidWorkpath, requireNArguments(1)),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "only",
						Usage: "Exports only the Terraform user resource; excludes the group and role resources when specified.",
					},
				},
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
		Description: "Generates Terraform configuration for Image and Video Manager resources.",
		Usage:       "export-imaging",
		ArgsUsage:   "<contract_id> <policy_set_id>",
		Action:      validatedAction(imaging.CmdCreateImaging, requireValidWorkpath, requireNArguments(2)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "policy-json-dir",
				Usage: "Path location for a policy in JSON format. Its default value is the same as for the '--tfworkpath' flag.",
			},
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.BoolFlag{
				Name:        "policy-as-hcl",
				Usage:       "Generates content of the policy using HCL format instead of JSON.",
				Destination: &tools.PolicyAsHCL,
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:        "export-cps",
		Description: "Generates Terraform configuration for CPS (Certificate Provisioning System) resources.",
		Usage:       "export-cps",
		ArgsUsage:   "<enrollment_id> <contract_id>",
		Action:      validatedAction(cps.CmdCreateCPS, requireValidWorkpath, requireNArguments(2)),
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
		Name:        "export-cloudaccess",
		Description: "Generates Terraform configuration for CAM (Cloud Access Manager) resources.",
		Usage:       "export-cloudaccess",
		ArgsUsage:   "<access_key_uid>",
		Action:      validatedAction(cloudaccess.CmdCreateCloudAccess, requireValidWorkpath, requireNArguments(1)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "tfworkpath",
				Usage:       "Directory used to store files created when running commands.",
				DefaultText: "current directory",
			},
			&cli.StringFlag{
				Name:  "group_id",
				Usage: "The unique identifier for the group (without the 'grp_' prefix) assigned to the access key.",
			},
			&cli.StringFlag{
				Name:  "contract_id",
				Usage: "The unique identifier for the contract (without the 'ctr_' prefix) assigned to the access key.",
			},
		},
		BashComplete: autocomplete.Default,
	})

	commands = append(commands, &cli.Command{
		Name:               "list",
		Description:        "List commands.",
		Action:             cmdList,
		CustomHelpTemplate: apphelp.SimplifiedHelpTemplate,
	})

	return commands, nil
}
