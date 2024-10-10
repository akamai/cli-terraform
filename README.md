# Akamai CLI for Akamai Terraform Provider

![Build Status](https://github.com/akamai/cli-terraform/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform)](https://goreportcard.com/report/github.com/akamai/cli-terraform)
![GitHub release](https://img.shields.io/github/v/release/akamai/cli-terraform)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/cli-terraform?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/cli-terraform)

## Getting Started

### Creating authentication credentials
Before you can use this CLI, you need to [Create authentication credentials.](https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials)

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install terraform
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-terraform/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compiling from Source

If you want to compile it from source, you will need Go 1.21 or later:

1. Create a clone of the target repository:
  `git clone https://github.com/akamai/cli-terraform.git`
2. Change to the package directory and compile the binary:
  - Linux/macOS/*nix: `go build -o akamai-terraform`
  - Windows: `go build -o akamai-terraform.exe`

## General Usage

```
  akamai terraform [global flags] command [command flags] [arguments...]

Description:
   Administer and manage available Akamai resources with Terraform

Built-In Commands:
  export-cloudwrapper
  export-domain (alias: create-domain)
  export-zone (alias: create-zone)
  export-appsec (alias: create-appsec)
  export-property (alias: create-property)
  export-cloudlets-policy (alias: create-cloudlets-policy)
  export-edgekv (alias: create-edgekv)
  export-edgeworker (alias: create-edgeworker)
  export-iam (alias: create-iam)
  export-imaging (alias: create-imaging)
  export-clientlist (alias: create-clientlist)
  list
  help

Global Flags:
   --help                                   show help (default: false)
   --edgerc value, -e value                 Location of the credentials file (default: "/home/user/.edgerc") [$AKAMAI_EDGERC]
   --section value, -s value                Section of the credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --accountkey value, --account-key value  Account switch key [$AKAMAI_EDGERC_ACCOUNT_KEY]
   --version                                Output CLI version (default: false)
```

## GTM Domains

### Usage

```
   akamai terraform [global flags] export-domain [flags] <domain>

Flags:
   --tfworkpath path       Directory used to store files created when running commands. (default: current directory)
   --resources             Creates a JSON-formatted resource file for import: <domain>_resources.json. The createconfig flag uses this file as an input. (default: false)
   --createconfig          Creates these Terraform configuration files based on the values in <domain>_resources.json: <domain>.tf and gtmvars.tf. Also creates this import script: <domain>_import.script. (default: false)
```

### Export list of all domain objects. Written in json format to <domain>_resources.json

```
$ akamai terraform export-domain --resources example.akadns.net
```

### Generate Terraform GTM Domain configuration file <domain>.tf, vars config file, gtmvars.tf, and import script, <domain>_resource_import.script

```
$ akamai terraform export-domain --createconfig example.akadns.net
```

### Domain Notes:
1. Mapping GTM entity names to TF resource names may require normalization. Invalid TF resource name characters will be replaced by underscores, '_' in config generation.


## EdgeDNS Zones

### Usage

```
   akamai terraform [global flags] export-zone [flags] <zone>

Flags:
   --tfworkpath path       Directory used to store files created when running commands. (default: current directory)
   --resources             Creates a JSON-formatted resource file for import: <zone>_resources.json. The createconfig flag uses this file as an input. (default: false)
   --createconfig          Creates these Terraform configuration files based on the values in <zone>_resources.json: <zone>.tf and dnsvars.tf. (default: false)
   --importscript          Creates import script for generated Terraform configuration script (<zone>_import.script) files. (default: false)
   --segmentconfig         Use with the createconfig flag to group and segment records by name into separate config files. (default: false)
   --configonly            Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for
                           importscript. Ignores any existing resource JSON file. (default: false)
   --namesonly             Directive for both resource gathering and config generation. All record set types assumed. (default: false)
   --recordname value      Used in resources gathering or with configonly to filter recordsets. Multiple recordname flags may be specified.
```

### Export List of Zone Recordsets. Written in json format to <zone>_resources.json

```
$ akamai terraform export-zone --resources testprimaryzone.com
```

### Generate Terraform Zone configuration file. Default args create <zone>.tf, vars config file, dnsvars.tf

```
$ akamai terraform export-zone --createconfig testprimaryzone.com
```

### Generate Zone import script, <zone>_resource_import.script

```
$ akamai terraform export-zone --importscript testprimaryzone.com
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
3. configonly - Generates a zone configuration without JSON itemization. The configuration generated varies based on which set of flags you use.

## Appsec

### Usage
```
   akamai terraform [global flags] export-appsec [flags] <name_of_security_config>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

## Property Manager Properties

Certain export conditions require the use of a particular property rule format. Verify your rule format matches the use case requirement and [update your rule format](https://techdocs.akamai.com/terraform/docs/set-up-includes#update-rule-format) as needed.

<table>
<thead>
  <tr>
    <th>Export condition</th>
    <th>Output</th>
    <th>Rule format</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td>General</td>
    <td>Your declarative property configuration and its JSON-formatted rules.</td>
    <td>Any supported format.</td>
  </tr>
  <tr>
    <td>Addition of <code>--rules-as-hcl</code> flag</td>
    <td>Your declarative property configuration and HCL-formatted rules. <strong>Does not return includes</strong> as includes are JSON-formatted.</td>
    <td>Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</td>
  </tr>
  <tr>
    <td>Addition of <code>include</code> subcommand</td>
    <td>Your property configuration and its JSON-formatted rules and includes.</td>
    <td>Any supported format, <em>but</em> your rules and includes must use the same rule format version.</td>
  </tr>
</tbody>
</table>

### Usage

```
   akamai terraform [global flags] export-property [subcommand] [flags] <property name>

Subcommand:
    include <contract_id> <include_name>    Generates Terraform configuration for Include resources. Deprecated, use `export-property-include` instead.

Flags:
   --tfworkpath path             Directory used to store files created when running commands. (default: current directory)
   --version value               Property version to import  (default: LATEST)
   --with-includes               Referenced includes will also be exported along with property. Deprecated.
   --rules-as-hcl                Rules will be exported as `akamai_property_rules_builder` data source in HCL format.
   --akamai-property-bootstrap   Referenced property will be exported using combination of `akamai-property-bootstrap` and `akamai-property` resources (default: false)
```

> Flag `rules-as-hcl` works now with `include` sub-command as well with `with-includes` flag.

### Export property manager property configuration.

```
$ akamai terraform export-property
```

## Property Manager Includes

Certain export conditions require the use of a particular property rule format. Verify your rule format matches the use case requirement and [update your rule format](https://techdocs.akamai.com/terraform/docs/set-up-includes#update-rule-format) as needed.

<table>
<thead>
  <tr>
    <th>Export condition</th>
    <th>Output</th>
    <th>Rule format</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td>General</td>
    <td>Your declarative include configuration and its JSON-formatted rules.</td>
    <td>Any supported format.</td>
  </tr>
  <tr>
    <td>Addition of <code>--rules-as-hcl</code> flag</td>
    <td>Your declarative include configuration and HCL-formatted rules. <strong>Does not return includes</strong> as includes are JSON-formatted.</td>
    <td>Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</td>
  </tr>
</tbody>
</table>

### Usage

```
   akamai terraform [global flags] export-property-include [flags] <contract_id> <include_name>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
   --rules-as-hcl         Rules will be exported as `akamai_property_rules_builder` data source in HCL format.
```

### Export property manager include configuration.

```
$ akamai terraform export-property-include
```

## Cloudlets

### Usage

```
   akamai terraform [global flags] export-cloudlets-policy [flags] <policy_name>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export Cloudlets Policy configuration.

```
$ akamai terraform export-cloudlets-policy
```

## CloudWrapper

### Usage

```
   akamai terraform [global flags] export-cloudwrapper [flags] <configuration_id>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export CloudWrapper configuration.

```
$ akamai terraform export-cloudwrapper
```

## EdgeWorkers

### Export EdgeKV Usage

```
   akamai terraform [global flags] export-edgekv [flags] <namespace_name> <network>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export edgekv configuration.

```
$ akamai terraform [global flags] export-edgekv [command flags] <namespace_name> <network>
```

### Export EdgeWorker Usage

```
   akamai terraform [global flags] export-edgeworker [command flags] <edgeworker_id>

Flags:
   --bundlepath path      Path location for placement of EdgeWorkers tgz code bundle. Default: same value as tfworkpath
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export EdgeWorker configuration.

```
$ akamai terraform [global flags] export-edgeworker [command flags] <edgeworker_id>
```

## Identity and Access Management

### Export Identity and Access Management usage

```
   akamai terraform [global flags] export-iam [subcommand]

Subcommands:
    all                     Exports all available Terraform Users, Groups, Roles, IP Allowlist and CIDR block resources
    allowlist               Exports Terraform IP Allowlist and CIDR block resources
    group [group id]        Exports group by id with relevant users and their roles
    role [role id]          Exports role by id with relevant users and their groups
    user [user's email]     Exports user by email with relevant user's groups and roles

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export Identity and Access Management configuration.

```
$ akamai terraform export-iam
```

## Image and Video Manager

### Export Image and Video policy usage

```
   akamai terraform [global flags] export-imaging [flags] <contract_id> <policy_set_id>

Flags:
   --tfworkpath path         Directory used to store files created when running commands. (default: current directory)
   --policy-json-dir path    Path location for placement of policy jsons. Default: same value as tfworkpath
   --policy-as-hcl           Generate content of the policy using HCL instead of JSON file (default: false)
```

### Export Image and Video policy configuration.

```
$ akamai terraform export-imaging
```

## Certificate Provisioning System (CPS)

### Export CPS usage

```
   akamai terraform [global flags] export-cps [flags] <enrollment_id> <contract_id>

Flags:
   --tfworkpath path                        Directory used to store files created when running commands. (default: current directory)
```

### Export CPS configuration.

```
$ akamai terraform export-cps
```

## Client Lists

### Export Client List.

```
akamai terraform [global flags] export-clientlist [flags] <list_id>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

## Cloud Access Manager

### Export Cloud Access key configuration

```
   akamai terraform [global flags] export-cloudaccess [flags] <cloud_access_key_uid>

Flags:
   --tfworkpath path         Directory used to store files created when running commands. (default: current directory)
```

## General Notes

1. Terraform variable configuration is generated in a separately named TF file for each Akamai entity type. These files
   will need to be merged by the Admin in the case where multiple entities are managed concurrently with the Terraform
   client.

## License

This package is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform?ref=badge_large)
