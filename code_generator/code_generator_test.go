package code_generator

import (
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"strings"
	"testing"
)

func TestBasicMethodDefineAndExecution(t *testing.T) {
	input := `
	def foo(x)
	  x + 10
	end

	foo(11)
	`

	expected := `
<Def:foo>
0 getlocal 0
1 putobject 10
2 opt_plus
3 leave
<ProgramStart>
0 putself
1 putstring foo
2 def_method 1
3 pop
4 putself
5 putobject 11
6 send foo
7 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestPopRedundantValue(t *testing.T) {
	input := `
	a = 10
	b = 11
	a + b # redundant value
	a
	`

	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0
2 putobject 11
3 setlocal 1
4 getlocal 0
5 getlocal 1
6 opt_plus
7 pop
8 getlocal 0
9 leave
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
2 opt_mult
3 putobject 100
4 opt_plus
5 putobject 2
6 opt_div
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
1 setlocal 0
2 putobject 100
3 setlocal 0
4 putobject 5
5 setlocal 1
6 getlocal 1
7 getlocal 0
8 opt_mult
9 putobject 100
10 opt_plus
11 putobject 2
12 opt_div
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
1 setlocal 0
2 putobject 5
3 setlocal 1
4 getlocal 0
5 getlocal 1
6 opt_gt
7 branchunless 11
8 putobject 10
9 setlocal 2
10 getlocal 2
11 putobject 1
12 opt_plus
13 leave
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
1 setlocal 0
2 putobject 5
3 setlocal 1
4 getlocal 0
5 getlocal 1
6 opt_gt
7 branchunless 11
8 putobject 10
9 setlocal 2
10 jump 13
11 putobject 5
12 setlocal 2
13 getlocal 2
14 putobject 1
15 opt_plus
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
	cg := New(program)
	return cg.GenerateByteCode(program)
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
