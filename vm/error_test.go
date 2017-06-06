package vm

import "testing"

func TestUndefinedMethod(t *testing.T) {
	evaluated := testEval(t, "a")
	obj, ok := evaluated.(*UndefinedMethodErrorObject)
	if !ok {
		t.Errorf("Expect UndefinedMethodError. got=%T (%+v)", obj, obj)
	}
}
