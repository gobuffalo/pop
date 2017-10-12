package slices

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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
