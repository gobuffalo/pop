package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Bool replaces sql.NullBool with an implementation
// that supports proper JSON encoding/decoding.
type Bool nulls.Bool

// NewBool returns a new, properly instantiated
// Bool object.
var NewBool = nulls.NewBool
