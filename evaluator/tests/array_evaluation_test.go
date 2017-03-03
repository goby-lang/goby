package evaluator_test

import (
	"github.com/st0012/Rooby/object"
	"testing"
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
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			[][1]
		`, nil},
		{`
			[1, 2, 10, 5][2]
		`, int64(10)},
		{`
			[1, "a", 10, 5][1]
		`, "a"},
		{`
			a = [1, "a", 10, 5]
			a[0]
		`, 1},
		{`
			a = [1, "a", 10, 5]
			a[2] = a[1]
			a[2]

		`, "a"},
		{`
			a = []
			a[10] = 100
			a[10]
		`, 100},
		{`
			a = []
			a[10] = 100
			a[0]
		`, nil},
		{`
			a = [1, 2 ,3 ,5 , 10]
			a[0] = a[1] + a[2] + a[3] * a[4]
			a[0]
		`, 100},
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

				t.Fatalf("expect input: \"%s\"'s result should be Null. got=%T(%s)", tt.input, evaluated, evaluated.Inspect())
			}
		}
	}
}
