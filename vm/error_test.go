package vm

import (
	"fmt"
	"strings"
	"testing"
)

type errorTestCase struct {
	input       string
	expected    string
	expectedCFP int
}

// Error mechanism test

func TestStackTraces(t *testing.T) {
	tests := []struct {
		input          string
		expectedMsg    string
		expectedTraces []string
		expectedCFP    int
		expectedSP     int
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
			2,
			2,
		},
		{`def foo(a, b, c)
		  a + b + c
		end

		def bar
		  arr = [1, 2, 3, 5]
		  foo(*arr)
		end

		def baz
		  bar
		end

		baz
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			[]string{
				fmt.Sprintf("from %s:7", getFilename()),
				fmt.Sprintf("from %s:11", getFilename()),
				fmt.Sprintf("from %s:14", getFilename()),
			},
			3,
			3,
		},
		{`def foo
		  10
		end

		[1, 2, 3].each do |i|
		  foo(i)
		end
		`,
			"ArgumentError: Expect at most 0 args for method 'foo'. got: 1",
			[]string{
				fmt.Sprintf("from %s:6", getFilename()),
				fmt.Sprintf("from %s:5", getFilename()),
			},
			4,
			2,
		},
		/*
			TODO: This case should have these stack traces:
			from /Users/stanlow/projects/go/src/github.com/goby-lang/goby/vm/error_test.go:9
			from /Users/stanlow/projects/go/src/github.com/goby-lang/goby/vm/error_test.go:2
			from /Users/stanlow/projects/go/src/github.com/goby-lang/goby/vm/error_test.go:8

			But currently we haven't been able to trace to the `yield` keyword.
		*/
		{`def foo
		  yield(10)
		end

		def bar
		end

		foo do |ten|
		  bar(ten)
		end
		`,
			"ArgumentError: Expect at most 0 args for method 'bar'. got: 1",
			[]string{
				fmt.Sprintf("from %s:9", getFilename()),
				fmt.Sprintf("from %s:8", getFilename()),
			},
			4,
			// receiver(mainObject), receiver, argument 10, errorObject
			4,
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expectedMsg)
		checkErrorTraces(t, i, evaluated, tt.expectedTraces)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, tt.expectedSP)
	}
}

// Error types test

func TestUndefinedMethodError(t *testing.T) {
	tests := []errorTestCase{
		{`a`, "UndefinedMethodError: Undefined Method 'a' for <Instance of: Object>", 1},
		{`class Foo
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>", 1},
		{`class Foo
		   attr_reader("foo")
		 end

		 a = Foo.new
		 a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>", 1},
		{`class Foo
		  attr_reader("bar")
		end

		a = Foo.new
		a.bar = "fuz"
		`, "UndefinedMethodError: Undefined Method 'bar=' for <Instance of: Foo>", 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestArgumentError(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		errorLine   int
		expectedCFP int
		expectedSP  int
	}{
		{`def foo(x)
		end

		foo
		`,
			"ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4, 1, 1},
		{`def foo(x)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4, 1, 1},
		{`def foo(x = 10)
		end

		foo(1, 2)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2",
			4, 1, 1},
		{`def foo(x, y = 10)
		end

		foo(1, 2, 3)
		`,
			"ArgumentError: Expect at most 2 args for method 'foo'. got: 3",
			4, 1, 1},
		{`"1234567890".include? "123", Class`,
			"ArgumentError: Expect 1 argument. got=2",
			1, 1, 1},
		{`"1234567890".include? "123", Class, String`,
			"ArgumentError: Expect 1 argument. got=3",
			1, 1, 1},
		{`def foo(a, *b)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4, 1, 1},
		{`def foo(a, b, *c)
		end

		foo(10)
		`, "ArgumentError: Expect at least 2 args for method 'foo'. got: 1",
			4, 1, 1},
		{`def foo(a, b = 10, *c)
		end

		foo
		`, "ArgumentError: Expect at least 1 args for method 'foo'. got: 0",
			4, 1, 1},
		{`def foo(a, b, c)
		  a + b + c
		end

		arr = [1, 2, 3, 5]
		foo(*arr)
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			6, 1, 1},
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
			// The two objects on the stack would be:
			// - the receiver of bar, because that call haven't been finished
			// - the error object
			6, 2, 2},
		{`def foo(a, b, c)
		  a + b + c
		end

		def bar
		  arr = [1, 2, 3, 5]
		  foo(*arr)
		end

		def baz
		  bar
		end

		baz
		`,
			"ArgumentError: Expect at most 3 args for method 'foo'. got: 4",
			// The three objects on the stack would be:
			// - the receiver of baz, because that call haven't been finished
			// - the receiver of bar, because that call haven't been finished
			// - the error object
			6, 3, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, tt.expectedSP)
	}
}

func TestKeywordArgumentError(t *testing.T) {
	tests := []errorTestCase{
		{`def foo(x:)
		  x
		end

		foo
		`,
			"ArgumentError: Method foo requires key argument x", 1},
		{`def foo
		  10
		end

		foo(y: 1)
		`,
			"ArgumentError: Expect at most 0 args for method 'foo'. got: 1", 1},
		{`def foo(x)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo", 1},
		{`def foo(x = 10)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo", 1},
		{`def foo(x:)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: Method foo requires key argument x", 1},
		{`def foo(x: 10)
		  x
		end

		foo(y: 1)
		`,
			"ArgumentError: unknown key y for method foo", 1},
		{`def foo(x: 10)
		  x
		end

		foo(y: 1, x: 100)
		`,
			"ArgumentError: Expect at most 1 args for method 'foo'. got: 2", 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConstantAlreadyInitializedError(t *testing.T) {
	tests := []errorTestCase{
		{`Foo = 10
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.", 1},
		{`class Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.", 1},
		{`module Foo; end
		Foo = 100
		`, "ConstantAlreadyInitializedError: Constant Foo already been initialized. Can't assign value to a constant twice.", 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Error test helper methods

func checkErrorMsg(t *testing.T, index int, evaluated Object, expectedErrMsg string) {
	err, ok := evaluated.(*Error)
	if !ok {
		t.Fatalf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

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

	if joinedTraces != joinedExpectedTraces {
		t.Fatalf("At test case %d: Expect traces to be:\n%s \n got: \n%s", index, joinedExpectedTraces, joinedTraces)
	}
}
