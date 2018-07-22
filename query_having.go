package pop

// Having will append a HAVING clause to the query
func (q *Query) Having(condition string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		Log("warn", "Query is setup to use raw SQL")
		return q
	}
	q.havingClauses = append(q.havingClauses, HavingClause{condition, args})

	return q
}
