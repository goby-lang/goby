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
