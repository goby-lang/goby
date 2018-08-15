package vm

import "testing"

func TestNullClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Null.class.name`, "Class"},
		{`Null.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestEvalNil(t *testing.T) {
	input := `nil`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, nil)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestNilInspect(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`nil.to_s`, ""},
		{`nil.inspect`, "nil"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestNullTypeConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`nil.to_i`, 0},
		{`nil.to_s`, ""},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestNullTypeConversionFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`nil.to_s(1)`, "ArgumentError: Expect 0 argument. got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Method test

func TestNullBangPrefixMethod(t *testing.T) {
	input := `
	a = nil
	!a
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, true)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestNullIsNilMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.nil?`, false},
		{`"Hello World".nil?`, false},
		{`(2..10).nil?`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.nil?`, false},
		{`[1, 2, 3, 4, 5].nil?`, false},
		{`true.nil?`, false},
		{`String.nil?`, false},
		{`nil.nil?`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestNullIsNilMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`nil.nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
