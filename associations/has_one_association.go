package associations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/markbates/inflect"
)

type hasOneAssociation struct {
	ownedModel reflect.Value
	ownedType  reflect.Type
	ownerID    interface{}
	ownerName  string
	fkID       string
}

func init() {
	associationBuilders["has_one"] = hasOneAssociationBuilder
}

func hasOneAssociationBuilder(p associationParams) (Association, error) {
	fval := p.modelValue.FieldByName(p.field.Name)
	return &hasOneAssociation{
		ownedModel: fval,
		ownedType:  fval.Type(),
		ownerID:    p.modelValue.FieldByName("ID").Interface(),
		ownerName:  p.modelType.Name(),
		fkID:       p.popTags.Find("fk_id").Value,
	}, nil
}

func (h *hasOneAssociation) TableName() string {
	return inflect.Tableize(h.ownedType.Name())
}

func (h *hasOneAssociation) Type() reflect.Kind {
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

// SQLConstraint returns the content for a where clause, and the args
// needed to execute it.
func (h *hasOneAssociation) SQLConstraint() (string, []interface{}) {
	tn := strings.ToLower(h.ownerName)
	condition := fmt.Sprintf("%s_id = ?", tn)
	if h.fkID != "" {
		condition = fmt.Sprintf("%s = ?", h.fkID)
	}

	return condition, []interface{}{h.ownerID}
}
