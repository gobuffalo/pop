package pop

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_genericDumpSchema(t *testing.T) {
	table := []struct {
		cmd *exec.Cmd
		err bool
	}{
		{exec.Command("ls"), false},
		{exec.Command("asdfasdf"), true},
	}

	for _, tt := range table {
		t.Run(strings.Join(tt.cmd.Args, " "), func(st *testing.T) {
			r := require.New(st)
			bb := &bytes.Buffer{}
			err := genericDumpSchema(&ConnectionDetails{}, tt.cmd, bb)
			if tt.err {
				r.Error(err)
				return
			}
			r.NoError(err)
			r.NotEqual(0, len(bb.Bytes()))
		})
	}
}
