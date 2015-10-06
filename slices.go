package pop

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// For reading in arrays from postgres
type FloatSlice []float64
type IntSlice []int
type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []byte"))
	}
	(*s) = strToStringSlice(string(b))
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	return fmt.Sprintf("{%s}", strings.Join(s, ",")), nil
}

func strToStringSlice(s string) []string {
	r := strings.Trim(s, "{}")
	return strings.Split(r, ",")
}

func (s *IntSlice) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []byte"))
	}
	str := string(b)
	(*s) = strToIntSlice(str)
	return nil
}

func (s IntSlice) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.Itoa(i)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func strToIntSlice(s string) []int {
	r := strings.Trim(s, "{}")
	a := make([]int, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.Atoi(t)
		a = append(a, i)
	}
	return a
}

func (s *FloatSlice) Scan(src interface{}) error {
	fmt.Printf("src: %s\n", src)
	b, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []byte"))
	}
	str := string(b)
	(*s) = strToFloatSlice(str, *s)
	return nil
}

func (s FloatSlice) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.FormatFloat(i, 'f', -1, 64)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func strToFloatSlice(s string, a []float64) []float64 {
	r := strings.Trim(s, "{}")
	if a == nil {
		a = make([]float64, 0, 10)
	}
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.ParseFloat(t, 64)
		a = append(a, i)
	}
	return a
}
