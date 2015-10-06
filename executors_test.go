package pop_test

import (
	"testing"

	"github.com/markbates/going/nulls"
	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Exec(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		tx.Create(&user)

		ctx, _ := tx.Count(user)
		a.Equal(1, ctx)

		q := tx.RawQuery("delete from users where id = ?", user.ID)
		err := q.Exec()
		a.NoError(err)

		ctx, _ = tx.Count(user)
		a.Equal(0, ctx)
	})
}

func Test_Save(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *pop.Connection) {
		u := &User{Name: nulls.NewString("Mark")}
		r.Zero(u.ID)
		tx.Save(u)
		r.NotZero(u.ID)

		uat := u.UpdatedAt.UnixNano()

		tx.Save(u)
		r.NotEqual(uat, u.UpdatedAt.UnixNano())
	})
}

func Test_Create(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		count, _ := tx.Count(&User{})
		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		err := tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, _ := tx.Count(&User{})
		a.Equal(count+1, ctx)

		u := User{}
		q := tx.Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		a.NoError(err)
		a.Equal(user.Name.String, "Mark 'Awesome' Bates")
	})
}

func Test_Create_Timestamps(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		a.Zero(user.CreatedAt)
		a.Zero(user.UpdatedAt)

		err := tx.Create(&user)
		a.NoError(err)

		a.NotZero(user.CreatedAt)
		a.NotZero(user.UpdatedAt)

		friend := Friend{FirstName: "Ross", LastName: "Gellar"}
		err = tx.Create(&friend)
		a.NoError(err)
	})
}

func Test_Update(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		tx.Create(&user)

		a.NotZero(user.CreatedAt)
		a.NotZero(user.UpdatedAt)

		user.Name.String = "Marky"
		err := tx.Update(&user)
		a.NoError(err)

		tx.Reload(&user)
		a.Equal(user.Name.String, "Marky")
	})
}

func Test_Destroy(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		count, err := tx.Count("users")
		user := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, err := tx.Count("users")
		a.Equal(count+1, ctx)

		err = tx.Destroy(&user)
		a.NoError(err)

		ctx, _ = tx.Count("users")
		a.Equal(count, ctx)
	})
}
