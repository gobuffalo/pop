package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int adds an implementation for int
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Int instead.
type Int = nulls.Int

// NewInt returns a new, properly instantiated
// Int object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewInt instead.
var NewInt = nulls.NewInt
