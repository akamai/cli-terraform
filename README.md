# Akamai CLI for Akamai Terraform Provider

![Build Status](https://github.com/akamai/cli-terraform/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform/v2)](https://goreportcard.com/report/github.com/akamai/cli-terraform/v2)
![GitHub release](https://img.shields.io/github/v/release/akamai/cli-terraform)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/cli-terraform?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/cli-terraform/v2)

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

If you want to compile it from the source, you need Go 1.22 or later:

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
  export-domain
  export-zone
  export-appsec
  export-clientlist
  export-property
  export-property-include
  export-cloudwrapper
  export-cloudlets-policy
  export-edgekv
  export-edgeworker
  export-iam
  export-imaging
  export-cps
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
   --resources             Creates a JSON-formatted resource file for import: `<zone>_resources.json`. The `--createconfig` flag uses this file as an input. (default: false)
   --createconfig          Creates these Terraform configuration files based on the values in `<zone>_resources.json`: `<zone>.tf` and `dnsvars.tf`. (default: false)
   --importscript          Creates an import script for the generated Terraform configuration script files (<zone>_import.script). (default: false)
   --segmentconfig         Use with the `--createconfig` flag to group and segment records by name into separate config files. (default: false)
   --configonly            Directive for createconfig. Create entire Terraform zone and recordsets configuration (<zone>.tf), dnsvars.tf. Saves zone config for
                           importscript. Ignores any existing resource JSON file. (default: false)
   --namesonly             Directive for both gathering resources and generating a config file. All record set types are assumed. (default: false)
   --recordname value      Used when gathering resources or with the `--configonly` to filter recordsets. You can specify the `--recordname` flag multiple times.
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
3. `configonly`. Generates a zone configuration without the JSON itemization. The configuration generated varies based on the set of flags you use.

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
    <td>Addition of the <code>--split-depth=X</code> flag</td>
    <td>Rules will be exported into a module. Each rule up to an <code>X</code> nesting level will be placed in dedicated file. 
For example, <code>--split-depth=1</code> means that the default/root rule and all its direct children will be placed in dedicated files. Rules with higher nesting levels will be placed in a file of their closest ancestor.</td>
    <td>Any supported format.</td>
  </tr>
</tbody>
</table>

### Usage

```shell
   akamai terraform [global flags] export-property [flags] <property name>

Flags:
   --tfworkpath path             Directory used to store files created when running commands. (default: current directory)
   --version value               Property version to import  (default: LATEST)
   --rules-as-hcl                Exports rules as the `akamai_property_rules_builder` data source in HCL format.
   --akamai-property-bootstrap   Exports the referenced property using a combination of `akamai-property-bootstrap` and `akamai-property` resources. (default: false)
   --split-depth value           Exports the rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level.
```

> You can use the `split-depth` flag only along with the `rules-as-hcl` flag.

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
  <tr>
    <td>Addition of the <code>--split-depth=X</code> flag</td>
    <td>Rules will be exported into a module. Each rule will be placed in a dedicated file up to a specified nesting level. 
For example, <code>--split-depth=1</code> means that the default/root rule and all its direct children will be placed in dedicated files. Rules with higher nesting levels will be placed in a file of their closest ancestor.</td>
    <td>Any supported format.</td>
  </tr>
</tbody>
</table>

### Usage

```shell
   akamai terraform [global flags] export-property-include [flags] <contract_id> <include_name>

Flags:
   --tfworkpath path      Directory used to store files created when running commands. (default: current directory)
   --rules-as-hcl         Exports rules as the `akamai_property_rules_builder` data source in HCL format.
   --split-depth value    Exports the rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level.
```

> You can use the `split-depth` flag only along with the `rules-as-hcl` flag.

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
   --bundlepath path      Path location of the EdgeWorkers `tgz` code bundle. Its default value is the same as for the `--tfworkpath` flag.
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
    all                     Exports all available Terraform users, groups, roles, IP allowlist and CIDR block resources.
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
   --policy-json-dir path    Path location for a policy in JSON format. Its default value is the same as for the `--tfworkpath` flag.
   --policy-as-hcl           Generates content of the policy using HCL format instead of JSON. (default: false)
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