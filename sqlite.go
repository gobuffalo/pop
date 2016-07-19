package pop

// SQLite is currently not supported due to cgo issues

import (
	"os"
	"path/filepath"

	. "github.com/markbates/pop/columns"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	ConnectionDetails *ConnectionDetails
}

func (m *SQLite) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *SQLite) URL() string {
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

func (m *SQLite) CreateDB() error {
	d := filepath.Dir(m.ConnectionDetails.URL)
	return os.MkdirAll(d, 0755)
}

func (m *SQLite) DropDB() error {
	return os.Remove(m.ConnectionDetails.URL)
}

func (m *SQLite) TranslateSQL(sql string) string {
	return sql
}

func NewSQLite(deets *ConnectionDetails) Dialect {
	cd := &SQLite{
		ConnectionDetails: deets,
	}

	return cd
}
