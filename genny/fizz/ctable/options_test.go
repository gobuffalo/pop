package ctable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.TableName = "widget"

	err = opts.Validate()
	r.NoError(err)

	r.Contains(opts.Name, "_create_widgets.fizz")
	r.Equal("migrations", opts.Path)
}
