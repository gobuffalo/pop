package generate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConfigCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := ConfigCmd
	c.SetArgs([]string{})

	tdir := t.TempDir()
	t.Chdir(tdir)

	err := c.Execute()
	r.NoError(err)
	r.FileExists("./database.yml")
}
