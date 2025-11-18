# Akamai CLI for Akamai Terraform Provider

![Build Status](https://github.com/akamai/cli-terraform/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform/v2)](https://goreportcard.com/report/github.com/akamai/cli-terraform/v2)
![GitHub release](https://img.shields.io/github/v/release/akamai/cli-terraform)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://pkg.go.dev/badge/github.com/akamai/cli-terraform?utm_source=godoc)](https://pkg.go.dev/github.com/akamai/cli-terraform/v2)

This library provides a command-line interface to export Akamai configuration assets that you can import later into your Terraform state.

Requires Go 1.24 or later.

## Install

To install this package, you can use the Akamai CLI.

```shell
akamai install terraform
```

If you already have the package installed on your system, run `akamai update terraform` to update it.

You can also use the package as a stand-alone command by downloading the [latest release binary](https://github.com/akamai/cli-terraform/releases) or cloning this repository and compiling the binary yourself.

```shell
# Linux/macOS/*nix
go build -o akamai-terraform

# Windows
go build -o akamai-terraform.exe
```

## Authenticate

Authentication credentials for our Terraform provider use a hash-based message authentication code or HMAC-SHA-256 created through an API client. Each member of your team should use their own client set up locally to prevent accidental exposure of credentials.

There are different types of API clients that grant access based on your need, role, or how many accounts you manage.

| API client type| Description |
|---|---|
| [Basic](#create-a-basic-api-client)| Access to the first 99 API associated with your account without any specific configuration. Individual service read/write permissions are based on your role. |
| [Advanced](https://techdocs.akamai.com/developer/docs/create-a-client-with-custom-permissions) | Configurable permissions to limit or narrow down the scope of the API for your account.|
| [Managed](https://techdocs.akamai.com/developer/docs/manage-many-accounts-with-one-api-client)| Configurable permissions that work for multiple accounts.|

To set up and use multiple clients, clients that use an account switch key, or clients as environment variables, see [Alternative authentication](https://techdocs.akamai.com/terraform/docs/gs-authentication).

### Create a basic API client

1. Navigate to the [Identity and Access Management](https://control.akamai.com/apps/identity-management/#/tabs/users/list) section of Akamai Control Center and click **Create API Client**.

2. Click **Quick** and then **Download** in the **Credentials** section.

3. Open the downloaded file with a text editor and add `[default]` as a header above all text.

   ```bash Shell
   [default]
   client_secret = C113nt53KR3TN6N90yVuAgICxIRwsObLi0E67/N8eRN=
   host = akab-h05tnam3wl42son7nktnlnnx-kbob3i3v.luna.akamaiapis.net
   access_token = akab-acc35t0k3nodujqunph3w7hzp7-gtm6ij
   client_token = akab-c113ntt0k3n4qtari252bfxxbsl-yvsdj
   ```

4. Save your credentials as an EdgeGrid resource file named `.edgerc` in your local home directory.

## Use

1. To use the library, provide the path to your `.edgerc` file, your credentials section header, and the export command for the asset you need.
   
    If you manage multiple accounts, pass your account switch key using the global `--accountkey` flag. Use other additional flags or arguments to further define the output.

    > **Notes:**
    > - If you pass the command without the `--edgerc` and `--section` global flags, the command, by default, will point to the local home directory of your `.edgerc` file and the `default` credentials section header of that file.
    > - By default, when you run the command, generated files are saved in your active directory. Use the `--tfworkpath` command flag in any export command to change the storage path.

    Syntax:

    ```shell
    akamai [global flags] terraform <command> [command flags] [arguments...]
    ```

    Example:

    ```shell
    akamai --edgerc "~/.edgerc" --section "default" --accountkey "A-CCT1234:A-CCT5432" terraform export-property --version 3 "my-property"
    ```

	  The export adds a subdirectory containing declarative asset and variable configurations, supplemental files in JSON, and an import script in your active directory.

    If you're using the binary, you invoke the package command from your active directory as `./akamai-terraform` along with the credential details, the export command name, and other additional arguments.

    ```shell
    ./akamai-terraform --edgerc "~/.edgerc" --section "default" --accountkey "A-CCT1234:A-CCT5432" export-property --version 3 "my-property"
    ```

2. When the export is complete, run the import script to add your assets to your state.

> **Note:** Exported variable configuration files are entity-specific. When exporting multiple assets, merge the variable file content, removing any duplicates.

## Export command help

To get an overview of the library, run one of these:

- `akamai terraform help`. Lists commands and global flags available in the library.
- `akamai terraform list`. Lists available commands with a short description of each command's usage.

To get help information for a given Terraform CLI command, pass the export command and the `--help` flag.

Help command:

```shell
akamai terraform export-iam --help
```

Help output:

```shell
Name:
  akamai terraform export-iam

Usage:
  akamai [global flags] terraform export-iam [command flags] <subcommand>

Description:
  Generates Terraform configuration for Identity and Access Management resources.

Subcommands:
  all
  allowlist
  client
  group
  role
  user

Command Flags:
  --tfworkpath value  Directory used to store files created when running commands. (default: current directory)
  --help              show help (default: false)

Global Flags:
  --edgerc value, -e value                 Location of the credentials file (default: "/home/user/.edgerc") [$AKAMAI_EDGERC]
  --section value, -s value                Section of the credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
  --accountkey value, --account-key value  Account switch key [$AKAMAI_EDGERC_ACCOUNT_KEY]
```

## Core export commands

<ul>
  <li><a href="#export-apidefinitions">export-apidefinitions</a></li>
  <li><a href="#exportappsec">export-appsec</a></li>
  <li><a href="#exportclientlist">export-clientlist</a></li>
  <li><a href="#export-cloudaccess">export-cloudaccess</a></li>
  <li><a href="#export-cloudcertificate">export-cloudcertificate</a></li>
  <li><a href="#exportcloudlets-policy">export-cloudlets-policy</a></li>
  <li><a href="#exportcloudwrapper">export-cloudwrapper</a></li>
  <li><a href="#exportcps">export-cps</a></li>
  <li><a href="#exportdomain">export-domain</a></li>
  <li><a href="#exportdomainownership">export-domainownership</a></li>
  <li><a href="#exportedgekv">export-edgekv</a></li>
  <li><a href="#exportedgeworker">export-edgeworker</a></li>
  <li><a href="#exportiam">export-iam</a></li>
  <li><a href="#exportimaging">export-imaging</a></li>
  <li><a href="#exportmtls-keystore">export-mtls-keystore</a></li>
  <li><a href="#export-mtls-truststore">export-mtls-truststore</a></li>
  <li><a href="#exportproperty">export-property</a></li>
  <li><a href="#exportproperty-include">export-property-include</a></li>
  <li><a href="#exportzone">export-zone</a></li>
</ul>

## Common flags

This section provides you with a list of global and command flags that you can use universally across all export commands.

For details of command flags specific to individual export commands, see the information provided in each export command section.

### Global flags

|Flag|Description|
|---|---|
| <code>&#x2011;&#x2011;edgerc</code> (string) | Alias `-e`. The location of your credentials file. The default is <code>$HOME/.edgerc</code>. |
| <code>&#x2011;&#x2011;section</code> (string) | Alias `-s`. A credential set's section name. The default is <code>default</code>. |
| <code>&#x2011;&#x2011;accountkey</code> (string) | Alias `--account-key`. An account switch key.|
| <code>&#x2011;&#x2011;version</code> (integer) | Outputs the CLI version number.|

### Command flags

|Flag|Description|
|---|---|
| <code>&#x2011;&#x2011;tfworkpath</code> (string) | Sets the path to a directory in which you want to store the files created when running the export commands. The default is your active directory. |
| <code>&#x2011;&#x2011;help</code> (boolean) | Outputs a specific command's available options and descriptions. |


## export-apidefinitions

Export a Terraform configuration for your API definitions.

### Syntax

```shell
akamai terraform [global flags] export-apidefinitions [flags] <api_id>
```

### Basic usage

```shell
akamai terraform export-apidefinitions 12345
```

### Command flags

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--version` (integer) | The API's version number. If not specified, it exports the API's `latest` version by default. | `akamai terraform export-apidefinitions --version 1 12345` |
| `--format` (string) | The format of the API file, either `openapi` or `json`. Defaults to `openapi` if not specified. | `akamai terraform export-apidefinitions --format "openapi" 12345` |

## export‑appsec

Export a Terraform declarative security configuration and its targets and policies in JSON.

### Syntax

```shell
akamai [global flags] terraform export-appsec [command flags] <security configuration name>
```

### Basic usage

```shell
akamai terraform export‑appsec "my-security-config"
```

## export‑clientlist

Export a Terraform configuration for your client list.

### Syntax

```shell
akamai [global flags] terraform export-clientlist [command flags] <list_id>
```

### Basic usage

```shell
akamai terraform export-clientlist "123456_MYLISTID"
```

## export-cloudaccess

Export a Terraform configuration for your could access key. Add additional flags to narrow down your result.

### Syntax

```shell
akamai [global flags] terraform export-cloudaccess [command flags] <access_key_uid>
```

### Basic usage

```shell
akamai terraform export-cloudaccess 98765
```

### Command flags

<table border="1">
    <thead>
        <tr>
            <th>Flag</th>
            <th>Description</th>
            <th>Example</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>--group_id</code> (string)</td>
            <td>Required along with the <code>--contract_id</code> flag. Exports a cloud access key for a specific group and contract. <br /><br /><blockquote><b>Note:</b> Provide your group ID without the <code>grp_</code> prefix.</blockquote></td>
            <td rowspan="2"><code>akamai terraform export-cloudaccess --group_id 12345 --contract_id "C-0N7RAC7" 98765</code></td>
        </tr>
        <tr>
            <td><code>--contract_id</code> (string)</td>
            <td>Required along with the <code>--group_id</code> flag. Exports a cloud access key for a specific group and contract. <br /><br /><blockquote><b>Note:</b> Provide your contract ID without the <code>ctr_</code> prefix.</blockquote></td>
        </tr>
    </tbody>
</table>

## export-cloudcertificate

Export a Terraform configuration for your cloud certificate.

> **Note:**
>
> If the certificate is in the `READY_FOR_USE` or `ACTIVE` status, the `akamai_cloudcertificates_upload_signed_certificate` resource will also be included in your configuration. 
> If the certificate is in the `CSR_READY` status, the `akamai_cloudcertificates_upload_signed_certificate` resource will be generated but commented out.

### Syntax

```shell
akamai [global flags] terraform export-cloudcertificate [command flags] <cloud_certificate_name>
```

### Basic usage

```shell
akamai terraform export-cloudcertificate "my-cloudcertificate"
```

## export‑cloudlets-policy

Export a Terraform configuration for your cloudlet policy.

### Syntax

```shell
akamai [global flags] terraform export-cloudlets-policy [command flags] <policy_name>
```

### Basic usage

```shell
akamai terraform export-cloudlets-policy "my-policy"
```

## export‑cloudwrapper

Export a Terraform configuration for your cloud wrapper.

### Syntax

```shell
akamai [global flags] terraform export-cloudwrapper [command flags] <config_id>
```

### Basic usage

```shell
akamai terraform export-cloudwrapper 12345
```

## export‑cps

Export a Terraform configuration for your certificate.

### Syntax

```shell
akamai [global flags] terraform export-cps [command flags] <enrollment_id> <contract_id>
```

### Basic usage

```shell
akamai terraform export-cps "12345" "C-0N7RAC7"
```

## export‑domain

Export a Terraform configuration for your domain.

> **Note:** Mapping GTM entity names to Terraform resource names may require normalization. Invalid Terraform resource name characters are replaced by underscores during config generation.

### Syntax

```shell
akamai [global flags] terraform export-domain [command flags] <domain>
```

### Basic usage

```shell
akamai terraform export-domain "my-domain"
```

## export‑domainownership

Export a Terraform configuration for your domain ownership.

> **Notes:**
> - Each domain with a validation scope must exist as an `FQDN` for the export to succeed.
> - If a domain doesn't have a validation scope, it should match only one type: `HOST`, `DOMAIN`, or `WILDCARD`.
> - Domains with a domain status other than `VALIDATED` are exported as commented out for the `akamai_property_domainownership_validation` resource.
> - You can export up to 1000 domains at once.
> - The names of the exported Terraform resources are taken from the first domain in the list. If that domain contains invalid characters, these characters are replaced with underscores in the resources’ names.

### Syntax

```shell
akamai [global flags] terraform export-domainownership [command flags] <domain_name>[:validation_scope][,<domain_name>[:validation_scope]...]
```

### Basic usage

```shell
akamai terraform export-domainownership "my-domain:HOST,another-domain"
```

## export‑edgekv

Export a Terraform configuration for your namespace and network's EdgeKV.

### Syntax

```shell
akamai [global flags] terraform export-edgekv [command flags] <namespace_name> <network>
```

### Basic usage

```shell
akamai terraform export-edgekv "my-edgekv" "staging"
```

## export‑edgeworker

Export a Terraform configuration for your edgeworker's code bundle in `.tgz` format.

### Syntax

```shell
akamai [global flags] terraform export-edgeworker [command flags] <edgeworker_id>
```

### Basic usage

```shell
akamai terraform export-edgeworker 12345
```

### Command flag

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--bundlepath` (string) | Sets the path to a directory in which you want to store the EdgeWorkers `.tgz` code bundle. The default is your active directory. | `akamai terraform export-edgeworker --bundlepath "path/to/your/directory" 12345` |

## export‑iam

Export a Terraform configuration for users, groups, roles, IP allowlist, and CIDR block resources. Use additional options to narrow down your results to specific resources.

### Syntax

```shell
akamai [global flags] terraform export-iam [command flags] <subcommand> [subcommand flags]
```

### Basic usage

```shell
akamai terraform export-iam all
```

### Subcommands

| Subcommand            | Description                                                                                          | Example                                               |
|-----------------------|------------------------------------------------------------------------------------------------------|-------------------------------------------------------|
| `all` (boolean)       | Exports all available Terraform users, groups, roles, IP allowlist, client and CIDR block resources. | `akamai terraform export-iam all`                     |
| `allowlist` (boolean) | Exports the Terraform IP allowlist and CIDR block resources.                                         | `akamai terraform export-iam allowlist`               |
| `client` (string)     | Exports an API client. If the API client's ID is not provided then self API client is exported.<br /><br /> <blockquote><b>Note:</b> You can export the API client only if it has assigned credentials.</blockquote> | `akamai terraform export-iam client 1zk2gv34gkx5crv6` |
| `group` (string)      | Exports a group by ID with relevant users and their roles.                                           | `akamai terraform export-iam group 12345`             |
| `role` (string)       | Exports a role by ID with relevant users and their groups.                                           | `akamai terraform export-iam role 12345`              |
| `user` (string)       | Exports a user by their email with a relevant user's groups and roles.                               | `akamai terraform export-iam user "jsmith@email.com"` |

### Subcommand flag

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--only` (boolean) | An advanced option for the `group`, `role`, and `user` subcommands. Exports only specific information. | `akamai terraform export-iam user --only "jsmith@email.com"` |

## export‑imaging

Export a Terraform configuration for your imaging policy in JSON.

### Syntax

```shell
akamai [global flags] terraform export-imaging [command flags] <contract_id> <policy_set_id>
```

### Basic usage

```shell
akamai terraform export-imaging "C-0N7RAC7" "my-policy-set_1234"
```

### Command flags

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--policy-as-hcl` (boolean) | Exports your imaging policy with rules in HCL format. | `akamai terraform export-imaging --policy-as-hcl "C-0N7RAC7" "my-policy-set_12345"` |
| `--policy-json-dir` (string) | Sets the path to a directory in which you want to store your policy in JSON format. The default is your active directory.| `akamai terraform export-imaging --policy-json-dir "path/to/your/directory" "C-0N7RAC7" "my-policy-set_12345"` |

## export‑mtls-keystore

Export a Terraform configuration for your mTLS client certificate.
 
> **Note:** For third-party certificates, version information is always exported without being commented out as there is no automatic rotation. You can manually update the version fields as needed before applying the configuration.

### Syntax

```shell
akamai [global flags] terraform export-mtls-keystore [command flags] <certificate_id> [<group_id>  <contract_id>]
```

### Basic usage

```shell
akamai terraform export-mtls-keystore 12345
```

## export-mtls-truststore

Export a Terraform configuration for your Mutual TLS Edge Truststore CA set, along with
associated activations if they exist.

> **Note:**
>
> If the CA set you're exporting hasn't been activated on any networks, staging or production,
> the `akamai_mtlstruststore_ca_set_activation` resource will still be included in your
> configuration but commented out. This is to avoid accidental activation. To activate the CA set,
> uncomment the `akamai_mtlstruststore_ca_set_activation` resource and run `terraform apply`.

### Syntax

```shell
akamai [global flags] terraform export-mtls-truststore [command flags] <CA set name>
```

### Basic usage

```shell
akamai terraform export-mtls-truststore "my-ca-set-name"
```

### Command flags

| Flag | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | Example                                                                  |
| ------- |----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| `--version` (string) | Exports your declarative CA set configuration with possible activations for a specific CA set version. If provided, must be a positive integer. <br /><br /> <blockquote><b>Notes:</b> <ul><li>If you don't provide the <code>--version</code> flag, by default, it exports the <code>latest</code> CA set version whether it is active or not.</li><li>Since the `akamai_mtlstruststore_ca_set` resource always represents the latest version of a CA set, generating configuration for an older version will result in a non-empty Terraform plan.</li></ul></blockquote> | `akamai terraform export-mtls-truststore --version "1" "my-ca-set-name"` |

## export‑property

Export a Terraform configuration for your property along with its JSON-formatted rules, but without the includes. Use the [`export-property-include`](#exportproperty-include) command to export your includes.

> **Notes:**
>
> - Certain export conditions require the use of a particular property rule format. Verify whether your rule format matches the use case requirement and [update your rule format](https://techdocs.akamai.com/terraform/docs/set-up-property-provisioning#update-rule-format) as needed.
> - If the property you're exporting hasn't been activated on any networks, staging or production, the `akamai_property_activation` resource will still be included in your configuration but commented out. This is to avoid accidental activation. To activate the property, uncomment the `akamai_property_activation` resource and run `terraform apply`.
> - The `akamai_edge_hostname` resource isn't generated for `CCM` related hostnames.

### Syntax

```shell
akamai [global flags] terraform export-property [command flags] <property name>
```

### Basic usage

```shell
akamai terraform export-property "my-property"
```

### Command flags

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--version` (string) | Exports your declarative property configuration and the property's JSON-formatted rules without includes for a specific property version. <br /><br /> <blockquote><b>Note:</b> If you don't provide the <code>--version</code> flag, by default, it exports the <code>latest</code> property version whether it is active or not.</blockquote> | `akamai terraform export-property --version "1" "my-property-name"` |
| `--rules-as-hcl` (boolean) | Exports your declarative property configuration and the property's rules as the `akamai_property_rules_builder` data source in HCL format. <br /><br /> <blockquote><b>Note:</b> Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</blockquote> | `akamai terraform export-property --rules-as-hcl "my-property-name"` |
| `--split-depth` (integer)| Exports rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level. For example, `--split-depth=1` means that the default/root rule and all its direct children will be placed in dedicated files. Rules with higher nesting levels will be placed in a file of their closest ancestor. <br /><br /> <blockquote><b>Note:</b> You can use this flag only along with the <code>--rules-as-hcl</code> flag.</blockquote> | `akamai terraform export-property --split-depth=1 --rules-as-hcl "my-property-name"` |

## export‑property-include

Export a Terraform configuration for your property's include and its JSON-formatted rules.

> **Notes:**
>
> - Certain export conditions require the use of a particular property rule format. Verify whether your rule format matches the use case requirement and [update your rule format](https://techdocs.akamai.com/terraform/docs/set-up-includes#update-rule-format) as needed.
>  - If the include you're exporting hasn't been activated on any networks, staging or production, the `akamai_property_include_activation` resource will still be included in your configuration but commented out. This is to avoid accidental activation. To activate the include, uncomment the `akamai_property_include_activation` resource and run `terraform apply`.

### Syntax

```shell
akamai [global flags] terraform export-property-include [command flags] <contract_id> <include_name>
```

### Basic usage

```shell
akamai terraform export-property-include "C-0N7RAC7" "my-property-include"
```

### Command flags

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--rules-as-hcl` (boolean) | Exports your property's include as the `akamai_property_rules_builder` data source in HCL format. <br /><br /> <blockquote><b>Note:</b> Must be a dated rule format ≥ <code>v2023-01-05</code>. Cannot use <code>latest</code>.</blockquote> | `akamai terraform export-property-include --rules-as-hcl "my-property-include"` |
| `--split-depth` (integer)| Exports rules into a dedicated module. Each rule will be placed in a separate file up to a specified nesting level. For example, `--split-depth=1` means that the default/root rule and all its direct children will be placed in dedicated files. Rules with higher nesting levels will be placed in a file of their closest ancestor. <br /><br /> <blockquote><b>Note:</b> You can use this flag only along with the <code>--rules-as-hcl</code> flag.</blockquote> | `akamai terraform export-property-include --split-depth=1 --rules-as-hcl "my-property-include"` |

## export‑zone

Export a Terraform configuration for your zone and its related resources.

### Syntax

```shell
akamai [global flags] terraform export-zone [command flags] <zone>
```

### Basic usage

```shell
akamai terraform export-zone --resources "my-dns-zone.com"
```

### Command flags

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--resources` (boolean) | Generates a JSON-formatted resource file. `‑‑createconfig` uses this file as input. | `akamai terraform export-zone --resources "my-dns-zone.com"` |
| `--createconfig` (boolean) | Generates configurations based on the values in the output files from `‑‑resources`. | `akamai terraform export-zone --createconfig "my-dns-zone.com"` |
| `--importscript` (boolean) | Generates an import script for the generated zone configuration. | `akamai terraform export-zone --importscript "my-dns-zone.com"`  |

### Command flags to use with `--resources` and `--createconfig` only

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--recordname` (string) | Filters the generated resource list by a given record name or record names. You can provide it multiple times to specify different record names. | `akamai terraform export-zone --resources --recordname "my-record-name.com" "my-dns-zone.com"` |
| `--namesonly` (boolean) | Generates a resource file with recordset names only for all associated record types. | `akamai terraform export-zone --resources --namesonly "my-dns-zone.com"` |

### Command flags to use with `--createconfig` only

| Flag | Description | Example |
| ------- | --------- | --------- |
| `--segmentconfig` (boolean) | Generates a modularized configuration. | `akamai terraform export-zone --createconfig --segmentconfig "my-dns-zone.com"` |
| `--configonly` (boolean) | Generates a zone configuration without JSON itemization. The configuration generated varies based on which set of flags you use. | `akamai terraform export-zone --createconfig --configonly "my-dns-zone.com"` |

## General notes

Terraform variable configuration is generated in a separately named Terraform file for each Akamai entity type. These files will need to be merged by the Admin in the case where multiple entities are managed concurrently with the Terraform client.

## License

Copyright 2025 Akamai Technologies, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use these files except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.