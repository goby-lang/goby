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
3 pop
4 putobject 100
5 setlocal 0 0
6 pop
7 pop
8 putobject 5
9 setlocal 0 1
10 pop
11 pop
12 getlocal 0 1
13 getlocal 0 0
14 send * 1
15 putobject 100
16 send + 1
17 putobject 2
18 send / 1
19 pop
20 pop
21 putself
22 send foo 0
23 pop
24 pop
25 leave
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
3 pop
4 putobject 5
5 setlocal 0 1
6 pop
7 pop
8 getlocal 0 0
9 getlocal 0 1
10 send > 1
11 branchunless 15
12 putobject 10
13 setlocal 0 2
14 jump 16
15 putnil
16 pop
17 pop
18 getlocal 0 2
19 putobject 1
20 send + 1
21 pop
22 pop
23 leave
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
3 pop
4 putobject 5
5 setlocal 0 1
6 pop
7 pop
8 getlocal 0 0
9 getlocal 0 1
10 send > 1
11 branchunless 15
12 putobject 10
13 setlocal 0 2
14 jump 17
15 putobject 5
16 setlocal 0 2
17 pop
18 pop
19 getlocal 0 2
20 putobject 1
21 send + 1
22 pop
23 pop
24 leave
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
3 pop
4 putself
5 send foo 0
6 expand_array 3
7 setlocal 0 0
8 pop
9 setinstancevariable @b
10 pop
11 setlocal 0 1
12 pop
13 pop
14 leave
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
3 pop
4 getconstant Foo false
5 setconstant Bar
6 pop
7 pop
8 getconstant Foo false
9 getconstant Bar false
10 send + 1
11 pop
12 pop
13 leave
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
3 pop
4 putobject false
5 setlocal 0 1
6 pop
7 pop
8 getlocal 0 0
9 send ! 0
10 getlocal 0 1
11 send == 1
12 pop
13 pop
14 leave
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
6 pop
7 getlocal 0 0
8 putobject 0
9 putstring "foo"
10 send []= 2
11 pop
12 pop
13 getlocal 0 0
14 putobject 0
15 send [] 1
16 setlocal 0 1
17 pop
18 pop
19 leave
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
0 putstring "bar"
1 putobject 5
2 putstring "foo"
3 putobject 1
4 newhash 4
5 setlocal 0 0
6 pop
7 pop
8 newhash 0
9 setlocal 0 1
10 pop
11 pop
12 getlocal 0 1
13 putstring "baz"
14 getlocal 0 0
15 putstring "bar"
16 send [] 1
17 getlocal 0 0
18 putstring "foo"
19 send [] 1
20 send - 1
21 send []= 2
22 pop
23 pop
24 getlocal 0 1
25 putstring "baz"
26 send [] 1
27 getlocal 0 0
28 putstring "bar"
29 send [] 1
30 send + 1
31 pop
32 pop
33 leave
`
	expected2 := `
<ProgramStart>
0 putstring "foo"
1 putobject 1
2 putstring "bar"
3 putobject 5
4 newhash 4
5 setlocal 0 0
6 pop
7 pop
8 newhash 0
9 setlocal 0 1
10 pop
11 pop
12 getlocal 0 1
13 putstring "baz"
14 getlocal 0 0
15 putstring "bar"
16 send [] 1
17 getlocal 0 0
18 putstring "foo"
19 send [] 1
20 send - 1
21 send []= 2
22 pop
23 pop
24 getlocal 0 1
25 putstring "baz"
26 send [] 1
27 getlocal 0 0
28 putstring "bar"
29 send [] 1
30 send + 1
31 pop
32 pop
33 leave
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
6 pop
7 pop
8 leave
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
3 pop
4 jump 15
5 putnil
6 pop
7 jump 15
8 putobject 10
9 pop
10 getlocal 0 0
11 putobject 1
12 send + 1
13 setlocal 0 0
14 pop
15 getlocal 0 0
16 putobject 100
17 send < 1
18 branchif 8
19 putnil
20 pop
21 pop
22 getlocal 0 0
23 pop
24 pop
25 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
