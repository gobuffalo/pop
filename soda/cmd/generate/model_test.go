package generate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_addAttribute(t *testing.T) {
	cases := []struct {
		AttrInput string
		HasID     bool
		HasNulls  bool
		Validable bool
	}{
		{AttrInput: "plate", HasID: false, Validable: true},
		{AttrInput: "id", HasID: true, Validable: true},
		{AttrInput: "id:int", HasID: true, Validable: true},
		{AttrInput: "optional:nulls.String", HasNulls: true},
	}

	for index, tcase := range cases {
		t.Run(fmt.Sprintf("%v", index), func(tt *testing.T) {
			r := require.New(tt)
			m, err := newModel("car", "json", "models")
			r.NoError(err)
			a, err := newAttribute(tcase.AttrInput, &m)
			r.NoError(err)
			err = m.addAttribute(a)
			r.NoError(err)

			r.Equal(tcase.HasID, m.HasID)
			r.Equal(tcase.HasNulls, m.HasNulls)

			if !tcase.Validable {
				log.Println(m.ValidatableAttributes)
				r.Equal(0, len(m.ValidatableAttributes))
				return
			}

			r.Equal(a.Name, m.ValidatableAttributes[0].Name)
		})

	}

}

func Test_model_addID(t *testing.T) {
	r := require.New(t)

	m, err := newModel("car", "json", "models")
	r.NoError(err)
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, true)
	r.Equal(m.Attributes[0].Name.String(), "id")
	r.Equal(string(m.Attributes[0].GoType), "uuid.UUID")

	m, err = newModel("car", "json", "models")
	r.NoError(err)
	a, err := newAttribute("id:int", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.NoError(err)
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, false)
	r.Equal(m.Attributes[0].Name.String(), "id")
	r.Equal(string(m.Attributes[0].GoType), "int")
}

func Test_model_addDuplicate(t *testing.T) {
	r := require.New(t)

	m, err := newModel("car", "json", "models")
	r.NoError(err)
	a, err := newAttribute("color:string", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.NoError(err)

	a, err = newAttribute("color:string", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.Error(err)

	a, err = newAttribute("color:int", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.Error(err)
}

func Test_package_paths(t *testing.T) {
	cases := map[string]struct {
		path            string
		expectedPackage string
	}{
		"default": {
			"models",
			"models",
		},
		"alternate name": {
			"entities",
			"entities",
		},
		"nested folders": {
			"pkg/models/admin",
			"admin",
		},
	}

	for desc, c := range cases {
		t.Run(desc, func(t *testing.T) {
			r := require.New(t)

			inTempDir(r, func() {
				m, err := newModel("car", "json", c.path)
				r.NoError(err)

				r.Equal(c.expectedPackage, m.Package)
				r.Equal(c.path, m.ModelPath)

				err = m.Generate()
				r.NoError(err)

				content, err := ioutil.ReadFile(path.Join(c.path, "car.go"))
				r.NoError(err)

				r.Contains(string(content), "package "+c.expectedPackage)
			})
		})
	}
}

func Test_testPkgName(t *testing.T) {
	r := require.New(t)
	m, err := newModel("car", "json", "models")
	r.NoError(err)

	r.Equal("models", m.testPkgName())

	err = os.Mkdir("./models", 0755)
	r.NoError(err)

	defer os.RemoveAll("./models")

	r.Equal("models", m.testPkgName())

	f, err := os.Create(filepath.Join("models", "foo_test.go"))
	r.NoError(err)
	_, err = f.Write([]byte("// some comment\npackage models"))
	f.Close()
	r.NoError(err)

	r.Equal("models", m.testPkgName())

	r.NoError(os.Remove(f.Name()))
	f, err = os.Create(filepath.Join("models", "foo_test.go"))
	r.NoError(err)
	_, err = f.Write([]byte("// some comment\npackage models_test"))
	f.Close()
	r.NoError(err)

	r.Equal("models_test", m.testPkgName())
}

func Test_model_Fizz(t *testing.T) {
	r := require.New(t)

	m, err := newModel("car", "json", "models")

	a, err := newAttribute("id:int", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.NoError(err)

	a, err = newAttribute("brand:string", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.NoError(err)

	a, err = newAttribute("owner:nulls.String", &m)
	r.NoError(err)
	err = m.addAttribute(a)
	r.NoError(err)

	expected := `create_table("cars") {
	t.Column("id", "integer", {primary: true})
	t.Column("brand", "string", {})
	t.Column("owner", "string", {null: true})
	t.Timestamps()
}`
	r.Equal(expected, m.Fizz())
}

func inTempDir(r *require.Assertions, fn func()) {
	dir, err := ioutil.TempDir("", "tests")
	r.NoError(err)
	defer os.RemoveAll(dir)

	cwd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(cwd)

	os.Chdir(dir)
	fn()
}
