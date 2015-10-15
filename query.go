package pop

type Query struct {
	RawSQL       *Clause
	LimitResults int
	WhereClauses Clauses
	OrderClauses Clauses
	FromClauses  FromClauses
	Paginator    *Paginator
	Connection   *Connection
}

func (c *Connection) RawQuery(stmt string, args ...interface{}) *Query {
	return Q(c).RawQuery(stmt, args...)
}

func (q *Query) RawQuery(stmt string, args ...interface{}) *Query {
	q.RawSQL = &Clause{stmt, args}
	return q
}

func (c *Connection) Where(stmt string, args ...interface{}) *Query {
	return Q(c).Where(stmt, args...)
}

func (q *Query) Where(stmt string, args ...interface{}) *Query {
	q.WhereClauses = append(q.WhereClauses, Clause{stmt, args})
	return q
}

func (c *Connection) Order(stmt string) *Query {
	return Q(c).Order(stmt)
}

func (q *Query) Order(stmt string) *Query {
	q.OrderClauses = append(q.OrderClauses, Clause{stmt, []interface{}{}})
	return q
}

func (c *Connection) Limit(limit int) *Query {
	return Q(c).Limit(limit)
}

func (q *Query) Limit(limit int) *Query {
	q.LimitResults = limit
	return q
}

func Q(c *Connection) *Query {
	return &Query{
		RawSQL:     &Clause{},
		Connection: c,
	}
}

func (q Query) ToSQL(model *Model, addColumns ...string) (string, []interface{}) {
	sb := NewSQLBuilder(q, model, addColumns...)
	return sb.String(), sb.Args()
}
