package config

import (
	"os"

	"errors"
)

// Options needed for the config generator
type Options struct {
	Root     string // Defaults to PWD
	FileName string // Defaults to database.yml
	Dialect  string // required
	Prefix   string // required - <prefix>_development
}

func (opts *Options) Validate() error {
	if opts.Root == "" {
		pwd, _ := os.Getwd()
		opts.Root = pwd
	}
	if opts.Prefix == "" {
		return errors.New("you must provide a database name prefix")
	}
	if opts.FileName == "" {
		opts.FileName = "database.yml"
	}
	if opts.Dialect == "" {
		return errors.New("you must provide a database dialect")
	}
	return nil
}
