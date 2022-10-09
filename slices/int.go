package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
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
	var str string
	switch t := src.(type) {
	case []byte:
		str = string(t)
	case string:
		str = t
	default:
		return fmt.Errorf("scan source was not []byte nor string but %T", src)
	}

	v, err := strToInt(str)
	*i = v
	return err
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
	var ss []int
	for _, x := range strings.Split(string(text), ",") {
		f, err := strconv.Atoi(x)
		if err != nil {
			return err
		}
		ss = append(ss, f)
	}
	*i = ss
	return nil
}

func strToInt(s string) ([]int, error) {
	r := strings.Trim(s, "{}")
	a := make([]int, 0, 10)

	split := strings.Split(r, ",")
	// Split returns [""] when splitting the empty string.
	if len(split) == 1 && split[0] == "" {
		return a, nil
	}

	for _, t := range split {
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		a = append(a, i)
	}

	return a, nil
}
