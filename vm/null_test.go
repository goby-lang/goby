package vm

import "testing"

func TestEvalNil(t *testing.T) {
	input := `nil`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, nil)
	vm.checkCFP(t, 0, 0)
}

func TestBangPrefix(t *testing.T) {
	input := `
	a = nil
	!a
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, true)
	vm.checkCFP(t, 0, 0)
}
