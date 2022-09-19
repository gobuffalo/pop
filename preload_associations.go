package pop

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6/internal/defaults"
	"github.com/gobuffalo/pop/v6/logging"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var validFieldRegexp = regexp.MustCompile(`^(([a-zA-Z0-9]*)(\.[a-zA-Z0-9]+)?)+$`)

// NewModelMetaInfo creates the meta info details for the model passed
// as a parameter.
func NewModelMetaInfo(model *Model) *ModelMetaInfo {
	mmi := &ModelMetaInfo{}
	mmi.Model = model
	mmi.init()
	return mmi
}

// NewAssociationMetaInfo creates the meta info details for the passed association.
func NewAssociationMetaInfo(fi *reflectx.FieldInfo) *AssociationMetaInfo {
	ami := &AssociationMetaInfo{}
	ami.FieldInfo = fi
	ami.init()
	return ami
}

// ModelMetaInfo a type to abstract all fields information regarding
// to a model. A model is representation of a table in the
// database.
type ModelMetaInfo struct {
	*reflectx.StructMap
	Model        *Model
	mapper       *reflectx.Mapper
	nestedFields map[string][]string
}

func (mmi *ModelMetaInfo) init() {
	m := reflectx.NewMapper("")
	mmi.mapper = m

	t := reflectx.Deref(reflect.TypeOf(mmi.Model.Value))
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = reflectx.Deref(t.Elem())
	}

	mmi.StructMap = m.TypeMap(t)
	mmi.nestedFields = make(map[string][]string)
}

func (mmi *ModelMetaInfo) iterate(fn func(reflect.Value)) {
	modelValue := reflect.Indirect(reflect.ValueOf(mmi.Model.Value))
	if modelValue.Kind() == reflect.Slice || modelValue.Kind() == reflect.Array {
		for i := 0; i < modelValue.Len(); i++ {
			fn(modelValue.Index(i))
		}
		return
	}
	fn(modelValue)
}

func (mmi *ModelMetaInfo) getDBFieldTaggedWith(value string) *reflectx.FieldInfo {
	for _, fi := range mmi.Index {
		if fi.Field.Tag.Get("db") == value {
			if len(fi.Children) > 0 {
				return fi.Children[0]
			}
			return fi
		}
	}
	return nil
}

func (mmi *ModelMetaInfo) preloadFields(fields ...string) ([]*reflectx.FieldInfo, error) {
	if len(fields) == 0 {
		return mmi.Index, nil
	}

	var preloadFields []*reflectx.FieldInfo
	for _, f := range fields {
		if !validFieldRegexp.MatchString(f) {
			return preloadFields, fmt.Errorf("association field '%s' does not match the format %s", f, "'<field>' or '<field>.<nested-field>'")
		}
		if strings.Contains(f, ".") {
			fname := f[:strings.Index(f, ".")]
			mmi.nestedFields[fname] = append(mmi.nestedFields[fname], f[strings.Index(f, ".")+1:])
			f = f[:strings.Index(f, ".")]
		}

		preloadField := mmi.GetByPath(f)
		if preloadField == nil {
			return preloadFields, fmt.Errorf("field %s does not exist in model %s", f, mmi.Model.TableName())
		}

		var exist bool
		for _, pf := range preloadFields {
			if pf.Path == preloadField.Path {
				exist = true
			}
		}
		if !exist {
			preloadFields = append(preloadFields, preloadField)
		}
	}
	return preloadFields, nil
}

// AssociationMetaInfo a type to abstract all field information
// regarding to an association. An association is a field
// that has defined a tag like 'has_many', 'belongs_to',
// 'many_to_many' and 'has_one'.
type AssociationMetaInfo struct {
	*reflectx.FieldInfo
	*reflectx.StructMap
}

func (ami *AssociationMetaInfo) init() {
	mapper := reflectx.NewMapper("")
	t := reflectx.Deref(ami.FieldInfo.Field.Type)
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = reflectx.Deref(t.Elem())
	}

	ami.StructMap = mapper.TypeMap(t)
}

func (ami *AssociationMetaInfo) toSlice() reflect.Value {
	ft := reflectx.Deref(ami.Field.Type)
	var vt reflect.Value
	if ft.Kind() == reflect.Slice || ft.Kind() == reflect.Array {
		vt = reflect.New(ft)
	} else {
		vt = reflect.New(reflect.SliceOf(ft))
	}
	return vt
}

func (ami *AssociationMetaInfo) getDBFieldTaggedWith(value string) *reflectx.FieldInfo {
	for _, fi := range ami.StructMap.Index {
		if fi.Field.Tag.Get("db") == value {
			if len(fi.Children) > 0 {
				return fi.Children[0]
			}
			return fi
		}
	}
	return nil
}

func (ami *AssociationMetaInfo) targetPrimaryID() string {
	pid := ami.Field.Tag.Get("primary_id")
	switch {
	case pid == "":
		return "id"
	case ami.getDBFieldTaggedWith(pid) != nil:
		return pid
	case ami.GetByPath(pid) != nil:
		return ami.GetByPath(pid).Field.Tag.Get("db")
	default:
		return ""
	}
}

func (ami *AssociationMetaInfo) fkName() string {
	t := ami.Field.Type
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = reflectx.Deref(t.Elem())
	}
	fkName := fmt.Sprintf("%s%s", flect.Underscore(flect.Singularize(t.Name())), "_id")
	fkNameTag := flect.Underscore(ami.Field.Tag.Get("fk_id"))
	return defaults.String(fkNameTag, fkName)
}

// preload is the query mode used to load associations from database
// similar to the active record default approach on Rails.
func preload(tx *Connection, model interface{}, fields ...string) error {
	mmi := NewModelMetaInfo(NewModel(model, tx.Context()))

	preloadFields, err := mmi.preloadFields(fields...)
	if err != nil {
		return err
	}

	var associations []*AssociationMetaInfo
	for _, fieldInfo := range preloadFields {
		if isFieldAssociation(fieldInfo.Field) && fieldInfo.Parent.Name == "" {
			associations = append(associations, NewAssociationMetaInfo(fieldInfo))
		}
	}

	for _, asoc := range associations {
		if asoc.Field.Tag.Get("has_many") != "" {
			err := preloadHasMany(tx, asoc, mmi)
			if err != nil {
				return err
			}
		}

		if asoc.Field.Tag.Get("has_one") != "" {
			err := preloadHasOne(tx, asoc, mmi)
			if err != nil {
				return err
			}
		}

		if asoc.Field.Tag.Get("belongs_to") != "" {
			err := preloadBelongsTo(tx, asoc, mmi)
			if err != nil {
				return err
			}
		}

		if asoc.Field.Tag.Get("many_to_many") != "" {
			err := preloadManyToMany(tx, asoc, mmi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isFieldAssociation(field reflect.StructField) bool {
	for _, associationLabel := range []string{"has_many", "has_one", "belongs_to", "many_to_many"} {
		if field.Tag.Get(associationLabel) != "" {
			return true
		}
	}
	return false
}

func preloadHasMany(tx *Connection, asoc *AssociationMetaInfo, mmi *ModelMetaInfo) error {
	// 1) get all associations ids.
	// 1.1) In here I pick ids from model meta info directly.
	idField := asoc.getDBFieldTaggedWith(asoc.targetPrimaryID())
	ids := []interface{}{}
	mmi.Model.iterate(func(m *Model) error {
		if idField.Path == "ID" {
			ids = append(ids, m.ID())
			return nil
		}

		v, err := m.fieldByName(idField.Path)
		if err != nil {
			return err
		}

		ids = append(ids, normalizeValue(v.Interface()))
		return nil
	})

	if len(ids) == 0 {
		return nil
	}

	// 2) load all associations constraint by model ids.
	fk := asoc.Field.Tag.Get("fk_id")
	if fk == "" {
		fk = mmi.Model.associationName()
	}

	q := tx.Q()
	q.eager = false
	q.eagerFields = []string{}

	slice := asoc.toSlice()

	if strings.TrimSpace(asoc.Field.Tag.Get("order_by")) != "" {
		q.Order(asoc.Field.Tag.Get("order_by"))
	}

	err := q.Where(fmt.Sprintf("%s in (?)", fk), ids).All(slice.Interface())
	if err != nil {
		return err
	}

	// 2.1) load all nested associations from this assoc.
	if asocNestedFields, ok := mmi.nestedFields[asoc.Path]; ok {
		for _, asocNestedField := range asocNestedFields {
			if err := preload(tx, slice.Interface(), asocNestedField); err != nil {
				return err
			}
		}
	}

	// 3) iterate over every model and fill it with the assoc.
	foreignField := asoc.getDBFieldTaggedWith(fk)
	mmi.iterate(func(mvalue reflect.Value) {
		for i := 0; i < slice.Elem().Len(); i++ {
			asocValue := slice.Elem().Index(i)
			valueField := reflect.Indirect(mmi.mapper.FieldByName(asocValue, foreignField.Path))
			if mmi.mapper.FieldByName(mvalue, idField.Path).Interface() == valueField.Interface() ||
				reflect.DeepEqual(mmi.mapper.FieldByName(mvalue, idField.Path), valueField) {
				// IMPORTANT
				//
				// FieldByName will initialize the value. It is important that this happens AFTER
				// we checked whether the field should be set. Otherwise, we'll set a zero value!
				//
				// This is most likely the reason for https://github.com/gobuffalo/pop/issues/139
				modelAssociationField := mmi.mapper.FieldByName(mvalue, asoc.Name)
				switch {
				case modelAssociationField.Kind() == reflect.Slice || modelAssociationField.Kind() == reflect.Array:
					modelAssociationField.Set(reflect.Append(modelAssociationField, asocValue))
				case modelAssociationField.Kind() == reflect.Ptr:
					modelAssociationField.Elem().Set(reflect.Append(modelAssociationField.Elem(), asocValue))
				default:
					modelAssociationField.Set(asocValue)
				}
			}
		}
	})

	return nil
}

func preloadHasOne(tx *Connection, asoc *AssociationMetaInfo, mmi *ModelMetaInfo) error {
	// 1) get all associations ids.
	idField := asoc.getDBFieldTaggedWith(asoc.targetPrimaryID())
	ids := []interface{}{}
	mmi.Model.iterate(func(m *Model) error {
		if idField.Path == "ID" {
			ids = append(ids, m.ID())
			return nil
		}

		v, err := m.fieldByName(idField.Path)
		if err != nil {
			return err
		}

		ids = append(ids, normalizeValue(v.Interface()))
		return nil
	})

	if len(ids) == 0 {
		return nil
	}

	// 2) load all associations constraint by model ids.
	fk := asoc.Field.Tag.Get("fk_id")
	if fk == "" {
		fk = mmi.Model.associationName()
	}

	q := tx.Q()
	q.eager = false
	q.eagerFields = []string{}

	slice := asoc.toSlice()
	err := q.Where(fmt.Sprintf("%s in (?)", fk), ids).All(slice.Interface())
	if err != nil {
		return err
	}

	// 2.1) load all nested associations from this assoc.
	if asocNestedFields, ok := mmi.nestedFields[asoc.Path]; ok {
		for _, asocNestedField := range asocNestedFields {
			if err := preload(tx, slice.Interface(), asocNestedField); err != nil {
				return err
			}
		}
	}

	//  3) iterate over every model and fill it with the assoc.
	foreignField := asoc.getDBFieldTaggedWith(fk)
	mmi.iterate(func(mvalue reflect.Value) {
		for i := 0; i < slice.Elem().Len(); i++ {
			asocValue := slice.Elem().Index(i)
			if mmi.mapper.FieldByName(mvalue, idField.Path).Interface() == mmi.mapper.FieldByName(asocValue, foreignField.Path).Interface() ||
				reflect.DeepEqual(mmi.mapper.FieldByName(mvalue, idField.Path), mmi.mapper.FieldByName(asocValue, foreignField.Path)) {
				// IMPORTANT
				//
				// FieldByName will initialize the value. It is important that this happens AFTER
				// we checked whether the field should be set. Otherwise, we'll set a zero value!
				//
				// This is most likely the reason for https://github.com/gobuffalo/pop/issues/139
				modelAssociationField := mmi.mapper.FieldByName(mvalue, asoc.Name)
				switch {
				case modelAssociationField.Kind() == reflect.Slice || modelAssociationField.Kind() == reflect.Array:
					modelAssociationField.Set(reflect.Append(modelAssociationField, asocValue))
				case modelAssociationField.Kind() == reflect.Ptr:
					modelAssociationField.Elem().Set(asocValue)
				default:
					modelAssociationField.Set(asocValue)
				}
			}
		}
	})

	return nil
}

func preloadBelongsTo(tx *Connection, asoc *AssociationMetaInfo, mmi *ModelMetaInfo) error {
	// 1) get all associations ids.
	fi := mmi.getDBFieldTaggedWith(asoc.fkName())
	if fi == nil {
		fi = mmi.getDBFieldTaggedWith(fmt.Sprintf("%s%s", flect.Underscore(asoc.Path), "_id"))
	}

	fkids := []interface{}{}
	mmi.iterate(func(val reflect.Value) {
		if !isFieldNilPtr(val, fi) {
			fkids = append(fkids, mmi.mapper.FieldByName(val, fi.Path).Interface())
		}
	})

	if len(fkids) == 0 {
		return nil
	}

	// 2) load all associations constraint by association fields ids.
	fk := asoc.targetPrimaryID()

	q := tx.Q()
	q.eager = false
	q.eagerFields = []string{}

	slice := asoc.toSlice()
	err := q.Where(fmt.Sprintf("%s in (?)", fk), fkids).All(slice.Interface())
	if err != nil {
		return err
	}

	// 2.1) load all nested associations from this assoc.
	if asocNestedFields, ok := mmi.nestedFields[asoc.Path]; ok {
		for _, asocNestedField := range asocNestedFields {
			if err := preload(tx, slice.Interface(), asocNestedField); err != nil {
				return err
			}
		}
	}

	// 3) iterate over every model and fill it with the assoc.
	idField := mmi.getDBFieldTaggedWith(asoc.targetPrimaryID())
	mmi.iterate(func(mvalue reflect.Value) {
		if isFieldNilPtr(mvalue, fi) {
			return
		}
		for i := 0; i < slice.Elem().Len(); i++ {
			asocValue := slice.Elem().Index(i)
			fkField := reflect.Indirect(mmi.mapper.FieldByName(mvalue, fi.Path))
			field := mmi.mapper.FieldByName(asocValue, idField.Path)
			if fkField.Interface() == field.Interface() || reflect.DeepEqual(fkField, field) {
				// IMPORTANT
				//
				// FieldByName will initialize the value. It is important that this happens AFTER
				// we checked whether the field should be set. Otherwise, we'll set a zero value!
				//
				// This is most likely the reason for https://github.com/gobuffalo/pop/issues/139
				modelAssociationField := mmi.mapper.FieldByName(mvalue, asoc.Name)
				switch {
				case modelAssociationField.Kind() == reflect.Slice || modelAssociationField.Kind() == reflect.Array:
					modelAssociationField.Set(reflect.Append(modelAssociationField, asocValue))
				case modelAssociationField.Kind() == reflect.Ptr:
					modelAssociationField.Elem().Set(asocValue)
				default:
					modelAssociationField.Set(asocValue)
				}
			}
		}
	})

	return nil
}

func preloadManyToMany(tx *Connection, asoc *AssociationMetaInfo, mmi *ModelMetaInfo) error {
	// 1) get all associations ids.
	// 1.1) In here I pick ids from model meta info directly.
	idField := asoc.getDBFieldTaggedWith(asoc.targetPrimaryID())
	ids := []interface{}{}
	mmi.Model.iterate(func(m *Model) error {
		if idField.Path == "ID" {
			ids = append(ids, m.ID())
			return nil
		}

		v, err := m.fieldByName(idField.Path)
		if err != nil {
			return err
		}

		ids = append(ids, normalizeValue(v.Interface()))
		return nil
	})

	if len(ids) == 0 {
		return nil
	}

	// 2) load all associations.
	// 2.1) In here I pick the label name from association.
	manyToManyTableName := asoc.Field.Tag.Get("many_to_many")
	modelAssociationName := mmi.Model.associationName()
	assocFkName := asoc.fkName()

	if asoc.Field.Tag.Get("primary_id") != "" {
		modelAssociationName = asoc.Field.Tag.Get("primary_id")
	}

	sql := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s in (?)", modelAssociationName, assocFkName, manyToManyTableName, modelAssociationName)
	sql, args, _ := sqlx.In(sql, ids)
	sql = tx.Dialect.TranslateSQL(sql)
	log(logging.SQL, sql, args...)

	cn, err := tx.Store.Transaction()
	if err != nil {
		return err
	}

	rows, err := cn.Queryx(sql, args...)
	if err != nil {
		return err
	}

	mapAssoc := map[string][]interface{}{}
	fkids := []interface{}{}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return err
		}
		if len(row) > 0 {
			if _, ok := row[0].([]uint8); ok { // -> it's UUID
				row[0] = string(row[0].([]uint8))
			}
			if _, ok := row[1].([]uint8); ok { // -> it's UUID
				row[1] = string(row[1].([]uint8))
			}
			key := fmt.Sprintf("%v", row[0])
			mapAssoc[key] = append(mapAssoc[key], row[1])
			fkids = append(fkids, row[1])
		}
	}

	q := tx.Q()
	q.eager = false
	q.eagerFields = []string{}

	if strings.TrimSpace(asoc.Field.Tag.Get("order_by")) != "" {
		q.Order(asoc.Field.Tag.Get("order_by"))
	}

	slice := asoc.toSlice()
	q.Where("id in (?)", fkids).All(slice.Interface())

	// 2.2) load all nested associations from this assoc.
	if asocNestedFields, ok := mmi.nestedFields[asoc.Path]; ok {
		for _, asocNestedField := range asocNestedFields {
			if err := preload(tx, slice.Interface(), asocNestedField); err != nil {
				return err
			}
		}
	}

	// 3) iterate over every model and fill it with the assoc.
	mmi.iterate(func(mvalue reflect.Value) {
		id := mmi.mapper.FieldByName(mvalue, "ID").Interface()
		if assocFkIds, ok := mapAssoc[fmt.Sprintf("%v", id)]; ok {
			for i := 0; i < slice.Elem().Len(); i++ {
				asocValue := slice.Elem().Index(i)
				for _, fkid := range assocFkIds {
					if fmt.Sprintf("%v", fkid) == fmt.Sprintf("%v", mmi.mapper.FieldByName(asocValue, "ID").Interface()) {
						// IMPORTANT
						//
						// FieldByName will initialize the value. It is important that this happens AFTER
						// we checked whether the field should be set. Otherwise, we'll set a zero value!
						//
						// This is most likely the reason for https://github.com/gobuffalo/pop/issues/139
						modelAssociationField := mmi.mapper.FieldByName(mvalue, asoc.Name)
						modelAssociationField.Set(reflect.Append(modelAssociationField, asocValue))
					}
				}
			}
		}
	})

	return nil
}

func isFieldNilPtr(val reflect.Value, fi *reflectx.FieldInfo) bool {
	fieldValue := reflectx.FieldByIndexesReadOnly(val, fi.Index)
	return fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()
}

func normalizeValue(val interface{}) interface{} {
	switch t := val.(type) {
	case nulls.String:
		return t.String
	case nulls.Float64:
		return t.Float64
	case nulls.Int:
		return t.Int
	case nulls.Time:
		return t.Time
	default:
		return t
	}
}
