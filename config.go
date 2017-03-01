package pop

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/markbates/going/defaults"
	"gopkg.in/yaml.v2"
)

var lookupPaths = []string{"", "./config", "/config", "../", "../config", "../..", "../../config"}
var ConfigName = "database.yml"

func init() {
	ap := os.Getenv("APP_PATH")
	if ap != "" {
		AddLookupPaths(ap)
	}
	ap = os.Getenv("POP_PATH")
	if ap != "" {
		AddLookupPaths(ap)
	}
}

// LoadConfig loads the configuration file and resets all of the connections pop knows about
func LoadConfig() (map[string]*Connection, error) {
	path, err := findConfigPath()
	if err != nil {
		return nil, err
	}

	if Debug {
		fmt.Printf("[POP]: Loading config file from %s\n", path)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read file %s", path)
	}

	tmpl := template.New("test")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(s1, s2 string) string {
			return defaults.String(os.Getenv(s1), s2)
		},
		"env": func(s1 string) string {
			return os.Getenv(s1)
		},
	})
	t, err := tmpl.Parse(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse config template")
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute config template")
	}

	deets := map[string]*ConnectionDetails{}
	err = yaml.Unmarshal(bb.Bytes(), &deets)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal config to yaml")
	}

	connections = make(map[string]*Connection)
	for n, d := range deets {
		con, err := NewConnection(d)
		if err != nil {
			return nil, err
		}
		connections[n] = con
	}
	return connections, nil
}

func LookupPaths() []string {
	return lookupPaths
}

func AddLookupPaths(paths ...string) {
	lookupPaths = append(paths, lookupPaths...)
	LoadConfig()
}

func findConfigPath() (string, error) {
	for _, p := range LookupPaths() {
		path, _ := filepath.Abs(filepath.Join(p, ConfigName))
		if _, err := os.Stat(path); err == nil {
			return path, err
		}
	}
	return "", errors.New("[POP]: Tried to load configuration file, but couldn't find it.")
}
