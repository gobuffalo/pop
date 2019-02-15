package slices

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_String_Scan(t *testing.T) {
	r := require.New(t)
	in := `{"This has a comma,","This has a double quote\"","Also a single'"}`
	s := &String{}
	r.NoError(s.Scan(in))
	ss := []string(*s)
	r.Len(ss, 3)
	r.Equal([]string{"This has a comma,", "This has a double quote\"", "Also a single'"}, ss)
}

func Test_String_Value(t *testing.T) {
	r := require.New(t)
	s := String{"This has a comma,", "This has a double quote\"", "Also a single'"}
	v, err := s.Value()
	r.NoError(err)
	r.Equal(`{"This has a comma,","This has a double quote\"","Also a single'"}`, v)
}

func Test_String_UnmarshalText(t *testing.T) {
	r := require.New(t)
	in := "foo,bar,\"baz,bax\""
	s := &String{}
	r.NoError(s.UnmarshalText([]byte(in)))

	ss := []string(*s)
	r.Len(ss, 3)
	r.Equal([]string{"foo", "bar", "baz,bax"}, ss)
}

func Test_String_JSON_Unmarshal(t *testing.T) {
	r := require.New(t)

	x := `["foo", "bar"]`
	s := String{}
	r.NoError(json.Unmarshal([]byte(x), &s))
	r.Equal(String{"foo", "bar"}, s)
}

func Test_String_JSON_Marshal(t *testing.T) {
	r := require.New(t)

	x := `["foo","bar"]`
	s := String{"foo", "bar"}
	b, err := json.Marshal(s)
	r.NoError(err)
	r.Equal(x, string(b))
}
