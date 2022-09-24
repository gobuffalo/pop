package pop

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gobuffalo/pop/v6/internal/defaults"
	"github.com/gobuffalo/pop/v6/internal/randx"
	"github.com/gobuffalo/pop/v6/logging"
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

// Context returns the connection's context set by "Context()" or context.TODO()
// if no context is set.
func (c *Connection) Context() context.Context {
	if c, ok := c.Store.(interface{ Context() context.Context }); ok {
		return c.Context()
	}

	return context.TODO()
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
		return nil, err
	}
	c := &Connection{
		ID: randx.String(30),
	}

	if nc, ok := newConnection[deets.Dialect]; ok {
		c.Dialect, err = nc(deets)
		if err != nil {
			return c, fmt.Errorf("could not create new connection: %w", err)
		}
		return c, nil
	}
	return nil, fmt.Errorf("could not found connection creator for %v", deets.Dialect)
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
			return nil, err
		}
	}
	e = defaults.String(e, "development")
	c := Connections[e]
	if c == nil {
		return c, fmt.Errorf("could not find connection named %s", e)
	}

	if err := c.Open(); err != nil {
		return c, fmt.Errorf("couldn't open connection for %s: %w", e, err)
	}
	return c, nil
}

// Open creates a new datasource connection
func (c *Connection) Open() error {
	if c.Store != nil {
		return nil
	}
	if c.Dialect == nil {
		return errors.New("invalid connection instance")
	}
	details := c.Dialect.Details()

	db, err := openPotentiallyInstrumentedConnection(c.Dialect, c.Dialect.URL())
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(details.Pool)
	if details.IdlePool != 0 {
		db.SetMaxIdleConns(details.IdlePool)
	}
	if details.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(details.ConnMaxLifetime)
	}
	if details.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(details.ConnMaxIdleTime)
	}
	if details.Unsafe {
		db = db.Unsafe()
	}
	c.Store = &dB{db}

	if d, ok := c.Dialect.(afterOpenable); ok {
		if err := d.AfterOpen(c); err != nil {
			c.Store = nil
			return fmt.Errorf("could not open database connection: %w", err)
		}
	}
	return nil
}

// Close destroys an active datasource connection
func (c *Connection) Close() error {
	if err := c.Store.Close(); err != nil {
		return fmt.Errorf("couldn't close connection: %w", err)
	}
	return nil
}

// Transaction will start a new transaction on the connection. If the inner function
// returns an error then the transaction will be rolled back, otherwise the transaction
// will automatically commit at the end.
func (c *Connection) Transaction(fn func(tx *Connection) error) error {
	return c.Dialect.Lock(func() (err error) {
		var dberr error

		log(logging.SQL, "--- BEGIN Transaction ---")
		cn, err := c.NewTransaction()
		if err != nil {
			return err
		}

		defer func() {
			if ex := recover(); ex != nil {
				log(logging.SQL, "--- ROLLBACK Transaction (inner function panic) ---")
				dberr = cn.TX.Rollback()
				if dberr != nil {
					err = fmt.Errorf("database error while inner panic rollback: %w", dberr)
				}
				err = fmt.Errorf("transaction was rolled back due to inner panic")
			}
		}()

		err = fn(cn)
		if err != nil {
			log(logging.SQL, "--- ROLLBACK Transaction ---")
			dberr = cn.TX.Rollback()
		} else {
			log(logging.SQL, "--- END Transaction ---")
			dberr = cn.TX.Commit()
		}

		if dberr != nil {
			return fmt.Errorf("database error on committing or rolling back transaction: %w", dberr)
		}

		return err
	})

}

// Rollback will open a new transaction and automatically rollback that transaction
// when the inner function returns, regardless. This can be useful for tests, etc...
func (c *Connection) Rollback(fn func(tx *Connection)) error {
	cn, err := c.NewTransaction()
	if err != nil {
		return err
	}
	fn(cn)
	return cn.TX.Rollback()
}

// NewTransaction starts a new transaction on the connection
func (c *Connection) NewTransaction() (*Connection, error) {
	return c.NewTransactionContextOptions(c.Context(), nil)
}

// NewTransactionContext starts a new transaction on the connection using the provided context
func (c *Connection) NewTransactionContext(ctx context.Context) (*Connection, error) {
	return c.NewTransactionContextOptions(ctx, nil)
}

// NewTransactionContextOptions starts a new transaction on the connection using the provided context and transaction options
func (c *Connection) NewTransactionContextOptions(ctx context.Context, options *sql.TxOptions) (*Connection, error) {
	var cn *Connection
	if c.TX == nil {
		tx, err := c.Store.TransactionContextOptions(ctx, options)
		if err != nil {
			return cn, fmt.Errorf("couldn't start a new transaction: %w", err)
		}

		cn = &Connection{
			ID:      randx.String(30),
			Store:   contextStore{store: tx, ctx: ctx},
			Dialect: c.Dialect,
			TX:      tx,
		}
	} else {
		cn = c
	}
	return cn, nil
}

// WithContext returns a copy of the connection, wrapped with a context.
func (c *Connection) WithContext(ctx context.Context) *Connection {
	cn := c.copy()
	cn.Store = contextStore{
		store: cn.Store,
		ctx:   ctx,
	}
	return cn
}

func (c *Connection) copy() *Connection {
	return &Connection{
		ID:      randx.String(30),
		Store:   c.Store,
		Dialect: c.Dialect,
		TX:      c.TX,
	}
}

// Q creates a new "empty" query for the current connection.
func (c *Connection) Q() *Query {
	return Q(c)
}

// disableEager disables eager mode for current connection.
func (c *Connection) disableEager() {
	// The check technically is not required, because (*Connection).Eager() creates a (shallow) copy.
	// When not reusing eager connections, this should be safe.
	// However, this write triggers the go race detector.
	if c.eager {
		c.eager = false
		c.eagerFields = []string{}
	}
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
		return err
	}
	return nil
}
