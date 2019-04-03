package fix

import (
	"io/ioutil"
	"testing"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func Test_AutoTimestampsOff(t *testing.T) {
	r := require.New(t)
	box := packr.New("./fixtures/auto_timestamps_off", "./fixtures/auto_timestamps_off")

	err := box.Walk(func(path string, info packr.File) error {
		t.Run(path, func(tt *testing.T) {
			rr := require.New(tt)
			b, err := ioutil.ReadAll(info)
			rr.NoError(err)

			body := string(b)
			_, err = AutoTimestampsOff(body)
			rr.NoError(err)
		})
		return nil
	})
	r.NoError(err)
}
