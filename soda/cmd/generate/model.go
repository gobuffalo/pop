package generate

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/flect"
	nflect "github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/pop"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

type model struct {
	Package               string
	ModelPath             string
	Imports               []string
	Name                  nflect.Ident
	attributesCache       map[string]struct{}
	Attributes            []attribute
	ValidatableAttributes []attribute
	StructTag             string

	HasNulls  bool
	HasUUID   bool
	HasSlices bool
	HasID     bool
}

func (m model) Generate() error {
	g := makr.New()
	defer g.Fmt(".")
	ctx := makr.Data{}
	ctx["model"] = m
	ctx["plural_model_name"] = m.modelNamePlural()
	ctx["model_name"] = m.modelName()
	ctx["package_name"] = m.Package

	ctx["test_package_name"] = m.testPkgName()

	ctx["char"] = m.Name.Char()
	ctx["encoding_type"] = m.StructTag
	ctx["encoding_type_char"] = nflect.Char(m.StructTag)

	fname := filepath.Join(m.ModelPath, m.Name.File(".go").String())
	g.Add(makr.NewFile(fname, modelTemplate))
	tfname := filepath.Join(m.ModelPath, m.Name.File("_test.go").String())
	g.Add(makr.NewFile(tfname, modelTestTemplate))
	return g.Run(".", ctx)
}

func (m model) modelName() string {
	x := strings.Split(m.Name.String(), "/")
	for i, s := range x {
		x[i] = flect.New(s).Singularize().Pascalize().String()
	}
	return strings.Join(x, "")
}

func (m model) modelNamePlural() string {
	return flect.New(m.modelName()).Pluralize().Pascalize().String()
}

func (m model) testPkgName() string {
	pkg := m.Package

	path, _ := os.Getwd()
	path = filepath.Join(path, m.ModelPath)

	if _, err := os.Stat(path); err != nil {
		return pkg
	}
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, "_test.go") {
			fset := token.NewFileSet()

			b, err := ioutil.ReadFile(p)
			if err != nil {
				return errors.WithStack(err)
			}
			f, err := parser.ParseFile(fset, p, string(b), 0)
			if err != nil {
				return errors.WithStack(err)
			}

			conf := types.Config{Importer: importer.Default()}
			p, err := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
			if err != nil {
				return errors.WithStack(err)
			}
			pkg = p.Name()

			return io.EOF
		}
		return nil
	})

	return pkg
}

func (m *model) addAttribute(a attribute) error {
	k := a.Name.String()
	if _, found := m.attributesCache[k]; found {
		return fmt.Errorf("duplicated field \"%s\"", k)
	}
	m.attributesCache[k] = struct{}{}
	if a.Name.String() == "id" {
		// No need to create a default ID
		m.HasID = true
		a.Primary = true
		// Ensure ID is the first attribute
		m.Attributes = append([]attribute{a}, m.Attributes...)
	} else {
		m.Attributes = append(m.Attributes, a)
	}

	if a.Nullable {
		return nil
	}

	if a.IsValidable() {
		if a.GoType == "time.Time" {
			a.GoType = "Time"
		}
		m.ValidatableAttributes = append(m.ValidatableAttributes, a)
	}
	return nil
}

func (m *model) addID() {
	if m.HasID {
		return
	}

	if !m.HasUUID {
		m.HasUUID = true
		m.Imports = append(m.Imports, "github.com/gofrs/uuid")
	}

	id := flect.New("id")
	a := attribute{Name: id, OriginalType: "uuid.UUID", GoType: "uuid.UUID", Primary: true}
	// Ensure ID is the first attribute
	m.Attributes = append([]attribute{a}, m.Attributes...)
	m.HasID = true
}

func (m model) generateModelFile() error {
	err := os.MkdirAll(m.ModelPath, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create folder %s", m.ModelPath)
	}

	return m.Generate()
}

func (m model) generateFizz(path string) error {
	migrationPath := defaults.String(path, "./migrations")
	return pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s", m.Name.Tableize()), "fizz", []byte(m.Fizz()), []byte(m.UnFizz()))
}

func (m model) generateSQL(path, env string) error {
	migrationPath := defaults.String(path, "./migrations")
	db, err := pop.Connect(env)
	if err != nil {
		return err
	}

	d := db.Dialect
	f := d.FizzTranslator()

	return pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s.%s", m.Name.Tableize(), d.Name()), "sql", []byte(m.GenerateSQLFromFizz(m.Fizz(), f)), []byte(m.GenerateSQLFromFizz(m.UnFizz(), f)))
}

// Fizz generates the create table instructions
func (m model) Fizz() string {
	s := []string{fmt.Sprintf("create_table(\"%s\") {", m.Name.Tableize())}
	for _, a := range m.Attributes {
		switch a.Name.String() {
		case "created_at", "updated_at":
		default:
			col := fizz.Column{
				Name:    a.Name.Underscore().String(),
				ColType: fizzColType(a.OriginalType),
				Options: map[string]interface{}{},
			}
			if a.Primary {
				col.Options["primary"] = true
			}
			if a.Nullable {
				col.Options["null"] = true
			}
			s = append(s, "\t"+col.String())
		}
	}
	s = append(s, "}")
	return strings.Join(s, "\n")
}

// UnFizz generates the drop table instructions
func (m model) UnFizz() string {
	return fmt.Sprintf("drop_table(\"%s\")", m.Name.Tableize())
}

// GenerateSQLFromFizz generates SQL instructions from fizz instructions
func (m model) GenerateSQLFromFizz(content string, f fizz.Translator) string {
	content, err := fizz.AString(content, f)
	if err != nil {
		return ""
	}
	return content
}

func newModel(name, structTag, modelPath string) (model, error) {
	m := model{
		Package:               filepath.Base(modelPath),
		ModelPath:             modelPath,
		Imports:               []string{"time", "github.com/gobuffalo/pop", "github.com/gobuffalo/validate"},
		Name:                  nflect.New(name),
		Attributes:            []attribute{},
		ValidatableAttributes: []attribute{},
		attributesCache:       map[string]struct{}{},
		StructTag:             structTag,
	}

	switch structTag {
	case "json":
		m.Imports = append(m.Imports, "encoding/json")
	case "xml":
		m.Imports = append(m.Imports, "encoding/xml")
	default:
		return model{}, errors.New("invalid struct tags (use xml or json)")
	}

	_ = m.addAttribute(attribute{Name: flect.New("created_at"), OriginalType: "time.Time", GoType: "time.Time", PreventValidation: true, StructTag: structTag})
	_ = m.addAttribute(attribute{Name: flect.New("updated_at"), OriginalType: "time.Time", GoType: "time.Time", PreventValidation: true, StructTag: structTag})

	return m, nil
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
