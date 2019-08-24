package ctable

import (
	"strings"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// New creates a generator to make files for  a table based
// on the given options.
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	t := fizz.NewTable(opts.TableName, map[string]interface{}{})
	for _, attr := range opts.Attrs {
		o := fizz.Options{}
		name := attr.Name.Underscore().String()
		colType := fizzColType(attr.CommonType())
		if name == "id" {
			o["primary"] = true
		}
		if strings.HasPrefix(attr.GoType(), "nulls.") {
			o["null"] = true
		}
		if err := t.Column(name, colType, o); err != nil {
			return g, err
		}
	}
	f := genny.NewFileS(opts.Name, t.String())
	g.File(f)
	return g, nil
}

func fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "slices.uuid", "[]string":
		return "varchar[]"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "numeric[]"
	case "slices.int":
		return "int[]"
	case "slices.map":
		return "jsonb"
	case "float32", "float64", "float":
		return "decimal"
	case "blob", "[]byte":
		return "blob"
	default:
		if strings.HasPrefix(s, "nulls.") {
			return fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
