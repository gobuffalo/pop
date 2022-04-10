package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ModelCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := ModelCmd
	c.SetArgs([]string{})

	tdir := t.TempDir()

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	err = c.Execute()
	r.EqualError(err, "you must set a name for your model")
}

func Test_ModelCmd_NameOnly(t *testing.T) {
	r := require.New(t)
	c := ModelCmd
	c.SetArgs([]string{"users"})

	tdir := t.TempDir()

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	err = c.Execute()
	r.NoError(err)

	r.DirExists(filepath.Join(tdir, "migrations"))
	r.DirExists(filepath.Join(tdir, "models"))
	r.FileExists(filepath.Join(tdir, "models", "user.go"))
	r.FileExists(filepath.Join(tdir, "models", "user_test.go"))
}
