package pop

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

type PostgreSQL struct {
	translateCache    map[string]string
	mu                sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (p *PostgreSQL) Details() *ConnectionDetails {
	return p.ConnectionDetails
}

func (p *PostgreSQL) Create(store Store, model *Model, cols Columns) error {
	cols.Remove("id")
	id := struct {
		ID int `db:"id"`
	}{}
	w := cols.Writeable()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) returning id", model.TableName(), w.String(), w.SymbolizedString())
	Log(query)
	stmt, err := store.PrepareNamed(query)
	if err != nil {
		return err
	}
	err = stmt.Get(&id, model.Value)
	if err == nil {
		model.setID(id.ID)
	}
	return err
}

func (p *PostgreSQL) Update(store Store, model *Model, cols Columns) error {
	return genericUpdate(store, model, cols)
}

func (p *PostgreSQL) Destroy(store Store, model *Model) error {
	return genericDestroy(store, model)
}

func (p *PostgreSQL) SelectOne(store Store, model *Model, query Query) error {
	return genericSelectOne(store, model, query)
}

func (p *PostgreSQL) SelectMany(store Store, models *Model, query Query) error {
	return genericSelectMany(store, models, query)
}

func (p *PostgreSQL) CreateDB() error {
	// createdb -h db -p 5432 -U postgres enterprise_development
	deets := p.ConnectionDetails
	cmd := exec.Command("createdb", "-e", "-h", deets.Host, "-p", deets.Port, "-U", deets.User, deets.Database)
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
}

func (p *PostgreSQL) DropDB() error {
	deets := p.ConnectionDetails
	cmd := exec.Command("dropdb", "-e", "-h", deets.Host, "-p", deets.Port, "-U", deets.User, deets.Database)
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
}

func (m *PostgreSQL) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}

	s := "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *PostgreSQL) MigrationURL() string {
	return m.URL()
}

func (p *PostgreSQL) TranslateSQL(sql string) string {
	defer p.mu.Unlock()
	p.mu.Lock()

	if csql, ok := p.translateCache[sql]; ok {
		return csql
	}
	curr := 1
	out := make([]byte, 0, len(sql))
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			str := "$" + strconv.Itoa(curr)
			for _, char := range str {
				out = append(out, byte(char))
			}
			curr += 1
		} else {
			out = append(out, sql[i])
		}
	}
	csql := string(out)
	p.translateCache[sql] = csql
	return csql
}

func (p *PostgreSQL) FizzTranslator() fizz.Translator {
	return translators.Postgres{}
}

func NewPostgreSQL(deets *ConnectionDetails) Dialect {
	deets.Parse("5432")
	cd := &PostgreSQL{
		ConnectionDetails: deets,
		translateCache:    map[string]string{},
		mu:                sync.Mutex{},
	}
	return cd
}
