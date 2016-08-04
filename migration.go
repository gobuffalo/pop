package pop

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/markbates/pop/fizz"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	_ "github.com/mattes/migrate/driver/sqlite3"
)

var mrx = regexp.MustCompile("(\\d+)_(.+)\\.(up|down)\\.(sql|fizz)")

var smSQL = `CREATE TABLE IF NOT EXISTS "schema_migrations"(
	"version" TEXT NOT NULL,
	PRIMARY KEY("version")
);

CREATE UNIQUE INDEX IF NOT EXISTS "version_idx" ON "schema_migrations"("version");`

func (c *Connection) MigrationCreate(path, name, ext string) error {
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	up := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	err := ioutil.WriteFile(up, []byte(""), 0666)
	if err != nil {
		return err
	}
	fmt.Printf("> %s\n", up)

	down := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	err = ioutil.WriteFile(down, []byte(""), 0666)
	if err != nil {
		return err
	}

	fmt.Printf("> %s\n", down)
	return err
}

type migrationFile struct {
	Path      string
	FileName  string
	Version   string
	Name      string
	Direction string
	FileType  string
}

func (m migrationFile) execute(c *Connection) error {
	b, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	return c.RawQuery(string(b)).Exec()
}

func findMigrations(dir string, direction string, fn func(migrationFile) error) error {
	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches := mrx.FindAllStringSubmatch(info.Name(), -1)
			if len(matches) >= 0 {
				m := matches[0]
				mf := migrationFile{
					Path:      p,
					FileName:  m[0],
					Version:   m[1],
					Name:      m[2],
					Direction: m[3],
					FileType:  m[4],
				}
				if mf.Direction == direction {
					return fn(mf)
				}
			}
		}
		return nil
	})
}

func (c *Connection) MigrateUp(path string) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return err
	}
	return wrapWithTemplates(path, c, func(dir string) error {
		return findMigrations(dir, "up", func(m migrationFile) error {
			i, err := c.Where("version = ?", m.Version).Count("schema_migrations")
			if err != nil {
				return err
			}
			if i == 0 {
				err = c.Transaction(func(tx *Connection) error {
					err := m.execute(tx)
					if err != nil {
						return err
					}
					_, err = tx.Store.Exec("insert into schema_migrations (version) values(?)", m.Version)
					return err
				})
				if err == nil {
					fmt.Printf("> %s\n", m.FileName)
				}
				return err
			}
			return nil
		})
	})
}

func (c *Connection) MigrateDown(path string) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return err
	}
	return wrapWithTemplates(path, c, func(dir string) error {
		return findMigrations(dir, "down", func(m migrationFile) error {
			i, err := c.Where("version = ?", m.Version).Count("schema_migrations")
			if err != nil {
				return err
			}
			if i > 0 {
				err = c.Transaction(func(tx *Connection) error {
					err := m.execute(tx)
					if err != nil {
						return err
					}
					_, err = tx.Store.Exec("delete from schema_migrations where version = ?", m.Version)
					return err
				})
				if err == nil {
					fmt.Printf("< %s\n", m.FileName)
				}
				return err
			}
			return nil
		})
	})
}

func (c *Connection) MigrateReset(path string) error {
	err := c.MigrateDown(path)
	if err != nil {
		return err
	}
	return c.MigrateUp(path)
}

func (c *Connection) createSchemaMigrations() error {
	return c.Transaction(func(tx *Connection) error {
		return tx.RawQuery(smSQL).Exec()
	})
}

func printTimer(timerStart time.Time) {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}

func wrapWithTemplates(p string, c *Connection, fn func(dir string) error) error {
	dir, err := ioutil.TempDir("", "pop")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}
	for _, file := range files {
		ext := path.Ext(file.Name())
		tfn := filepath.Join(dir, strings.Replace(file.Name(), ".fizz", ".sql", 1))
		content, _ := ioutil.ReadFile(filepath.Join(p, file.Name()))
		f, err := os.Create(tfn)
		if err != nil {
			return err
		}
		switch ext {
		case ".fizz":
			s, _ := fizz.AString(string(content), c.Dialect.FizzTranslator())
			if err != nil {
				return err
			}
			fmt.Fprint(f, s)
		case ".sql":
			t := template.Must(template.New("letter").Parse(string(content)))
			err = t.Execute(f, c.Dialect.Details())
			if err != nil {
				return err
			}
		}
	}
	err = fn(dir)
	if err != nil {
		return err
	}
	return nil
}
