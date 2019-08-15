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

	cases := []interface{}{
		tn{},
		[]tn{},
	}
	for _, tc := range cases {
		m := Model{Value: tc}
		r.Equal("this is my table name", m.TableName())
	}
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

	// Override time.Now()
	t0, _ := time.Parse(time.RFC3339, "2019-07-14T00:00:00Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*TimeTimestamp)
	r.Equal(t0, v.CreatedAt)
	r.Equal(t0, v.UpdatedAt)
}

func Test_Touch_Time_Timestamp_With_Existing_Value(t *testing.T) {
	r := require.New(t)

	// Override time.Now()
	t0, _ := time.Parse(time.RFC3339, "2019-07-14T00:00:00Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	createdAt := nowFunc().Add(-36 * time.Hour)

	m := Model{Value: &TimeTimestamp{CreatedAt: createdAt}}
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*TimeTimestamp)
	r.Equal(createdAt, v.CreatedAt)
	r.Equal(t0, v.UpdatedAt)
}

func Test_Touch_Unix_Timestamp(t *testing.T) {
	r := require.New(t)

	m := Model{Value: &UnixTimestamp{}}

	// Override time.Now()
	t0, _ := time.Parse(time.RFC3339, "2019-07-14T00:00:00Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*UnixTimestamp)
	r.Equal(int(t0.Unix()), v.CreatedAt)
	r.Equal(int(t0.Unix()), v.UpdatedAt)
}

func Test_Touch_Unix_Timestamp_With_Existing_Value(t *testing.T) {
	r := require.New(t)

	// Override time.Now()
	t0, _ := time.Parse(time.RFC3339, "2019-07-14T00:00:00Z")
	nowFunc = func() time.Time { return t0 }
	defer func() { nowFunc = time.Now }()

	createdAt := int(time.Now().Add(-36 * time.Hour).Unix())

	m := Model{Value: &UnixTimestamp{CreatedAt: createdAt}}
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*UnixTimestamp)
	r.Equal(createdAt, v.CreatedAt)
	r.Equal(int(t0.Unix()), v.UpdatedAt)
}
