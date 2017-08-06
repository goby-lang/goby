// +build linux

package vm

import (
	"os"
	"testing"
)

func TestCallingPluginFunction(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "github.com/goby-lang/goby/test_fixtures/import_test/plugin/plugin.go"
	p.send("Foo", "!")
	p.send("Baz")
	`

	v := initTestVM()
	// We don't test the result here for two reasons:
	// - If it doesn't work it'll returns error or panic
	// - It's hard to test a plugin obj
	v.testEval(t, input)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingPluginFunctionWithReturnValue(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "github.com/goby-lang/goby/test_fixtures/import_test/plugin/plugin.go"
	p.send("Bar")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, "Bar")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingLibFuncFromPlugin(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "github.com/goby-lang/goby/test_fixtures/import_test/plugin/plugin.go"
	p.send("ReturnLibName")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, "lib")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestNewPluginUsage(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.generate("db") do |p|
	  p.import_pkg("", "database/sql")
	  p.import_pkg("_", "github.com/lib/pq")
	  p.link_function("sql", "Open")
	end

	conn, err = p.send("Open", "postgres", "")
	err = conn.send("Ping")
	err.is_nil
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, false)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func skipPluginTestIfEnvNotSet(t *testing.T) {
	if os.Getenv("TEST_PLUGIN") == "" {
		t.Skip("skipping plugin related tests")
	}
}
