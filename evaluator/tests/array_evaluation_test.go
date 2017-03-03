package evaluator_test

import (
	"testing"
	"github.com/st0012/Rooby/object"
)

func TestEvalArrayExpression(t *testing.T) {
	input := `
	[1, "234", true]
	`

	evaluated := testEval(t, input)

	arr, ok := evaluated.(*object.ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be an array. got=%T", evaluated)
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testStringObject(t, arr.Elements[1], "234")
	testBooleanObject(t, arr.Elements[2], true)
}

func TestEvalArrayIndex(t *testing.T) {
	tests := []struct{
		input string
		expected interface{}
	}{
		{`
			[1, 2, 10, 5][2]
		`, int64(10)},
		{`
			[1, "a", 10, 5][1]
		`, "a"},
		{`
			[][1]
		`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case nil:
			_, ok := evaluated.(*object.Null)

			if !ok {
				t.Fatalf("expect input: \"%s\"'s result should be Null. got=%T", tt.input, evaluated)
			}
		}
	}
}
