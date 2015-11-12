package commands

import (
	"github.com/codegangsta/cli"
	_ "github.com/mattes/migrate/driver/mysql"
	_ "github.com/mattes/migrate/driver/postgres"
	_ "github.com/mattes/migrate/driver/sqlite3"
)

var EnvFlag = cli.StringFlag{
	Name:  "e",
	Value: "development",
	Usage: "The environment you want to run migrations against. Will use $GO_ENV if set.",
}
