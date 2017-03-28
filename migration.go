package pop

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/markbates/pop/fizz"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var mrx = regexp.MustCompile("(\\d+)_(.+)\\.(up|down)\\.(sql|fizz)")

func init() {
	MapTableName("schema_migrations", "schema_migration")
	MapTableName("schema_migration", "schema_migration")
}

var schemaMigrations = fizz.Table{
	Name: "schema_migration",
	Columns: []fizz.Column{
		{Name: "version", ColType: "string"},
	},
	Indexes: []fizz.Index{
		{Name: "version_idx", Columns: []string{"version"}, Unique: true},
	},
}

func MigrationCreate(path, name, ext string, up, down []byte) error {
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	err := os.MkdirAll(path, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create migrations path %s", path)
	}

	upf := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	err = ioutil.WriteFile(upf, up, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", upf)
	}
	fmt.Printf("> %s\n", upf)

	downf := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	err = ioutil.WriteFile(downf, down, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", downf)
	}

	fmt.Printf("> %s\n", downf)
	return nil
}

func (c *Connection) MigrateUp(path string) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "migration up: problem creating schema migrations")
	}
	return findMigrations(path, "up", 0, func(m migrationFile) error {
		exists, err := c.Where("version = ?", m.Version).Exists("schema_migration")
		if err != nil || exists {
			return errors.Wrapf(err, "problem checking for migration version %s", m.Version)
		}
		err = c.Transaction(func(tx *Connection) error {
			err := m.Execute(tx)
			if err != nil {
				return err
			}
			_, err = tx.Store.Exec(fmt.Sprintf("insert into schema_migration (version) values ('%s')", m.Version))
			return errors.Wrapf(err, "problem inserting migration version %s", m.Version)
		})
		if err == nil {
			fmt.Printf("> %s\n", m.FileName)
		}
		return err
	}, -1)
}

func (c *Connection) MigrateDown(path string, step int) error {
	now := time.Now()
	defer printTimer(now)

	err := c.createSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "migration down: problem creating schema migrations")
	}

	//increase skip by
	count, err := c.Count("schema_migration")
	if err != nil {
		return errors.Wrap(err, "migration down: unable count existing migration")
	}

	return findMigrations(path, "down", count, func(m migrationFile) error {
		exists, err := c.Where("version = ?", m.Version).Exists("schema_migration")
		if err != nil || !exists {
			fmt.Errorf("migration missing: %s", m.Version)
			return errors.Wrapf(err, "problem checking for migration version %s", m.Version)
		}
		err = c.Transaction(func(tx *Connection) error {
			err := m.Execute(tx)
			if err != nil {
				return err
			}
			err = tx.RawQuery("delete from schema_migration where version = ?", m.Version).Exec()
			return errors.Wrapf(err, "problem deleting migration version %s", m.Version)
		})
		if err == nil {
			fmt.Printf("< %s\n", m.FileName)
		}
		return err
	}, step)
}

func (c *Connection) MigrateReset(path string) error {
	err := c.MigrateDown(path, -1)
	if err != nil {
		return err
	}
	return c.MigrateUp(path)
}

func (c *Connection) createSchemaMigrations() error {
	err := c.Open()
	if err != nil {
		return errors.Wrap(err, "could not open connection")
	}
	_, err = c.Store.Exec("select * from schema_migration")
	if err == nil {
		return nil
	}

	return c.Transaction(func(tx *Connection) error {
		smSQL, err := c.Dialect.FizzTranslator().CreateTable(schemaMigrations)
		if err != nil {
			return errors.Wrap(err, "could not build SQL for schema migration table")
		}
		return errors.Wrap(tx.RawQuery(smSQL).Exec(), "could not create schema migration table")
	})
}

func findMigrations(dir string, direction string, runned int, fn func(migrationFile) error, step int) error {
	mfs := migrationFiles{}
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches := mrx.FindAllStringSubmatch(info.Name(), -1)
			if matches == nil || len(matches) == 0 {
				return nil
			}
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
		return nil
	})
	if direction == "down" {
		sort.Sort(sort.Reverse(mfs))
		// skip all runned migration
		mfs = mfs[len(mfs)-runned:]
		// run only required steps
		if step > 0 && len(mfs) >= step {
			mfs = mfs[:step]
		}
	} else {
		sort.Sort(mfs)
	}
	for _, mf := range mfs {
		err := fn(mf)
		if err != nil {
			return errors.Wrap(err, "error from called function")
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
