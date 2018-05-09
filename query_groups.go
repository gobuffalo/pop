package pop

import "github.com/gobuffalo/pop/log"

// GroupBy will append a GROUP BY clause to the query
func (q *Query) GroupBy(field string, fields ...string) *Query {
	if q.RawSQL.Fragment != "" {
		log.DefaultLogger.WithField("raw", q.RawSQL.Fragment).WithField("field", field).WithField("fields", fields).Warn("Query is setup to use raw SQL, not adding GROUP BY clause")
		return q
	}
	q.groupClauses = append(q.groupClauses, GroupClause{field})
	if len(fields) > 0 {
		for i := range fields {
			q.groupClauses = append(q.groupClauses, GroupClause{fields[i]})
		}
	}
	return q
}
