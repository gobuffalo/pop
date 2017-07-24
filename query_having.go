package pop

func (q *Query) Having(condition string, args ...interface{}) *Query {
	q.havingClauses = append(q.havingClauses, HavingClause{condition, args})

	return q
}
