package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_model_addID(t *testing.T) {
	r := require.New(t)

	m, _ := newModel("car")
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, true)
	r.Equal(string(m.Attributes[0].Name), "id")
	r.Equal(string(m.Attributes[0].GoType), "uuid.UUID")

	m, _ = newModel("car")
	newAttribute("id:int", &m)
	m.addID()

	r.Equal(m.HasID, true)
	r.Equal(m.HasUUID, false)
	r.Equal(string(m.Attributes[0].Name), "id")
	r.Equal(string(m.Attributes[0].GoType), "int")
}

func Test_testPkgName(t *testing.T) {
	r := require.New(t)
	m, _ := newModel("car")

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
