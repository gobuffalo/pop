package pop

import (
	"fmt"
	"reflect"

	"github.com/kr/pretty"
	"github.com/markbates/going/defaults"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

func (q *Query) findAssociations(m *Model) error {
	for _, as := range q.withAssociations {
		rv := reflect.Indirect(reflect.ValueOf(m.Value))
		rt := rv.Type()
		id := rv.FieldByName("ID")
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			tag := defaults.String(field.Tag.Get("assoc"), field.Name)
			assocKey := field.Tag.Get("assockey")
			if assocKey == "" {
				assocKey = fmt.Sprintf("%s_id", inflect.Underscore(inflect.Singularize(rt.Name())))
			}
			if tag == as {
				rf := rv.FieldByName(field.Name)
				if rf.Kind() == reflect.Slice {
					fmt.Println("a slice of ", field.Name)
					qq := q.Connection.Where(fmt.Sprintf("%s = ?", assocKey), id.Interface())
					so := reflect.SliceOf(rf.Type())
					// err := qq.All(&so)
					// i := so..Interface()
					pretty.Println("### so.Name() ->", so.Name())
					err := qq.All(so)
					if err != nil {
						return errors.WithStack(err)
					}
				}

				// we found the thing we were looking for, so move on with other assocations
				break
			}
		}
	}

	return m.afterFind(q.Connection)
}
