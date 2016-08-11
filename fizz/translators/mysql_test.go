package translators_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/stretchr/testify/require"
)

var _ fizz.Translator = (*translators.MySQL)(nil)
var myt = translators.NewMySQL("", "")

func init() {
	myconn, err := pop.Connect("mysql")
	if err != nil {
		panic(err.Error())
	}
	deets := myconn.Dialect.Details()
	myt = translators.NewMySQL(myconn.URL(), deets.Database)
}

func Test_MySQL_SchemaMigration(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE TABLE schema_migrations (
version VARCHAR (255) NOT NULL
) ENGINE=InnoDB;
CREATE UNIQUE INDEX version_idx ON schema_migrations (version);`

	res, err := myt.CreateTable(fizz.Table{
		Name: "schema_migrations",
		Columns: []fizz.Column{
			{Name: "version", ColType: "string"},
		},
		Indexes: []fizz.Index{
			{Name: "version_idx", Columns: []string{"version"}, Unique: true},
		},
	})
	r.NoError(err)
	r.Equal(ddl, res)
}

func Test_MySQL_CreateTable(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE TABLE users (
id integer NOT NULL AUTO_INCREMENT,
PRIMARY KEY(id),
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
first_name VARCHAR (255) NOT NULL,
last_name VARCHAR (255) NOT NULL,
email VARCHAR (20) NOT NULL,
permissions text,
age integer DEFAULT 40
) ENGINE=InnoDB;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
	})
	`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_DropTable(t *testing.T) {
	r := require.New(t)

	ddl := `DROP TABLE "users";`

	res, _ := fizz.AString(`drop_table("users")`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_RenameTable(t *testing.T) {
	r := require.New(t)

	ddl := `ALTER TABLE "users" RENAME TO "people";`

	res, _ := fizz.AString(`rename_table("users", "people")`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_RenameTable_NotEnoughValues(t *testing.T) {
	r := require.New(t)

	_, err := myt.RenameTable([]fizz.Table{})
	r.Error(err)
}

func Test_MySQL_AddColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE users ADD COLUMN mycolumn VARCHAR (50) NOT NULL DEFAULT 'foo';`

	res, _ := fizz.AString(`add_column("users", "mycolumn", "string", {"default": "foo", "size": 50})`, myt)

	r.Equal(ddl, res)
}

func Test_MySQL_DropColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE users DROP COLUMN mycolumn;`

	res, _ := fizz.AString(`drop_column("users", "mycolumn")`, myt)

	r.Equal(ddl, res)
}

func Test_MySQL_RenameColumn(t *testing.T) {
	r := require.New(t)
	ddl := `ALTER TABLE users CHANGE email email_address varchar(50) NOT NULL DEFAULT 'foo@example.com';`

	res, _ := fizz.AString(`rename_column("users", "email", "email_address")`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_AddIndex(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX users_email_idx ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {})`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_AddIndex_Unique(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE UNIQUE INDEX users_email_idx ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {"unique": true})`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_AddIndex_MultiColumn(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX users_id_email_idx ON users (id, email);`

	res, _ := fizz.AString(`add_index("users", ["id", "email"], {})`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_AddIndex_CustomName(t *testing.T) {
	r := require.New(t)
	ddl := `CREATE INDEX email_index ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {"name": "email_index"})`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_DropIndex(t *testing.T) {
	r := require.New(t)
	ddl := `DROP INDEX email_idx ON users;`

	res, _ := fizz.AString(`drop_index("users", "email_idx")`, myt)
	r.Equal(ddl, res)
}

func Test_MySQL_RenameIndex(t *testing.T) {
	r := require.New(t)

	ddl := `ALTER TABLE users RENAME INDEX email_idx TO email_address_ix;`

	res, _ := fizz.AString(`rename_index("users", "email_idx", "email_address_ix")`, myt)
	r.Equal(ddl, res)
}
