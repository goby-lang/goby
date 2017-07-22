package bytecode

import (
	"bytes"
	"strings"

	"github.com/goby-lang/goby/compiler/ast"
	"github.com/looplab/fsm"
)

type scope struct {
	self       ast.Statement
	program    *ast.Program
	localTable *localTable
	line       int
	anchor     *anchor
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
	fsm             *fsm.FSM
}

const (
	removeExp = "removeExp"
	keepExp   = "keepExp"
)

// NewGenerator initializes new Generator with complete AST tree.
func NewGenerator() *Generator {
	return &Generator{
		fsm: fsm.NewFSM(
			keepExp,

			/*
				This is for deciding if we should remove the expression.
				For example, these expression should be ignored when show up alone like:

				```
				a
				```

				```
				1 + a
				```

				```
				Foo
				```

				Because in these cases they are useless and will keep stack growing unnecessarily.

				Following expressions should be removed when declared but not used

				- Variable expressions like identifier, instance variable or constant
				- Data type expressions like string, integer, array...etc.
				- Self expression
				- Prefix expression like !true or -5
				- Not assignment infix expressions like: 1 + a * 5


				But only when those they are inside following places:
				- block argument
				- method definition
				- if expression's consequence or alternative block
				- while statement

				So if we know we are having those expressions in above places,
				we should switch the state to removeExp and compile function will ignore them.
			*/

			fsm.Events{
				{Name: removeExp, Src: []string{keepExp}, Dst: removeExp},
				{Name: keepExp, Src: []string{removeExp, keepExp}, Dst: keepExp},
			},
			fsm.Callbacks{},
		),
	}
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

	return strings.TrimSpace(strings.Replace(out.String(), "\n\n", "\n", -1))
}

// GenerateInstructions returns compiled instructions
func (g *Generator) GenerateInstructions(stmts []ast.Statement) []*InstructionSet {
	g.compileStatements(stmts, g.scope, g.scope.localTable)
	return g.instructionSets
}

func (g *Generator) compileCodeBlock(is *InstructionSet, stmt *ast.BlockStatement, scope *scope, table *localTable) {
	for i, s := range stmt.Statements {
		/*
			We shouldn't remove last expression since it would be the method's return value. Example:

			```
			def foo
			  10 <- should be removed
			  100 <- shouldn't be removed
			end
			```
		*/
		if i == len(stmt.Statements)-1 && g.fsm.Is(removeExp) {
			g.fsm.Event(keepExp)
			g.compileStatement(is, s, scope, table)
			g.fsm.Event(removeExp)
			continue
		}

		g.compileStatement(is, s, scope, table)
	}
}

func (g *Generator) endInstructions(is *InstructionSet) {
	if g.REPL && is.name == Program {
		return
	}
	is.define(Leave)
}
