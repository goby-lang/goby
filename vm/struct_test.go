// +build linux

package vm

import "testing"

func TestCallingStructFunctionWithReturnValue(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	result, err = bar.send("Name", "!")
	result
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, "xyz!")
	vm.checkCFP(t, 0, 0)
}

func TestCallingStructFunctionWithReturnError(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	result, err = bar.send("Name", "!")
	err
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, nil)
	vm.checkCFP(t, 0, 0)
}

func TestCallingStructFuncWithDifferentType(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	bar.send("Add", 10, 100.to_int64) # Add is func(int, int64) int64
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, 110)
	vm.checkCFP(t, 0, 0)
}
