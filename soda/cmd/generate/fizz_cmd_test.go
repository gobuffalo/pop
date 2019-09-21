package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_FizzCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := FizzCmd
	c.SetArgs([]string{})

	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, "testapp")
	r.NoError(err)
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	err = c.Execute()
	r.EqualError(err, "you must set a name for your table")
}
