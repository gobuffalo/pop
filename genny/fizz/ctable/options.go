package ctable

import (
	"fmt"
	"time"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

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
		opts.Name = fmt.Sprintf("%s_create_%s", time.Now().Format("YYYYMMDDHHmmSS"), name.New(opts.TableName).Tableize())
	}
	return nil
}
