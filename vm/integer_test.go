package vm

import (
	"testing"
)

func TestInitilaize(t *testing.T) {
	expected := 101
	i := initIntegerObject(expected)
	checkExpected(t, i, expected)
}

func TestEvalInteger(t *testing.T) {
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
		{`8 == 8`, true},
		{`3 == 4`, false},
		{`3 != 4`, true},
		{`4 != 4`, false},
		{`10++`, 11},
		{`1--`, 0},
		{`100.to_s`, "100"},
		{`100.to_i`, 100},
		{`1.even`, false},
		{`2.even`, true},
		{`1.odd`, true},
		{`2.odd`, false},
		{`1 / 1 + 1`, 2},
		{`0 / (1 + 1000)`, 0},
		{`5 ** (3 * 2) + 21`, 15646},
		{`(3 - 1) ** (3++) / 2`, 8},
		{`(25 / 5 + 5) * (2++)`, 30},
		{`(25 / 5 + 5) * 2++`, 21},
		{`2.next`, 3},
		{`1.next`, 2},
		{`1.pred`, 0},
		{`0.pred`, -1},
		{`	a = 0
		  	3.times do
		  		a++
			end
			a
			`, 3},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		1 + "p"
		`, newError("expect argument to be Integer type")},
		{`
		1 - "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 ** "p"
		`, newError("expect argument to be Integer type")},
		{`
		1 / "t"
		`, newError("expect argument to be Integer type")},
		{`
		1 > "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 >= "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 < "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 <= "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 <=> "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 == "m"
		`, newError("expect argument to be Integer type")},
		{`
		1 != "m"
		`, newError("expect argument to be Integer type")},
		{`
		(-2).times
		`, newError("Expect paramentr to be greater 0. got=-2")},
		{`
		2.times
		`, newError("Can't yield without a block")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}
