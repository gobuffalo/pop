// Package main is the entry point for the soda command-line tool, which is used for database management tasks such as
// migrations and seeding in the Buffalo framework.
package main

import (
	"github.com/gobuffalo/pop/v6/soda/cmd"
)

func main() {
	cmd.RootCmd.Use = "soda"
	cmd.Execute()
}
