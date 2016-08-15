package pop

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/going/defaults"
)

type ConnectionDetails struct {
	Dialect  string
	Database string
	Host     string
	Port     string
	User     string
	Password string
	URL      string
	Options  map[string]string
}

// Parse extracts the various components of a connection string.
func (cd *ConnectionDetails) Parse(port string) error {
	if cd.URL != "" {
		u, err := url.Parse(cd.URL)
		if err != nil {
			return err
		}
		cd.Dialect = u.Scheme
		cd.Database = u.Path
		if cd.Dialect != "sqlite3" {
			cd.Database = strings.TrimPrefix(u.Path, "/")
		}
		hp := strings.Split(u.Host, ":")
		cd.Host = hp[0]
		if len(hp) > 1 {
			cd.Port = defaults.String(hp[1], port)
		}
		if u.User != nil {
			cd.User = u.User.Username()
			cd.Password, _ = u.User.Password()
		}
	}
	return nil
}

func (cd *ConnectionDetails) RetrySleep() time.Duration {
	d, err := time.ParseDuration(defaults.String(cd.Options["retry_sleep"], "1ms"))
	if err != nil {
		return 1 * time.Millisecond
	}
	return d
}

func (cd *ConnectionDetails) RetryLimit() int {
	i, err := strconv.Atoi(defaults.String(cd.Options["retry_limit"], "1000"))
	if err != nil {
		return 100
	}
	return i
}
