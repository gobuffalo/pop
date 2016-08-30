package pop

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/markbates/going/defaults"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
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
		return errors.Wrap(genericCreate(s, model, cols), "sqlite create")
	})
}

func (m *sqlite) Update(s store, model *Model, cols Columns) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericUpdate(s, model, cols), "sqlite update")
	})
}

func (m *sqlite) Destroy(s store, model *Model) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericDestroy(s, model), "sqlite destroy")
	})
}

func (m *sqlite) SelectOne(s store, model *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericSelectOne(s, model, query), "sqlite select one")
	})
}

func (m *sqlite) SelectMany(s store, models *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericSelectMany(s, models, query), "sqlite select many")
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
	return errors.Wrapf(os.MkdirAll(d, 0766), "could not create SQLite database %s", m.ConnectionDetails.Database)
}

func (m *sqlite) DropDB() error {
	return errors.Wrapf(os.Remove(m.ConnectionDetails.Database), "could not drop SQLite database %s", m.ConnectionDetails.Database)
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
