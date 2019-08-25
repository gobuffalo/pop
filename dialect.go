package pop

import (
	"io"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/pop/columns"
)

type crudable interface {
	SelectOne(store, *Model, Query) error
	SelectMany(store, *Model, Query) error
	Create(store, *Model, columns.Columns) error
	Update(store, *Model, columns.Columns) error
	Destroy(store, *Model) error
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
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(string) string
	CreateDB() error
	DropDB() error
	DumpSchema(io.Writer) error
	LoadSchema(io.Reader) error
	Lock(func() error) error
	TruncateAll(*Connection) error
}

type afterOpenable interface {
	AfterOpen(*Connection) error
}
