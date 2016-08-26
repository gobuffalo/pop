package pop

import (
	"fmt"
	"os"
)

func CreateDB(c *Connection) error {
	var err error
	deets := c.Dialect.Details()
	if deets.Database != "" {
		Log(fmt.Sprintf("Create %s (%s)", deets.Database, c.URL()))
		err = c.Dialect.CreateDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}

func DropDB(c *Connection) error {
	var err error
	deets := c.Dialect.Details()
	if deets.Database != "" {
		Log(fmt.Sprintf("Drop %s (%s)", deets.Database, c.URL()))
		err = c.Dialect.DropDB()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return err
}
