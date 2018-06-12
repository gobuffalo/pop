# POP [![GoDoc](https://godoc.org/github.com/gobuffalo/pop?status.svg)](https://godoc.org/github.com/gobuffalo/pop) [![Build Status](https://travis-ci.org/gobuffalo/pop.svg)](https://travis-ci.org/gobuffalo/pop)

## A Tasty Treat For All Your Database Needs

So what does Pop do exactly? Well, it wraps the absolutely amazing [https://github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) library. It cleans up some of the common patterns and work flows usually associated with dealing with databases in Go.

Pop makes it easy to do CRUD operations, run migrations, and build/execute queries. Is Pop an ORM? I'll leave that up to you, the reader, to decide.

Pop, by default, follows conventions that were influenced by the [ActiveRecord](http://www.rubyonrails.org) Ruby gem. What does this mean?

* Tables must have an "id" column and a corresponding "ID" field on the `struct` being used.
* If there is a `timestamp` column named `created_at`, and a `CreatedAt time.Time` attribute on the `struct`, it will be set with the current time when the record is created.
* If there is a `timestamp` column named `updated_at`, and a `UpdatedAt time.Time` attribute on the `struct`, it will be set with the current time when the record is updated.
* Default database table names are lowercase, plural, and underscored versions of the `struct` name. Examples: User{} is "users", FooBar{} is "foo_bars", etc...

## Supported Databases

* PostgreSQL (>= 9.3)
* MySQL (>= 5.7)
* SQLite (>= 3.x)
* CockroachDB (>= 1.1.1)

## Connecting to Databases

Pop is easily configured using a YAML file. The configuration file should be stored in `config/database.yml` or `database.yml`.

#### Example Configuration File

```yaml
development:
  dialect: "postgres"
  database: "your_db_development"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"

test:
  dialect: "mysql"
  database: "your_db_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

staging:
  dialect: "sqlite3"
  database: "./staging.sqlite"

production:
  dialect: "postgres"
  url: {{ env "DATABASE_URL" }}
```

Note that the `database.yml` file is also a Go template, so you can use Go template syntax. There are two special functions that are included, `env` and `envOr`.

* `env` - This function will look for the named environment variable and insert it into your file. This is useful for configuring production databases without having to store secret information in your repository. `{{ env "DATABASE_URL" }}`
* `envOr` - This function will look for the named environment variable and use it. If the variable can not be found a default value will be used. `{{ envOr "MYSQL_HOST" "localhost" }}`

You can generate a default configuration file using the `init` command:

```
$ soda g config
```

The default will generate a `database.yml` file in the current directory for a PostgreSQL database. You can override the type of database using the `-t` flag and passing in any of the supported database types: `postgres`, `cockroach`, `mysql`, or `sqlite3`.

CockroachDB currently works best if you DO NOT use a URL and instead define each key item. Because CockroachDB more or less uses the same driver as postgres you have the same configuration options for both. In production you will also want to make sure you are using a [secure cluster](https://www.cockroachlabs.com/docs/stable/manual-deployment.html) and have set all the needed [connection parameters](https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters) for said secure connection. If you do not set the sslmode or set it to `disable` this will put dump and load commands into `--insecure` mode.

### In Your Code

Once you have a configuration file defined you can easily connect to one of these connections in your application.

```go
db, err := pop.Connect("development")
if err != nil {
  log.Panic(err)
}
```

Now that you have your connection to the database you can start executing queries against it.

## CLI Support

Pop features CLI support via the `soda` command for the following operations:

```bash
Available Commands:
  create      Creates databases for you
  drop        Drops databases for you
  generate
  help        Help about any command
  migrate     Runs migrations against your database.
  schema      Tools for working with your database schema
```

### Installing CLI Support

**Without** sqlite 3 support:

```bash
$ go get github.com/gobuffalo/pop/...
$ go install github.com/gobuffalo/pop/soda
```

**With** sqlite 3 support:

```bash
$ go get -u -v -tags sqlite github.com/gobuffalo/pop/...
$ go install -tags sqlite github.com/gobuffalo/pop/soda
```

If you're not building your code with `buffalo build`, you'll also have to pass `-tags sqlite` to `go build` when building your program.

## Managing Databases

### Creating Databases

Assuming you defined a configuration file like that described in the above section you can automatically create those databases using the `soda` command:

#### Create All Databases

```bash
$ soda create -a
```

#### Create a Specific Database

```bash
$ soda create -e development
```

### Dropping Databases

Assuming you defined a configuration file like that described in the above section you can automatically drop those databases using the `soda` command:

#### Drop All Databases

```bash
$ soda drop -a
```

#### Drop a Specific Database

```bash
$ soda drop -e development
```

## Models

The `soda` command supports the generation of models.

A full list of commands available for model generation can be found by asking for help:

```bash
$ soda generate help
```

### Generating Models

The `soda` command will generate Go models and, optionally, the associated migrations for you.

```bash
$ soda generate model user name:text email:text
```

Running this command will generate the following files:

```text
models/user.go
models/user_test.go
migrations/20170115024143_create_users.up.fizz
migrations/20170115024143_create_users.down.fizz
```

The `models/user.go` file contains a structure named `User` with fields `ID`, `CreatedAt`, `UpdatedAt`, `Name`, and `Email`. The first three correspond to the columns commonly found in ActiveRecord models as mentioned before, and the last two correspond to the additional fields specified on the command line. The known types are:

* `text` (`string` in Go)
* `blob` (`[]byte` in Go)
* `time` or `timestamp` (`time.Time`)
* `nulls.Text` (`nulls.String`) which corresponds to a nullifyable string, which can be distinguished from an empty string
* `uuid` (`uuid.UUID`)
* Other types are passed thru and are used as [Fizz](./fizz/README.md) types.

The `models/user_test.go` contains tests for the User model and they must be implemented by you.

The other two files correspond to the migrations as explained below. By default, it generates `.fizz` files but you can also generate `.sql` files by adding the flag `--migration-type sql` to the command. Be aware, that you will need to specify the appropriate environment, because `.sql` files are for specific databases.

```bash
$ soda generate model user name:text email:text --migration-type sql -e development
```

If `development`is associated with a `postgresql`configuration, running this command will generate the following files:

```text
models/user.go
models/user_test.go
migrations/20170115024143_create_users.postgres.up.sql
migrations/20170115024143_create_users.postgres.down.sql
```

## Migrations

The `soda` command supports the creation and running of migrations.

A full list of commands available for migration can be found by asking for help:

```bash
$ soda migrate --help
```

### Generating Migrations

The `soda` command will generate SQL migrations (both the up and down) files for you.

```bash
$ soda generate fizz name_of_migration
```

Running this command will generate the following files:

```text
./migrations/20160815134952_name_of_migration.up.fizz
./migrations/20160815134952_name_of_migration.down.fizz
```

The generated files are `fizz` files. [Fizz](./fizz/README.md) lets you use a common DSL for generating migrations. This means the same `.fizz` file can be run against any of the supported dialects of Pop! Find out more about [Fizz](./fizz/README.md)

If you want to generate old fashion `.sql` files you can use the `-t` flag for that:

```bash
$ soda generate sql name_of_migration
```

Running this command will generate the following files:

```text
./migrations/20160815134952_name_of_migration.up.sql
./migrations/20160815134952_name_of_migration.down.sql
```

The `soda migrate` command supports both `.fizz` and `.sql` files, so you can mix and match them to suit your needs.

### Running Migrations

The `soda` command will run the migrations using the following command:

```bash
$ soda migrate up
```

Migrations will be run in sequential order.

Migrations can also be run in reverse to rollback the schema.

```bash
$ soda migrate down
```

The previously run migrations will be kept track of in a table named `schema_migration` in the database. The table name can be configured by setting `migration_table_name` of the configuration options. The example below will use `migrations` as the table name.

```yaml
development:
  dialect: "postgres"
  url: "your_db_development"
  options:
    migration_table_name: migrations
```

## Querying

### Find By ID

```go
user := User{}
err := db.Find(&user, id)
```

### Find All

```go
users := []User{}
err := db.All(&users)
err = db.Where("id in (?)", 1, 2, 3).All(&users)
```

#### Find Last

```go
// Last() orders by created_at
user := models.User{}
err := tx.Last(&user)
```

### Find Where

```go
users := []models.User{}
query := db.Where("id = 1").Where("name = 'Mark'")
err := query.All(&users)

err = tx.Where("id in (?)", 1, 2, 3).All(&users)
```

#### Using `in` Queries

```go
err = db.Where("id in (?)", 1, 2, 3).All(&users)
err = db.Where("id in (?)", 1, 2, 3).Where("foo = ?", "bar").All(&users)
```

Unfortunately, for a variety of reasons you can't use an `and` query in the same `Where` call as an `in` query.

```go
// does not work:
err = db.Where("id in (?) and foo = ?", 1, 2, 3, "bar").All(&users)
// works:
err = db.Where("id in (?)", 1, 2, 3).Where("foo = ?", "bar").All(&users)
```

### Select specific columns
`Select` allows you to load specific columns from a table. Useful when you don't want all columns from a table to be loaded in a query.
```go
err = db.Select("name").All(&users)
// SELECT name FROM users

err = db.Select("max(age)").All(&users)
// SELECT max(age) FROM users

err = db.Select("age", "name").All(&users)
// SELECT age, name FROM users
```

### Join Query

```go
// page: page number
// perpage: limit
roles := []models.UserRole{}
query := models.DB.LeftJoin("roles", "roles.id=user_roles.role_id").
  LeftJoin("users u", "u.id=user_roles.user_id").
  Where(`roles.name like ?`, name).Paginate(page, perpage)

count, _ := query.Count(models.UserRole{})
count, _ := query.CountByField(models.UserRole{}, "*")
sql, args := query.ToSQL(&pop.Model{Value: models.UserRole{}}, "user_roles.*",
  "roles.name as role_name", "u.first_name", "u.last_name")
err := models.DB.RawQuery(sql, args...).All(&roles)
```

## Creating New Records

```go
// Create one record.
user := models.User{}
user.Name = "Mark"
err := tx.Create(&user)

// Create many records.
users := models.Users{
  {Name:"Mark"},
  {Name: "Larry"},
}

err := tx.Create(&users)
```

## Saving Records

The `Save` method will attempt to create the record if the `ID` is empty. If there is an `ID` set it will attempt to update the record with that ID in the database.

```go
// Save one record.
user := models.User{}
user.Name = "Mark"
err := tx.Save(&user)

// Save many records.
users := models.Users{
  {Name:"Mark"},
  {Name: "Larry"},
}

err := tx.Save(&users)
```

## Updating Records

```go
// Update one record.
user := models.User{}
user.Name = "Mark"
err := tx.Create(&user)

user.Name = "Mark Bates"
err = tx.Update(&user)

// Update many records.
users := models.Users{
  {Name:"Mark"},
  {Name: "Larry"},
}

err := tx.Create(&users)

users[0].Name = "Mark Bates"
users[1].Name = "Larry Morales"
err := tx.Update(&users)
```

## Destroy

```go
// Destroy one record.
user := models.User{}
user.Name = "Mark"
err := tx.Create(&user)

err = tx.Destroy(&user)

// Destroy many records.
users := models.Users{
  {Name:"Mark"},
  {Name: "Larry"},
}
err := tx.Create(&users)

err = tx.Destroy(&users)
```

## Eager Loading

Pop allows you to perform an eager loading for associations defined in a model. By using `pop.Connection.Eager()` function plus some fields tags predefined in your model you can extract associated data from a model.

```go
type User struct {
  ID           uuid.UUID
  Email        string
  Password     string
  Books        Books     `has_many:"books" order_by:"title asc"`
  FavoriteSong Song      `has_one:"song" fk_id:"u_id"`
  Houses       Addresses `many_to_many:"users_addresses"`
}
```

```go
type Book struct {
  ID      uuid.UUID
  Title   string
  Isbn    string
  User    User        `belongs_to:"user"`
  UserID  uuid.UUID
  Writers Writers     `has_many:"writers"`
}
```

```go
type Writer struct {
   ID     uuid.UUID   `db:"id"`
   Name   string      `db:"name"`
   BookID uuid.UUID   `db:"book_id"`
   Book   Book        `belongs_to:"book"`
}
```

```go
type Song struct {
  ID      uuid.UUID
  Title   string
  UserID  uuid.UUID   `db:"u_id"`
}
```

```go
type Address struct {
  ID           uuid.UUID
  Street       string
  HouseNumber  int
}

type Addresses []Address
```

* **has_many**: will load all records from the `books` table that have a column named `user_id` or the column specified with **fk_id** that matches the `User.ID` value.

* **belongs_to**: will load a record from `users` table that have a column named `id` that matches with `Book.UserID` value.

* **has_one**: will load a record from the `songs` table that have a column named `user_id` or the column specified with **fk_id** that matches the `User.ID` value.

* **many_to_many**: will load all records from the `addresses` table through the table `users_addresses`. Table `users_addresses` MUST define `address_id`  and `user_id` columns to match `User.ID` and `Address.ID` values. You can also define a **fk_id** tag that will be used in the target association i.e `addresses` table.

* **fk_id**: defines the column name in the target association that matches model `ID`. In the example above `Song` has a column named `u_id` that represents `id` of `users` table. When loading `FavoriteSong`, `u_id` will be used instead of `user_id`.

* **order_by**: used in `has_many` and `many_to_many` to indicate the order for the association when loading. The format to use is  `order_by:"<column_name> <asc | desc>"`


```go
u := Users{}
// preload all associations for user with name 'Mark', i.e Books, Houses and FavoriteSong
err := tx.Eager().Where("name = 'Mark'").All(&u)
// preload only Books association for user with name 'Mark'.
err  = tx.Eager("Books").Where("name = 'Mark'").All(&u)
```

### Eager Loading Nested Associations

Pop allows you to eager loading nested associations by using `.` character to concatenate them. Take a look at the example bellow.

```go
// will load all Books for u and for every Book will load the user which will be the same as u.
tx.Eager("Books.User").First(&u)
```

```go
// will load all Books for u and for every Book will load all Writers.
 tx.Eager("Books.Writers").First(&u)
```

```go
// will load all Books for u and for every Book will load all Writers and for every writer will load the Book association.
tx.Eager("Books.Writers.Book").First(&u)
```

```go
// will load all Books for u and for every Book will load all Writers. And Also it will load the favorite song for user.
tx.Eager("Books.Writers").Eager("FavoriteSong").First(&u)
```

### Eager Creation

Pop allows you to create models and their associations in one step. You no longer need to create every association separately anymore. Pop will even create join table records for `many_to_many` associations.

Assuming the following pieces of pseudo-code:

```go
user := User{
  Name: "Mark Bates",
  Books: Books{{Title: "Pop Book", Description: "Pop Book", Isbn: "PB1"}},
  FavoriteSong: Song{Title: "Don't know the title"},
  Houses: Addresses{
    Address{HouseNumber: 1, Street: "Golang"},
  },
}
```

```go
err := tx.Eager().Create(&user)
```

1. It will notice `Books` is a `has_many` association and it will realize that to actually store every book it will need to get the `User ID` first. So, it proceeds to store first `User` data so it can retrieve an **ID** and then use that ID to fill `UserID` field in every `Book` in `Books`. Later it stores all books in database.

2. `FavoriteSong` is a `has_one` association and it uses same logic described in `has_many` association. Since `User` data was previously saved before creating all books, it already knows that `User` got an `ID` so it fills its `UserID` field with that value and `FavoriteSong` is then stored in database.

3. `Houses` in this example is a `many_to_many` relationship and it will have to deal with two tables in this case: `users` and `addresses`. It will need to store all addresses first in `addresses` table before save them in the many to many table. Because `User` was already stored, it already have an `ID`. * This is a special case to deal with, since this behavior is different to all other associations, it is solved by implementing the `AssociationCreatableStatement` interface, all other associations implement by default `AssociationCreatable` interface.

For a `belongs_to` association like shown in the example below, it will need first to create `User` to retrieve **ID** value and then fill its `UserID` field before be saved in database.

```go
book := Book{
   Title:      "Pop Book",
   Description: "Pop Book",
   Isbn:        "PB1",
   User: User{
        Name: nulls.NewString("Larry"),
   },
}
```

```go
tx.Eager().Create(&book)
```

All these cases are assuming that none of models and associations has previously been saved in database.

## Callbacks

Pop provides a means to execute code before and after database operations. This is done by defining specific methods on your models. For example, to hash a user password you may want to define the following method:

```go
type User struct {
  ID       uuid.UUID
  Email    string
  Password string
}

func (u *User) BeforeSave(tx *pop.Connection) error {
  hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
  if err != nil {
    return errors.WithStack(err)
  }

  u.Password = string(hash)

  return nil
}
```

### Available Callbacks:

* BeforeSave
* BeforeCreate
* BeforeUpdate
* BeforeDestroy
* AfterSave
* AfterCreate
* AfterUpdate
* AfterDestroy
* AfterFind

## Further Reading

[The Unofficial pop Book:](https://andrew-sledge.gitbooks.io/the-unofficial-pop-book/content/) a gentle introduction to new users.
