package pop

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortingMigrations(t *testing.T) {
	t.Run("case=enforces precedence for specific migrations", func(t *testing.T) {
		migrations := Migrations{
			{
				Version: "1",
				DBType:  "all",
			},
			{
				Version: "1",
				DBType:  "postgres",
			},
			{
				Version: "2",
				DBType:  "cockroach",
			},
			{
				Version: "2",
				DBType:  "all",
			},
			{
				Version: "3",
				DBType:  "all",
			},
			{
				Version: "3",
				DBType:  "mysql",
			},
		}
		expectedOrder := Migrations{
			migrations[1],
			migrations[0],
			migrations[2],
			migrations[3],
			migrations[5],
			migrations[4],
		}

		sort.Sort(migrations)

		assert.Equal(t, expectedOrder, migrations)
	})
}
