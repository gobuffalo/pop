package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	pipep "github.com/mattes/migrate/pipe"
)

var timerStart time.Time
var url string
var path string

/*
Migrate runs migration commands against the database. The code for this command was inspired by
the code found in http://github.com/mattes/migrate. Thanks to this code we don't need to rely on
the `migrate` binary being installed, and this makes life a lot nicer.
*/
func Migrate() cli.Command {
	return cli.Command{
		Name: "migrate",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Value: "./migrations",
				Usage: "Path to the migrations folder",
			},
			EnvFlag,
		},
		Usage: "Runs migrations against your database.",
		Action: func(c *cli.Context) {
			cmd := "up"
			if len(c.Args()) > 0 {
				cmd = c.Args().Get(0)
			}

			env := defaults.String(os.Getenv("GO_ENV"), c.String("e"))

			conn := pop.Connections[env]
			if conn == nil {
				if cmd == "help" {
					helpCmd()
					return
				} else {
					fmt.Printf("The database connection '%s' is not defined!\n", env)
					os.Exit(1)
				}
			}

			fmt.Printf("Database: %s\n", conn)
			url = conn.String()

			path = c.String("path")
			os.Mkdir(path, 0755)
			fmt.Printf("Migrations path: %s\n", path)

			switch cmd {
			case "create":
				name := c.Args().Get(1)
				if name == "" {
					fmt.Println("Please specify name.")
					os.Exit(1)
				}

				mf, err := migrate.Create(url, path, name)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Printf("Version %v migration files created in %v:\n", mf.Version, path)
				fmt.Println(mf.UpFile.FileName)
				fmt.Println(mf.DownFile.FileName)

			case "up":
				handlePipeFunc(migrate.Up)

			case "down":
				handlePipeFunc(migrate.Down)

			case "redo":
				handlePipeFunc(migrate.Redo)

			case "reset":
				handlePipeFunc(migrate.Reset)

			case "version":
				version, err := migrate.Version(url, path)
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

type pipeFunc func(chan interface{}, string, string)

func handlePipeFunc(fn pipeFunc) {
	timerStart = time.Now()
	pipe := pipep.New()
	go fn(pipe, url, path)
	ok := writePipe(pipe)
	printTimer()
	if !ok {
		os.Exit(1)
	}
}

func writePipe(pipe chan interface{}) (ok bool) {
	okFlag := true
	if pipe != nil {
		for {
			select {
			case item, more := <-pipe:
				if !more {
					return okFlag
				} else {
					switch item.(type) {

					case string:
						fmt.Println(item.(string))

					case error:
						c := color.New(color.FgRed)
						c.Println(item.(error).Error(), "\n")
						okFlag = false

					case file.File:
						f := item.(file.File)
						c := color.New(color.FgBlue)
						if f.Direction == direction.Up {
							c.Print(">")
						} else if f.Direction == direction.Down {
							c.Print("<")
						}
						fmt.Printf(" %s\n", f.FileName)

					default:
						text := fmt.Sprint(item)
						fmt.Println(text)
					}
				}
			}
		}
	}
	return okFlag
}

func printTimer() {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
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
