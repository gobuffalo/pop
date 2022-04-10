package generate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FizzCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := FizzCmd
	c.SetArgs([]string{})

	tdir := t.TempDir()

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	err = c.Execute()
	r.EqualError(err, "you must set a name for your migration")
}
