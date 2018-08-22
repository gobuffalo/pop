package pop

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Model_TableName(t *testing.T) {
	r := require.New(t)

	m := Model{Value: User{}}
	r.Equal(m.TableName(), "users")

	m = Model{Value: &User{}}
	r.Equal(m.TableName(), "users")

	m = Model{Value: &Users{}}
	r.Equal(m.TableName(), "users")

	m = Model{Value: []User{}}
	r.Equal(m.TableName(), "users")

	m = Model{Value: &[]User{}}
	r.Equal(m.TableName(), "users")

	m = Model{Value: []*User{}}
	r.Equal(m.TableName(), "users")

}

type tn struct{}

func (tn) TableName() string {
	return "this is my table name"
}

func Test_TableName(t *testing.T) {
	r := require.New(t)

	m := Model{Value: tn{}}
	r.Equal("this is my table name", m.TableName())
}

func Test_TableName_With_Array(t *testing.T) {
	r := require.New(t)

	m := Model{Value: []tn{}}
	r.Equal("this is my table name", m.TableName())
}

type TimeTimestamp struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UnixTimestamp struct {
	ID        int `db:"id"`
	CreatedAt int `db:"created_at"`
	UpdatedAt int `db:"updated_at"`
}

func Test_Touch_Time_Timestamp(t *testing.T) {
	r := require.New(t)

	m := Model{Value: &TimeTimestamp{}}
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*TimeTimestamp)
	r.NotZero(v.CreatedAt)
	r.NotZero(v.UpdatedAt)
}

func Test_Touch_Unix_Timestamp(t *testing.T) {
	r := require.New(t)

	m := Model{Value: &UnixTimestamp{}}
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*UnixTimestamp)
	r.NotZero(v.CreatedAt)
	r.NotZero(v.UpdatedAt)
}
