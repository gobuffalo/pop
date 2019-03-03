package generate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Attribute_String(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		exp  string
		name string
	}{
		{
			name: "id",
			exp:  "\tID string `json:\"id\" db:\"id\"`",
		},
		{
			name: "user_id",
			exp:  "\tUserID string `json:\"user_id\" db:\"user_id\"`",
		},
		{
			name: "UserID",
			exp:  "\tUserID string `json:\"user_id\" db:\"user_id\"`",
		},
		{
			name: "userid",
			exp:  "\tUserid string `json:\"userid\" db:\"userid\"`",
		},
		{
			name: "userId",
			exp:  "\tUserID string `json:\"user_id\" db:\"user_id\"`",
		},
		{
			name: "expires",
			exp:  "\tExpires string `json:\"expires\" db:\"expires\"`",
		},
		{
			name: "message_headers",
			exp:  "\tMessageHeaders string `json:\"message_headers\" db:\"message_headers\"`",
		},
	}

	for _, c := range cases {
		model, err := newModel("car", "json", "models")
		r.NoError(err)
		a, err := newAttribute(c.name, &model)
		r.NoError(err)
		r.Equal(c.exp, a.String())
	}
}

func Test_newAttribute(t *testing.T) {
	cases := []struct {
		AttributeInput string
		ResultType     string
		Nullable       bool

		ModelHasUUID   bool
		ModelHasNulls  bool
		ModelHasSlices bool
		Invalid        bool
	}{
		{
			AttributeInput: "name",
			ResultType:     "string",
		},

		{
			AttributeInput: "name:text",
			ResultType:     "string",
		},
		{
			AttributeInput: "id:uuid.UUID",
			ResultType:     "uuid.UUID",
		},
		{
			AttributeInput: "other:uuid",
			ResultType:     "uuid.UUID",
			ModelHasUUID:   true,
		},
		{
			AttributeInput: "optional:nulls.String",
			ResultType:     "nulls.String",
			ModelHasNulls:  true,
			Nullable:       true,
		},
		{
			AttributeInput: "optional:slices.float",
			ResultType:     "slices.Float",
			ModelHasSlices: true,
		},
		{
			AttributeInput: "raw:blob",
			ResultType:     "[]byte",
		},
		{
			AttributeInput: "raw:[]byte",
			ResultType:     "[]byte",
		},
		{
			AttributeInput: "age:int",
			ResultType:     "int",
		},
		{
			AttributeInput: "age:int:int64",
			ResultType:     "int64",
		},
		{
			AttributeInput: "111:int",
			Invalid:        true,
		},
		{
			AttributeInput: "admin/user",
			Invalid:        true,
		},
		{
			AttributeInput: "admin;user",
			Invalid:        true,
		},
		{
			AttributeInput: "_bread",
			Invalid:        true,
		},
	}

	for index, tcase := range cases {
		t.Run(fmt.Sprintf("%d-%s", index, tcase.AttributeInput), func(tt *testing.T) {
			r := require.New(tt)
			model, err := newModel("car", "json", "models")
			r.NoError(err)
			a, err := newAttribute(tcase.AttributeInput, &model)
			if tcase.Invalid {
				r.Errorf(err, "%s should be an invalid attribute", tcase.AttributeInput)
				return
			}
			r.NoError(err)

			r.Equal(a.GoType, tcase.ResultType)
			r.Equal(a.Nullable, tcase.Nullable)

			r.Equal(model.HasUUID, tcase.ModelHasUUID)
			r.Equal(model.HasNulls, tcase.ModelHasNulls)
			r.Equal(model.HasSlices, tcase.ModelHasSlices)
		})
	}

}
