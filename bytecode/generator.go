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
	out        *scope
	localTable *localTable
	line       int
}

func newScope(s *scope, stmt ast.Statement) *scope {
	return &scope{out: s, localTable: newLocalTable(0), self: stmt, line: 0}
}

// Generator contains program's AST and will store generated instruction sets
type Generator struct {
	program         *ast.Program
	instructionSets []*instructionSet
	blockCounter    int
}

// NewGenerator initializes new Generator with complete AST tree.
func NewGenerator(program *ast.Program) *Generator {
	return &Generator{program: program}
}

// GenerateByteCode returns compiled bytecodes
func (g *Generator) GenerateByteCode(program *ast.Program) string {
	scope := &scope{program: program, localTable: newLocalTable(0)}
	g.compileStatements(program.Statements, scope, scope.localTable)
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
