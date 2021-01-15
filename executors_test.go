package pop

import (
	"testing"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func Test_IsZeroOfUnderlyingType(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		car := &ValidatableCar{Name: "VW"}
		r.True(IsZeroOfUnderlyingType(car.ID))
		err := tx.Save(car)
		r.NoError(err)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		r.False(IsZeroOfUnderlyingType(car.ID))

		var i int
		r.True(IsZeroOfUnderlyingType(i))
		i = 32
		r.False(IsZeroOfUnderlyingType(i))

		var s string
		r.True(IsZeroOfUnderlyingType(s))
		s = "42"
		r.False(IsZeroOfUnderlyingType(s))

		var u uuid.UUID
		r.True(IsZeroOfUnderlyingType(u))
		u, err = uuid.NewV1()
		r.NoError(err)
		r.False(IsZeroOfUnderlyingType(u))
	})
}

func Test_ValidateAndSave(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

func Test_Create_Single_Incremental_ID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *Connection) {
		singleID := &SingleID{}
		err := tx.Create(singleID)
		r.NoError(err)
		r.NotZero(singleID.ID)
	})
}

func Test_ValidateAndCreate_With_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		u := &User{Name: nulls.NewString("Mark")}
		r.Zero(u.ID)
		r.NoError(tx.Save(u))
		r.NotZero(u.ID)

		uat := u.UpdatedAt.UnixNano()

		r.NoError(tx.Save(u))
		time.Sleep(1 * time.Second)
		r.NotEqual(uat, u.UpdatedAt.UnixNano())
	})
}

func Test_Save_With_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		u := Users{
			{Name: nulls.NewString("Mark")},
			{Name: nulls.NewString("Larry")},
		}
		r.Zero(u[0].ID)
		r.Zero(u[1].ID)

		r.NoError(tx.Save(&u))
		r.NotZero(u[0].ID)
		r.NotZero(u[1].ID)

		uat := u[0].UpdatedAt.UnixNano()

		r.NoError(tx.Save(u))
		r.NotEqual(uat, u[0].UpdatedAt.UnixNano())
	})
}

func Test_Create(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

func Test_Create_With_Non_ID_PK(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, _ := tx.Count(&CrookedColour{})
		djs := []CrookedColour{
			{Name: "Phil Slabber"},
			{Name: "Leon Debaughn"},
			{Name: "Liam Merrett-Park"},
		}
		err := tx.Create(&djs)
		r.NoError(err)

		ctx, _ := tx.Count(&CrookedColour{})
		r.Equal(count+3, ctx)
		r.NotEqual(djs[0].ID, djs[1].ID)
		r.NotEqual(djs[1].ID, djs[2].ID)
	})
}

func Test_Create_With_Non_ID_PK_String(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		count, _ := tx.Count(&CrookedSong{})
		djs := []CrookedSong{
			{ID: "Flow"},
			{ID: "Do It Like You"},
			{ID: "I C Light"},
		}
		err := tx.Create(&djs)
		r.NoError(err)

		ctx, _ := tx.Count(&CrookedSong{})
		r.Equal(count+3, ctx)
		r.NotEqual(djs[0].ID, djs[1].ID)
		r.NotEqual(djs[1].ID, djs[2].ID)
	})
}

func Test_Create_Non_PK_ID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		r.NoError(tx.Create(&NonStandardID{OutfacingID: "make sure the tested entry does not have pk=0"}))

		count, err := tx.Count(&NonStandardID{})
		entry := &NonStandardID{
			OutfacingID: "beautiful to the outside ID",
		}
		r.NoError(tx.Create(entry))

		ctx, err := tx.Count(&NonStandardID{})
		r.NoError(err)
		r.Equal(count+1, ctx)
		r.NotZero(entry.ID)
	})
}

func Test_Eager_Create_Has_Many(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
		r.Equal(1, len(u.Books))
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(1, len(u.Houses))
		r.Equal(u.Houses[0].Street, "Modelo")
	})
}

func Test_Eager_Create_Has_Many_With_Existing(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		addr := Address{HouseNumber: 42, Street: "Life"}
		addrVerrs, addrErr := tx.ValidateAndCreate(&addr)
		r.NoError(addrErr)
		addrCount, _ := tx.Count(&Address{})
		r.Zero(addrVerrs.Count())
		r.Equal(1, addrCount)
		r.NotZero(addr.ID)

		count, _ := tx.Count(&User{})
		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
				addr,
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
		r.Equal(addrCount+1, ctx)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(1, len(u.Books))
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(2, len(u.Houses))
		if u.Houses[0].ID == addr.ID {
			r.Equal(u.Houses[0].Street, "Life")
			r.Equal(u.Houses[1].Street, "Modelo")
		} else {
			r.Equal(u.Houses[0].Street, "Modelo")
			r.Equal(u.Houses[1].Street, "Life")
		}
	})
}

func Test_Eager_Create_Has_Many_Reset_Eager_Mode_Connection(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

func Test_Eager_Validate_And_Create_Parental_With_Existing(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		addr := Address{HouseNumber: 42, Street: "Life"}
		addrVerrs, addrErr := tx.ValidateAndCreate(&addr)
		r.NoError(addrErr)
		addrCount, _ := tx.Count(&Address{})
		r.Zero(addrVerrs.Count())
		r.Equal(1, addrCount)
		r.NotZero(addr.ID)

		m2mCount, m2mErr := tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Zero(m2mCount)

		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
				addr,
			},
		}
		count, _ := tx.Count(&User{})

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)
		r.Equal(0, verrs.Count())

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		r.Equal(addrCount+1, ctx)

		m2mCount, m2mErr = tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Equal(2, m2mCount)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(1, len(u.Books))
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(2, len(u.Houses))
		if u.Houses[0].ID == addr.ID {
			r.Equal(u.Houses[0].Street, "Life")
			r.Equal(u.Houses[1].Street, "Modelo")
		} else {
			r.Equal(u.Houses[1].ID, addr.ID)
			r.Equal(u.Houses[0].Street, "Modelo")
			r.Equal(u.Houses[1].Street, "Life")
		}
	})
}

func Test_Eager_Validate_And_Create_Parental_With_Partial_Existing(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		addr := Address{HouseNumber: 42, Street: "Life"}
		addrVerrs, addrErr := tx.ValidateAndCreate(&addr)
		r.NoError(addrErr)
		addrCount, _ := tx.Count(&Address{})
		r.Zero(addrVerrs.Count())
		r.Equal(1, addrCount)
		r.NotZero(addr.ID)

		m2mCount, m2mErr := tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Zero(m2mCount)

		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}},
			FavoriteSong: Song{Title: "Hook - Blues Traveler"},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
				Address{ID: addr.ID},
			},
		}
		count, _ := tx.Count(&User{})

		verrs, err := tx.Eager().ValidateAndCreate(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)
		r.Equal(0, verrs.Count())

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		r.Equal(addrCount+1, ctx)

		m2mCount, m2mErr = tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Equal(2, m2mCount)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(1, len(u.Books))
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(2, len(u.Houses))
		if u.Houses[0].ID == addr.ID {
			r.Equal("Life", u.Houses[0].Street) // Street is blanked out
			r.Equal("Modelo", u.Houses[1].Street)
		} else {
			r.Equal(addr.ID, u.Houses[1].ID)
			r.Equal("Modelo", u.Houses[0].Street)
			r.Equal("Life", u.Houses[1].Street) // Street is blanked out
		}
	})
}

func Test_Flat_Validate_And_Create_Parental_With_Existing(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		addr := Address{HouseNumber: 42, Street: "Life"}
		addrVerrs, addrErr := tx.ValidateAndCreate(&addr)
		r.NoError(addrErr)
		addrCount, _ := tx.Count(&Address{})
		r.Zero(addrVerrs.Count())
		r.Equal(1, addrCount)
		r.NotZero(addr.ID)

		book := Book{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}
		bookVerrs, bookErr := tx.ValidateAndCreate(&book)
		r.NoError(bookErr)
		r.Zero(bookVerrs.Count())
		r.NotZero(book.ID)

		book2 := Book{Title: "Pop Book2", Isbn: "PB2", Description: "Awesome Book Also!"}
		bookVerrs, bookErr = tx.ValidateAndCreate(&book2)
		r.NoError(bookErr)
		r.Zero(bookVerrs.Count())
		r.NotZero(book2.ID)

		bookCount, _ := tx.Count(&Book{})
		r.Equal(2, bookCount)

		song := Song{Title: "Hook - Blues Traveler"}
		songVerrs, songErr := tx.ValidateAndCreate(&song)
		r.NoError(songErr)
		songCount, _ := tx.Count(&Song{})
		r.Zero(songVerrs.Count())
		r.Equal(1, songCount)
		r.NotZero(song.ID)

		m2mCount, m2mErr := tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Zero(m2mCount)

		user := User{
			Name:         nulls.NewString("Mark 'Awesome' Bates"),
			Books:        Books{book, book2},
			FavoriteSong: song,
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
				addr,
			},
		}
		count, _ := tx.Count(&User{})

		verrs, err := tx.ValidateAndCreate(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)
		r.Equal(0, verrs.Count())

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		r.Equal(addrCount, ctx)

		ctx, _ = tx.Count(&Book{})
		r.Equal(bookCount, ctx)

		ctx, _ = tx.Count(&Song{})
		r.Equal(songCount, ctx)

		m2mCount, m2mErr = tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Equal(1, m2mCount)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(2, len(u.Books))
		if u.Books[0].ID == book.ID {
			r.Equal(u.Books[0].Title, "Pop Book")
			r.Equal(u.Books[1].Title, "Pop Book2")
		} else {
			r.Equal(u.Books[1].Title, "Pop Book")
			r.Equal(u.Books[0].Title, "Pop Book2")
		}
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(1, len(u.Houses))
		r.Equal(addr.ID, u.Houses[0].ID)
		r.Equal("Life", u.Houses[0].Street)
	})
}

func Test_Flat_Validate_And_Create_Parental_With_Partial_Existing(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	r := require.New(t)
	transaction(func(tx *Connection) {
		addr := Address{HouseNumber: 42, Street: "Life"}
		addrVerrs, addrErr := tx.ValidateAndCreate(&addr)
		r.NoError(addrErr)
		addrCount, _ := tx.Count(&Address{})
		r.Zero(addrVerrs.Count())
		r.Equal(1, addrCount)
		r.NotZero(addr.ID)

		book := Book{Title: "Pop Book", Isbn: "PB1", Description: "Awesome Book!"}
		bookVerrs, bookErr := tx.ValidateAndCreate(&book)
		r.NoError(bookErr)
		bookCount, _ := tx.Count(&Book{})
		r.Zero(bookVerrs.Count())
		r.Equal(1, bookCount)
		r.NotZero(book.ID)

		song := Song{Title: "Hook - Blues Traveler"}
		songVerrs, songErr := tx.ValidateAndCreate(&song)
		r.NoError(songErr)
		songCount, _ := tx.Count(&Song{})
		r.Zero(songVerrs.Count())
		r.Equal(1, songCount)
		r.NotZero(song.ID)

		m2mCount, m2mErr := tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Zero(m2mCount)

		user := User{
			Name: nulls.NewString("Mark 'Awesome' Bates"),
			//TODO: add another existing here and test for it to make sure this works with multiples (books)
			Books:        Books{Book{ID: book.ID}},
			FavoriteSong: Song{ID: song.ID},
			Houses: Addresses{
				Address{HouseNumber: 86, Street: "Modelo"},
				Address{ID: addr.ID},
			},
		}
		count, _ := tx.Count(&User{})

		verrs, err := tx.ValidateAndCreate(&user)
		r.NoError(err)
		r.NotEqual(user.ID, 0)
		r.Equal(0, verrs.Count())

		ctx, _ := tx.Count(&User{})
		r.Equal(count+1, ctx)

		ctx, _ = tx.Count(&Address{})
		r.Equal(addrCount, ctx)

		ctx, _ = tx.Where("user_id = ?", user.ID).Count(&Book{})
		r.Equal(bookCount, ctx)

		ctx, _ = tx.Count(&Song{})
		r.Equal(songCount, ctx)

		m2mCount, m2mErr = tx.Count(&UsersAddress{})
		r.NoError(m2mErr)
		r.Equal(1, m2mCount)

		u := User{}
		q := tx.Eager().Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		r.NoError(err)
		r.Equal(u.Name.String, "Mark 'Awesome' Bates")
		r.Equal(1, len(u.Books))
		r.Equal(u.Books[0].Title, "Pop Book")
		r.Equal(u.FavoriteSong.Title, "Hook - Blues Traveler")
		r.Equal(1, len(u.Houses))
		r.Equal(addr.ID, u.Houses[0].ID)
		r.Equal("Life", u.Houses[0].Street)
	})
}

func Test_Eager_Create_Belongs_To(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
			Driver: &User{
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

func Test_Eager_Create_Belongs_To_Pointers(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		// Create a body with a head
		body := Body{
			Head: &Head{},
		}

		err := tx.Eager().Create(&body)
		r.NoError(err)
		r.NotZero(body.ID)
		r.NotZero(body.Head.ID)

		ctx, _ := tx.Count(&Body{})
		r.Equal(1, ctx)

		ctx, _ = tx.Count(&Head{})
		r.Equal(1, ctx)

		// Create a body without a head:
		body = Body{
			Head: nil,
		}

		err = tx.Eager().Create(&body)
		r.NoError(err)
		r.NotZero(body.ID)
		r.Nil(body.Head)

		ctx, _ = tx.Count(&Body{})
		r.Equal(2, ctx)

		ctx, _ = tx.Count(&Head{})
		r.Equal(1, ctx)

		err = tx.Eager().Create(&Head{
			BodyID: body.ID,
			Body:   nil,
		})
		r.NoError(err)
	})
}

func Test_Create_Belongs_To_Pointers(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		// Create a body without a head:
		body := Body{
			Head: nil,
		}

		err := tx.Create(&body)
		r.NoError(err)
		r.NotZero(body.ID)
		r.Nil(body.Head)

		// Create a head with the associated model set but not the ID
		created := HeadPtr{
			Body: &body,
		}
		err = tx.Create(&created)
		r.NoError(err)

		found := HeadPtr{}
		err = tx.Find(&found, created.ID)
		r.NoError(err)
		r.Equal(body.ID, *found.BodyID)
	})
}

func Test_Flat_Create_Belongs_To(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)
		user := User{
			Name: nulls.NewString("Larry"),
		}

		err := tx.Create(&user)
		r.NoError(err)
		ctx, _ := tx.Count(&User{})
		r.Equal(1, ctx)

		book := Book{
			Title:       "Pop Book",
			Description: "Pop Book",
			Isbn:        "PB1",
			User:        user,
		}

		err = tx.Create(&book)
		r.NoError(err)

		ctx, _ = tx.Count(&Book{})
		r.Equal(1, ctx)

		err = tx.Eager().Find(&book, book.ID)
		r.NoError(err)

		r.Equal(nulls.NewString("Larry"), book.User.Name)

		car := Taxi{
			Model:  "Fancy car",
			Driver: &user,
		}

		err = tx.Create(&car)
		r.NoError(err)

		ctx, _ = tx.Count(&Taxi{})
		r.Equal(1, ctx)

		err = tx.Eager().Find(&car, car.ID)
		r.NoError(err)

		r.Equal(nulls.NewString("Larry"), car.Driver.Name)
	})
}

func Test_Eager_Creation_Without_Associations(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		user.Name.String = "Marky"
		err := tx.Update(&user)
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.Equal(user.Name.String, "Marky")
	})
}

func Test_UpdateColumns(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		user.Name.String = "Fulano"
		user.UserName = "Fulano"
		err := tx.UpdateColumns(&user, "user_name") // Update UserName field/column only
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.Equal(user.Name.String, "Mark") // Name column should not be updated
		r.Equal(user.UserName, "Fulano")
	})
}

func Test_UpdateColumns_UpdatedAt(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Foo")}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)
		updatedAtBefore := user.UpdatedAt

		user.Name.String = "Bar"
		err := tx.UpdateColumns(&user, "name", "updated_at") // Update name and updated_at
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.NotEqual(user.UpdatedAt, updatedAtBefore) // UpdatedAt should be updated automatically
	})
}

func Test_UpdateColumns_MultipleColumns(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark"), UserName: "Sagan", Email: "test@example.com"}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		user.Name.String = "Ping"
		user.UserName = "Pong"
		user.Email = "fulano@example"
		err := tx.UpdateColumns(&user, "name", "user_name") // Update multiple columns
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.Equal(user.Name.String, "Ping")
		r.Equal(user.UserName, "Pong")
		r.Equal(user.Email, "test@example.com") // Email should not be updated
	})
}

func Test_UpdateColumns_All(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := User{Name: nulls.NewString("Mark"), UserName: "Sagan"}
		tx.Create(&user)

		r.NotZero(user.CreatedAt)
		r.NotZero(user.UpdatedAt)

		user.Name.String = "Ping"
		user.UserName = "Pong"
		user.Email = "ping@pong.com"
		err := tx.UpdateColumns(&user) // Update all columns
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.Equal(user.Name.String, "Ping")
		r.Equal(user.UserName, "Pong")
		r.Equal(user.Email, "ping@pong.com")
	})
}

func Test_UpdateColumns_With_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		user := Users{
			{
				Name:     nulls.NewString("Mark"),
				UserName: "Ping",
			},
			{
				Name:     nulls.NewString("Larry"),
				UserName: "Pong",
			},
		}
		tx.Create(&user)

		r.NotZero(user[0].CreatedAt)
		r.NotZero(user[0].UpdatedAt)

		r.NotZero(user[1].CreatedAt)
		r.NotZero(user[1].UpdatedAt)

		user[0].Name.String = "Fulano"
		user[0].UserName = "Thor"
		user[1].Name.String = "Fulana"
		user[1].UserName = "Freya"

		err := tx.UpdateColumns(&user, "name") // Update Name field/column only
		r.NoError(err)

		r.NoError(tx.Reload(&user))
		r.Equal(user[0].Name.String, "Fulano")
		r.Equal(user[0].UserName, "Ping") // UserName should not be updated
		r.Equal(user[1].Name.String, "Fulana")
		r.Equal(user[1].UserName, "Pong") // UserName should not be updated
	})
}

func Test_Update_With_Slice(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

		r.NoError(tx.Reload(&user))
		r.Equal(user[0].Name.String, "Marky")
		r.Equal(user[1].Name.String, "Lawrence")
	})
}

func Test_Update_UUID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

func Test_Update_With_Non_ID_PK(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		r.NoError(tx.Create(&CrookedColour{Name: "cc is not the first one"}))

		cc := CrookedColour{
			Name: "You?",
		}
		err := tx.Create(&cc)
		r.NoError(err)
		r.NotZero(cc.ID)
		id := cc.ID

		updatedName := "Me!"
		cc.Name = updatedName
		r.NoError(tx.Update(&cc))
		r.Equal(id, cc.ID)

		r.NoError(tx.Reload(&cc))
		r.Equal(updatedName, cc.Name)
		r.Equal(id, cc.ID)
	})
}

func Test_Update_Non_PK_ID(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
	transaction(func(tx *Connection) {
		r := require.New(t)

		client := &NonStandardID{
			OutfacingID: "my awesome hydra client",
		}
		r.NoError(tx.Create(client))

		updatedID := "your awesome hydra client"
		client.OutfacingID = updatedID
		r.NoError(tx.Update(client))
		r.NoError(tx.Reload(client))
		r.Equal(updatedID, client.OutfacingID)
	})
}

func Test_Destroy(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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
	if PDB == nil {
		t.Skip("skipping integration tests")
	}
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

		err = tx.TruncateAll()
		r.NoError(err)

		ctx, _ = tx.Count("users")
		r.Equal(count, ctx)
	})
}
