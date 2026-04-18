// Package randx provides functions for generating random data.
package randx

import (
	"math/rand/v2"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterBytesLen = len(letterBytes)
)

// String generates a random string with the given length.
func String(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.IntN(letterBytesLen)]
	}

	return string(b)
}
