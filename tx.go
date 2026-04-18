package pop

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
)

// Tx stores a transaction with an ID to keep track.
type Tx struct {
	ID int
	*sqlx.Tx
}

func newTX(ctx context.Context, db *dB, opts *sql.TxOptions) (*Tx, error) {
	t := &Tx{
		ID: rand.Int(),
	}
	tx, err := db.BeginTxx(ctx, opts)
	t.Tx = tx
	if err != nil {
		return nil, fmt.Errorf("could not create new transaction: %w", err)
	}
	return t, nil
}

// TransactionContext simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *Tx) TransactionContext(_ context.Context) (*Tx, error) {
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

// NamedQueryContext binds a named query and then runs Query on the result using the Transaction.
// It works with both structs and with map[string]any types.
func (tx *Tx) NamedQueryContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error) {
	// Workaround for https://github.com/jmoiron/sqlx/issues/447
	return sqlx.NamedQueryContext(ctx, tx, query, arg)
}
