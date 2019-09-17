package pop

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type dB struct {
	*sqlx.DB
}

func (db *dB) TransactionContext(ctx context.Context) (*Tx, error) {
	return newTX(ctx, db)
}

func (db *dB) Transaction() (*Tx, error) {
	return newTX(context.Background(), db)
}

func (db *dB) Rollback() error {
	return nil
}

func (db *dB) Commit() error {
	return nil
}
