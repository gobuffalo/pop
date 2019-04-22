package ctable

import (
	"fmt"
	"time"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

// Options for the table create generator
type Options struct {
	TableName string
	Name      string
	Path      string
	Attrs     attrs.Attrs
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
		opts.Name = fmt.Sprintf("%s_create_%s.fizz", time.Now().Format("YYYYMMDDHHmmSS"), name.New(opts.TableName).Tableize())
	}
	return nil
}
