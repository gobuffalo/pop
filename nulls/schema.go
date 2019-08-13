package nulls

import (
	"github.com/gobuffalo/nulls"
)

// RegisterWithSchema allows for the nulls package to be used with http://www.gorillatoolkit.org/pkg/schema#Converter
var RegisterWithSchema = nulls.RegisterWithSchema
