package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Assoc_HasMany(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		pop.Debug = true
		r := require.New(t)

		g := &Game{}
		r.NoError(tx.Create(g))

		r.NoError(tx.Create(&Player{GameID: g.ID}))
		r.NoError(tx.Create(&Player{GameID: g.ID}))
		r.NoError(tx.Create(&Player{GameID: g.ID}))

		r.NoError(tx.With("Players").Find(g, g.ID))
		r.Len(g.Players, 3)
	})
}
