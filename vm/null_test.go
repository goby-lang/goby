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

func TestNullComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123 == nil`, false},
		{`nil == nil`, true},
		{`nil == 123`, false},
		{`123 != nil`, true},
		{`nil != nil`, false},
		{`nil != 123`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}
