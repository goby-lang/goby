package vm

import "testing"

func TestEvalNil(t *testing.T) {
	input := `nil`

	evaluated := testEval(t, input)

	if _, ok := evaluated.(*NullObject); !ok {
		t.Fatalf("Expect result to be Null. got=%T", evaluated)
	}
}

func TestBangPrefix(t *testing.T) {
	input := `
	a = nil
	!a
	`

	evaluated := testEval(t, input)
	testBooleanObject(t, evaluated, true)
}
