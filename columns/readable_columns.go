package columns

import (
	"sort"
	"strings"
)

type ReadableColumns struct {
	Columns
}

// SelectString returns the SQL column list part of the SELECT
// query.
func (c ReadableColumns) SelectString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.SelectSQL)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}
