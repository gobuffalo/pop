package pop

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"strings"
	"sync"

	mysqld "github.com/go-sql-driver/mysql"
	pgx "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/ory/pop/v6/logging"
)

const instrumentedDriverName = "instrumented-sql-driver"

var sqlDriverLock = sync.Mutex{}

func instrumentDriver(deets *ConnectionDetails, defaultDriverName string) (driverName, dialect string, err error) {
	driverName = defaultDriverName
	if deets.Driver != "" {
		driverName = deets.Driver
	}
	dialect = driverName

	if !deets.UseInstrumentedDriver {
		if len(deets.InstrumentedDriverOptions) > 0 {
			log(logging.Warn, "SQL driver instrumentation is disabled but `ConnectionDetails.InstrumentedDriverOptions` is not empty. Please double-check if this is a error.")
		}

		// If instrumentation is disabled, we just return the driver name we got (e.g. "pgx").
		return driverName, dialect, nil
	}

	if len(deets.InstrumentedDriverOptions) == 0 {
		log(logging.Warn, "SQL driver instrumentation was enabled but no options have been passed to `ConnectionDetails.InstrumentedDriverOptions`. Instrumentation will therefore not result in any output.")
	}

	var dr driver.Driver
	var newDriverName string
	switch CanonicalDialect(driverName) {
	case nameCockroach:
		fallthrough
	case namePostgreSQL:
		dr = new(pgx.Driver)
		newDriverName = instrumentedDriverName + "-" + namePostgreSQL
	case nameMariaDB:
		fallthrough
	case nameMySQL:
		dr = mysqld.MySQLDriver{}
		newDriverName = instrumentedDriverName + "-" + nameMySQL
	case nameSQLite3:
		var err error
		dr, err = newSQLiteDriver()
		if err != nil {
			return "", "", err
		}
		newDriverName = instrumentedDriverName + "-" + nameSQLite3
	}

	sqlDriverLock.Lock()
	defer sqlDriverLock.Unlock()

	var found bool
	for _, n := range sql.Drivers() {
		if n == newDriverName {
			found = true
			break
		}
	}

	if !found {
		sql.Register(newDriverName, instrumentedsql.WrapDriver(dr, deets.InstrumentedDriverOptions...))
	}

	return newDriverName, dialect, nil
}

// openPotentiallyInstrumentedConnection first opens a raw SQL connection and then wraps it with `sqlx`.
//
// We do this because `sqlx` needs the database type in order to properly
// translate arguments (e.g. `?` to `$1`) in SQL queries. Because we use
// a custom driver name when using instrumentation, this detection would fail
// otherwise.
func openPotentiallyInstrumentedConnection(ctx context.Context, c dialect, dsn string) (*sqlx.DB, *pgxpool.Pool, error) {
	driverName, dialect, err := instrumentDriver(c.Details(), c.DefaultDriver())
	if err != nil {
		return nil, nil, err
	}

	// If "pool_min_conns" is set in the DSN, it means that we use the pgx pool feature flag.
	if strings.Contains(dsn, "pool_min_conns=") {
		// But of course only on Cockroach and PostgreSQL.
		switch CanonicalDialect(c.DefaultDriver()) {
		case nameCockroach:
			fallthrough
		case namePostgreSQL:
			pool, err := pgxpool.New(ctx, dsn)
			if err != nil {
				return nil, nil, err
			}

			// We don't need to configure the database/sql connection pool because pgxpool already does that.
			// Reference: https://github.com/jackc/pgx/discussions/2222#discussioncomment-11772064
			return sqlx.NewDb(stdlib.OpenDBFromPool(pool), dialect), pool, nil
		}
	}

	con, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open database connection: %w", err)
	}

	details := c.Details()
	if details.Pool != 0 {
		con.SetMaxOpenConns(details.Pool)
	}
	if details.IdlePool != 0 {
		con.SetMaxIdleConns(details.IdlePool)
	}
	if details.ConnMaxLifetime > 0 {
		con.SetConnMaxLifetime(details.ConnMaxLifetime)
	}
	if details.ConnMaxIdleTime > 0 {
		con.SetConnMaxIdleTime(details.ConnMaxIdleTime)
	}

	return sqlx.NewDb(con, dialect), nil, nil
}
