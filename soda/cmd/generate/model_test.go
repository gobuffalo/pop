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
	r := require.New(t)

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
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			m := newModel("car")
			a := newAttribute(tcase.AttrInput, &m)
			m.addAttribute(a)

			r.Equal(m.HasID, tcase.HasID)
			r.Equal(m.HasNulls, tcase.HasNulls)

			if !tcase.Validable {
				log.Println(m.ValidatableAttributes)
				r.Equal(len(m.ValidatableAttributes), 0)
				return
			}

			r.Equal(m.ValidatableAttributes[0].Name, a.Name)
		})

	}

}

func Test_model_addID(t *testing.T) {
	r := require.New(t)

	m := newModel("car")
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, true)
	r.Equal(string(m.Attributes[0].Name), "id")
	r.Equal(string(m.Attributes[0].GoType), "uuid.UUID")

	m = newModel("car")
	m.addAttribute(newAttribute("id:int", &m))
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, false)
	r.Equal(string(m.Attributes[0].Name), "id")
	r.Equal(string(m.Attributes[0].GoType), "int")
}

func Test_testPkgName(t *testing.T) {
	r := require.New(t)
	m := newModel("car")

	r.Equal("models", m.testPkgName())

	os.Mkdir("./models", 0755)
	defer os.RemoveAll("./models")

	r.Equal("models", m.testPkgName())

	f, err := os.Create(filepath.Join("models", "foo_test.go"))
	r.NoError(err)
	_, err = f.Write([]byte("// some comment\npackage models"))
	f.Close()

	r.Equal("models", m.testPkgName())

	r.NoError(os.Remove(f.Name()))
	f, err = os.Create(filepath.Join("models", "foo_test.go"))
	r.NoError(err)
	_, err = f.Write([]byte("// some comment\npackage models_test"))
	f.Close()

	r.Equal("models_test", m.testPkgName())
}
