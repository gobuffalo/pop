package pop

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/stretchr/testify/require"
)

func Test_MySQL_URL_As_Is(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "mysql://user:pass@(host:port)/dbase?opt=value",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{commonDialect{ConnectionDetails: cd}}
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.URL())
	r.Equal("user:pass@(host:port)/?opt=value", m.urlWithoutDb())
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.MigrationURL())
}

func Test_MySQL_URL_Override_withURL(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Database: "xx",
		Host:     "xx",
		Port:     "xx",
		User:     "xx",
		Password: "xx",
		URL:      "mysql://user:pass@(host:port)/dbase?opt=value",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{commonDialect{ConnectionDetails: cd}}
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.URL())
	r.Equal("user:pass@(host:port)/?opt=value", m.urlWithoutDb())
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.MigrationURL())
}

func Test_MySQL_URL_With_Values(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{
		Database: "dbase",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass",
		Options:  map[string]string{"opt": "value"},
	}}}
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.URL())
	r.Equal("user:pass@(host:port)/?opt=value", m.urlWithoutDb())
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.MigrationURL())
}

func Test_MySQL_URL_Without_User(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{
		Password: "pass",
		Database: "dbase",
	}}}
	// finalizerMySQL fills address part in real world.
	// without user, password cannot live alone
	r.Equal("(:)/dbase?", m.URL())
}

func Test_MySQL_URL_Without_Password(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{
		User:     "user",
		Database: "dbase",
	}}}
	// finalizerMySQL fills address part in real world.
	r.Equal("user@(:)/dbase?", m.URL())
}

func Test_MySQL_URL_urlParserMySQL_Standard(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "mysql://user:pass@(host:port)/database?collation=utf8&param2=value2",
	}
	err := urlParserMySQL(cd)
	r.NoError(err)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
	r.Equal("host", cd.Host)
	r.Equal("port", cd.Port)
	r.Equal("database", cd.Database)
	// only collation is stored as options by urlParserMySQL()
	r.Equal("utf8", cd.Options["collation"])
	r.Equal("", cd.Options["param2"])
}

func Test_MySQL_URL_urlParserMySQL_With_Protocol(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "user:pass@tcp(host:port)/dbase",
	}
	err := urlParserMySQL(cd)
	r.NoError(err)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
	r.Equal("host", cd.Host)
	r.Equal("port", cd.Port)
	r.Equal("dbase", cd.Database)
}

func Test_MySQL_URL_urlParserMySQL_Socket(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "unix(/tmp/socket)/dbase",
	}
	err := urlParserMySQL(cd)
	r.NoError(err)
	r.Equal("/tmp/socket", cd.Host)
	r.Equal("socket", cd.Port)

	// additional test without URL
	cd.URL = ""
	m := &mysql{commonDialect{ConnectionDetails: cd}}
	r.True(strings.HasPrefix(m.URL(), "unix(/tmp/socket)/dbase?"))
	r.True(strings.HasPrefix(m.urlWithoutDb(), "unix(/tmp/socket)/?"))
}

func Test_MySQL_URL_urlParserMySQL_Unsupported(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "mysql://user:pass@host:port/dbase?opt=value",
	}
	err := urlParserMySQL(cd)
	r.Error(err)
}

func Test_MySQL_Database_Open_Failure(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{}}}
	err := m.CreateDB()
	r.Error(err)
	err = m.DropDB()
	r.Error(err)
}

func Test_MySQL_FizzTranslator(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{}
	m := &mysql{commonDialect{ConnectionDetails: cd}}
	ft := m.FizzTranslator()
	r.IsType(&translators.MySQL{}, ft)
	r.Implements((*fizz.Translator)(nil), ft)
}

func Test_MySQL_Finalizer_Default_CD(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{}}}
	finalizerMySQL(m.ConnectionDetails)
	r.Equal(hostMySQL, m.ConnectionDetails.Host)
	r.Equal(portMySQL, m.ConnectionDetails.Port)
}

func Test_MySQL_Finalizer_Default_Options(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{}}}
	finalizerMySQL(m.ConnectionDetails)
	r.Contains(m.URL(), "multiStatements=true")
	r.Contains(m.URL(), "parseTime=true")
	r.Contains(m.URL(), "readTimeout=3s")
	r.Contains(m.URL(), "collation=utf8mb4_general_ci")
}

func Test_MySQL_Finalizer_Preserve_User_Defined_Options(t *testing.T) {
	r := require.New(t)
	m := &mysql{commonDialect{ConnectionDetails: &ConnectionDetails{
		Options: map[string]string{
			"multiStatements": "false",
			"parseTime":       "false",
			"readTimeout":     "1h",
			"collation":       "utf8",
		},
	}}}
	finalizerMySQL(m.ConnectionDetails)
	r.Contains(m.URL(), "multiStatements=false")
	r.Contains(m.URL(), "parseTime=false")
	r.Contains(m.URL(), "readTimeout=1h")
	r.Contains(m.URL(), "collation=utf8")
}

func (s *MySQLSuite) Test_MySQL_DDL_Operations() {
	r := s.Require()

	origDatabase := PDB.Dialect.Details().Database
	PDB.Dialect.Details().Database = "pop_test_mysql_extra"
	defer func() {
		PDB.Dialect.Details().Database = origDatabase
	}()

	PDB.Dialect.DropDB() // clean up
	err := PDB.Dialect.CreateDB()
	r.NoError(err)
	err = PDB.Dialect.CreateDB()
	r.Error(err)
	err = PDB.Dialect.DropDB()
	r.NoError(err)
	err = PDB.Dialect.DropDB()
	r.Error(err)
}

func (s *MySQLSuite) Test_MySQL_DDL_Schema() {
	r := s.Require()

	f, err := ioutil.TempFile(os.TempDir(), "pop_test_mysql_dump")
	r.NoError(err)
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	// do it against "pop_test"
	err = PDB.Dialect.DumpSchema(f)
	r.NoError(err)
	_, err = f.Seek(0, 0)
	r.NoError(err)
	err = PDB.Dialect.LoadSchema(f)
	r.NoError(err)

	origDatabase := PDB.Dialect.Details().Database
	PDB.Dialect.Details().Database = "pop_test_not_exist"
	defer func() {
		PDB.Dialect.Details().Database = origDatabase
	}()

	// do it against "pop_test_not_exist"
	_, err = f.Seek(0, 0)
	r.NoError(err)
	err = PDB.Dialect.LoadSchema(f)
	r.Error(err)
	err = PDB.Dialect.DumpSchema(f)
	r.Error(err)
}

//** DEPRECATED: preserve test cases below while deprecated codes alive
func Test_MySQL_Deprecated_CD_Encoding(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{Encoding: "myEncoding"}
	finalizerMySQL(cd)
	r.NotNil(cd.Options)
	r.Equal("myEncoding", cd.Encoding)
	r.Equal("myEncoding", cd.Options["collation"])
}
