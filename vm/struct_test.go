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

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, "xyz!")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingStructFunctionWithReturnError(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	result, err = bar.send("Name", "!")
	err
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, nil)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingStructFuncWithInt64(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	bar.send("Add", 10, 100.to_int64) # Add is func(int, int64) int64
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, 110)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingStructFuncWithGoObject(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/struct/struct.go"
	bar, err = p.send("NewBar", "xyz") # multiple result, so result is an array
	p.send("GetBarName", bar) # GetBarName is func(*Bar) string
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, "xyz")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}
