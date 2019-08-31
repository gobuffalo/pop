package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Time replaces sql.NullTime with an implementation
// that supports proper JSON encoding/decoding.
//
// Deprecated: use github.com/gobuffalo/nulls#Time instead.
type Time = nulls.Time

// NewTime returns a new, properly instantiated
// Time object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewTime instead.
var NewTime = nulls.NewTime
