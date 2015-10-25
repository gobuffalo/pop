package pop

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
)

var Connections = map[string]*Connection{}

type Connection struct {
	Store   Store
	Dialect Dialect
	Timings []time.Duration
}

func (c *Connection) String() string {
	return c.Dialect.URL()
}

func NewConnection(deets *ConnectionDetails) *Connection {
	c := &Connection{
		Timings: []time.Duration{},
	}
	switch deets.Dialect {
	case "postgres":
		c.Dialect = NewPostgreSQL(deets)
	case "mysql":
		c.Dialect = NewMySQL(deets)
	case "sqlite3":
		c.Dialect = NewSQLite(deets)
	}
	return c
}

func Connect(e string) (*Connection, error) {
	e = defaults.String(e, "development")
	c := Connections[e]
	if c == nil {
		return c, fmt.Errorf("Could not find connection named %s!", e)
	}
	db, err := sqlx.Open(c.Dialect.Details().Dialect, c.Dialect.URL())
	if err == nil {
		c.Store = &dB{db}
	}
	return c, nil
}

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
func (c *Connection) Q() *Query {
	return Q(c)
}

func (c *Connection) timeFunc(name string, fn func() error) error {
	now := time.Now()
	err := fn()
	c.Timings = append(c.Timings, time.Now().Sub(now))
	return err
}
