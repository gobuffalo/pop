package pop

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TX struct {
	ID int
	*sqlx.Tx
}

func NewTX(db *DB) (*TX, error) {
	t := &TX{
		ID: rand.Int(),
	}
	tx, err := db.Beginx()
	t.Tx = tx
	return t, err
}

func (tx *TX) Transaction() (*TX, error) {
	return tx, nil
}
