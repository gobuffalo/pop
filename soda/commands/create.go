package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
)

func Create() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			EnvFlag,
			cli.BoolFlag{
				Name:  "all",
				Usage: "Creates all of the databases in the database.yml",
			},
		},
		Usage: "Creates databases for you",
		Action: func(c *cli.Context) {
			env := defaults.String(os.Getenv("GO_ENV"), c.String("e"))
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
