package vm

import (
	"github.com/goby-lang/goby/compiler"
	"testing"
)

func TestVM_REPLExec(t *testing.T) {
	tests := []struct {
		inputs   []string
		expected interface{}
	}{
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
	}

	for i, test := range tests {
		v := New("./", []string{})
		v.InitForREPL()

		for _, input := range test.inputs {
			sets, err := compiler.CompileToInstructions(input)

			if err != nil {
				t.Fatalf(err.Error())
			}

			v.REPLExec(sets)
		}

		evaluated := v.GetExecResult()
		checkExpected(t, i, evaluated, test.expected)
	}
}

func testEval(t *testing.T, input string) Object {
	iss, err := compiler.CompileToInstructions(input)

	if err != nil {
		t.Fatal(err.Error())
	}

	v := New("./", []string{})

	v.ExecInstructions(iss, "./")

	return v.mainThread.stack.top().Target
}

func testIntegerObject(t *testing.T, i int, obj Object, expected int) bool {
	switch result := obj.(type) {
	case *IntegerObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf("At test case %d: %s", i, result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not Integer. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testNullObject(t *testing.T, i int, obj Object) bool {
	switch result := obj.(type) {
	case *NullObject:
		return true
	case *Error:
		t.Fatalf("At test case %d: %s", i, result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not NULL. got=%T (%+v)", i, obj, obj)
		return false
	}
}

func testStringObject(t *testing.T, i int, obj Object, expected string) bool {
	switch result := obj.(type) {
	case *StringObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%s, got=%s", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf(result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not String. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testBooleanObject(t *testing.T, i int, obj Object, expected bool) bool {
	switch result := obj.(type) {
	case *BooleanObject:
		if result.Value != expected {
			t.Fatalf("At test case %d: object has wrong value. expect=%d, got=%d", i, expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Fatalf(result.Message)
		return false
	default:
		t.Fatalf("At test case %d: object is not Boolean. got=%T (%+v).", i, obj, obj)
		return false
	}
}

func testArrayObject(t *testing.T, index int, obj Object, expected *ArrayObject) bool {
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Fatalf("At test case %d: object is not Array. got=%T (%+v)", index, obj, obj)
		return false
	}

	if len(result.Elements) != len(expected.Elements) {
		t.Fatalf("Don't equals length of array. Expect %d, got=%d", len(expected.Elements), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		intObj, ok := expected.Elements[i].(*IntegerObject)
		if ok {
			testIntegerObject(t, index, result.Elements[i], intObj.Value)
			continue
		}
		str, ok := expected.Elements[i].(*StringObject)
		if ok {
			testStringObject(t, index, result.Elements[i], str.Value)
			continue
		}

		b, ok := expected.Elements[i].(*BooleanObject)
		if ok {
			testBooleanObject(t, index, result.Elements[i], b.Value)
			continue
		}
		t.Fatalf("At test case %d: object is wrong type %T", index, expected.Elements[i])
	}

	return true
}

func checkExpected(t *testing.T, i int, evaluated Object, expected interface{}) {
	if isError(evaluated) {
		t.Fatalf("At test case %d: %s", i, evaluated.toString())
		return
	}

	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, i, evaluated, expected)
	case string:
		testStringObject(t, i, evaluated, expected)
	case bool:
		testBooleanObject(t, i, evaluated, expected)
	case nil:
		testNullObject(t, i, evaluated)
	}
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}
