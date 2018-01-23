package columns

import (
	"reflect"
)

// ColumnsForStruct returns a Columns instance for
// the struct passed in.

func ColumnsForStruct(s interface{}, tableName string) (columns Columns) {
	return ColumnsForStructWithAlias(s, tableName, "")
}

func ColumnsForStructWithAlias(s interface{}, tableName string, tableAlias string) (columns Columns) {
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

	fieldCount := st.NumField()

	for i := 0; i < fieldCount; i++ {
		field := st.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = field.Name
		}

		if tag != "-" {
			rw := field.Tag.Get("rw")
			if rw != "" {
				tag = tag + "," + rw
			}

			if tag == "preload" {
				cs := columns.Add(field.Name)
				c := cs[0]
				c.SetPreload()
				continue
			}

			cs := columns.Add(tag)
			c := cs[0]
			tag = field.Tag.Get("select")
			if tag != "" {
				c.SetSelectSQL(tag)
			}
		}
	}

	return columns
}
