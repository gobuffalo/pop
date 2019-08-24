package ctable

import (
	"testing"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("id:uuid", "created_at:timestamp", "updated_at:timestamp", "name", "description:text", "age:int", "bar:nulls.String")
	r.NoError(err)

	cases := []struct {
		Options *Options
		Result  string
	}{
		{
			&Options{
				TableName: "widgets",
				Name:      "create_widgets.fizz",
			},
			`create_table("widgets") {
	t.Timestamps()
}`,
		},
		{
			&Options{
				TableName: "widgets",
				Name:      "create_widgets.fizz",
				Attrs:     ats,
			},
			`create_table("widgets") {
	t.Column("id", "uuid", {primary: true})
	t.Column("name", "string", {})
	t.Column("description", "text", {})
	t.Column("age", "integer", {})
	t.Column("bar", "string", {null: true})
	t.Timestamps()
}`,
		},
	}

	for _, c := range cases {
		g, err := New(c.Options)
		r.NoError(err)

		run := gentest.NewRunner()
		run.With(g)

		r.NoError(run.Run())

		res := run.Results()

		r.Len(res.Commands, 0)
		r.Len(res.Files, 1)

		f := res.Files[0]
		r.Equal("create_widgets.fizz", f.Name())
		r.Equal(c.Result, f.String())
	}
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
