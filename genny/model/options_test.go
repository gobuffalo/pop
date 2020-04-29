package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.Name = "widget"

	err = opts.Validate()
	r.NoError(err)
	r.Equal(0, len(opts.Attrs))
	r.Equal("models", opts.Path)
	r.Equal("models", opts.Package)
	r.Equal("models", opts.TestPackage)
	r.Equal("json", opts.Encoding)

	// Force ID opt
	opts.ForceDefaultID = true
	err = opts.Validate()
	r.NoError(err)
	r.Equal(1, len(opts.Attrs))
	r.Equal("id", opts.Attrs[0].Name.String())

	// Force default timestamps
	opts.ForceDefaultTimestamps = true
	err = opts.Validate()
	r.NoError(err)

	r.Equal(3, len(opts.Attrs))
	r.Equal("created_at", opts.Attrs[1].Name.String())
	r.Equal("updated_at", opts.Attrs[2].Name.String())
}
