package validate

import (
	"testing"
	"os"
	"path/filepath"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"strconv"
)

var declrTmp = `package models

import (
	"github.com/gobuffalo/uuid"
	"time"
)`

var structTmp string = `

type %v struct {
	ID        uuid.UUID`+
	"`" +
	`json:"id" db:"id"` +
	"`\n" +
	`CreatedAt time.Time ` +
	"`" +
	`json:"created_at" db:"%v"` +
	"`\n" + `UpdatedAt time.Time ` +
	"`" +
	`json:"updated_at" db:"%v"` +
	"`\n}"

type structTpl struct {
	structName string
	createdAt string
	updatedAt string
	duplicateField string
}


func createModel(fileName string, structs []structTpl)  {
	os.Mkdir("./models", 0755)

	var tmp string = declrTmp

	for _, structTp := range structs {
		tmp = strings.Join([]string{
			tmp,
			strings.Join(
				[]string{
					fmt.Sprintf(
						structTmp,
						structTp.structName,
						structTp.createdAt,
						structTp.updatedAt,
					)},
				"\n"),
		}, "")
	}

	f, _ := os.Create(filepath.Join("models", fileName))
	f.WriteString(tmp)
	f.Close()
}

func Test_testValidate(t *testing.T) {
	r := require.New(t)

	createModel("customer.go", []structTpl{
		{
			"Customer",
			"created_at",
			"updated_at",
			"",
		},
		{
			"Customer1",
			"created_at",
			"updated_at",
			"",
		},
	})
	defer os.RemoveAll("./models")

	m := NewModel()

	errs := m.Validate()
	r.Empty(errs)
}

func Test_testValidateDuplicates(t *testing.T) {
	r := require.New(t)
	structs := []structTpl{
		{
			"Customer",
			"created_at",
			"created_at",
			"created_at",
		},
		{
			"Customer1",
			"created_at",
			"updated_at",
			"updated_at",
		},
		{
			"Customer2",
			"created_at",
			"created_at",
			"created_at",
		},
	}

	createModel("customer.go", structs)
	defer os.RemoveAll("./models")

	m := NewModel()

	errs := m.Validate()

	r.Len(errs, 2)

	for _, ers := range errs {
		for _, structTp := range structs {
			if ers.structName == structTp.structName {
				r.True(ers.duplicate)
				r.Equal(structTp.duplicateField, ers.field)
				r.Equal(structTp.structName, ers.structName)
			}
		}
	}
}

func BenchmarkModel_ValidateNoErrors(b *testing.B) {
	var cnt int = 10000

	for i := 0; i < cnt; i++ {
		structs := []structTpl{{
			"Customer" + strconv.Itoa(i),
			"created_at" + strconv.Itoa(i),
			"updated_at" + strconv.Itoa(i),
			"updated_at" + strconv.Itoa(i),
			},
		}

		createModel("Customer" + strconv.Itoa(i) + ".go", structs)
		defer os.RemoveAll("./models")
	}



	//We don't want to add the struct creation time into the benchmark
	//so we reset the timer
	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		m := NewModel()
		m.Validate()
	}
}

func BenchmarkModel_ValidateWithErrors(b *testing.B) {
	var cnt int = 10000

	for i := 0; i < cnt ; i++ {
		structs := []structTpl{{
			"Customer" + strconv.Itoa(i),
			"updated_at" + strconv.Itoa(i),
			"updated_at" + strconv.Itoa(i),
			"updated_at" + strconv.Itoa(i),
			},
		}

		createModel("Customer" + strconv.Itoa(i) + ".go", structs)
		defer os.RemoveAll("./models")
	}


	//We don't want to add the struct creation time into the benchmark
	//so we reset the timer
	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		m := NewModel()
		m.Validate()
	}
}
