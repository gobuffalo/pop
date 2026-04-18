package generate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FizzCmd_NoArg(t *testing.T) {
	r := require.New(t)
	c := FizzCmd
	c.SetArgs([]string{})

	tdir := t.TempDir()
	t.Chdir(tdir)

	err := c.Execute()
	r.EqualError(err, "you must set a name for your migration")
}
