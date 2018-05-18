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

	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/fizz"
	"github.com/markbates/inflect"
)

type model struct {
	Package               string
	Imports               []string
	Name                  inflect.Name
	Attributes            []attribute
	ValidatableAttributes []attribute

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
	ctx["plural_model_name"] = m.Name.ModelPlural()
	ctx["model_name"] = m.Name.Model()
	ctx["package_name"] = m.Package

	ctx["test_package_name"] = m.testPkgName()

	ctx["char"] = strings.ToLower(string([]byte(m.Name)[0]))
	ctx["encoding_type"] = structTag
	ctx["encoding_type_char"] = strings.ToLower(string([]byte(structTag)[0]))

	fname := filepath.Join(m.Package, m.Name.File()+".go")
	g.Add(makr.NewFile(fname, modelTemplate))
	tfname := filepath.Join(m.Package, m.Name.File()+"_test.go")
	g.Add(makr.NewFile(tfname, modelTestTemplate))
	return g.Run(".", ctx)
}

func (m model) testPkgName() string {
	pkg := m.Package

	path, _ := os.Getwd()
	path = filepath.Join("models")

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

func (m *model) addID() {
	if m.HasID {
		return
	}

	if !m.HasUUID {
		m.HasUUID = true
		m.Imports = append(m.Imports, "github.com/gobuffalo/uuid")
	}

	id := inflect.Name("id")
	a := attribute{Name: id, OriginalType: "uuid.UUID", GoType: "uuid.UUID"}
	// Ensure ID is the first attribute
	m.Attributes = append([]attribute{a}, m.Attributes...)
	m.HasID = true
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

func (m model) generateFizz(migrationPath string) error {
	err := pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s", m.Name.Table()), "fizz", []byte(m.Fizz()), []byte(fmt.Sprintf("drop_table(\"%s\")", m.Name.Table())))
	if err != nil {
		return err
	}

	return nil
}

func (m model) generateSQL(migrationPath, env string) error {
	db, err := pop.Connect(env)
	if err != nil {
		return err
	}

	err = pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s.%s", m.Name.Table(), db.Dialect.Name()), "sql", []byte(m.GenerateSQLFromFizz(m.Fizz(), db)), []byte(m.GenerateSQLFromFizz(m.UnFizz(), db)))
	if err != nil {
		return err
	}

	return nil
}

// Fizz generates the create table instructions
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
func newModel(name string) (model, error) {
	var (
		encodingImport string
		m              model
	)

	switch structTag {
	case "json":
		encodingImport = "encoding/json"
	case "xml":
		encodingImport = "encoding/xml"
	default:
		return m, errors.New("Invalid struct tags (use xml or json)")
	}

	m = model{
		Package: "models",
		Imports: []string{"time", "github.com/gobuffalo/pop", "github.com/gobuffalo/validate", encodingImport},
		Name:    inflect.Name(name),
		Attributes: []attribute{
			{Name: inflect.Name("created_at"), OriginalType: "time.Time", GoType: "time.Time"},
			{Name: inflect.Name("updated_at"), OriginalType: "time.Time", GoType: "time.Time"},
		},
		ValidatableAttributes: []attribute{},
	}
	return m, nil
}

func newModelFromArgs(args []string) (model, error) {
	var m model

	if len(args) == 0 {
		return m, errors.New("You must supply a name for your model")
	}

	m, err := newModel(args[0])

	if err != nil {
		return m, err
	}

	for _, def := range args[1:] {
		newAttribute(def, &m)
	}

	// Add a default UUID, if no custom ID is provided
	m.addID()

	return m, nil
}

// UnFizz generates the drop table instructions
func (m model) UnFizz() string {
	return fmt.Sprintf("drop_table(\"%s\")", m.Name.Table())
}

// GenerateSQLFromFizz generates SQL instructions from fizz instructions
func (m model) GenerateSQLFromFizz(content string, c *pop.Connection) string {
	content, err := fizz.AString(content, c.Dialect.FizzTranslator())
	if err != nil {
		return ""
	}
	return content
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
		if nrx.MatchString(s) {
			return fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
