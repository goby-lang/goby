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
<ProgramStart>
<Def:foo>
0 putobject 123
1 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
`, `
<ProgramStart>
0 putself
1 send foo 0
`},
			123},
		{
			[]string{
				`
<ProgramStart>
<Def:bar>
0 getlocal 0 0
1 putobject 10
2 send + 1
3 leave
<DefClass:Foo>
0 putself
1 putstring "bar"
2 def_method 1
3 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
`,
				`
<ProgramStart>
0 getconstant Foo
1 send new 0
2 putobject 90
3 send bar 1
`},
			100},
		{
			[]string{
				`
<Def:foo>
0 putobject 123
1 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
`,
				`
<ProgramStart>
0 putself
1 send foo 0
`,
				`
<Def:foo>
0 putobject 345
1 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
`,

				`
<ProgramStart>
0 putself
1 send foo 0
`,
			}, 345},
	}

	for _, test := range tests {
		v := New("./", []string{})
		v.InitForREPL()

		for _, input := range test.inputs {
			v.REPLExec(input)
		}

		evaluated := v.GetExecResult()
		checkExpected(t, evaluated, test.expected)
	}
}

func testEval(t *testing.T, input string) Object {
	is, err := compiler.CompileToInstructions(input)

	if err != nil {
		t.Fatal(err.Error())
	}

	v := New("./", []string{})
	v.ExecInstructions(is, "./")

	return v.mainThread.stack.top().Target
}

func testIntegerObject(t *testing.T, obj Object, expected int) bool {
	switch result := obj.(type) {
	case *IntegerObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Integer. got=%T (%+v).", obj, obj)
		return false
	}
}

func testNullObject(t *testing.T, obj Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	switch result := obj.(type) {
	case *StringObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%s, got=%s", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not String. got=%T (%+v).", obj, obj)
		return false
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	switch result := obj.(type) {
	case *BooleanObject:
		if result.Value != expected {
			t.Errorf("object has wrong value. expect=%d, got=%d", expected, result.Value)
			return false
		}

		return true
	case *Error:
		t.Error(result.Message)
		return false
	default:
		t.Errorf("object is not Boolean. got=%T (%+v).", obj, obj)
		return false
	}
}

func testArrayObject(t *testing.T, obj Object, expected *ArrayObject) bool {
	result, ok := obj.(*ArrayObject)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(result.Elements) != len(expected.Elements) {
		t.Fatalf("Don't equals length of array. Expect %d, got=%d", len(expected.Elements), len(result.Elements))
	}

	for i := 0; i < len(result.Elements); i++ {
		intObj, ok := expected.Elements[i].(*IntegerObject)
		if ok {
			testIntegerObject(t, result.Elements[i], intObj.Value)
			continue
		}
		str, ok := expected.Elements[i].(*StringObject)
		if ok {
			testStringObject(t, result.Elements[i], str.Value)
			continue
		}

		b, ok := expected.Elements[i].(*BooleanObject)
		if ok {
			testBooleanObject(t, result.Elements[i], b.Value)
			continue
		}

		t.Fatalf("object is wrong type %T", expected.Elements[i])
	}

	return true
}

func checkExpected(t *testing.T, evaluated Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, evaluated, expected)
	case string:
		testStringObject(t, evaluated, expected)
	case bool:
		testBooleanObject(t, evaluated, expected)
	case nil:
		_, ok := evaluated.(*NullObject)

		if !ok {
			t.Fatalf("expect result should be Null. got=%T", evaluated)
		}
	}
}

func isError(obj Object) bool {
	if obj != nil {
		_, ok := obj.(*Error)
		return ok
	}
	return false
}
