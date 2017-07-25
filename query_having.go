package pop

func (q *Query) Having(condition string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		panic("Query is setup to use raw SQL")
	}
	q.havingClauses = append(q.havingClauses, HavingClause{condition, args})

	return q
}
