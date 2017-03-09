package evaluator_test

import (
	"github.com/st0012/Rooby/evaluator"
	"testing"
)

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

func testStringObject(t *testing.T, obj evaluator.Object, expected string) bool {
	result, ok := obj.(*evaluator.StringObject)
	if !ok {
		t.Errorf("object is not a String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. expect=%s, got=%s", expected, result.Value)
		return false
	}

	return true
}
