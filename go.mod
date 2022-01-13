module github.com/akamai/cli-terraform

go 1.16

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.1.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.10.0
	github.com/akamai/cli-common-golang v0.0.0-20210716202303-5a2a24172430
	github.com/briandowns/spinner v1.16.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/shirou/gopsutil v2.20.4+incompatible
	github.com/stretchr/testify v1.7.0
	github.com/tj/assert v0.0.3
	github.com/urfave/cli/v2 v2.3.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
)

//replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 => ../akamaiopen-edgegrid-golang
//replace github.com/akamai/cli-common-golang => ../cli-common-golang
