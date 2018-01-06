package slices

import (
	"encoding/json"
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_UUID_JSON_Unmarshal(t *testing.T) {
	r := require.New(t)

	a, err := uuid.NewV4()
	r.NoError(err)
	b, err := uuid.NewV4()
	r.NoError(err)
	x := fmt.Sprintf("[%q, %q]", a, b)
	s := UUID{}
	r.NoError(json.Unmarshal([]byte(x), &s))
	r.Equal(UUID{a, b}, s)
}

func Test_UUID_JSON_Marshal(t *testing.T) {
	r := require.New(t)

	a, err := uuid.NewV4()
	r.NoError(err)
	b, err := uuid.NewV4()
	r.NoError(err)
	x := fmt.Sprintf("[%q,%q]", a, b)

	s := UUID{a, b}
	bb, err := json.Marshal(s)
	r.NoError(err)
	r.Equal(x, string(bb))
}
