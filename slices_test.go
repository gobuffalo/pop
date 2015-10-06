package pop_test

import (
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

type Cake struct {
	IntSlice    pop.IntSlice    `db:"int_slice"`
	FloatSlice  pop.FloatSlice  `db:"float_slice"`
	StringSlice pop.StringSlice `db:"string_slice"`
}

func Test_StringSlice(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				StringSlice: pop.StringSlice{"a", "b", "c"},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(pop.StringSlice{"a", "b", "c"}, c.StringSlice)
		})
	} else {
		t.SkipNow()
	}
}

func Test_IntSlice(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				IntSlice: pop.IntSlice{1, 2, 3},
				// FloatSlice: pop.FloatSlice{1.0, 2.1, 3.2},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(pop.IntSlice{1, 2, 3}, c.IntSlice)
		})
	} else {
		t.SkipNow()
	}
}

func Test_FloatSlice(t *testing.T) {
	if os.Getenv("SODA_DIALECT") == "postgres" {
		transaction(func(tx *pop.Connection) {
			r := require.New(t)

			c := &Cake{
				FloatSlice: pop.FloatSlice{1.0, 2.1, 3.2},
			}
			err := tx.Create(c)
			r.NoError(err)

			err = tx.Reload(c)
			r.Equal(pop.FloatSlice{1.0, 2.1, 3.2}, c.FloatSlice)
		})
	} else {
		t.SkipNow()
	}
}
