package pop

// BelongsTo adds a "where" clause based on the "ID" of the
// "model" passed into it.
func (c *Connection) BelongsTo(model any) *Query {
	return Q(c).BelongsTo(model)
}

// BelongsToAs adds a "where" clause based on the "ID" of the
// "model" passed into it using an alias.
func (c *Connection) BelongsToAs(model any, as string) *Query {
	return Q(c).BelongsToAs(model, as)
}

// BelongsTo adds a "where" clause based on the "ID" of the
// "model" passed into it.
func (q *Query) BelongsTo(model any) *Query {
	m := NewModel(model, q.Connection.Context())
	q.Where(m.associationName()+" = ?", m.ID())
	return q
}

// BelongsToAs adds a "where" clause based on the "ID" of the
// "model" passed into it, using an alias.
func (q *Query) BelongsToAs(model any, as string) *Query {
	m := NewModel(model, q.Connection.Context())
	q.Where(as+" = ?", m.ID())
	return q
}

// BelongsToThrough adds a "where" clause that connects the "bt" model
// through the associated "thru" model.
func (c *Connection) BelongsToThrough(bt, thru any) *Query {
	return Q(c).BelongsToThrough(bt, thru)
}

// BelongsToThrough adds a "where" clause that connects the "bt" model
// through the associated "thru" model.
func (q *Query) BelongsToThrough(bt, thru any) *Query {
	q.belongsToThroughClauses = append(q.belongsToThroughClauses, belongsToThroughClause{
		BelongsTo: NewModel(bt, q.Connection.Context()),
		Through:   NewModel(thru, q.Connection.Context()),
	})
	return q
}
