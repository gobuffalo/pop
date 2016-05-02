package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

func Drop() cli.Command {
	return cli.Command{
		Name: "drop",
		Flags: []cli.Flag{
			EnvFlag,
			ConfigFlag,
			cli.BoolFlag{
				Name:  "all",
				Usage: "Drops all of the databases in the database.yml",
			},
			DebugFlag,
		},
		Usage: "Drops databases for you",
		Action: func(c *cli.Context) {
			pop.Debug = c.Bool("d")
			if c.Bool("all") {
				for _, conn := range pop.Connections {
					dropDB(conn)
				}
			} else {
				dropDB(getConn(c))
			}
		},
	}
}

func dropDB(c *pop.Connection) error {
	var err error
	if c.Dialect.Details().Database != "" {
		err = c.Dialect.DropDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}
