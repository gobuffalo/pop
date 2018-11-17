package nulls_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	ID         Int64     `json:"id" db:"id"`
	Name       String    `json:"name" db:"name"`
	Alive      Bool      `json:"alive" db:"alive"`
	Price      Float64   `json:"price" db:"price"`
	Birth      Time      `json:"birth" db:"birth"`
	Price32    Float32   `json:"price32" db:"price32"`
	Bytes      ByteSlice `json:"bytes" db:"bytes"`
	IntType    Int       `json:"intType" db:"int_type"`
	Int32Type  Int32     `json:"int32Type" db:"int32_type"`
	UInt32Type UInt32    `json:"uint32Type" db:"uint32_type"`
	UID        UUID      `json:"uid" db:"uid"`
}

const schema = `CREATE TABLE "main"."foos" (
	 "id" integer,
	 "name" text,
	 "alive" integer,
	 "price" float,
	 "birth" timestamp,
	 "price32" float,
	 "bytes"  blob,
	 "int_type" integer,
	 "int32_type" integer,
	 "uint32_type" integer,
	 "uid" uuid
);`

var uid, _ = uuid.NewV4()
var now = time.Now()

func newValidFoo() Foo {
	return Foo{
		ID:         NewInt64(1),
		Name:       NewString("Mark"),
		Alive:      NewBool(true),
		Price:      NewFloat64(9.99),
		Birth:      NewTime(now),
		Price32:    NewFloat32(3.33),
		Bytes:      NewByteSlice([]byte("Byte Slice")),
		IntType:    NewInt(2),
		Int32Type:  NewInt32(3),
		UInt32Type: NewUInt32(5),
		UID:        NewUUID(uid),
	}
}

func Test_TypesMarshalProperly(t *testing.T) {
	t.Parallel()

	a := require.New(t)
	f := newValidFoo()

	ti, _ := json.Marshal(now)
	ba, _ := json.Marshal(f.Bytes)
	jsonString := fmt.Sprintf(`{"id":1,"name":"Mark","alive":true,"price":9.99,"birth":%s,"price32":3.33,"bytes":%s,"intType":2,"int32Type":3,"uint32Type":5,"uid":"%s"}`, ti, ba, uid.String())

	// check marshalling to json works:
	data, _ := json.Marshal(f)
	a.Equal(string(data), jsonString)

	// check unmarshalling from json works:
	f = Foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	a.Equal(f.ID.Int64, int64(1))
	a.Equal(f.Name.String, "Mark")
	a.Equal(f.Alive.Bool, true)
	a.Equal(f.Price.Float64, 9.99)
	a.Equal(f.Birth.Time.Nanosecond(), now.Nanosecond())
	a.Equal(f.Price32.Float32, float32(3.33))
	a.Equal(f.Bytes.ByteSlice, ba)
	a.Equal(f.IntType.Int, 2)
	a.Equal(f.Int32Type.Int32, int32(3))
	a.Equal(f.UInt32Type.UInt32, uint32(5))
	a.Equal(uid.String(), f.UID.UUID.String())

	// check marshalling nulls works:
	f = Foo{}
	jsonString = `{"id":null,"name":null,"alive":null,"price":null,"birth":null,"price32":null,"bytes":null,"intType":null,"int32Type":null,"uint32Type":null,"uid":null}`
	data, _ = json.Marshal(f)
	a.Equal(string(data), jsonString)

	f = Foo{}
	json.NewDecoder(strings.NewReader(jsonString)).Decode(&f)
	a.Equal(f.ID.Int64, int64(0))
	a.False(f.ID.Valid)
	a.Equal(f.Name.String, "")
	a.False(f.Name.Valid)
	a.Equal(f.Alive.Bool, false)
	a.False(f.Alive.Valid)
	a.Equal(f.Price.Float64, float64(0))
	a.False(f.Price.Valid)
	a.Equal(f.Birth.Time, time.Time{})
	a.False(f.Birth.Valid)
	a.Equal(f.Price32.Float32, float32(0))
	a.False(f.Price32.Valid)
	a.Equal(f.Bytes.ByteSlice, []byte(nil))
	a.False(f.Bytes.Valid)
	a.Equal(f.IntType.Int, 0)
	a.False(f.IntType.Valid)
	a.Equal(f.Int32Type.Int32, int32(0))
	a.False(f.Int32Type.Valid)
	a.Equal(f.UInt32Type.UInt32, uint32(0))
	a.False(f.UInt32Type.Valid)
	a.Equal(f.UID.UUID, uuid.Nil)
	a.False(f.UID.Valid)
}

func Test_TypeSaveAndRetrieveProperly(t *testing.T) {
	t.Parallel()

	a := require.New(t)

	initDB(func(db *sqlx.DB) {
		// Test with invalid INSERT query
		tx, err := db.Beginx()
		a.NoError(err)
		_, err = tx.Exec("insert into foos")
		a.Error(err)

		f := Foo{}
		a.Equal(sql.ErrNoRows, tx.Get(&f, "select * from foos"))
		a.False(f.Alive.Valid)
		a.False(f.Birth.Valid)
		a.False(f.ID.Valid)
		a.False(f.Name.Valid)
		a.False(f.Price.Valid)
		a.False(f.Alive.Bool)
		a.False(f.Price32.Valid)
		a.False(f.Bytes.Valid)
		a.False(f.IntType.Valid)
		a.False(f.Int32Type.Valid)
		a.False(f.UInt32Type.Valid)
		a.Equal(f.Birth.Time.UnixNano(), time.Time{}.UnixNano())
		a.Equal(f.ID.Int64, int64(0))
		a.Equal(f.Name.String, "")
		a.Equal(f.Price.Float64, float64(0))
		a.Equal(f.Price32.Float32, float32(0))
		a.Equal(f.Bytes.ByteSlice, []byte(nil))
		a.Equal(f.IntType.Int, 0)
		a.Equal(f.Int32Type.Int32, int32(0))
		a.Equal(f.UInt32Type.UInt32, uint32(0))
		a.NoError(tx.Rollback())

		// Test with valid INSERT query
		tx, err = db.Beginx()
		a.NoError(err)

		f = newValidFoo()
		_, err = tx.NamedExec("INSERT INTO foos (id, name, alive, price, birth, price32, bytes, int_type, int32_type, uint32_type, uid) VALUES (:id, :name, :alive, :price, :birth, :price32, :bytes, :int_type, :int32_type, :uint32_type, :uid)", &f)
		a.NoError(err)
		f = Foo{}
		a.NoError(tx.Get(&f, "select * from foos"))
		a.True(f.Alive.Valid)
		a.True(f.Birth.Valid)
		a.True(f.ID.Valid)
		a.True(f.Name.Valid)
		a.True(f.Price.Valid)
		a.True(f.Alive.Bool)
		a.True(f.Price32.Valid)
		a.True(f.Bytes.Valid)
		a.True(f.IntType.Valid)
		a.True(f.Int32Type.Valid)
		a.True(f.UInt32Type.Valid)
		a.Equal(f.Birth.Time.UnixNano(), now.UnixNano())
		a.Equal(f.ID.Int64, int64(1))
		a.Equal(f.Name.String, "Mark")
		a.Equal(f.Price.Float64, 9.99)
		a.Equal(f.Price32.Float32, float32(3.33))
		a.Equal(f.Bytes.ByteSlice, []byte("Byte Slice"))
		a.Equal(f.IntType.Int, 2)
		a.Equal(f.Int32Type.Int32, int32(3))
		a.Equal(f.UInt32Type.UInt32, uint32(5))
		a.Equal(uid.String(), f.UID.UUID.String())

		a.NoError(tx.Rollback())
	})
}

func initDB(f func(db *sqlx.DB)) {
	os.Remove("./foo.db")
	db, _ := sqlx.Open("sqlite3", "./foo.db")
	db.MustExec(schema)
	f(db)
	os.Remove("./foo.db")
}
