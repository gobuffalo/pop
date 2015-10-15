package pop_test

import (
	"os"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Scopes(t *testing.T) {
	r := require.New(t)
	oql := "SELECT name as full_name, users.alive, users.bio, users.birth_date, users.created_at, users.id, users.name, users.price, users.updated_at FROM users AS users"

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
