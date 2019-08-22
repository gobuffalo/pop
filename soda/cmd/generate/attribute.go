package generate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobuffalo/flect"
)

var attrNamePattern = regexp.MustCompile(`^\p{L}[\p{L}\d_]*$`)

type attribute struct {
	Name              flect.Ident
	OriginalType      string
	GoType            string
	Nullable          bool
	Primary           bool
	PreventValidation bool
	StructTag         string
}

func (a attribute) String() string {
	if len(a.StructTag) == 0 {
		a.StructTag = "json"
	}
	return fmt.Sprintf("\t%s %s `%s:\"%s\" db:\"%s\"`", a.Name.Pascalize(), a.GoType, a.StructTag, a.Name.Underscore(), a.Name.Underscore())
}

func (a attribute) IsValidable() bool {
	return !a.PreventValidation && (a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int")
}

func newAttribute(base string, model *model) (attribute, error) {
	col := strings.Split(base, ":")
	if len(col) == 1 {
		col = append(col, "string")
	}

	if !attrNamePattern.MatchString(col[0]) {
		return attribute{}, fmt.Errorf("%s is not a valid attribute name", col[0])
	}

	nullable := strings.HasPrefix(col[1], "nulls.")
	if !model.HasNulls && nullable {
		model.HasNulls = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/nulls")
	} else if !model.HasSlices && strings.HasPrefix(col[1], "slices.") {
		model.HasSlices = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/pop/slices")
	} else if !model.HasUUID && col[1] == "uuid" {
		model.HasUUID = true
		model.Imports = append(model.Imports, "github.com/gofrs/uuid")
	}

	got := colType(col[1])
	if len(col) > 2 {
		got = col[2]
	}
	name := flect.New(col[0])
	primary := false
	if name.String() == "id" {
		primary = true
	}
	a := attribute{
		Name:         name,
		OriginalType: col[1],
		GoType:       got,
		Nullable:     nullable,
		StructTag:    model.StructTag,
		Primary:      primary,
	}

	return a, nil
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
	case "[]byte", "blob":
		return "[]byte"
	default:
		return s
	}
}
