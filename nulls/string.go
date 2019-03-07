package nulls

import (
	"github.com/gobuffalo/nulls"
)

// String replaces sql.NullString with an implementation
// that supports proper JSON encoding/decoding.
type String nulls.String

// NewString returns a new, properly instantiated
// String object.
var NewString = nulls.NewString
