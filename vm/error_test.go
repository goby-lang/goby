package vm

import "testing"

func TestUndefinedMethodError(t *testing.T) {
	evaluated := testEval(t, "a")
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("Expect Error. got=%T (%+v)", evaluated, evaluated)
	}
	if err.Class.ReturnName() != UndefinedMethodError {
		t.Errorf("Expect %s. got=%T (%+v)", UndefinedMethodError, evaluated, evaluated)
	}
}

func TestArgumentError(t *testing.T) {
	evaluated := testEval(t, "[].count(5,4,3)")
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("Expect Error. got=%T (%+v)", evaluated, evaluated)
	}
	if err.Class.ReturnName() != ArgumentError {
		t.Errorf("Expect %s. got=%T (%+v)", ArgumentError, evaluated, evaluated)
	}
}

func TestTypeError(t *testing.T) {
	evaluated := testEval(t, "10 * \"foo\"")
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("Expect Error. got=%T (%+v)", evaluated, evaluated)
	}
	if err.Class.ReturnName() != TypeError {
		t.Errorf("Expect %s. got=%T (%+v)", TypeError, evaluated, evaluated)
	}
}
