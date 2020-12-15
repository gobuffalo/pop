package slices

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Map_UnmarshalText(t *testing.T) {
	r := require.New(t)

	m := Map{}
	err := m.UnmarshalText([]byte(`{"a":"b"}`))
	r.NoError(err)
	r.Equal("b", m["a"])
}

func Test_Map_MarshalJSON(t *testing.T) {
	r := require.New(t)

	m := Map{"a": "b"}
	b, err := json.Marshal(m)
	r.NoError(err)
	r.Equal([]byte(`{"a":"b"}`), b)
}

func Test_Map_UnMarshalJSON_uninitialized_map_does_not_panic(t *testing.T) {
	r := require.New(t)

	maps := make([]Map, 0)
	r.NotPanics(func() {
		err := json.Unmarshal([]byte(`[{"a": "b"}]`), &maps)
		r.NoError(err)
		r.Len(maps, 1)
	})
}
