package pop

// SQLite is currently not supported due to cgo issues

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/markbates/going/defaults"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

type SQLite struct {
	gil               *sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (m *SQLite) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *SQLite) URL() string {
	return m.ConnectionDetails.Database + "?cache=shared&mode=rwc"
}

func (m *SQLite) MigrationURL() string {
	return m.ConnectionDetails.URL
}

func (m *SQLite) Create(store Store, model *Model, cols Columns) error {
	return genericCreate(store, model, cols)

}

func (m *SQLite) Update(store Store, model *Model, cols Columns) error {
	return genericUpdate(store, model, cols)

}

func (m *SQLite) Destroy(store Store, model *Model) error {
	return genericDestroy(store, model)

}

func (m *SQLite) SelectOne(store Store, model *Model, query Query) error {
	return genericSelectOne(store, model, query)
}

func (m *SQLite) SelectMany(store Store, models *Model, query Query) error {
	return genericSelectMany(store, models, query)
}

func (m *SQLite) Lock(fn func() error) error {
	if defaults.String(m.Details().Options["lock"], "true") == "true" {
		defer m.gil.Unlock()
		m.gil.Lock()
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

func (m *SQLite) CreateDB() error {
	d := filepath.Dir(m.ConnectionDetails.Database)
	err := os.MkdirAll(d, 0766)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (m *SQLite) DropDB() error {
	return os.Remove(m.ConnectionDetails.Database)
}

func (m *SQLite) TranslateSQL(sql string) string {
	return sql
}

func (m *SQLite) FizzTranslator() fizz.Translator {
	return translators.NewSQLite(m.Details().Database)
}

func NewSQLite(deets *ConnectionDetails) Dialect {
	deets.URL = fmt.Sprintf("sqlite3://%s", deets.Database)
	cd := &SQLite{
		gil:               &sync.Mutex{},
		ConnectionDetails: deets,
	}

	return cd
}
