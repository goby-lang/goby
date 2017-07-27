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

func TestNullAssignmentByOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`a = nil; a ||= 123;       a;`, 123},
		{`a = nil; a ||= "string";  a;`, "string"},
		{`a = nil; a ||= nil;     a;`, nil},
		{`a = nil; a ||= (1..4);    a.to_s;`, "(1..4)"},
		{`a = nil; a ||= { b: 1 };  a["b"];`, 1},
		{`a = nil; a ||= Object;    a.name;`, "Object"},
		{`a = nil; a ||= [1, 2, 3]; a[0];`, 1},
		{`a = nil; a ||= [1, 2, 3]; a[1];`, 2},
		{`a = nil; a ||= [1, 2, 3]; a[2];`, 3},
		{`a = nil; a ||= nil;       a;`, nil},
		{`a = nil; a ||= nil || 1;  a;`, 1},
		{`a = nil; a ||= 1 || nil;  a;`, 1},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestNullIsNilMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`nil.is_nil("Hello")`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}
