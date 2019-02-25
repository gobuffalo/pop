package config

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

var templates = packr.New("pop:genny:config", "../config/templates")

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

	t := gogen.TemplateTransformer(data, gogen.TemplateHelpers)
	g.Transformer(t)
	g.Transformer(genny.Dot())

	return g, nil
}
