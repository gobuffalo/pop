package associations

import (
	"database/sql"
	"fmt"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsZeroOfUnderlyingType(t *testing.T) {
	for k, tc := range []struct {
		in   interface{}
		zero bool
	}{
		{in: nil, zero: true},
		{in: 0, zero: true},
		{in: 1, zero: false},
		{in: false, zero: true},
		{in: "", zero: true},
		{in: interface{}(nil), zero: true},
		{in: uuid.NullUUID{}, zero: true},
		{in: uuid.UUID{}, zero: true},
		{in: uuid.NullUUID{Valid: true}, zero: false},
		{in: nulls.Int{}, zero: true},
		{in: nulls.String{}, zero: true},
		{in: nulls.Bool{}, zero: true},
		{in: nulls.Float64{}, zero: true},
		{in: sql.NullString{}, zero: true},
		{in: sql.NullString{Valid: true}, zero: false},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.EqualValues(t, tc.zero, IsZeroOfUnderlyingType(tc.in))
		})
	}
}
