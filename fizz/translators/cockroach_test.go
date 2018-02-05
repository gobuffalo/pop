package translators_test

import (
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

var _ fizz.Translator = (*translators.Cockroach)(nil)

func (p *CockroachSuite) crdbt() *translators.Cockroach {
	ret := translators.NewCockroach("", "pop_test")
	schema := map[string]*fizz.Table{}
	ta := &fizz.Table{Name: "users"}
	ta.Column("testColumn", "type", nil)
	ta.Indexes = append(ta.Indexes, fizz.Index{Name: "testIndex"})
	schema["users"] = ta
	ta = &fizz.Table{Name: "table_name"}
	ta.Column("old_column", "type", nil)
	schema["table_name"] = ta
	ta = &fizz.Table{Name: "mytable"}
	ta.Column("mycolumn", "type", nil)
	schema["mytable"] = ta
	ta = &fizz.Table{Name: "table"}
	ta.Indexes = []fizz.Index{fizz.Index{Name: "old_ix"}}
	schema["table"] = ta
	ta = &fizz.Table{Name: "profiles"}
	schema["profiles"] = ta

	ret.Schema.ReplaceSchema(schema)
	return ret
}

func (p *CockroachSuite) Test_Cockroach_CreateTable() {
	r := p.Require()
	ddl := `CREATE TABLE "users" (
"id" SERIAL PRIMARY KEY,
"first_name" VARCHAR (255) NOT NULL,
"last_name" VARCHAR (255) NOT NULL,
"email" VARCHAR (20) NOT NULL,
"permissions" jsonb,
"age" integer DEFAULT '40',
"company_id" UUID NOT NULL DEFAULT uuid_generate_v1(),
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "jsonb", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
		t.Column("company_id", "uuid", {"default_raw": "uuid_generate_v1()"})
	})
	`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_CreateTable_UUID() {
	r := p.Require()
	ddl := `CREATE TABLE "users" (
"first_name" VARCHAR (255) NOT NULL,
"last_name" VARCHAR (255) NOT NULL,
"email" VARCHAR (20) NOT NULL,
"permissions" jsonb,
"age" integer DEFAULT '40',
"uuid" UUID PRIMARY KEY,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "jsonb", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
		t.Column("uuid", "uuid", {"primary": true})
	})
	`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_CreateTables_WithForeignKeys() {
	r := p.Require()
	ddl := `CREATE TABLE "users" (
"id" SERIAL PRIMARY KEY,
"email" VARCHAR (20) NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL
);COMMIT TRANSACTION;BEGIN TRANSACTION;
CREATE TABLE "profiles" (
"id" SERIAL PRIMARY KEY,
"user_id" INT NOT NULL,
"first_name" VARCHAR (255) NOT NULL,
"last_name" VARCHAR (255) NOT NULL,
"created_at" timestamp NOT NULL,
"updated_at" timestamp NOT NULL,
CONSTRAINT profiles_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id)
);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("id", "INT", {"primary": true})
		t.Column("email", "string", {"size":20})
	})
	create_table("profiles", func(t) {
		t.Column("id", "INT", {"primary": true})
		t.Column("user_id", "INT", {})
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.ForeignKey("user_id", {"users": ["id"]}, {})
	})
	`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_DropTable() {
	r := p.Require()

	ddl := `DROP TABLE "users";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`drop_table("users")`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_RenameTable() {
	r := p.Require()

	ddl := `ALTER TABLE "users" RENAME TO "people";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`rename_table("users", "people")`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_RenameTable_NotEnoughValues() {
	r := p.Require()

	_, err := p.crdbt().RenameTable([]fizz.Table{})
	r.Error(err)
}

func (p *CockroachSuite) Test_Cockroach_ChangeColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "mytable" RENAME COLUMN "mycolumn" TO "_mycolumn_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "mytable" ADD COLUMN "mycolumn" VARCHAR (50) NOT NULL DEFAULT 'foo';COMMIT TRANSACTION;BEGIN TRANSACTION;
UPDATE "mytable" SET "mycolumn" = "_mycolumn_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;
ALTER TABLE "mytable" DROP COLUMN "_mycolumn_tmp";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`change_column("mytable", "mycolumn", "string", {"default": "foo", "size": 50})`, p.crdbt())

	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_AddColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "mytable" ADD COLUMN "mycolumn" VARCHAR (50) NOT NULL DEFAULT 'foo';COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_column("mytable", "mycolumn", "string", {"default": "foo", "size": 50})`, p.crdbt())

	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_DropColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "table_name" DROP COLUMN "column_name";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`drop_column("table_name", "column_name")`, p.crdbt())

	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_RenameColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "table_name" RENAME COLUMN "old_column" TO "new_column";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`rename_column("table_name", "old_column", "new_column")`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_AddIndex() {
	r := p.Require()
	ddl := `CREATE INDEX "table_name_column_name_idx" ON "table_name" (column_name);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {})`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_AddIndex_Unique() {
	r := p.Require()
	ddl := `CREATE UNIQUE INDEX "table_name_column_name_idx" ON "table_name" (column_name);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"unique": true})`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_AddIndex_MultiColumn() {
	r := p.Require()
	ddl := `CREATE INDEX "table_name_col1_col2_col3_idx" ON "table_name" (col1, col2, col3);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_index("table_name", ["col1", "col2", "col3"], {})`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_AddIndex_CustomName() {
	r := p.Require()
	ddl := `CREATE INDEX "custom_name" ON "table_name" (column_name);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"name": "custom_name"})`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_DropIndex() {
	r := p.Require()
	ddl := `DROP INDEX IF EXISTS "my_idx";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`drop_index("users", "my_idx")`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_RenameIndex() {
	r := p.Require()

	ddl := `ALTER INDEX "table"@"old_ix" RENAME TO "new_ix";COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`rename_index("table", "old_ix", "new_ix")`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) buildSchema() translators.Schema {
	schema := map[string]*fizz.Table{}
	ta := &fizz.Table{Name: "testTable"}
	ta.Column("testColumn", "type", nil)
	ta.Indexes = append(ta.Indexes, fizz.Index{Name: "testIndex"})
	schema["testTable"] = ta
	return translators.CreateSchema("name", "url", schema)
}

func (p *CockroachSuite) Test_Cockroach_AddForeignKey() {
	r := p.Require()

	ddl := `ALTER TABLE profiles ADD CONSTRAINT profiles_users_id_fk FOREIGN KEY (user_id) REFERENCES users (id);COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`add_foreign_key("profiles", "user_id", {"users": ["id"]}, {})`, p.crdbt())
	r.Equal(ddl, res)
}

func (p *CockroachSuite) Test_Cockroach_DropForeignKey() {
	r := p.Require()

	ddl := `ALTER TABLE profiles DROP CONSTRAINT  profiles_users_id_fk;COMMIT TRANSACTION;BEGIN TRANSACTION;`

	res, _ := fizz.AString(`drop_foreign_key("profiles", "profiles_users_id_fk", {})`, p.crdbt())
	r.Equal(ddl, res)
}
