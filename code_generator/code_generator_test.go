package code_generator

import (
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"strings"
	"testing"
)

func TestClassMethodDefinition(t *testing.T) {
	input :=`
class Foo
  def self.bar
    10
  end
end

Foo.bar
`
	expected := `
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

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestClassDefinition(t *testing.T) {
	input := `
class Bar
  def bar
    10
  end
end

class Foo < Bar
end

Foo.new.bar
`
	expected := `
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
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestBasicMethodReDefineAndExecution(t *testing.T) {
	input := `
	def foo(x)
	  x + 100
	end

	def foo(x)
	  x + 10
	end

	foo(11)
	`

	expected := `
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

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestBasicMethodDefineAndExecution(t *testing.T) {
	input := `
	def foo(x, y)
	  z = 10
	  x - y + z
	end

	foo(11, 1)
	`

	expected := `
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
4 putobject 1
5 putobject 11
6 send foo 2
7 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestArithmeticCompilation(t *testing.T) {
	input := `
	(1 * 10 + 100) / 2
	`

	expected := `
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

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestLocalVariableAccessInCurrentScope(t *testing.T) {
	input := `
	a = 10
	a = 100
	b = 5
	(b * a + 100) / 2
	`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0
2 putobject 100
3 setlocal 0
4 putobject 5
5 setlocal 1
6 getlocal 1
7 getlocal 0
8 send * 1
9 putobject 100
10 send + 1
11 putobject 2
12 send / 1
13 leave`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestConditionWithoutAlternativeCompilation(t *testing.T) {
	input := `
	a = 10
	b = 5
	if a > b
	  c = 10
	end

	c + 1
	`

	expected := `
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

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestConditionWithAlternativeCompilation(t *testing.T) {
	input := `
	a = 10
	b = 5
	if a > b
	  c = 10
	else
	  c = 5
	end

	c + 1
	`

	expected := `
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
10 jump 13
11 putobject 5
12 setlocal 2
13 getlocal 2
14 putobject 1
15 send + 1
16 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func compileToBytecode(input string) string {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()
	cg := New(program)
	return cg.GenerateByteCode(program)
}

func compareBytecode(t *testing.T, value, expected string) {
	value = strings.TrimSpace(value)
	expected = strings.TrimSpace(expected)
	if value != expected {
		t.Fatalf(`
Bytecode compare failed
Expect:
"%s"

Got:
"%s"
`, expected, value)
	}
}
