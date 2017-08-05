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
		c.pkgs.first[:name]
	`, "database/sql"},
		{`
		require "plugin"

		p = Plugin.new do |c|
		  c.add_pkg("", "database/sql")
		  c.add_func("sql", "Open")
		end

		c = p.context
		c.funcs.first[:prefix]
	`, "sql"},
		{`
		require "plugin"

		p = Plugin.new do |c|
		  c.add_pkg("", "database/sql")
		  c.add_func("sql", "Open")
		end

		c = p.context
		c.funcs.first[:name]
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
