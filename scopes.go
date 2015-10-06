package pop

type ScopeFunc func(q *Query) *Query

func (q *Query) Scope(sf ScopeFunc) *Query {
	return sf(q)
}

func (c *Connection) Scope(sf ScopeFunc) *Query {
	return Q(c).Scope(sf)
}
