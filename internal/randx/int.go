package randx

import (
	"crypto/rand"
	"math"
	"math/big"
)

// NonNegativeInt returns a non-negative random int from a crypto-safe source.
func NonNegativeInt() int {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
