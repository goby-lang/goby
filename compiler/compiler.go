package compiler

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/bytecode"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/parser"
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

func CompileToInstructions(input string) ([]*bytecode.InstructionSet, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	g := bytecode.NewGenerator()
	g.InitTopLevelScope(program)
	return g.GenerateInstructions(program.Statements), nil
}
