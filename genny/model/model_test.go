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

func Test_New_Standard(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("id:uuid", "created_at:timestamp", "updated_at:timestamp", "name", "description:text", "age:int", "bar:nulls.String")
	r.NoError(err)
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

func Test_New_No_Attrs(t *testing.T) {
	r := require.New(t)
	g, err := New(&Options{
		Name: "widget",
	})

	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	f, err := res.Find("models/widget.go")
	r.NoError(err)

	tf := gogen.FmtTransformer()
	f, err = tf.Transform(f)
	r.NoError(err)

	box := packr.New("Test_New_No_Attrs", "../model/_fixtures")
	bf, err := box.FindString("models/widget_empty.go")
	r.NoError(err)
	r.Equal(bf, f.String())
}

func Test_New_XML(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("id:uuid", "created_at:timestamp", "updated_at:timestamp", "name", "description:text", "age:int", "bar:nulls.String")
	r.NoError(err)
	g, err := New(&Options{
		Name:     "widget",
		Encoding: "xml",
		Attrs:    ats,
	})

	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.NoError(gentest.CompareFiles([]string{"models/widget.go", "models/widget_test.go"}, res.Files))

	box := packr.New("Test_New_XML", "../model/_fixtures")

	f, err := res.Find("models/widget_test.go")
	r.NoError(err)
	bf, err := box.FindString(f.Name())
	r.NoError(err)
	r.Equal(bf, f.String())

	f, err = res.Find("models/widget.go")

	tf := gogen.FmtTransformer()
	f, err = tf.Transform(f)
	r.NoError(err)

	bf, err = box.FindString("models/widget_xml.go")
	r.NoError(err)
	r.Equal(bf, f.String())
}
