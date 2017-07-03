package vm

import "testing"

func TestEvalNil(t *testing.T) {
	input := `nil`

	evaluated := testEval(t, input)
	checkExpected(t, evaluated, nil)
}

func TestBangPrefix(t *testing.T) {
	input := `
	a = nil
	!a
	`

	evaluated := testEval(t, input)
	checkExpected(t, evaluated, true)
}
