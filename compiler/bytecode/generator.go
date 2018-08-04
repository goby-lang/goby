package bytecode

import (
	"github.com/goby-lang/goby/compiler/ast"
)

type scope struct {
	program    *ast.Program
	localTable *localTable
	anchors    map[string]*anchor
}

func newScope() *scope {
	return &scope{localTable: newLocalTable(0), anchors: make(map[string]*anchor)}
}

// Generator contains program's AST and will store generated instruction sets
type Generator struct {
	REPL                   bool
	instructionSets        []*InstructionSet
	blockCounter           int
	scope                  *scope
	instructionsWithAnchor []*Instruction
}

// NewGenerator initializes new Generator with complete AST tree.
func NewGenerator() *Generator {
	return &Generator{instructionsWithAnchor: []*Instruction{}}
}

// ResetInstructionSets clears generator's instruction sets
func (g *Generator) ResetInstructionSets() {
	g.instructionSets = []*InstructionSet{}
}

// InitTopLevelScope sets generator's scope with program node, which means it's the top level scope
func (g *Generator) InitTopLevelScope(program *ast.Program) {
	g.scope = &scope{program: program, localTable: newLocalTable(0), anchors: make(map[string]*anchor)}
}

// GenerateInstructions returns compiled instructions
func (g *Generator) GenerateInstructions(stmts []ast.Statement) []*InstructionSet {
	g.compileStatements(stmts, g.scope, g.scope.localTable)
	for _, i := range g.instructionsWithAnchor {
		i.Params[0] = i.AnchorLine()
	}
	//fmt.Println(g.instructionsToString())
	//fmt.Print()
	return g.instructionSets
}

func (g *Generator) compileCodeBlock(is *InstructionSet, stmt *ast.BlockStatement, scope *scope, table *localTable) {
	for _, s := range stmt.Statements {
		g.compileStatement(is, s, scope, table)
	}
}

func (g *Generator) endInstructions(is *InstructionSet, sourceLine int) {
	if g.REPL && is.name == Program {
		return
	}
	is.define(Leave, sourceLine)
}
