package vm

import (
	"os"
	"testing"
)

func TestCallingPluginFunctionNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test/plugin/plugin.go"
	p.go_func("Foo", "!")
	p.go_func("Baz")
	`

	v := initTestVM()
	// We don't test the result here for two reasons:
	// - If it doesn't work it'll returns error or panic
	// - It's hard to test a plugin obj
	v.testEval(t, input, getFilename())
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingPluginFunctionWithReturnValueNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test/plugin/plugin.go"
	p.go_func("Bar")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, "Bar")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingLibFuncFromPluginNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test/plugin/plugin.go"
	p.go_func("ReturnLibName")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, "lib")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestPluginGenerationNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.generate("db") do |p|
	  p.import_pkg("", "database/sql")
	  p.import_pkg("_", "github.com/lib/pq")
	  p.link_function("sql", "Open")
	end

	conn, err = p.go_func("Open", "postgres", "")
	err = conn.go_func("Ping")
	err.nil?
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, true)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func skipPluginTestIfEnvNotSet(t *testing.T) {
	if os.Getenv("NO_RACE_DETECTION") == "" {
		t.Skip("skipping plugin related tests")
	}
}
