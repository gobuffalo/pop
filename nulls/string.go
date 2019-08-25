package nulls

import (
	"github.com/gobuffalo/nulls"
)

// String replaces sql.NullString with an implementation
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#String instead.
type String = nulls.String

// NewString returns a new, properly instantiated
// String object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewString instead.
var NewString = nulls.NewString
