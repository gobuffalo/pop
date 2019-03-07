package nulls

import (
	"github.com/gobuffalo/nulls"
)

// UInt32 adds an implementation for int
// that supports proper JSON encoding/decoding.
type UInt32 nulls.UInt32

// NewUInt32 returns a new, properly instantiated
// Int object.
var NewUInt32 = nulls.NewUInt32
