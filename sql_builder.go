package pop

import (
	"fmt"
	"log"
	"strings"

	. "github.com/markbates/pop/columns"
)

type SQLBuilder struct {
	Query      Query
	Model      *Model
	AddColumns []string
	sql        string
	args       []interface{}
}

func NewSQLBuilder(q Query, m *Model, addColumns ...string) *SQLBuilder {
	return &SQLBuilder{
		Query:      q,
		Model:      m,
		AddColumns: addColumns,
		args:       []interface{}{},
	}
}

func (sq *SQLBuilder) String() string {
	if sq.sql == "" {
		sq.compile()
	}
	sq.log()
	return sq.sql
}

func (sq *SQLBuilder) Args() []interface{} {
	if len(sq.args) == 0 {
		if len(sq.Query.RawSQL.Arguments) > 0 {
			sq.args = sq.Query.RawSQL.Arguments
		} else {
			sq.compile()
		}
	}
	return sq.args
}

func (sq *SQLBuilder) compile() {
	if sq.sql == "" {
		if sq.Query.RawSQL.Fragment != "" {
			sq.sql = sq.Query.RawSQL.Fragment
		} else {
			sq.sql = sq.buildSQL()
		}
		sq.sql = sq.Query.Connection.Dialect.TranslateSQL(sq.sql)
	}
	sq.log()
}

func (sq *SQLBuilder) log() {
	if Debug {
		args := sq.args
		x := fmt.Sprintf("[POP]: %s", sq.sql)

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
}

func (sq *SQLBuilder) buildSQL() string {
	cols := sq.buildColumns()

	fc := sq.buildFromClauses()

	sql := fmt.Sprintf("SELECT %s FROM %s", cols.Readable().SelectString(), fc)

	sql = sq.buildWhereClauses(sql)
	sql = sq.buildOrderClauses(sql)
	sql = sq.buildPaginationClauses(sql)

	return sql
}

func (sq *SQLBuilder) buildFromClauses() FromClauses {
	models := []*Model{
		sq.Model,
	}
	for _, mc := range sq.Query.BelongsToThroughClauses {
		models = append(models, mc.Through)
	}

	fc := sq.Query.FromClauses
	for _, m := range models {
		tableName := m.TableName()
		fc = append(fc, FromClause{
			From: tableName,
			As:   strings.Replace(tableName, ".", "_", -1),
		})
	}

	return fc
}

func (sq *SQLBuilder) buildWhereClauses(sql string) string {
	mcs := sq.Query.BelongsToThroughClauses
	for _, mc := range mcs {
		sq.Query.Where(fmt.Sprintf("%s.%s = ?", mc.Through.TableName(), mc.BelongsTo.AssociationName()), mc.BelongsTo.ID())
		sq.Query.Where(fmt.Sprintf("%s.id = %s.%s", sq.Model.TableName(), mc.Through.TableName(), sq.Model.AssociationName()))
	}

	wc := sq.Query.WhereClauses
	if len(wc) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, wc.Join(" AND "))
		for _, arg := range wc.Args() {
			sq.args = append(sq.args, arg)
		}
	}
	return sql
}

func (sq *SQLBuilder) buildOrderClauses(sql string) string {
	oc := sq.Query.OrderClauses
	if len(oc) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, oc.Join(", "))
		for _, arg := range oc.Args() {
			sq.args = append(sq.args, arg)
		}
	}
	return sql
}

func (sq *SQLBuilder) buildPaginationClauses(sql string) string {
	if sq.Query.LimitResults > 0 && sq.Query.Paginator == nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, sq.Query.LimitResults)
	}
	if sq.Query.Paginator != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, sq.Query.Paginator.PerPage)
		sql = fmt.Sprintf("%s OFFSET %d", sql, sq.Query.Paginator.Offset)
	}
	return sql
}

var columnCache = map[string]Columns{}

func (sq *SQLBuilder) buildColumns() Columns {
	tableName := sq.Model.TableName()
	acl := len(sq.AddColumns)
	if acl <= 0 {
		cols, ok := columnCache[tableName]
		if ok {
			return cols
		}
		cols = ColumnsForStruct(sq.Model.Value, tableName)
		columnCache[tableName] = cols
		return cols
	} else {
		cols := NewColumns("")
		cols.Add(sq.AddColumns...)
		return cols
	}
}
