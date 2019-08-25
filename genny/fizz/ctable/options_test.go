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

	r.Contains(opts.Name, "_create_widgets")
	r.Equal("migrations", opts.Path)
}

func Test_Options_Validate_Errors(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		TableName: "widget",
		Type:      "sql",
	}
	err := opts.Validate()
	r.EqualError(err, "sql migrations require a fizz translator")

	opts.Translator = mockTranslator{}
	opts.Type = "aaa"
	err = opts.Validate()
	r.EqualError(err, "aaa migration type is not allowed")
}
