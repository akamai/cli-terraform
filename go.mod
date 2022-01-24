module github.com/akamai/cli-terraform

go 1.16

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.1.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.9.0
	github.com/akamai/cli-common-golang v0.0.0-20210716202303-5a2a24172430
	github.com/briandowns/spinner v1.16.0 // indirect
	github.com/fatih/color v1.13.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/shirou/gopsutil v2.20.4+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
)

//replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 => ../akamaiopen-edgegrid-golang
