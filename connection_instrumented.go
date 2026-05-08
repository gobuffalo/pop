package pop

import (
	"database/sql"
	"database/sql/driver"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"

	mysqld "github.com/go-sql-driver/mysql"
	pgx "github.com/jackc/pgx/v5/stdlib"

	"github.com/gobuffalo/pop/v6/logging"
)

func instrumentDriver(deets *ConnectionDetails, defaultDriverName string) (drv driver.Driver, dialect string, err error) {
	driverName := defaultDriverName
	if deets.Driver != "" {
		driverName = deets.Driver
	}
	dialect = driverName

	if !deets.UseInstrumentedDriver {
		if len(deets.InstrumentedDriverOptions) > 0 {
			log(
				logging.Warn,
				"SQL driver instrumentation is disabled but `ConnectionDetails.InstrumentedDriverOptions` is not empty. Please double-check if this is a error.",
			)
		}

		// If instrumentation is disabled, return nil driver to signal non-instrumented path.
		return nil, dialect, nil
	}

	if len(deets.InstrumentedDriverOptions) == 0 {
		log(
			logging.Warn,
			"SQL driver instrumentation was enabled but no options have been passed to `ConnectionDetails.InstrumentedDriverOptions`. Instrumentation will therefore not result in any output.",
		)
	}

	switch CanonicalDialect(driverName) {
	case nameCockroach:
		fallthrough
	case namePostgreSQL:
		drv = new(pgx.Driver)
	case nameMariaDB:
		fallthrough
	case nameMySQL:
		drv = mysqld.MySQLDriver{}
	case nameSQLite3:
		drv, err = newSQLiteDriver()
		if err != nil {
			return nil, "", err
		}
	}

	return instrumentedsql.WrapDriver(drv, deets.InstrumentedDriverOptions...), dialect, nil
}

// openPotentiallyInstrumentedConnection first opens a raw SQL connection and then wraps it with `sqlx`.
//
// We do this because `sqlx` needs the database type in order to properly
// translate arguments (e.g. `?` to `$1`) in SQL queries. Because we use
// a custom driver name when using instrumentation, this detection would fail
// otherwise.
func openPotentiallyInstrumentedConnection(c dialect, dsn string) (*sqlx.DB, error) {
	drv, dialect, err := instrumentDriver(c.Details(), c.DefaultDriver())
	if err != nil {
		return nil, err
	}

	var con *sql.DB
	if drv != nil {
		// Use sql.OpenDB with the per-connection wrapped driver instead of
		// sql.Register + sql.Open, which only registers one driver per driver type
		// regardless of per-connection options.
		con = sql.OpenDB(&driverConnector{drv, dsn})
	} else {
		con, err = sql.Open(c.DefaultDriver(), dsn)
		if err != nil {
			return nil, fmt.Errorf("could not open database connection: %w", err)
		}
	}

	return sqlx.NewDb(con, dialect), nil
}

// driverConnector wraps a driver.Driver with a DSN to implement driver.Connector.
type driverConnector struct {
	driver driver.Driver
	dsn    string
}

func (dc *driverConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return dc.driver.Open(dc.dsn)
}

func (dc *driverConnector) Driver() driver.Driver {
	return dc.driver
}
