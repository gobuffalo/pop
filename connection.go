package pop

import (
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
	"github.com/markbates/going/randx"
	"github.com/pkg/errors"
)

// Connections contains all available connections
var Connections = map[string]*Connection{}

// Connection represents all necessary details to talk with a datastore
type Connection struct {
	ID          string
	Store       store
	Dialect     dialect
	Elapsed     int64
	TX          *Tx
	eager       bool
	eagerFields []string
}

func (c *Connection) String() string {
	return c.URL()
}

// URL returns the datasource connection string
func (c *Connection) URL() string {
	return c.Dialect.URL()
}

// MigrationURL returns the datasource connection string used for running the migrations
func (c *Connection) MigrationURL() string {
	return c.Dialect.MigrationURL()
}

// MigrationTableName returns the name of the table to track migrations
func (c *Connection) MigrationTableName() string {
	return c.Dialect.Details().MigrationTableName()
}

// NewConnection creates a new connection, and sets it's `Dialect`
// appropriately based on the `ConnectionDetails` passed into it.
func NewConnection(deets *ConnectionDetails) (*Connection, error) {
	err := deets.Finalize()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c := &Connection{
		ID: randx.String(30),
	}
	switch deets.Dialect {
	case "postgres":
		c.Dialect = newPostgreSQL(deets)
	case "cockroach":
		c.Dialect = newCockroach(deets)
	case "mysql":
		c.Dialect = newMySQL(deets)
	case "sqlite3":
		c.Dialect, err = newSQLite(deets)
		if err != nil {
			return c, errors.WithStack(err)
		}
	}
	return c, nil
}

// Connect takes the name of a connection, default is "development", and will
// return that connection from the available `Connections`. If a connection with
// that name can not be found an error will be returned. If a connection is
// found, and it has yet to open a connection with its underlying datastore,
// a connection to that store will be opened.
func Connect(e string) (*Connection, error) {
	if len(Connections) == 0 {
		err := LoadConfigFile()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	e = defaults.String(e, "development")
	c := Connections[e]
	if c == nil {
		return c, errors.Errorf("could not find connection named %s", e)
	}
	err := c.Open()
	return c, errors.Wrapf(err, "couldn't open connection for %s", e)
}

// Open creates a new datasource connection
func (c *Connection) Open() error {
	if c.Store != nil {
		return nil
	}
	details := c.Dialect.Details()
	db, err := sqlx.Open(details.Dialect, c.Dialect.URL())
	db.SetMaxOpenConns(details.Pool)
	db.SetMaxIdleConns(details.IdlePool)
	if err == nil {
		c.Store = &dB{db}
	}
	return errors.Wrap(err, "couldn't connect to database")
}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	return errors.Wrap(c.Store.Close(), "couldn't close connection")
}

// Transaction will start a new transaction on the connection. If the inner function
// returns an error then the transaction will be rolled back, otherwise the transaction
// will automatically commit at the end.
func (c *Connection) Transaction(fn func(tx *Connection) error) error {
	return c.Dialect.Lock(func() error {
		var dberr error
		cn, err := c.NewTransaction()
		if err != nil {
			return err
		}
		err = fn(cn)
		if err != nil {
			dberr = cn.TX.Rollback()
		} else {
			dberr = cn.TX.Commit()
		}
		if err != nil {
			return errors.WithStack(err)
		}
		return errors.Wrap(dberr, "error committing or rolling back transaction")
	})

}

// NewTransaction starts a new transaction on the connection
func (c *Connection) NewTransaction() (*Connection, error) {
	var cn *Connection
	if c.TX == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return cn, errors.Wrap(err, "couldn't start a new transaction")
		}
		cn = &Connection{
			ID:      randx.String(30),
			Store:   tx,
			Dialect: c.Dialect,
			TX:      tx,
		}
	} else {
		cn = c
	}
	return cn, nil
}

func (c *Connection) copy() *Connection {
	return &Connection{
		ID:      randx.String(30),
		Store:   c.Store,
		Dialect: c.Dialect,
		TX:      c.TX,
	}
}

// Rollback will open a new transaction and automatically rollback that transaction
// when the inner function returns, regardless. This can be useful for tests, etc...
func (c *Connection) Rollback(fn func(tx *Connection)) error {
	var cn *Connection
	if c.TX == nil {
		tx, err := c.Store.Transaction()
		if err != nil {
			return errors.Wrap(err, "couldn't start a new transaction")
		}
		cn = &Connection{
			ID:      randx.String(30),
			Store:   tx,
			Dialect: c.Dialect,
			TX:      tx,
		}
	} else {
		cn = c
	}
	fn(cn)
	return cn.TX.Rollback()
}

// Q creates a new "empty" query for the current connection.
func (c *Connection) Q() *Query {
	return Q(c)
}

// disableEager disables eager mode for current connection.
func (c *Connection) disableEager() {
	c.eager = false
	c.eagerFields = []string{}
}

// TruncateAll truncates all data from the datasource
func (c *Connection) TruncateAll() error {
	return c.Dialect.TruncateAll(c)
}

func (c *Connection) timeFunc(name string, fn func() error) error {
	start := time.Now()
	err := fn()
	atomic.AddInt64(&c.Elapsed, int64(time.Since(start)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
