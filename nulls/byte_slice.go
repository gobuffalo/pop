package nulls

import (
	"github.com/gobuffalo/nulls"
)

// ByteSlice adds an implementation for []byte
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#ByteSlice instead.
type ByteSlice = nulls.ByteSlice

// NewByteSlice returns a new, properly instantiated
// ByteSlice object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewByteSlice instead.
var NewByteSlice = nulls.NewByteSlice
