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
					Dropper(conn)
				}
			} else {
				Dropper(getConn(c))
			}
		},
	}
}

func Dropper(c *pop.Connection) error {
	var err error
	if c.Dialect.Details().Database != "" {
		err = c.Dialect.DropDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}
