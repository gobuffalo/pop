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

// Importer defines the options for the pop CSV import feature.
type Importer struct {
	cnx *pop.Connection
	// Comma is the field delimiter.
	// It is set to comma (',') by default.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune
}

// Import loads a CSV file from its name, and insert rows into the datasource
// using the given connection.
//
// Example usage:
// 	imp := csv.NewImporter(tx)
// 	err := imp.Import("./csv/files/my_table_data.csv", "my_table")
func (ci *Importer) Import(f string, model interface{}) error {
	m := &pop.Model{Value: model}
	// Load CSV entries
	fd, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fd.Close()
	r := csv.NewReader(fd)
	r.Comma = ci.Comma
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
		qBuffer.WriteString(", ")
	}

	// Write the last field, so we don't need to add a condition in the loop
	qBuffer.WriteString(h[hl])
	qBuffer.WriteString(") VALUES(")

	// Write values tokens
	qBuffer.WriteString(strings.Repeat("?, ", hl))
	qBuffer.WriteString("?)")

	// Prepare and use the query for each value
	q := ci.cnx.Dialect.TranslateSQL(qBuffer.String())
	stmt, err := ci.cnx.Store.Preparex(q)
	if err != nil {
		return err
	}
	values := records[1:len(records)]
	for _, v := range values {
		iv := make([]interface{}, len(v))
		for i := range v {
			iv[i] = v[i]
		}
		pop.Log(q, iv...)
		_, err := stmt.Exec(iv...)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewImporter creates a new CSV importer.
func NewImporter(c *pop.Connection) *Importer {
	return &Importer{
		cnx:   c,
		Comma: ',',
	}
}
