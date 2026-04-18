package columns_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gobuffalo/pop/v6/columns"
)

func Test_Column_UpdateString(t *testing.T) {
	r := require.New(t)
	c := columns.Column{Name: "foo"}
	r.Equal(c.UpdateString(), "foo = :foo")
}
