package ctable

import (
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	t := fizz.NewTable(opts.TableName, map[string]interface{}{})
	f := genny.NewFileS(opts.Name, t.String())
	g.File(f)
	return g, nil
}
