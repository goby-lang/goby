package vm

import (
	"testing"

	"github.com/goby-lang/goby/bytecode"
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
)

func TestCunstomConstructorAndInstanceVariable(t *testing.T) {
	input := `
<Def:initialize>
0 getlocal 0
1 setinstancevariable @x
2 getlocal 1
3 setinstancevariable @y
4 getlocal 0
5 getlocal 1
6 send - 1
7 setinstancevariable @z
8 leave
<Def:bar>
0 getinstancevariable @x
1 getinstancevariable @y
2 send + 1
3 getinstancevariable @z
4 send + 1
5 leave
<DefClass:Foo>
0 putself
1 putstring "initialize"
2 def_method 2
3 putself
4 putstring "bar"
5 def_method 0
6 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo
4 putobject 100
5 putobject 50
6 send new 2
7 send bar 0
8 leave
`
	expected := 200
	result := testExec(input).(*IntegerObject).Value
	if result != expected {
		t.Fatalf("Expect result to be %d. got=%d", expected, result)
	}
}

func TestCodeSectionOverrideIssue(t *testing.T) {
	input := `
<Def:foo>
0 putobject 60
1 leave
<Def:foo>
0 putobject 40
1 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
3 putself
4 send foo 0
5 setlocal 0
6 putself
7 putstring "foo"
8 def_method 0
9 putself
10 send foo 0
11 setlocal 1
12 putself
13 send foo 0
14 setlocal 2
15 getlocal 0
16 getlocal 1
17 send + 1
18 getlocal 2
19 send + 1
20 leave
`
	expected := 140
	result := testExec(input).(*IntegerObject).Value
	if result != expected {
		t.Fatalf("Expect result to be %d. got=%d", expected, result)
	}
}

func TestHash(t *testing.T) {
	input := `
<ProgramStart>
0 putstring "foo"
1 putobject 1
2 putstring "bar"
3 putobject 5
4 newhash 4
5 setlocal 0
6 newhash 0
7 setlocal 1
8 getlocal 1
9 putstring "baz"
10 getlocal 0
11 putstring "bar"
12 send [] 1
13 getlocal 0
14 putstring "foo"
15 send [] 1
16 send - 1
17 send []= 2
18 getlocal 1
19 putstring "baz"
20 send [] 1
21 getlocal 0
22 putstring "bar"
23 send [] 1
24 send + 1
25 leave
`
	expected := 9
	result := testExec(input).(*IntegerObject).Value
	if result != expected {
		t.Fatalf("Expect result to be %d. got=%d", expected, result)
	}
}

func TestArrayCompilation(t *testing.T) {
	input := `
<ProgramStart>
0 putobject 1
1 putobject 2
2 putstring "bar"
3 newarray 3
4 setlocal 0
5 getlocal 0
6 putobject 0
7 putstring "foo"
8 send []= 2
9 getlocal 0
10 putobject 0
11 send [] 1
12 setlocal 1
13 leave
`
	expected := "foo"
	result := testExec(input).(*StringObject).Value
	if result != expected {
		t.Fatalf("Expect result to be %s. got=%s", expected, result)
	}
}

func TestParsingComplexString(t *testing.T) {
	input := `
<ProgramStart>
0 putstring "Hello World! "
1 leave
	`

	result := testExec(input).(*StringObject).Value
	if result != "Hello World! " {
		t.Fatalf("Expect result to be \"Hello World! \". got=%s", result)
	}
}

func TestBuiltInInstanceMethod1(t *testing.T) {
	input := `
<ProgramStart>
0 putstring "123"
1 send class 0
2 send name 0
3 putstring ", "
4 send + 1
5 putobject 123
6 send class 0
7 send name 0
8 send + 1
9 putstring ", "
10 send + 1
11 putobject true
12 send class 0
13 send name 0
14 send + 1
15 leave
`
	expected := "String, Integer, Boolean"
	result := testExec(input).(*StringObject).Value
	if result != expected {
		t.Fatalf("Expect result to be %s. got=%s", expected, result)
	}
}

func TestClassDefinitionWithInheritance(t *testing.T) {
	input := `
<Def:bar>
0 putobject 10
1 leave
<DefClass:Bar>
0 putself
1 putstring "bar"
2 def_method 0
3 leave
<DefClass:Foo>
0 leave
<ProgramStart>
0 putself
1 def_class class:Bar
2 pop
3 putself
4 def_class class:Foo Bar
5 pop
6 getconstant Foo
7 send new 0
8 send bar 0
9 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 10 {
		t.Fatalf("Expect result to be 10. got=%d", result)
	}
}

func TestClassMethodDefinition(t *testing.T) {
	input := `
<Def:bar>
0 putobject 10
1 leave
<DefClass:Foo>
0 putself
1 putstring "bar"
2 def_singleton_method 0
3 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo
4 send bar 0
5 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 10 {
		t.Fatalf("Expect result to be 10. got=%d", result)
	}
}

func TestClassDefinition(t *testing.T) {
	input := `
<Def:bar>
0 putobject 11
1 leave
<DefClass:Foo>
0 putself
1 putstring "bar"
2 def_method 0
3 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo
4 send new 0
5 send bar 0
6 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 11 {
		t.Fatalf("Expect result to be 11. got=%d", result)
	}
}

func TestBasicMethodReDefineAndExecution(t *testing.T) {
	input := `
<Def:foo>
0 getlocal 0
1 putobject 100
2 send + 1
3 leave
<Def:foo>
0 getlocal 0
1 putobject 10
2 send + 1
3 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 1
3 putself
4 putstring "foo"
5 def_method 1
6 putself
7 putobject 11
8 send foo 1
9 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 21 {
		t.Fatalf("Expect result to be 21. got=%d", result)
	}
}

func TestBasicMethodDefineAndExecution(t *testing.T) {
	input := `
<Def:foo>
0 putobject 10
1 setlocal 2
2 getlocal 0
3 getlocal 1
4 send - 1
5 getlocal 2
6 send + 1
7 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 2
3 putself
4 putobject 11
5 putobject 1
6 send foo 2
7 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 20 {
		t.Fatalf("Expect result to be 20. got=%d", result)
	}
}

func TestArithmeticCalculation(t *testing.T) {
	input := `
<ProgramStart>
0 putobject 1
1 putobject 10
2 send * 1
3 putobject 100
4 send + 1
5 putobject 2
6 send / 1
7 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 55 {
		t.Fatalf("Expect result to be 55. got=%d", result)
	}
}

func TestConditionWithAlternativeCompilation(t *testing.T) {
	input := `
<ProgramStart>
0 putobject 10
1 setlocal 0
2 putobject 5
3 setlocal 1
4 getlocal 0
5 getlocal 1
6 send > 1
7 branchunless 11
8 putobject 10
9 setlocal 2
10 getlocal 2
11 putobject 1
12 send + 1
13 leave
`
	result := testExec(input).(*IntegerObject).Value
	if result != 11 {
		t.Fatalf("Expect result to be 11. got=%d", result)
	}
}

func testEval(t *testing.T, input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	g := bytecode.NewGenerator(program)
	bytecodes := g.GenerateByteCode(program)
	return testExec(bytecodes)
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testExec(bytecodes string) Object {
	v := New("./", []string{})
	v.ExecBytecodes(bytecodes, "./")

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
