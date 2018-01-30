package generate

import (
	"fmt"
	"strings"

	"github.com/markbates/inflect"
)

type attribute struct {
	Name         inflect.Name
	OriginalType string
	GoType       string
	Nullable     bool
}

func (a attribute) String() string {
	return fmt.Sprintf("\t%s %s `%s:\"%s\" db:\"%s\"`", a.Name.Model(), a.GoType, structTag, a.Name, a.Name)
}

func newAttribute(base string, model *model) attribute {
	col := strings.Split(base, ":")
	if len(col) == 1 {
		col = append(col, "string")
	}

	nullable := nrx.MatchString(col[1])
	if !model.HasNulls && nullable {
		model.HasNulls = true
		model.Imports = append(model.Imports, "github.com/markbates/pop/nulls")
	}

	if strings.HasPrefix(col[1], "slices.") {
		model.Imports = append(model.Imports, "github.com/markbates/pop/slices")
	}

	if !model.HasUUID && col[1] == "uuid" {
		model.HasUUID = true
		model.Imports = append(model.Imports, "github.com/satori/go.uuid")
	}

	a := attribute{
		Name:         inflect.Name(col[0]),
		OriginalType: col[1],
		GoType:       colType(col[1]),
		Nullable:     nullable,
	}

	return a
}
