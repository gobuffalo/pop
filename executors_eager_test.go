package pop_test

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
)

//Finished TODO Need to test Update has many update
func Test_Eager_Update_Has_Many_Create(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)
		count, _ := tx.Count(&User{})
		user := User{
			Name: nulls.NewString("Mark 'Awesome' Bates"),
		}

		err := tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, _ := tx.Count(&User{})
		a.Equal(count+1, ctx)

		//Set the Related Models
		user.Books = Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}}
		user.FavoriteSong = Song{Title: "Hook - Blues Traveler"}
		user.Houses = Addresses{
			Address{HouseNumber: 86, Street: "Modelo"},
		}

		err = tx.Eager().Update(&user)

		a.NoError(err)

		ctx, _ = tx.Count(&Book{})
		a.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Song{})
		a.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		a.Equal(count+1, ctx)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		a.NoError(err)

		a.Equal(u.Name.String, "Mark 'Awesome' Bates")
		a.Equal(u.Books[0].Title, "Pop Book")
		a.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		a.Equal(u.Houses[0].Street, "Modelo")
	})
}

// Finished TODO Need to test Validate Update has many update
func Test_Eager_Validate_And_Update_Has_Many_Create(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		count, _ := tx.Count(&User{})
		user := User{
			Name: nulls.NewString("Mark 'Awesome' Bates"),
		}

		err := tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, _ := tx.Count(&User{})
		a.Equal(count+1, ctx)

		//Set the Related Models
		user.Books = Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}}
		user.FavoriteSong = Song{Title: "Hook - Blues Traveler"}
		user.Houses = Addresses{
			Address{HouseNumber: 86, Street: "Modelo"},
		}

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		a.NoError(err)
		ctx, _ = tx.Count(&User{})
		a.Zero(ctx)
		a.Equal(1, verrs.Count()) // Missing Books.Description.
	})
}

func Test_Eager_Validate_And_Update_Parental(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := User{
			Name:         nulls.NewString(""),
			Books:        Books{{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		a.NoError(err)
		ctx, _ := tx.Count(&User{})
		a.Zero(ctx)
		a.Equal(1, verrs.Count()) // Missing Books.Description.
	})
}

func Test_Eager_Update_Belongs_To_Create(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)
		book := Book{
			Title:       "Pop Book",
			Description: "Pop Book",
			Isbn:        "PB1",
			User: User{
				Name: nulls.NewString("Larry"),
			},
		}

		err := tx.Eager().Create(&book)
		a.NoError(err)

		ctx, _ := tx.Count(&Book{})
		a.Equal(1, ctx)

		ctx, _ = tx.Count(&User{})
		a.Equal(1, ctx)

		car := Taxi{
			Model: "Fancy car",
			Driver: User{
				Name: nulls.NewString("Larry 2"),
			},
		}

		err = tx.Eager().Create(&car)
		a.NoError(err)

		ctx, _ = tx.Count(&Taxi{})
		a.Equal(1, ctx)

		err = tx.Eager().Find(&car, car.ID)
		a.NoError(err)

		a.Equal(nulls.NewString("Larry 2"), car.Driver.Name)
	})
}

//Finished
func Test_Eager_Update_Without_Associations(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)
		code := CourseCode{
			Course: Course{},
		}

		tx.Create(&code)

		c1 := code.Course.ID

		code.Course = Course{}

		err := tx.Eager().Update(&code)
		a.NoError(err)

		ctx, _ := tx.Count(&CourseCode{})
		a.Equal(1, ctx)
		a.NotEqual(c1, code.Course.ID)
	})
}
