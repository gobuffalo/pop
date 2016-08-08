package pop

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

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

func MigrationCreate(path, name, ext string) error {
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

func (c *Connection) MigrateUp(path string) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return err
	}
	return findMigrations(path, "up", func(m migrationFile) error {
		i, err := c.Where("version = ?", m.Version).Count("schema_migrations")
		if err != nil {
			return err
		}
		if i == 0 {
			err = c.Transaction(func(tx *Connection) error {
				err := m.Execute(tx)
				if err != nil {
					return err
				}
				_, err = tx.Store.Exec(fmt.Sprintf("insert into schema_migrations (version) values ('%s')", m.Version))
				return err
			})
			if err == nil {
				fmt.Printf("> %s\n", m.FileName)
			}
			return err
		}
		return nil
	})
}

func (c *Connection) MigrateDown(path string) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return err
	}
	return findMigrations(path, "down", func(m migrationFile) error {
		i, err := c.Where("version = ?", m.Version).Count("schema_migrations")
		if err != nil {
			return err
		}
		if i > 0 {
			err = c.Transaction(func(tx *Connection) error {
				err := m.Execute(tx)
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
}

func (c *Connection) MigrateReset(path string) error {
	err := c.MigrateDown(path)
	if err != nil {
		return err
	}
	return c.MigrateUp(path)
}

func (c *Connection) createSchemaMigrations() error {
	err := c.Open()
	if err != nil {
		return err
	}
	return c.Transaction(func(tx *Connection) error {
		return tx.RawQuery(smSQL).Exec()
	})
}

func findMigrations(dir string, direction string, fn func(migrationFile) error) error {
	mfs := migrationFiles{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
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
					mfs = append(mfs, mf)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if direction == "down" {
		sort.Sort(sort.Reverse(mfs))
	} else {
		sort.Sort(mfs)
	}
	for _, mf := range mfs {
		err = fn(mf)
		if err != nil {
			return err
		}
	}
	return nil
}

func printTimer(timerStart time.Time) {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}
