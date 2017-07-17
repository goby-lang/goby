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

func TestNullIsNilMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.is_nil`, false},
		{`"Hello World".is_nil`, false},
		{`(2..10).is_nil`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_nil`, false},
		{`[1, 2, 3, 4, 5].is_nil`, false},
		{`true.is_nil`, false},
		{`String.is_nil`, false},
		{`nil.is_nil`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}
