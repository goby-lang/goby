package vm

import (
	"testing"
)

func TestPluginInitialization(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "plugin"

		p = Plugin.config("db") do |c|
		  c.import_pkg("", "database/sql")
		  c.link_function("sql", "Open")
		end

		c = p.context
		c.packages.first[:name]
	`, "database/sql"},
		{`
		require "plugin"

		p = Plugin.config("db") do |c|
		  c.import_pkg("", "database/sql")
		  c.link_function("sql", "Open")
		end

		c = p.context
		c.functions.first[:prefix]
	`, "sql"},
		{`
		require "plugin"

		p = Plugin.config("db") do |c|
		  c.import_pkg("", "database/sql")
		  c.link_function("sql", "Open")
		end

		c = p.context
		c.functions.first[:name]
	`, "Open"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
