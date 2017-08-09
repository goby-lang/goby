package vm

import (
	"testing"
)

func TestIntegerArithmeticOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 + 2`, 3},
		{`10 + 0`, 10},
		{`22 - 10`, 12},
		{`2 - 10`, -8},
		{`5 * 20`, 100},
		{`4 % 2`, 0},
		{`10 % 3`, 1},
		{`6 % 4`, 2},
		{`25 / 5`, 5},
		{`25 > 5`, true},
		{`4 > 6`, false},
		{`-5 < -4`, true},
		{`100 < 81`, false},
		{`5 ** 4`, 625},
		{`25 / 5`, 5},
		{`1 / 1 + 1`, 2},
		{`0 / (1 + 1000)`, 0},
		{`5 ** (3 * 2) + 21`, 15646},
		{`(3 - 1) ** 4 / 2`, 8},
		{`(25 / 5 + 5) * 3`, 30},
		{`(25 / 5 + 5) * 2`, 20},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerArithmeticOperationFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected string
	}{
		{`
		1 + "p"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 - "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 ** "p"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 / "t"
		`, "TypeError: Expect argument to be Integer. got: String"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, TypeError, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`25 > 5`, true},
		{`4 > 6`, false},
		{`-5 < -4`, true},
		{`100 < 81`, false},
		{`25 > 5`, true},
		{`4 > 6`, false},
		{`4 >= 4`, true},
		{`2 >= 5`, false},
		{`-5 < -4`, true},
		{`100 < 81`, false},
		{`10 <= 10`, true},
		{`10 <= 0`, false},
		{`10 <=> 0`, 1},
		{`1 <=> 2`, -1},
		{`4 <=> 4`, 0},
		{`123 == 123`, true},
		{`123 == 124`, false},
		{`123 == "123"`, false},
		{`123 == (1..3)`, false},
		{`123 == { a: 1, b: 2 }`, false},
		{`123 == [1, "String", true, 2..5]`, false},
		{`123 == Integer`, false},
		{`123 != 123`, false},
		{`123 != 124`, true},
		{`123 != "123"`, true},
		{`123 != (1..3)`, true},
		{`123 != { a: 1, b: 2 }`, true},
		{`123 != [1, "String", true, 2..5]`, true},
		{`123 != Integer`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerComparisonFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected string
	}{
		{`
		1 > "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 >= "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 < "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 <= "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
		{`
		1 <=> "m"
		`, "TypeError: Expect argument to be Integer. got: String"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, TypeError, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`100.to_i`, 100},
		{`100.to_s`, "100"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

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
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
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
		  		a++
			end
			a
			`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIntegerTimesMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		(-2).times
		`, newError("Expect paramentr to be greater 0. got=-2")},
		{`
		2.times
		`, newError("Can't yield without a block")},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}
