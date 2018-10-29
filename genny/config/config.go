package config

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

var templates = packr.NewBox("../config/templates")

// New generator to create a database.yml file
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	f, err := templates.Open(opts.Dialect + ".yml.tmpl")
	if err != nil {
		return g, errors.Errorf("unable to find database.yml template for dialect %s", opts.Dialect)
	}

	name := filepath.Join(opts.Root, opts.FileName+".tmpl")
	gf := genny.NewFile(name, f)
	g.File(gf)

	data := map[string]interface{}{
		"opts": opts,
	}

	t := gotools.TemplateTransformer(data, gotools.TemplateHelpers)
	g.Transformer(t)
	g.Transformer(genny.Dot())

	return g, nil
}
