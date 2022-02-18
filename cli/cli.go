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

package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/akamai/cli-terraform/pkg/commands"
	akacli "github.com/akamai/cli/pkg/app"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
)

var (
	// VERSION holds current version of cli
	VERSION = "0.4.0"
)

// Run initializes the cli and runs it
func Run() error {
	term := terminal.Color()
	ctx := context.Background()
	ctx = terminal.Context(ctx, term)

	app := akacli.CreateAppTemplate(ctx, "terraform",
		"A CLI Plugin for Akamai Terraform Provider",
		"Administer and Manage Supported Akamai Feature resources with Terraform",
		VERSION)

	cmds, err := commands.CommandLocator()
	if err != nil {
		return fmt.Errorf(color.RedString("An error occurred initializing commands: %s"), err)
	}
	if len(cmds) > 0 {
		app.Commands = cmds
	}

	return app.RunContext(ctx, os.Args)
}
