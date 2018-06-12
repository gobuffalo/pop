package pop

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gobuffalo/pop/nulls"
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
	Ptr           bool
	tags          reflect.StructTag
}

type AssociationMetaType struct {
	*MetaType
	Root *MetaType
}

// Constraint returns the sql where clause to restrict query.
// for this association.
func (amt *AssociationMetaType) Constraint() string {
	column := (&Model{Value: amt.Root.Value.Interface()}).associationName()

	if amt.tags.Get("belongs_to") != "" {
		column = "id"
	}

	if amt.tags.Get("fk_id") != "" {
		column = amt.tags.Get("fk_id")
	}

	return fmt.Sprintf("%s in (?)", column)
}

// DependencyField returns the field that relates this association
// with another.
func (amt *AssociationMetaType) DependencyField() string {
	if amt.Root.IndirectType.Kind() == reflect.Slice {
		return fmt.Sprintf("%s%s", amt.Root.IndirectType.Elem().Name(), "ID")
	}

	return fmt.Sprintf("%s%s", amt.Root.IndirectType.Name(), "ID")
}

// Interface analogous to Interface method reflect.Value.
func (t *MetaType) Interface() interface{} {
	return t.Value.Interface()
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

// FieldByName returns field specified with name. If metatype is
// for an array and all elements are struct it will return all
// field values for all elements.
func (t *MetaType) FieldByName(field string) ReflectValues {
	allValues := ReflectValues{}
	if t.IndirectType.Kind() == reflect.Slice {
		for i := 0; i < t.IndirectValue.Len(); i++ {
			s := t.IndirectValue.Index(i)
			if s.Kind() == reflect.Struct {
				v := s.FieldByName(field)
				allValues = append(allValues, v)
			}
		}
		return allValues
	}

	if t.IndirectType.Kind() == reflect.Struct {
		v := t.IndirectValue.FieldByName(field)
		allValues = append(allValues, v)
	}

	return allValues
}

// Associations return metadata for associations defined in model t.
func (t *MetaType) Associations() []*AssociationMetaType {
	if t.IndirectType.Kind() == reflect.Slice {
		v := innerValueForSlice(t)
		return buildMetaType(v.Interface(), "").Associations()
	}

	atag := strings.Fields(associationTags)
	fields := t.ffields(func(tm *MetaType) bool {
		for _, tag := range atag {
			if tm.tags.Get(tag) != "" {
				return true
			}
		}
		return false
	})

	assos := []*AssociationMetaType{}
	for _, f := range fields {
		assos = append(assos, &AssociationMetaType{f, t})
	}
	return assos
}

// Association returns only association with the specified tag name.
func (t *MetaType) Association(kind string) []*AssociationMetaType {
	if t.IndirectType.Kind() == reflect.Slice {
		v := innerValueForSlice(t)
		return buildMetaType(v.Interface(), "").Association(kind)
	}

	fields := t.ffields(func(tm *MetaType) bool {
		return tm.tags.Get(kind) != ""
	})

	assos := []*AssociationMetaType{}
	for _, f := range fields {
		assos = append(assos, &AssociationMetaType{f, t})
	}
	return assos
}

// MakeSlice makes a new slice based on current MetaType type.
func (t *MetaType) MakeSlice() *MetaType {
	switch t.Type.Kind() {
	case reflect.Slice, reflect.Array:
		sliceVal := reflect.MakeSlice(t.Type, t.Value.Len(), t.Value.Cap())
		return buildMetaType(reflect.New(sliceVal.Type()).Interface(), "")
	default:
		sliceVal := reflect.MakeSlice(reflect.SliceOf(t.Type), 0, 0)
		return buildMetaType(reflect.New(sliceVal.Type()).Interface(), "")
	}
}

// MakeMap creates a map for a struct model type.
func (t *MetaType) MakeMap() *MetaType {
	switch t.IndirectType.Kind() {
	case reflect.Struct:
		return t.makeMapForStruct()
	case reflect.Slice, reflect.Array:
		return t.makeMapForSlice()
	default:
		mapMade := buildMap(reflect.TypeOf(""), t.Type)
		return buildMetaType(mapMade.Interface(), "")
	}
}

func (t *MetaType) makeMapForStruct() *MetaType {
	fIDs := t.ffields(func(f *MetaType) bool {
		return f.Name == "ID"
	})

	typ := t.IndirectType
	if t.Ptr {
		typ = reflect.PtrTo(t.IndirectType)
	}
	mapMade := buildMap(fIDs[0].Type, typ)

	return buildMetaType(mapMade.Interface(), "")
}

func (t *MetaType) makeMapForSlice() *MetaType {
	if t.IndirectType.Kind() != reflect.Slice {
		panic(fmt.Sprintf("%s is not an array type", t.Type))
	}

	v := innerValueForSlice(t)
	m := buildMetaType(v.Interface(), "").MakeMap()

	if reflect.Indirect(v).Kind() == reflect.Struct {
		if t.IndirectValue.Len() > 0 {
			for i := 0; i < t.IndirectValue.Len(); i++ {
				v := t.IndirectValue.Index(i)
				m.AppendMap(reflect.Indirect(v).FieldByName("ID").Interface(), v.Addr().Interface())
			}
		}
		return m
	}

	if t.IndirectValue.Len() > 0 {
		for i := 0; i < t.IndirectValue.Len(); i++ {
			v := t.IndirectValue.Index(i).Interface()
			m.AppendMap(fmt.Sprintf("%s", v), v)
		}
	}
	return m
}

// MakeMapWithField create a map for the specified MetaType with the key
// indicated by the field passed as parameter.
func (t *MetaType) MakeMapWithField(field string) *MetaType {
	if t.IndirectType.Kind() == reflect.Slice || t.IndirectType.Kind() == reflect.Array {
		val := innerValueForSlice(t)
		fval := reflect.Indirect(val).FieldByName(field)

		var m reflect.Value
		if n := nulls.New(fval.Interface()); n != nil {
			m = buildMap(reflect.TypeOf(n.Interface()), val.Type())
		} else {
			m = buildMap(fval.Type(), val.Type())
		}

		for i := 0; i < t.IndirectValue.Len(); i++ {
			v := t.IndirectValue.Index(0)
			f := v.FieldByName(field)
			if n := nulls.New(f.Interface()); n != nil {
				m.SetMapIndex(reflect.ValueOf(n.Interface()), v.Addr())
			} else {
				m.SetMapIndex(f, v.Addr())
			}
		}
		return buildMetaType(m.Interface(), "")
	}

	if t.IndirectType.Kind() == reflect.Struct {
		f := t.IndirectValue.FieldByName(field)
		m := buildMap(f.Type(), t.Type)
		m.SetMapIndex(f, t.Value)
		return buildMetaType(m, "")
	}
	return nil
}

func innerValueForSlice(t *MetaType) reflect.Value {
	var elemType reflect.Type
	var v reflect.Value

	// validates which slice inner element is.
	if t.Ptr {
		elemType = t.IndirectType.Elem()
	} else {
		elemType = t.Type.Elem()
	}

	// creates a new object type.
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		v = reflect.New(elemType)
	} else {
		v = reflect.New(elemType)
	}

	return v
}

// AppendMap appends to this map a value specified by key.
func (t *MetaType) AppendMap(key, value interface{}) {
	if t.Type.Kind() != reflect.Map {
		panic(fmt.Sprintf("%s is not a map", t.Type))
	}
	t.Value.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}

// MapKeys returns all keys registered in this map.
func (t *MetaType) MapKeys() ReflectValues {
	if t.Type.Kind() != reflect.Map {
		panic(fmt.Sprintf("%s is not a map", t.Type))
	}
	return ReflectValues(t.Value.MapKeys())
}

// MapValue returns a value for the specified key.
func (t *MetaType) MapValue(key interface{}) reflect.Value {
	if t.Type.Kind() != reflect.Map {
		panic(fmt.Sprintf("%s is not a map", t.Type))
	}
	return reflect.Indirect(t.Value.MapIndex(reflect.ValueOf(key)))
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

	mm.Ptr = mm.Type.Kind() == reflect.Ptr

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

func buildMap(key, val reflect.Type) reflect.Value {
	mapType := reflect.MapOf(key, val)
	mapVal := reflect.MakeMap(mapType)
	return mapVal
}

// ReflectValue is a helper to wrap some util functions for
// reflect.Value slice.
type ReflectValues []reflect.Value

// Interface parses all values to interface. It is an analougus
// to Interface method for reflect.Value.
func (rv ReflectValues) Interface() []interface{} {
	values := []interface{}{}
	for _, v := range rv {
		values = append(values, v.Interface())
	}
	return values
}
