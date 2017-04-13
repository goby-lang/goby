package vm

import (
	"testing"
)

//func TestCallBlock(t *testing.T) {
//	input := `
//<Def:initialize>
//0 putself
//1 putself
//2 invokeblock 1
//3 leave
//<Def:color=>
//0 setinstancevariable @color
//1 leave
//<Def:color>
//0 getinstancevariable @color
//1 leave
//<Def:doors=>
//0 setinstancevariable @doors
//1 leave
//<Def:doors>
//0 getinstancevariable @doors
//1 leave
//<DefClass:Car>
//0 putself
//1 putstring "initialize"
//2 def_method 0
//3 putself
//4 putstring "color="
//5 def_method 1
//6 putself
//7 putstring "color"
//8 def_method 0
//9 putself
//10 putstring "doors="
//11 def_method 1
//12 putself
//13 putstring "doors"
//14 def_method 0
//15 leave
//<Block>
//0 getlocal 1
//1 putstring "Red"
//2 send color= 1
//3 getlocal 1
//4 putobject 4
//5 send doors= 1
//6 leave
//<ProgramStart>
//0 putself
//1 def_class Car
//2 pop
//3 getconstant Car
//4 send new 0 block
//5 setlocal 0
//6 putstring "My car's color is "
//7 getlocal 0
//8 send color 0
//9 send + 1
//10 putstring " and it's got "
//11 send + 1
//12 getlocal 0
//13 send doors 0
//14 send to_s 0
//15 send + 1
//16 putstring " doors."
//17 send + 1
//18 leave
//`
//	expected := "My car's color is Red and it's got 4 doors."
//	result := testExec(input).(*StringObject).Value
//	if result != expected {
//		t.Fatalf("Expect result to be %d. got=%d", expected, result)
//	}
//}

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
1 def_class Foo
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
1 def_class Bar
2 pop
3 putself
4 def_class Foo Bar
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
1 def_class Foo
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
1 def_class Foo
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

func testExec(bytecodes string) interface{} {
	p := NewBytecodeParser()
	v := New()
	p.VM = v
	p.Parse(bytecodes)
	cf := NewCallFrame(v.LabelTable[PROGRAM]["ProgramStart"][0])
	cf.Self = MainObj
	v.CallFrameStack.Push(cf)
	v.Exec()

	return v.Stack.Top()
}

