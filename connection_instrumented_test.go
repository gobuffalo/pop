package pop

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/luna-duclos/instrumentedsql"
	"github.com/stretchr/testify/suite"
)

func testInstrumentedDriver(p *suite.Suite) {
	r := p.Require()
	deets := *Connections[os.Getenv("SODA_DIALECT")].Dialect.Details()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	// The WaitGroup and channel ensures that the logger is properly called. This can only happen
	// when the instrumented driver is working as expected and returns the expected query.
	var (
		queryMySQL = "SELECT 1 FROM DUAL WHERE 1=?"
		queryOther = "SELECT 1 WHERE 1=?"
		mc         = make(chan string)
		wg         sync.WaitGroup
		expected   = []string{
			"SELECT 1 FROM DUAL WHERE 1=?",
			"SELECT 1 FROM DUAL WHERE 1=$1",
			"SELECT 1 WHERE 1=?",
			"SELECT 1 WHERE 1=$1",
		}
	)

	query := queryOther
	if os.Getenv("SODA_DIALECT") == "mysql" {
		query = queryMySQL
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var messages []string
		var found bool
		for {
			select {
			case m := <-mc:
				p.T().Logf("Received message: %s", m)
				messages = append(messages, m)
				for _, e := range expected {
					if strings.Contains(m, e) {
						p.T().Logf("Found part %s in %s", e, m)
						found = true
						break
					}
				}
			case <-ctx.Done():
				if !found {
					r.FailNow(fmt.Sprintf("Expected tracer to return the \"%s\" query but only the following messages have been received:\n\n\t%s", query, strings.Join(messages, "\n\t")))
					return
				}
				return
			}
		}
	}()

	var checker = instrumentedsql.LoggerFunc(func(ctx context.Context, msg string, keyvals ...interface{}) {
		p.T().Logf("Instrumentation received message: %s - %+v", msg, keyvals)
		mc <- fmt.Sprintf("%s - %+v", msg, keyvals)
	})

	deets.UseInstrumentedDriver = true
	deets.InstrumentedDriverOptions = []instrumentedsql.Opt{instrumentedsql.WithLogger(checker)}

	c, err := NewConnection(&deets)
	r.NoError(err)
	r.NoError(c.Open())

	err = c.WithContext(context.TODO()).RawQuery(query, 1).Exec()
	r.NoError(err)

	wg.Wait()
}

func (s *PostgreSQLSuite) Test_Instrumentation() {
	testInstrumentedDriver(&s.Suite)
}

func (s *MySQLSuite) Test_Instrumentation() {
	testInstrumentedDriver(&s.Suite)
}

func (s *SQLiteSuite) Test_Instrumentation() {
	testInstrumentedDriver(&s.Suite)
}

func (s *CockroachSuite) Test_Instrumentation() {
	testInstrumentedDriver(&s.Suite)
}
