package model

import (
	"path"
	"sort"
	"strings"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("github.com/gobuffalo/pop/genny/model/templates", "../model/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	m := model{
		Name:                  name.New(opts.Name),
		Encoding:              name.New(opts.Encoding),
		ValidatableAttributes: validatable(opts.Attrs),
		Imports:               buildImports(opts),
	}

	sort.Strings(m.Imports)

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

type model struct {
	Name                  name.Ident
	Encoding              name.Ident
	Imports               []string
	ValidatableAttributes attrs.Attrs
}

func validatable(ats attrs.Attrs) attrs.Attrs {
	var xats attrs.Attrs
	for _, a := range ats {
		n := a.Name.Proper().String()
		if n == "CreatedAt" || n == "UpdatedAt" {
			continue
		}
		switch a.GoType() {
		case "string", "time.Time", "int":
			xats = append(xats, a)
		}

	}
	return xats
}

func buildImports(opts *Options) []string {
	imps := map[string]bool{
		"github.com/gobuffalo/validate": true,
	}
	imps[path.Join("encoding", strings.ToLower(opts.Encoding))] = true
	ats := opts.Attrs
	for _, a := range ats {
		switch a.GoType() {
		case "uuid":
			imps["github.com/gofrs/uuid"] = true
		case "time.Time":
			imps["time"] = true
		default:
			if strings.HasPrefix(a.GoType(), "nulls") {
				imps["github.com/gobuffalo/pop/nulls"] = true
			}
			if strings.HasPrefix(a.GoType(), "slices") {
				imps["github.com/gobuffalo/pop/slices"] = true
			}
		}
	}
	i := make([]string, 0, len(imps))
	for k := range imps {
		i = append(i, k)
	}
	return i
}
