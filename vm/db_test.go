package vm

import "testing"

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

func TestDBQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			require "db"

			db = DB.open("postgres", "user=postgres sslmode=disable")
			results = db.query("SELECT datname FROM pg_database WHERE datname='postgres'")
			results.first[:datname]
			`,
			"postgres"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
