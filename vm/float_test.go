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

func TestFloatArithmeticOperationWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'13.5'.to_f  +  '3.2'.to_f`, 16.7},
		{`'13.5'.to_f  -  '3.2'.to_f`, 10.3},
		{`'13.5'.to_f  *  '3.2'.to_f`, 43.2},
		{`'13.5'.to_f  %  '3.75'.to_f`, 2.25},
		{`'13.5'.to_f  /  '3.75'.to_f`, 3.6},
		{`'16.0'.to_f  ** '3.5'.to_f`, 16384.0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperationWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'13.5'.to_f  +  3`, 16.5},
		{`'13.5'.to_f  -  3`, 10.5},
		{`'13.5'.to_f  *  3`, 40.5},
		{`'13.5'.to_f  %  3`, 1.5},
		{`'13.5'.to_f  /  3`, 4.5},
		{`'13.5'.to_f  ** 3`, 2460.375},
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
		{`'1'.to_f + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f ** "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1.5'.to_f >   '2.5'.to_f`, false},
		{`'2.5'.to_f >   '1.5'.to_f`, true},
		{`'3.5'.to_f >   '3.5'.to_f`, false},
		{`'1.5'.to_f <   '2.5'.to_f`, true},
		{`'2.5'.to_f <   '1.5'.to_f`, false},
		{`'3.5'.to_f <   '3.5'.to_f`, false},
		{`'1.5'.to_f >=  '2.5'.to_f`, false},
		{`'2.5'.to_f >=  '1.5'.to_f`, true},
		{`'3.5'.to_f >=  '3.5'.to_f`, true},
		{`'1.5'.to_f <=  '2.5'.to_f`, true},
		{`'2.5'.to_f <=  '1.5'.to_f`, false},
		{`'3.5'.to_f <=  '3.5'.to_f`, true},
		{`'1.5'.to_f <=> '2.5'.to_f`, -1},
		{`'2.5'.to_f <=> '1.5'.to_f`, 1},
		{`'3.5'.to_f <=> '3.5'.to_f`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1'.to_f >   2`, false},
		{`'2'.to_f >   1`, true},
		{`'3'.to_f >   3`, false},
		{`'1'.to_f <   2`, true},
		{`'2'.to_f <   1`, false},
		{`'3'.to_f <   3`, false},
		{`'1'.to_f >=  2`, false},
		{`'2'.to_f >=  1`, true},
		{`'3'.to_f >=  3`, true},
		{`'1'.to_f <=  2`, true},
		{`'2'.to_f <=  1`, false},
		{`'3'.to_f <=  3`, true},
		{`'1'.to_f <=> 2`, -1},
		{`'2'.to_f <=> 1`, 1},
		{`'3'.to_f <=> 3`, 0},
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
		{`'1'.to_f > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_f <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFloatEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'123.5'.to_f  ==  '123.5'.to_f`, true},
		{`'123'.to_f    ==  123`, true},
		{`'123.5'.to_f  ==  '124'.to_f`, false},
		{`'123.5'.to_f  ==  "123.5"`, false},
		{`'123.5'.to_f  ==  (1..3)`, false},
		{`'123.5'.to_f  ==  { a: 1 }`, false},
		{`'123.5'.to_f  ==  [1]`, false},
		{`'123.5'.to_f  ==  Float`, false},
		{`'123.5'.to_f  !=  '123.5'.to_f`, false},
		{`'123.5'.to_f  !=  123`, true},
		{`'123.5'.to_f  !=  '124'.to_f`, true},
		{`'123.5'.to_f  !=  "123.5"`, true},
		{`'123.5'.to_f  !=  (1..3)`, true},
		{`'123.5'.to_f  !=  { a: 1 }`, true},
		{`'123.5'.to_f  !=  [1]`, true},
		{`'123.5'.to_f  !=  Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
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
		{`'100.3'.to_f.to_d.to_s`, "100.3"},
		{`
		'3.14159265358979'.to_f.to_d.to_s`,
			"3.14159265358979"},
		{`
		'-273.150000000'.to_f.to_d.to_s`,
			"-273.15"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
