package pop

import (
	"io"

	"github.com/gobuffalo/fizz"

	"github.com/gobuffalo/pop/v6/columns"
)

type crudable interface {
	SelectOne(conn *Connection, m *Model, q Query) error
	SelectMany(conn *Connection, m *Model, q Query) error
	Create(conn *Connection, m *Model, cols columns.Columns) error
	Update(conn *Connection, m *Model, cols columns.Columns) error
	UpdateQuery(conn *Connection, m *Model, cols columns.Columns, q Query) (int64, error)
	Destroy(conn *Connection, m *Model) error
	Delete(conn *Connection, m *Model, q Query) error
}

type fizzable interface {
	FizzTranslator() fizz.Translator
}

type quotable interface {
	Quote(key string) string
}

type dialect interface {
	crudable
	fizzable
	quotable

	Name() string
	DefaultDriver() string
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(sql string) string
	CreateDB() error
	DropDB() error
	DumpSchema(w io.Writer) error
	LoadSchema(r io.Reader) error
	Lock(fn func() error) error
	TruncateAll(conn *Connection) error
}

type afterOpenable interface {
	AfterOpen(conn *Connection) error
}
