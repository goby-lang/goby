package code_generator

import (
	"github.com/st0012/Rooby/lexer"
	"github.com/st0012/Rooby/parser"
	"strings"
	"testing"
)

func TestArithmeticCompilation(t *testing.T) {
	input := `
	(1 * 10 + 100) / 2
	`

	expected := `<ProgramStart>
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
	expected := `<ProgramStart>
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
