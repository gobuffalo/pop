package generate

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConfigCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := ConfigCmd
	c.SetArgs([]string{})

	tdir, err := ioutil.TempDir("", "testapp")
	r.NoError(err)
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	err = c.Execute()
	r.NoError(err)
	r.FileExists("./database.yml")
}
