package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// ErrInvalidCSV is returned when the input CSV file is not valid.
var ErrInvalidCSV = errors.New("invalid CSV file")

// Importer defines the options for the pop CSV import feature.
type Importer struct {
	cnx     *pop.Connection
	elapsed int64
	// Comma is the field delimiter.
	// It is set to comma (',') by default.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune
}

func (ci *Importer) timeFunc(name string, fn func() error) error {
	now := time.Now()
	err := fn()
	atomic.AddInt64(&ci.elapsed, int64(time.Now().Sub(now)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// FromReader loads a CSV file from an io.Reader, and insert rows into the datasource
// using the given connection. It allows to use a custom file reader, such as a remote
// one.
//
// Example usage:
// 	imp := csv.NewImporter(tx)
//  	fd, err := os.Open("./csv/files/my_table_data.csv")
//	if err != nil {
//		return err
//	}
//	defer fd.Close()
//	return ci.FromReader(fd, "my_table")
func (ci *Importer) FromReader(fd io.Reader, model interface{}) error {
	m := &pop.Model{Value: model}
	return ci.timeFunc("CSV import", func() error {
		r := csv.NewReader(fd)
		r.Comma = ci.Comma
		h, err := r.Read()
		if err != nil {
			if err == io.EOF {
				// Header not found!
				return ErrInvalidCSV
			}
			return err
		}

		// Prepare the query from the header values
		var qBuffer bytes.Buffer
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

		for {
			// Read next line
			v, err := r.Read()
			if err != nil {
				if err == io.EOF {
					// No more lines to read
					break
				}
				return err
			}
			iv := make([]interface{}, len(v))
			for i := range v {
				iv[i] = v[i]
			}
			pop.Log(q, iv...)
			_, err = stmt.Exec(iv...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// FromPath loads a CSV file from its name, and insert rows into the datasource
// using the given connection.
//
// Example usage:
// 	imp := csv.NewImporter(tx)
// 	err := imp.FromPath("./csv/files/my_table_data.csv", "my_table")
func (ci *Importer) FromPath(f string, model interface{}) error {
	fd, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fd.Close()
	return ci.FromReader(fd, model)
}

// NewImporter creates a new CSV importer.
func NewImporter(c *pop.Connection) *Importer {
	return &Importer{
		cnx:   c,
		Comma: ',',
	}
}
