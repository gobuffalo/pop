package pop

import (
	"github.com/gobuffalo/pop/log"
	"github.com/pkg/errors"
)

// CreateDB creates a database, given a connection definition
func CreateDB(c *Connection) error {
	deets := c.Dialect.Details()
	if deets.Database != "" {
		log.DefaultLogger.WithField("database", deets.Database).WithField("url", c.URL()).Debug("Creating database")
		return errors.Wrapf(c.Dialect.CreateDB(), "couldn't create database %s", deets.Database)
	}
	return nil
}

// DropDB drops an existing database, given a connection definition
func DropDB(c *Connection) error {
	deets := c.Dialect.Details()
	if deets.Database != "" {
		log.DefaultLogger.WithField("database", deets.Database).WithField("url", c.URL()).Debug("Dropping database")
		return errors.Wrapf(c.Dialect.DropDB(), "couldn't drop database %s", deets.Database)
	}
	return nil
}
