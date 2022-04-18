// +build sqlite

package pop

import (
	"database/sql/driver"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/v6/columns"
	"github.com/gobuffalo/pop/v6/internal/defaults"
	"github.com/gobuffalo/pop/v6/logging"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Load SQLite3 CGo driver
)

const nameSQLite3 = "sqlite3"

func init() {
	AvailableDialects = append(AvailableDialects, nameSQLite3)
	dialectSynonyms["sqlite"] = nameSQLite3
	urlParser[nameSQLite3] = urlParserSQLite3
	newConnection[nameSQLite3] = newSQLite
	finalizer[nameSQLite3] = finalizerSQLite
}

var _ dialect = &sqlite{}

type sqlite struct {
	commonDialect
	gil   *sync.Mutex
	smGil *sync.Mutex
}

func (m *sqlite) Name() string {
	return nameSQLite3
}

func (m *sqlite) DefaultDriver() string {
	return nameSQLite3
}

func (m *sqlite) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *sqlite) URL() string {
	c := m.ConnectionDetails
	return c.Database + "?" + c.OptionsString("")
}

func (m *sqlite) MigrationURL() string {
	return m.ConnectionDetails.URL
}

func (m *sqlite) Create(s store, model *Model, cols columns.Columns) error {
	return m.locker(m.smGil, func() error {
		keyType, err := model.PrimaryKeyType()
		if err != nil {
			return err
		}
		switch keyType {
		case "int", "int64":
			var id int64
			cols.Remove(model.IDField())
			w := cols.Writeable()
			var query string
			if len(w.Cols) > 0 {
				query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", m.Quote(model.TableName()), w.QuotedString(m), w.SymbolizedString())
			} else {
				query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", m.Quote(model.TableName()))
			}
			log(logging.SQL, query, model.Value)
			res, err := s.NamedExec(query, model.Value)
			if err != nil {
				return err
			}
			id, err = res.LastInsertId()
			if err == nil {
				model.setID(id)
			}
			if err != nil {
				return err
			}
			return nil
		}
		if err := genericCreate(s, model, cols, m); err != nil {
			return fmt.Errorf("sqlite create: %w", err)
		}
		return nil
	})
}

func (m *sqlite) Update(s store, model *Model, cols columns.Columns) error {
	return m.locker(m.smGil, func() error {
		if err := genericUpdate(s, model, cols, m); err != nil {
			return fmt.Errorf("sqlite update: %w", err)
		}
		return nil
	})
}

func (m *sqlite) Destroy(s store, model *Model) error {
	return m.locker(m.smGil, func() error {
		if err := genericDestroy(s, model, m); err != nil {
			return fmt.Errorf("sqlite destroy: %w", err)
		}
		return nil
	})
}

func (m *sqlite) Delete(s store, model *Model, query Query) error {
	return genericDelete(s, model, query)
}

func (m *sqlite) SelectOne(s store, model *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		if err := genericSelectOne(s, model, query); err != nil {
			return fmt.Errorf("sqlite select one: %w", err)
		}
		return nil
	})
}

func (m *sqlite) SelectMany(s store, models *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		if err := genericSelectMany(s, models, query); err != nil {
			return fmt.Errorf("sqlite select many: %w", err)
		}
		return nil
	})
}

func (m *sqlite) Lock(fn func() error) error {
	return m.locker(m.gil, fn)
}

func (m *sqlite) locker(l *sync.Mutex, fn func() error) error {
	if defaults.String(m.Details().Options["lock"], "true") == "true" {
		defer l.Unlock()
		l.Lock()
	}
	err := fn()
	attempts := 0
	for err != nil && err.Error() == "database is locked" && attempts <= m.Details().RetryLimit() {
		time.Sleep(m.Details().RetrySleep())
		err = fn()
		attempts++
	}
	return err
}

func (m *sqlite) CreateDB() error {
	_, err := os.Stat(m.ConnectionDetails.Database)
	if err == nil {
		return fmt.Errorf("could not create SQLite database '%s'; database exists", m.ConnectionDetails.Database)
	}
	dir := filepath.Dir(m.ConnectionDetails.Database)
	err = os.MkdirAll(dir, 0766)
	if err != nil {
		return fmt.Errorf("could not create SQLite database '%s': %w", m.ConnectionDetails.Database, err)
	}
	f, err := os.Create(m.ConnectionDetails.Database)
	if err != nil {
		return fmt.Errorf("could not create SQLite database '%s': %w", m.ConnectionDetails.Database, err)
	}
	_ = f.Close()

	log(logging.Info, "created database '%s'", m.ConnectionDetails.Database)
	return nil
}

func (m *sqlite) DropDB() error {
	err := os.Remove(m.ConnectionDetails.Database)
	if err != nil {
		return fmt.Errorf("could not drop SQLite database %s: %w", m.ConnectionDetails.Database, err)
	}
	log(logging.Info, "dropped database '%s'", m.ConnectionDetails.Database)
	return nil
}

func (m *sqlite) TranslateSQL(sql string) string {
	return sql
}

func (m *sqlite) FizzTranslator() fizz.Translator {
	return translators.NewSQLite(m.Details().Database)
}

func (m *sqlite) DumpSchema(w io.Writer) error {
	cmd := exec.Command("sqlite3", m.Details().Database, ".schema")
	return genericDumpSchema(m.Details(), cmd, w)
}

func (m *sqlite) LoadSchema(r io.Reader) error {
	cmd := exec.Command("sqlite3", m.ConnectionDetails.Database)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.Copy(in, r)
	}()
	log(logging.SQL, strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	log(logging.Info, "loaded schema for %s", m.Details().Database)
	return nil
}

func (m *sqlite) TruncateAll(tx *Connection) error {
	const tableNames = `SELECT name FROM sqlite_master WHERE type = "table"`
	names := []struct {
		Name string `db:"name"`
	}{}

	err := tx.RawQuery(tableNames).All(&names)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return nil
	}
	stmts := []string{}
	for _, n := range names {
		stmts = append(stmts, fmt.Sprintf("DELETE FROM %s", m.Quote(n.Name)))
	}
	return tx.RawQuery(strings.Join(stmts, "; ")).Exec()
}

func newSQLite(deets *ConnectionDetails) (dialect, error) {
	deets.URL = fmt.Sprintf("sqlite3://%s", deets.Database)
	cd := &sqlite{
		gil:           &sync.Mutex{},
		smGil:         &sync.Mutex{},
		commonDialect: commonDialect{ConnectionDetails: deets},
	}

	return cd, nil
}

func urlParserSQLite3(cd *ConnectionDetails) error {
	db := strings.TrimPrefix(cd.URL, "sqlite://")
	db = strings.TrimPrefix(db, "sqlite3://")

	dbparts := strings.Split(db, "?")
	cd.Database = dbparts[0]

	if len(dbparts) != 2 {
		return nil
	}

	q, err := url.ParseQuery(dbparts[1])
	if err != nil {
		return fmt.Errorf("unable to parse sqlite query: %w", err)
	}

	if cd.Options == nil { // prevent panic
		cd.Options = make(map[string]string)
	}
	for k := range q {
		cd.Options[k] = q.Get(k)
	}

	return nil
}

func finalizerSQLite(cd *ConnectionDetails) {
	defs := map[string]string{
		"_busy_timeout": "5000",
	}
	forced := map[string]string{
		"_fk": "true",
	}
	if cd.Options == nil { // prevent panic
		cd.Options = make(map[string]string)
	}

	for k, v := range defs {
		cd.Options[k] = defaults.String(cd.Options[k], v)
	}

	for k, v := range forced {
		// respect user specified options but print warning!
		cd.Options[k] = defaults.String(cd.Options[k], v)
		if cd.Options[k] != v { // when user-defined option exists
			log(logging.Warn, "IMPORTANT! '%s: %s' option is required to work properly but your current setting is '%v: %v'.", k, v, k, cd.Options[k])
			log(logging.Warn, "It is highly recommended to remove '%v: %v' option from your config!", k, cd.Options[k])
		} // or override with `cd.Options[k] = v`?
		if cd.URL != "" && !strings.Contains(cd.URL, k+"="+v) {
			log(logging.Warn, "IMPORTANT! '%s=%s' option is required to work properly. Please add it to the database URL in the config!", k, v)
		} // or fix user specified url?
	}
}

func newSQLiteDriver() (driver.Driver, error) {
	return new(sqlite3.SQLiteDriver), nil
}
