package nulls

import (
	"github.com/gobuffalo/nulls"
)

// UUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
type UUID nulls.UUID

// NewUUID returns a new, properly instantiated
// UUID object.
var NewUUID = nulls.NewUUID
