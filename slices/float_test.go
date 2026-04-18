package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Float_Scan(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		in := "{}"
		v := &Float{}
		require.NoError(t, v.Scan(in))
		require.Len(t, *v, 0)
	})

	t.Run("non-empty slice", func(t *testing.T) {
		in := "{3.14,9.999}"
		v := &Float{}
		require.NoError(t, v.Scan(in))
		require.Equal(t, []float64(*v), []float64{3.14, 9.999})
	})

	t.Run("invalid entry", func(t *testing.T) {
		in := "{44,word}"
		v := &Float{}
		require.Error(t, v.Scan(in))
	})
}
