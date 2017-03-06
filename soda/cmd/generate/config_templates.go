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
  url: {{"{{"}}envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/{{.name}}_production"}}`

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
  url: {{"{{"}}envOr "DATABASE_URL" "mysql://root:root@(localhost:3306)/{{.name}}_production"}}`

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
