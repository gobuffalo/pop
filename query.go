package pop

type Query struct {
	RawSQL                  *clause
	limitResults            int
	whereClauses            clauses
	orderClauses            clauses
	fromClauses             fromClauses
	belongsToThroughClauses belongsToThroughClauses
	Paginator               *Paginator
	Connection              *Connection
}

func (c *Connection) RawQuery(stmt string, args ...interface{}) *Query {
	return Q(c).RawQuery(stmt, args...)
}

func (q *Query) RawQuery(stmt string, args ...interface{}) *Query {
	q.RawSQL = &clause{stmt, args}
	return q
}

func (c *Connection) Where(stmt string, args ...interface{}) *Query {
	return Q(c).Where(stmt, args...)
}

func (q *Query) Where(stmt string, args ...interface{}) *Query {
	q.whereClauses = append(q.whereClauses, clause{stmt, args})
	return q
}

func (c *Connection) Order(stmt string) *Query {
	return Q(c).Order(stmt)
}

func (q *Query) Order(stmt string) *Query {
	q.orderClauses = append(q.orderClauses, clause{stmt, []interface{}{}})
	return q
}

func (c *Connection) Limit(limit int) *Query {
	return Q(c).Limit(limit)
}

func (q *Query) Limit(limit int) *Query {
	q.limitResults = limit
	return q
}

func Q(c *Connection) *Query {
	return &Query{
		RawSQL:     &clause{},
		Connection: c,
	}
}

func (q Query) ToSQL(model *Model, addColumns ...string) (string, []interface{}) {
	sb := NewSQLBuilder(q, model, addColumns...)
	return sb.String(), sb.Args()
}

func (q Query) ToSQLBuilder(model *Model, addColumns ...string) *SQLBuilder {
	return NewSQLBuilder(q, model, addColumns...)
}
