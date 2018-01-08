package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_UUIDIsPresent(t *testing.T) {
	r := require.New(t)

	id := uuid.NewV4()
	v := UUIDIsPresent{"Name", id}
	errors := validate.NewErrors()
	v.IsValid(errors)
	r.Equal(errors.Count(), 0)

	v = UUIDIsPresent{"Name", uuid.UUID{}}
	v.IsValid(errors)
	r.Equal(errors.Count(), 1)
	r.Equal(errors.Get("name"), []string{"Name can not be blank."})
}
