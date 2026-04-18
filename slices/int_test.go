package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Int_Scan(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		in := "{}"
		v := &Int{}
		require.NoError(t, v.Scan(in))
		require.Len(t, *v, 0)
	})

	t.Run("non-empty slice", func(t *testing.T) {
		in := "{44,55}"
		v := &Int{}
		require.NoError(t, v.Scan(in))
		require.Equal(t, []int(*v), []int{44, 55})
	})

	t.Run("invalid entry", func(t *testing.T) {
		in := "{44,word}"
		v := &Int{}
		require.Error(t, v.Scan(in))
	})
}
