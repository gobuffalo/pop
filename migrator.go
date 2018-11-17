package pop

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gobuffalo/pop/logging"
	"github.com/pkg/errors"
)

var mrx = regexp.MustCompile(`^(\d+)_([^.]+)(\.[a-z0-9]+)?\.(up|down)\.(sql|fizz)$`)

// NewMigrator returns a new "blank" migrator. It is recommended
// to use something like MigrationBox or FileMigrator. A "blank"
// Migrator should only be used as the basis for a new type of
// migration system.
func NewMigrator(c *Connection) Migrator {
	return Migrator{
		Connection: c,
		Migrations: map[string]Migrations{
			"up":   {},
			"down": {},
		},
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Connection *Connection
	SchemaPath string
	Migrations map[string]Migrations
}

// UpLogOnly insert pending "up" migrations logs only, without applying the patch.
// It's used when loading the schema dump, instead of the migrations.
func (m Migrator) UpLogOnly() error {
	c := m.Connection
	return m.exec(func() error {
		mtn := c.MigrationTableName()
		mfs := m.Migrations["up"]
		sort.Sort(mfs)
		return c.Transaction(func(tx *Connection) error {
			for _, mi := range mfs {
				if mi.DBType != "all" && mi.DBType != c.Dialect.Name() {
					// Skip migration for non-matching dialect
					continue
				}
				exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
				if err != nil {
					return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
				}
				if exists {
					continue
				}
				_, err = tx.Store.Exec(fmt.Sprintf("insert into %s (version) values ('%s')", mtn, mi.Version))
				if err != nil {
					return errors.Wrapf(err, "problem inserting migration version %s", mi.Version)
				}
			}
			return nil
		})
	})
}

// Up runs pending "up" migrations and applies them to the database.
func (m Migrator) Up() error {
	c := m.Connection
	return m.exec(func() error {
		mtn := c.MigrationTableName()
		mfs := m.Migrations["up"]
		sort.Sort(mfs)
		applied := 0
		for _, mi := range mfs {
			if mi.DBType != "all" && mi.DBType != c.Dialect.Name() {
				// Skip migration for non-matching dialect
				continue
			}
			exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			if exists {
				continue
			}
			err = c.Transaction(func(tx *Connection) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				_, err = tx.Store.Exec(fmt.Sprintf("insert into %s (version) values ('%s')", mtn, mi.Version))
				return errors.Wrapf(err, "problem inserting migration version %s", mi.Version)
			})
			if err != nil {
				return errors.WithStack(err)
			}
			log(logging.Info, "> %s", mi.Name)
			applied++
		}
		if applied == 0 {
			log(logging.Info, "Migrations already up to date, nothing to apply")
		}
		return nil
	})
}

// Down runs pending "down" migrations and rolls back the
// database by the specified number of steps.
func (m Migrator) Down(step int) error {
	c := m.Connection
	return m.exec(func() error {
		mtn := c.MigrationTableName()
		count, err := c.Count(mtn)
		if err != nil {
			return errors.Wrap(err, "migration down: unable count existing migration")
		}
		mfs := m.Migrations["down"]
		sort.Sort(sort.Reverse(mfs))
		// skip all runned migration
		if len(mfs) > count {
			mfs = mfs[len(mfs)-count:]
		}
		// run only required steps
		if step > 0 && len(mfs) >= step {
			mfs = mfs[:step]
		}
		for _, mi := range mfs {
			exists, err := c.Where("version = ?", mi.Version).Exists(mtn)
			if err != nil || !exists {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			err = c.Transaction(func(tx *Connection) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				err = tx.RawQuery(fmt.Sprintf("delete from %s where version = ?", mtn), mi.Version).Exec()
				return errors.Wrapf(err, "problem deleting migration version %s", mi.Version)
			})
			if err != nil {
				return err
			}

			log(logging.Info, "< %s", mi.Name)
		}
		return nil
	})
}

// Reset the database by running the down migrations followed by the up migrations.
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return errors.WithStack(err)
	}
	return m.Up()
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func (m Migrator) CreateSchemaMigrations() error {
	c := m.Connection
	mtn := c.MigrationTableName()
	err := c.Open()
	if err != nil {
		return errors.Wrap(err, "could not open connection")
	}
	_, err = c.Store.Exec(fmt.Sprintf("select * from %s", mtn))
	if err == nil {
		return nil
	}

	return c.Transaction(func(tx *Connection) error {
		schemaMigrations := newSchemaMigrations(mtn)
		smSQL, err := c.Dialect.FizzTranslator().CreateTable(schemaMigrations)
		if err != nil {
			return errors.Wrap(err, "could not build SQL for schema migration table")
		}
		err = tx.RawQuery(smSQL).Exec()
		if err != nil {
			return errors.WithStack(errors.Wrap(err, smSQL))
		}
		return nil
	})
}

// Status prints out the status of applied/pending migrations.
func (m Migrator) Status() error {
	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.WithStack(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Version\tName\tStatus\t")
	for _, mf := range m.Migrations["up"] {
		exists, err := m.Connection.Where("version = ?", mf.Version).Exists(m.Connection.MigrationTableName())
		if err != nil {
			return errors.Wrapf(err, "problem with migration")
		}
		state := "Pending"
		if exists {
			state = "Applied"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", mf.Version, mf.Name, state)
	}
	return w.Flush()
}

// DumpMigrationSchema will generate a file of the current database schema
// based on the value of Migrator.SchemaPath
func (m Migrator) DumpMigrationSchema() error {
	if m.SchemaPath == "" {
		return nil
	}
	c := m.Connection
	schema := filepath.Join(m.SchemaPath, "schema.sql")
	f, err := os.Create(schema)
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.Dialect.DumpSchema(f)
	if err != nil {
		os.RemoveAll(schema)
		return errors.WithStack(err)
	}
	return nil
}

func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer func() {
		err := m.DumpMigrationSchema()
		if err != nil {
			log(logging.Warn, "Migrator: unable to dump schema: %v", err)
		}
	}()
	defer printTimer(now)

	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "Migrator: problem creating schema migrations")
	}
	return fn()
}

func printTimer(timerStart time.Time) {
	diff := time.Since(timerStart).Seconds()
	if diff > 60 {
		log(logging.Info, "%.4f minutes", diff/60)
	} else {
		log(logging.Info, "%.4f seconds", diff)
	}
}
