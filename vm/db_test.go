package vm

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

var db *sqlx.DB

func init() {
	setupDB()
}

func setupDB() {
	db, _ = sqlx.Open("postgres", "user=postgres dbname=goby_test sslmode=disable")
	db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	name varchar(40),
	age integer
)`)
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
	v := initTestVM()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			results = db.query("SELECT EXISTS(SELECT * FROM users)")
			results.count
			`,
			1},
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			db.exec("INSERT INTO users (name, age) VALUES ('Stan', 23)")
			db.exec("UPDATE users SET age=10 WHERE name='Stan'")
			results = db.query("SELECT * FROM users")
			results.last[:age]
			`,
			10},
	}

	for i, tt := range tests {
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDBQuery(t *testing.T) {
	db.Exec("INSERT INTO users (name, age) VALUES ('Goby', 0)")
	v := initTestVM()

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			results = db.query("SELECT datname FROM pg_database WHERE datname='postgres'")
			results.last[:datname]
			`,
			"postgres"},
		{`
			require "db"

			db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
			results = db.query("SELECT * FROM users WHERE name='Goby'")
			results.last[:age]
			`,
			0},
	}

	for i, tt := range tests {
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
