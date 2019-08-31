package cempty

import (
	"fmt"
	"time"

	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

var nowFunc = time.Now

type nameable interface {
	Name() string
}

// Options for the empty migration generator.
type Options struct {
	// TableName is the name of the table.
	TableName string
	// Name is the name of the generated file.
	Name string
	// Path is the dir path where to generate the migration files.
	Path string
	// Translator is a Fizz translator to use when asking for SQL migrations.
	Translator nameable
	// Type is the type of migration to generate (sql or fizz).
	Type string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.TableName) == 0 {
		return errors.New("you must set a name for your table")
	}
	if len(opts.Path) == 0 {
		opts.Path = "migrations"
	}
	if len(opts.Name) == 0 {
		timestamp := nowFunc().UTC().Format("20060102150405")
		opts.Name = fmt.Sprintf("%s_create_%s", timestamp, name.New(opts.TableName).Tableize())
	}
	if len(opts.Type) == 0 {
		opts.Type = "fizz"
	}
	if opts.Type != "fizz" && opts.Type != "sql" {
		return errors.Errorf("%s migration type is not allowed", opts.Type)
	}
	if opts.Type == "sql" && opts.Translator == nil {
		return errors.New("sql migrations require a fizz translator")
	}
	return nil
}
