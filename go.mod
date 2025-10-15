module github.com/akamai/cli-terraform/v2

go 1.23.6

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v12 v12.1.0
	github.com/akamai/cli/v2 v2.0.2
	github.com/fatih/color v1.18.0
	github.com/hashicorp/hcl/v2 v2.23.0
	github.com/jinzhu/copier v0.4.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v2 v2.19.3
	github.com/wk8/go-ordered-map/v2 v2.1.8
)

require (
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/briandowns/spinner v1.23.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zclconf/go-cty v1.15.1 // indirect
	go.uber.org/ratelimit v0.3.1 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/akamai/AkamaiOPEN-edgegrid-golang/v12 => ../akamaiopen-edgegrid-golang
//replace github.com/akamai/cli/v2 => ../cli
replace gopkg.in/yaml.v2 v2.2.2 => gopkg.in/yaml.v2 v2.4.0 // Fix security vulnerability
