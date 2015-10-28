# POP

## A Tasty Treat For All Your Database Needs

So what does Pop do exactly? Well, it wraps the absolutely amazing [https://github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) library. It cleans up some of the common patterns and workflows usually associated with dealing with databases in Go.

Pop makes it easy to do CRUD operations, run migrations, and build/execute queries. Is Pop an ORM? I'll leave that up to you, the reader, to decide.

## Connecting to Databases

Pop is easily configured using a YAML file. The configuration file should be stored in `config/database.yml` or `database.yml`.

#### Example Configuration File

```yaml
development:
  dialect: "postgres"
  database: "your_db_development"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"

test:
  dialect: "mysql"
  database: "your_db_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

staging:
  dialect: "sqlite3"
  url: "./staging.sqlite"

production:
  dialect: "postgres"
  url: {{ env "DATABASE_URL" }}
```

asdf
