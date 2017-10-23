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
		p = Plugin.config("db") do |c|
		  c.import_pkg("", "database/sql")
		  c.link_function("sql", "Open")
		end

		c = p.context
		c.packages.first[:name]
	`, "database/sql"},
		{`
		p = Plugin.config("db") do |c|
		  c.import_pkg("", "database/sql")
		  c.link_function("sql", "Open")
		end

		c = p.context
		c.functions.first[:prefix]
	`, "sql"},
		{`
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
		evaluated := v.testEvalWithRequire(t, tt.input, getFilename(), "plugin")
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
