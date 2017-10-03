package vm

import (
	"testing"
)

func TestFloatClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Float.class.name`, "Class"},
		{`Float.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1.5'.to_f + '2'.to_f`, 3.5},
		{`'22.5'.to_f - '10'.to_f`, 12.5},
		{`'5.5'.to_f * '20'.to_f`, 110.0},
		{`'4.5'.to_f % '2'.to_f`, 0.5},
		{`'10.5'.to_f % '3.5'.to_f`, 0.0},
		{`'25.5'.to_f / '5'.to_f`, 5.1},
		{`'5.5'.to_f ** '2'.to_f`, 30.25},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_f + "p"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f - "m"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f ** "p"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f / "t"`, "TypeError: Expect argument to be Float. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'25.5'.to_f > '5'.to_f`, true},
		{`'4.5'.to_f > '6'.to_f`, false},
		{`'-5.5'.to_f < '-4'.to_f`, true},
		{`'100.5'.to_f < '81'.to_f`, false},
		{`'25.5'.to_f > '5'.to_f`, true},
		{`'4.5'.to_f > '6'.to_f`, false},
		{`'4.5'.to_f >= '4'.to_f`, true},
		{`'2.5'.to_f >= '5'.to_f`, false},
		{`'-5.5'.to_f < '-4'.to_f`, true},
		{`'100.5'.to_f < '81'.to_f`, false},
		{`'10.5'.to_f <= '10.5'.to_f`, true},
		{`'10.5'.to_f <= '0'.to_f`, false},
		{`'10.5'.to_f <=> '0'.to_f`, 1},
		{`'1.5'.to_f <=> '2'.to_f`, -1},
		{`'4.5'.to_f <=> '4.5'.to_f`, 0},
		{`'123.5'.to_f == '123.5'.to_f`, true},
		{`'123.5'.to_f == '124'.to_f`, false},
		{`'123.5'.to_f == "'123'.to_f"`, false},
		{`'123.5'.to_f == (1..3)`, false},
		{`'123.5'.to_f == { a: '1'.to_f, b: '2'.to_f }`, false},
		{`'123.5'.to_f == ['1'.to_f, "String", true, 2..5]`, false},
		{`'123.5'.to_f == Float`, false},
		{`'123.5'.to_f != '123.5'.to_f`, false},
		{`'123.5'.to_f != '124'.to_f`, true},
		{`'123.5'.to_f != "'123'.to_f"`, true},
		{`'123.5'.to_f != (1..3)`, true},
		{`'123.5'.to_f != { a: '1'.to_f, b: '2'.to_f }`, true},
		{`'123.5'.to_f != ['1'.to_f, "String", true, 2..5]`, true},
		{`'123.5'.to_f != Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_f > "m"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f >= "m"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f < "m"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f <= "m"`, "TypeError: Expect argument to be Float. got: String", 1},
		{`'1'.to_f <=> "m"`, "TypeError: Expect argument to be Float. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestFloatConversions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'100.3'.to_f.to_i`, 100},
		{`'100.3'.to_f.to_s`, "100.3"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
