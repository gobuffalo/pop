package pop

import (
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/stretchr/testify/require"
)

func Test_New_Implementation_For_Nplus1(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			a.NoError(tx.Create(&user))

			book := Book{UserID: nulls.NewInt(user.ID)}
			a.NoError(tx.Create(&book))

			writer := Writer{Name: "Larry", BookID: book.ID}
			a.NoError(tx.Create(&writer))

			if name == "Mark" {
				song := Song{UserID: user.ID}
				a.NoError(tx.Create(&song))

				address := Address{Street: "Pop"}
				a.NoError(tx.Create(&address))

				home := UsersAddress{UserID: user.ID, AddressID: address.ID}
				a.NoError(tx.Create(&home))
			}
		}

		users := []User{}
		a.NoError(tx.All(&users))

		// FILL THE HAS-MANY and HAS_ONE
		a.NoError(preload(tx, &users))

		a.Len(users[0].Books, 1)
		a.Len(users[1].Books, 1)
		a.Len(users[2].Books, 1)
		a.Equal(users[0].FavoriteSong.UserID, users[0].ID)
		a.Len(users[0].Houses, 1)

		book := Book{}
		a.NoError(tx.First(&book))
		a.NoError(preload(tx, &book))
		a.Len(book.Writers, 1)
		a.Equal("Larry", book.Writers[0].Name)
		a.Equal("Mark", book.User.Name.String)
	})
}

func Test_New_Implementation_For_Nplus1_With_UUID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)

		courses := []Course{}
		for i := 0; i < 3; i++ {
			course := Course{}
			a.NoError(tx.Create(&course))
			courses = append(courses, course)
			if i == 0 {
				a.NoError(tx.Create(&CourseCode{
					CourseID: course.ID,
				}))
			}
		}

		courseCodes := []CourseCode{}
		a.NoError(tx.All(&courseCodes))
		a.Len(courseCodes, 1)

		// FILL THE HAS-MANY and HAS_ONE
		a.NoError(preload(tx, &courseCodes))
		a.Equal(courses[0].ID, courseCodes[0].Course.ID)

		student := Student{}
		a.NoError(tx.Create(&student))

		parent := Parent{}
		a.NoError(tx.Create(&parent))

		a.NoError(tx.RawQuery("insert into parents_students(parent_id, student_id) values(?,?)", parent.ID.String(), student.ID.String()).Exec())

		parents := []Parent{}
		a.NoError(tx.All(&parents))

		a.NoError(preload(tx, &parents))
		a.Len(parents, 1)
		a.Len(parents[0].Students, 1)
		a.Equal(student.ID, parents[0].Students[0].ID)
	})
}

func Test_New_Implementation_For_Nplus1_Single(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			a.NoError(tx.Create(&user))

			book := Book{UserID: nulls.NewInt(user.ID)}
			a.NoError(tx.Create(&book))

			writer := Writer{Name: "Larry", BookID: book.ID}
			a.NoError(tx.Create(&writer))

			if name == "Mark" {
				song := Song{UserID: user.ID}
				a.NoError(tx.Create(&song))

				address := Address{Street: "Pop"}
				a.NoError(tx.Create(&address))

				home := UsersAddress{UserID: user.ID, AddressID: address.ID}
				a.NoError(tx.Create(&home))
			}
		}

		users := []User{}
		a.NoError(tx.All(&users))

		// FILL THE HAS-MANY and HAS_ONE
		a.NoError(preload(tx, &users, "Books"))

		a.Len(users[0].Books, 1)
		a.Len(users[1].Books, 1)
		a.Len(users[2].Books, 1)
		a.Zero(users[0].FavoriteSong.UserID)
		a.Len(users[0].Houses, 0)
	})
}

func Test_New_Implementation_For_Nplus1_Nested(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		var song Song
		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			a.NoError(tx.Create(&user))

			book := Book{UserID: nulls.NewInt(user.ID)}
			a.NoError(tx.Create(&book))

			if name == "Mark" {
				song = Song{UserID: user.ID}
				a.NoError(tx.Create(&song))

				address := Address{Street: "Pop"}
				a.NoError(tx.Create(&address))

				home := UsersAddress{UserID: user.ID, AddressID: address.ID}
				a.NoError(tx.Create(&home))
			}
		}

		SetEagerMode(EagerPreload)
		users := []User{}
		a.NoError(tx.Eager("Houses", "Books", "Books.User.FavoriteSong").All(&users))
		a.Len(users[0].Books, 1)
		a.Len(users[1].Books, 1)
		a.Len(users[2].Books, 1)
		a.Len(users[0].Houses, 1)

		a.Equal(users[0].ID, users[0].Books[0].User.ID)
		a.Equal(song.ID, users[0].Books[0].User.FavoriteSong.ID)
		SetEagerMode(EagerDefault)
	})
}

func Test_New_Implementation_For_Nplus1_BelongsTo_Not_Underscore(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		user := User{Name: nulls.NewString("Mark")}
		a.NoError(tx.Create(&user))

		taxi := Taxi{UserID: nulls.NewInt(user.ID)}
		a.NoError(tx.Create(&taxi))

		SetEagerMode(EagerPreload)
		taxis := []Taxi{}
		a.NoError(tx.EagerPreload().All(&taxis))
		a.Len(taxis, 1)
		a.Equal("Mark", taxis[0].Driver.Name.String)
		SetEagerMode(EagerDefault)
	})
}

func Test_New_Implementation_For_BelongsTo_Multiple_Fields(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		user := User{Name: nulls.NewString("Mark")}
		a.NoError(tx.Create(&user))

		address := Address{HouseNumber: 2, Street: "Street One"}
		a.NoError(tx.Create(&address))

		taxi := Taxi{UserID: nulls.NewInt(user.ID), AddressID: nulls.NewInt(address.ID)}
		a.NoError(tx.Create(&taxi))

		book := Book{TaxiID: nulls.NewInt(taxi.ID), Title: "My Book"}
		a.NoError(tx.Create(&book))

		SetEagerMode(EagerPreload)
		books := []Book{}
		a.NoError(tx.EagerPreload("Taxi.Driver", "Taxi.Address").All(&books))
		a.Len(books, 1)
		a.Equal(user.Name.String, books[0].Taxi.Driver.Name.String)
		a.Equal(address.Street, books[0].Taxi.Address.Street)
		SetEagerMode(EagerDefault)
	})
}

func Test_New_Implementation_For_BelongsTo_Ptr_Field(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		a := require.New(t)
		toAddress := Address{HouseNumber: 1, Street: "Destination Ave"}
		a.NoError(tx.Create(&toAddress))

		taxi := Taxi{ToAddressID: &toAddress.ID}
		a.NoError(tx.Create(&taxi))

		book1 := Book{TaxiID: nulls.NewInt(taxi.ID), Title: "My Book"}
		a.NoError(tx.Create(&book1))

		taxiNilToAddress := Taxi{ToAddressID: nil}
		a.NoError(tx.Create(&taxiNilToAddress))

		book2 := Book{TaxiID: nulls.NewInt(taxiNilToAddress.ID), Title: "Another Book"}
		a.NoError(tx.Create(&book2))

		SetEagerMode(EagerPreload)
		books := []Book{}
		a.NoError(tx.EagerPreload("Taxi.ToAddress").Order("created_at").All(&books))
		a.Len(books, 2)
		a.Equal(toAddress.Street, books[0].Taxi.ToAddress.Street)
		a.NotNil(books[0].Taxi.ToAddressID)
		a.Nil(books[1].Taxi.ToAddress)
		a.Nil(books[1].Taxi.ToAddressID)
		SetEagerMode(EagerDefault)
	})
}
