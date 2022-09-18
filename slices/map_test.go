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

func Test_Map_UnmarshalJSON(t *testing.T) {
	r := require.New(t)

	m := Map{}
	err := json.Unmarshal([]byte(`{"a":"b"}`), &m)
	r.NoError(err)
	r.Equal("b", m["a"])
}

// for the next test case
type Ship struct {
	Name string
	Crew Map `json:"crew"`
}

func Test_Map_UnmarshalJSON_Nested(t *testing.T) {
	r := require.New(t)

	p := &Ship{}
	err := json.Unmarshal([]byte(`{"crew":{"captain":"Han", "navigator":"Chewie"}}`), p)
	r.NoError(err)
	r.Equal("Han", p.Crew["captain"])
	r.Equal("Chewie", p.Crew["navigator"])
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

func Test_Map_Scan(t *testing.T) {
	r := require.New(t)
	in := []byte(`{"a":"b"}`)
	m := Map{}
	r.NoError(m.Scan(in))
	r.Equal("b", m["a"])
}

func Test_Map_Null_Scan(t *testing.T) {
	r := require.New(t)
	m := Map{}
	r.NoError(m.Scan(nil))
	r.Equal(Map{}, m)
}
