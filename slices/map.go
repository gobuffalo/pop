package slices

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Map is a map[string]interface.
type Map map[string]interface{}

// Interface implements the nulls.nullable interface.
func (m Map) Interface() interface{} {
	return map[string]interface{}(m)
}

// Scan implements the sql.Scanner interface.
// It allows to read the map from the database value.
func (m *Map) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	err := json.Unmarshal(b, m)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Value implements the driver.Valuer interface.
// It allows to convert the map to a driver.value.
func (m Map) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(b), nil
}

// UnmarshalJSON will unmarshall JSON value into
// the map representation of this value.
func (m Map) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	for key, value := range stuff {
		m[key] = value
	}
	return nil
}

// UnmarshalText will unmarshall text value into
// the map representation of this value.
func (m Map) UnmarshalText(text []byte) error {
	fmt.Println(string(text))
	err := json.Unmarshal(text, &m)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
