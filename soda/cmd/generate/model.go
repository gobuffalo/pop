package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/markbates/going/defaults"
	"github.com/markbates/inflect"
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var skipMigration bool

func init() {
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")
}

var nrx = regexp.MustCompile(`^nulls.(.+)`)

type names struct {
	Original string
	Table    string
	Proper   string
	File     string
	Plural   string
	Char     string
}

func newName(name string) names {
	return names{
		Original: name,
		File:     name,
		Table:    inflect.Tableize(name),
		Proper:   inflect.Camelize(name),
		Plural:   inflect.Pluralize(inflect.Camelize(name)),
		Char:     strings.ToLower(string([]byte(name)[0])),
	}
}

type attribute struct {
	Names        names
	OriginalType string
	GoType       string
	Nullable     bool
}

func (a attribute) String() string {
	return fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s\"`", a.Names.Proper, a.GoType, a.Names.Original, a.Names.Original)
}

type model struct {
	Package    string
	Imports    []string
	Names      names
	Attributes []attribute
}

func (m model) String() string {
	tmp := strings.Replace(modelTemplate, "PLURAL_MODEL_NAME", m.Names.Plural, -1)
	tmp = strings.Replace(tmp, "MODEL_NAME", m.Names.Proper, -1)
	tmp = strings.Replace(tmp, "PACKAGE_NAME", m.Package, -1)
	tmp = strings.Replace(tmp, "CHAR", m.Names.Char, -1)
	ims := []string{}
	for _, im := range m.Imports {
		ims = append(ims, fmt.Sprintf("\t\"%s\"", im))
	}
	tmp = strings.Replace(tmp, "IMPORTS", strings.Join(ims, "\n"), -1)
	ats := []string{}
	for _, a := range m.Attributes {
		ats = append(ats, a.String())
	}
	tmp = strings.Replace(tmp, "ATTRIBUTES", strings.Join(ats, "\n"), -1)

	return tmp
}

func (m model) Fizz() string {
	s := []string{fmt.Sprintf("create_table(\"%s\", func(t) {", m.Names.Table)}
	for _, a := range m.Attributes {
		switch a.Names.Original {
		case "id", "created_at", "updated_at":
		default:
			x := fmt.Sprintf("\tt.Column(\"%s\", \"%s\", {})", a.Names.Original, fizzColType(a.OriginalType))
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
	id := newName("id")
	id.Proper = "ID"
	return model{
		Package: "models",
		Imports: []string{"time", "encoding/json"},
		Names:   newName(name),
		Attributes: []attribute{
			{Names: id, OriginalType: "int", GoType: "int"},
			{Names: newName("created_at"), OriginalType: "time.Time", GoType: "time.Time"},
			{Names: newName("updated_at"), OriginalType: "time.Time", GoType: "time.Time"},
		},
	}
}

var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your model!")
		}

		model := newModel(args[0])

		hasNulls := false
		for _, def := range args[1:] {
			col := strings.Split(def, ":")
			if len(col) == 1 {
				col = append(col, "string")
			}
			nullable := nrx.MatchString(col[1])
			if !hasNulls && nullable {
				hasNulls = true
				model.Imports = append(model.Imports, "github.com/markbates/pop/nulls")
			}
			model.Attributes = append(model.Attributes, attribute{
				Names:        newName(col[0]),
				OriginalType: col[1],
				GoType:       colType(col[1]),
				Nullable:     nullable,
			})
		}

		err := os.MkdirAll(model.Package, 0766)
		if err != nil {
			return errors.Wrapf(err, "couldn't create folder %s", model.Package)
		}

		fname := filepath.Join(model.Package, model.Names.File+".go")
		err = ioutil.WriteFile(fname, []byte(model.String()), 0666)
		if err != nil {
			return errors.Wrapf(err, "couldn't write to file %s", fname)
		}
		fmt.Printf("> %s\n", fname)

		tfname := filepath.Join(model.Package, model.Names.File+"_test.go")
		tmp := strings.Replace(modelTestTemplate, "MODEL_NAME", model.Names.Proper, -1)
		tmp = strings.Replace(tmp, "PACKAGE_NAME", model.Package, -1)
		err = ioutil.WriteFile(tfname, []byte(tmp), 0666)
		if err != nil {
			return errors.Wrapf(err, "couldn't write to file %s", tfname)
		}
		fmt.Printf("> %s\n", tfname)

		md, _ := filepath.Abs(fname)
		goi := exec.Command("gofmt", "-w", md)
		out, err := goi.CombinedOutput()
		if err != nil {
			fmt.Printf("Received an error when trying to run gofmt -> %#v\n", err)
			fmt.Println(out)
		}

		if !skipMigration {
			cflag := cmd.Flag("path")
			migrationPath := defaults.String(cflag.Value.String(), "./migrations")
			err = pop.MigrationCreate(migrationPath, fmt.Sprintf("create_%s", model.Names.Table), "fizz", []byte(model.Fizz()), []byte(fmt.Sprintf("drop_table(\"%s\")", model.Names.Table)))
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func colType(s string) string {
	switch s {
	case "text":
		return "string"
	case "time", "timestamp":
		return "time.Time"
	case "nulls.Text":
		return "nulls.String"
	default:
		return s
	}
}

func fizzColType(s string) string {
	if nrx.MatchString(s) {
		return fizzColType(strings.Replace(s, "nulls.", "", -1))
	}
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time":
		return "timestamp"
	default:
		return strings.ToLower(s)
	}
}
