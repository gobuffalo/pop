package translators_test

import (
	"github.com/gobuffalo/pop/fizz"
	"github.com/gobuffalo/pop/fizz/translators"
)

func (s *SchemaSuite) buildSchema() translators.Schema {
	schema := map[string]*fizz.Table{}
	ta := &fizz.Table{Name: "testTable"}
	ta.Column("testColumn", "type", nil)
	ta.Indexes = append(ta.Indexes, fizz.Index{Name: "testIndex"})
	schema["testTable"] = ta
	return translators.CreateSchema("name", "url", schema)
}

func (s *SchemaSuite) Test_Schema_TableInfo() {
	r := s.Require()
	ts := s.buildSchema()
	t, err := ts.TableInfo("testTable")
	r.NoError(err)
	r.Equal("testTable", t.Name)
}

func (s *SchemaSuite) Test_Schema_ColumnInfo() {
	r := s.Require()
	ts := s.buildSchema()
	c, err := ts.ColumnInfo("testTable", "testCOLUMN")
	r.NoError(err)
	r.Equal("testColumn", c.Name)
}

func (s *SchemaSuite) Test_Schema_IndexInfo() {
	r := s.Require()
	ts := s.buildSchema()
	c, err := ts.IndexInfo("testTable", "testindEX")
	r.NoError(err)
	r.Equal("testIndex", c.Name)
}
