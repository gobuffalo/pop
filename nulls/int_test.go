package nulls_test

import (
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/stretchr/testify/require"
)

func Test_IntUnmarshalJSON(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		Input []byte
		Value int
		Valid bool
	}{
		{
			Input: []byte{'0'},
			Value: 0,
			Valid: true,
		},
		{
			Input: []byte{'4', '2'},
			Value: 42,
			Valid: true,
		},
		{
			Input: []byte{'n', 'u', 'l', 'l'},
			Value: 0,
			Valid: false,
		},
	}

	for _, c := range cases {
		i := nulls.Int{}
		r.NoError(i.UnmarshalJSON(c.Input))
		r.Equal(c.Value, i.Int)
		r.Equal(c.Valid, i.Valid)
	}
}

func Test_IntUnmarshalJSON_Errors(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		Input []byte
	}{
		{
			Input: []byte{'a'},
		},
		{
			Input: []byte{},
		},
	}

	for _, c := range cases {
		i := nulls.Int{}
		r.Error(i.UnmarshalJSON(c.Input))
		r.Equal(0, i.Int)
		r.False(i.Valid)
	}
}
