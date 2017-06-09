package bytecode

import (
	"bytes"
	"github.com/goby-lang/goby/ast"
	"regexp"
	"strings"
)

type localTable struct {
	store map[string]int
	count int
	depth int
	upper *localTable
}

func (lt *localTable) get(v string) (int, bool) {
	i, ok := lt.store[v]

	return i, ok
}

func (lt *localTable) set(val string) int {
	c, ok := lt.store[val]

	if !ok {
		c = lt.count
		lt.store[val] = c
		lt.count++
		return c
	}

	return c
}

func (lt *localTable) setLCL(v string, d int) (index, depth int) {
	index, depth, ok := lt.getLCL(v, d)

	if !ok {
		index = lt.set(v)
		depth = lt.depth
		return index, depth
	}

	return index, depth
}

func (lt *localTable) getLCL(v string, d int) (index, depth int, ok bool) {
	index, ok = lt.get(v)

	if ok {
		return index, d - lt.depth, ok
	}

	if lt.upper != nil {
		index, depth, ok = lt.upper.getLCL(v, d)
		return
	}

	return -1, 0, false
}

func newLocalTable(depth int) *localTable {
	s := make(map[string]int)
	return &localTable{store: s, depth: depth}
}

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
