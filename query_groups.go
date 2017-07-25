package pop

func (q *Query) GroupBy(field string, fields ...string) *Query {
	if q.RawSQL.Fragment != "" {
		panic("Query is setup to use raw SQL")
	}
	q.groupClauses = append(q.groupClauses, GroupClause{field})
	if len(fields) > 0 {
		for i := range fields {
			q.groupClauses = append(q.groupClauses, GroupClause{fields[i]})
		}
	}
	return q
}
