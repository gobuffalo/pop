package pop

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortingMigrations(t *testing.T) {
	examples := Migrations{
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

	t.Run("case=enforces precedence for specific up migrations", func(t *testing.T) {
		migrations := make(Migrations, len(examples))
		copy(migrations, examples)

		expectedOrder := Migrations{
			examples[1],
			examples[0],
			examples[2],
			examples[3],
			examples[5],
			examples[4],
		}

		sort.Sort(UpMigrations{migrations})

		assert.Equal(t, expectedOrder, migrations)
	})

	t.Run("case=enforces precedence for specific down migrations", func(t *testing.T) {
		migrations := make(Migrations, len(examples))
		copy(migrations, examples)

		expectedOrder := Migrations{
			examples[5],
			examples[4],
			examples[2],
			examples[3],
			examples[1],
			examples[0],
		}

		sort.Sort(DownMigrations{migrations})

		assert.Equal(t, expectedOrder, migrations)
	})
}
