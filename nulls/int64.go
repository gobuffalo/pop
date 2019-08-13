package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int64 replaces sql.Int64 with an implementation
// that supports proper JSON encoding/decoding.
type Int64 nulls.Int64

// NewInt64 returns a new, properly instantiated
// Int64 object.
var NewInt64 = nulls.NewInt64
