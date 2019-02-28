package model

import (
	"testing"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/gogen"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name: "widget",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	r.NoError(gentest.CompareFiles([]string{"models/widget.go", "models/widget_test.go"}, res.Files))
}

// b g r widget name description:text age:int bar:nulls.String

func Test_New_Standard(t *testing.T) {
	r := require.New(t)

	ats := attrs.Attrs{}
	for _, x := range []string{"id:uuid", "created_at:timestamp", "updated_at:timestamp"} {
		a, err := attrs.Parse(x)
		r.NoError(err)
		ats = append(ats, a)
	}
	xats, err := attrs.ParseArgs("name", "description:text", "age:int", "bar:nulls.String")
	r.NoError(err)
	ats = append(ats, xats...)
	g, err := New(&Options{
		Name:  "widget",
		Attrs: ats,
	})

	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.NoError(gentest.CompareFiles([]string{"models/widget.go", "models/widget_test.go"}, res.Files))

	box := packr.New("Test_New_Standard", "../model/_fixtures")

	f, err := res.Find("models/widget_test.go")
	r.NoError(err)
	bf, err := box.FindString(f.Name())
	r.NoError(err)
	r.Equal(bf, f.String())

	f, err = res.Find("models/widget.go")

	tf := gogen.FmtTransformer()
	f, err = tf.Transform(f)
	r.NoError(err)

	bf, err = box.FindString(f.Name())
	r.NoError(err)
	r.Equal(bf, f.String())

}
