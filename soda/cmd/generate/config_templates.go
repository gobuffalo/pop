package generate

var configTemplates = map[string]string{
	"postgres": `development:
  dialect: postgres
  database: {{.}}_development
  username: postgres
  password: postgres
  host: 127.0.0.1

test:
  dialect: postgres
  database: {{.}}_test
  username: postgres
  password: postgres
  host: 127.0.0.1

production:
  dialect: postgres
  database: {{.}}_production
  username: postgres
  password: postgres
  host: 127.0.0.1
`,
	"mysql": `development:
  dialect: "mysql"
  database: "{{.}}_development"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

test:
  dialect: "mysql"
  database: "{{.}}_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

production:
  dialect: "mysql"
  database: "{{.}}_production"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"
	`,
	"sqlite3": `development:
	dialect: "sqlite3"
	database: "./{{.}}_development.sqlite"

test:
	dialect: "sqlite3"
	database: "./{{.}}_test.sqlite"

production:
	dialect: "sqlite3"
	database: "./{{.}}_production.sqlite"
`,
}
