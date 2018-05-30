package pop

import (
	"reflect"
	"strings"
)

var associationTags = "belongs_to has_one has_many many_to_many"

// MetaType is a meta data for a model.
type MetaType struct {
	Name          string
	Type          reflect.Type
	Value         reflect.Value
	IndirectValue reflect.Value
	IndirectType  reflect.Type
	Fields        []*MetaType
	tags          reflect.StructTag
}

func (t *MetaType) ffields(fn func(*MetaType) bool) []*MetaType {
	ffields := []*MetaType{}
	for _, f := range t.Fields {
		if fn(f) {
			ffields = append(ffields, f)
		}
	}
	return ffields
}

// Associations return metadata for associations defined in model t.
func (t *MetaType) Associations() []*MetaType {
	atag := strings.Fields(associationTags)
	return t.ffields(func(tm *MetaType) bool {
		for _, tag := range atag {
			if tm.tags.Get(tag) != "" {
				return true
			}
		}
		return false
	})
}

// MakeSlice makes a new slice based on current MetaType type.
func (t *MetaType) MakeSlice() *MetaType {
	switch t.Type.Kind() {
	case reflect.Slice, reflect.Array:
		sliceVal := reflect.MakeSlice(t.Type, t.Value.Len(), t.Value.Cap())
		return buildMetaType(sliceVal.Interface(), "")
	default:
		sliceVal := reflect.MakeSlice(reflect.SliceOf(t.Type), 0, 0)
		return buildMetaType(sliceVal.Interface(), "")
	}
}

// MakeMap creates a map for a struct model type.
func (t *MetaType) MakeMap() *MetaType {
	fIDs := t.ffields(func(f *MetaType) bool {
		return f.Name == "ID"
	})

	switch t.IndirectType.Kind() {
	case reflect.Struct:
		mapType := reflect.MapOf(fIDs[0].Type, t.IndirectType)
		mapVal := reflect.MakeMap(mapType)
		return buildMetaType(mapVal.Interface(), "")
	default:
		return nil
	}
}

// Meta a helper function that wraps reflection data type
// for a model.
func (m *Model) Meta() *MetaType {
	return buildMetaType(m.Value, "")
}

func buildMetaType(model interface{}, name string) *MetaType {
	v := reflect.ValueOf(model)
	iv := reflect.Indirect(v)

	mm := &MetaType{
		Type:          v.Type(),
		IndirectType:  iv.Type(),
		Name:          name,
		Value:         v,
		IndirectValue: iv,
	}
	if mm.Name == "" {
		mm.Name = mm.IndirectType.Name()
	}

	if iv.Kind() == reflect.Struct {
		mm.Fields = []*MetaType{}
		for i := 0; i < iv.NumField(); i++ {
			ft := iv.Type().Field(i)
			fv := iv.Field(i)
			if fv.CanInterface() {
				field := buildMetaType(fv.Interface(), ft.Name)
				field.tags = ft.Tag
				mm.Fields = append(mm.Fields, field)
			}
		}
	}

	return mm
}
