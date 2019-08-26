package pop

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/pop/nulls"
)

func Test_Eager_Update_Many_Many_Create(t *testing.T) {

	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		count2, _ := tx.Count(&User{})
		println("Count of user s in database: ", count2)

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

		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(count+1, len(u2.Houses))

	})
}

func Test_Eager_Update_Has_Many_Add_Existing(t *testing.T) {

	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
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

		// Create Book

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
		err = tx.Eager().Update(&u)
		r.NoError(err)

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

	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
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

		// Create Book

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		// update Address

		address := u.Houses[0]

		address.HouseNumber = 43

		u.Houses[0] = address

		// Update user
		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(u2.Houses[0].HouseNumber, 43)

	})
}

func Test_Eager_Update_Many_2_Many_Update_Existing(t *testing.T) {

	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
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

		// Create Address

		addy2 := Address{HouseNumber: 33, Street: "Broad"}

		err = tx.Create(&addy2)

		r.NoError(err)

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		u.Houses = append(u.Houses, addy2)

		// Update user
		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		count := len(u.Houses)

		r.NoError(err)
		r.Equal(count, len(u2.Houses))

	})
}

func Test_Eager_Update_Has_One(t *testing.T) {

	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
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
		count, _ := tx.Count(&Songs{})

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q.First(&u)

		// update Song

		u.FavoriteSong = Song{Title: "Body - Brando"}

		// Update user
		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)

		songs := Songs{}

		numSongs, _ := tx.Count(&songs)

		r.NoError(err)
		r.Equal(count+1, numSongs)

		// Favorite Song should equal "Body - Brando"
		r.Equal("Body - Brando", u2.FavoriteSong.Title)

	})
}

func Test_Eager_Update_Many_To_Many(t *testing.T) {
	t.Skip("skipping broken test")
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {

		user := User{
			Name: nulls.NewString("Carl Lewis2"),
			Books: Books{
				{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"},
			},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		// Create User
		err := tx.Eager().Create(&user)

		u := User{}
		q := tx.Eager().Where("name = ?", "Carl Lewis2")
		err = q.First(&u)

		count := len(u.Houses)

		// Add two houses
		u.Houses = append(u.Houses, Address{HouseNumber: 43, Street: "Bryden"})
		u.Houses = append(u.Houses, Address{HouseNumber: 455, Street: "Broad"})

		// Remove the first house

		u.Houses = u.Houses[1:]

		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis2")
		err = q2.First(&u2)

		r.NoError(err)
		r.Equal(count+1, len(u2.Houses))

	})
}

func Test_Eager_Update_Has_Many_Transfer(t *testing.T) {
	// t.Skip("I am skipping this test")
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		// Create Users
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

		user2 := User{
			Name:         nulls.NewString("Carl McKinney"),
			FavoriteSong: Song{Title: "Anything - Goldlink"},
			Houses: Addresses{
				{HouseNumber: 105, Street: "Jump Street"},
			},
		}

		err := tx.Eager().Create(&user)
		r.NoError(err)

		err = tx.Eager().Create(&user2)
		r.NoError(err)

		// Change book owner

		// Find user
		u := User{}
		q := tx.Eager().Where("name = ?", "Carl McKinney")
		err = q.First(&u)

		u.Books = user.Books

		// Update user
		err = tx.Eager().Update(&u)
		r.NoError(err)

		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)
		r.NoError(err)

		u3 := User{}
		q3 := tx.Eager().Where("name = ?", "Carl McKinney")
		err = q3.First(&u3)
		r.NoError(err)

		r.Equal(0, len(u2.Books))
		r.Equal(1, len(u3.Books))

		//	Book UserID should equal u3.ID

		book := Book{}

		err = tx.Where("title = ?", "Pop Book").First(&book)

		r.NoError(err)

		r.Equal(u3.ID, book.UserID.Int)

	})
}

func Test_Eager_Update_Belongs_To(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		// Create Users
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

		user2 := User{
			Name: nulls.NewString("Bran Stark"),
		}

		err := tx.Eager().Create(&user)
		err = tx.Eager().Create(&user2)
		r.NoError(err)

		// Find Book
		book := Book{}
		err = tx.Where("title = ?", "Pop Book").First(&book)
		r.NoError(err)

		// Change Attribute of book owner
		book.User = user
		book.User.Alive = nulls.NewBool(true)

		// Update book
		err = tx.Eager().Update(&book)
		r.NoError(err)

		// Find the book's user directly
		u2 := User{}
		q2 := tx.Eager().Where("name = ?", "Carl Lewis")
		err = q2.First(&u2)
		r.NoError(err)

		// assert that the user has been changes to alive from dead
		r.Equal(true, u2.Alive.Bool)

	})
}
