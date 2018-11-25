package pop

import (
	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Eager_Update_Has_Many_Create(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		count2, _ := tx.Count(&User{})
		println("Count of users in Eager: ", count2)

		user := User{
			Name: nulls.NewString("Carl Lewis"),
			Books: Books{
				{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"},
			},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		err := tx.Eager().Create(&user)

		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		count := len(u.Houses)

		u.Houses = append(u.Houses, Address{HouseNumber: 43, Street: "Bryden"})

		tx.Eager().Update(&u)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(count+1, len(u2.Houses))

	})
}
