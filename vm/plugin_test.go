// +build linux

package vm

import (
	"os"
	"testing"
)

func TestCallingPluginFunction(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/plugin.go"
	p.send("Foo", "!")
	p.send("Baz")
	`

	vm := initTestVM()
	// We don't test the result here for two reasons:
	// - If it doesn't work it'll returns error or panic
	// - It's hard to test a plugin obj
	vm.testEval(t, input)
	vm.checkCFP(t, 0, 0)
}

func TestCallingPluginFunctionWithReturnValue(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/plugin.go"
	result = p.send("NewBar", "xyz") # multiple result, so result is an array
	bar = result[0]
	bar.send("Name")
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, "xyz")
	vm.checkCFP(t, 0, 0)
}

func skipPluginTestIfEnvNotSet(t *testing.T) {
	if os.Getenv("TEST_PLUGIN") == "" {
		t.Skip("skipping plugin related tests")
	}
}
