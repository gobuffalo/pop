package model

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/attrs"
)

// Options for generating a new model
type Options struct {
	Name                   string      `json:"name"`
	Attrs                  attrs.Attrs `json:"props"`
	Path                   string      `json:"path"`
	Package                string      `json:"package"`
	TestPackage            string      `json:"test_package"`
	Encoding               string      `json:"encoding"`
	ForceDefaultID         bool        `json:"force_default_id"`
	ForceDefaultTimestamps bool        `json:"force_default_timestamps"`
}

// Validate that options are usable
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
	if opts.Encoding != "json" && opts.Encoding != "jsonapi" && opts.Encoding != "xml" {
		return fmt.Errorf("unsupported encoding option %s", opts.Encoding)
	}

	return opts.forceDefaults()
}

func (opts *Options) forceDefaults() error {
	var idFound, createdAtFound, updatedAtFound bool
	for _, a := range opts.Attrs {
		switch a.Name.Underscore().String() {
		case "id":
			idFound = true
		case "created_at":
			createdAtFound = true
		case "updated_at":
			updatedAtFound = true
		}
	}
	// Add a default UUID, if no custom ID is provided
	if opts.ForceDefaultID && !idFound {
		id, err := attrs.Parse("id:uuid")
		if err != nil {
			return err
		}
		opts.Attrs = append([]attrs.Attr{id}, opts.Attrs...)
	}

	// Add default timestamp columns if they were not provided
	if opts.ForceDefaultTimestamps && !createdAtFound {
		createdAt, err := attrs.Parse("created_at:time")
		if err != nil {
			return err
		}
		opts.Attrs = append(opts.Attrs, createdAt)
	}
	if opts.ForceDefaultTimestamps && !updatedAtFound {
		updatedAt, err := attrs.Parse("updated_at:time")
		if err != nil {
			return err
		}
		opts.Attrs = append(opts.Attrs, updatedAt)
	}
	return nil
}
