package config

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v6"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	for _, d := range pop.AvailableDialects {
		run := genny.DryRunner(context.Background())

		g, err := New(&Options{
			Prefix:  "foo",
			Dialect: d,
		})

		r.NoError(err)
		run.With(g)
		r.NoError(run.Run())

		res := run.Results()
		r.Len(res.Commands, 0)
		r.Len(res.Files, 1)

		f := res.Files[0]
		r.Equal("database.yml", filepath.Base(f.Name()))
		body := f.String()
		r.Contains(body, d)
		r.Contains(body, "foo_development")
		r.Contains(body, "foo_production")
		r.Contains(body, "foo_test")
	}
}

func Test_New_No_Dialect(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{
		Prefix: "foo",
	})

	r.Error(err)
}

func Test_New_No_Prefix(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{
		Dialect: "postgres",
	})

	r.Error(err)
}

func Test_New_BadDialect(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{
		Prefix:  "foo",
		Dialect: "unknown",
	})

	r.EqualError(err, "unable to find database.yml template for dialect unknown")
}
