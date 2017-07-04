package compiler

import (
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/parser"
	"fmt"
)

func CompileToBytecode(input string) (string, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return "", fmt.Errorf(err.Message)
	}
	g := bytecode.NewGenerator()
	g.InitTopLevelScope(program)
	return g.GenerateByteCode(program.Statements), nil
}