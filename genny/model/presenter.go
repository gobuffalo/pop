package model

import (
	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/flect/name"
)

type presenter struct {
	Name        name.Ident
	Encoding    name.Ident
	Imports     []string
	Validations attrs.Attrs
}
