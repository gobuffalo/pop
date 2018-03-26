package slices

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// For reading in arrays from postgres

type UUID []uuid.UUID

func (s UUID) Interface() interface{} {
	return []uuid.UUID(s)
}

func (s *UUID) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	us, err := strSliceToUUIDSlice(strToUUID(string(b)))
	if err != nil {
		return errors.WithStack(err)
	}
	(*s) = UUID(us)
	return nil
}

func (s UUID) Value() (driver.Value, error) {
	ss := make([]string, len(s))
	for i, u := range s {
		ss[i] = u.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(ss, ",")), nil
}

func (s *UUID) UnmarshalText(text []byte) error {
	ss := []string{}
	for _, x := range strings.Split(string(text), ",") {
		ss = append(ss, strings.TrimSpace(x))
	}
	us, err := strSliceToUUIDSlice(ss)
	if err != nil {
		return errors.WithStack(err)
	}
	(*s) = us
	return nil
}

func (s *UUID) UnmarshalJSON(data []byte) error {
	ss := []string{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	us, err := strSliceToUUIDSlice(ss)
	if err != nil {
		return errors.WithStack(err)
	}
	(*s) = us
	return nil
}

func (s UUID) TagValue() string {
	return s.Format(",")
}

func (s UUID) Format(sep string) string {
	ss := make([]string, len(s))
	for i, u := range s {
		ss[i] = u.String()
	}
	return strings.Join(ss, sep)
}

func strToUUID(s string) []string {
	r := strings.Trim(s, "{}")
	return strings.Split(r, ",")
}

func strSliceToUUIDSlice(ss []string) (UUID, error) {
	us := make([]uuid.UUID, len(ss))
	for i, s := range ss {
		if s == "" {
			continue
		}
		u, err := uuid.FromString(s)
		if err != nil {
			return UUID{}, errors.WithStack(err)
		}
		us[i] = u
	}
	return UUID(us), nil
}
