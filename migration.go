package pop

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	_ "github.com/mattes/migrate/driver/sqlite3"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	pipep "github.com/mattes/migrate/pipe"
)

func (c *Connection) MigrationCreate(path, name string) error {
	if name == "" {
		return errors.New("Please specify name.")
	}

	mf, err := migrate.Create(c.MigrationURL(), path, name)
	if err != nil {
		return err
	}

	fmt.Printf("Version %v migration files created in %v:\n", mf.Version, path)
	fmt.Println(mf.UpFile.FileName)
	fmt.Println(mf.DownFile.FileName)
	return nil
}

func (c *Connection) MigrateUp(path string) bool {
	return wrapWithTemplates(path, c, func(dir string) bool {
		return handlePipeFunc(migrate.Up, c.MigrationURL(), dir)
	})
}

func (c *Connection) MigrateDown(path string) bool {
	return wrapWithTemplates(path, c, func(dir string) bool {
		return handlePipeFunc(migrate.Down, c.MigrationURL(), dir)
	})
}

func (c *Connection) MigrateRedo(path string) bool {
	return wrapWithTemplates(path, c, func(dir string) bool {
		return handlePipeFunc(migrate.Redo, c.MigrationURL(), dir)
	})
}

func (c *Connection) MigrateReset(path string) bool {
	return wrapWithTemplates(path, c, func(dir string) bool {
		return handlePipeFunc(migrate.Reset, c.MigrationURL(), dir)
	})
}

func (c *Connection) MigrationVersion(path string) (uint64, error) {
	return migrate.Version(c.MigrationURL(), path)
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

func wrapWithTemplates(path string, c *Connection, fn func(dir string) bool) bool {
	dir, err := ioutil.TempDir("", "pop")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		tfn := filepath.Join(dir, file.Name())
		content, _ := ioutil.ReadFile(filepath.Join(path, file.Name()))
		t := template.Must(template.New("letter").Parse(string(content)))
		f, err := os.Create(tfn)
		if err != nil {
			fmt.Println(err)
			return false
		}
		err = t.Execute(f, c.Dialect.Details())
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return fn(dir)
}
