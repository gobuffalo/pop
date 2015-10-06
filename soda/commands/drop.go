package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

var Drop = cli.Command{
	Name: "drop",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "e",
			Value: "development",
			Usage: "Drops the specified database from the database.yml file",
		},
		cli.BoolFlag{
			Name:  "all",
			Usage: "Drops all of the databases in the database.yml",
		},
	},
	Usage: "Drops databases for you",
	Action: func(c *cli.Context) {
		env := c.String("e")
		if c.Bool("all") {
			for _, conn := range pop.Connections {
				dropDB(conn)
			}
		} else {
			conn := pop.Connections[env]
			if conn == nil {
				fmt.Fprintf(os.Stderr, "%s is not a valid environment!\n", env)
				return
			}
			dropDB(conn)
		}
	},
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
