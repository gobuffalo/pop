package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int32 adds an implementation for int32
// that supports proper JSON encoding/decoding.
type Int32 nulls.Int32

// NewInt32 returns a new, properly instantiated
// Int object.
var NewInt32 = nulls.NewInt32
