package pop

import "github.com/jmoiron/sqlx"

type DB struct {
	*sqlx.DB
}

func (db *DB) Transaction() (*TX, error) {
	return NewTX(db)
}

func (db *DB) Rollback() error {
	return nil
}

func (db *DB) Commit() error {
	return nil
}
