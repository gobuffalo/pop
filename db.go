package pop

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

func (db *Database) TransactionContext(ctx context.Context) (*Tx, error) {
	return newTX(ctx, db, nil)
}

func (db *Database) Transaction() (*Tx, error) {
	return newTX(context.Background(), db, nil)
}

func (db *Database) TransactionContextOptions(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	return newTX(ctx, db, opts)
}

func (db *Database) Rollback() error {
	return nil
}

func (db *Database) Commit() error {
	return nil
}

func (db *Database) Database() *Database {
	return db
}
