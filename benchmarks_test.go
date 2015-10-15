package pop_test

import (
	"testing"

	"github.com/markbates/going/nulls"
	"github.com/markbates/pop"
)

func Benchmark_Create_Pop(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		for n := 0; n < b.N; n++ {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			tx.Create(u)
		}
	})
}

func Benchmark_Create_Raw(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		for n := 0; n < b.N; n++ {
			u := &User{
				Name: nulls.NewString("Mark Bates"),
			}
			q := "INSERT INTO users (alive, bio, birth_date, created_at, name, price, updated_at) VALUES (:alive, :bio, :birth_date, :created_at, :name, :price, :updated_at)"
			tx.Store.NamedExec(q, u)
		}
	})
}

func Benchmark_Update(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Update(u)
		}
	})
}

func Benchmark_Find_Pop(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Find(u, u.ID)
		}
	})
}

func Benchmark_Find_Raw(b *testing.B) {
	transaction(func(tx *pop.Connection) {
		u := &User{
			Name: nulls.NewString("Mark Bates"),
		}
		tx.Create(u)
		for n := 0; n < b.N; n++ {
			tx.Store.Get(u, "select * from users where id = ?", u.ID)
		}
	})
}
