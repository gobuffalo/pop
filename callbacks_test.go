package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Callbacks(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		user := &CallbacksUser{
			BeforeS: "BS",
			BeforeC: "BC",
			BeforeU: "BU",
			BeforeD: "BD",
			AfterS:  "AS",
			AfterC:  "AC",
			AfterU:  "AU",
			AfterD:  "AD",
			AfterF:  "AF",
		}

		r.NoError(tx.Save(user))

		r.Equal("BeforeSave", user.BeforeS)
		r.Equal("BeforeCreate", user.BeforeC)
		r.Equal("AfterSave", user.AfterS)
		r.Equal("AfterCreate", user.AfterC)
		r.Equal("BU", user.BeforeU)
		r.Equal("AU", user.AfterU)

		r.NoError(tx.Update(user))

		r.Equal("BeforeUpdate", user.BeforeU)
		r.Equal("AfterUpdate", user.AfterU)
		r.Equal("BD", user.BeforeD)
		r.Equal("AD", user.AfterD)

		r.Equal("AF", user.AfterF)
		r.NoError(tx.Find(user, user.ID))
		r.Equal("AfterFind", user.AfterF)

		r.NoError(tx.Destroy(user))

		r.Equal("BeforeDestroy", user.BeforeD)
		r.Equal("AfterDestroy", user.AfterD)

	})
}

func Test_Callbacks_on_Slice(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)
		for i := 0; i < 2; i++ {
			r.NoError(tx.Create(&CallbacksUser{}))
		}

		users := CallbacksUsers{}
		r.NoError(tx.All(&users))

		r.Len(users, 2)

		for _, u := range users {
			r.Equal("AfterFind", u.AfterF)
		}
	})
}
