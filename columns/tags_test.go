package columns_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/pop/v6/columns"
)

func Test_Tags_TagsFor(t *testing.T) {
	r := require.New(t)

	typ := reflect.TypeFor[foo]()
	f, _ := typ.FieldByName("FirstName")
	tags := columns.TagsFor(f)

	r.Equal(len(tags), 2)
	r.Equal(tags.Find("db").Value, "first_name")
	r.Equal(tags.Find("select").Value, "first_name as f")
}
