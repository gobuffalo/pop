package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int32 adds an implementation for int32
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Int32 instead.
type Int32 = nulls.Int32

// NewInt32 returns a new, properly instantiated
// Int object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewInt32 instead.
var NewInt32 = nulls.NewInt32
