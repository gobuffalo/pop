package associations

import (
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/inflect"
)

type hasOneAssociation struct {
	ownedTableName string
	ownedModel     reflect.Value
	ownedType      reflect.Type
	ownerID        interface{}
	ownerName      string
	owner          interface{}
	fkID           string
	*associationSkipable
	*associationComposite
}

func init() {
	associationBuilders["has_one"] = hasOneAssociationBuilder
}

func hasOneAssociationBuilder(p associationParams) (Association, error) {
	// Validates if ownerIDField is nil, this association will be skipped.
	var skipped bool
	ownerID := p.modelValue.FieldByName("ID")
	if fieldIsNil(ownerID) {
		skipped = true
	}

	fval := p.modelValue.FieldByName(p.field.Name)
	return &hasOneAssociation{
		owner:          p.model,
		ownedTableName: p.popTags.Find("has_one").Value + "s",
		ownedModel:     fval,
		ownedType:      fval.Type(),
		ownerID:        ownerID.Interface(),
		ownerName:      p.modelType.Name(),
		fkID:           p.popTags.Find("fk_id").Value,
		associationSkipable: &associationSkipable{
			skipped: skipped,
		},
		associationComposite: &associationComposite{innerAssociations: p.innerAssociations},
	}, nil
}

func (h *hasOneAssociation) Kind() reflect.Kind {
	return h.ownedType.Kind()
}

func (h *hasOneAssociation) Interface() interface{} {
	if h.ownedModel.Kind() == reflect.Ptr {
		val := reflect.New(h.ownedType.Elem())
		h.ownedModel.Set(val)
		return h.ownedModel.Interface()
	}
	return h.ownedModel.Addr().Interface()
}

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (h *hasOneAssociation) Constraint() (string, []interface{}) {
	tn := flect.Underscore(h.ownerName)
	condition := fmt.Sprintf("%s_id = ?", tn)
	if h.fkID != "" {
		condition = fmt.Sprintf("%s = ?", h.fkID)
	}

	return condition, []interface{}{h.ownerID}
}

func (h *hasOneAssociation) AfterInterface() interface{} {
	if h.ownedModel.Kind() == reflect.Ptr {
		return h.ownedModel.Interface()
	}

	currentVal := h.ownedModel.Interface()
	zeroVal := reflect.Zero(h.ownedModel.Type()).Interface()
	if reflect.DeepEqual(zeroVal, currentVal) {
		return nil
	}

	return h.ownedModel.Addr().Interface()
}

func (h *hasOneAssociation) AfterSetup() error {
	ownerID := reflect.Indirect(reflect.ValueOf(h.owner)).FieldByName("ID").Interface()

	fval := h.ownedModel.FieldByName(h.ownerName + "ID")
	if fval.CanSet() {
		if n := nulls.New(fval.Interface()); n != nil {
			fval.Set(reflect.ValueOf(n.Parse(ownerID)))
		} else {
			fval.Set(reflect.ValueOf(ownerID))
		}
		return nil
	}

	return fmt.Errorf("could not set '%s' to '%s'", ownerID, fval)
}

func (h *hasOneAssociation) AfterProcess() AssociationStatement {
	otherIDField := "ID"
	fval := h.ownedModel.FieldByName(otherIDField)

	ownerIDField := "ID"
	ownerID := reflect.Indirect(reflect.ValueOf(h.owner)).FieldByName(ownerIDField).Interface()

	ids := []interface{}{ownerID}

	id := ""
	if fval.Type().Name() == "UUID" {
		id = fval.Interface().(uuid.UUID).String()
	} else {
		id = fmt.Sprint(fval.Interface())
	}
	if id != "0" && id != emptyUUID {
		ids = append(ids, id)
	} else {
		return AssociationStatement{
			Statement: "",
			Args:      []interface{}{},
		}
	}
	spew.Printf("has_one AfterProcess (ids):%v\n %v(%v) has_one %v:%v\n", ids, h.ownerName, ownerID, h.ownedType, h.ownedTableName)

	fk := h.fkID
	if fk == "" {
		fk = inflect.Underscore(h.ownerName) + "_id"
	}

	ret := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s = ?", h.ownedTableName, fk, otherIDField)

	return AssociationStatement{
		Statement: ret,
		Args:      ids,
	}
}
