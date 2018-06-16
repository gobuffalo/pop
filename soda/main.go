package main

import "github.com/gobuffalo/pop/soda/cmd"

func main() {
	cmd.RootCmd.Use = "soda"
	cmd.Execute()
}
