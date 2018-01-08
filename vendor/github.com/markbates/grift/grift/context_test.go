package grift

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewContext(t *testing.T) {
	r := require.New(t)
	c := NewContext("foo")
	r.Equal("foo", c.Name)
	c.Set("bar", "baz")
	r.Equal("baz", c.Value("bar"))
	r.Nil(c.Value("bax"))
}

func Test_NewContextWithContext(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "bax", "xab")
	c := NewContextWithContext("foo", ctx)
	r.Equal("foo", c.Name)
	c.Set("bar", "baz")
	r.Equal("baz", c.Value("bar"))
	r.Equal("xab", c.Value("bax"))
}
