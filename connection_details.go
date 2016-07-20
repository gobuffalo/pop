package pop

import (
	"net/url"
	"strings"

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
}

// Parse extracts the various components of a connection string.
func (cd *ConnectionDetails) Parse(port string) error {
	if cd.URL != "" {
		u, err := url.Parse(cd.URL)
		if err != nil {
			return err
		}
		cd.Dialect = u.Scheme
		cd.Database = strings.TrimPrefix(u.Path, "/")
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
