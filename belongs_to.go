package pop

import (
	"fmt"

	"github.com/markbates/inflect"
)

func (c *Connection) BelongsTo(model interface{}) *Query {
	return Q(c).BelongsTo(model)
}

func (q *Query) BelongsTo(model interface{}) *Query {
	m := &Model{Value: model}
	tn := m.TableName()
	tn = inflect.Singularize(tn)
	args := []interface{}{m.ID()}
	q.WhereClauses = append(q.WhereClauses, Clause{fmt.Sprintf("%s_id = ?", tn), args})
	return q
}
