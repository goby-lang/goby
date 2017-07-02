package bytecode

import (
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
	"strings"
	"testing"
)

func TestRequireRelativeCompilation(t *testing.T) {
	input := `
	require_relative("foo")

	Foo.bar
	`

	expected := `
<ProgramStart>
0 putself
1 putstring "foo"
2 send require_relative 1
3 getconstant Foo
4 send bar 0
5 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestRequireCompilation(t *testing.T) {
	input := `
	require("foo")

	Foo.bar
	`

	expected := `
<ProgramStart>
0 putself
1 putstring "foo"
2 send require 1
3 getconstant Foo
4 send bar 0
5 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNestedBlockCompilation(t *testing.T) {
	input := `
class Foo
  def bar
    yield
  end
end

a = 100
i = 10
b = 1000

f = Foo.new

f.bar do
  i = 3 * a
  f.bar do
    i = 3 + i
  end
end
i
`
	expected := `
<Def:bar>
0 putself
1 invokeblock 0
2 leave
<DefClass:Foo>
0 putself
1 putstring "bar"
2 def_method 0
3 leave
<Block:1>
0 putobject 3
1 getlocal 2 1
2 send + 1
3 setlocal 2 1
4 leave
<Block:0>
0 putobject 3
1 getlocal 1 0
2 send * 1
3 setlocal 1 1
4 getlocal 1 3
5 send bar 0 block:1
6 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 putobject 100
4 setlocal 0 0
5 putobject 10
6 setlocal 0 1
7 putobject 1000
8 setlocal 0 2
9 getconstant Foo
10 send new 0
11 setlocal 0 3
12 getlocal 0 3
13 send bar 0 block:0
14 getlocal 0 1
15 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestCallBlockCompilation(t *testing.T) {
	input := `
x = 1

foo do
  puts(x)
  y = 1
  puts(y)
  x = 2
  puts(x)
end

puts(x)
`
	expected := `
<Block:0>
0 putself
1 getlocal 1 0
2 send puts 1
3 putobject 1
4 setlocal 0 0
5 putself
6 getlocal 0 0
7 send puts 1
8 putobject 2
9 setlocal 1 0
10 putself
11 getlocal 1 0
12 send puts 1
13 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 putself
3 send foo 0 block:0
4 putself
5 getlocal 0 0
6 send puts 1
7 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestCumstomConstructor(t *testing.T) {
	input := `
class Foo
  def initialize(x, y)
    @x = x
    @y = y
    @z = x - y
  end

  def bar
    @x + @y + @z
  end
end

Foo.new(100, 50).bar
`

	expected := `
<Def:initialize>
0 getlocal 0 0
1 setinstancevariable @x
2 getlocal 0 1
3 setinstancevariable @y
4 getlocal 0 0
5 getlocal 0 1
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
0 getlocal 0 0
1 putobject 100
2 send + 1
3 leave
<Def:foo>
0 getlocal 0 0
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
1 setlocal 0 2
2 getlocal 0 0
3 getlocal 0 1
4 send - 1
5 getlocal 0 2
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

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestMethodDefWithDefaultValueArgument(t *testing.T) {
	input := `
	def foo(x, y = 10)
	  x + y
	end

	foo(100)
	`

	expected := `
<Def:foo>
0 putobject 10
1 setlocal 0 1 1
2 getlocal 0 0
3 getlocal 0 1
4 send + 1
5 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 2
3 putself
4 putobject 100
5 send foo 1
6 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func compileToBytecode(input string) string {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		panic(err.Message)
	}
	g := NewGenerator()
	g.InitTopLevelScope(program)
	return g.GenerateByteCode(program.Statements)
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
