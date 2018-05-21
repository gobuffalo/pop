package columns

import (
	"sort"
	"strings"
)

type WriteableColumns struct {
	Columns
}

// UpdateString returns the SQL column list part of the UPDATE
// query.
func (c WriteableColumns) UpdateString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.UpdateString())
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}
