package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int64 replaces sql.Int64 with an implementation
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Int64 instead.
type Int64 = nulls.Int64

// NewInt64 returns a new, properly instantiated
// Int64 object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewInt64 instead.
var NewInt64 = nulls.NewInt64
