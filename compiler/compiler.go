package compiler

import (
	"fmt"

	"github.com/gooby-lang/gooby/compiler/bytecode"
	"github.com/gooby-lang/gooby/compiler/lexer"
	"github.com/gooby-lang/gooby/compiler/parser"
)

// CompileToInstructions compiles input source code into instruction set data structures
func CompileToInstructions(input string, pm parser.ParserMode) ([]*bytecode.InstructionSet, error) {
	l := lexer.New(input)
	p := parser.New(l)
	p.Mode = pm
	program, err := p.ParseProgram()
	if err != nil {
		return nil, fmt.Errorf(err.Message)
	}
	g := bytecode.NewGenerator()
	g.InitTopLevelScope(program)
	return g.GenerateInstructions(program.Statements), nil
}
