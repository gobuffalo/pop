package pop

import (
	"fmt"
	"os"
)

func CreateDB(c *Connection) error {
	var err error
	if c.Dialect.Details().Database != "" {
		err = c.Dialect.CreateDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}

func DropDB(c *Connection) error {
	var err error
	if c.Dialect.Details().Database != "" {
		err = c.Dialect.DropDB()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
	return err
}
