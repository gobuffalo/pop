package pop

import (
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Find(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		tCases := []string{
			"Mark",
			"💩",
		}

		for _, tc := range tCases {
			user := User{Name: nulls.NewString(tc)}
			err := tx.Create(&user)
			r.NoError(err)

			u := User{}
			err = tx.Find(&u, user.ID)
			r.NoError(err)

			r.NotEqual(u.ID, 0)
			r.Equal(u.Name.String, tc)
		}
	})
}

func Test_Create_MissingID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		client := Client{ClientID: "client-0001"}
		err := tx.Create(&client)
		r.Error(err)
		r.Contains(err.Error(), "model *pop.Client is missing required field ID")
	})
}

func Test_Find_MissingID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		r.NoError(tx.RawQuery("INSERT INTO clients (id) VALUES (?)", "client-0001").Exec())

		u := Client{}
		r.EqualError(tx.Find(&u, "client-0001"), "model *pop.Client is missing required field ID")
	})
}

func Test_Find_LeadingZeros(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		labels := []string{
			"0",
			"01",
			"001",
			"123",
		}

		for _, v := range labels {
			label := Label{ID: v}
			err := tx.Create(&label)
			r.NoError(err)

			l := Label{}
			err = tx.Find(&l, v)
			r.NoError(err)

			r.Equal(l.ID, v)
		}
	})
}

func Test_Select(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark"), Email: "mark@gobuffalo.io"}
		err := tx.Create(&user)
		r.NoError(err)

		q := tx.Select("name", "email", "\n", "\t\n", "")

		sm := NewModel(new(User), tx.Context())
		sql, _ := q.ToSQL(sm)
		r.Equal(tx.Dialect.TranslateSQL("SELECT email, name FROM users AS users"), sql)

		u := User{}
		err = q.Find(&u, user.ID)
		r.NoError(err)

		r.Equal(u.Email, "mark@gobuffalo.io")
		r.Equal(u.Name.String, "Mark")
		r.Zero(u.ID)
	})
}

func Test_Find_Eager_Has_Many(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		u := User{}
		err = tx.Eager("Books").Find(&u, user.ID)
		r.NoError(err)

		r.NotEqual(u.ID, 0)
		r.Equal(u.Name.String, "Mark")
		books := u.Books
		r.NotEqual(len(books), 0)
		r.Equal(books[0].Title, book.Title)
	})
}

func Test_All_Eager_Preload_Mode(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			r.NoError(err)

			if name == "Mark" {
				for _, title := range []string{"Pop Book", "My Pop Book", "Asociations Book"} {
					book := Book{Title: title, Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
					err = tx.Create(&book)
					r.NoError(err)
				}
			}
		}

		u := Users{}
		err := tx.EagerPreload().All(&u)
		r.NoError(err)
		r.Equal(len(u), 3)
		r.Equal(len(u[0].Books), 3)
		r.Equal(u[0].Books[0].Title, "Asociations Book")
		r.Equal(u[0].Books[1].Title, "My Pop Book")
		r.Equal(u[0].Books[2].Title, "Pop Book")
	})
}

func Test_Find_Eager_Has_Many_Order_By(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book1 := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book1)
		r.NoError(err)

		book2 := Book{Title: "New Pop Book", Isbn: "PB2", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book2)
		r.NoError(err)

		u := User{}
		err = tx.Eager().Find(&u, user.ID)
		r.NoError(err)

		r.Equal(len(u.Books), 2)
		r.Equal(book2.Title, u.Books[0].Title)
	})
}

func Test_Find_Eager_Belongs_To(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark"), UserName: "mark"}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		b := Book{}
		err = tx.Eager().Find(&b, book.ID)
		r.NoError(err)

		r.NotEqual(b.ID, 0)
		r.NotEqual(b.User.ID, 0)
		r.Equal(b.User.ID, user.ID)

		userAttr := UserAttribute{UserName: "mark", NickName: "Mark Z."}
		err = tx.Create(&userAttr)
		r.NoError(err)

		uA := UserAttribute{}
		err = tx.Eager().Find(&uA, userAttr.ID)
		r.NoError(err)
		r.Equal(uA.User.ID, user.ID)
	})
}

func Test_Find_Eager_Belongs_To_Nulls(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1"}
		err = tx.Create(&book)
		r.NoError(err)

		b := Book{}
		err = tx.Eager().Find(&b, book.ID)
		r.NoError(err)
	})
}

func Test_Find_Eager_Belongs_To_Pointers(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		body := Body{}
		err := tx.Create(&body)
		r.NoError(err)

		head := Head{BodyID: body.ID}
		err = tx.Create(&head)
		r.NoError(err)

		b := Body{}
		err = tx.Eager().Find(&b, body.ID)
		r.NoError(err)
	})
}

func Test_Find_Eager_Has_One(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		coolSong := Song{Title: "Hook - Blues Traveler", UserID: user.ID}
		err = tx.Create(&coolSong)
		r.NoError(err)

		u := User{}
		err = tx.Eager().Find(&u, user.ID)
		r.NoError(err)

		r.NotEqual(u.ID, 0)
		r.Equal(u.Name.String, "Mark")
		r.Equal(u.FavoriteSong.ID, coolSong.ID)

		// eager should work with rawquery
		uid := u.ID
		u = User{}
		err = tx.RawQuery("select * from users where id=?", uid).First(&u)
		r.NoError(err)
		r.Equal(u.FavoriteSong.ID, uuid.Nil)

		err = tx.RawQuery("select * from users where id=?", uid).Eager("FavoriteSong").First(&u)
		r.NoError(err)
		r.Equal(u.FavoriteSong.ID, coolSong.ID)
	})
}

func Test_Find_Eager_Has_One_With_Inner_Associations_Struct(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		composer := Composer{Name: "Blues Traveler"}
		err = tx.Create(&composer)
		r.NoError(err)

		coolSong := Song{Title: "Hook", UserID: user.ID, ComposedByID: composer.ID}
		err = tx.Create(&coolSong)
		r.NoError(err)

		u := User{}
		err = tx.Eager("FavoriteSong.ComposedBy").Find(&u, user.ID)
		r.NoError(err)

		r.NotEqual(u.ID, 0)
		r.Equal(u.Name.String, "Mark")
		r.Equal(u.FavoriteSong.ID, coolSong.ID)
		r.Equal(u.FavoriteSong.ComposedBy.Name, composer.Name)
	})
}

func Test_Find_Eager_Has_One_With_Inner_Associations_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		writer := Writer{Name: "Mark Bates", BookID: book.ID}
		err = tx.Create(&writer)
		r.NoError(err)

		u := User{}
		err = tx.Eager("Books.Writers").Find(&u, user.ID)
		r.NoError(err)

		r.NotEqual(u.ID, 0)
		r.Equal(u.Name.String, "Mark")
		r.Equal(len(u.Books), 1)
		r.Equal(len(u.Books[0].Writers), 1)

		r.Equal(u.Books[0].Title, book.Title)
		r.Equal(u.Books[0].Writers[0].Name, writer.Name)
		r.Zero(u.Books[0].Writers[0].Book.ID)
	})
}

func Test_Eager_Bad_Format(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		u := User{}
		err = tx.Eager("Books.").First(&u)
		r.Error(err)

		err = tx.Eager("Books.*").First(&u)
		r.Error(err)

		err = tx.Eager(".*").First(&u)
		r.Error(err)

		err = tx.Eager(".").First(&u)
		r.Error(err)
	})
}

func Test_Find_Eager_Many_To_Many(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		address := Address{Street: "Pop Avenue", HouseNumber: 1}
		err = tx.Create(&address)
		r.NoError(err)

		ownerProperty := UsersAddress{UserID: user.ID, AddressID: address.ID}
		err = tx.Create(&ownerProperty)
		r.NoError(err)

		u := User{}
		err = tx.Eager("Houses").Find(&u, user.ID)
		r.NoError(err)

		r.NotEqual(u.ID, 0)
		r.Equal(u.Name.String, "Mark")

		r.Equal(len(u.Houses), 1)
		r.Equal(u.Houses[0].Street, address.Street)

		address2 := Address{Street: "Pop Avenue 2", HouseNumber: 1}
		err = tx.Create(&address2)
		r.NoError(err)

		user2 := User{Name: nulls.NewString("Mark 2")}
		err = tx.Create(&user2)
		r.NoError(err)

		ownerProperty2 := UsersAddress{UserID: user2.ID, AddressID: address2.ID}
		err = tx.Create(&ownerProperty2)
		r.NoError(err)

		// eager should work with rawquery
		uid := u.ID
		u = User{}
		err = tx.RawQuery("select * from users where id=?", uid).Eager("Houses").First(&u)
		r.NoError(err)
		r.Equal(1, len(u.Houses))

		// eager ALL
		var users []User
		err = tx.RawQuery("select * from users order by created_at asc").Eager("Houses").All(&users)
		r.NoError(err)
		r.Equal(2, len(users))

		u = users[0]
		r.Equal(u.Name.String, "Mark")
		r.Equal(1, len(u.Houses))
		r.Equal(u.Houses[0].Street, "Pop Avenue")

		u = users[1]
		r.Equal(u.Name.String, "Mark 2")
		r.Equal(1, len(u.Houses))
		r.Equal(u.Houses[0].Street, "Pop Avenue 2")
	})
}

func Test_Load_Associations_Loaded_Model(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		u := User{}
		err = tx.Find(&u, user.ID)

		r.NoError(err)
		r.Zero(len(u.Books))

		err = tx.Load(&u)

		r.NoError(err)
		r.Equal(len(u.Books), 1)
		r.Equal(u.Books[0].Title, book.Title)
	})
}

func Test_First(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		first := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&first)
		r.NoError(err)

		last := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&last)
		r.NoError(err)

		u := User{}
		err = tx.Where("name = 'Mark'").First(&u)
		r.NoError(err)

		r.Equal(first.ID, u.ID)
	})
}

func Test_Last(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		first := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&first)
		r.NoError(err)

		last := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&last)
		r.NoError(err)

		u := User{}
		err = tx.Where("name = 'Mark'").Last(&u)
		r.NoError(err)

		r.Equal(last.ID, u.ID)
	})
}

func Test_All(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			r.NoError(err)
		}

		u := Users{}
		err := tx.All(&u)
		r.NoError(err)
		r.Equal(len(u), 3)

		u = Users{}
		err = tx.Where("name = 'Mark'").All(&u)
		r.NoError(err)
		r.Equal(len(u), 1)
	})
}

func Test_All_Eager_Slice_With_All(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			r.NoError(err)

			book := Book{Title: "Book of " + user.Name.String, UserID: nulls.NewInt(user.ID)}
			err = tx.Create(&book)
			r.NoError(err)
		}

		u := Users{}
		err := tx.Eager("Books.User").All(&u)
		r.NoError(err)
		r.Equal(len(u), 3)

		r.Equal(u[0].ID, u[0].Books[0].User.ID)
		r.Equal(u[1].ID, u[1].Books[0].User.ID)
		r.Equal(u[2].ID, u[2].Books[0].User.ID)
	})
}

func Test_All_Eager(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			r.NoError(err)

			if name == "Mark" {
				book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
				err = tx.Create(&book)
				r.NoError(err)
			}
		}

		u := Users{}
		err := tx.Eager(" Books ", " ").Where("name = 'Mark'").All(&u)
		r.NoError(err)
		r.Equal(len(u), 1)
		r.Equal(len(u[0].Books), 1)
	})
}

func Test_All_Eager_For_Query(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		u := Users{}
		q := tx.Q()
		err = q.Eager("Books").Where("name = 'Mark'").All(&u)
		r.NoError(err)
		r.Equal(len(u), 1)
		r.Equal(len(u[0].Books), 1)
	})
}

func Test_All_Eager_Field_Not_Found_Error(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		u := Users{}
		err = tx.Eager("FieldNotFound").Where("name = 'Mark'").All(&u)
		r.Error(err)
		r.Equal("could not retrieve associations: field FieldNotFound does not exist in model User", err.Error())
	})
}

func Test_All_Eager_Allow_Chain_Call(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		coolSong := Song{Title: "Hook - Blues Traveler", UserID: user.ID}
		err = tx.Create(&coolSong)
		r.NoError(err)

		book := Book{Title: "Pop Book", Isbn: "PB1", UserID: nulls.NewInt(user.ID)}
		err = tx.Create(&book)
		r.NoError(err)

		u := Users{}
		err = tx.Eager("Books").Eager("FavoriteSong").Where("name = 'Mark'").All(&u)
		r.NoError(err)
		r.Equal(len(u), 1)
		r.Equal(len(u[0].Books), 1)
		r.Equal(u[0].FavoriteSong.Title, coolSong.Title)
	})
}

func Test_Count(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)
		c, err := tx.Count(&user)
		r.NoError(err)
		r.Equal(c, 1)

		c, err = tx.Where("1=1").CountByField(&user, "distinct id")
		r.NoError(err)
		r.Equal(c, 1)
		// should ignore order in count

		c, err = tx.Order("id desc").Count(&user)
		r.NoError(err)
		r.Equal(c, 1)

		var uAQ []UsersAddressQuery
		_, err = Q(tx).Select("users_addresses.*").LeftJoin("users", "users.id=users_addresses.user_id").Count(&uAQ)
		r.NoError(err)

		_, err = Q(tx).Select("users_addresses.*", "users.name", "users.email").LeftJoin("users", "users.id=users_addresses.user_id").Count(&uAQ)
		r.NoError(err)
	})
}

func Test_Count_Disregards_Pagination(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		names := []string{
			"Jack",
			"Hurley",
			"Charlie",
			"Desmond",
			"Juliet",
			"Locke",
			"Sawyer",
			"Kate",
			"Benjamin Linus",
		}

		for _, name := range names {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			r.NoError(err)
		}

		firstUsers := Users{}
		secondUsers := Users{}

		q := tx.Paginate(1, 3)
		r.NoError(q.All(&firstUsers))
		r.Equal(len(names), q.Paginator.TotalEntriesSize) // ensure paginator populates count
		r.Equal(3, len(firstUsers))

		firstUsers = Users{}
		q = tx.RawQuery("select * from users").Paginate(1, 3)
		r.NoError(q.All(&firstUsers))
		r.Equal(1, q.Paginator.Page)
		r.Equal(3, q.Paginator.PerPage)
		r.Equal(len(names), q.Paginator.TotalEntriesSize) // ensure paginator populates count

		r.Equal(3, len(firstUsers))
		totalFirstPage := q.Paginator.TotalPages

		q = tx.Paginate(2, 3)
		r.NoError(q.All(&secondUsers))

		r.Equal(3, len(secondUsers))
		totalSecondPage := q.Paginator.TotalPages

		r.NotEqual(0, totalFirstPage)
		r.NotEqual(0, totalSecondPage)
		r.Equal(totalFirstPage, totalSecondPage)

		firstUsers = Users{}
		q = tx.RawQuery("select * from users limit  2").Paginate(1, 5)
		err := q.All(&firstUsers)
		r.NoError(err)
		r.Equal(2, len(firstUsers)) // raw query limit applies

		firstUsers = Users{}
		q = tx.RawQuery("select * from users limit 2 offset 1").Paginate(1, 5)
		err = q.All(&firstUsers)
		r.NoError(err)
		r.Equal(2, len(firstUsers))

		firstUsers = Users{}
		q = tx.RawQuery("select * from users limit 2 offset\t1").Paginate(1, 5)
		err = q.All(&firstUsers)
		r.NoError(err)
		r.Equal(2, len(firstUsers))

		firstUsers = Users{}
		q = tx.RawQuery(`select * from users limit 2 offset
			1`).Paginate(1, 5)
		err = q.All(&firstUsers)
		r.NoError(err)
		r.Equal(2, len(firstUsers))

		firstUsers = Users{}
		q = tx.RawQuery(`select * from users limit 2 offset
			1
			`).Paginate(1, 5) // ending space and tab
		err = q.All(&firstUsers)
		r.NoError(err)
		r.Equal(2, len(firstUsers))

		if tx.Dialect.Name() == "sqlite" {
			firstUsers = Users{}
			q = tx.RawQuery("select * from users limit 2,1").Paginate(1, 5)
			err = q.All(&firstUsers)
			r.NoError(err)
			r.Equal(2, len(firstUsers))

			firstUsers = Users{}
			q = tx.RawQuery("select * from users limit 2 , 1").Paginate(1, 5)
			err = q.All(&firstUsers)
			r.NoError(err)
			r.Equal(2, len(firstUsers))
		}

		if tx.Dialect.Name() == "postgresql" {
			firstUsers = Users{}
			q = tx.RawQuery("select * from users FETCH FIRST 3 rows only").Paginate(1, 5)
			err = q.All(&firstUsers)
			r.NoError(err)
			r.Equal(3, len(firstUsers)) // should fetch only 3
		}
	})
}

func Test_Count_RawQuery(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		c, err := tx.RawQuery("select * from users as users").Count(nil)
		r.NoError(err)
		r.Equal(c, 1)

		c, err = tx.RawQuery("select * from users as users where id = -1").Count(nil)
		r.NoError(err)
		r.Equal(c, 0)

		c, err = tx.RawQuery("select name, max(created_at) from users as users group by name").Count(nil)
		r.NoError(err)
		r.Equal(c, 1)

		c, err = tx.RawQuery("select name from users order by name asc limit 5").Count(nil)
		r.NoError(err)
		r.Equal(c, 1)

		c, err = tx.RawQuery("select name from users order by name asc limit 5 offset 0").Count(nil)
		r.NoError(err)
		r.Equal(c, 1)
	})
}

func Test_Exists(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		t, _ := tx.Where("id = ?", 0).Exists("users")
		r.False(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		r.NoError(err)

		t, _ = tx.Where("id = ?", user.ID).Exists("users")
		r.True(t)
	})
}

func Test_FindManyToMany(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		parent := &Parent{}
		r.NoError(tx.Create(parent))

		student := &Student{}
		r.NoError(tx.Create(student))

		r.NoError(tx.RawQuery("INSERT INTO parents_students (student_id, parent_id) VALUES(?,?)", student.ID, parent.ID).Exec())

		p := &Parent{}
		err := tx.Eager("Students").Find(p, parent.ID)
		r.NoError(err)
	})
}
