package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Float64 replaces sql.NullFloat64 with an implementation
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Float64 instead.
type Float64 = nulls.Float64

// NewFloat64 returns a new, properly instantiated
// Float64 object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewFloat64 instead.
var NewFloat64 = nulls.NewFloat64
