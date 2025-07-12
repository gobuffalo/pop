package models

import (
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

// Pet is a model for testing relationships
type Pet struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (p Pet) TableName() string {
	return "pets"
}

// String implements the Stringer interface
func (p Pet) String() string {
	return p.Name
}

// User is the test model used across various tests
type User struct {
	ID        uuid.UUID    `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Contact   nulls.String `db:"contact"`
	Provider  string       `db:"provider"`
	Password  string       `db:"password"`
	Status    string       `db:"status"`
	Pets      []Pet        `has_many:"pets"`
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
