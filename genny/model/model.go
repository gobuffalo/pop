package model

import (
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

// New returns a generator for creating a new model
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("github.com/gobuffalo/pop/genny/model/templates", "../model/templates")); err != nil {
		return g, errors.WithStack(err)
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
	}

	t := gogen.TemplateTransformer(ctx, help)
	g.Transformer(t)
	g.Transformer(genny.Replace("-name-", opts.Name))
	g.Transformer(genny.Replace("-package-", opts.Package))
	return g, nil
}
