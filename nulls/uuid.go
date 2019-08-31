package nulls

import (
	"github.com/gobuffalo/nulls"
)

// UUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
//
// Deprecated: use github.com/gobuffalo/nulls#UUID instead.
type UUID = nulls.UUID

// NewUUID returns a new, properly instantiated
// UUID object.
//
// Deprecated: use github.com/gobuffalo/nulls#NewUUID instead.
var NewUUID = nulls.NewUUID
