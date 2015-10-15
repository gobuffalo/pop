package pop

import "fmt"

func (c *Connection) BelongsTo(model interface{}) *Query {
	return Q(c).BelongsTo(model)
}

func (q *Query) BelongsTo(model interface{}) *Query {
	m := &Model{Value: model}
	q.Where(fmt.Sprintf("%s = ?", m.AssociationName()), m.ID())
	return q
}

func (c *Connection) BelongsToThrough(bt, thru interface{}) *Query {
	return Q(c).BelongsToThrough(bt, thru)
}

func (q *Query) BelongsToThrough(bt, thru interface{}) *Query {
	q.BelongsToThroughClauses = append(q.BelongsToThroughClauses, BelongsToThroughClause{
		BelongsTo: &Model{Value: bt},
		Through:   &Model{Value: thru},
	})
	return q
}
