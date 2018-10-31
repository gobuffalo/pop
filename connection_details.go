package pop

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	_mysql "github.com/go-sql-driver/mysql"
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

// overrideWithURL parses and overrides all connection details
// with values form URL except Dialect.
func (cd *ConnectionDetails) overrideWithURL() error {
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
		cd.Dialect = u.Scheme
	}

	// warning message is required to prevent confusion
	// even though this behavior was documented.
	if cd.Database+cd.Host+cd.Port+cd.User+cd.Password != "" {
		log(logging.Warn, "One or more of connection parameters are specified in database.yml. Override them with values in URL.")
	}

	if strings.HasPrefix(cd.Dialect, "sqlite") {
		cd.Database = u.Path
		return nil
	} else if strings.HasPrefix(cd.Dialect, "mysql") {
		return cd.overrideWithMySQLURL()
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

func (cd *ConnectionDetails) overrideWithMySQLURL() error {
	// parse and verify whether URL is supported by underlying driver or not.
	cfg, err := _mysql.ParseDSN(strings.TrimPrefix(cd.URL, "mysql://"))
	if err != nil {
		return errors.Wrapf(err, "the URL '%s' is not supported by MySQL driver", cd.URL)
	}

	cd.User = cfg.User
	cd.Password = cfg.Passwd
	cd.Database = cfg.DBName
	cd.Encoding = defaults.String(cfg.Collation, "utf8_general_ci")
	addr := strings.TrimSuffix(strings.TrimPrefix(cfg.Addr, "("), ")")
	if cfg.Net == "unix" {
		cd.Port = "socket"
		cd.Host = addr
	} else {
		tmp := strings.Split(addr, ":")
		cd.Host = tmp[0]
		if len(tmp) > 1 {
			cd.Port = tmp[1]
		}
	}

	return nil
}

// Finalize cleans up the connection details by normalizing names,
// filling in default values, etc...
func (cd *ConnectionDetails) Finalize() error {
	cd.Dialect = normalizeSynonyms(cd.Dialect)
	// PostgreSQL connection string can't be parsed as URL, so let's skip finallization.
	// Example string: "user=pqgotest dbname=pqgotest sslmode=verify-full"
	// More information about the format: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	// TODO: When pg connection string recognized, parse it and fill in the data.
	if cd.URL != "" && !strings.Contains(cd.URL, "dbname=") {
		if err := cd.overrideWithURL(); err != nil {
			return err
		}
	}

	switch cd.Dialect {
	case "postgres":
		cd.Port = defaults.String(cd.Port, "5432")
	case "cockroach":
		cd.Port = defaults.String(cd.Port, "26257")
	case "mysql":
		cd.Port = defaults.String(cd.Port, "3306")
	case "sqlite3":
		// Nothing more to do here
	default:
		return errors.Errorf("unknown dialect %s", cd.Dialect)
	}
	return nil
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
