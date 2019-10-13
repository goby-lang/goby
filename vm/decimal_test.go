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
		VerifyExpected(t, i, evaluated, tt.expected)
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
		{`('13.5'.to_d  /  '3.75'.to_d).to_s`, "3.6"},
		{`('16.0'.to_d  ** '3.5'.to_d).to_s`, "16384"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('13.5'.to_d  +  3).to_s`, "16.5"},
		{`('13.5'.to_d  -  3).to_s`, "10.5"},
		{`('13.5'.to_d  *  3).to_s`, "40.5"},
		{`('13.5'.to_d  /  3).to_s`, "4.5"},
		{`('13.5'.to_d  ** 3).to_s`, "2460.375"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`('16.1'.to_d  + "1.1".to_d).to_s`, "17.2"},
		{`('16.1'.to_d  + "1.1".to_f).to_s[0..13]`, "17.20000000000"},
		{`('16.1'.to_d  - "1.1".to_d).to_s`, "15"},
		{`('16.1'.to_d  - "1.1".to_f).to_s[0..13]`, "14.99999999999"},
		{`('16.1'.to_d  * "1.1".to_d).to_s`, "17.71"},
		{`('16.1'.to_d  * "1.1".to_f).to_s[0..13]`, "17.71000000000"},
		{`('16.1'.to_d  / "1.1".to_d).to_s[0..13]`, "14.63636363636"},
		{`('16.1'.to_d  / "1.1".to_f).to_s[0..13]`, "14.63636363636"},
		{`('16.1'.to_d  ** "1.1".to_d).to_s[0..13]`, "21.25731771584"},
		{`('16.1'.to_d  ** "1.1".to_f).to_s[0..13]`, "21.25731771584"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'1'.to_d >   2`, false},
		{`'2'.to_d >   1`, true},
		{`'3'.to_d >   3`, false},
		{`'1'.to_d <   2`, true},
		{`'2'.to_d <   1`, false},
		{`'3'.to_d <   3`, false},
		{`'1'.to_d >=  2`, false},
		{`'2'.to_d >=  1`, true},
		{`'3'.to_d >=  3`, true},
		{`'1'.to_d <=  2`, true},
		{`'2'.to_d <=  1`, false},
		{`'3'.to_d <=  3`, true},
		{`'1'.to_d <=> 2`, -1},
		{`'2'.to_d <=> 1`, 1},
		{`'3'.to_d <=> 3`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1'.to_d > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`'1'.to_d <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'123'.to_d    ==  123`, false},
		{`'123'.to_d.to_i    ==  123`, true},
		{`'123.5'.to_d  ==  '123.5'.to_d`, true},
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// Test type conversion

func TestDecimalToArray(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		"129.30928304982039482039842".to_d.to_a[0].to_s
		`, "6465464152491019741019921"},
		{`
		"129.30928304982039482039842".to_d.to_a[1].to_s
		`, "50000000000000000000000"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalToIntegerStringConversions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'0.3'.to_d.to_s`, "0.3"},
		{`'-0.3'.to_d.to_s`, "-0.3"},
		{`'100.3'.to_d.to_s`, "100.3"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalToInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`'0.3'.to_d.to_i`, 0},
		{`'0.5'.to_d.to_i`, 0},
		{`'0.9'.to_d.to_i`, 0},
		{`'1.1'.to_d.to_i`, 1},
		{`'-0.3'.to_d.to_i`, 0},
		{`'-0.5'.to_d.to_i`, 0},
		{`'-0.9'.to_d.to_i`, 0},
		{`'-1.3'.to_d.to_i`, -1},
		{`'100.3'.to_d.to_i`, 100},
		{`'-100.3'.to_d.to_i`, -100},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalToNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`'-13.5'.to_d.to_i`, -13},
		{`'-13.5'.to_d.to_f`, -13.5},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalToStringMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'13.5'.to_d.to_s`, "13.5"},
		{`'13.5'.to_d.fraction`, "27/2"},
		{`'13.5'.to_d.inverse.fraction`, "2/27"},
		{`'20/40'.to_d.reduction`, "1/2"},
		{`'40/20'.to_d.fraction`, "2/1"},
		{`'40/20'.to_d.reduction`, "2"},
		{`'-13.5'.to_d.numerator.to_s`, "-27"},
		{`'-13.5'.to_d.denominator.to_s`, "2"},
		{`'129.30928304982039482039842'.to_d.numerator.to_s`, "6465464152491019741019921"},
		{`'129.30928304982039482039842'.to_d.denominator.to_s`, "50000000000000000000000"},
		// The followings are permissible
		{`'1.'.to_d.to_s`, "1"},
		{`'.1'.to_d.to_s`, "0.1"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalToStringFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'1.1.1'.to_d`, "ArgumentError: Invalid numeric string. got: 1.1.1", 1},
		{`'1/1/1'.to_d`, "ArgumentError: Invalid numeric string. got: 1/1/1", 1},
		{`'1.1/1'.to_d`, "ArgumentError: Invalid numeric string. got: 1.1/1", 1},
		{`'1/1.1'.to_d`, "ArgumentError: Invalid numeric string. got: 1/1.1", 1},
		{`'1..1'.to_d`, "ArgumentError: Invalid numeric string. got: 1..1", 1},
		{`'..1'.to_d`, "ArgumentError: Invalid numeric string. got: ..1", 1},
		{`'1//1'.to_d`, "ArgumentError: Invalid numeric string. got: 1//1", 1},
		{`'abc'.to_d`, "ArgumentError: Invalid numeric string. got: abc", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalZeroDivisionFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`(6.0).to_d / 0`, "ZeroDivisionError: Divided by 0", 1},
		{`(6.0).to_d / -0`, "ZeroDivisionError: Divided by 0", 1},
		{`(6.0).to_d / "0".to_d`, "ZeroDivisionError: Divided by 0", 1},
		{`(6.0).to_d / "-0".to_d`, "ZeroDivisionError: Divided by 0", 1},
		{`(6.0).to_d / "0".to_f`, "ZeroDivisionError: Divided by 0", 1},
		{`(6.0).to_d / "-0".to_f`, "ZeroDivisionError: Divided by 0", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalTruncation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`["9.0".to_d][0].to_s`, "9"},
		{`["9.00".to_d][0].to_s`, "9"},
		{`["1.0".to_d][0].to_s`, "1"},
		{`["1.00".to_d][0].to_s`, "1"},
		{`["0.0".to_d][0].to_s`, "0"},
		{`["0.00".to_d][0].to_s`, "0"},
		{`["0.00".to_d][0].class.to_s`, "Decimal"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalNumberOfDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"3.14".to_d.to_s`, "3.14"},
		{`"3.141".to_d.to_s`, "3.141"},
		{`"3.1415".to_d.to_s`, "3.1415"},
		{`"3.14159".to_d.to_s`, "3.14159"},
		{`"3.141592".to_d.to_s`, "3.141592"},
		{`"3.1415926".to_d.to_s`, "3.1415926"},
		{`"3.14159265".to_d.to_s`, "3.14159265"},
		{`"3.141592653".to_d.to_s`, "3.141592653"},
		{`"3.1415926535".to_d.to_s`, "3.1415926535"},
		{`"3.14159265358".to_d.to_s`, "3.14159265358"},
		{`"3.141592653589".to_d.to_s`, "3.141592653589"},
		{`"3.1415926535897".to_d.to_s`, "3.1415926535897"},
		{`"3.14159265358979".to_d.to_s`, "3.14159265358979"},
		{`"3.141592653589793".to_d.to_s`, "3.141592653589793"},
		{`"3.1415926535897932".to_d.to_s`, "3.1415926535897932"},
		{`"3.1415926535897932384626".to_d.to_s`, "3.1415926535897932384626"},
		{`"-3.14".to_d.to_s`, "-3.14"},
		{`"-3.141".to_d.to_s`, "-3.141"},
		{`"-3.1415".to_d.to_s`, "-3.1415"},
		{`"-3.14159".to_d.to_s`, "-3.14159"},
		{`"-3.141592".to_d.to_s`, "-3.141592"},
		{`"-3.1415926".to_d.to_s`, "-3.1415926"},
		{`"-3.14159265".to_d.to_s`, "-3.14159265"},
		{`"-3.141592653".to_d.to_s`, "-3.141592653"},
		{`"-3.1415926535".to_d.to_s`, "-3.1415926535"},
		{`"-3.14159265358".to_d.to_s`, "-3.14159265358"},
		{`"-3.141592653589".to_d.to_s`, "-3.141592653589"},
		{`"-3.1415926535897".to_d.to_s`, "-3.1415926535897"},
		{`"-3.14159265358979".to_d.to_s`, "-3.14159265358979"},
		{`"-3.141592653589793".to_d.to_s`, "-3.141592653589793"},
		{`"-3.1415926535897932".to_d.to_s`, "-3.1415926535897932"},
		{`"-3.1415926535897932384626".to_d.to_s`, "-3.1415926535897932384626"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalMinusZero(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"0.0".to_d.to_s`, "0"},
		{`"-0.0".to_d.to_s`, "0"},

		{`("-0.0".to_d == "0.0".to_d).to_s`, "true"},

		{`("0.0".to_d + "0.0".to_d).to_s`, "0"},
		{`("0.0".to_d + "-0.0".to_d).to_s`, "0"},
		{`("0.0".to_d - "0.0".to_d).to_s`, "0"},
		{`("0.0".to_d - "-0.0".to_d).to_s`, "0"},

		{`("-0.0".to_d + "0.0".to_d).to_s`, "0"},
		{`("-0.0".to_d + "-0.0".to_d).to_s`, "0"},
		{`("-0.0".to_d - "0.0".to_d).to_s`, "0"},
		{`("-0.0".to_d - "-0.0".to_d).to_s`, "0"},

		{`("0.0".to_d * "0.0".to_d).to_s`, "0"},
		{`("0.0".to_d * "-0.0".to_d).to_s`, "0"},
		{`("-0.0".to_d * "0.0".to_d).to_s`, "0"},
		{`("-0.0".to_d * "-0.0".to_d).to_s`, "0"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDecimalDupMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"1.1".to_d.dup == "1.1".to_d`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
