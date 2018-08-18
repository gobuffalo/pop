package pop

import (
	"testing"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/require"
)

func Test_ValidateAndSave(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndSave(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateSave"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car = &ValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndSave_With_Slice(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := []ValidatableCar{
			{Name: "VW"},
			{Name: "AU"},
		}
		verrs, err := tx.ValidateAndSave(&car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 4)
		r.Equal([]string{"Validate", "ValidateSave", "Validate", "ValidateSave"}, validationLogs)

		r.NotZero(car[0].ID)
		r.NotZero(car[0].CreatedAt)
		r.NotZero(car[1].ID)
		r.NotZero(car[1].CreatedAt)

		validationLogs = []string{}
		car = []ValidatableCar{
			{Name: ""},
			{Name: "AU"},
		}
		verrs, err = tx.ValidateAndSave(&car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := []NotValidatableCar{
			{Name: ""},
			{Name: "AU"},
		}
		verrs, err = tx.ValidateAndSave(&ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndCreate(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndCreate(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car = &ValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndCreate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndCreate_With_Slice(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := []ValidatableCar{
			{Name: "VW"},
			{Name: "AU"},
		}
		verrs, err := tx.ValidateAndCreate(&car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 4)
		r.Equal([]string{"Validate", "ValidateCreate", "Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car[0].ID)
		r.NotZero(car[0].CreatedAt)
		r.NotZero(car[1].ID)
		r.NotZero(car[1].CreatedAt)

		validationLogs = []string{}
		car = []ValidatableCar{
			{Name: ""},
			{Name: "AU"},
		}
		verrs, err = tx.ValidateAndSave(&car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := []NotValidatableCar{
			{Name: ""},
			{Name: "AU"},
		}
		verrs, err = tx.ValidateAndCreate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndUpdate(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndCreate(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car.Name = ""
		verrs, err = tx.ValidateAndUpdate(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndCreate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)

		validationLogs = []string{}
		ncar.Name = ""
		verrs, err = tx.ValidateAndUpdate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndUpdate_With_Slice(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		car := []ValidatableCar{
			{Name: "VW"},
			{Name: "AU"},
		}
		verrs, err := tx.ValidateAndCreate(&car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 4)
		r.Equal([]string{"Validate", "ValidateCreate", "Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car[0].ID)
		r.NotZero(car[0].CreatedAt)
		r.NotZero(car[1].ID)
		r.NotZero(car[1].CreatedAt)

		validationLogs = []string{}
		car[0].Name = ""
		verrs, err = tx.ValidateAndUpdate(&car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := []NotValidatableCar{
			{Name: ""},
			{Name: "AU"},
		}
		verrs, err = tx.ValidateAndCreate(&ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)

		validationLogs = []string{}
		ncar[1].Name = ""
		verrs, err = tx.ValidateAndUpdate(&ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_Exec(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		tx.Create(&user)

		ctx, _ := tx.Count(user)
		r.Equal(1, ctx)

		q := tx.RawQuery("delete from users where id = ?", user.ID)
		err := q.Exec()
		r.NoError(err)

		ctx, _ = tx.Count(user)
		r.Equal(0, ctx)
	})
}

func Test_ExecCount(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		tx.Create(&user)

		ctx, _ := tx.Count(user)
		r.Equal(1, ctx)

		q := tx.RawQuery("delete from users where id = ?", user.ID)
		count, err := q.ExecWithCount()
		r.NoError(err)

		r.Equal(1, count)

		ctx, _ = tx.Count(user)
		r.Equal(0, ctx)
	})
}

func Test_Save(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *Connection) {
		u := &User{Name: nulls.NewString("Mark")}
		r.Zero(u.ID)
		tx.Save(u)
		r.NotZero(u.ID)

		uat := u.UpdatedAt.UnixNano()

		tx.Save(u)
		r.NotEqual(uat, u.UpdatedAt.UnixNano())
	})
}

func Test_Save_With_Slice(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *Connection) {
		u := Users{
			{Name: nulls.NewString("Mark")},
			{Name: nulls.NewString("Larry")},
		}
		r.Zero(u[0].ID)
		r.Zero(u[1].ID)

		tx.Save(&u)
		r.NotZero(u[0].ID)
		r.NotZero(u[1].ID)

		uat := u[0].UpdatedAt.UnixNano()

		tx.Save(u)
		r.NotEqual(uat, u[0].UpdatedAt.UnixNano())
	})
}

func Test_Create(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, _ := tx.Count(&User{})
		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		err := tx.Create(&user)
		r.NoError(err)
		r.NotEqual(0, user.ID)

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		u := User{}
		q := tx.Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal("Mark 'Awesome' Bates", user.Name.String)
	})
}

func Test_Create_stringID(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, err := tx.Count(&Label{})
		r.NoError(err)
		label := Label{ID: "red"}
		err = tx.Create(&label)
		r.NoError(err)
		r.Equal("red", label.ID)

		ctx, err := tx.Count(&Label{})
		r.NoError(err)
		r.Equal(count+1, ctx)

		l := Label{}
		err = tx.Find(&l, "red")
		r.NoError(err)
		r.Equal("red", l.ID)
	})
}

func Test_Create_With_Slice(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, _ := tx.Count(&User{})
		users := Users{
			{Name: nulls.NewString("Mark Bates")},
			{Name: nulls.NewString("Larry M. Jordan")},
			{Name: nulls.NewString("Pop")},
		}
		err := tx.Create(&users)
		r.NoError(err)

		ctx, _ := tx.Count(&User{})
		r.Equal(count+3, ctx)
	})
}

func Test_Eager_Create_Has_Many(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		count, _ := tx.Count(&User{})
		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		err := tx.Eager().Create(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Book{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Song{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		r.Equal(count+1, ctx)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(u.Houses[0].Street, "Modelo")
	})
}

func Test_Eager_Create_Has_Many_Reset_Eager_Mode_Connection(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		count, _ := tx.Count(&User{})
		user1 := User{
			Name:  nulls.NewString("Mark 'Awesome' Bates"),
			Books: Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}},
		}

		err := tx.Eager("Books").Create(&user1)
		r.NoError(err)
		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)
		ctx, _ = tx.Count(&Book{})
		r.Equal(count+1, ctx)

		book := Book{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}

		err = tx.Eager().Create(&book)
		r.NoError(err)
		ctx, _ = tx.Count(&Book{})
		r.Equal(count+2, ctx)
	})
}

func Test_Eager_Validate_And_Create_Has_Many(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *Connection) {
		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{{Title: "Pop Book", Isbn: "PB1"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		r.NoError(err)
		ctx, _ := tx.Count(&User{})
		r.Zero(ctx)
		r.Equal(1, verrs.Count()) // Missing Books.Description.
	})
}

func Test_Eager_Validate_And_Create_Parental(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *Connection) {
		user := User{
			Name:         nulls.NewString(""),
			Books:        Books{{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
			},
		}

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		r.NoError(err)
		ctx, _ := tx.Count(&User{})
		r.Zero(ctx)
		r.Equal(1, verrs.Count()) // Missing Books.Description.
	})
}

func Test_Eager_Create_Belongs_To(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		book := Book{
			Title:       "Pop Book",
			Description: "Pop Book",
			Isbn:        "PB1",
			User: User{
				Name: nulls.NewString("Larry"),
			},
		}

		err := tx.Eager().Create(&book)
		r.NoError(err)

		ctx, _ := tx.Count(&Book{})
		r.Equal(1, ctx)

		ctx, _ = tx.Count(&User{})
		r.Equal(1, ctx)

		car := Taxi{
			Model: "Fancy car",
			Driver: User{
				Name: nulls.NewString("Larry 2"),
			},
		}

		err = tx.Eager().Create(&car)
		r.NoError(err)

		ctx, _ = tx.Count(&Taxi{})
		r.Equal(1, ctx)

		err = tx.Eager().Find(&car, car.ID)
		r.NoError(err)

		r.Equal(nulls.NewString("Larry 2"), car.Driver.Name)
	})
}

func Test_Eager_Creation_Without_Associations(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		code := CourseCode{
			Course: Course{},
		}

		err := tx.Eager().Create(&code)
		r.NoError(err)

		ctx, _ := tx.Count(&CourseCode{})
		r.Equal(1, ctx)
	})
}

func Test_Create_UUID(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, _ := tx.Count(&Song{})
		song := Song{Title: "Automatic Buffalo"}
		err := tx.Create(&song)
		r.NoError(err)
		r.NotZero(song.ID)

		ctx, _ := tx.Count(&Song{})
		r.Equal(count+1, ctx)

		u := Song{}
		q := tx.Where("title = ?", "Automatic Buffalo")
		err = q.First(&u)
		r.NoError(err)
	})
}

func Test_Create_Existing_UUID(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)
		id, err := uuid.NewV4()
		r.NoError(err)

		count, _ := tx.Count(&Song{})
		song := Song{
			ID:    id,
			Title: "Automatic Buffalo",
		}

		err = tx.Create(&song)
		r.NoError(err)
		r.NotZero(song.ID)
		r.Equal(id.String(), song.ID.String())

		ctx, _ := tx.Count(&Song{})
		r.Equal(count+1, ctx)

	})
}

func Test_Create_Timestamps(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		r.Zero(user.CreatedAt)
		r.Zero(user.UpdatedAt)

		err := tx.Create(&user)
		r.NoError(err)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		friend := Friend{FirstName: "Ross", LastName: "Gellar"}
		err = tx.Create(&friend)
		r.NoError(err)
	})
}

func Test_Update(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		user.Name.String = "Marky"
		err := tx.Update(&user)
		r.NoError(err)

		tx.Reload(&user)
		r.Equal(user.Name.String, "Marky")
	})
}

func Test_Update_With_Slice(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := Users{
			{Name: nulls.NewString("Mark")},
			{Name: nulls.NewString("Larry")},
		}
		tx.Create(&user)

		r.NotZero(user[0].CreatedAt)
		r.NotZero(user[0].UpdatedAt)

		r.NotZero(user[1].CreatedAt)
		r.NotZero(user[1].UpdatedAt)

		user[0].Name.String = "Marky"
		user[1].Name.String = "Lawrence"

		err := tx.Update(&user)
		r.NoError(err)

		tx.Reload(&user)
		r.Equal(user[0].Name.String, "Marky")
		r.Equal(user[1].Name.String, "Lawrence")
	})
}

func Test_Update_UUID(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		song := Song{Title: "Automatic Buffalo"}
		err := tx.Create(&song)
		r.NoError(err)

		r.NotZero(song.CreatedAt)
		r.NotZero(song.UpdatedAt)

		song.Title = "Hum"
		err = tx.Update(&song)
		r.NoError(err)

		err = tx.Reload(&song)
		r.NoError(err)
		r.Equal("Hum", song.Title)
	})
}

func Test_Destroy(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, err := tx.Count("users")
		r.NoError(err)
		user := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)

		ctx, err := tx.Count("users")
		r.NoError(err)
		r.Equal(count+1, ctx)

		err = tx.Destroy(&user)
		r.NoError(err)

		ctx, _ = tx.Count("users")
		r.Equal(count, ctx)
	})
}

func Test_Destroy_With_Slice(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, err := tx.Count("users")
		r.NoError(err)
		user := Users{
			{Name: nulls.NewString("Mark")},
			{Name: nulls.NewString("Larry")},
		}
		err = tx.Create(&user)
		r.NoError(err)
		r.NotEqual(user[0].ID, 0)
		r.NotEqual(user[1].ID, 0)

		ctx, err := tx.Count("users")
		r.NoError(err)
		r.Equal(count+2, ctx)

		err = tx.Destroy(&user)
		r.NoError(err)

		ctx, _ = tx.Count("users")
		r.Equal(count, ctx)
	})
}

func Test_Destroy_UUID(t *testing.T) {
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, err := tx.Count("songs")
		r.NoError(err)
		song := Song{Title: "Automatic Buffalo"}
		err = tx.Create(&song)
		r.NoError(err)
		r.NotZero(song.ID)

		ctx, err := tx.Count("songs")
		r.NoError(err)
		r.Equal(count+1, ctx)

		err = tx.Destroy(&song)
		r.NoError(err)

		ctx, _ = tx.Count("songs")
		r.Equal(count, ctx)
	})
}

func Test_TruncateAll(t *testing.T) {
	count := int(0)
	transaction(func(tx *Connection) {
		r := require.New(t)

		var err error
		count, err = tx.Count("users")
		r.NoError(err)
		user := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)

		ctx, err := tx.Count("users")
		r.NoError(err)
		r.Equal(count+1, ctx)
	})

	transaction(func(tx *Connection) {
		r := require.New(t)

		err := tx.TruncateAll()
		r.NoError(err)

		ctx, _ := tx.Count("users")
		r.Equal(count, ctx)
	})
}
