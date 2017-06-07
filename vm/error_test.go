package vm

import "testing"

func TestUndefinedMethodError(t *testing.T) {
	evaluated := testEval(t, "a")
	_, ok := evaluated.(*UndefinedMethodErrorObject)
	if !ok {
		t.Errorf("Expect UndefinedMethodError. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestUndefinedMethodErrorStopsProgram(t *testing.T) {
	evaluated := testEval(t, `a; 1 + 1`)
	_, ok := evaluated.(*UndefinedMethodErrorObject)
	if !ok {
		t.Errorf("Expect UndefinedMethodError. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestArgumentError(t *testing.T) {
	evaluated := testEval(t, "[].count(5,4,3)")
	_, ok := evaluated.(*ArgumentErrorObject)
	if !ok {
		t.Errorf("Expect ArgumentError. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestTypeError(t *testing.T) {
	evaluated := testEval(t, "10 * \"foo\"")
	_, ok := evaluated.(*TypeErrorObject)
	if !ok {
		t.Errorf("Expect TypeError. got=%T (%+v)", evaluated, evaluated)
	}
}
