package ctable

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/flect/name"
)

var nowFunc = time.Now

// Options for the table create generator.
type Options struct {
	// TableName is the name of the table.
	TableName string
	// Name is the name of the generated file.
	Name string
	// Path is the dir path where to generate the migration files.
	Path string
	// Attrs is a slice of columns to add to the table.
	Attrs attrs.Attrs
	// Translator is a Fizz translator to use when asking for SQL migrations.
	Translator fizz.Translator
	// Type is the type of migration to generate (sql or fizz).
	// For sql migrations, you'll have to provide a valid Translator too.
	Type string
	// ForceDefaultTimestamps enables auto timestamping for the generated table.
	ForceDefaultTimestamps bool `json:"force_default_timestamps"`
	// ForceDefaultID enables auto UUID for the generated table.
	ForceDefaultID bool `json:"force_default_id"`
}

// Validate that options are usable
func (opts *Options) Validate() error {
	if len(opts.TableName) == 0 {
		return errors.New("you must set a name for your table")
	}
	opts.TableName = name.New(opts.TableName).Tableize().String()
	if len(opts.Path) == 0 {
		opts.Path = "migrations"
	}
	if len(opts.Name) == 0 {
		timestamp := nowFunc().UTC().Format("20060102150405")
		opts.Name = fmt.Sprintf("%s_create_%s", timestamp, opts.TableName)
	}
	if len(opts.Type) == 0 {
		opts.Type = "fizz"
	}
	if opts.Type != "fizz" && opts.Type != "sql" {
		return fmt.Errorf("%s migration type is not allowed", opts.Type)
	}
	if opts.Type == "sql" && opts.Translator == nil {
		return errors.New("sql migrations require a fizz translator")
	}
	if opts.ForceDefaultID {
		var idFound bool
		for _, a := range opts.Attrs {
			switch a.Name.Underscore().String() {
			case "id":
				idFound = true
			}
		}
		if !idFound {
			// Add a default UUID
			id, err := attrs.Parse("id:uuid")
			if err != nil {
				return err
			}
			opts.Attrs = append([]attrs.Attr{id}, opts.Attrs...)
		}
	}
	return nil
}
