package cempty

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

// New creates a generator to make empty migration files.
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	var f genny.File
	if opts.Type == "sql" {
		f = genny.NewFileS(filepath.Join(opts.Path, fmt.Sprintf("%s.%s.up.sql", opts.Name, opts.Translator.Name())), "")
		g.File(f)
		f = genny.NewFileS(filepath.Join(opts.Path, fmt.Sprintf("%s.%s.down.sql", opts.Name, opts.Translator.Name())), "")
		g.File(f)
	} else {
		f = genny.NewFileS(filepath.Join(opts.Path, opts.Name+".up.fizz"), "")
		g.File(f)
		f = genny.NewFileS(filepath.Join(opts.Path, opts.Name+".down.fizz"), "")
		g.File(f)
	}
	return g, nil
}
