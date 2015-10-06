package pop

import (
	"fmt"
	"log"
	"strings"
)

type Query struct {
	RawSQL       *Clause
	LimitResults int
	WhereClauses Clauses
	OrderClauses Clauses
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
	sql := q.RawSQL.Fragment
	args := q.RawSQL.Arguments
	if sql == "" {
		sql, args = q.buildSQL(model, addColumns...)
	}
	sql = q.Connection.Dialect.TranslateSQL(sql)
	if Debug {
		x := fmt.Sprintf("[POP]: %s", sql)

		if len(args) > 0 {
			xargs := make([]string, len(args))
			for i, a := range args {
				switch a.(type) {
				case string:
					xargs[i] = fmt.Sprintf("%q", a)
				default:
					xargs[i] = fmt.Sprintf("%v", a)
				}
			}
			x = fmt.Sprintf("%s | %s", x, xargs)
		}

		log.Println(x)

	}
	return sql, args
}

func (q Query) buildSQL(model *Model, addColumns ...string) (sql string, args []interface{}) {
	tableName := model.TableName()
	cols := q.buildColumns(model, addColumns...)
	sql = fmt.Sprintf("SELECT %s FROM %s as %s", cols.Readable().SelectString(), tableName, strings.Replace(tableName, ".", "_", -1))
	if len(q.WhereClauses) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, q.WhereClauses.Join(" AND "))
		for _, arg := range q.WhereClauses.Args() {
			args = append(args, arg)
		}
	}
	if len(q.OrderClauses) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, q.OrderClauses.Join(", "))
		for _, arg := range q.OrderClauses.Args() {
			args = append(args, arg)
		}
	}
	if q.LimitResults > 0 && q.Paginator == nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, q.LimitResults)
	}
	if q.Paginator != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, q.Paginator.PerPage)
		sql = fmt.Sprintf("%s OFFSET %d", sql, q.Paginator.Offset)
	}
	return sql, args
}

var columnCache = map[string]Columns{}

func (q Query) buildColumns(model *Model, addColumns ...string) Columns {
	tableName := model.TableName()
	acl := len(addColumns)
	if acl <= 0 {
		cols, ok := columnCache[tableName]
		if ok {
			return cols
		}
		cols = ColumnsForStruct(model.Value)
		columnCache[tableName] = cols
		return cols
	} else {
		cols := NewColumns()
		cols.Add(addColumns...)
		return cols
	}
}
