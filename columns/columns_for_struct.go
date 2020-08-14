package columns

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

// ForStruct returns a Columns instance for
// the struct passed in.
func ForStruct(s interface{}, tableName string) (columns Columns) {
	return ForStructWithAlias(s, tableName, "")
}

// ForStructWithAlias returns a Columns instance for the struct passed in.
// If the tableAlias is not empty, it will be used.
func ForStructWithAlias(s interface{}, tableName string, tableAlias string) (columns Columns) {
	columns = NewColumnsWithAlias(tableName, tableAlias)
	defer func() {
		if r := recover(); r != nil {
			columns = NewColumnsWithAlias(tableName, tableAlias)
			columns.Add("*")
		}
	}()
	st := reflect.TypeOf(s)
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() == reflect.Slice {
		st = st.Elem()
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
	}

	fieldMap := reflectx.NewMapper("").TypeMap(st)

	for _, i := range fieldMap.Index {
		popTags := TagsFor(i.Field)
		tag := popTags.Find("db")

		if !tag.Ignored() && !tag.Empty() {
			col := tag.Value

			// add writable or readable.
			tag := popTags.Find("rw")
			if !tag.Empty() {
				col = col + "," + tag.Value
			}

			cs := columns.Add(col)

			// add select clause.
			tag = popTags.Find("select")
			if !tag.Empty() {
				c := cs[0]
				c.SetSelectSQL(tag.Value)
			}
		}
	}

	return columns
}
