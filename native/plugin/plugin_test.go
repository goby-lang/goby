package plugin

import (
	"testing"

	"github.com/goby-lang/goby/vm"
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
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}
