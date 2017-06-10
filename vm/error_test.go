package vm

import "testing"

func TestUndefinedMethodError(t *testing.T) {
	tests := []struct{
		input string
		errorMsg string
	}{
		{`a`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>"},
		{`
		 class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>"},
		{`
		 class Foo
		   attr_reader("foo")
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>"},
		{`
		class Foo
		  attr_reader("bar")
		end

		a = Foo.new
		a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)

		if !ok {
			t.Errorf("Expect Error. got=%T (%+v)", evaluated, evaluated)
		}

		if err.Class.Name != UndefinedMethodError {
			t.Errorf("Expect error to be %s. got=%s", UndefinedMethodError, err.Class.Name)
		}

		if err.Message != tt.errorMsg {
			t.Errorf("Expected error message: %s\nGot: %s\n", tt.errorMsg, err.Message)
		}
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
