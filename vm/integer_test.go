package vm

import (
	"testing"
)

func TestInitilaize(t *testing.T) {
	expected := 101
	i := initilaizeInteger(expected)
	if expected != i.Value {
		t.Fatalf("Expect: %d. got=%d", expected, i.Value)
	}
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
		{`25 / 5`, 5},
		{`25 > 5`, true},
		{`4 > 6`, false},
		{`-5 < -4`, true},
		{`100 < 81`, false},
		{`8 == 8`, true},
		{`3 == 4`, false},
		{`3 != 4`, true},
		{`4 != 4`, false},
		{`10++`, 11},
		{`1--`, 0},
		{`100.to_s`, "100"},
		{`1.even`, false},
		{`2.even`, true},
		{`1.odd`, true},
		{`2.odd`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case int:
			testIntegerObject(t, evaluated, tt.expected.(int))
		case string:
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}
