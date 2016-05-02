package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

func Create() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			EnvFlag,
			ConfigFlag,
			cli.BoolFlag{
				Name:  "all",
				Usage: "Creates all of the databases in the database.yml",
			},
			DebugFlag,
		},
		Usage: "Creates databases for you",
		Action: func(c *cli.Context) {
			pop.Debug = c.Bool("d")
			if c.Bool("all") {
				for _, conn := range pop.Connections {
					createDB(conn)
				}
			} else {
				createDB(getConn(c))
			}
		},
	}
}

func createDB(c *pop.Connection) error {
	var err error
	if c.Dialect.Details().Database != "" {
		err = c.Dialect.CreateDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}
