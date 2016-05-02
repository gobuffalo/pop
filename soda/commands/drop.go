package commands

import (
	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

func Drop() cli.Command {
	return cli.Command{
		Name: "drop",
		Flags: append(commonFlags,
			cli.BoolFlag{
				Name:  "all",
				Usage: "Drops all of the databases in the database.yml",
			},
		),
		Usage: "Drops databases for you",
		Action: func(c *cli.Context) {
			commandInit(c)
			if c.Bool("all") {
				for _, conn := range pop.Connections {
					pop.DropDB(conn)
				}
			} else {
				pop.DropDB(getConn(c))
			}
		},
	}
}
