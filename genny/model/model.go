package model

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/*
var templates embed.FS

// New returns a generator for creating a new model
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, err
	}

	if err := g.FS(sub); err != nil {
		return g, err
	}

	m := presenter{
		Name:        name.New(opts.Name),
		Encoding:    name.New(opts.Encoding),
		Validations: validatable(opts.Attrs),
		Imports:     buildImports(opts),
	}

	ctx := map[string]interface{}{
		"opts":  opts,
		"model": m,
	}
	help := map[string]interface{}{
		"capitalize": flect.Capitalize,
		"trim_package": func(t string) string {
			i := strings.LastIndex(t, ".")
			if i == -1 {
				return t
			}
			return t[i+1:]
		},
	}

	t := gogen.TemplateTransformer(ctx, help)
	g.Transformer(t)
	g.Transformer(genny.Replace("name-", flect.Singularize(opts.Name)))
	g.Transformer(genny.Replace("path-", opts.Path))
	return g, nil
}
