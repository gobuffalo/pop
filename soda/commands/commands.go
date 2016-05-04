package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
)

var commonFlags = []cli.Flag{
	EnvFlag,
	ConfigFlag,
	DebugFlag,
}

var EnvFlag = cli.StringFlag{
	Name:  "e",
	Value: "development",
	Usage: "The environment you want to run migrations against. Will use $GO_ENV if set.",
}

var ConfigFlag = cli.StringFlag{
	Name:  "c",
	Value: "config/database.yml",
	Usage: "The configuration file you would like to use.",
}

var DebugFlag = cli.BoolFlag{
	Name:  "d",
	Usage: "Debug/verbose mode",
}

func commandInit(c *cli.Context) {
	pop.Debug = c.Bool("d")
	fmt.Printf("Soda %s\n", c.App.Version)
	setConfigLocation(c)
}

func getEnv(c *cli.Context) string {
	return defaults.String(os.Getenv("GO_ENV"), c.String("e"))
}

func getConn(c *cli.Context) *pop.Connection {
	setConfigLocation(c)
	env := getEnv(c)
	conn := pop.Connections[getEnv(c)]
	if conn == nil {
		fmt.Printf("The database connection '%s' is not defined!\n", env)
		os.Exit(1)
	}
	return conn
}

func setConfigLocation(c *cli.Context) {
	abs, err := filepath.Abs(defaults.String(c.String("c"), "config/database.yml"))
	if err != nil {
		return
	}
	dir, file := filepath.Split(abs)
	pop.AddLookupPaths(dir)
	pop.ConfigName = file
}
