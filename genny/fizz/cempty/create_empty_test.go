package cempty

import (
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		TableName: "widgets",
		Name:      "create_widgets",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("migrations/create_widgets.down.fizz", f.Name())
	r.Equal("", f.String())

	f = res.Files[1]
	r.Equal("migrations/create_widgets.up.fizz", f.Name())
	r.Equal("", f.String())
}
