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
3 pop
4 getconstant Foo false
5 putobject 100
6 putobject 50
7 send new 2
8 send bar 0
9 pop
10 pop
11 leave
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
3 pop
4 getconstant Foo true
5 getconstant Bar true
6 getconstant Baz false
7 send new 0
8 send bar 0
9 pop
10 pop
11 leave
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
3 pop
4 getconstant Foo false
5 send bar 0
6 pop
7 pop
8 leave
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
3 pop
4 putself
5 getconstant Bar false
6 def_class class:Foo Bar
7 pop
8 pop
9 getconstant Foo false
10 send new 0
11 send bar 0
12 pop
13 pop
14 leave
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
3 pop
4 leave
<ProgramStart>
0 putself
1 def_class module:Bar
2 pop
3 pop
4 putself
5 def_class class:Foo
6 pop
7 pop
8 getconstant Foo false
9 send new 0
10 send bar 0
11 pop
12 pop
13 leave
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
4 jump 17
5 putnil
6 pop
7 jump 17
8 putself
9 getlocal 1 0
10 send puts 1
11 pop
12 getlocal 1 0
13 putobject 1
14 send + 1
15 setlocal 1 0
16 pop
17 getlocal 1 0
18 putobject 1000
19 send <= 1
20 branchif 8
21 putnil
22 pop
23 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 pop
3 pop
4 putself
5 send thread 0 block:0
6 pop
7 pop
8 leave
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
3 pop
4 jump 13
5 putnil
6 pop
7 jump 13
8 getlocal 0 0
9 putobject 1
10 send - 1
11 setlocal 0 0
12 pop
13 getlocal 0 0
14 putobject 0
15 send > 1
16 branchif 8
17 putnil
18 pop
19 pop
20 getlocal 0 0
21 pop
22 pop
23 leave
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
3 pop
4 putobject 1
5 putobject 2
6 putobject 3
7 newarray 3
8 setlocal 0 1
9 pop
10 pop
11 jump 20
12 putnil
13 pop
14 jump 20
15 getlocal 0 0
16 putobject 1
17 send - 1
18 setlocal 0 0
19 pop
20 getlocal 0 0
21 getlocal 0 1
22 send length 0
23 send > 1
24 branchif 15
25 putnil
26 pop
27 pop
28 getlocal 0 0
29 pop
30 pop
31 leave
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
3 pop
4 putobject 0
5 setlocal 0 1
6 pop
7 pop
8 jump 30
9 putnil
10 pop
11 jump 30
12 getlocal 0 0
13 putobject 1
14 send + 1
15 setlocal 0 0
16 pop
17 getlocal 0 0
18 putobject 5
19 send == 1
20 branchunless 23
21 jump 30
22 jump 24
23 putnil
24 pop
25 getlocal 0 1
26 putobject 1
27 send + 1
28 setlocal 0 1
29 pop
30 getlocal 0 0
31 putobject 10
32 send < 1
33 branchif 12
34 putnil
35 pop
36 pop
37 leave
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
3 pop
4 putobject 0
5 setlocal 0 1
6 pop
7 pop
8 putobject 0
9 setlocal 0 2
10 pop
11 pop
12 jump 51
13 putnil
14 pop
15 jump 51
16 getlocal 0 0
17 putobject 1
18 send + 1
19 setlocal 0 0
20 pop
21 jump 45
22 putnil
23 pop
24 jump 45
25 getlocal 0 1
26 putobject 1
27 send + 1
28 setlocal 0 1
29 pop
30 getlocal 0 1
31 putobject 3
32 send == 1
33 branchunless 36
34 jump 51
35 jump 37
36 putnil
37 pop
38 getlocal 0 2
39 getlocal 0 0
40 getlocal 0 1
41 send * 1
42 send + 1
43 setlocal 0 2
44 pop
45 getlocal 0 1
46 putobject 5
47 send < 1
48 branchif 25
49 putnil
50 pop
51 getlocal 0 0
52 putobject 10
53 send < 1
54 branchif 16
55 putnil
56 pop
57 pop
58 getlocal 0 2
59 putobject 10
60 send * 1
61 setlocal 0 3
62 pop
63 pop
64 getlocal 0 3
65 putobject 100
66 send + 1
67 pop
68 pop
69 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
