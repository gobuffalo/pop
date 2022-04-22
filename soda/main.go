package main

import (
	"github.com/gobuffalo/pop/v6/soda/cmd"
)

func main() {
	cmd.RootCmd.Use = "soda"
	cmd.Execute()
}
