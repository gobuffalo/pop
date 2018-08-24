package config

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

var templates = packr.NewBox("../config/templates")

// New generator to create a database.yml file
func New(opts *Options) (*genny.Generator, error) {
	if opts.Root == "" {
		pwd, _ := os.Getwd()
		opts.Root = pwd
	}
	if opts.Prefix == "" {
		return nil, errors.New("you must provide a database name prefix")
	}
	if opts.FileName == "" {
		opts.FileName = "database.yml"
	}
	if opts.Dialect == "" {
		return nil, errors.New("you must provide a database dialect")
	}
	g := genny.New()

	f, err := templates.Open(opts.Dialect + ".yml.tmpl")
	if err != nil {
		return g, errors.Errorf("unknown dialect %s", opts.Dialect)
	}

	name := filepath.Join(opts.Root, opts.FileName+".tmpl")
	gf := genny.NewFile(name, f)
	g.File(gf)

	h := template.FuncMap{}
	data := map[string]interface{}{
		"opts": opts,
	}

	t := gotools.TemplateTransformer(data, h)
	g.Transformer(t)

	g.Transformer(genny.Replace("-dot-", "."))

	return g, nil
}
