package commands

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/markbates/pop"
)

/*
Migrate runs migration commands against the database. The code for this command was inspired by
the code found in http://github.com/mattes/migrate. Thanks to this code we don't need to rely on
the `migrate` binary being installed, and this makes life a lot nicer.
*/
func Migrate() cli.Command {
	return cli.Command{
		Name: "migrate",
		Flags: append(commonFlags,
			cli.StringFlag{
				Name:  "path",
				Value: "./migrations",
				Usage: "Path to the migrations folder",
			},
		),
		Usage: "Runs migrations against your database.",
		Action: func(c *cli.Context) {
			commandInit(c)
			cmd := "up"
			if len(c.Args()) > 0 {
				cmd = c.Args().Get(0)
			}

			if cmd == "help" {
				helpCmd()
				return
			}

			conn := getConn(c)
			if pop.Debug {
				fmt.Printf("[POP] Database: %s\n", conn)
			}

			path := c.String("path")
			os.Mkdir(path, 0755)
			if pop.Debug {
				fmt.Printf("[POP] Migrations path: %s\n", path)
			}

			switch cmd {
			case "create":
				name := c.Args().Get(1)
				err := conn.MigrationCreate(path, name)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			case "up":
				ok := conn.MigrateUp(path)
				if !ok {
					fmt.Println("Migrations did not run successfully!")
					os.Exit(1)
				}

			case "down":
				ok := conn.MigrateDown(path)
				if !ok {
					fmt.Println("Migrations did not run successfully!")
					os.Exit(1)
				}

			case "redo":
				ok := conn.MigrateRedo(path)
				if !ok {
					fmt.Println("Migrations did not run successfully!")
					os.Exit(1)
				}

			case "reset":
				ok := conn.MigrateReset(path)
				if !ok {
					fmt.Println("Migrations did not run successfully!")
					os.Exit(1)
				}

			case "version":
				version, err := conn.MigrationVersion(path)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println(version)

			default:
				fallthrough
			case "help":
				helpCmd()
			}
		},
	}
}

func helpCmd() {
	os.Stderr.WriteString(
		`usage: pop migrate [-path=<path> -e=<connection>]  <command> [<args>]

Commands:
   create <name>  Create a new migration
   up             Apply all -up- migrations
   down           Apply all -down- migrations
   reset          Down followed by Up
   redo           Roll back most recent migration, then apply it again
   version        Show current migration version
   help           Show this help

'-path' defaults to "./migrations".
'-e' defaults to "development".
`)
}
