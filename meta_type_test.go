package pop_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Model_Meta(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: &User{}}
	mm := m.Meta()

	r.Equal(mm.Type, reflect.TypeOf(m.Value))
	r.Equal(mm.IndirectType, reflect.Indirect(reflect.ValueOf(m.Value)).Type())
}

func Test_Model_Meta_Slice(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: &User{}}
	mm := m.Meta()
	sl := mm.MakeSlice()

	r.Equal(sl.IndirectType.Kind(), reflect.Slice)
	r.Equal(sl.IndirectValue.Len(), 0)
}

func Test_Model_Meta_Map_For_Struct(t *testing.T) {
	r := require.New(t)

	u := User{}
	m := pop.Model{Value: &u}
	mm := m.Meta()
	sl := mm.MakeMap()

	r.Equal(sl.Type.Kind(), reflect.Map)
	r.Equal(sl.Value.Type(), reflect.MapOf(reflect.TypeOf(u.ID), mm.Type))
}

func Test_Model_Meta_Map_For_Slice(t *testing.T) {
	r := require.New(t)

	u := []*User{
		&User{Email: "user@email.com"},
	}
	m := pop.Model{Value: &u}
	mm := m.Meta()
	sl := mm.MakeMap()

	r.Equal(sl.Type, reflect.MapOf(reflect.TypeOf(0), reflect.TypeOf(&User{})))
	r.Equal(1, len(sl.Value.Interface().(map[int]*User)))

	// Map for non-struct with pointer.
	n := 1
	v := []*int{&n}
	m = pop.Model{Value: &v}
	mm = m.Meta()
	sl = mm.MakeMap()
	r.Equal(sl.Type, reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(v[0])))
	r.Equal(1, len(sl.Value.Interface().(map[string]*int)))

	// Map for non-struct without pointer.
	v2 := []int{1}
	m = pop.Model{Value: &v2}
	mm = m.Meta()
	sl = mm.MakeMap()
	r.Equal(sl.Type, reflect.MapOf(reflect.TypeOf(""), reflect.PtrTo(reflect.TypeOf(v2[0]))))
	r.Equal(1, len(sl.Value.Interface().(map[string]*int)))
}

func Test_Model_Meta_Associations(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: &User{}}
	mm := m.Meta()

	mAssociations := mm.Associations()
	r.Equal(3, len(mAssociations))
}

func Test_Model_Meta_Associations_Loading(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			a.NoError(err)

			book := Book{UserID: nulls.NewInt(user.ID)}
			err = tx.Create(&book)
			a.NoError(err)

			if name == "Mark" {
				song := Song{UserID: user.ID}
				err = tx.Create(&song)
				a.NoError(err)

				address := Address{Street: "Pop"}
				err = tx.Create(&address)

				home := UsersAddress{UserID: user.ID, AddressID: address.ID}
				err = tx.Create(&home)
			}
		}

		users := Users{}
		tx.All(&users)

		mt := (&pop.Model{Value: &users}).Meta()
		err := mt.LoadDirect(tx, "has_many")
		a.NoError(err)

		err = mt.LoadDirect(tx, "has_one")
		a.NoError(err)

		// err = pop.LoadBidirect(&users, tx, "many_to_many")
		// a.NoError(err)

		a.Equal(1, len(users[0].Books))
		a.Equal(users[0].ID, users[0].FavoriteSong.UserID)
		a.Zero(users[1].FavoriteSong.UserID)

		books := Books{}
		err = tx.All(&books)
		a.NoError(err)

		mt = (&pop.Model{Value: &books}).Meta()
		mt.LoadIndirect(tx, "belongs_to")
		a.Equal(users[0].ID, books[0].User.ID)
		a.Equal(users[1].ID, books[1].User.ID)
		a.Equal(users[2].ID, books[2].User.ID)
	})
}
