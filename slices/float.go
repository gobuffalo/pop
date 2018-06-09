package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Float is a slice of float64.
type Float []float64

// Interface implements the nulls.nullable interface.
func (f Float) Interface() interface{} {
	return []float64(f)
}

// Scan implements the sql.Scanner interface.
// It allows to read the float slice from the database value.
func (f *Float) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*f) = strToFloat(str)
	return nil
}

// Value implements the driver.Valuer interface.
// It allows to convert the float slice to a driver.value.
func (f Float) Value() (driver.Value, error) {
	sa := make([]string, len(f))
	for x, i := range f {
		sa[x] = strconv.FormatFloat(i, 'f', -1, 64)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

// UnmarshalText will unmarshall text value into
// the float slice representation of this value.
func (f *Float) UnmarshalText(text []byte) error {
	ss := []float64{}
	for _, x := range strings.Split(string(text), ",") {
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return errors.WithStack(err)
		}
		ss = append(ss, f)
	}
	(*f) = ss
	return nil
}

func strToFloat(s string) []float64 {
	r := strings.Trim(s, "{}")
	a := make([]float64, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.ParseFloat(t, 64)
		a = append(a, i)
	}
	return a
}
