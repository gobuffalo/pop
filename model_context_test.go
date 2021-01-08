package pop

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ContextTable struct {
	ID        string    `db:"id"`
	Value     string    `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (t ContextTable) TableName(ctx context.Context) string {
	// This is singular on purpose! It will checck if the TableName is properly
	// Respected in slices as well.
	return "context_prefix_" + ctx.Value("prefix").(string) + "_table"
}

func Test_ModelContext(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}

	t.Run("contextless", func(t *testing.T) {
		r := require.New(t)
		r.Panics(func() {
			var c ContextTable
			r.NoError(PDB.Create(&c))
		}, "panics if context prefix is not set")
	})

	for _, prefix := range []string{"a", "b"} {
		t.Run("prefix="+prefix, func(t *testing.T) {
			r := require.New(t)

			expected := ContextTable{ID: prefix, Value: prefix}
			c := PDB.WithContext(context.WithValue(context.Background(), "prefix", prefix))
			r.NoError(c.Create(&expected))

			var actual ContextTable
			r.NoError(c.Find(&actual, expected.ID))
			r.EqualValues(prefix, actual.Value)
			r.EqualValues(prefix, actual.ID)

			exists, err := c.Where("id = ?", actual.ID).Exists(new(ContextTable))
			r.NoError(err)
			r.True(exists)

			count, err := c.Where("id = ?", actual.ID).Count(new(ContextTable))
			r.NoError(err)
			r.EqualValues(1, count)

			expected.Value += expected.Value
			r.NoError(c.Update(&expected))

			r.NoError(c.Find(&actual, expected.ID))
			r.EqualValues(prefix+prefix, actual.Value)
			r.EqualValues(prefix, actual.ID)

			var results []ContextTable
			require.NoError(t, c.All(&results))

			require.NoError(t, c.First(&expected))
			require.NoError(t, c.Last(&expected))

			r.NoError(c.Destroy(&expected))
		})
	}

	t.Run("prefix=unknown", func(t *testing.T) {
		r := require.New(t)
		c := PDB.WithContext(context.WithValue(context.Background(), "prefix", "unknown"))
		err := c.Create(&ContextTable{ID: "unknown"})
		r.Error(err)

		if !strings.Contains(err.Error(), "context_prefix_unknown_table") { // All other databases
			t.Fatalf("Expected error to contain indicator that table does not exist but got: %s", err.Error())
		}
	})

	t.Run("cache_busting", func(t *testing.T) {
		r := require.New(t)

		var expectedA, expectedB ContextTable
		expectedA.ID = "expectedA"
		expectedB.ID = "expectedB"

		cA := PDB.WithContext(context.WithValue(context.Background(), "prefix", "a"))
		r.NoError(cA.Create(&expectedA))

		cB := PDB.WithContext(context.WithValue(context.Background(), "prefix", "b"))
		r.NoError(cB.Create(&expectedB))

		var actualA, actualB []ContextTable
		r.NoError(cA.All(&actualA))
		r.NoError(cB.All(&actualB))

		r.Len(actualA, 1)
		r.Len(actualB, 1)

		r.NotEqual(actualA[0].ID, actualB[0].ID, "if these are equal context switching did not work")
	})
}
