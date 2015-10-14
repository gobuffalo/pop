package pop

import (
	"fmt"
	"strings"
)

type Clause struct {
	Fragment  string
	Arguments []interface{}
}

type Clauses []Clause

func (c Clauses) Join(sep string) string {
	out := make([]string, 0, len(c))
	for _, clause := range c {
		out = append(out, clause.Fragment)
	}
	return strings.Join(out, sep)
}

func (c Clauses) Args() (args []interface{}) {
	for _, clause := range c {
		for _, arg := range clause.Arguments {
			args = append(args, arg)
		}
	}
	return
}

type FromClause struct {
	From string
	As   string
}

type FromClauses []FromClause

func (c FromClause) String() string {
	return fmt.Sprintf("%s AS %s", c.From, c.As)
}

func (c FromClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, ", ")
}
