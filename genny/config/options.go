package config

// Options needed for the config generator
type Options struct {
	Root     string // Defaults to PWD
	FileName string // Defaults to database.yml
	Dialect  string // required
	Prefix   string // required - <prefix>_development
}
