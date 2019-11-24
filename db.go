package pop

import (
	"context"

	"github.com/jmoiron/sqlx"
)

var _ store = new(dB)

type dB struct {
	*sqlx.DB
}

func (db *dB) TransactionContext(context.Context) (*Tx, error) {
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
