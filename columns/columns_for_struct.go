package columns

import (
	"reflect"

	"github.com/markbates/oncer"
)

// ColumnsForStruct returns a Columns instance for
// the struct passed in.
//
// Deprecated: use ForStruct instead.
func ColumnsForStruct(s interface{}, tableName string) (columns Columns) {
	oncer.Deprecate(0, "columns.ColumnsForStruct", "Use columns.ForStruct instead.")
	return ForStruct(s, tableName)
}

// ColumnsForStructWithAlias returns a Columns instance for the struct passed in.
// If the tableAlias is not empty, it will be used.
//
// Deprecated: use ForStructWithAlias instead.
func ColumnsForStructWithAlias(s interface{}, tableName string, tableAlias string) (columns Columns) {
	oncer.Deprecate(0, "columns.ColumnsForStructWithAlias", "Use columns.ForStructWithAlias instead.")
	return ForStructWithAlias(s, tableName, tableAlias)
}

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

	do := func(reflect.Type) {}
	do = func(v reflect.Type) {
		fieldCount := v.NumField()

		for i := 0; i < fieldCount; i++ {
			field := v.Field(i)

			// Is bad approach fill for all structure data db tag,
			// this not allow to us use composition at all as we fill not existen data
			// from time.Time struct or others.
			//
			// For better we can support json tag withoud not filling db tag.
			// if we want make restrict access to filed, we should use db tag with - value
			popTags := TagsForReal(field)
			tagDB := popTags.Find("db")
			tagJSON := popTags.Find("json")

			tag := tagDB
			if tagDB.Empty() {
				tag = tagJSON
			}

			// support composition structures.
			// If no tag json and db inside structures this run just ignor not important fields.
			if field.Type.Kind() == reflect.Struct {
				do(v.Field(i).Type)
			}

			if tag.Ignored() || tag.Empty() {
				continue
			}

			col := tag.Value

			// add writable or readable.
			tag = popTags.Find("rw")
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
	do(st)

	return columns
}
