package pop_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop"
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

	r.Equal(sl.Type.Kind(), reflect.Slice)
	r.Equal(sl.Value.Len(), 0)
}

func Test_Model_Meta_Map(t *testing.T) {
	r := require.New(t)

	u := &User{}
	m := pop.Model{Value: u}
	mm := m.Meta()
	sl := mm.MakeMap()

	r.Equal(sl.Type.Kind(), reflect.Map)
	r.Equal(sl.Value.Type(), reflect.MapOf(reflect.TypeOf(u.ID), mm.IndirectType))
}

func Test_Model_Meta_Associations(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: &User{}}
	mm := m.Meta()

	mAssociations := mm.Associations()
	r.Equal(3, len(mAssociations))
}
