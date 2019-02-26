package nulls_test

import (
	"testing"

	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Int64UnmarshalJSON(t *testing.T) {
	r := require.New(t)
	cases := []struct {
		Input []byte
		Value int64
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
		i := nulls.Int64{}
		r.NoError(i.UnmarshalJSON(c.Input))
		r.Equal(c.Value, i.Int64)
		r.Equal(c.Valid, i.Valid)
	}
}

func Test_Int64UnmarshalJSON_Errors(t *testing.T) {
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
		i := nulls.Int64{}
		r.Error(i.UnmarshalJSON(c.Input))
		r.Equal(int64(0), i.Int64)
		r.False(i.Valid)
	}
}
