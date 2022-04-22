package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/google/jsonapi"
)

// Widget is used by pop to map your widgets database table to your go code.
type Widget struct {
	ID          uuid.UUID    `jsonapi:"primary,id" db:"id"`
	CreatedAt   time.Time    `jsonapi:"attr,created_at" db:"created_at"`
	UpdatedAt   time.Time    `jsonapi:"attr,updated_at" db:"updated_at"`
	Name        string       `jsonapi:"attr,name" db:"name"`
	Description string       `jsonapi:"attr,description" db:"description"`
	Age         int          `jsonapi:"attr,age" db:"age"`
	Bar         nulls.String `jsonapi:"attr,bar" db:"bar"`
}

// String is not required by pop and may be deleted
func (w Widget) String() string {
	var jb strings.Builder
	_ = jsonapi.MarshalPayload(&jb, &w)
	return jb.String()
}

// Widgets is not required by pop and may be deleted
type Widgets []Widget

// String is not required by pop and may be deleted
func (w Widgets) String() string {
	var jb strings.Builder
	_ = jsonapi.MarshalPayload(&jb, &w)
	return jb.String()
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (w *Widget) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: w.Name, Name: "Name"},
		&validators.StringIsPresent{Field: w.Description, Name: "Description"},
		&validators.IntIsPresent{Field: w.Age, Name: "Age"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (w *Widget) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (w *Widget) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
