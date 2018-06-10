package db

import (
	"testing"

	"github.com/goby-lang/goby/vm"
	"github.com/jmoiron/sqlx"
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

func TestDBPingMethod(t *testing.T) {
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
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestDBCloseMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			require "db"

			db = DB.open("postgres", "user=postgres sslmode=disable")
			first_ping = db.ping # This should be successful

			db.close
			second_ping = db.ping # This should be failed

			first_ping && !second_ping
			`,
			true},
	}

	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestDBRunMethod(t *testing.T) {
	input := `
	require "db"

	db = DB.open("postgres", "user=postgres dbname=goby_test sslmode=disable")
	db.run("create table if not exists test_items (
	  id   serial primary key,
	  title varchar(40)
	)")

	id = db.exec("INSERT INTO test_items (title) VALUES ('Stan')")
	results = db.query("SELECT EXISTS(SELECT * FROM test_items WHERE id = $1)", id)

	db.run("drop table test_items")

	results.first[:exists]
	`

	evaluated := vm.ExecAndReturn(t, input)
	vm.VerifyExpected(t, 0, evaluated, true)
}

func TestDBExecMethod(t *testing.T) {
	setupDB(t)

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
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}
