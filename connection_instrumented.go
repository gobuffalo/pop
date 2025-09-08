package pop

import (
	"cmp"
	"context"
	"fmt"
	"strings"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxstdlib "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// openPotentiallyInstrumentedConnection first opens a raw SQL connection and then wraps it with `sqlx`.
//
// We do this because `sqlx` needs the database type in order to properly
// translate arguments (e.g. `?` to `$1`) in SQL queries. Because we use
// a custom driver name when using instrumentation, this detection would fail
// otherwise.
func openPotentiallyInstrumentedConnection(ctx context.Context, c dialect, dsn string) (*sqlx.DB, *pgxpool.Pool, error) {
	var (
		driver   = cmp.Or(c.Details().Driver, c.DefaultDriver())
		otelopts = []otelsql.Option{
			otelsql.WithSpanOptions(otelsql.SpanOptions{
				Ping:     false,
				RowsNext: false,
				OmitRows: true,
			}),
			otelsql.WithTracerProvider(
				cmp.Or(
					c.Details().TracerProvider,
					trace.TracerProvider(noop.NewTracerProvider()),
				)),
		}
	)
	// If "pool_min_conns" is set in the DSN, it means that we use the pgx pool feature flag.
	if strings.Contains(dsn, "pool_min_conns=") && (CanonicalDialect(driver) == nameCockroach || CanonicalDialect(driver) == namePostgreSQL) {
		pool, err := pgxpool.New(ctx, dsn)
		if err != nil {
			return nil, nil, err
		}

		con := otelsql.OpenDB(pgxstdlib.GetPoolConnector(pool), append(otelopts, otelsql.WithSpanOptions(otelsql.SpanOptions{
			Ping:                 false,
			RowsNext:             false,
			OmitRows:             true,
			OmitConnectorConnect: true, // this does not work correctly with pgxpool
		}))...)
		con.SetMaxIdleConns(0) // important: see documentation for pgxstdlib.GetPoolConnector
		return sqlx.NewDb(con, driver), pool, nil
	} else {
		con, err := otelsql.Open(driver, dsn, otelopts...)
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
		return sqlx.NewDb(con, driver), nil, nil
	}
}
