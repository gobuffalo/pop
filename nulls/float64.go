package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Float64 replaces sql.NullFloat64 with an implementation
// that supports proper JSON encoding/decoding.
type Float64 nulls.Float64

// NewFloat64 returns a new, properly instantiated
// Float64 object.
var NewFloat64 = nulls.NewFloat64
