package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Time replaces sql.NullTime with an implementation
// that supports proper JSON encoding/decoding.
type Time nulls.Time

// NewTime returns a new, properly instantiated
// Time object.
var NewTime = nulls.NewTime
