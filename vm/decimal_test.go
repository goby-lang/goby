package vm

import (
	"testing"
)

func TestDecimalClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Decimal.class.name`, "Class"},
		{`Decimal.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithDecimal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('13.5'.to_d  +  '3.5'.to_d).to_s`, "17"},
		{`('13.5'.to_d  +  '3.2'.to_d).to_s`, "16.7"},
		{`('13.5'.to_d  -  '3.2'.to_d).to_s`, "10.3"},
		{`('13.5'.to_d  *  '3.2'.to_d).to_s`, "43.2"},
		//{`'13.5'.to_d  %  '3.75'.to_d`, 2.25},
		{`('13.5'.to_d  /  '3.75'.to_d).to_s`, "3.6"},
		//{`'16.0'.to_d  ** '3.5'.to_d`, 16384.0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

//func TestDecimalArithmeticOperationWithInteger(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected interface{}
//	}{
//		{`'13.5'.to_d  +  3`, 16.5},
//		{`'13.5'.to_d  -  3`, 10.5},
//		{`'13.5'.to_d  *  3`, 40.5},
//		{`'13.5'.to_d  %  3`, 1.5},
//		{`'13.5'.to_d  /  3`, 4.5},
//		{`'13.5'.to_d  ** 3`, 2460.375},
//	}
//
//	for i, tt := range tests {
//		v := initTestVM()
//		evaluated := v.testEval(t, tt.input, getFilename())
//		checkExpected(t, i, evaluated, tt.expected)
//		v.checkCFP(t, i, 0)
//		v.checkSP(t, i, 1)
//	}
//}

func TestDecimalArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		//{`'1'.to_d ** "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1.5'.to_d >   '2.5'.to_d`, false},
		{`'2.5'.to_d >   '1.5'.to_d`, true},
		{`'3.5'.to_d >   '3.5'.to_d`, false},
		{`'1.5'.to_d <   '2.5'.to_d`, true},
		{`'2.5'.to_d <   '1.5'.to_d`, false},
		{`'3.5'.to_d <   '3.5'.to_d`, false},
		{`'1.5'.to_d >=  '2.5'.to_d`, false},
		{`'2.5'.to_d >=  '1.5'.to_d`, true},
		{`'3.5'.to_d >=  '3.5'.to_d`, true},
		{`'1.5'.to_d <=  '2.5'.to_d`, true},
		{`'2.5'.to_d <=  '1.5'.to_d`, false},
		{`'3.5'.to_d <=  '3.5'.to_d`, true},
		{`'1.5'.to_d <=> '2.5'.to_d`, -1},
		{`'2.5'.to_d <=> '1.5'.to_d`, 1},
		{`'3.5'.to_d <=> '3.5'.to_d`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

//func TestDecimalComparisonWithInteger(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected interface{}
//	}{
//		{`'1'.to_d >   2`, false},
//		{`'2'.to_d >   1`, true},
//		{`'3'.to_d >   3`, false},
//		{`'1'.to_d <   2`, true},
//		{`'2'.to_d <   1`, false},
//		{`'3'.to_d <   3`, false},
//		{`'1'.to_d >=  2`, false},
//		{`'2'.to_d >=  1`, true},
//		{`'3'.to_d >=  3`, true},
//		{`'1'.to_d <=  2`, true},
//		{`'2'.to_d <=  1`, false},
//		{`'3'.to_d <=  3`, true},
//		{`'1'.to_d <=> 2`, -1},
//		{`'2'.to_d <=> 1`, 1},
//		{`'3'.to_d <=> 3`, 0},
//	}
//
//	for i, tt := range tests {
//		v := initTestVM()
//		evaluated := v.testEval(t, tt.input, getFilename())
//		checkExpected(t, i, evaluated, tt.expected)
//		v.checkCFP(t, i, 0)
//		v.checkSP(t, i, 1)
//	}
//}

func TestDecimalComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`'1'.to_d <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'123.5'.to_d  ==  '123.5'.to_d`, true},
		//{`'123'.to_d    ==  123`, true},
		{`'123.5'.to_d  ==  '124'.to_d`, false},
		{`'123.5'.to_d  ==  "123.5"`, false},
		{`'123.5'.to_d  ==  (1..3)`, false},
		{`'123.5'.to_d  ==  { a: 1 }`, false},
		{`'123.5'.to_d  ==  [1]`, false},
		{`'123.5'.to_d  ==  Float`, false},
		{`'123.5'.to_d  !=  '123.5'.to_d`, false},
		{`'123.5'.to_d  !=  123`, true},
		{`'123.5'.to_d  !=  '124'.to_d`, true},
		{`'123.5'.to_d  !=  "123.5"`, true},
		{`'123.5'.to_d  !=  (1..3)`, true},
		{`'123.5'.to_d  !=  { a: 1 }`, true},
		{`'123.5'.to_d  !=  [1]`, true},
		{`'123.5'.to_d  !=  Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalConversions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//{`'100.3'.to_d.to_i`, 100},
		{`'100.3'.to_d.to_s`, "100.3"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
