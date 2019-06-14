package randx

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(1)
}

func Test_String(t *testing.T) {
	r := require.New(t)
	r.Len(String(5), 5)
	r.Len(String(50), 50)
}
