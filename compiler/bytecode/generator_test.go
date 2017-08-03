package bytecode

import (
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
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
3 pop
4 pop
5 getconstant Foo false
6 send bar 0
7 pop
8 pop
9 leave
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
3 pop
4 pop
5 getconstant Foo false
6 send bar 0
7 pop
8 pop
9 leave
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
4 pop
5 getlocal 1 3
6 send bar 0 block:1
7 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 pop
4 putobject 100
5 setlocal 0 0
6 pop
7 pop
8 putobject 10
9 setlocal 0 1
10 pop
11 pop
12 putobject 1000
13 setlocal 0 2
14 pop
15 pop
16 getconstant Foo false
17 send new 0
18 setlocal 0 3
19 pop
20 pop
21 getlocal 0 3
22 send bar 0 block:0
23 pop
24 pop
25 getlocal 0 1
26 pop
27 pop
28 leave
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
3 pop
4 putobject 1
5 setlocal 0 0
6 pop
7 putself
8 getlocal 0 0
9 send puts 1
10 pop
11 putobject 2
12 setlocal 1 0
13 pop
14 putself
15 getlocal 1 0
16 send puts 1
17 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 pop
3 pop
4 putself
5 send foo 0 block:0
6 pop
7 pop
8 putself
9 getlocal 0 0
10 send puts 1
11 pop
12 pop
13 leave
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
7 pop
8 pop
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
3 pop
4 putself
5 putstring "foo"
6 def_method 1
7 pop
8 putself
9 putobject 11
10 send foo 1
11 pop
12 pop
13 leave
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
2 pop
3 getlocal 0 0
4 getlocal 0 1
5 send - 1
6 getlocal 0 2
7 send + 1
8 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 2
3 pop
4 putself
5 putobject 11
6 putobject 1
7 send foo 2
8 pop
9 pop
10 leave
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
3 pop
4 putself
5 putobject 100
6 send foo 1
7 pop
8 pop
9 leave
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
