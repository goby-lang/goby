package vm

import (
	"fmt"
	"github.com/goby-lang/goby/compiler"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
	"github.com/goby-lang/goby/vm/errors"
	"os"
	"runtime"
	"testing"
)

func TestVM_REPLExec(t *testing.T) {
	tests := []struct {
		inputs   []string
		expected interface{}
	}{
		{
			[]string{
				`
				a, b = [3, 6]
				a + b
				`,
			},
			9,
		}, {
			[]string{
				`
				a, _, c = [1, 2, 3]
				c
				`,
			},
			3,
		}, {
			[]string{`
def foo(x)
  yield(x + 10)
end
def bar(y)
  foo(y) do |f|
	yield(f)
  end
end
def baz(z)
  bar(z + 100) do |b|
	yield(b)
  end
end
a = 0
baz(100) do |b|
  a = b
end
a
			`,
				`
class Foo
  def bar
	100
  end
end
module Baz
  class Bar
	def bar
	  Foo.new.bar
	end
  end
end
				`,
				`
Baz::Bar.new.bar + a
`,
			},
			310,
		},
		{
			[]string{`

def foo
  123
end
`, `
foo
`},
			123},
		{
			[]string{
				`
class Foo
  def bar(x)
    x + 10
  end
end
`,
				`
Foo.new.bar(90)
`},
			100},
		{
			[]string{
				`
def foo
  123
end
`,
				`
foo
`,
				`
def foo
  345
end
`,

				`
foo
`,
			}, 345},
		{
			[]string{
				`raise ArgumentError`,
			}, "error",
		},
	}

	for i, test := range tests {
		v := initTestVM()
		v.InitForREPL()

		// Initialize parser, lexer is not important here
		p := parser.New(lexer.New(""))
		p.Mode = parser.REPLMode

		program, _ := p.ParseProgram()

		// Initialize code generator, and it will behavior a little different in REPL mode.
		g := bytecode.NewGenerator()
		g.REPL = true
		g.InitTopLevelScope(program)

		for _, input := range test.inputs {
			p := parser.New(lexer.New(input))
			p.Mode = parser.REPLMode

			program, _ := p.ParseProgram()
			sets := g.GenerateInstructions(program.Statements)

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()
		verifyExpected(t, i, evaluated, test.expected)
		// Because REPL should maintain a base call frame so that the whole program won't exit
		v.checkCFP(t, i, 1)
	}
}

func TestVM_REPLExecFail(t *testing.T) {

	tests := []struct {
		inputs   []string
		expected string
	}{
		{
			[]string{
				`raise ArgumentError`,
			},
			fmt.Sprintf("InternalError: '%s'", errors.ArgumentError),
		},
		{
			[]string{
				"NonExistentBuiltinMethod",
			},
			"NameError: uninitialized constant NonExistentBuiltinMethod",
		},
		{
			[]string{
				"Hash.notExist",
			},
			"UndefinedMethodError: Undefined Method 'notExist' for Hash",
		},
	}

	for i, test := range tests {
		v := initTestVM()
		v.InitForREPL()

		// Initialize parser, lexer is not important here
		p := parser.New(lexer.New(""))
		p.Mode = parser.REPLMode

		program, _ := p.ParseProgram()

		// Initialize code generator, and it will behavior a little different in REPL mode.
		g := bytecode.NewGenerator()
		g.REPL = true
		g.InitTopLevelScope(program)

		for _, input := range test.inputs {
			p := parser.New(lexer.New(input))
			p.Mode = parser.REPLMode

			// prevent parse errors from panicking tests
			program, Err := p.ParseProgram()
			if Err != nil {
				t.Fatalf("At case %d unexpected parse error %q", i, Err.Message)
			}
			sets := g.GenerateInstructions(program.Statements)

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()

		if evaluated.toString() != test.expected {
			t.Fatalf("At case %d expected %s got %s", i,
				test.expected, evaluated.toString())
		}

		// Because REPL should maintain a base call frame so that the whole program won't exit
		v.checkCFP(t, i, 1)
	}
}

func initTestVM() *VM {
	fn, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	v, err := New(fn, []string{})

	if err != nil {
		panic(err)
	}

	v.mode = TestMode
	return v
}

func (v *VM) testEval(t *testing.T, input, filepath string) Object {
	iss, err := compiler.CompileToInstructions(input, parser.TestMode)

	if err != nil {
		t.Helper()
		t.Errorf("Error when compiling input: %s", input)
		t.Fatal(err.Error())
	}

	v.ExecInstructions(iss, filepath)

	return v.mainThread.stack.top().Target
}

func (v *VM) checkCFP(t *testing.T, index, expectedCFP int) {
	t.Helper()
	if v.mainThread.cfp != expectedCFP {
		t.Errorf("At case %d expect main thread's cfp to be %d. got: %d", index, expectedCFP, v.mainThread.cfp)
	}
}

func (v *VM) checkSP(t *testing.T, index, expectedSp int) {
	t.Helper()
	if v.mainThread.sp != expectedSp {
		fmt.Println(v.mainThread.stack.inspect())
		t.Errorf("At case %d expect main thread's sp to be %d. got: %d", index, expectedSp, v.mainThread.sp)
	}

}

// Verification helpers

func verifyExpected(t *testing.T, i int, evaluated Object, expected interface{}) {
	t.Helper()
	if isError(evaluated) {
		t.Errorf("At test case %d: %s", i, evaluated.toString())
		return
	}

	switch expected := expected.(type) {
	case int:
		verifyIntegerObject(t, i, evaluated, expected)
	case float64:
		verifyFloatObject(t, i, evaluated, expected)
	case string:
		verifyStringObject(t, i, evaluated, expected)
	case bool:
		verifyBooleanObject(t, i, evaluated, expected)
	case []interface{}:
		verifyArrayObject(t, i, evaluated, expected)
	case nil:
		verifyNullObject(t, i, evaluated)
	default:
		t.Errorf("Unknown type %T at case %d", expected, i)
	}
}

func verifyIntegerObject(t *testing.T, i int, obj Object, expected int) bool {
	t.Helper()
	switch result := obj.(type) {
	case *IntegerObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Integer. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyFloatObject(t *testing.T, i int, obj Object, expected float64) bool {
	t.Helper()
	switch result := obj.(type) {
	case *FloatObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%f, got=%f", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Float. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyNullObject(t *testing.T, i int, obj Object) bool {
	t.Helper()
	switch result := obj.(type) {
	case *NullObject:
		return true
	case *Error:
		t.Errorf("At test case %d: %s", i, result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not NULL. got=%s (%+v)", i, obj.Class().Name, obj)
		return false
	}
}

func verifyStringObject(t *testing.T, i int, obj Object, expected string) bool {
	t.Helper()
	switch result := obj.(type) {
	case *StringObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%q, got=%q", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf(result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not String. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyBooleanObject(t *testing.T, i int, obj Object, expected bool) bool {
	t.Helper()
	switch result := obj.(type) {
	case *BooleanObject:
		if result.value != expected {
			t.Errorf("At test case %d: object has wrong value. expect=%t, got=%t", i, expected, result.value)
			return false
		}

		return true
	case *Error:
		t.Errorf(result.Message())
		return false
	default:
		t.Errorf("At test case %d: object is not Boolean. got=%s (%+v).", i, obj.Class().Name, obj)
		return false
	}
}

func verifyArrayObject(t *testing.T, index int, obj Object, expected []interface{}) bool {
	t.Helper()
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not Array. got=%s (%+v)", index, obj.Class().Name, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("Don't equals length of array. Expect %d, got=%d", len(expected), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		verifyExpected(t, index, result.Elements[i], expected[i])
	}

	return true
}

// Same as testHashObject(), but expects a ConcurrentArray.
func verifyConcurrentArrayObject(t *testing.T, index int, obj Object, expected []interface{}) bool {
	t.Helper()
	result, ok := obj.(*ConcurrentArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not ConcurrentArray. got=%s (%+v)", index, obj.Class().Name, obj)
		return false
	}

	if len(result.InternalArray.Elements) != len(expected) {
		t.Errorf("Don't equals length of array. Expect %d, got=%d", len(expected), len(result.InternalArray.Elements))
	}

	for i := 0; i < len(result.InternalArray.Elements); i++ {
		verifyExpected(t, i, result.InternalArray.Elements[i], expected[i])
	}

	return true
}

// Same as testHashObject(), but expects a ConcurrentHash.
//
func verifyConcurrentHashObject(t *testing.T, index int, objectResult Object, expected map[string]interface{}) bool {
	t.Helper()
	result, ok := objectResult.(*ConcurrentHashObject)

	if !ok {
		t.Errorf("At test case %d: result is not ConcurrentHash. got=%s", index, objectResult.Class().Name)
		return false
	}

	pairs := make(map[string]Object)

	iterator := func(key, value interface{}) bool {
		pairs[key.(string)] = value.(Object)
		return true
	}

	result.internalMap.Range(iterator)

	return _checkHashPairs(t, pairs, expected)
}

// Tests a Hash Object, with a few limitations:
//
// - the tested hash must be shallow (no nested objects as values);
// - the test hash must have strings as keys;
// - the error message won't mention the key - only the value.
//
// The second limitation is currently the only Hash format in Goby, anyway.
//
func verifyHashObject(t *testing.T, index int, objectResult Object, expected map[string]interface{}) bool {
	t.Helper()
	result, ok := objectResult.(*HashObject)

	if !ok {
		t.Errorf("At test case %d: result is not Hash. got=%s", index, objectResult.Class().Name)
		return false
	}

	return _checkHashPairs(t, result.Pairs, expected)
}

// Testing API like testArrayObject(), but performed on bidimensional arrays.
//
// Input example:
//
//		evaluated = '[["a", 1], ["b", "2"]]'
//		expected = [][]interface{}{{"a", 1}, {"b", "2"}}
//		testBidimensionalArrayObject(t, i, evaluated, expected)
//
func verifyBidimensionalArrayObject(t *testing.T, index int, obj Object, expected [][]interface{}) bool {
	t.Helper()
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("At test case %d: object is not Array. got=%T (%+v)", index, obj, obj)
		return false
	}

	if len(result.Elements) != len(expected) {
		t.Errorf("Unexpected result size. Expect %d, got=%d", len(expected), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		resultRow := result.Elements[i]
		expectedRow := expected[i]

		verifyArrayObject(t, index, resultRow, expectedRow)
	}

	return true
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}

// Internal helpers -----------------------------------------------------

func _checkHashPairs(t *testing.T, actual map[string]Object, expected map[string]interface{}) bool {
	if len(actual) != len(expected) {
		t.Errorf("Unexpected result size. Expected %d, got=%d", len(expected), len(actual))
	}

	for expectedKey, expectedValue := range expected {
		resultValue := actual[expectedKey]

		verifyExpected(t, i, resultValue, expectedValue)
	}

	return true
}

func getFilename() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}
