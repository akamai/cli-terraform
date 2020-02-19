# Akamai CLI for Global Traffic Management (GTM) Terraform Provider

[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-terraform-gtm)](https://goreportcard.com/report/github.com/akamai/cli-terraform-gtm) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform-gtm.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform-gtm?ref=badge_shield)

An [Akamai CLI](https://developer.akamai.com/cli) package for managing Terraform GTM Domain configurations and associated objects.

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install terraform-gtm
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-terraform-gtm/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compiling from Source

If you want to compile it from source, you will need Go 1.7 or later, and the [Dep](https://golang.github.io/dep/) package manager installed:

1. Fetch the package:  
  `go get github.com/akamai/cli-terraform-gtm`
2. Change to the package directory:  
  `cd $GOPATH/src/github.com/akamai/cli-terraform-gtm`
3. Install dependencies using `dep`:  
  `dep ensure`
4. Compile the binary:
  - Linux/macOS/*nix: `go build -o akamai-terraform-gtm`
  - Windows: `go build -o akamai-terraform-gtm.exe`
5. Move the binary (`akamai-terraform-gtm` or `akamai-terraform-gtm.exe`) in to your `PATH`

## Usage

```
  akamai-terraform-gtm [--edgerc] [--section] <command> [sub-command]

Description:
   Manage Terraform GTM Domain configurations and assoc objects

Global Flags:
   --tfworkpath value      file path location for placement of created and/or modified artifacts. Default: current directory
   --resources             Create json formatted resource import list file, <domain>_resources.json. Used as input by createconfig.
   --createconfig          Create Terraform configuration (<domain>.tf), gtmvars.tf, and import command script (<domain>_import.script) files

Built-In Commands:
  create-domain
  list
  help
```

### Create list of all domain objects. Written in json format to <domain>_resources.json

```
$ akamai terraform-gtm create-domain example.akadns.net --resources
```

### Generate Terraform GTM Domain configuration file <domain>.tf, vars config file, gtmvars.tf, and import script, <domain>_resource_import.script

``
$ akamai terraform-gtm create-domain example.akadns.net --createconfig
```

Notes:
1. Mapping GTM entity names to TF resource names may require normalization. Invalid TF resource name characters will be replaced by underscores, '_' in config generation.
 
## License

This package is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform-gtm.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-terraform-gtm?ref=badge_large)
