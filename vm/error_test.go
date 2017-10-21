package vm

import (
	"fmt"
	"strings"
	"testing"
)

func TestErrorLineNumber(t *testing.T) {
	tests := []errorTestCase{
		{`a
		123
		`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>", 1, 1},
		{`class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		 a.z
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			5, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestUndefinedMethodError(t *testing.T) {
	tests := []errorTestCase{
		{`a`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>", 1, 1},
		{`class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			5, 1},
		{`class Foo
		   attr_reader("foo")
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			6, 1},
		{`class Foo
		  attr_reader("bar")
		end

		a = Foo.new
		a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>",
			6, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}

}

func TestUnsupportedMethodError(t *testing.T) {
	tests := []errorTestCase{
		{`String.new`, "UnsupportedMethodError: Unsupported Method #new for String", 1, 1},
		{`Integer.new`, "UnsupportedMethodError: Unsupported Method #new for Integer", 1, 1},
		{`Hash.new`, "UnsupportedMethodError: Unsupported Method #new for Hash", 1, 1},
		{`Array.new`, "UnsupportedMethodError: Unsupported Method #new for Array", 1, 1},
		{`Boolean.new`, "UnsupportedMethodError: Unsupported Method #new for Boolean", 1, 1},
		{`Null.new`, "UnsupportedMethodError: Unsupported Method #new for Null", 1, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
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
			4, 1},
		{`def foo(x)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4, 1},
		{`def foo(x = 10)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4, 1},
		{`def foo(x, y = 10)
		end

		foo(1, 2, 3)
		`,
			"ArgumentError: Expect at most 2 args for method 'foo'. got: 3",
			4, 1},
		{`"1234567890".include? "123", Class`,
			"ArgumentError: Expect 1 argument. got=2",
			1, 1},
		{`"1234567890".include? "123", Class, String`,
			"ArgumentError: Expect 1 argument. got=3",
			1, 1},
		{`def foo(a, *b)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4, 1},
		{`def foo(a, b, *c)
		end

		foo(10)
		`, "ArgumentError: Expect at least 2 args for method 'foo'. got: 1",
			4, 1},
		{`def foo(a, b = 10, *c)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4, 1},
		{`def foo(a, b, c)
		  a + b + c
		end

		arr = [1, 2, 3, 5]
		foo(*arr)
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			6, 1},
		{`def foo(a, b, c)
		  a + b + c
		end

		def bar
		  arr = [1, 2, 3, 5]
		  foo(*arr)
		end

		bar
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			6, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStackTraces(t *testing.T) {
	tests := []struct {
		input          string
		expectedMsg    string
		expectedTraces []string
	}{
		{`def foo(a, b, c)
		  a + b + c
		end

		def bar
		  arr = [1, 2, 3, 5]
		  foo(*arr)
		end

		bar
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			[]string{
				fmt.Sprintf("from %s:7", getFilename()),
				fmt.Sprintf("from %s:10", getFilename()),
			},
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expectedMsg)
		checkErrorTraces(t, i, evaluated, tt.expectedTraces)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestKeywordArgumentError(t *testing.T) {
	tests := []errorTestCase{
		{`def foo(x:)
		  x
		end

		foo
		`,
			"ArgumentError: Method foo requires key argument x",
			5, 1},
		{`def foo
		  10
		end

		foo(y: 1)
		`,
			"ArgumentError: Expect at most 0 args for method 'foo'. got: 1",
			5, 1},
		{`def foo(x)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo",
			5, 1},
		{`def foo(x = 10)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo",
			5, 1},
		{`def foo(x:)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: Method foo requires key argument x",
			5, 1},
		{`def foo(x: 10)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo",
			5, 1},
		{`def foo(x: 10)
		  x
		end

		foo(y: 1, x: 100)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			5, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConstantAlreadyInitializedError(t *testing.T) {
	tests := []errorTestCase{
		{`Foo = 10
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2, 1},
		{`class Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2, 1},
		{`module Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.",
			2, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func checkError(t *testing.T, index int, evaluated Object, expectedErrMsg, fn string, line int) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

	expectedErrMsg = fmt.Sprintf("%s. At %s:%d", expectedErrMsg, fn, line)
	if err.message != expectedErrMsg {
		t.Fatalf("At test case %d: Expect error message to be:\n  %s. got: \n%s", index, expectedErrMsg, err.Message())
	}
}

func checkErrorMsg(t *testing.T, index int, evaluated Object, expectedErrMsg string) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

	fmt.Println(err.Message())
	if err.message != expectedErrMsg {
		t.Fatalf("At test case %d: Expect error message to be:\n  %s. got: \n%s", index, expectedErrMsg, err.message)
	}
}

func checkErrorTraces(t *testing.T, index int, evaluated Object, expectedTraces []string) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

	joinedExpectedTraces := strings.Join(expectedTraces, "\n")
	joinedTraces := strings.Join(err.stackTraces, "\n")

	if joinedTraces != joinedTraces {
		t.Fatalf("At test case %d: Expect traces to be:\n  %s. got: \n%s", index, joinedExpectedTraces, joinedTraces)
	}
}
