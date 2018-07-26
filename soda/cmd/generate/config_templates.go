package generate

var pgConfig = `development:
  dialect: postgres
  database: {{.name}}_development
  user: postgres
  password: postgres
  host: 127.0.0.1
  pool: 5

test:
  url: {{"{{"}}envOr "TEST_DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/{{.name}}_test?sslmode=disable"}}

production:
  url: {{"{{"}}envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/{{.name}}_production?sslmode=disable"}}`

var mysqlConfig = `development:
  dialect: "mysql"
  database: "{{.name}}_development"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"
  encoding: "utf8mb4_general_ci"

test:
  url: {{"{{"}}envOr "TEST_DATABASE_URL" "mysql://root:root@(localhost:3306)/{{.name}}_test?parseTime=true&multiStatements=true&readTimeout=1s"}}

production:
  url: {{"{{"}}envOr "DATABASE_URL" "mysql://root:root@(localhost:3306)/{{.name}}_production?parseTime=true&multiStatements=true&readTimeout=1s"}}`

var sqliteConfig = `development:
  dialect: "sqlite3"
  database: {{.sqlitePath}}_development.sqlite

test:
  dialect: "sqlite3"
  database: {{.sqlitePath}}_test.sqlite

production:
  dialect: "sqlite3"
  database: {{.sqlitePath}}_production.sqlite`

var cockroachConfig = `development:
  dialect: cockroach
  database: {{.name}}_development
  host: {{"{{"}}envOr "COCKROACH_HOST" "127.0.0.1"  }}
  port: {{"{{"}} envOr "COCKROACH_PORT" "26257"  }}
  user: {{"{{"}} envOr "COCKROACH_USER" "root"  }}
  password: {{"{{"}} envOr "COCKROACH_PASSWORD" ""  }}
  pool: 5

test:
  dialect: cockroach
  database: {{.name}}_test
  host: {{"{{"}} envOr "COCKROACH_HOST" "127.0.0.1"  }}
  port: {{"{{"}} envOr "COCKROACH_PORT" "26257"  }}
  user: {{"{{"}} envOr "COCKROACH_USER" "root"  }}
  password: {{"{{"}} envOr "COCKROACH_PASSWORD" ""  }}
  pool: 5

production:
  dialect: cockroach
  database: {{.name}}_production
  host: {{"{{"}} envOr "COCKROACH_HOST" "127.0.0.1"  }}
  port: {{"{{"}} envOr "COCKROACH_PORT" "26257"  }}
  user: {{"{{"}} envOr "COCKROACH_USER" "root"  }}
  password: {{"{{"}} envOr "COCKROACH_PASSWORD" ""  }}
  pool: 5
  `

var configTemplates = map[string]string{
	"postgres":   pgConfig,
	"postgresql": pgConfig,
	"pg":         pgConfig,
	"mysql":      mysqlConfig,
	"sqlite3":    sqliteConfig,
	"sqlite":     sqliteConfig,
	"cockroach":  cockroachConfig,
	"crdb":       cockroachConfig,
}
