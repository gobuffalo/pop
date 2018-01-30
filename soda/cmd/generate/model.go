package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/markbates/going/defaults"
	"github.com/markbates/inflect"
	"github.com/markbates/pop"
)

type model struct {
	Package               string
	Imports               []string
	Name                  inflect.Name
	Attributes            []attribute
	ValidatableAttributes []attribute

	HasNulls bool
	HasUUID  bool
	HasID    bool
}

func (m model) Generate() error {
	g := makr.New()
	defer g.Fmt(".")
	ctx := makr.Data{}
	ctx["model"] = m
	ctx["plural_model_name"] = m.Name.Plural()
	ctx["model_name"] = m.Name.Model()
	ctx["package_name"] = m.Package
	ctx["char"] = strings.ToLower(string([]byte(m.Name)[0]))
	ctx["encoding_type"] = structTag
	ctx["encoding_type_char"] = strings.ToLower(string([]byte(structTag)[0]))

	fname := filepath.Join(m.Package, m.Name.File()+".go")
	g.Add(makr.NewFile(fname, modelTemplate))
	tfname := filepath.Join(m.Package, m.Name.File()+"_test.go")
	g.Add(makr.NewFile(tfname, modelTestTemplate))
	return g.Run(".", ctx)
}

func (m model) AppendAttribute(a attribute) {
	if a.Name == "id" {
		// No need to create a default ID
		m.HasID = true
		// Ensure ID is the first attribute
		m.Attributes = append([]attribute{a}, m.Attributes...)
	} else {
		m.Attributes = append(m.Attributes, a)
	}

	if a.Nullable {
		return
	}

	if a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int" {
		if a.GoType == "time.Time" {
			a.GoType = "Time"
		}
		m.ValidatableAttributes = append(m.ValidatableAttributes, a)
	}
}

func (m model) addID() {
	if m.HasID {
		return
	}

	if !m.HasUUID {
		m.Imports = append(m.Imports, "github.com/satori/go.uuid")
	}

	id := inflect.Name("id")
	a := attribute{Name: id, OriginalType: "uuid.UUID", GoType: "uuid.UUID"}
	// Ensure ID is the first attribute
	m.Attributes = append([]attribute{a}, m.Attributes...)
}

func (m model) generateModelFile() error {
	err := os.MkdirAll(m.Package, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create folder %s", m.Package)
	}

	err = m.Generate()
	if err != nil {
		return err
	}

	return nil
}

func (m model) generateFizz(cflag *pflag.Flag) error {
	migrationPath := defaults.String(cflag.Value.String(), "./migrations")
	err := pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s", m.Name.Table()), "fizz", []byte(m.Fizz()), []byte(fmt.Sprintf("drop_table(\"%s\")", m.Name.Table())))
	if err != nil {
		return err
	}

	return nil
}

func (m model) Fizz() string {
	s := []string{fmt.Sprintf("create_table(\"%s\", func(t) {", m.Name.Table())}
	for _, a := range m.Attributes {
		switch a.Name {
		case "created_at", "updated_at":
		case "id":
			s = append(s, fmt.Sprintf("\tt.Column(\"id\", \"%s\", {\"primary\": true})", fizzColType(a.OriginalType)))
		default:
			x := fmt.Sprintf("\tt.Column(\"%s\", \"%s\", {})", a.Name.Underscore(), fizzColType(a.OriginalType))
			if a.Nullable {
				x = strings.Replace(x, "{}", `{"null": true}`, -1)
			}
			s = append(s, x)
		}
	}
	s = append(s, "})")
	return strings.Join(s, "\n")
}

func newModel(name string) model {
	m := model{
		Package: "models",
		Imports: []string{"time", "github.com/markbates/pop", "github.com/markbates/validate"},
		Name:    inflect.Name(name),
		Attributes: []attribute{
			{Name: inflect.Name("created_at"), OriginalType: "time.Time", GoType: "time.Time"},
			{Name: inflect.Name("updated_at"), OriginalType: "time.Time", GoType: "time.Time"},
		},
		ValidatableAttributes: []attribute{},
	}
	return m
}
