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
	akamai "github.com/akamai/cli-common-golang"
	"os"
)

var (
	VERSION = "0.0.2"
)

func main() {
	akamai.CreateApp(
		"gtm",
		"A CLI for GTM",
		"Manage GTM Domains and assoc objects",
		VERSION,
		"gtm",
		commandLocator,
	)

	setHelpTemplates()
	akamai.App.Run(os.Args)
}
