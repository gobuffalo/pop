package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Int_Scan(t *testing.T) {
	r := require.New(t)
	t.Run("empty slice", func(t *testing.T) {
		in := "{}"
		v := &Int{}
		r.NoError(v.Scan(in))
		r.Len(*v, 0)
	})

	t.Run("non-empty slice", func(t *testing.T) {
		in := "{44,55}"
		v := &Int{}
		r.NoError(v.Scan(in))
		r.Equal([]int(*v), []int{44, 55})
	})

	t.Run("invalid entry", func(t *testing.T) {
		in := "{44,word}"
		v := &Int{}
		r.Error(v.Scan(in))
	})
}
