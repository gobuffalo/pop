package pop

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gobuffalo/pop/v5/testdata/models/ac"
	"github.com/gobuffalo/pop/v5/testdata/models/bc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/pop/v5/testdata/models/a"
	"github.com/gobuffalo/pop/v5/testdata/models/b"
)

func Test_Model_TableName(t *testing.T) {
	for k, v := range []interface{}{
		User{},
		&User{},

		&Users{},
		Users{},

		[]*User{},
		&[]*User{},

		[]User{},
		&[]User{},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			r := require.New(t)
			m := Model{Value: v}
			r.Equal("users", m.TableName())
		})
	}
}

type tn struct{}

func (tn) TableName() string {
	return "this is my table name"
}

type tnc struct{}

func (tnc) TableName(ctx context.Context) string {
	return ctx.Value("name").(string)
}

// A failing test case for #477
func Test_TableNameCache(t *testing.T) {
	r := assert.New(t)
	r.Equal("usera", (&Model{Value: a.User{}}).TableName())
	r.Equal("userb", (&Model{Value: b.User{}}).TableName())
	r.Equal("usera", (&Model{Value: []a.User{}}).TableName())
	r.Equal("userb", (&Model{Value: []b.User{}}).TableName())
}

// A failing test case for #477
func Test_TableNameContextCache(t *testing.T) {
	ctx := context.WithValue(context.Background(), "name", "context_table")

	r := assert.New(t)
	r.Equal("context_table_useras", (&Model{Value: ac.User{}, ctx: ctx}).TableName())
	r.Equal("context_table_userbs", (&Model{Value: bc.User{}, ctx: ctx}).TableName())
	r.Equal("context_table_useras", (&Model{Value: []ac.User{}, ctx: ctx}).TableName())
	r.Equal("context_table_userbs", (&Model{Value: []bc.User{}, ctx: ctx}).TableName())
}

func Test_TableName(t *testing.T) {
	r := require.New(t)

	cases := []interface{}{
		tn{},
		&tn{},
		[]tn{},
		&[]tn{},
		[]*tn{},
		&[]*tn{},
	}
	for _, tc := range cases {
		m := Model{Value: tc}
		r.Equal("this is my table name", m.TableName())
	}
}

func Test_TableNameContext(t *testing.T) {
	r := require.New(t)

	tn := "context_table_names"
	ctx := context.WithValue(context.Background(), "name", tn)

	cases := []interface{}{
		tnc{},
		[]tnc{},
	}
	for _, tc := range cases {
		m := Model{Value: tc, ctx: ctx}
		r.Equal(tn, m.TableName())
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

	m := NewModel(&TimeTimestamp{}, context.Background())

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

	m := NewModel(&TimeTimestamp{CreatedAt: createdAt}, context.Background())
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*TimeTimestamp)
	r.Equal(createdAt, v.CreatedAt)
	r.Equal(t0, v.UpdatedAt)
}

func Test_Touch_Unix_Timestamp(t *testing.T) {
	r := require.New(t)

	m := NewModel(&UnixTimestamp{}, context.Background())

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

	m := NewModel(&UnixTimestamp{CreatedAt: createdAt}, context.Background())
	m.touchCreatedAt()
	m.touchUpdatedAt()
	v := m.Value.(*UnixTimestamp)
	r.Equal(createdAt, v.CreatedAt)
	r.Equal(int(t0.Unix()), v.UpdatedAt)
}

func Test_IDField(t *testing.T) {
	r := require.New(t)

	type testCustomID struct {
		ID int `db:"custom_id"`
	}
	m := Model{Value: &testCustomID{ID: 1}}
	r.Equal("custom_id", m.IDField())

	type testNormalID struct {
		ID int
	}
	m = Model{Value: &testNormalID{ID: 1}}
	r.Equal("id", m.IDField())
}

type testPrefixID struct {
	ID int `db:"custom_id"`
}

func (t testPrefixID) TableName() string {
	return "foo.bar"
}

func Test_WhereID(t *testing.T) {
	r := require.New(t)
	m := Model{Value: &testPrefixID{ID: 1}}

	r.Equal("foo_bar.custom_id = ?", m.whereID())
	r.Equal("foo_bar.custom_id = :custom_id", m.whereNamedID())

	type testNormalID struct {
		ID int
	}
	m = Model{Value: &testNormalID{ID: 1}}
	r.Equal("id", m.IDField())
}
