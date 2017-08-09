package vm

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

func setupDB(t *testing.T) *sqlx.DB {
	db, _ := sqlx.Open("postgres", "user=postgres dbname=goby_test sslmode=disable")
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id   serial PRIMARY KEY,
		name varchar(40),
		age integer
  	)
`)

	if err != nil {
		t.Fatalf(err.Error())
	}

	return db
}

func cleanTable() {
	db, _ := sqlx.Open("postgres", "user=postgres dbname=goby_test sslmode=disable")
	db.Exec(`DELETE (SELECT * FROM users)`)
}

func TestPGConnectionPing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			require "db"

			db = DB.open("postgres", "user=postgres sslmode=disable")
			db.ping
			`,
			true},
		{`
			require "db"

			db = DB.open("postgres", "user=test sslmode=disable")
			db.ping
			`,
			false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDBExec(t *testing.T) {
	setupDB(t)

	v := initTestVM()

	tests := []struct {
		input    string
		expected interface{}
	}{
		// Insert and query
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			results = db.query("SELECT * FROM users WHERE id = $1", id)
			results.first[:name]
			`,
			"Stan"},
		// Insert and delete
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			id = db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			db.exec("DELETE FROM users WHERE id = $1", id)
			results = db.query("SELECT EXISTS(SELECT * FROM users WHERE id = $1)", id)
			results.first[:exists]
			`,
			false},
		// Insert and update and query
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			id = db.exec("INSERT INTO users (name, age) VALUES ('John', 20)")
			id2 = db.exec("UPDATE users SET age=10 WHERE id = $1", id)
			# See if update returns usable id, too
			results = db.query("SELECT * FROM users WHERE id = $1", id2)
			results.first[:age]
			`,
			10},
	}

	for i, tt := range tests {
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
		cleanTable()
	}
}
