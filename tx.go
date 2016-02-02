package pop

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tX struct {
	ID int
	*sqlx.Tx
}

func newTX(db *dB) (*tX, error) {
	t := &tX{
		ID: rand.Int(),
	}
	tx, err := db.Beginx()
	t.Tx = tx
	return t, err
}

// Transaction simply returns the current transaction,
// this is defined so it implements the `Store` interface.
func (tx *tX) Transaction() (*tX, error) {
	return tx, nil
}
