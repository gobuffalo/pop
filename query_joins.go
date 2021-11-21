package pop

import (
	"github.com/gobuffalo/pop/v6/logging"
)

// Join will append a JOIN clause to the query
func (q *Query) Join(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"JOIN", table, on, args})
	return q
}

// LeftJoin will append a LEFT JOIN clause to the query
func (q *Query) LeftJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"LEFT JOIN", table, on, args})
	return q
}

// RightJoin will append a RIGHT JOIN clause to the query
func (q *Query) RightJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"RIGHT JOIN", table, on, args})
	return q
}

// LeftOuterJoin will append a LEFT OUTER JOIN clause to the query
func (q *Query) LeftOuterJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"LEFT OUTER JOIN", table, on, args})
	return q
}

// RightOuterJoin will append a RIGHT OUTER JOIN clause to the query
func (q *Query) RightOuterJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"RIGHT OUTER JOIN", table, on, args})
	return q
}

// InnerJoin will append an INNER JOIN clause to the query
func (q *Query) InnerJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		log(logging.Warn, "Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"INNER JOIN", table, on, args})
	return q
}
