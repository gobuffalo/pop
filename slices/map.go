package slices

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type Map map[string]interface{}

func (s *Map) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	err := json.Unmarshal(b, s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s Map) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(b), nil
}
