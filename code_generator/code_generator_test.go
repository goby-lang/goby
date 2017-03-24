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

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	p.CheckErrors()

	bytecodes := GenerateByteCode(program)
	compareBytecode(t, bytecodes, expected)
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
