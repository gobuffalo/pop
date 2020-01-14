package model

import (
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/packr/v2"
)

// New returns a generator for creating a new model
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if err := g.Box(packr.New("github.com/gobuffalo/pop/genny/model/templates", "../model/templates")); err != nil {
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
	g.Transformer(genny.Replace("-name-", flect.Singularize(opts.Name)))
	g.Transformer(genny.Replace("-path-", opts.Path))
	return g, nil
}
