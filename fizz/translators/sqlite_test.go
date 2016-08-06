package translators_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/stretchr/testify/require"
)

var _ fizz.Translator = (*translators.SQLite)(nil)
var sqt = translators.SQLite{}

func Test_SQLite_CreateTable(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE TABLE "users" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"first_name" TEXT NOT NULL,
"last_name" TEXT NOT NULL,
"email" TEXT NOT NULL,
"permissions" text,
"age" integer DEFAULT '40'
);`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
	})
	`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_DropTable(t *testing.T) {
	r := require.New(t)

	ddl := `DROP TABLE "users";`

	res, _ := fizz.AString(`drop_table("users")`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_RenameTable(t *testing.T) {
	r := require.New(t)

	ddl := `ALTER TABLE "users" RENAME TO "people";`

	res, _ := fizz.AString(`rename_table("users", "people")`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_RenameTable_NotEnoughValues(t *testing.T) {
	r := require.New(t)

	_, err := sqt.RenameTable([]fizz.Table{})
	r.Error(err)
}

func Test_SQLite_AddColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "mytable" ADD COLUMN "mycolumn" TEXT NOT NULL DEFAULT 'foo';`

	res, _ := fizz.AString(`add_column("mytable", "mycolumn", "string", {"default": "foo", "size": 50})`, sqt)

	r.Equal(ddl, res)
}

func Test_SQLite_DropColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "table_name" DROP COLUMN "column_name";`

	res, _ := fizz.AString(`drop_column("table_name", "column_name")`, sqt)

	r.Equal(ddl, res)
}

func Test_SQLite_RenameColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "table_name" RENAME COLUMN "old_column" TO "new_column";`

	res, _ := fizz.AString(`rename_column("table_name", "old_column", "new_column")`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_AddIndex(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {})`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_AddIndex_Unique(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE UNIQUE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"unique": true})`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_AddIndex_MultiColumn(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "table_name_col1_col2_col3_idx" ON "table_name" (col1, col2, col3);`

	res, _ := fizz.AString(`add_index("table_name", ["col1", "col2", "col3"], {})`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_AddIndex_CustomName(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "custom_name" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"name": "custom_name"})`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_DropIndex(t *testing.T) {
	r := require.New(t)
	ddl := `DROP INDEX "my_idx";`

	res, _ := fizz.AString(`drop_index("my_idx")`, sqt)
	r.Equal(ddl, res)
}

func Test_SQLite_RenameIndex(t *testing.T) {
	r := require.New(t)

	ddl := `ALTER INDEX "old_ix" RENAME TO "new_ix";`

	res, _ := fizz.AString(`rename_index("old_ix", "new_ix")`, sqt)
	r.Equal(ddl, res)
}
