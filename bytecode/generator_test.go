package bytecode

import (
	"github.com/rooby-lang/Rooby/lexer"
	"github.com/rooby-lang/Rooby/parser"
	"strings"
	"testing"
)

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
1 getlocal 1 2
2 send + 1
3 setlocal 1 2
4 leave
<Block:0>
0 putobject 3
1 getlocal 0 1
2 send * 1
3 setlocal 1 1
4 getlocal 3 1
5 send bar 0 block:1
6 leave
<ProgramStart>
0 putself
1 def_class Foo
2 pop
3 putobject 100
4 setlocal 0 0
5 putobject 10
6 setlocal 1 0
7 putobject 1000
8 setlocal 2 0
9 getconstant Foo
10 send new 0
11 setlocal 3 0
12 getlocal 3 0
13 send bar 0 block:0
14 getlocal 1 0
15 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestCallBlockCompilation(t *testing.T) {
	input := `
def foo
  yield(20, 10)
end

x = 100

self.foo do |x, y|
  x - y
end
`
	expected := `
<Def:foo>
0 putself
1 putobject 20
2 putobject 10
3 invokeblock 2
4 leave
<Block:0>
0 getlocal 0 0
1 getlocal 1 0
2 send - 1
3 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
3 putobject 100
4 setlocal 0 0
5 putself
6 send foo 0 block:0
7 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestHashCompilation(t *testing.T) {
	input := `
	a = { foo: 1, bar: 5 }
	b = {}
	b["baz"] = a["bar"] - a["foo"]
	b["baz"] + a["bar"]
`

	expected1 := `
<ProgramStart>
0 putstring "foo"
1 putobject 1
2 putstring "bar"
3 putobject 5
4 newhash 4
5 setlocal 0 0
6 newhash 0
7 setlocal 1 0
8 getlocal 1 0
9 putstring "baz"
10 getlocal 0 0
11 putstring "bar"
12 send [] 1
13 getlocal 0 0
14 putstring "foo"
15 send [] 1
16 send - 1
17 send []= 2
18 getlocal 1 0
19 putstring "baz"
20 send [] 1
21 getlocal 0 0
22 putstring "bar"
23 send [] 1
24 send + 1
25 leave
`
	expected2 := `
<ProgramStart>
0 putstring "bar"
1 putobject 5
2 putstring "foo"
3 putobject 1
4 newhash 4
5 setlocal 0 0
6 newhash 0
7 setlocal 1 0
8 getlocal 1 0
9 putstring "baz"
10 getlocal 0 0
11 putstring "bar"
12 send [] 1
13 getlocal 0 0
14 putstring "foo"
15 send [] 1
16 send - 1
17 send []= 2
18 getlocal 1 0
19 putstring "baz"
20 send [] 1
21 getlocal 0 0
22 putstring "bar"
23 send [] 1
24 send + 1
25 leave
`
	bytecode := strings.TrimSpace(compileToBytecode(input))

	// This is because hash stores data using map.
	// And map's keys won't be sorted when running in for loop.
	// So we can get 2 possible results.
	expected1 = strings.TrimSpace(expected1)
	expected2 = strings.TrimSpace(expected2)
	if bytecode != expected1 && bytecode != expected2 {
		t.Fatalf(`
Bytecode compare failed
Expect:
"%s"

Or:

"%s"

Got:
"%s"
`, expected1, expected2, bytecode)
	}

}

func TestArrayCompilation(t *testing.T) {
	input := `
	a = [1, 2, "bar"]
	a[0] = "foo"
	c = a[0]
`

	expected := `
<ProgramStart>
0 putobject 1
1 putobject 2
2 putstring "bar"
3 newarray 3
4 setlocal 0 0
5 getlocal 0 0
6 putobject 0
7 putstring "foo"
8 send []= 2
9 getlocal 0 0
10 putobject 0
11 send [] 1
12 setlocal 1 0
13 leave
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
2 getlocal 1 0
3 setinstancevariable @y
4 getlocal 0 0
5 getlocal 1 0
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
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestClassMethodDefinition(t *testing.T) {
	input := `
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
1 setlocal 2 0
2 getlocal 0 0
3 getlocal 1 0
4 send - 1
5 getlocal 2 0
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
1 setlocal 0 0
2 putobject 100
3 setlocal 0 0
4 putobject 5
5 setlocal 1 0
6 getlocal 1 0
7 getlocal 0 0
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
1 setlocal 0 0
2 putobject 5
3 setlocal 1 0
4 getlocal 0 0
5 getlocal 1 0
6 send > 1
7 branchunless 10
8 putobject 10
9 setlocal 2 0
10 putnil
11 getlocal 2 0
12 putobject 1
13 send + 1
14 leave
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
1 setlocal 0 0
2 putobject 5
3 setlocal 1 0
4 getlocal 0 0
5 getlocal 1 0
6 send > 1
7 branchunless 11
8 putobject 10
9 setlocal 2 0
10 jump 13
11 putobject 5
12 setlocal 2 0
13 getlocal 2 0
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
	g := NewGenerator(program)
	return g.GenerateByteCode(program)
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
