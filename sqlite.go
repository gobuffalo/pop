package pop

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

type sqlite struct {
	gil               *sync.Mutex
	smGil             *sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (m *sqlite) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *sqlite) URL() string {
	return m.ConnectionDetails.Database + "?_busy_timeout=5000"
}

func (m *sqlite) MigrationURL() string {
	return m.ConnectionDetails.URL
}

func (m *sqlite) Create(s store, model *Model, cols Columns) error {
	return m.locker(m.smGil, func() error {
		return genericCreate(s, model, cols)
	})
}

func (m *sqlite) Update(s store, model *Model, cols Columns) error {
	return m.locker(m.smGil, func() error {
		return genericUpdate(s, model, cols)
	})
}

func (m *sqlite) Destroy(s store, model *Model) error {
	return m.locker(m.smGil, func() error {
		return genericDestroy(s, model)
	})
}

func (m *sqlite) SelectOne(s store, model *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return genericSelectOne(s, model, query)
	})
}

func (m *sqlite) SelectMany(s store, models *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return genericSelectMany(s, models, query)
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
	d := filepath.Dir(m.ConnectionDetails.Database)
	err := os.MkdirAll(d, 0766)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (m *sqlite) DropDB() error {
	return os.Remove(m.ConnectionDetails.Database)
}

func (m *sqlite) TranslateSQL(sql string) string {
	return sql
}

func (m *sqlite) FizzTranslator() fizz.Translator {
	return translators.NewSQLite(m.Details().Database)
}

func newSQLite(deets *ConnectionDetails) dialect {
	deets.URL = fmt.Sprintf("sqlite3://%s", deets.Database)
	cd := &sqlite{
		gil:               &sync.Mutex{},
		smGil:             &sync.Mutex{},
		ConnectionDetails: deets,
	}

	return cd
}
