package nulls_test

import (
	"testing"

	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Int32UnmarshalJSON(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		Input []byte
		Value int32
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
			Input: []byte{'4', '2', '.', '1', '2'},
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
		i := nulls.Int32{}
		r.NoError(i.UnmarshalJSON(c.Input))
		r.Equal(c.Value, i.Int32)
		r.Equal(c.Valid, i.Valid)
	}
}

func Test_Int32UnmarshalJSON_Errors(t *testing.T) {
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
		i := nulls.Int32{}
		r.Error(i.UnmarshalJSON(c.Input))
		r.Equal(int32(0), i.Int32)
		r.False(i.Valid)
	}
}
