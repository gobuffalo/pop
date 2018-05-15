package pop_test

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/csv"
	"github.com/stretchr/testify/require"
)

func Test_CSV_Import(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		composer := &Composer{}
		imp := csv.NewImporter(tx)
		r.NoError(imp.FromPath("./csv/files/composers.csv", composer))

		c, err := tx.Count(composer)
		r.NoError(err)
		r.Equal(5, c)
	})
}
