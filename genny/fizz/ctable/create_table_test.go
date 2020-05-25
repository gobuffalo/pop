package ctable

import (
	"testing"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("id:uuid", "created_at:timestamp", "updated_at:timestamp", "name", "description:text", "age:int", "bar:nulls.String", "started_at:time.Time", "finished_at:nulls.Time")
	r.NoError(err)

	cases := []struct {
		Options *Options
		Result  string
	}{
		{
			&Options{
				TableName:              "widgets",
				Name:                   "create_widgets",
				ForceDefaultTimestamps: true,
			},
			`create_table("widgets") {
	t.Timestamps()
}`,
		},
		{
			&Options{
				TableName:              "widget",
				Name:                   "create_widgets",
				ForceDefaultTimestamps: true,
			},
			`create_table("widgets") {
	t.Timestamps()
}`,
		},
		{
			&Options{
				TableName:              "widgets",
				Name:                   "create_widgets",
				Attrs:                  ats,
				ForceDefaultTimestamps: true,
			},
			`create_table("widgets") {
	t.Column("id", "uuid", {primary: true})
	t.Column("name", "string", {})
	t.Column("description", "text", {})
	t.Column("age", "integer", {})
	t.Column("bar", "string", {null: true})
	t.Column("started_at", "timestamp", {})
	t.Column("finished_at", "timestamp", {null: true})
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
		r.Len(res.Files, 2)

		f := res.Files[0]
		r.Equal("migrations/create_widgets.down.fizz", f.Name())
		r.Equal(`drop_table("widgets")`, f.String())

		f = res.Files[1]
		r.Equal("migrations/create_widgets.up.fizz", f.Name())
		r.Equal(c.Result, f.String())
	}
}

func Test_New_SQL(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("id:uuid", "created_at:timestamp", "updated_at:timestamp", "name", "description:text", "age:int", "bar:nulls.String")
	r.NoError(err)

	g, err := New(&Options{
		TableName:  "widgets",
		Name:       "create_widgets",
		Type:       "sql",
		Translator: mockTranslator{},
		Attrs:      ats,
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("migrations/create_widgets.test.down.sql", f.Name())
	r.Equal("drop table;", f.String())

	f = res.Files[1]
	r.Equal("migrations/create_widgets.test.up.sql", f.Name())
	r.Equal("create table;", f.String())
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
