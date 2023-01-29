package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Float_Scan(t *testing.T) {
	r := require.New(t)
	t.Run("empty slice", func(t *testing.T) {
		in := "{}"
		v := &Float{}
		r.NoError(v.Scan(in))
		r.Len(*v, 0)
	})

	t.Run("non-empty slice", func(t *testing.T) {
		in := "{3.14,9.999}"
		v := &Float{}
		r.NoError(v.Scan(in))
		r.Equal([]float64(*v), []float64{3.14, 9.999})
	})

	t.Run("invalid entry", func(t *testing.T) {
		in := "{44,word}"
		v := &Float{}
		r.Error(v.Scan(in))
	})
}
