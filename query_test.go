package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Where(t *testing.T) {
	a := require.New(t)

	query := PDB.Where("id = ?", 1)
	a.Equal(query.WhereClauses, pop.Clauses{pop.Clause{"id = ?", []interface{}{1}}})

	query.Where("first_name = ? and last_name = ?", "Mark", "Bates")
	a.Equal(query.WhereClauses, pop.Clauses{
		pop.Clause{"id = ?", []interface{}{1}},
		pop.Clause{"first_name = ? and last_name = ?", []interface{}{"Mark", "Bates"}},
	})

	query = PDB.Where("name = ?", "Mark 'Awesome' Bates")
	a.Equal(query.WhereClauses, pop.Clauses{
		pop.Clause{"name = ?", []interface{}{"Mark 'Awesome' Bates"}},
	})

	query = PDB.Where("name = ?", "'; truncate users; --")
	a.Equal(query.WhereClauses, pop.Clauses{
		pop.Clause{"name = ?", []interface{}{"'; truncate users; --"}},
	})
}

func Test_Order(t *testing.T) {
	a := require.New(t)

	query := PDB.Order("id desc")
	a.Equal(query.OrderClauses, pop.Clauses{
		pop.Clause{"id desc", []interface{}{}},
	})

	query.Order("name desc")
	a.Equal(query.OrderClauses, pop.Clauses{
		pop.Clause{"id desc", []interface{}{}},
		pop.Clause{"name desc", []interface{}{}},
	})
}

func Test_ToSQL(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := &pop.Model{Value: &User{}}

		query := pop.Q(tx)
		q, _ := query.ToSQL(user)
		a.Equal("SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users", q)

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users ORDER BY id desc")

		query = tx.Where("id = 1")
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users WHERE id = 1")

		query = tx.Where("id = 1").Where("name = 'Mark'")
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users WHERE id = 1 AND name = 'Mark'")

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users WHERE id = 1 AND name = 'Mark' ORDER BY id desc")

		query.Order("name asc")
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users WHERE id = 1 AND name = 'Mark' ORDER BY id desc, name asc")

		query = tx.Limit(10)
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users LIMIT 10")

		query = tx.Paginate(3, 10)
		q, _ = query.ToSQL(user)
		a.Equal(q, "SELECT alive, bio, birth_date, created_at, id, name, name as full_name, price, updated_at FROM users AS users LIMIT 10 OFFSET 20")
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
