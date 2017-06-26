package bytecode

import (
	"bytes"
	"github.com/goby-lang/goby/ast"
	"regexp"
	"strings"
)

type scope struct {
	self       ast.Statement
	program    *ast.Program
	localTable *localTable
	line       int
	anchor     *anchor
}

func newScope(s *scope, stmt ast.Statement) *scope {
	return &scope{localTable: newLocalTable(0), self: stmt, line: 0}
}

// Generator contains program's AST and will store generated instruction sets
type Generator struct {
	REPL            bool
	instructionSets []*instructionSet
	blockCounter    int
	scope           *scope
}

// NewGenerator initializes new Generator with complete AST tree.
func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) ResetInstructionSets() {
	g.instructionSets = []*instructionSet{}
}

// GenerateByteCode returns compiled bytecodes
func (g *Generator) GenerateByteCode(program *ast.Program, newScope bool) string {
	if newScope {
		g.scope = &scope{program: program, localTable: newLocalTable(0)}
	}

	g.compileStatements(program.Statements, g.scope, g.scope.localTable)
	var out bytes.Buffer

	for _, is := range g.instructionSets {
		out.WriteString(is.compile())
	}

	return strings.TrimSpace(removeEmptyLine(out.String()))
}

func (g *Generator) compileCodeBlock(is *instructionSet, stmt *ast.BlockStatement, scope *scope, table *localTable) {
	for _, s := range stmt.Statements {
		g.compileStatement(is, s, scope, table)
	}
}

func (g *Generator) endInstructions(is *instructionSet) {
	if g.REPL && is.label.Name == Program {
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
