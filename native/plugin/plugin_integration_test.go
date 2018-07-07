package plugin

import (
	"os"
	"testing"

	"github.com/goby-lang/goby/vm"
)

func TestCallingPluginFunctionNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test.VM()/plugin.go"
	p.go_func("Foo", "!")
	p.go_func("Baz")
	`

	// We don't test the result here for two reasons:
	// - If it doesn't work it'll returns error or panic
	// - It's hard to test a plugin obj
	vm.ExecAndReturn(t, input)
}

func TestCallingPluginFunctionWithReturnValueNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test.VM()/plugin.go"
	p.go_func("Bar")
	`

	evaluated := vm.ExecAndReturn(t, input)
	vm.VerifyExpected(t, 0, evaluated, "Bar")
}

func TestCallingLibFuncFromPluginNoRaceDetection(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	require "plugin"

	p = Plugin.use "../test_fixtures/import_test.VM()/plugin.go"
	p.go_func("ReturnLibName")
	`

	evaluated := vm.ExecAndReturn(t, input)
	vm.VerifyExpected(t, 0, evaluated, "lib")
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
	!err.nil? && err.go_func("Error").is_a?(String)
	`

	evaluated := vm.ExecAndReturn(t, input)
	vm.VerifyExpected(t, 0, evaluated, true)
}

func skipPluginTestIfEnvNotSet(t *testing.T) {
	if os.Getenv("NO_RACE_DETECTION") == "" {
		t.Skip("skipping plugin related tests")
	}
}
