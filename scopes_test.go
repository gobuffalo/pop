package pop_test

import (
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Scopes(t *testing.T) {
	r := require.New(t)
	oql := "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users as users"

	transaction(func(tx *pop.Connection) {
		u := &pop.Model{Value: &User{}}
		q := tx.Q()

		s, _ := q.ToSQL(u)
		r.Equal(oql, s)

		q.Scope(func(qy *pop.Query) *pop.Query {
			return qy.Where("id = ?", 1)
		})

		s, _ = q.ToSQL(u)
		if os.Getenv("SODA_DIALECT") == "postgres" {
			r.Equal(oql+" WHERE id = $1", s)
		} else {
			r.Equal(oql+" WHERE id = ?", s)
		}
	})
}
