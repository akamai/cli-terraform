// Package main contains entry code for the application.
package main

import (
	"fmt"
	"os"

	"github.com/akamai/cli-terraform/v2/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
