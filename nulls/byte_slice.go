package nulls

import (
	"github.com/gobuffalo/nulls"
)

// ByteSlice adds an implementation for []byte
// that supports proper JSON encoding/decoding.
type ByteSlice nulls.ByteSlice

// NewByteSlice returns a new, properly instantiated
// ByteSlice object.
var NewByteSlice = nulls.NewByteSlice
