package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// For reading in arrays from postgres
type Float []float64
type Int []int
type String []string

func (s *String) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	(*s) = strToString(string(b))
	return nil
}

func (s String) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(s, ",")), nil
}

func strToString(s string) []string {
	r := strings.Trim(s, "{}")
	return strings.Split(r, ",")
}

func (s *Int) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*s) = strToInt(str)
	return nil
}

func (s Int) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.Itoa(i)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func strToInt(s string) []int {
	r := strings.Trim(s, "{}")
	a := make([]int, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.Atoi(t)
		a = append(a, i)
	}
	return a
}

func (s *Float) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*s) = strToFloat(str, *s)
	return nil
}

func (s Float) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.FormatFloat(i, 'f', -1, 64)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func strToFloat(s string, a []float64) []float64 {
	r := strings.Trim(s, "{}")
	a = make([]float64, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.ParseFloat(t, 64)
		a = append(a, i)
	}
	return a
}
