module github.com/akamai/cli-terraform

go 1.16

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.1.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.11.0
	github.com/akamai/cli v1.4.1
	github.com/briandowns/spinner v1.16.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/shirou/gopsutil v2.20.4+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tj/assert v0.0.3
	github.com/urfave/cli/v2 v2.3.0
)

//replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 => ../akamaiopen-edgegrid-golang
//replace github.com/akamai/cli => ../cli
