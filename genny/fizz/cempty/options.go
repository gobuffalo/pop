package cempty

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/flect/name"
)

var nowFunc = time.Now

type nameable interface {
	Name() string
}

// Options for the empty migration generator.
type Options struct {
	// TableName is the name of the table.
	// Deprecated: use Name directly since TableName doesn't make sense in this generator.
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

// Validate that options are usable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must set a name for your migration")
	}
	if len(opts.Path) == 0 {
		opts.Path = "migrations"
	}
	timestamp := nowFunc().UTC().Format("20060102150405")
	opts.Name = fmt.Sprintf("%s_%s", timestamp, name.New(opts.Name).Underscore())

	if len(opts.Type) == 0 {
		opts.Type = "fizz"
	}
	if opts.Type != "fizz" && opts.Type != "sql" {
		return fmt.Errorf("%s migration type is not allowed", opts.Type)
	}
	if opts.Type == "sql" && opts.Translator == nil {
		return errors.New("sql migrations require a fizz translator")
	}
	return nil
}
