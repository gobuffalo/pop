# POP [![GoDoc](https://godoc.org/github.com/markbates/pop?status.svg)](https://godoc.org/github.com/markbates/pop) [![Build Status](https://travis-ci.org/markbates/pop.svg)](https://travis-ci.org/markbates/pop)

## A Tasty Treat For All Your Database Needs

So what does Pop do exactly? Well, it wraps the absolutely amazing [https://github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) library. It cleans up some of the common patterns and workflows usually associated with dealing with databases in Go.

Pop makes it easy to do CRUD operations, run migrations, and build/execute queries. Is Pop an ORM? I'll leave that up to you, the reader, to decide.

Pop, by default, follows conventions that were defined by the ActiveRecord Ruby gem, http://www.rubyonrails.org. What does this mean?

* Tables must have an "id" column and a corresponding "ID" field on the `struct` being used.
* If there is a timestamp column named "created_at", "CreatedAt" on the `struct`, it will be set with the current time when the record is created.
* If there is a timestamp column named "updated_at", "UpdatedAt" on the `struct`, it will be set with the current time when the record is updated.
* Default databases are lowercase, underscored versions of the `struct` name. Examples: User{} is "users", FooBar{} is "foo_bars", etc...

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

Note that the `database.yml` file is also a Go template, so you can use Go template syntax. There are two special functions that are included, `env` and `envOr`.

* `env` - This function will look for the named environment variable and insert it into your file. This is useful for configuration production databases without having to store secret information in your repository. `{{ env "DATABASE_URL" }}`
* `envOr` - This function will look for the named environment variable and use it. If the variable can not be found a default value will be used. `{{ envOr "MYSQL_HOST" "localhost" }}`

### In your code

Once you have a configuration file defined you can easily connect to one of these connections in your application.

```go
db, err := pop.Connect("development")
if err != nil {
  log.Panic(err)
}
```

Now that you have your connection to the database you can start executing queries against it.

## CLI Support

Pop features CLI support via the `soda` for the following operations:

* creating databases
* dropping databases
* migrating databases

### Installing CLI Support

```bash
$ go get -d -t -u github.com/markbates/pop/...
$ go install github.com/markbates/pop/soda
```

### Creating Databases

Assuming you defined a configuration file like that described in the above section you can automatically create those databases using the `soda` command:

#### Create All Databases

```bash
$ soda create -all
```

#### Create a Specific Database

```bash
$ soda create -e development
```

### Dropping Databases

Assuming you defined a configuration file like that described in the above section you can automatically drop those databases using the `soda` command:

#### Drop All Databases

```bash
$ soda drop -all
```

#### Drop a Specific Database

```bash
$ soda drop -e development
```

### Migrations

The `soda` command supports the creation and running of SQL migrations. Yes, I said SQL! There is no fancy DSL for writing these files, just plain old SQL, which, personally, I think is DSL enough.

A full list of commands available for migration can be found by asking for help:

```bash
$ soda migrate help
```

#### Create Migrations

The `soda` command will generate SQL migrations (both the up and down) files for you.

```bash
$ soda migrate create name_of_migration
```

Running this command with generate the following files:

```text
./migrations/0001_name_of_migration.up.sql
./migrations/0001_name_of_migration.down.sql
```

It is up to you to fill these files with the appropriate SQL to do whatever it is you need done.

#### Running Migrations

The `soda` command will run the migrations using the following command:

```bash
$ soda migrate up
```

Migrations will be run in sequential order. The previously run migrations will be kept track of in a table named `schema_migrations` in the database.

Migrations can also be run reverse to rollback the schema.

```bash
$ soda migrate down
```
