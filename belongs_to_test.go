package pop_test

import (
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
