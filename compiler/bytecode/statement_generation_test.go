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
10 pop
11 leave
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
1 putstring "bar"
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
4 getconstant Bar false
5 def_class class:Foo Bar
6 pop
7 getconstant Foo false
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
1 getconstant Bar false
2 send include 1
3 leave
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
3 jump 15
4 putnil
5 pop
6 jump 15
7 putself
8 getlocal 1 0
9 send puts 1
10 getlocal 1 0
11 putobject 1
12 send + 1
13 setlocal 1 0
14 pop
15 getlocal 1 0
16 putobject 1000
17 send <= 1
18 branchif 7
19 putnil
20 pop
21 leave
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
3 jump 12
4 putnil
5 pop
6 jump 12
7 getlocal 0 0
8 putobject 1
9 send - 1
10 setlocal 0 0
11 pop
12 getlocal 0 0
13 putobject 0
14 send > 1
15 branchif 7
16 putnil
17 pop
18 getlocal 0 0
19 leave
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
9 jump 18
10 putnil
11 pop
12 jump 18
13 getlocal 0 0
14 putobject 1
15 send - 1
16 setlocal 0 0
17 pop
18 getlocal 0 0
19 getlocal 0 1
20 send length 0
21 send > 1
22 branchif 13
23 putnil
24 pop
25 getlocal 0 0
26 leave
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
7 putnil
8 pop
9 jump 25
10 getlocal 0 0
11 putobject 1
12 send + 1
13 setlocal 0 0
14 pop
15 getlocal 0 0
16 putobject 5
17 send == 1
18 branchunless 20
19 jump 25
20 getlocal 0 1
21 putobject 1
22 send + 1
23 setlocal 0 1
24 pop
25 getlocal 0 0
26 putobject 10
27 send < 1
28 branchif 10
29 putnil
30 pop
31 leave
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
9 jump 45
10 putnil
11 pop
12 jump 45
13 getlocal 0 0
14 putobject 1
15 send + 1
16 setlocal 0 0
17 pop
18 jump 39
19 putnil
20 pop
21 jump 39
22 getlocal 0 1
23 putobject 1
24 send + 1
25 setlocal 0 1
26 pop
27 getlocal 0 1
28 putobject 3
29 send == 1
30 branchunless 32
31 jump 45
32 getlocal 0 2
33 getlocal 0 0
34 getlocal 0 1
35 send * 1
36 send + 1
37 setlocal 0 2
38 pop
39 getlocal 0 1
40 putobject 5
41 send < 1
42 branchif 22
43 putnil
44 pop
45 getlocal 0 0
46 putobject 10
47 send < 1
48 branchif 13
49 putnil
50 pop
51 getlocal 0 2
52 putobject 10
53 send * 1
54 setlocal 0 3
55 pop
56 getlocal 0 3
57 putobject 100
58 send + 1
59 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
