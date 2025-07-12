package models

import (
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

// User is the test model used across various tests
type User struct {
	ID        uuid.UUID    `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Contact   nulls.String `db:"contact"`
	Provider  string       `db:"provider"`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (u User) TableName() string {
	return "users"
}

// String implements the Stringer interface
func (u User) String() string {
	return u.Name
}

// Users is a slice of User models
type Users []User
