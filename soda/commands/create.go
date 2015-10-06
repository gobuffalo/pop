package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

var Create = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "e",
			Value: "development",
			Usage: "Creates the specified database from the database.yml file",
		},
		cli.BoolFlag{
			Name:  "all",
			Usage: "Creates all of the databases in the database.yml",
		},
	},
	Usage: "Creates databases for you",
	Action: func(c *cli.Context) {
		env := c.String("e")
		if c.Bool("all") {
			for _, conn := range pop.Connections {
				createDB(conn)
			}
		} else {
			conn := pop.Connections[env]
			if conn == nil {
				fmt.Fprintf(os.Stderr, "%s is not a valid environment!\n", env)
				return
			}
			createDB(conn)
		}
	},
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
