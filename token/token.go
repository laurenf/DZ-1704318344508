package token

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

// Token is a simple struct for the db and default limit
type Token struct {
	db    *sqlx.DB
	limit int
}

// NewToken creates a Token instance
func NewToken() *Token {
	db := sqlx.MustConnect("postgres", "postgres://foo:foo@localhost:5433/foo?sslmode=disable")
	return &Token{
		db:    db,
		limit: 2,
	}
}

func foo(in bool) bool {
	return in
}

// Foo is the best data type
type Foo struct {
	ID        int       `db:"id" json:"id"`
	Data      string    `db:"data" json:"data"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// FooPlusToken is for serializing Foo and the token
type FooPlusToken struct {
	Foo       []Foo  `json:"foo"`
	LastToken string `json:"lastToken"`
	NewToken  string `json:"newToken"`
}

var schema = `
		CREATE TABLE foo (
		id SERIAL PRIMARY KEY,
		data varchar(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
		CREATE INDEX ON foo (updated_at)
`
var answers = []string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

func (t *Token) createDB() {
	t.db.Exec("DROP TABLE foo")
	t.db.MustExec(schema)
}

func (t *Token) destroyDB() {
	t.db.MustExec("DROP TABLE foo")
}

func (t *Token) insertData(data string, delay int) {
	tx := t.db.MustBegin()
	q := `INSERT INTO foo (data) VALUES ($1)`
	if data == "" {
		data = answers[rand.Intn(len(answers))]
	}
	tx.MustExec(q, data)
	time.Sleep(time.Duration(delay) * time.Second)
	tx.Commit()
}

func (t *Token) updateData(id int, data string) {
	t.db.MustExec("UPDATE foo SET data = $1, updated_at = NOW() WHERE id = $2", data, id)
}

func (t *Token) selectData(token string) ([]Foo, string) {
	q := `SELECT * FROM foo
	WHERE (
	  updated_at >= $1
		  OR (updated_at = $1 AND id >= $2)
			)
	AND updated_at < now()
	ORDER BY updated_at asc, id asc
	LIMIT $3`
	foos := []Foo{}
	var id int
	var updatedAt time.Time
	if token != "" {
		id, updatedAt = t.tokenToArgs(token)
	}
	err := t.db.Select(&foos, q, updatedAt, id, t.limit+1)
	if err == sql.ErrNoRows {
		return nil, ""
	}
	if len(foos) <= t.limit {
		return foos, ""
	}
	// only call this if len(foos) > t.limit
	newToken := t.createToken(foos)
	return foos[:t.limit], newToken
}

func (t *Token) tokenToArgs(token string) (int, time.Time) {
	i := strings.Split(token, "_")
	id, err := strconv.Atoi(i[0])
	if err != nil {
		panic(err)
	}
	time, err := time.Parse(time.RFC3339Nano, i[1])
	if err != nil {
		panic(err)
	}
	return id, time
}

func (t *Token) createToken(foo []Foo) string {
	// grab the last element which is where the next
	// request will start
	f := foo[t.limit]
	s := []string{strconv.Itoa(f.ID), f.UpdatedAt.Format(time.RFC3339Nano)}
	newToken := strings.Join(s, "_")
	return newToken
}

func toJSON(foo []Foo, lastToken, newToken string) []byte {
	f := FooPlusToken{
		Foo:       foo,
		LastToken: lastToken,
		NewToken:  newToken,
	}
	bytes, _ := json.Marshal(f)
	return bytes
}
