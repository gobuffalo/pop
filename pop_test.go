package pop

import (
	stdlog "log"
	"os"
	"testing"
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/stretchr/testify/suite"
)

var PDB *Connection

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
	Debug = false
	AddLookupPaths("./")

	dialect := os.Getenv("SODA_DIALECT")

	var err error
	PDB, err = Connect(dialect)
	if err != nil {
		stdlog.Panic(err)
	}
}

func transaction(fn func(tx *Connection)) {
	err := PDB.Rollback(func(tx *Connection) {
		fn(tx)
	})
	if err != nil {
		stdlog.Fatal(err)
	}
}

func ts(s string) string {
	return PDB.Dialect.TranslateSQL(s)
}

type User struct {
	ID           int           `db:"id"`
	UserName     string        `db:"user_name"`
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

// Validate gets run every time you call a "Validate*" (ValidateAndSave, ValidateAndCreate, ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name.String, Name: "Name"},
	), nil
}

type Users []User

type UserAttribute struct {
	ID       int    `db:"id"`
	UserName string `db:"user_name"`
	NickName string `db:"nick_name"`

	User User `json:"user" belongs_to:"user" fk_id:"UserName" primary_id:"UserName"`
}

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

type Taxi struct {
	ID        int       `db:"id"`
	Model     string    `db:"model"`
	UserID    nulls.Int `db:"user_id"`
	Driver    User      `belongs_to:"user" fk_id:"UserID"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate gets run every time you call a "Validate*" (ValidateAndSave, ValidateAndCreate, ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Book) Validate(tx *Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.Description, Name: "Description"},
	), nil
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

type UsersAddressQuery struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	AddressID int       `db:"address_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	UserName  *string `db:"name" json:"user_name"`
	UserEmail *string `db:"email" json:"user_email"`
}

func (UsersAddressQuery) TableName() string {
	return "users_addresses"
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

type Family struct {
	ID        int       `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Family) TableName() string {
	// schema.table_name
	return "family.members"
}

type Enemy struct {
	A string
}

type Song struct {
	ID           uuid.UUID `db:"id"`
	Title        string    `db:"title"`
	UserID       int       `db:"u_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	ComposedByID int       `json:"composed_by_id" db:"composed_by_id"`
	ComposedBy   Composer  `belongs_to:"composer"`
}

type Composer struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Course struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CourseCode struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CourseID  uuid.UUID `json:"course_id" db:"course_id"`
	Course    Course    `json:"-" db:"-"`
	// Course Course `belongs_to:"course"`
}

type ValidatableCar struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var validationLogs = []string{}

func (v *ValidatableCar) Validate(tx *Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "Validate")
	verrs := validate.Validate(&validators.StringIsPresent{Field: v.Name, Name: "Name"})
	return verrs, nil
}

func (v *ValidatableCar) ValidateSave(tx *Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "ValidateSave")
	return nil, nil
}

func (v *ValidatableCar) ValidateUpdate(tx *Connection) (*validate.Errors, error) {
	validationLogs = append(validationLogs, "ValidateUpdate")
	return nil, nil
}

func (v *ValidatableCar) ValidateCreate(tx *Connection) (*validate.Errors, error) {
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

func (u *CallbacksUser) BeforeSave(tx *Connection) error {
	u.BeforeS = "BeforeSave"
	return nil
}

func (u *CallbacksUser) BeforeUpdate(tx *Connection) error {
	u.BeforeU = "BeforeUpdate"
	return nil
}

func (u *CallbacksUser) BeforeCreate(tx *Connection) error {
	u.BeforeC = "BeforeCreate"
	return nil
}

func (u *CallbacksUser) BeforeDestroy(tx *Connection) error {
	u.BeforeD = "BeforeDestroy"
	return nil
}

func (u *CallbacksUser) AfterSave(tx *Connection) error {
	u.AfterS = "AfterSave"
	return nil
}

func (u *CallbacksUser) AfterUpdate(tx *Connection) error {
	u.AfterU = "AfterUpdate"
	return nil
}

func (u *CallbacksUser) AfterCreate(tx *Connection) error {
	u.AfterC = "AfterCreate"
	return nil
}

func (u *CallbacksUser) AfterDestroy(tx *Connection) error {
	u.AfterD = "AfterDestroy"
	return nil
}

func (u *CallbacksUser) AfterFind(tx *Connection) error {
	u.AfterF = "AfterFind"
	return nil
}

type Label struct {
	ID string `db:"id"`
}
