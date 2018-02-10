package pop_test

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_Model_TableName(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: User{}}
	r.Equal(m.TableName(), "users")

	m = pop.Model{Value: &User{}}
	r.Equal(m.TableName(), "users")

	m = pop.Model{Value: &Users{}}
	r.Equal(m.TableName(), "users")

	m = pop.Model{Value: []User{}}
	r.Equal(m.TableName(), "users")

	m = pop.Model{Value: &[]User{}}
	r.Equal(m.TableName(), "users")

	m = pop.Model{Value: []*User{}}
	r.Equal(m.TableName(), "users")

}

type tn struct{}

func (tn) TableName() string {
	return "this is my table name"
}

func Test_TableName(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: tn{}}
	r.Equal("this is my table name", m.TableName())
}

func Test_TableName_With_Array(t *testing.T) {
	r := require.New(t)

	m := pop.Model{Value: []tn{}}
	r.Equal("this is my table name", m.TableName())
}
