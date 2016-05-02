package commands

import (
	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

func Create() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: append(commonFlags,
			cli.BoolFlag{
				Name:  "all",
				Usage: "Creates all of the databases in the database.yml",
			},
		),
		Usage: "Creates databases for you",
		Action: func(c *cli.Context) {
			commandInit(c)
			if c.Bool("all") {
				for _, conn := range pop.Connections {
					pop.CreateDB(conn)
				}
			} else {
				pop.CreateDB(getConn(c))
			}
		},
	}
}
