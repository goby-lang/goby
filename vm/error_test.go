package vm

import (
	"fmt"
	"testing"
)

func TestErrorLineNumber(t *testing.T) {
	tests := []errorTestCase{
		{`a
		123
		`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>", 1},
		{`class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		 a.z
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			5},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestUndefinedMethodError(t *testing.T) {
	tests := []errorTestCase{
		{`a`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>", 1},
		{`class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			5},
		{`class Foo
		   attr_reader("foo")
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			6},
		{`class Foo
		  attr_reader("bar")
		end

		a = Foo.new
		a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			6},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}

}

func TestUnsupportedMethodError(t *testing.T) {
	tests := []errorTestCase{
		{`String.new`, "UnsupportedMethodError: Unsupported Method #new for String", 1},
		{`Integer.new`, "UnsupportedMethodError: Unsupported Method #new for Integer", 1},
		{`Hash.new`, "UnsupportedMethodError: Unsupported Method #new for Hash", 1},
		{`Array.new`, "UnsupportedMethodError: Unsupported Method #new for Array", 1},
		{`Boolean.new`, "UnsupportedMethodError: Unsupported Method #new for Boolean", 1},
		{`Null.new`, "UnsupportedMethodError: Unsupported Method #new for Null", 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestArgumentError(t *testing.T) {
	tests := []errorTestCase{
		{`def foo(x)
		end

		foo
		`,
			"ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4},
		{`def foo(x)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4},
		{`def foo(x = 10)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4},
		{`def foo(x, y = 10)
		end

		foo(1, 2, 3)
		`,
			"ArgumentError: Expect at most 2 args for method 'foo'. got: 3",
			4},
		{`"1234567890".include? "123", Class`,
			"ArgumentError: Expect 1 argument. got=2",
			1},
		{`"1234567890".include? "123", Class, String`,
			"ArgumentError: Expect 1 argument. got=3",
			1},
		{`def foo(a, *b)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4},
		{`def foo(a, b, *c)
		end

		foo(10)
		`, "ArgumentError: Expect at least 2 args for method 'foo'. got: 1",
			4},
		{`def foo(a, b = 10, *c)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4},
		{`def foo(a, b, c)
		  a + b + c
		end

		arr = [1, 2, 3, 5]
		foo(*arr)
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			6},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestConstantAlreadyInitializedError(t *testing.T) {
	tests := []errorTestCase{
		{`Foo = 10
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2},
		{`class Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2},
		{`module Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func checkError(t *testing.T, index int, evaluated Object, expectedErrMsg, fn string, line int) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Errorf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

	expectedErrMsg = fmt.Sprintf("%s. At %s:%d", expectedErrMsg, fn, line)
	if err.Message != expectedErrMsg {
		t.Fatalf("At test case %d: Expect error message to be:\n  %s. got: \n%s", index, expectedErrMsg, err.Message)
	}
}
