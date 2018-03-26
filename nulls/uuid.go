package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// UUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
type UUID struct {
	UUID  uuid.UUID
	Valid bool
}

func (ns UUID) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.UUID
}

// NewUUID returns a new, properly instantiated
// UUID object.
func NewUUID(u uuid.UUID) UUID {
	return UUID{UUID: u, Valid: true}
}

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return u.UUID.Value()
}

// Scan implements the sql.Scanner interface.
func (u *UUID) Scan(src interface{}) error {
	if src == nil {
		u.UUID, u.Valid = uuid.Nil, false
		return nil
	}

	// Delegate to UUID Scan function
	u.Valid = true
	return u.UUID.Scan(src)
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns UUID) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.UUID.String())
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *UUID) UnmarshalJSON(text []byte) error {
	ns.Valid = false
	ns.UUID = uuid.Nil
	if string(text) == "null" {
		return nil
	}

	s := string(text)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")

	u, err := uuid.FromString(s)
	if err != nil {
		return errors.WithStack(err)
	}
	ns.UUID = u
	ns.Valid = true

	return nil
}

func (ns *UUID) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
