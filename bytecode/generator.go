package bytecode

import (
	"bytes"
	"fmt"
	"github.com/rooby-lang/rooby/ast"
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

type scope struct {
	self       ast.Statement
	program    *ast.Program
	out        *scope
	localTable *localTable
	line       int
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

func (g *Generator) compileStatements(stmts []ast.Statement, scope *scope, table *localTable) {
	is := &instructionSet{label: &label{Name: Program}}

	for _, statement := range stmts {
		g.compileStatement(is, statement, scope, table)
	}

	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileStatement(is *instructionSet, statement ast.Statement, scope *scope, table *localTable) {
	scope.line++
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		g.compileExpression(is, stmt.Expression, scope, table)
	case *ast.DefStatement:
		is.define(PutSelf)
		is.define(PutString, fmt.Sprintf("\"%s\"", stmt.Name.Value))
		switch stmt.Receiver.(type) {
		case *ast.SelfExpression:
			is.define(DefSingletonMethod, len(stmt.Parameters))
		case nil:
			is.define(DefMethod, len(stmt.Parameters))
		}

		g.compileDefStmt(stmt, scope)
	case *ast.AssignStatement:
		g.compileAssignStmt(is, stmt, scope, table)
	case *ast.ClassStatement:
		is.define(PutSelf)

		if stmt.SuperClass != nil {
			is.define(DefClass, stmt.Name.Value, stmt.SuperClass.Value)
		} else {
			is.define(DefClass, stmt.Name.Value)
		}

		is.define(Pop)
		g.compileClassStmt(stmt, scope)
	case *ast.ReturnStatement:
		g.compileExpression(is, stmt.ReturnValue, scope, table)
		g.endInstructions(is)
	case *ast.RequireRelativeStatement:
		is.define(RequireRelative, stmt.Filepath)
	}
}

func (g *Generator) compileClassStmt(stmt *ast.ClassStatement, scope *scope) {
	scope = newScope(scope, stmt)
	is := &instructionSet{}
	is.setLabel(fmt.Sprintf("%s:%s", LabelDefClass, stmt.Name.Value))

	g.compileBlockStatement(is, stmt.Body, scope, scope.localTable)
	is.define(Leave)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileAssignStmt(is *instructionSet, stmt *ast.AssignStatement, scope *scope, table *localTable) {
	g.compileExpression(is, stmt.Value, scope, table)

	switch name := stmt.Name.(type) {
	case *ast.Identifier:
		index, depth := table.setLCL(name.Value, table.depth)
		is.define(SetLocal, index, depth)
	case *ast.InstanceVariable:
		is.define(SetInstanceVariable, name.Value)
	case *ast.Constant:
		is.define(SetConstant, name.Value)
	}
}

func (g *Generator) compileDefStmt(stmt *ast.DefStatement, scope *scope) {
	scope = newScope(scope, stmt)

	is := &instructionSet{}
	is.setLabel(fmt.Sprintf("%s:%s", LabelDef, stmt.Name.Value))

	for i := 0; i < len(stmt.Parameters); i++ {
		scope.localTable.setLCL(stmt.Parameters[i].Value, scope.localTable.depth)
	}

	g.compileBlockStatement(is, stmt.BlockStatement, scope, scope.localTable)
	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileExpression(is *instructionSet, exp ast.Expression, scope *scope, table *localTable) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		index, depth, ok := table.getLCL(exp.Value, table.depth)

		// it's local variable
		if ok {
			is.define(GetLocal, index, depth)
			return
		}

		// otherwise it's a method call
		is.define(PutSelf)
		is.define(Send, exp.Value, 0)

	case *ast.Constant:
		is.define(GetConstant, exp.Value)
	case *ast.InstanceVariable:
		is.define(GetInstanceVariable, exp.Value)
	case *ast.IntegerLiteral:
		is.define(PutObject, fmt.Sprint(exp.Value))
	case *ast.StringLiteral:
		is.define(PutString, fmt.Sprintf("\"%s\"", exp.Value))
	case *ast.BooleanExpression:
		is.define(PutObject, fmt.Sprint(exp.Value))
	case *ast.NilExpression:
		is.define(PutNull)
	case *ast.ArrayExpression:
		for _, elem := range exp.Elements {
			g.compileExpression(is, elem, scope, table)
		}
		is.define(NewArray, len(exp.Elements))
	case *ast.HashExpression:
		for key, value := range exp.Data {
			is.define(PutString, fmt.Sprintf("\"%s\"", key))
			g.compileExpression(is, value, scope, table)
		}
		is.define(NewHash, len(exp.Data)*2)
	case *ast.InfixExpression:
		g.compileInfixExpression(is, exp, scope, table)
	case *ast.PrefixExpression:
		switch exp.Operator {
		case "!":
			g.compileExpression(is, exp.Right, scope, table)
			is.define(Send, exp.Operator, 0)
		case "-":
			is.define(PutObject, 0)
			g.compileExpression(is, exp.Right, scope, table)
			is.define(Send, exp.Operator, 1)
		}

	case *ast.IfExpression:
		g.compileIfExpression(is, exp, scope, table)
	case *ast.SelfExpression:
		is.define(PutSelf)
	case *ast.YieldExpression:
		is.define(PutSelf)

		for _, arg := range exp.Arguments {
			g.compileExpression(is, arg, scope, table)
		}

		is.define(InvokeBlock, len(exp.Arguments))
	case *ast.CallExpression:
		g.compileExpression(is, exp.Receiver, scope, table)

		for _, arg := range exp.Arguments {
			g.compileExpression(is, arg, scope, table)
		}

		if exp.Block != nil {
			newTable := newLocalTable(table.depth + 1)
			newTable.upper = table
			blockIndex := g.blockCounter
			g.blockCounter++
			g.compileBlockArgExpression(blockIndex, exp, scope, newTable)
			is.define("send", exp.Method, len(exp.Arguments), fmt.Sprintf("block:%d", blockIndex))
			return
		}
		is.define("send", exp.Method, len(exp.Arguments))
	}
}

func (g *Generator) compileBlockArgExpression(index int, exp *ast.CallExpression, scope *scope, table *localTable) {
	is := &instructionSet{}
	is.setLabel(fmt.Sprintf("%s:%d", Block, index))

	for i := 0; i < len(exp.BlockArguments); i++ {
		table.set(exp.BlockArguments[i].Value)
	}

	g.compileBlockStatement(is, exp.Block, scope, table)
	g.endInstructions(is)
	g.instructionSets = append(g.instructionSets, is)
}

func (g *Generator) compileIfExpression(is *instructionSet, exp *ast.IfExpression, scope *scope, table *localTable) {
	g.compileExpression(is, exp.Condition, scope, table)

	anchor1 := &anchor{}
	is.define(BranchUnless, anchor1)

	g.compileBlockStatement(is, exp.Consequence, scope, table)

	anchor1.line = is.Count + 1

	if exp.Alternative == nil {
		anchor1.line--
		is.define(PutNull)
		return
	}

	anchor2 := &anchor{}
	is.define(Jump, anchor2)

	g.compileBlockStatement(is, exp.Alternative, scope, table)

	anchor2.line = is.Count
}

func (g *Generator) compileInfixExpression(is *instructionSet, node *ast.InfixExpression, scope *scope, table *localTable) {
	g.compileExpression(is, node.Left, scope, table)
	g.compileExpression(is, node.Right, scope, table)
	is.define(Send, node.Operator, "1")
}

func (g *Generator) compileBlockStatement(is *instructionSet, stmt *ast.BlockStatement, scope *scope, table *localTable) {
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

func newLocalTable(depth int) *localTable {
	s := make(map[string]int)
	return &localTable{store: s, depth: depth}
}

func newScope(s *scope, stmt ast.Statement) *scope {
	return &scope{out: s, localTable: newLocalTable(0), self: stmt, line: 0}
}
