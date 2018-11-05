package pop

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/pop/logging"
	"github.com/markbates/going/defaults"
	"github.com/markbates/oncer"
	"github.com/pkg/errors"
)

// ConnectionDetails stores the data needed to connect to a datasource
type ConnectionDetails struct {
	// Example: "postgres" or "sqlite3" or "mysql"
	Dialect string
	// The name of your database. Example: "foo_development"
	Database string
	// The host of your database. Example: "127.0.0.1"
	Host string
	// The port of your database. Example: 1234
	// Will default to the "default" port for each dialect.
	Port string
	// The username of the database user. Example: "root"
	User string
	// The password of the database user. Example: "password"
	Password string
	// The encoding to use to create the database and communicate with it.
	Encoding string
	// Instead of specifying each individual piece of the
	// connection you can instead just specify the URL of the
	// database. Example: "postgres://postgres:postgres@localhost:5432/pop_test?sslmode=disable"
	URL string
	// Defaults to 0 "unlimited". See https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	Pool int
	// Defaults to 0 "unlimited". See https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	IdlePool int
	Options  map[string]string
	// Query string encoded options from URL. Example: "sslmode=disable"
	RawOptions string
}

var dialectX = regexp.MustCompile(`\S+://`)

// withURL parses and overrides all connection details with values
// from standard URL except Dialect. It also calls dialect specific
// URL parser if exists.
func (cd *ConnectionDetails) withURL() error {
	ul := cd.URL
	if cd.Dialect != "" && !dialectX.MatchString(ul) {
		ul = cd.Dialect + "://" + ul
	}
	u, err := url.Parse(ul)
	if err != nil {
		return errors.Wrapf(err, "couldn't parse %s", ul)
	}

	//! dialect should not be overrided here (especially for cockroach)
	if cd.Dialect == "" {
		cd.Dialect = normalizeSynonyms(u.Scheme)
	}
	if !DialectSupported(cd.Dialect) {
		return errors.Errorf("unsupported dialect '%v'", cd.Dialect)
	}

	// warning message is required to prevent confusion
	// even though this behavior was documented.
	if cd.Database+cd.Host+cd.Port+cd.User+cd.Password != "" {
		log(logging.Warn, "One or more of connection details are specified in database.yml. Override them with values in URL.")
	}

	if up, ok := urlParser[cd.Dialect]; ok {
		return up(cd)
	}

	cd.Database = strings.TrimPrefix(u.Path, "/")

	hp := strings.Split(u.Host, ":")
	cd.Host = hp[0]
	if len(hp) > 1 {
		cd.Port = hp[1]
	}

	if u.User != nil {
		cd.User = u.User.Username()
		cd.Password, _ = u.User.Password()
	}
	cd.RawOptions = u.RawQuery

	return nil
}

// Finalize cleans up the connection details by normalizing names,
// filling in default values, etc...
func (cd *ConnectionDetails) Finalize() error {
	cd.Dialect = normalizeSynonyms(cd.Dialect)

	if cd.Options == nil { // for safety
		cd.Options = make(map[string]string)
	}

	if cd.URL != "" {
		if err := cd.withURL(); err != nil {
			return err
		}
	}

	if fin, ok := finalizer[cd.Dialect]; ok {
		fin(cd)
	}

	if DialectSupported(cd.Dialect) {
		if cd.Database != "" || cd.URL != "" {
			return nil
		}
		return errors.New("no database or URL specified")
	}
	return errors.Errorf("unsupported dialect '%v'", cd.Dialect)
}

// Parse cleans up the connection details by normalizing names,
// filling in default values, etc...
// Deprecated: use ConnectionDetails.Finalize() instead.
func (cd *ConnectionDetails) Parse(port string) error {
	oncer.Deprecate(0, "pop.ConnectionDetails#Parse", "pop.ConnectionDetails#Finalize")
	return cd.Finalize()
}

// RetrySleep returns the amount of time to wait between two connection retries
func (cd *ConnectionDetails) RetrySleep() time.Duration {
	d, err := time.ParseDuration(defaults.String(cd.Options["retry_sleep"], "1ms"))
	if err != nil {
		return 1 * time.Millisecond
	}
	return d
}

// RetryLimit returns the maximum number of accepted connection retries
func (cd *ConnectionDetails) RetryLimit() int {
	i, err := strconv.Atoi(defaults.String(cd.Options["retry_limit"], "1000"))
	if err != nil {
		return 100
	}
	return i
}

// MigrationTableName returns the name of the table to track migrations
func (cd *ConnectionDetails) MigrationTableName() string {
	return defaults.String(cd.Options["migration_table_name"], "schema_migration")
}

// OptionsString returns URL parameter encoded string from options.
func (cd *ConnectionDetails) OptionsString(s string) string {
	if cd.RawOptions != "" {
		return cd.RawOptions
	}
	if cd.Options != nil {
		for k, v := range cd.Options {
			s = fmt.Sprintf("%s&%s=%s", s, k, v)
		}
	}
	return strings.TrimLeft(s, "&")
}
