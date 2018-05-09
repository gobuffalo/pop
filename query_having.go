package pop

import "github.com/gobuffalo/pop/log"

// Having will append a HAVING clause to the query
func (q *Query) Having(condition string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log.DefaultLogger.WithField("raw", q.RawSQL.Fragment).WithField("condition", condition).WithField("args", args).Warn("Query is setup to use raw SQL, not adding HAVING clause")
		return q
	}
	q.havingClauses = append(q.havingClauses, HavingClause{condition, args})

	return q
}
