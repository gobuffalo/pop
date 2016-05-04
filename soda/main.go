package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop/soda/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "soda"
	app.Usage = "A tasty treat for all your database needs"
	app.Version = "2.7.4"

	app.Commands = []cli.Command{
		commands.Migrate(),
		commands.Create(),
		commands.Drop(),
	}

	app.Run(os.Args)
}
