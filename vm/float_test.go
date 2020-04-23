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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperationWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`13.5  +  3.2`, 16.7},
		{`13.5  -  3.2`, 10.3},
		{`13.5  *  3.2`, 43.2},
		{`13.5  %  3.75`, 2.25},
		{`13.5  /  3.75`, 3.6},
		{`16.0  ** 3.5`, 16384.0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperationWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`13.5  +  3`, 16.5},
		{`13.5  -  3`, 10.5},
		{`13.5  *  3`, 40.5},
		{`13.5  %  3`, 1.5},
		{`13.5  /  3`, 4.5},
		{`13.5  ** 3`, 2460.375},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`1.1 + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1.1 - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1.1 ** "p"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1.1 / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.5 >   2.5`, false},
		{`2.5 >   1.5`, true},
		{`3.5 >   3.5`, false},
		{`1.5 <   2.5`, true},
		{`2.5 <   1.5`, false},
		{`3.5 <   3.5`, false},
		{`1.5 >=  2.5`, false},
		{`2.5 >=  1.5`, true},
		{`3.5 >=  3.5`, true},
		{`1.5 <=  2.5`, true},
		{`2.5 <=  1.5`, false},
		{`3.5 <=  3.5`, true},
		{`1.5 <=> 2.5`, -1},
		{`2.5 <=> 1.5`, 1},
		{`3.5 <=> 3.5`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 >   2`, false},
		{`2 >   1`, true},
		{`3 >   3`, false},
		{`1 <   2`, true},
		{`2 <   1`, false},
		{`3 <   3`, false},
		{`1 >=  2`, false},
		{`2 >=  1`, true},
		{`3 >=  3`, true},
		{`1 <=  2`, true},
		{`2 <=  1`, false},
		{`3 <=  3`, true},
		{`1 <=> 2`, -1},
		{`2 <=> 1`, 1},
		{`3 <=> 3`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`1 > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1 >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1 < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1 <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
		{`1 <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFloatEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`123.5  ==  123.5`, true},
		{`123    ==  123`, true},
		{`123.5  ==  124`, false},
		{`123.5  ==  "123.5"`, false},
		{`123.5  ==  (1..3)`, false},
		{`123.5  ==  { a: 1 }`, false},
		{`123.5  ==  [1]`, false},
		{`123.5  ==  Float`, false},
		{`123.5  !=  123.5`, false},
		{`123.5  !=  123`, true},
		{`123.5  !=  124`, true},
		{`123.5  !=  "123.5"`, true},
		{`123.5  !=  (1..3)`, true},
		{`123.5  !=  { a: 1 }`, true},
		{`123.5  !=  [1]`, true},
		{`123.5  !=  Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatConversions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`(100.3).to_i`, 100},
		{`(100.3).to_s`, "100.3"},
		{`(100.3).to_d.to_s`, "100.3"},
		{`
		(3.14159265358979).to_d.to_s`,
			"3.14159265358979"},
		{`
		(-273.150000000).to_d.to_s`,
			"-273.15"},
		{`100.3.to_i`, 100},
		{`100.3.to_s`, "100.3"},
		{`100.3.to_d.to_s`, "100.3"},
		{`
		3.14159265358979.to_d.to_s`,
			"3.14159265358979"},
		{`
		-273.150000000.to_d.to_s`,
			"-273.15"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`(0.1 + 0.2).to_s`, "0.30000000000000004"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatZeroDivisionFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`6.0 / 0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 / -0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 / 0.0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 / -0.0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 % 0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 % -0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 % 0.0`, "ZeroDivisionError: Divided by 0", 1, 1},
		{`6.0 % -0.0`, "ZeroDivisionError: Divided by 0", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestFloatTruncation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`[9.0][0].to_s`, "9.0"},
		{`[9.00][0].to_s`, "9.0"},
		{`[1.0][0].to_s`, "1.0"},
		{`[1.00][0].to_s`, "1.0"},
		{`[0.0][0].to_s`, "0.0"},
		{`[0.00][0].to_s`, "0.0"},
		{`[0.0][0].class.to_s`, "Float"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatNumberOfDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`3.14.to_s`, "3.14"},
		{`3.141.to_s`, "3.141"},
		{`3.1415.to_s`, "3.1415"},
		{`3.14159.to_s`, "3.14159"},
		{`3.141592.to_s`, "3.141592"},
		{`3.1415926.to_s`, "3.1415926"},
		{`3.14159265.to_s`, "3.14159265"},
		{`3.141592653.to_s`, "3.141592653"},
		{`3.1415926535.to_s`, "3.1415926535"},
		{`3.14159265358.to_s`, "3.14159265358"},
		{`3.141592653589.to_s`, "3.141592653589"},
		{`3.1415926535897.to_s`, "3.1415926535897"},
		{`3.14159265358979.to_s`, "3.14159265358979"},
		{`3.141592653589793.to_s`, "3.141592653589793"},
		{`3.1415926535897932.to_s`, "3.141592653589793"},
		{`3.1415926535897932384626.to_s`, "3.141592653589793"},
		{`-3.14.to_s`, "-3.14"},
		{`-3.141.to_s`, "-3.141"},
		{`-3.1415.to_s`, "-3.1415"},
		{`-3.14159.to_s`, "-3.14159"},
		{`-3.141592.to_s`, "-3.141592"},
		{`-3.1415926.to_s`, "-3.1415926"},
		{`-3.14159265.to_s`, "-3.14159265"},
		{`-3.141592653.to_s`, "-3.141592653"},
		{`-3.1415926535.to_s`, "-3.1415926535"},
		{`-3.14159265358.to_s`, "-3.14159265358"},
		{`-3.141592653589.to_s`, "-3.141592653589"},
		{`-3.1415926535897.to_s`, "-3.1415926535897"},
		{`-3.14159265358979.to_s`, "-3.14159265358979"},
		{`-3.141592653589793.to_s`, "-3.141592653589793"},
		{`-3.1415926535897932.to_s`, "-3.141592653589793"},
		{`-3.1415926535897932384626.to_s`, "-3.141592653589793"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// API tests

func TestFloatAbs(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"34.56.abs", 34.56},
		{"-34.56.abs", 34.56},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatMinusZero(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`0.0.to_s`, "0.0"},
		{`-0.0.to_s`, "-0.0"},
		{`(-0.0 == 0.0).to_s`, "true"},
		{`(0.0 + 0.0).to_s`, "0.0"},
		{`(0.0 + -0.0).to_s`, "0.0"},
		{`(0.0 - 0.0).to_s`, "0.0"},
		{`(0.0 - -0.0).to_s`, "0.0"},
		{`(-0.0 + 0.0).to_s`, "0.0"},
		{`(-0.0 + -0.0).to_s`, "-0.0"},
		{`(-0.0 - 0.0).to_s`, "-0.0"},
		{`(-0.0 - -0.0).to_s`, "0.0"},
		{`(0.0 * 0.0).to_s`, "0.0"},
		{`(0.0 * -0.0).to_s`, "-0.0"},
		{`(-0.0 * 0.0).to_s`, "-0.0"},
		{`(-0.0 * -0.0).to_s`, "0.0"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatCeil(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"1.2.ceil", 2},
		{"2.0.ceil", 2},
		{"-1.2.ceil", -1},
		{"-2.0.ceil", -2},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatDupMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.1.dup == 1.1`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
func TestFloatFloor(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"1.2.floor", 1},
		{"2.0.floor", 2},
		{"-1.2.floor", -2},
		{"-2.0.floor", -2},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatNegative(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"-1.0.negative?", true},
		{"0.0.negative?", false},
		{"1.0.negative?", false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatPositive(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"-1.0.positive?", false},
		{"0.0.positive?", false},
		{"1.0.positive?", true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatRound(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"1.115.round", 1.0},
		{"1.115.round(1)", 1.1},
		{"1.115.round(2)", 1.12},
		{"-1.115.round", -1.0},
		{"-1.115.round(1)", -1.1},
		{"-1.115.round(2)", -1.12},
		{"1.115.round(-1)", 0.0},
		{"-1.115.round(-1)", 0.0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFloatZero(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"0.0.zero?", true},
		{"1.0.zero?", false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
