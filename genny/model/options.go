package model

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/attrs"
	"github.com/pkg/errors"
)

// Options for generating a new model
type Options struct {
	Name        string      `json:"name"`
	Attrs       attrs.Attrs `json:"props"`
	Path        string      `json:"path"`
	Package     string      `json:"package"`
	TestPackage string      `json:"test_package"`
	Encoding    string      `json:"encoding"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must set a name for your model")
	}
	if len(opts.Path) == 0 {
		opts.Path = "models"
	}
	if len(opts.Package) == 0 {
		opts.Package = filepath.Base(opts.Path)
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
