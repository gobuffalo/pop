package ctable

import (
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		TableName: "widgets",
		Name:      "create_widgets.fizz",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("create_widgets.fizz", f.Name())
	r.Equal(`create_table("widgets") {
	t.Timestamps()
}`, f.String())
}

func Test_New_Fail(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		TableName: "",
		Name:      "create_widgets.fizz",
	})
	r.Error(err)
	r.NotEmpty(g)
}
