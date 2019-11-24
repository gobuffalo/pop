package main

import (
	"github.com/gobuffalo/pop/v4/soda/cmd"
)

func main() {
	cmd.RootCmd.Use = "soda"
	cmd.Execute()
}
