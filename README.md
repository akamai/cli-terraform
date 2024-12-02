# Akamai CLI for Akamai Terraform Provider

![Build Status](https://github.com/akamai/cli-terraform/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform)](https://goreportcard.com/report/github.com/akamai/cli-terraform)
![GitHub release](https://img.shields.io/github/v/release/akamai/cli-terraform)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/cli-terraform?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/cli-terraform)

## Get started

### Create authentication credentials

Before you can use this CLI, you need to [Create authentication credentials.](https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials)

### Install

To install this package, use Akamai CLI:

```shell
$ akamai install terraform
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-terraform/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compile from source

If you want to compile it from the source, you need Go 1.21 or later:

1. Create a clone of the target repository:
  `git clone https://github.com/akamai/cli-terraform.git`
1. Change to the package directory and compile the binary:
   - Linux/macOS/*nix: `go build -o akamai-terraform`
   - Windows: `go build -o akamai-terraform.exe`

## General usage

```shell
Usage:
  akamai terraform [global flags] command [command flags] [arguments...]

Description:
  Export selected resources for faster adoption in Terraform.

Commands:
  export-domain (alias: create-domain)
  export-zone (alias: create-zone)
  export-appsec (alias: create-appsec)
  export-clientlist
  export-property (alias: create-property)
  export-property-include
  export-cloudwrapper
  export-cloudlets-policy (alias: create-cloudlets-policy)
  export-edgekv (alias: create-edgekv)
  export-edgeworker (alias: create-edgeworker)
  export-iam (alias: create-iam)
  export-imaging (alias: create-imaging)
  export-cps (alias: create-cps)
  export-cloudaccess
  list
  help

Global Flags:
  --edgerc value, -e value                 Location of the credentials file (default: "/home/user/.edgerc") [$AKAMAI_EDGERC]
  --section value, -s value                Section of the credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
  --accountkey value, --account-key value  Account switch key [$AKAMAI_EDGERC_ACCOUNT_KEY]
  --help                                   show help (default: false)
  --version                                Output CLI version (default: false)
```

## GTM Domains

### Usage

```shell
   akamai terraform [global flags] export-domain [flags] <domain>

Flags:
   --tfworkpath path       Directory used to store files created when running commands. (default: current directory)
```

### Export Terraform GTM domain configuration

```shell
$ akamai terraform export-domain example.akadns.net
```

### Domain notes

Mapping GTM entity names to TF resource names may require normalization. Invalid TF resource name characters are replaced by underscores, '_' in config generation.


## Edge DNS Zones

### Usage

```shell
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

### Export list of zone record sets

The command creates a `<zone>_resources.json` file.

```shell
$ akamai terraform export-zone --resources testprimaryzone.com
```

### Generate Terraform zone configuration file

Default arguments create these files: `<zone>.tf`, `<zone>_zoneconfig.json`, and `dnsvars.tf`.

```shell
$ akamai terraform export-zone --createconfig testprimaryzone.com
```

### Generate zone import script

The command creates a `<zone>_resource_import.script` file.

```shell
$ akamai terraform export-zone --importscript testprimaryzone.com
```


### Zone notes

1. The resources directive generates a `<zone>_resources.json` file for consumption by `createconfig`.
2. The `createconfig` directive generates a `<zone>_zoneconfig.json` file for consumption by `importscript`.

####  Advanced options for --resources

1. `recordname`. Filters the generated resources list by a record name or record names.
2. `namesonly`. Generates a resource file with record set names only. All associated Types are represented.

#### Advanced options for --createconfig

1. `namesonly`. Generates resources for all associated Types.
2. `segmentconfig`. Generates a modularized configuration.
3. `configonly`. Generates a zone configuration without the JSON itemization. The configuration generated varies based on which set of flags you use.

## AppSec

### Usage

```shell
   akamai terraform [global flags] export-appsec [flags] <name_of_security_config>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

## Property Manager Properties

Certain export conditions require the use of a particular property rule format. Verify whether your rule format matches the use case requirement and [update your rule format](https://techdocs.akamai.com/terraform/docs/set-up-includes#update-rule-format) as needed.

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
    <td>Addition of the <code>--rules-as-hcl</code> flag</td>
    <td>Your declarative property configuration and HCL-formatted rules. <strong>Does not return includes</strong> as includes are JSON-formatted.</td>
    <td>Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</td>
  </tr>
  <tr>
    <td>Addition of the <code>include</code> subcommand</td>
    <td>Your property configuration and its JSON-formatted rules and includes.</td>
    <td>Any supported format, <em>but</em> your rules and includes must use the same rule format version.</td>
  </tr>
</tbody>
</table>

### Usage

```shell
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

> The `rules-as-hcl` flag works now with the `include` subcommand and the `with-includes` flag.

### Export property manager property configuration

```shell
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
    <td>Addition of the <code>--rules-as-hcl</code> flag</td>
    <td>Your declarative include configuration and HCL-formatted rules. <strong>Does not return includes</strong> as includes are JSON-formatted.</td>
    <td>Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</td>
  </tr>
</tbody>
</table>

### Usage

```shell
   akamai terraform [global flags] export-property-include [flags] <contract_id> <include_name>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
   --rules-as-hcl         Rules will be exported as `akamai_property_rules_builder` data source in HCL format.
```

### Export property manager include configuration

```shell
$ akamai terraform export-property-include
```

## Cloudlets

### Usage

```shell
   akamai terraform [global flags] export-cloudlets-policy [flags] <policy_name>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export Cloudlets policy configuration

```shell
$ akamai terraform export-cloudlets-policy
```

## CloudWrapper

### Usage

```shell
   akamai terraform [global flags] export-cloudwrapper [flags] <configuration_id>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export CloudWrapper configuration

```shell
$ akamai terraform export-cloudwrapper
```

## EdgeWorkers

### EdgeKV usage

```shell
   akamai terraform [global flags] export-edgekv [flags] <namespace_name> <network>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export EdgeKV configuration

```shell
$ akamai terraform [global flags] export-edgekv [flags] <namespace_name> <network>
```

### EdgeWorker usage

```shell
   akamai terraform [global flags] export-edgeworker [flags] <edgeworker_id>

Flags:
   --bundlepath path      Path location for placement of EdgeWorkers tgz code bundle. Default: same value as tfworkpath
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

### Export EdgeWorker configuration

```shell
$ akamai terraform [global flags] export-edgeworker [flags] <edgeworker_id>
```

## Identity and Access Management

### Usage

```
   akamai terraform [global flags] export-iam [subcommand] [flags]

Subcommands:
    all                     Exports all available Terraform Users, Groups, Roles, IP Allowlist and CIDR block resources
    allowlist               Exports Terraform IP Allowlist and CIDR block resources
    group [group id]        Exports group by id with relevant users and their roles
    role [role id]          Exports role by id with relevant users and their groups
    user [user's email]     Exports user by email with relevant user's groups and roles

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
   --only                 Exports only the specified Identity and Access Management resource, excluding additional details when present. (Can only be used with `group`, `role`, or `user` subcommands.)
```

### Export Identity and Access Management configuration

```shell
$ akamai terraform export-iam
```

## Image and Video Manager

### Usage

```shell
   akamai terraform [global flags] export-imaging [flags] <contract_id> <policy_set_id>

Flags:
   --tfworkpath path         Directory used to store files created when running commands. (default: current directory)
   --policy-json-dir path    Path location for placement of policy jsons. Default: same value as tfworkpath
   --policy-as-hcl           Generate content of the policy using HCL instead of JSON file (default: false)
```

### Export Image and Video policy configuration

```shell
$ akamai terraform export-imaging
```

## Certificate Provisioning System (CPS)

### Usage

```shell
   akamai terraform [global flags] export-cps [flags] <enrollment_id> <contract_id>

Flags:
   --tfworkpath path                        Directory used to store files created when running commands. (default: current directory)
```

### Export CPS configuration

```shell
$ akamai terraform export-cps
```

## Client Lists

### Usage

```shell
akamai terraform [global flags] export-clientlist [flags] <list_id>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
```

## Cloud Access Manager

### Usage

```shell
   akamai terraform [global flags] export-cloudaccess [flags] <cloud_access_key_uid>

Flags:
   --tfworkpath path         Directory used to store files created when running commands. (default: current directory)
   --group_id                The unique identifier for the group (without the `grp_` prefix) assigned to the access key.
   --contract_id             The unique identifier for the contract (without the `ctr_` prefix) assigned to the access key.
```

## General notes

Terraform variable configuration is generated in a separately named TF file for each Akamai entity type. These files will need to be merged by the Admin in the case where multiple entities are managed concurrently with the Terraform client.

## License

Copyright 2024 Akamai Technologies, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use these files except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.