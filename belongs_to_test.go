package pop_test

import (
	"os"
	"strings"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_BelongsTo(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		q := tx.Q()
		q = q.BelongsTo(&User{ID: 1})
		cls := q.WhereClauses
		r.Equal(1, len(cls))

		cl := cls[0]
		r.Equal("user_id = ?", cl.Fragment)
		r.Equal(cl.Arguments, []interface{}{1})
	})
}

func Test_BelongsToThrough(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		q := tx.BelongsToThrough(&User{ID: 1}, &Friend{})
		qs := "SELECT enemies.A FROM enemies AS enemies, good_friends AS good_friends WHERE good_friends.user_id = ? AND enemies.id = good_friends.enemy_id"
		if os.Getenv("SODA_DIALECT") == "postgres" {
			qs = strings.Replace(qs, "?", "$1", -1)
		}

		sql, args := q.ToSQL(&pop.Model{Value: &Enemy{}})
		r.Equal(qs, sql)
		r.Equal([]interface{}{1}, args)
	})
}
