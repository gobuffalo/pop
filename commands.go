package pop

import (
	"fmt"
	"os"
)

func CreateDB(c *Connection) error {
	err := c.Open()
	if err != nil {
		return err
	}
	deets := c.Dialect.Details()
	if deets.Database != "" {
		fmt.Printf("[POP] Create %s (%s)\n", deets.Database, c.URL())
		err = c.Dialect.CreateDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}

func DropDB(c *Connection) error {
	err := c.Open()
	if err != nil {
		return err
	}
	deets := c.Dialect.Details()
	if deets.Database != "" {
		fmt.Printf("[POP] Drop %s (%s)\n", deets.Database, c.URL())
		err = c.Dialect.DropDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}
