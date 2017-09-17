package bytecode

import "testing"

func TestCustomClassConstructor(t *testing.T) {
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
2 pop
3 getlocal 0 1
4 setinstancevariable @y
5 pop
6 getlocal 0 0
7 getlocal 0 1
8 send - 1
9 setinstancevariable @z
10 leave
<Def:bar>
0 getinstancevariable @x
1 getinstancevariable @y
2 send + 1
3 getinstancevariable @z
4 send + 1
5 leave
<DefClass:Foo>
0 putself
1 putstring initialize
2 def_method 2
3 putself
4 putstring bar
5 def_method 0
6 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo false
4 putobject 100
5 putobject 50
6 send new 2
7 send bar 0
8 leave
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
1 putstring bar
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
3 getconstant Foo true
4 getconstant Bar true
5 getconstant Baz false
6 send new 0
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
1 putstring bar
2 def_singleton_method 0
3 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo false
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
1 putstring bar
2 def_method 0
3 leave
<DefClass:Foo>
0 leave
<ProgramStart>
0 putself
1 def_class class:Bar
2 pop
3 putself
4 getconstant Bar false
5 def_class class:Foo Bar
6 pop
7 pop
8 getconstant Foo false
9 send new 0
10 send bar 0
11 leave
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
1 putstring bar
2 def_method 0
3 leave
<DefClass:Foo>
0 putself
1 getconstant Bar false
2 send include 1
3 pop
4 leave
<ProgramStart>
0 putself
1 def_class module:Bar
2 pop
3 putself
4 def_class class:Foo
5 pop
6 getconstant Foo false
7 send new 0
8 send bar 0
9 leave
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
3 pop
4 jump 14
5 putself
6 getlocal 1 0
7 send puts 1
8 pop
9 getlocal 1 0
10 putobject 1
11 send + 1
12 setlocal 1 0
13 pop
14 getlocal 1 0
15 putobject 1000
16 send <= 1
17 branchif 5
18 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 pop
3 putself
4 send thread 0 block:0
5 leave
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
2 pop
3 jump 9
4 getlocal 0 0
5 putobject 1
6 send - 1
7 setlocal 0 0
8 pop
9 getlocal 0 0
10 putobject 0
11 send > 1
12 branchif 4
13 getlocal 0 0
14 leave
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
2 pop
3 putobject 1
4 putobject 2
5 putobject 3
6 newarray 3
7 setlocal 0 1
8 pop
9 jump 15
10 getlocal 0 0
11 putobject 1
12 send - 1
13 setlocal 0 0
14 pop
15 getlocal 0 0
16 getlocal 0 1
17 send length 0
18 send > 1
19 branchif 10
20 getlocal 0 0
21 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNextStatementCompilation(t *testing.T) {
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
2 pop
3 putobject 0
4 setlocal 0 1
5 pop
6 jump 25
7 getlocal 0 0
8 putobject 1
9 send + 1
10 setlocal 0 0
11 pop
12 getlocal 0 0
13 putobject 5
14 send == 1
15 branchunless 18
16 jump 25
17 jump 19
18 putnil
19 pop
20 getlocal 0 1
21 putobject 1
22 send + 1
23 setlocal 0 1
24 pop
25 getlocal 0 0
26 putobject 10
27 send < 1
28 branchif 7
29 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestBreakStatementCompilation(t *testing.T) {
	input := `
x = 0
y = 0
i = 0

while x < 10 do
  x = x + 1
  while y < 5 do
	y = y + 1

	if y == 3
	  break
	end

	i = i + x * y
  end
end

a = i * 10
a + 100
`
	expected := `
<ProgramStart>
0 putobject 0
1 setlocal 0 0
2 pop
3 putobject 0
4 setlocal 0 1
5 pop
6 putobject 0
7 setlocal 0 2
8 pop
9 jump 40
10 getlocal 0 0
11 putobject 1
12 send + 1
13 setlocal 0 0
14 pop
15 jump 36
16 getlocal 0 1
17 putobject 1
18 send + 1
19 setlocal 0 1
20 pop
21 getlocal 0 1
22 putobject 3
23 send == 1
24 branchunless 27
25 jump 40
26 jump 28
27 putnil
28 pop
29 getlocal 0 2
30 getlocal 0 0
31 getlocal 0 1
32 send * 1
33 send + 1
34 setlocal 0 2
35 pop
36 getlocal 0 1
37 putobject 5
38 send < 1
39 branchif 16
40 getlocal 0 0
41 putobject 10
42 send < 1
43 branchif 10
44 getlocal 0 2
45 putobject 10
46 send * 1
47 setlocal 0 3
48 pop
49 getlocal 0 3
50 putobject 100
51 send + 1
52 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
