package pop

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/stretchr/testify/require"
)

func Benchmark_Create_Pop(b *testing.B) {
	transaction(func(tx *Connection) {
		for b.Loop() {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			require.NoError(b, tx.Create(u))
		}
	})
}

func Benchmark_Create_Raw(b *testing.B) {
	transaction(func(tx *Connection) {
		for b.Loop() {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			q := "INSERT INTO users (alive, bio, birth_date, created_at, name, price, updated_at) VALUES (:alive, :bio, :birth_date, :created_at, :name, :price, :updated_at)"
			_, err := tx.Store.NamedExec(q, u)
			require.NoError(b, err)
		}
	})
}

func Benchmark_Update(b *testing.B) {
	transaction(func(tx *Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		require.NoError(b, tx.Create(u))
		for b.Loop() {
			require.NoError(b, tx.Update(u))
		}
	})
}

func Benchmark_Find_Pop(b *testing.B) {
	transaction(func(tx *Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		require.NoError(b, tx.Create(u))
		for b.Loop() {
			require.NoError(b, tx.Find(u, u.ID))
		}
	})
}

func Benchmark_Find_Raw(b *testing.B) {
	transaction(func(tx *Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		require.NoError(b, tx.Create(u))
		for b.Loop() {
			require.NoError(b, tx.Store.Get(u, "select * from users where id = ?", u.ID))
		}
	})
}

func Benchmark_translateOne(b *testing.B) {
	q := "select * from users where id = ? and name = ? and email = ? and a = ? and b = ? and c = ? and d = ? and e = ? and f = ?"
	for b.Loop() {
		translateOne(q)
	}
}

func Benchmark_translateTwo(b *testing.B) {
	q := "select * from users where id = ? and name = ? and email = ? and a = ? and b = ? and c = ? and d = ? and e = ? and f = ?"
	for b.Loop() {
		translateTwo(q)
	}
}

func translateOne(sql string) string {
	curr := 1
	out := make([]byte, 0, len(sql))
	for i := range len(sql) {
		if sql[i] == '?' {
			str := "$" + strconv.Itoa(curr)
			for _, char := range str {
				out = append(out, byte(char))
			}
			curr++
		} else {
			out = append(out, sql[i])
		}
	}
	return string(out)
}

func translateTwo(sql string) string {
	curr := 1
	csql := ""
	for i := range len(sql) {
		x := sql[i]
		if x == '?' {
			csql = fmt.Sprintf("%s$%d", csql, curr)
			curr++
		} else {
			csql += string(x)
		}
	}
	return csql
}
