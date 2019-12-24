# Akamai CLI for Global Traffic Management (GTM)

[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-gtm)](https://goreportcard.com/report/github.com/akamai/cli-gtm) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-gtm.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-gtm?ref=badge_shield)

An [Akamai CLI](https://developer.akamai.com/cli) package for managing GTM Domains and associated objects.

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install gtm
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-gtm/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compiling from Source

If you want to compile it from source, you will need Go 1.7 or later, and the [Dep](https://golang.github.io/dep/) package manager installed:

1. Fetch the package:  
  `go get github.com/akamai/cli-gtm`
2. Change to the package directory:  
  `cd $GOPATH/src/github.com/akamai/cli-gtm`
3. Install dependencies using `dep`:  
  `dep ensure`
4. Compile the binary:
  - Linux/macOS/*nix: `go build -o akamai-gtm`
  - Windows: `go build -o akamai-gtm.exe`
5. Move the binary (`akamai-gtm` or `akamai-gtm.exe`) in to your `PATH`

## Usage

```
  akamai-gtm [--edgerc] [--section] <command> [sub-command]

Description:
   Manage GTM Domains and assoc objects

Global Flags:
   --edgerc value  Location of the credentials file (default: "/home/elynes/.edgerc") [$AKAMAI_EDGERC]
   --section value     Section of the credentials file (default: "gtm") [$AKAMAI_EDGERC_SECTION]

Built-In Commands:
  update-datacenter
  update-property
  query-status
  list
  help
```

### Enable datacenters in domain

To enable one or more datacenters:

```
$ akamai gtm update-datacenter example.akadns.net --datacenter 3131 --datacenter 3132 --enable
```

### Update datacenters in property

To enable datacenters in a property:

```
$ akamai gtm update-property example.akadns.net testproperty --datacenter 3131 --disable
```

To modify a datacenter's weight:                                    

```
$ akamai gtm update-property example.akadns.net testproperty --datacenter 3131 --weight 20
```

To midify a datacenter's servers:

```
$ akamai gtm update-property example.akadns.net testproperty --datacenter 3131 --server 1.2.3.6 --server 1.2.1.1
```

### Query Status 

Query a datacenter's status:

```
$ akamai gtm query-status example.akadns.net --datacenter 3132
```

To query a property's status:

```
$ akamai gtm query-status example.akadns.net --property testproperty
```

## License

This package is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-gtm.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-gtm?ref=badge_large)
