package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_RawSQL(t *testing.T) {
	r := require.New(t)

	ddl := `raw("select * from users")`

	ch := fizz.AString(ddl)
	b := <-ch
	s := b.Data.(string)
	r.Equal("select * from users", s)
}
