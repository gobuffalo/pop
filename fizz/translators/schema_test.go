package translators_test

import (
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

func (s *SchemaSuite) Test_Schema_TableInfo() {
	r := s.Require()
	schema := map[string]*fizz.Table{}
	ta := &fizz.Table{Name: "testTable"}
	ta.Column("testColumn", "type", nil)
	schema["testTable"] = ta
	ts := translators.CreateSchema("name", "url", schema)
	t, err := ts.TableInfo("testTable")
	r.NoError(err)
	r.Equal("testTable", t.Name)
}

func (s *SchemaSuite) Test_Schema_ColumnInfo() {
	r := s.Require()
	schema := map[string]*fizz.Table{}
	ta := &fizz.Table{Name: "testTable"}
	ta.Column("testColumn", "type", nil)
	schema["testTable"] = ta
	ts := translators.CreateSchema("name", "url", schema)
	c, err := ts.ColumnInfo("testTable", "testCOLUMN")
	r.NoError(err)
	r.Equal("testColumn", c.Name)
}
