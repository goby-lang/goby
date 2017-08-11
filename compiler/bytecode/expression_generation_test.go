package bytecode

import (
	"strings"
	"testing"
)

func TestLocalVariableAccessCompilation(t *testing.T) {
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
2 pop
3 putobject 100
4 setlocal 0 0
5 pop
6 putobject 5
7 setlocal 0 1
8 pop
9 putself
10 send foo 0
11 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestIfExpressionWithoutAlternativeCompilation(t *testing.T) {
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
2 pop
3 putobject 5
4 setlocal 0 1
5 pop
6 getlocal 0 0
7 getlocal 0 1
8 send > 1
9 branchunless 13
10 putobject 10
11 setlocal 0 2
12 jump 14
13 putnil
14 pop
15 getlocal 0 2
16 putobject 1
17 send + 1
18 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestIfExpressionWithAlternativeCompilation(t *testing.T) {
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
2 pop
3 putobject 5
4 setlocal 0 1
5 pop
6 getlocal 0 0
7 getlocal 0 1
8 send > 1
9 branchunless 13
10 putobject 10
11 setlocal 0 2
12 jump 15
13 putobject 5
14 setlocal 0 2
15 pop
16 getlocal 0 2
17 putobject 1
18 send + 1
19 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestMultipleVariableAssignmentCompilation(t *testing.T) {
	input := `

	def foo
	  [1, 2, 3]
	end

	a, @b, c = foo
	`

	expected := `
<Def:foo>
0 putobject 1
1 putobject 2
2 putobject 3
3 newarray 3
4 leave
<ProgramStart>
0 putself
1 putstring "foo"
2 def_method 0
3 putself
4 send foo 0
5 expand_array 3
6 setlocal 0 0
7 setinstancevariable @b
8 setlocal 0 1
9 leave
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
2 pop
3 getconstant Foo false
4 setconstant Bar
5 pop
6 getconstant Foo false
7 getconstant Bar false
8 send + 1
9 leave
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
2 pop
3 putobject false
4 setlocal 0 1
5 pop
6 getlocal 0 0
7 send ! 0
8 getlocal 0 1
9 send == 1
10 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
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
5 pop
6 getlocal 0 0
7 putobject 0
8 putstring "foo"
9 send []= 2
10 pop
11 getlocal 0 0
12 putobject 0
13 send [] 1
14 setlocal 0 1
15 leave
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
6 pop
7 newhash 0
8 setlocal 0 1
9 pop
10 getlocal 0 1
11 putstring "baz"
12 getlocal 0 0
13 putstring "bar"
14 send [] 1
15 getlocal 0 0
16 putstring "foo"
17 send [] 1
18 send - 1
19 send []= 2
20 pop
21 getlocal 0 1
22 putstring "baz"
23 send [] 1
24 getlocal 0 0
25 putstring "bar"
26 send [] 1
27 send + 1
28 leave
`
	expected2 := `
<ProgramStart>
0 putstring "bar"
1 putobject 5
2 putstring "foo"
3 putobject 1
4 newhash 4
5 setlocal 0 0
6 pop
7 newhash 0
8 setlocal 0 1
9 pop
10 getlocal 0 1
11 putstring "baz"
12 getlocal 0 0
13 putstring "bar"
14 send [] 1
15 getlocal 0 0
16 putstring "foo"
17 send [] 1
18 send - 1
19 send []= 2
20 pop
21 getlocal 0 1
22 putstring "baz"
23 send [] 1
24 getlocal 0 0
25 putstring "bar"
26 send [] 1
27 send + 1
28 leave
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

%s
`, expected1, expected2, bytecode)
	}

}

func TestRangeCompilation(t *testing.T) {
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

func TestUnusedExpressionRemoval(t *testing.T) {
	input := `
	i = 0

	while i < 100 do
	  10
	  i += 1
	end

	i
	`

	expected := `
<ProgramStart>
0 putobject 0
1 setlocal 0 0
2 pop
3 jump 12
4 putnil
5 pop
6 jump 12
7 getlocal 0 0
8 putobject 1
9 send + 1
10 setlocal 0 0
11 pop
12 getlocal 0 0
13 putobject 100
14 send < 1
15 branchif 7
16 putnil
17 pop
18 getlocal 0 0
19 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
