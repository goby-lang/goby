package bytecode

import (
	"bytes"
	"github.com/goby-lang/goby/compiler/ast"
	"regexp"
	"strings"
)

type scope struct {
	self       ast.Statement
	program    *ast.Program
	localTable *localTable
	line       int
	anchor     *Anchor
}

func newScope(stmt ast.Statement) *scope {
	return &scope{localTable: newLocalTable(0), self: stmt, line: 0}
}

// Generator contains program's AST and will store generated instruction sets
type Generator struct {
	REPL            bool
	instructionSets []*InstructionSet
	blockCounter    int
	scope           *scope
}

// NewGenerator initializes new Generator with complete AST tree.
func NewGenerator() *Generator {
	return &Generator{}
}

// ResetInstructionSets clears generator's instruction sets
func (g *Generator) ResetInstructionSets() {
	g.instructionSets = []*InstructionSet{}
}

// InitTopLevelScope sets generator's scope with program node, which means it's the top level scope
func (g *Generator) InitTopLevelScope(program *ast.Program) {
	g.scope = &scope{program: program, localTable: newLocalTable(0)}
}

// GenerateByteCode returns compiled instructions in string format
func (g *Generator) GenerateByteCode(stmts []ast.Statement) string {
	g.compileStatements(stmts, g.scope, g.scope.localTable)
	var out bytes.Buffer

	for _, is := range g.instructionSets {
		out.WriteString(is.compile())
	}

	return strings.TrimSpace(removeEmptyLine(out.String()))
}

// GenerateInstructions returns compiled instructions
func (g *Generator) GenerateInstructions(stmts []ast.Statement) []*InstructionSet {
	g.compileStatements(stmts, g.scope, g.scope.localTable)
	return g.instructionSets
}

func (g *Generator) compileCodeBlock(is *InstructionSet, stmt *ast.BlockStatement, scope *scope, table *localTable) {
	for _, s := range stmt.Statements {
		g.compileStatement(is, s, scope, table)
	}
}

func (g *Generator) endInstructions(is *InstructionSet) {
	if g.REPL && is.Label.Name == Program {
		return
	}
	is.define(Leave)
}

func removeEmptyLine(s string) string {
	regex, err := regexp.Compile("\n+")
	if err != nil {
		panic(err)
	}
	s = regex.ReplaceAllString(s, "\n")

	return s
}
