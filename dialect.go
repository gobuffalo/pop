package pop

import (
	"io"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/pop/columns"
)

type dialect interface {
	Name() string
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(string) string
	Create(store, *Model, columns.Columns) error
	Update(store, *Model, columns.Columns) error
	Destroy(store, *Model) error
	SelectOne(store, *Model, Query) error
	SelectMany(store, *Model, Query) error
	CreateDB() error
	DropDB() error
	DumpSchema(io.Writer) error
	LoadSchema(io.Reader) error
	FizzTranslator() fizz.Translator
	Lock(func() error) error
	TruncateAll(*Connection) error
	Quote(key string) string
}

type afterOpenable interface {
	AfterOpen(*Connection) error
}
