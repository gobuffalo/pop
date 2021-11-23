package config

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/*
var templates embed.FS

// New generator to create a database.yml file
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, err
	}

	f, err := templates.Open("templates/" + opts.Dialect + ".yml.tmpl")
	if err != nil {
		return g, fmt.Errorf("unable to find database.yml template for dialect %s", opts.Dialect)
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
