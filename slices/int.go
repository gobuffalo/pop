package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Int is a slice of int.
type Int []int

// Interface implements the nulls.nullable interface.
func (i Int) Interface() interface{} {
	return []int(i)
}

// Scan implements the sql.Scanner interface.
// It allows to read the int slice from the database value.
func (i *Int) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*i) = strToInt(str)
	return nil
}

// Value implements the driver.Valuer interface.
// It allows to convert the int slice to a driver.value.
func (i Int) Value() (driver.Value, error) {
	sa := make([]string, len(i))
	for x, v := range i {
		sa[x] = strconv.Itoa(v)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

// UnmarshalText will unmarshall text value into
// the int slice representation of this value.
func (i *Int) UnmarshalText(text []byte) error {
	ss := []int{}
	for _, x := range strings.Split(string(text), ",") {
		f, err := strconv.Atoi(x)
		if err != nil {
			return errors.WithStack(err)
		}
		ss = append(ss, f)
	}
	(*i) = ss
	return nil
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
