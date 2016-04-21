package pop

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
)

// Connections contains all of the available connections
var Connections = map[string]*Connection{}

// Connection represents all of the necessary details for
// talking with a datastore
type Connection struct {
	Store   Store
	Dialect Dialect
	Timings []time.Duration
}

func (c *Connection) String() string {
	return c.Dialect.URL()
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `ConnectionDetails` passed into it.
func NewConnection(deets *ConnectionDetails) *Connection {
	c := &Connection{
		Timings: []time.Duration{},
	}
	switch deets.Dialect {
	case "postgres":
		c.Dialect = NewPostgreSQL(deets)
	case "mysql":
		c.Dialect = NewMySQL(deets)
		// case "sqlite3":
		// 	c.Dialect = NewSQLite(deets)
	}
	return c
}

// Connect takes the name of a connection, default is "development", and will
// return that connection from the available `Connections`. If a connection with
// that name can not be found an error will be returned. If a connection is
// found, and it has yet to open a connection with its underlying datastore,
// a connection to that store will be opened.
func Connect(e string) (*Connection, error) {
	e = defaults.String(e, "development")
	c := Connections[e]
	if c == nil {
		return c, fmt.Errorf("Could not find connection named %s!", e)
	}
	if c.Store != nil {
		return c, nil
	}
	db, err := sqlx.Open(c.Dialect.Details().Dialect, c.Dialect.URL())
	if err == nil {
		c.Store = &dB{db}
	}
	return c, nil
}

// Transaction will start a new transaction on the connection. If the inner function
// returns an error then the transaction will be rolled back, otherwise the transaction
// will automatically commit at the end.
func (c *Connection) Transaction(fn func(tx *Connection) error) error {
	tx, err := c.Store.Transaction()
	if err != nil {
		return err
	}
	cn := &Connection{
		Store:   tx,
		Dialect: c.Dialect,
		Timings: []time.Duration{},
	}
	err = fn(cn)
	if err != nil {
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return err
}

// Rollback will open a new transaction and automatically rollback that transaction
// when the inner function returns, regardless. This can be useful for tests, etc...
func (c *Connection) Rollback(fn func(tx *Connection)) error {
	tx, err := c.Store.Transaction()
	if err != nil {
		return err
	}
	cn := &Connection{
		Store:   tx,
		Dialect: c.Dialect,
		Timings: []time.Duration{},
	}
	fn(cn)
	return tx.Rollback()
}

// Q creates a new "empty" query for the current connection.
func (c *Connection) Q() *Query {
	return Q(c)
}

func (c *Connection) timeFunc(name string, fn func() error) error {
	now := time.Now()
	err := fn()
	c.Timings = append(c.Timings, time.Now().Sub(now))
	return err
}
