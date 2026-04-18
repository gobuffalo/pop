package randx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	r := require.New(t)
	r.Len(String(5), 5)
	r.Len(String(50), 50)
}

func Test_String_Parallel(t *testing.T) {
	wg := sync.WaitGroup{}
	for range 100 {
		wg.Go(func() {
			id := String(30)
			if len(id) != 30 {
				t.Fail()
			}
		})
	}
	wg.Wait()
}
