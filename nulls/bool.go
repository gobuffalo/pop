package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Bool replaces sql.NullBool with an implementation
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Bool instead.
type Bool = nulls.Bool

// NewBool returns a new, properly instantiated
// Bool object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewBool instead.
var NewBool = nulls.NewBool
