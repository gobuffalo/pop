package ctable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	t0, _ := time.Parse(time.RFC3339, "2019-08-28T07:46:02Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.TableName = "widget"

	err = opts.Validate()
	r.NoError(err)

	r.Equal(opts.Name, "20190828074602_create_widgets")
	r.Equal("migrations", opts.Path)

	// Custom migration name
	opts.Name = "custom_migration"
	err = opts.Validate()
	r.NoError(err)

	r.Equal(opts.Name, "custom_migration")
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
