package fizz

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Exec(t *testing.T) {
	r := require.New(t)

	b := NewBubbler(nil)
	f := fizzer{b}
	bb := &bytes.Buffer{}
	f.Exec(bb).(func(string))("echo hello")
	r.Equal("hello\n", bb.String())
}
