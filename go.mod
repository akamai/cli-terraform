module github.com/akamai/cli-terraform

go 1.17

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v1.1.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.12.0
	github.com/akamai/cli v1.4.1
	github.com/fatih/color v1.13.0
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/shirou/gopsutil v2.20.4+incompatible
	github.com/stretchr/testify v1.7.1
	github.com/tj/assert v0.0.3
	github.com/urfave/cli/v2 v2.3.0
)

require (
	github.com/AlecAivazis/survey/v2 v2.2.7 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apex/log v1.9.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/briandowns/spinner v1.16.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/zclconf/go-cty v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
)

//replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 => ../akamaiopen-edgegrid-golang
//replace github.com/akamai/cli => ../cli
