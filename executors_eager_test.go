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

func Test_Eager_Update_Has_Many_Add_Existing(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		// Create User
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

		//Create Book

		book := Book{Title: "The Life of Pi", Description: "Fiction", Isbn: "PB2"}

		err = tx.Create(&book)

		r.NoError(err)

		count2, _ := tx.Count(&Book{})
		println("Count of Books: ", count2)

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		count, _ := tx.Count(&Book{})

		// Add book
		u.Books = append(u.Books, book)

		// Update user
		tx.Eager().Update(&u)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(count, len(u2.Books))
		for _, value := range u2.Books {

			r.Equal(value.UserID.Int, u2.ID)
		}

	})
}

func Test_Eager_Update_Has_Many_Update_Existing(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		// Create User
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

		//Create Book

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		count, _ := tx.Count(&Book{})

		// update Address

		address := u.Houses[0]

		address.HouseNumber = 43

		u.Houses[0] = address

		// Update user
		tx.Eager().Update(&u)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(count, len(u2.Books))
		r.Equal(u2.Houses[0].HouseNumber, 43)

	})
}
