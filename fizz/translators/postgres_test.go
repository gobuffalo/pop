package translators_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/stretchr/testify/require"
)

func Test_Postgres_CreateTable(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE TABLE IF NOT EXISTS "users" (
"id" SERIAL PRIMARY KEY,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
"first_name" VARCHAR (255) NOT NULL,
"last_name" VARCHAR (255) NOT NULL,
"email" VARCHAR (20) NOT NULL,
"permissions" jsonb,
"age" integer
);`

	bub, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "jsonb", {"null": true})
		t.Column("age", "integer", {"null": true})
	})
	`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.CreateTable(tl)
	r.Equal(ddl, res)

}

func Test_Postgres_DropTable(t *testing.T) {
	r := require.New(t)

	ddl := `DROP TABLE IF EXISTS "users";`

	bub, _ := fizz.AString(`drop_table("users")`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.DropTable(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_RenameTable(t *testing.T) {
	r := require.New(t)

	ddl := `ALTER TABLE "users" RENAME TO "people";`

	bub, _ := fizz.AString(`rename_table("users", "people")`)

	b := bub.Bubbles[0]
	tl := b.Data.([]fizz.Table)

	p := translators.Postgres{}
	res, _ := p.RenameTable(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_RenameTable_NotEnoughValues(t *testing.T) {
	r := require.New(t)

	p := translators.Postgres{}
	_, err := p.RenameTable([]fizz.Table{})
	r.Error(err)
}

func Test_Postgres_AddColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "mytable" ADD COLUMN "mycolumn" VARCHAR (50) NOT NULL DEFAULT 'foo';`

	bub, _ := fizz.AString(`add_column("mytable", "mycolumn", "string", {"default": "foo", "size": 50})`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.AddColumn(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_DropColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "table_name" DROP COLUMN "column_name";`

	bub, _ := fizz.AString(`drop_column("table_name", "column_name")`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.DropColumn(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_RenameColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE "table_name" RENAME COLUMN "old_column" TO "new_column";`

	bub, _ := fizz.AString(`rename_column("table_name", "old_column", "new_column")`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.RenameColumn(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_AddIndex(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	bub, _ := fizz.AString(`add_index("table_name", "column_name", {})`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.AddIndex(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_AddIndex_Unique(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE UNIQUE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	bub, _ := fizz.AString(`add_index("table_name", "column_name", {"unique": true})`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.AddIndex(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_AddIndex_MultiColumn(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "table_name_col1_col2_col3_idx" ON "table_name" (col1, col2, col3);`

	bub, _ := fizz.AString(`add_index("table_name", ["col1", "col2", "col3"], {})`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.AddIndex(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_AddIndex_CustomName(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX "custom_name" ON "table_name" (column_name);`

	bub, _ := fizz.AString(`add_index("table_name", "column_name", {"name": "custom_name"})`)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	p := translators.Postgres{}
	res, _ := p.AddIndex(tl)
	r.Equal(ddl, res)
}

func Test_Postgres_DropIndex(t *testing.T) {
	r := require.New(t)
	ddl := `DROP INDEX IF EXISTS "my_idx";`

	bub, _ := fizz.AString(`drop_index("my_idx")`)
	b := bub.Bubbles[0]

	tl := b.Data.(fizz.Index)

	p := translators.Postgres{}
	res, _ := p.DropIndex(tl)
	r.Equal(ddl, res)
}
