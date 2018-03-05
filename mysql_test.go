package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MySQL_URL(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "mysql://dbase:dbase@(dbase:dbase)/dbase?dbase=dbase",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{ConnectionDetails: cd}
	r.Equal("dbase:dbase@(dbase:dbase)/dbase?dbase=dbase", m.URL())
	r.Equal("dbase:dbase@(dbase:dbase)/?dbase=dbase", m.urlWithoutDb())
}
