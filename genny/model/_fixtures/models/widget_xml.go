package models

import (
	"encoding/xml"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Widget struct {
	ID          uuid.UUID    `xml:"id" db:"id"`
	CreatedAt   time.Time    `xml:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `xml:"updated_at" db:"updated_at"`
	Name        string       `xml:"name" db:"name"`
	Description string       `xml:"description" db:"description"`
	Age         int          `xml:"age" db:"age"`
	Bar         nulls.String `xml:"bar" db:"bar"`
}

// String is not required by pop and may be deleted
func (w Widget) String() string {
	xw, _ := xml.Marshal(w)
	return string(xw)
}

// Widgets is not required by pop and may be deleted
type Widgets []Widget

// String is not required by pop and may be deleted
func (w Widgets) String() string {
	xw, _ := xml.Marshal(w)
	return string(xw)
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
