package pop

import (
	"os"
	"slices"

	"github.com/stretchr/testify/suite"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func testInstrumentedDriver(p *suite.Suite) {
	var (
		queryMySQL = "SELECT 1 FROM DUAL WHERE (1+1)=?"
		query      = "SELECT 1 WHERE (1+1)=?"
		expected   = []string{
			"SELECT 1 FROM DUAL WHERE (1+1)=?",
			"SELECT 1 FROM DUAL WHERE (1+1)=$1",
			"SELECT 1 WHERE (1+1)=?",
			"SELECT 1 WHERE (1+1)=$1",
		}
		r        = p.Require()
		recorder = tracetest.NewSpanRecorder()
		provider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
		tracer   = provider.Tracer("test")
		deets    = *Connections[os.Getenv("SODA_DIALECT")].Dialect.Details()
	)
	deets.TracerProvider = provider
	if os.Getenv("SODA_DIALECT") == "mysql" {
		query = queryMySQL
	}

	c, err := NewConnection(&deets)
	r.NoError(err)
	r.NoError(c.Open())

	ctx, span := tracer.Start(p.T().Context(), "parent")

	err = c.WithContext(ctx).RawQuery(query, 2).Exec()
	r.NoError(err)

	span.End()

	var foundStatement bool
	for _, span := range recorder.Ended() {
		for _, attr := range span.Attributes() {
			if slices.ContainsFunc(expected, func(s string) bool { return attr.Key == "db.statement" && attr.Value.AsString() == s }) {
				foundStatement = true
			}
			r.NotContains(attr.Value.AsString(), "2", "expected attributes to not contain query arguments")
		}
	}
	r.True(foundStatement, "expected to find db.statement attribute with query inside span")
}

func (s *PostgreSQLSuite) TestInstrumentationPostgreSQL() {
	testInstrumentedDriver(&s.Suite)
}

func (s *MySQLSuite) TestInstrumentationMySQL() {
	testInstrumentedDriver(&s.Suite)
}

func (s *SQLiteSuite) TestInstrumentationSQLite() {
	testInstrumentedDriver(&s.Suite)
}

func (s *CockroachSuite) TestInstrumentationCockroachDB() {
	testInstrumentedDriver(&s.Suite)
}
