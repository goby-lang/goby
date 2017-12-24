package vm

import (
	"testing"
)

func TestIntegerClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Integer.class.name`, "Class"},
		{`Integer.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerArithmeticOperationWithInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`13  +  3`, 16},
		{`13  -  3`, 10},
		{`13  *  3`, 39},
		{`13  %  3`, 1},
		{`13  /  3`, 4},
		{`13  ** 3`, 2197},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerArithmeticOperationWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`13  +  '3.5'.to_f`, 16.5},
		{`13  -  '3.5'.to_f`, 9.5},
		{`13  *  '3.5'.to_f`, 45.5},
		{`13  %  '3.5'.to_f`, 2.5},
		{`13  /  '6.5'.to_f`, 2.0},
		{`4   ** '3.5'.to_f`, 128.0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerArithmeticOperationsPriority(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 / 1 + 1`, 2},
		{`0 / (1 + 1000)`, 0},
		{`5 ** (3 * 2) + 21`, 15646},
		{`(3 - 1) ** 4 / 2`, 8},
		{`(25 / 5 + 5) * 3`, 30},
		{`(25 / 5 + 5) * 2`, 20},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerArithmeticOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`1 + "p"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 - "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 ** "p"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 / "t"`, "TypeError: Expect argument to be Numeric. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerComparisonWithInteger(t *testing.T) {
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
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerComparisonWithFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 >   '2'.to_f`, false},
		{`2 >   '1'.to_f`, true},
		{`3 >   '3'.to_f`, false},
		{`1 <   '2'.to_f`, true},
		{`2 <   '1'.to_f`, false},
		{`3 <   '3'.to_f`, false},
		{`1 >=  '2'.to_f`, false},
		{`2 >=  '1'.to_f`, true},
		{`3 >=  '3'.to_f`, true},
		{`1 <=  '2'.to_f`, true},
		{`2 <=  '1'.to_f`, false},
		{`3 <=  '3'.to_f`, true},
		{`1 <=> '2'.to_f`, -1},
		{`2 <=> '1'.to_f`, 1},
		{`3 <=> '3'.to_f`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`1 > "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 >= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 < "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 <= "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
		{`1 <=> "m"`, "TypeError: Expect argument to be Numeric. got: String", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`4  ==  4`, true},
		{`4  ==  '4'.to_f`, true},
		{`4  ==  '5'.to_f`, false},
		{`4  ==  '4'`, false},
		{`4  ==  (1..3)`, false},
		{`4  ==  { a: 1 }`, false},
		{`4  ==  [1]`, false},
		{`4  ==  Float`, false},
		{`4  !=  4`, false},
		{`4  !=  '4'.to_f`, false},
		{`4  !=  '5'.to_f`, true},
		{`4  !=  '4'`, true},
		{`4  !=  (1..3)`, true},
		{`4  !=  { a: 1 }`, true},
		{`4  !=  [1]`, true},
		{`4  !=  Float`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`100.to_i`, 100},
		{`100.to_f`, 100.0},
		{`100.to_s`, "100"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// Method test

func TestIntegerEvenMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.even?`, false},
		{`2.even?`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerNextMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`2.next`, 3},
		{`1.next`, 2},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerOddMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.odd?`, true},
		{`2.odd?`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerPredMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1.pred`, 0},
		{`0.pred`, -1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerTimesMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`	a = 0
		  	3.times do
		  		a += 1
			end
			a
			`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerTimesMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`(-2).times`, "InternalError: Expect integer greater than or equal 0. got: -2", 1},
		{`2.times`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerZeroDivisionFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`6 / 0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 / -0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 / 0.0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 / -0.0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 % 0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 % -0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 % 0.0`, "ZeroDivisionError: Divided by 0", 1},
		{`6 % -0.0`, "ZeroDivisionError: Divided by 0", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
