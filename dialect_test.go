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

func Test_genericDumpSchema_stripsBackslashCommands(t *testing.T) {
	r := require.New(t)

	// Use printf to simulate pg_dump output containing backslash commands
	cmd := exec.Command("printf", "\\\\restrict abc123\\nCREATE TABLE t (id int);\\n\\\\unrestrict abc123\\n")
	bb := &bytes.Buffer{}
	err := genericDumpSchema(&ConnectionDetails{}, cmd, bb)
	r.NoError(err)
	r.Equal("CREATE TABLE t (id int);\n", bb.String())
}

func Test_stripPsqlBackslashCommands(t *testing.T) {
	table := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "no backslash commands",
			input:  "CREATE TABLE foo (id int);\nINSERT INTO foo VALUES (1);\n",
			output: "CREATE TABLE foo (id int);\nINSERT INTO foo VALUES (1);\n",
		},
		{
			name:   "restrict and unrestrict",
			input:  "-- header\n\\restrict abc123\nSET statement_timeout = 0;\n\\unrestrict abc123\n",
			output: "-- header\nSET statement_timeout = 0;\n",
		},
		{
			name:   "various backslash commands",
			input:  "\\connect mydb\nCREATE TABLE t (id int);\n\\echo done\n",
			output: "CREATE TABLE t (id int);\n",
		},
		{
			name:   "empty input",
			input:  "",
			output: "",
		},
		{
			name:   "only backslash commands",
			input:  "\\restrict key1\n\\unrestrict key1\n",
			output: "",
		},
		{
			name:   "backslash in middle of line is kept",
			input:  "SELECT E'hello\\nworld';\n",
			output: "SELECT E'hello\\nworld';\n",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(st *testing.T) {
			r := require.New(st)
			result := stripPsqlBackslashCommands([]byte(tt.input))
			r.Equal(tt.output, string(result))
		})
	}
}
