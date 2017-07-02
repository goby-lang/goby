package bytecode

import (
	"testing"
	"strings"
)

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