package nulls

import (
	"github.com/gobuffalo/nulls"
)

// Int adds an implementation for int
// that supports proper JSON encoding/decoding.
type Int nulls.Int

// NewInt returns a new, properly instantiated
// Int object.
var NewInt = nulls.NewInt
