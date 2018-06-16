package translators_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CockroachSuite struct {
	suite.Suite
}

type PostgreSQLSuite struct {
	suite.Suite
}

type MySQLSuite struct {
	suite.Suite
}

type SQLiteSuite struct {
	suite.Suite
}

type SchemaSuite struct {
	suite.Suite
}

func TestSpecificSuites(t *testing.T) {
	switch os.Getenv("SODA_DIALECT") {
	case "postgres":
		suite.Run(t, &PostgreSQLSuite{})
	case "cockroach":
		suite.Run(t, &CockroachSuite{})
	case "mysql":
		suite.Run(t, &MySQLSuite{})
	case "sqlite":
		suite.Run(t, &SQLiteSuite{})
	}

	suite.Run(t, &SchemaSuite{})
}
