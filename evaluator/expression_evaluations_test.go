package evaluator

import "testing"

func TestEvalInfixIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`if (10 > 5) {
				100
			} else {
				-10
			}
			`,
			100,
		},
		{
			`if (5 != 5) {
				false
			} else {
				true
			}
			`,
			true,
		},
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, tt.expected.(int64))
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case nil:
			testNullObject(t, evaluated)
		}

	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Stan " + "Lo"`, "Stan Lo"},
		{`"Dog" + "&" + "Cat"`, "Dog&Cat"},
		{`"Dog" == "Dog"`, true},
		{`"1234" > "123"`, true},
		{`"1234" < "123"`, false},
		{`"1234" != "123"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case string:
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestEvalInfixBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBangPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!5", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalMinusPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"-10", -10},
		{"--10", 10},
		{"--5", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"st0012"`, "st0012"},
		{`'Monkey'`, "Monkey"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
