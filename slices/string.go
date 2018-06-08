package slices

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// For reading in arrays from postgres

// String is a slice of strings.
type String []string

func (s String) Interface() interface{} {
	return []string(s)
}

// Scan implements the sql.Scanner interface.
// It allows to read the string slice from the database value.
func (s *String) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	(*s) = strToString(string(b))
	return nil
}

// Value implements the driver.Valuer interface.
// It allows to convert the string slice to a driver.value.
func (s String) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(s, ",")), nil
}

func (s *String) UnmarshalText(text []byte) error {
	ss := []string{}
	for _, x := range strings.Split(string(text), ",") {
		ss = append(ss, strings.TrimSpace(x))
	}
	(*s) = ss
	return nil
}

func (s *String) UnmarshalJSON(data []byte) error {
	ss := []string{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	(*s) = String(ss)
	return nil
}

func (s String) TagValue() string {
	return s.Format(",")
}

func (s String) Format(sep string) string {
	return strings.Join([]string(s), sep)
}

func strToString(s string) []string {
	r := strings.Trim(s, "{}")
	return strings.Split(r, ",")
}
