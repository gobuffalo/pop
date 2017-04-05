package pop

import (
	"fmt"
	"strings"
)

type joinClause struct {
	JoinType  string
	Table     string
	As        string
	On        []string
	Arguments []interface{}
}

type joinClauses []joinClause

func (c joinClause) String() string {
	sql := fmt.Sprintf("%s %s", c.JoinType, c.Table)

	if c.As != "" && c.As != c.Table {
		sql += fmt.Sprintf(" %s", c.As)
	}

	if len(c.On) > 0 {
		sql += " ON " + strings.Join(c.On, " AND ")
	}

	return sql
}

func (c joinClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, " ")
}
