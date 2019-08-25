package nulls_test

import (
	"encoding/json"
	"testing"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func Test_UUID_UnmarshalJSON(t *testing.T) {
	r := require.New(t)
	id, err := uuid.NewV4()
	r.NoError(err)

	b, err := json.Marshal(id)
	r.NoError(err)

	nid := &nulls.UUID{}

	r.NoError(nid.UnmarshalJSON(b))

	r.True(nid.Valid)
	r.Equal(id.String(), nid.UUID.String())
}
