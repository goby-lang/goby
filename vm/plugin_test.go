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

		p = Plugin.new do |c|
		  c.add_pkg("", "database/sql")
		  c.add_func("sql", "Open")
		end

		c = p.context
		c.pkgs.first[1]
	`, "database/sql"},
		{`
		require "plugin"

		p = Plugin.new do |c|
		  c.add_pkg("", "database/sql")
		  c.add_func("sql", "Open")
		end

		c = p.context
		c.funcs.first[0]
	`, "sql"},
		{`
		require "plugin"

		p = Plugin.new do |c|
		  c.add_pkg("", "database/sql")
		  c.add_func("sql", "Open")
		end

		c = p.context
		c.funcs.first[1]
	`, "Open"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
