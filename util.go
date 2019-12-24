// Copyright 2019. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func setHelpTemplates() {
	cli.AppHelpTemplate =
		color.YellowString("Usage: \n") +
			`{{if or (or (eq .HelpName "akamai-gtm update-datacenter") (eq .HelpName "akamai gtm update-datacenter")) (or (eq .HelpName "akamai-gtm update-property") (eq .HelpName "akamai gtm update-property")) (or (eq .HelpName "akamai-gtm query-status") (eq .HelpName "akamai gtm query-status"))}}` +
			color.BlueString(`	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{end}}`) +
			`{{else}}` +
			color.BlueString(`	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}}{{range .VisibleFlags}} [--{{.Name}}]{{end}}{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{if .Commands}} <command> [sub-command]{{end}}{{end}}`) +
			`{{end}}` +

			"{{if .Description}}\n\n" +
			color.YellowString("Description:\n") +
			"   {{.Description}}" +
			"\n\n{{end}}" +

			"{{if .VisibleFlags}}" +
			color.YellowString("Global Flags:\n") +
			"{{range $index, $option := .VisibleFlags}}" +
			"{{if $index}}\n{{end}}" +
			"   {{$option}}" +
			"{{end}}" +
			"\n\n{{end}}" +

			"{{if .VisibleCommands}}" +
			`{{if or (or (eq .HelpName "akamai-gtm update-datacenter") (eq .HelpName "akamai gtm update-datacenter")) (or (eq .HelpName "akamai-gtm update-property") (eq .HelpName "akamai gtm update-property")) (or (eq .HelpName "akamai-gtm query-status") (eq .HelpName "akamai gtm query-status"))}}` +
			`{{else}}` +
			color.YellowString("Built-In Commands:\n") +
			`{{end}}` +
			"{{range .VisibleCategories}}" +
			"{{if .Name}}" +
			"\n{{.Name}}\n" +
			"{{end}}" +
			"{{range .VisibleCommands}}" +
			color.GreenString("  {{.Name}}") +
			"{{if .Aliases}} ({{ $length := len .Aliases }}{{if eq $length 1}}alias:{{else}}aliases:{{end}} " +
			"{{range $index, $alias := .Aliases}}" +
			"{{if $index}}, {{end}}" +
			color.GreenString("{{$alias}}") +
			"{{end}}" +
			"){{end}}\n" +
			"{{end}}" +
			"{{end}}" +
			"{{end}}\n" +

			"{{if .Copyright}}" +
			color.HiBlackString("{{.Copyright}}") +
			"{{end}}\n"

	cli.CommandHelpTemplate =
		color.YellowString("Name: \n") +
			"   {{.HelpName}}\n\n" +

			`{{if .Description}}` +
			color.YellowString("Description: \n") +
			"   {{.Description}}\n\n" +
			`{{end}}` +

			color.YellowString("Usage: \n") +
			color.BlueString("   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{end}} {{if .VisibleFlags}}{{range .VisibleFlags}}[--{{.Name}}] {{end}}{{end}}{{end}}\n\n") +

			"{{if .Category}}" +
			color.YellowString("Type: \n") +
			"   {{.Category}}\n\n{{end}}" +

			"{{if .VisibleFlags}}" +
			color.YellowString("Flags: \n") +
			"{{range .VisibleFlags}}   {{.}}\n{{end}}\n{{end}}" +

			"{{if .Subcommands}}" +
			"{{range .Subcommands}}   {{.Name}}\n{{end}}{{end}}"

	cli.SubcommandHelpTemplate = cli.CommandHelpTemplate
}
