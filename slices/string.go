package slices

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"github.com/lib/pq"
)

// For reading in arrays from postgres

// String is a slice of strings.
type String []string

// Interface implements the nulls.nullable interface.
func (s String) Interface() interface{} {
	return []string(s)
}

// Scan implements the sql.Scanner interface.
// It allows to read the string slice from the database value.
func (s *String) Scan(src interface{}) error {
	// Still relying on pq driver to help with string arrays.
	ss := pq.StringArray(*s)
	err := ss.Scan(src)
	*s = String(ss)
	return err
}

// Value implements the driver.Valuer interface.
// It allows to convert the string slice to a driver.value.
func (s String) Value() (driver.Value, error) {
	ss := pq.StringArray(s)
	return ss.Value()
}

// UnmarshalJSON will unmarshall JSON value into
// the string slice representation of this value.
func (s *String) UnmarshalJSON(data []byte) error {
	var ss pq.StringArray
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	*s = String(ss)
	return nil
}

// UnmarshalText will unmarshall text value into
// the string slice representation of this value.
func (s *String) UnmarshalText(text []byte) error {
	r := csv.NewReader(bytes.NewReader(text))

	var words []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		words = append(words, record...)
	}

	*s = String(words)
	return nil
}

// TagValue implements the tagValuer interface, to work with https://github.com/gobuffalo/tags.
func (s String) TagValue() string {
	return s.Format(",")
}

// Format presents the slice as a string, using a given separator.
func (s String) Format(sep string) string {
	return strings.Join([]string(s), sep)
}
