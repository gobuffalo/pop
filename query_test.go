package pop_test

import (
	"fmt"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Where(t *testing.T) {
	a := require.New(t)
	m := &pop.Model{Value: &Enemy{}}

	q := PDB.Where("id = ?", 1)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE id = ?"), q.ToSQLBuilder(m).String())

	q.Where("first_name = ? and last_name = ?", "Mark", "Bates")
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE id = ? AND first_name = ? and last_name = ?"), q.ToSQLBuilder(m).String())

	q = PDB.Where("name = ?", "Mark 'Awesome' Bates")
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE name = ?"), q.ToSQLBuilder(m).String())

	q = PDB.Where("name = ?", "'; truncate users; --")
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE name = ?"), q.ToSQLBuilder(m).String())
}

func Test_Order(t *testing.T) {
	a := require.New(t)

	m := &pop.Model{Value: &Enemy{}}
	q := PDB.Order("id desc")
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies ORDER BY id desc"), q.ToSQLBuilder(m).String())

	q.Order("name desc")
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies ORDER BY id desc, name desc"), q.ToSQLBuilder(m).String())
}

func Test_ToSQL(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := &pop.Model{Value: &User{}}

		s := "SELECT name as full_name, users.alive, users.bio, users.birth_date, users.created_at, users.id, users.name, users.price, users.updated_at FROM users AS users"

		query := pop.Q(tx)
		q, _ := query.ToSQL(user)
		a.Equal(s, q)

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s ORDER BY id desc", s), q)

		query = tx.Where("id = 1")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1", s), q)

		query = tx.Where("id = 1").Where("name = 'Mark'")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark'", s), q)

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark' ORDER BY id desc", s), q)

		query.Order("name asc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark' ORDER BY id desc, name asc", s), q)

		query = tx.Limit(10)
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s LIMIT 10", s), q)

		query = tx.Paginate(3, 10)
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s LIMIT 10 OFFSET 20", s), q)
	})
}

func Test_ToSQLInjection(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := &pop.Model{Value: &User{}}
		query := tx.Where("name = '?'", "\\\u0027 or 1=1 limit 1;\n-- ")
		q, _ := query.ToSQL(user)
		a.NotEqual("SELECT * FROM users AS users WHERE name = '\\'' or 1=1 limit 1;\n-- '", q)
	})
}

func Test_ToSQL_RawQuery(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		query := tx.RawQuery("this is some ? raw ?", "random", "query")
		q, args := query.ToSQL(nil)
		a.Equal(q, tx.Dialect.TranslateSQL("this is some ? raw ?"))
		a.Equal(args, []interface{}{"random", "query"})
	})
}
