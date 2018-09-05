package generate

import (
	"fmt"
	"log"
	"os"
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
			m, err := newModel("car", "json")
			r.NoError(err)
			a, err := newAttribute(tcase.AttrInput, &m)
			r.NoError(err)
			m.addAttribute(a)

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

	m, err := newModel("car", "json")
	r.NoError(err)
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, true)
	r.Equal(m.Attributes[0].Name.String(), "id")
	r.Equal(string(m.Attributes[0].GoType), "uuid.UUID")

	m, err = newModel("car", "json")
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

	m, err := newModel("car", "json")
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

func Test_testPkgName(t *testing.T) {
	r := require.New(t)
	m, err := newModel("car", "json")
	r.NoError(err)

	r.Equal("models", m.testPkgName())

	os.Mkdir("./models", 0755)
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
