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

	expected := `
<ProgramStart>
putobject 1
putobject 10
opt_mult
putobject 100
opt_plus
putobject 2
opt_div
leave
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
putobject 10
setlocal 0
putobject 100
setlocal 0
putobject 5
setlocal 1
getlocal 1
getlocal 0
opt_mult
putobject 100
opt_plus
putobject 2
opt_div
leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestConditionCompilation(t *testing.T) {
	input := `
	a = 10
	b = 5
	if a > b
	  a
	else
	  b
	end
	`

	expected := `
putobject 10
setlocal 0
putobject 5
setlocal 1
getlocal 0
getlocal 1
opt_gl
branchunless

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
%s

Got:
%s
`, expected, value)
	}
}
