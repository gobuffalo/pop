package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scopes(t *testing.T) {
	r := require.New(t)
	oql := "SELECT enemies.A FROM enemies AS enemies"

	m := &Model{Value: &Enemy{}}

	q := PDB.Q()
	s, _ := q.ToSQL(m)
	r.Equal(oql, s)

	q.Scope(func(qy *Query) *Query {
		return qy.Where("id = ?", 1)
	})

	s, _ = q.ToSQL(m)
	r.Equal(ts(oql+" WHERE id = ?"), s)
}
