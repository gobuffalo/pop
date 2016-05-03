package pop

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	// _ "github.com/mattes/migrate/driver/sqlite3"
	pipep "github.com/mattes/migrate/pipe"
)

func (c *Connection) MigrationCreate(path, name string) error {
	if name == "" {
		return errors.New("Please specify name.")
	}

	mf, err := migrate.Create(c.URL(), path, name)
	if err != nil {
		return err
	}

	fmt.Printf("Version %v migration files created in %v:\n", mf.Version, path)
	fmt.Println(mf.UpFile.FileName)
	fmt.Println(mf.DownFile.FileName)
	return nil
}

func (c *Connection) MigrateUp(path string) bool {
	return handlePipeFunc(migrate.Up, c.URL(), path)
}

func (c *Connection) MigrateDown(path string) bool {
	return handlePipeFunc(migrate.Down, c.URL(), path)
}

func (c *Connection) MigrateRedo(path string) bool {
	return handlePipeFunc(migrate.Redo, c.URL(), path)
}

func (c *Connection) MigrateReset(path string) bool {
	return handlePipeFunc(migrate.Reset, c.URL(), path)
}

func (c *Connection) MigrationVersion(path string) (uint64, error) {
	return migrate.Version(c.URL(), path)
}

type pipeFunc func(chan interface{}, string, string)

func handlePipeFunc(fn pipeFunc, url string, path string) bool {
	timerStart := time.Now()
	pipe := pipep.New()
	go fn(pipe, url, path)
	ok := writePipe(pipe)
	printTimer(timerStart)
	return ok
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

func printTimer(timerStart time.Time) {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}
