package commands

import "github.com/codegangsta/cli"

var EnvFlag = cli.StringFlag{
	Name:  "e",
	Value: "development",
	Usage: "The environment you want to run migrations against. Will use $GO_ENV if set.",
}
