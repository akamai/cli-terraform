# Akamai CLI for Akamai Terraform Provider

[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform)](https://goreportcard.com/report/github.com/akamai/cli-terraform) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform?ref=badge_shield)

An [Akamai CLI](https://developer.akamai.com/cli) package for administering and managing Akamai Terraform configurations

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install terraform
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-terraform/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compiling from Source

If you want to compile it from source, you will need Go 1.12 or later:

1. Fetch the package:  
  `go get github.com/akamai/cli-terraform`
2. Change to the package directory:  
  `cd $GOPATH/src/github.com/akamai/cli-terraform`
3. Compile the binary:
  - Linux/macOS/*nix: `go build -o akamai-terraform`
  - Windows: `go build -o akamai-terraform.exe`
4. Move the binary (`akamai-terraform` or `akamai-terraform.exe`) in to your `PATH`

## General Usage

```
  akamai-terraform [--edgerc] [--section] <command> [sub-command]

Description:
   Manage Akamai Terraform configurations and assoc objects. Current support includes Akamai GTM domains and EdgeDNS zones.

Global Flags:
   --edgerc value  Location of the credentials file (default: "/home/elynes/.edgerc") [$AKAMAI_EDGERC]
   --section value     Section of the credentials file (default: "terraform") [$AKAMAI_EDGERC_SECTION]
   --accountkey value  Account switch key [$AKAMAI_EDGERC_ACCOUNT_KEY]

Built-In Commands:
  create-domain
  create-zone
  create-property
  create-cloudlets-policy
  create-edgekv
  list
  help
```

## GTM Domains

### Usage

```
   akamai-terraform create-domain [domain] [--tfworkpath path] [--resources] [--createconfig] 

Flags: 
   --tfworkpath path       file path location for placement of created and/or modified artifacts. Default: current directory
   --resources             Create json formatted resource import list file, <domain>_resources.json. Used as input by createconfig.
   --createconfig          Create Terraform configuration (<domain>.tf), gtmvars.tf, and import command script (<domain>_import.script) files using resources json
```

### Create list of all domain objects. Written in json format to <domain>_resources.json

```
$ akamai terraform create-domain example.akadns.net --resources
```

### Generate Terraform GTM Domain configuration file <domain>.tf, vars config file, gtmvars.tf, and import script, <domain>_resource_import.script

```
$ akamai terraform create-domain example.akadns.net --createconfig
```

### Domain Notes:
1. Mapping GTM entity names to TF resource names may require normalization. Invalid TF resource name characters will be replaced by underscores, '_' in config generation.
 

## EdgeDNS Zones

### Usage

```
   akamai-terraform create-zone [zone] [--tfworkpath path] [--resources] [--createconfig] [--importscript] [--segmentconfig] [--configonly] [--namesonly] [--recordname]

Flags: 
   --tfworkpath path       file path location for placement of created and/or modified artifacts. Default: current directory
   --resources             Create json formatted resource import list file, <zone>_resources.json. Used as input by createconfig.
   --createconfig          Create Terraform configuration (<zone>.tf), dnsvars.tf from generated resources file. Saves zone config for import.
   --importscript          Create import script for generated Terraform configuration script (<zone>_import.script) files
   --segmentconfig         Directive for createconfig. Group and segment records by name into separate config files.
   --configonly            Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for 
                           importscript. Ignores any existing resource json file.
   --namesonly             Directive for both resource gathering and config generation. All record set types assumed.
   --recordname value      Used in resources gathering or with configonly to filter recordsets. Multiple recordname flags may be specified.
```

### Create List of Zone Recordsets. Written in json format to <zone>_resources.json

```
$ akamai terraform create-zone testprimaryzone.com --resources
```

### Generate Terraform Zone configuration file. Default args create <zone>.tf, vars config file, dnsvars.tf

```
$ akamai terraform create-zone testprimaryzone.com --createconfig
```

### Generate Zone import script, <zone>_resource_import.script

```
$ akamai terraform create-zone testprimaryzone.com --importscript
```


### Zone Notes

1. The resources directive generates a <zone>_resources.json file for consumption by createconfig
2. The createconfig directive generates a <zone>_zoneconfig.json file for consumption by importscript

####  Advanced options for --resources

1. recordname - filters generated resources list by record name(s)
2. namesonly - Generates resource file with recordset names only. All associated Types will be represented.

#### Advanced options for --createconfig

1. namesonly - Resources for all associated Types will be generated
2. segmentconfig - Generate a modularized configuration. 
3. configonly. Generate zone configuration directly without json itemization. Scope limited by additional specified flags.

## Property Manager Properties

### Usage

```
   akamai-terraform create-property [property name] [--tfworkpath path] 

Flags:
   --tfworkpath path      file path location for placement of created and/or modified artifacts. Default: current directory
```

### Create property manager property configuration.

```
$ akamai terraform create-property
```

## Cloudlets

### Usage

```
   akamai-terraform create-cloudlets-policy [policy name] [--tfworkpath path] 

Flags:
   --tfworkpath path      path location for placement of created artifacts. Default: current directory
```

### Create policy configuration.

```
$ akamai terraform create-cloudlets-policy
```

## Edgeworkers

### Usage

```
   akamai-terraform create-edgekv [namespace_name] [network] [--tfworkpath path] 

Flags:
   --tfworkpath path      path location for placement of created artifacts. Default: current directory
```

### Create edgekv configuration.

```
$ akamai terraform create-edgekv
```

## General Notes
1. Terraform variable configuration is generated in a separately named TF file for each Akamai entity type. These files will need to be merged by the Admin in the case where multiple entities are managed concurrently with the Terraform client.

## License

This package is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform?ref=badge_large)
