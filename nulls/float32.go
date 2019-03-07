package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Float32 adds an implementation for float32
// that supports proper JSON encoding/decoding.
type Float32 nulls.Float32

// NewFloat32 returns a new, properly instantiated
// Float32 object.
var NewFloat32 = nulls.NewFloat32
