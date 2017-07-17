package vm

import "testing"

func TestTypeError(t *testing.T) {
	vm := initTestVM()
	evaluated := vm.testEval(t, "10 * \"foo\"")
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("Expect Error. got=%T (%+v)", evaluated, evaluated)
	}
	if err.Class().ReturnName() != TypeError {
		t.Errorf("Expect %s. got=%T (%+v)", TypeError, evaluated, evaluated)
	}
	vm.checkCFP(t, 0, 1)
}

func TestUndefinedMethodError(t *testing.T) {
	tests := []struct {
		input    string
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

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, UndefinedMethodError, tt.errorMsg)
		vm.checkCFP(t, i, 1)
	}

}

func TestUnsupportedMethodError(t *testing.T) {
	tests := []struct {
		input    string
		errorMsg string
	}{
		{`String.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:String>"},
		{`Integer.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:Integer>"},
		{`Hash.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:Hash>"},
		{`Array.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:Array>"},
		{`Boolean.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:Boolean>"},
		{`Null.new`, "UnsupportedMethodError: Unsupported Method #new for <Class:Null>"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, UnsupportedMethodError, tt.errorMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestArgumentError(t *testing.T) {
	tests := []struct {
		input  string
		errMsg string
	}{
		{`
		def foo(x)
		end

		foo
		`,
			"ArgumentError: Expect at least 1 args for method 'foo'. got: 0"},
		{`
		def foo(x)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2"},
		{`
		def foo(x = 10)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2"},
		{`
		def foo(x, y = 10)
		end

		foo(1, 2, 3)
		`,
			"ArgumentError: Expect at most 2 args for method 'foo'. got: 3"},
		{`
		"1234567890".include "123", Class
		`,
			"ArgumentError: Expect 1 argument. got=2",
		},
		{`
		"1234567890".include "123", Class, String
		`,
			"ArgumentError: Expect 1 argument. got=3",
		},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, ArgumentError, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func checkError(t *testing.T, index int, evaluated Object, expectedErrType, expectedErrMsg string) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}
	if err.class.Name != expectedErrType {
		t.Errorf("At test case %d: Expect %s. got=%T (%+v)", index, expectedErrType, evaluated, err)
	}
	if err.Message != expectedErrMsg {
		t.Fatalf("At test case %d: Expect error message to be:\n  %s. got: \n%s", index, expectedErrMsg, err.Message)
	}
}
