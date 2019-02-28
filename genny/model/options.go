package model

import (
	"strings"

	"github.com/gobuffalo/attrs"
	"github.com/pkg/errors"
)

type Options struct {
	Name        string      `json:"name"`
	Attrs       attrs.Attrs `json:"props"`
	Package     string      `json:"package"`
	TestPackage string      `json:"test_package"`
	Encoding    string      `json:"encoding"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must set a name for your model")
	}
	if len(opts.Package) == 0 {
		opts.Package = "models"
	}
	if len(opts.TestPackage) == 0 {
		opts.TestPackage = opts.Package
	}
	if len(opts.Encoding) == 0 {
		opts.Encoding = "json"
	}
	opts.Encoding = strings.ToLower(opts.Encoding)
	return nil
}
