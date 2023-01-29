package pop

import (
	"time"

	"github.com/gobuffalo/pop/v6/slices"
)

type Cake struct {
	ID        int           `db:"id"`
	Int       slices.Int    `db:"int_slice"`
	Float     slices.Float  `db:"float_slice"`
	String    slices.String `db:"string_slice"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

func (s *PostgreSQLSuite) Test_String() {
	transaction(func(tx *Connection) {
		r := s.Require()

		c := &Cake{
			String: slices.String{"a", "b", "c"},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.NoError(err)
		r.Equal(slices.String{"a", "b", "c"}, c.String)
	})
}

func (s *PostgreSQLSuite) Test_Int() {
	transaction(func(tx *Connection) {
		r := s.Require()

		c := &Cake{
			Int: slices.Int{1, 2, 3},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.NoError(err)
		r.Equal(slices.Int{1, 2, 3}, c.Int)
		r.Equal(slices.Float{}, c.Float)
	})
}

func (s *PostgreSQLSuite) Test_Float() {
	transaction(func(tx *Connection) {
		r := s.Require()

		c := &Cake{
			Float: slices.Float{1.0, 2.1, 3.2},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.NoError(err)
		r.Equal(slices.Int{}, c.Int)
		r.Equal(slices.Float{1.0, 2.1, 3.2}, c.Float)
	})
}
