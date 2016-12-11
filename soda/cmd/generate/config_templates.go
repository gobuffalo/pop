package generate

var pgConfig = `development:
  dialect: postgres
  database: {{.name}}_development
  user: postgres
  password: postgres
  host: 127.0.0.1
  pool: 5

test:
  dialect: postgres
  database: {{.name}}_test
  user: postgres
  password: postgres
  host: 127.0.0.1

production:
  dialect: postgres
  database: {{.name}}_production
  user: postgres
  password: postgres
  host: 127.0.0.1
  pool: 25`

var mysqlConfig = `development:
  dialect: "mysql"
  database: "{{.name}}_development"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

test:
  dialect: "mysql"
  database: "{{.name}}_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

production:
  dialect: "mysql"
  database: "{{.name}}_production"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"`

var sqliteConfig = `development:
  dialect: "sqlite3"
  database: "./{{.name}}_development.sqlite"

test:
  dialect: "sqlite3"
  database: "./{{.name}}_test.sqlite"

production:
  dialect: "sqlite3"
  database: "./{{.name}}_production.sqlite"`

var configTemplates = map[string]string{
	"postgres":   pgConfig,
	"postgresql": pgConfig,
	"pg":         pgConfig,
	"mysql":      mysqlConfig,
	"sqlite3":    sqliteConfig,
	"sqlite":     sqliteConfig,
}
