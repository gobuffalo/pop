package pop

import "strings"

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
