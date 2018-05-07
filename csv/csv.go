package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
)

// ErrInvalidCSV is returned when the input CSV file is not valid.
var ErrInvalidCSV = errors.New("invalid CSV file")

// Import loads a CSV file from its name, and insert rows into the datasource
// using the given connection.
func Import(c *pop.Connection, f string, model interface{}) error {
	m := &pop.Model{Value: model}
	// Load CSV entries
	fd, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fd.Close()
	r := csv.NewReader(fd)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	if len(records) <= 1 {
		return ErrInvalidCSV
	}

	// Create the insert query

	// The first row are the columns
	var qBuffer bytes.Buffer
	h := records[0]
	qBuffer.WriteString(fmt.Sprintf("INSERT INTO %s (", m.TableName()))
	hl := len(h) - 1
	for i := 0; i < hl; i++ {
		qBuffer.WriteString(h[i])
		qBuffer.WriteString(",")
	}

	// Write the last field, so we don't need to add a condition in the loop
	qBuffer.WriteString(h[hl])
	qBuffer.WriteString(") VALUES(")

	// Write values tokens
	qBuffer.WriteString(strings.Repeat("?, ", hl))
	qBuffer.WriteString("?)")

	// Prepare and use the query for each value
	q := qBuffer.String()
	stmt, err := c.Store.Preparex(q)
	if err != nil {
		return err
	}
	values := records[1:len(records)]
	for _, v := range values {
		pop.Log(q, v)
		_, err := stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}
