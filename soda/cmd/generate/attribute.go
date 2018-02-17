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
	return fmt.Sprintf("\t%s %s `%s:\"%s\" db:\"%s\"`", a.Name.Camel(), a.GoType, structTag, a.Name.UnderSingular(), a.Name.UnderSingular())
}

func (a attribute) IsValidable() bool {
	return a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int"
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

	if !model.HasSlices && strings.HasPrefix(col[1], "slices.") {
		model.HasSlices = true
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

func colType(s string) string {
	switch strings.ToLower(s) {
	case "text":
		return "string"
	case "time", "timestamp", "datetime":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	case "decimal", "float":
		return "float64"
	default:
		return s
	}
}
