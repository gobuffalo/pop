package pop

import (
	"strings"
)

// GroupClause holds the field to apply the GROUP clause on
type GroupClause struct {
	Field string
}

type groupClauses []GroupClause

func (c GroupClause) String() string {
	return c.Field
}

func (c groupClauses) String() string {
	var cs []string
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, ", ")
}
