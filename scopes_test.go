package pop

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scopes(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	oql := "SELECT enemies.A FROM enemies AS enemies"

	m := NewModel(new(Enemy), context.Background())

	q := PDB.Q()
	s, _ := q.ToSQL(m)
	r.Equal(oql, s)

	q.Scope(func(qy *Query) *Query {
		return qy.Where("id = ?", 1)
	})

	s, _ = q.ToSQL(m)
	r.Equal(ts(oql+" WHERE id = ?"), s)
}
