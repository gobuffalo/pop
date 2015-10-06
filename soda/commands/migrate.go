package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/markbates/going/clam"
	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
	_ "github.com/mattes/migrate/migrate"
)

var Migrate = cli.Command{
	Name: "migrate",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Value: "./migrations",
			Usage: "Path to the migrations folder",
		},
		cli.StringFlag{
			Name:  "e",
			Value: "development",
			Usage: "The environment you want to run migrations against. Will use $GO_ENV if set.",
		},
	},
	Usage: "Runs migrations against your database. This just wraps https://github.com/mattes/migrate with good defaults.",
	Action: func(c *cli.Context) {
		env := defaults.String(os.Getenv("GO_ENV"), c.String("e"))

		conn := pop.Connections[env]
		fmt.Printf("database: %s\n", conn)
		os.Mkdir("migrations", 0755)

		name := "migrate"
		args := []string{"-url", conn.String(), "-path", c.String("path")}
		if len(c.Args()) == 0 {
			args = append(args, "up")
		} else {
			args = append(args, c.Args()...)
		}
		cmd := exec.Command(name, args...)

		err := clam.RunAndListen(cmd, func(s string) {
			fmt.Println(s)
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	},
}
