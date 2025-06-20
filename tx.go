package pop

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jmoiron/sqlx"
	"github.com/ory/pop/v6/internal/randx"
)

// Tx stores a transaction with an ID to keep track.
type Tx struct {
	ID int
	*sqlx.Tx
	db   *sql.DB
	pool *pgxpool.Pool
}

func newTX(ctx context.Context, db *dB, pool *pgxpool.Pool, opts *sql.TxOptions) (*Tx, error) {
	t := &Tx{
		ID:   randx.NonNegativeInt(),
		db:   db.SQLDB(),
		pool: pool,
	}
	tx, err := db.BeginTxx(ctx, opts)
	t.Tx = tx
	if err != nil {
		return nil, fmt.Errorf("could not create new transaction: %w", err)
	}
	return t, nil
}

func (tx *Tx) SQLDB() *sql.DB {
	return tx.db
}

func (tx *Tx) PGXPool() *pgxpool.Pool {
	return tx.pool
}

func (tx *Tx) PingContext(ctx context.Context) error {
	return tx.db.PingContext(ctx)
}

// TransactionContext simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *Tx) TransactionContext(ctx context.Context) (*Tx, error) {
	return tx, nil
}

// TransactionContextOptions simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *Tx) TransactionContextOptions(_ context.Context, _ *sql.TxOptions) (*Tx, error) {
	return tx, nil
}

// Transaction simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *Tx) Transaction() (*Tx, error) {
	return tx, nil
}

// Close does nothing. This is defined so it implements the `Store` interface.
func (tx *Tx) Close() error {
	return nil
}

// Workaround for https://github.com/jmoiron/sqlx/issues/447
func (tx *Tx) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return sqlx.NamedQueryContext(ctx, tx, query, arg)
}
