package pop_test

import (
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

var PDB *pop.Connection

type PostgreSQLSuite struct {
	suite.Suite
}

type MySQLSuite struct {
	suite.Suite
}

type SQLiteSuite struct {
	suite.Suite
}

func TestSpecificSuites(t *testing.T) {
	switch os.Getenv("SODA_DIALECT") {
	case "postgres":
		suite.Run(t, &PostgreSQLSuite{})
	case "mysql":
		suite.Run(t, &MySQLSuite{})
	case "sqlite":
		suite.Run(t, &SQLiteSuite{})
	}
}

func init() {
	pop.Debug = false
	pop.AddLookupPaths("./")

	dialect := os.Getenv("SODA_DIALECT")

	var err error
	PDB, err = pop.Connect(dialect)
	if err != nil {
		log.Panic(err)
	}
}

func transaction(fn func(tx *pop.Connection)) {
	err := PDB.Rollback(func(tx *pop.Connection) {
		fn(tx)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ts(s string) string {
	return PDB.Dialect.TranslateSQL(s)
}

type User struct {
	ID           int           `db:"id"`
	Email        string        `db:"email"`
	Name         nulls.String  `db:"name"`
	Alive        nulls.Bool    `db:"alive"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
	BirthDate    nulls.Time    `db:"birth_date"`
	Bio          nulls.String  `db:"bio"`
	Price        nulls.Float64 `db:"price"`
	FullName     nulls.String  `db:"full_name" select:"name as full_name"`
	Books        Books         `has_many:"books" order_by:"title asc"`
	FavoriteSong Song          `has_one:"song" fk_id:"u_id"`
	Houses       Addresses     `many_to_many:"users_addresses"`
}

type Users []User

type Book struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Isbn        string    `db:"isbn"`
	UserID      nulls.Int `db:"user_id"`
	User        User      `belongs_to:"user"`
	Description string    `db:"description"`
	Writers     Writers   `has_many:"writers"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Books []Book

type Writer struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	BookID    int       `db:"book_id"`
	Book      Book      `belongs_to:"book"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Writers []Writer

type Address struct {
	ID          int       `db:"id"`
	Street      string    `db:"street"`
	HouseNumber int       `db:"house_number"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Addresses []Address

type UsersAddress struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	AddressID int       `db:"address_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Friend struct {
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Friend) TableName() string {
	return "good_friends"
}

type Friends []Friend

type Enemy struct {
	A string
}

type Song struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	UserID     int       `db:"u_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ComposerID int       `json:"composer_id" db:"composer_id"`
	ComposedBy Composer  `belongs_to:"composer"`
}

type Composer struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ValidatableCar struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var validationLogs = []string{}

func (v *ValidatableCar) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "Validate")
	verrs := validate.Validate(&validators.StringIsPresent{Field: v.Name, Name: "Name"})
	return verrs, nil
}

func (v *ValidatableCar) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "ValidateSave")
	return nil, nil
}

func (v *ValidatableCar) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "ValidateUpdate")
	return nil, nil
}

func (v *ValidatableCar) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "ValidateCreate")
	return nil, nil
}

type NotValidatableCar struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CallbacksUser struct {
	ID        int       `db:"id"`
	BeforeS   string    `db:"before_s"`
	BeforeC   string    `db:"before_c"`
	BeforeU   string    `db:"before_u"`
	BeforeD   string    `db:"before_d"`
	AfterS    string    `db:"after_s"`
	AfterC    string    `db:"after_c"`
	AfterU    string    `db:"after_u"`
	AfterD    string    `db:"after_d"`
	AfterF    string    `db:"after_f"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CallbacksUsers []CallbacksUser

func (u *CallbacksUser) BeforeSave(tx *pop.Connection) error {
	u.BeforeS = "BeforeSave"
	return nil
}

func (u *CallbacksUser) BeforeUpdate(tx *pop.Connection) error {
	u.BeforeU = "BeforeUpdate"
	return nil
}

func (u *CallbacksUser) BeforeCreate(tx *pop.Connection) error {
	u.BeforeC = "BeforeCreate"
	return nil
}

func (u *CallbacksUser) BeforeDestroy(tx *pop.Connection) error {
	u.BeforeD = "BeforeDestroy"
	return nil
}

func (u *CallbacksUser) AfterSave(tx *pop.Connection) error {
	u.AfterS = "AfterSave"
	return nil
}

func (u *CallbacksUser) AfterUpdate(tx *pop.Connection) error {
	u.AfterU = "AfterUpdate"
	return nil
}

func (u *CallbacksUser) AfterCreate(tx *pop.Connection) error {
	u.AfterC = "AfterCreate"
	return nil
}

func (u *CallbacksUser) AfterDestroy(tx *pop.Connection) error {
	u.AfterD = "AfterDestroy"
	return nil
}

func (u *CallbacksUser) AfterFind(tx *pop.Connection) error {
	u.AfterF = "AfterFind"
	return nil
}
