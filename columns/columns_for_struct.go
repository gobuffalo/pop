package columns

import "reflect"

// ColumnsForStruct returns a Columns instance for
// the struct passed in.
func ColumnsForStruct(s interface{}, tableName string) (columns Columns) {
	columns = NewColumns(tableName)
	defer func() {
		if r := recover(); r != nil {
			columns = NewColumns(tableName)
			columns.Add("*")
		}
	}()
	st := reflect.TypeOf(s)
	if st.Kind() == reflect.Ptr {
		st = reflect.ValueOf(s).Elem().Type()
	}
	if st.Kind() == reflect.Slice {
		v := reflect.ValueOf(s)
		t := v.Type()
		x := t.Elem().Elem()

		n := reflect.New(x)
		return ColumnsForStruct(n.Interface(), tableName)
	}

	field_count := st.NumField()

	for i := 0; i < field_count; i++ {
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
