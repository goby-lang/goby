package bytecode

import (
	"github.com/goby-lang/goby/lexer"
	"github.com/goby-lang/goby/parser"
	"strings"
	"testing"
)

func TestRangeExpression(t *testing.T) {
	input := `
	(1..(1+4)).each do |i|
	  puts(i)
	end
	`

	expected := `
<Block:0>
0 putself
1 getlocal 0 0
2 send puts 1
3 leave
<ProgramStart>
0 putobject 1
1 putobject 1
2 putobject 4
3 send + 1
4 newrange 0
5 send each 0 block:0
6 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestWhileStatementInBlock(t *testing.T) {
	input := `
	i = 1
	thread do
	  puts(i)
	  while i <= 1000 do
		puts(i)
		i = i + 1
	  end
	end
	`

	expected := `
<Block:0>
0 putself
1 getlocal 1 0
2 send puts 1
3 jump 14
4 putnil
5 pop
6 jump 14
7 putself
8 getlocal 1 0
9 send puts 1
10 getlocal 1 0
11 putobject 1
12 send + 1
13 setlocal 1 0
14 getlocal 1 0
15 putobject 1000
16 send <= 1
17 branchif 7
18 putnil
19 pop
20 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 putself
3 send thread 0 block:0
4 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNextStatement(t *testing.T) {
	input := `
	x = 0
	y = 0

	while x < 10 do
	  x = x + 1
	  if x == 5
	    next
	  end
	  y = y + 1
	end
	`

	expected := `
<ProgramStart>
0 putobject 0
1 setlocal 0 0
2 putobject 0
3 setlocal 0 1
4 jump 22
5 putnil
6 pop
7 jump 22
8 getlocal 0 0
9 putobject 1
10 send + 1
11 setlocal 0 0
12 getlocal 0 0
13 putobject 5
14 send == 1
15 branchunless 17
16 jump 22
17 putnil
18 getlocal 0 1
19 putobject 1
20 send + 1
21 setlocal 0 1
22 getlocal 0 0
23 putobject 10
24 send < 1
25 branchif 8
26 putnil
27 pop
28 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNamespacedClass(t *testing.T) {
	input := `
	module Foo
	  class Bar
	    class Baz
	      def bar
	      end
	    end
	  end
	end

	Foo::Bar::Baz.new.bar
	`

	expected := `
<Def:bar>
0 putnil
1 leave
<DefClass:Baz>
0 putself
1 putstring "bar"
2 def_method 0
3 leave
<DefClass:Bar>
0 putself
1 def_class class:Baz
2 pop
3 leave
<DefClass:Foo>
0 putself
1 def_class class:Bar
2 pop
3 leave
<ProgramStart>
0 putself
1 def_class module:Foo
2 pop
3 getconstant Foo
4 getconstant Bar
5 getconstant Baz
6 send new 0
7 send bar 0
8 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

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
7 setlocal 0 1
8 getlocal 0 1
9 putstring "baz"
10 getlocal 0 0
11 putstring "bar"
12 send [] 1
13 getlocal 0 0
14 putstring "foo"
15 send [] 1
16 send - 1
17 send []= 2
18 getlocal 0 1
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
7 setlocal 0 1
8 getlocal 0 1
9 putstring "baz"
10 getlocal 0 0
11 putstring "bar"
12 send [] 1
13 getlocal 0 0
14 putstring "foo"
15 send [] 1
16 send - 1
17 send []= 2
18 getlocal 0 1
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
12 setlocal 0 1
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
1 def_class class:Foo
2 pop
3 getconstant Foo
4 send bar 0
5 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestClassCompilation(t *testing.T) {
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
1 def_class class:Bar
2 pop
3 putself
4 getconstant Bar
5 def_class class:Foo Bar
6 pop
7 getconstant Foo
8 send new 0
9 send bar 0
10 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestModuleCompilation(t *testing.T) {
	input := `


module Bar
  def bar
    10
  end
end

class Foo
  include(Bar)
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
0 putself
1 getconstant Bar
2 send include 1
3 leave
<ProgramStart>
0 putself
1 def_class module:Bar
2 pop
3 putself
4 def_class class:Foo
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
	foo # This should be a method lookup
	`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0 0
2 putobject 100
3 setlocal 0 0
4 putobject 5
5 setlocal 0 1
6 getlocal 0 1
7 getlocal 0 0
8 send * 1
9 putobject 100
10 send + 1
11 putobject 2
12 send / 1
13 putself
14 send foo 0
15 leave`

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
3 setlocal 0 1
4 getlocal 0 0
5 getlocal 0 1
6 send > 1
7 branchunless 10
8 putobject 10
9 setlocal 0 2
10 putnil
11 getlocal 0 2
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
3 setlocal 0 1
4 getlocal 0 0
5 getlocal 0 1
6 send > 1
7 branchunless 11
8 putobject 10
9 setlocal 0 2
10 jump 13
11 putobject 5
12 setlocal 0 2
13 getlocal 0 2
14 putobject 1
15 send + 1
16 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestWhileStatementWithoutMethodCallInCondition(t *testing.T) {
	input := `
	i = 10

	while i > 0 do
	  i = i - 1
	end

	i
`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0 0
2 jump 10
3 putnil
4 pop
5 jump 10
6 getlocal 0 0
7 putobject 1
8 send - 1
9 setlocal 0 0
10 getlocal 0 0
11 putobject 0
12 send > 1
13 branchif 6
14 putnil
15 pop
16 getlocal 0 0
17 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestWhileStatementWithMethodCallInCondition(t *testing.T) {
	input := `
	i = 10
	a = [1, 2, 3]

	while i > a.length do
	  i = i - 1
	end

	i
`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0 0
2 putobject 1
3 putobject 2
4 putobject 3
5 newarray 3
6 setlocal 0 1
7 jump 15
8 putnil
9 pop
10 jump 15
11 getlocal 0 0
12 putobject 1
13 send - 1
14 setlocal 0 0
15 getlocal 0 0
16 getlocal 0 1
17 send length 0
18 send > 1
19 branchif 11
20 putnil
21 pop
22 getlocal 0 0
23 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestConstantCompilation(t *testing.T) {
	input := `
	Foo = 10
	Bar = Foo
	Foo + Bar
	`

	expected := `
<ProgramStart>
0 putobject 10
1 setconstant Foo
2 getconstant Foo
3 setconstant Bar
4 getconstant Foo
5 getconstant Bar
6 send + 1
7 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestBooleanCompilation(t *testing.T) {
	input := `
	a = true
	b = false
	!a == b
`
	expected := `
<ProgramStart>
0 putobject true
1 setlocal 0 0
2 putobject false
3 setlocal 0 1
4 getlocal 0 0
5 send ! 0
6 getlocal 0 1
7 send == 1
8 leave
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
