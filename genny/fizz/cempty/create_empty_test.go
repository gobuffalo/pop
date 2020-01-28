package cempty

import (
	"testing"
	"time"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	t0, _ := time.Parse(time.RFC3339, "2019-08-28T07:46:02Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	g, err := New(&Options{
		Name: "create_widgets",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal("migrations/20190828074602_create_widgets.down.fizz", f.Name())
	r.Equal("", f.String())

	f = res.Files[1]
	r.Equal("migrations/20190828074602_create_widgets.up.fizz", f.Name())
	r.Equal("", f.String())
}
