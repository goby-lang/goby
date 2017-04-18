package code_generator

import (
	"bytes"
	"fmt"
	"github.com/st0012/Rooby/ast"
	"regexp"
	"strings"
)

type LocalTable struct {
	store map[string]int
	count int
	depth int
	upper *LocalTable
}

func (lt *LocalTable) Get(v string) (int, bool) {
	i, ok := lt.store[v]

	return i, ok
}

func (lt *LocalTable) Set(val string) int {
	c, ok := lt.store[val]

	if !ok {
		c = lt.count
		lt.store[val] = c
		lt.count += 1
		return c
	}

	return c
}

type Scope struct {
	Self        ast.Statement
	Program     *ast.Program
	Out         *Scope
	LocalTable  *LocalTable
	Line        int
}

func (lt *LocalTable) setLCL(v string, d int) (index, depth int) {
	index, depth, ok := lt.getLCL(v, d)

	if !ok {
		index = lt.Set(v)
		depth = lt.depth
		return index, depth
	}

	return index, depth
}

func (lt *LocalTable) getLCL(v string, d int) (index, depth int, ok bool) {
	index, ok = lt.Get(v)

	if ok {
		return index, d - lt.depth, ok
	}

	if lt.upper != nil {
		index, depth, ok = lt.upper.getLCL(v, d)
		return
	}

	return -1, 0, false
}

type CodeGenerator struct {
	program         *ast.Program
	instructionSets []*InstructionSet
	blockCounter    int
}

// Initialize new CodeGenerator with complete AST tree.
func New(program *ast.Program) *CodeGenerator {
	return &CodeGenerator{program: program}
}

// Return compiled bytecodes
func (cg *CodeGenerator) GenerateByteCode(program *ast.Program) string {
	scope := &Scope{Program: program, LocalTable: newLocalTable(0)}
	cg.compileStatements(program.Statements, scope, scope.LocalTable)
	var out bytes.Buffer

	for _, is := range cg.instructionSets {
		out.WriteString(is.Compile())
	}

	return strings.TrimSpace(removeEmptyLine(out.String()))
}

func (cg *CodeGenerator) compileStatements(stmts []ast.Statement, scope *Scope, table *LocalTable) {
	is := &InstructionSet{Label: &Label{Name: "ProgramStart"}}

	for _, statement := range stmts {
		cg.compileStatement(is, statement, scope, table)
	}

	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileStatement(is *InstructionSet, statement ast.Statement, scope *Scope, table *LocalTable) {
	scope.Line += 1
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		cg.compileExpression(is, stmt.Expression, scope, table)
	case *ast.DefStatement:
		is.Define("putself")
		is.Define("putstring", fmt.Sprintf("\"%s\"", stmt.Name.Value))
		switch stmt.Receiver.(type) {
		case *ast.SelfExpression:
			is.Define("def_singleton_method", len(stmt.Parameters))
		case nil:
			is.Define("def_method", len(stmt.Parameters))
		}

		cg.compileDefStmt(stmt, scope)
	case *ast.AssignStatement:
		cg.compileAssignStmt(is, stmt, scope, table)
	case *ast.ClassStatement:
		is.Define("putself")

		if stmt.SuperClass != nil {
			is.Define("def_class", stmt.Name.Value, stmt.SuperClass.Value)
		} else {
			is.Define("def_class", stmt.Name.Value)
		}

		is.Define("pop")
		cg.compileClassStmt(stmt, scope)
	case *ast.ReturnStatement:
		cg.compileExpression(is, stmt.ReturnValue, scope, table)
		cg.endInstructions(is)
	}
}

func (cg *CodeGenerator) compileClassStmt(stmt *ast.ClassStatement, scope *Scope) {
	scope = newScope(scope, stmt)
	is := &InstructionSet{}
	is.SetLabel(fmt.Sprintf("DefClass:%s", stmt.Name.Value))

	cg.compileBlockStatement(is, stmt.Body, scope, scope.LocalTable)
	is.Define("leave")
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileAssignStmt(is *InstructionSet, stmt *ast.AssignStatement, scope *Scope, table *LocalTable) {
	cg.compileExpression(is, stmt.Value, scope, table)

	switch name := stmt.Name.(type) {
	case *ast.Identifier:
		index, depth := table.setLCL(name.Value, table.depth)
		is.Define("setlocal", index, depth)
	case *ast.InstanceVariable:
		is.Define("setinstancevariable", name.Value)
	case *ast.Constant:
		is.Define("setconstant", name.Value)
	}
}

func (cg *CodeGenerator) compileDefStmt(stmt *ast.DefStatement, scope *Scope) {
	scope = newScope(scope, stmt)

	is := &InstructionSet{}
	is.SetLabel(fmt.Sprintf("Def:%s", stmt.Name.Value))

	for i := 0; i < len(stmt.Parameters); i++ {
		scope.LocalTable.setLCL(stmt.Parameters[i].Value, scope.LocalTable.depth)
	}

	cg.compileBlockStatement(is, stmt.BlockStatement, scope, scope.LocalTable)
	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileExpression(is *InstructionSet, exp ast.Expression, scope *Scope, table *LocalTable) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		index, depth, ok := table.getLCL(exp.Value, table.depth)

		// it's local variable
		if ok {
			is.Define("getlocal", index, depth)
			return
		}

		// otherwise it's a method call
		is.Define("putself")
		is.Define("send", exp.Value, 0)

	case *ast.Constant:
		is.Define("getconstant", exp.Value)
	case *ast.InstanceVariable:
		is.Define("getinstancevariable", exp.Value)
	case *ast.IntegerLiteral:
		is.Define("putobject", fmt.Sprint(exp.Value))
	case *ast.StringLiteral:
		is.Define("putstring", fmt.Sprintf("\"%s\"", exp.Value))
	case *ast.Boolean:
		is.Define("putobject", fmt.Sprint(exp.Value))
	case *ast.ArrayExpression:
		for _, elem := range exp.Elements {
			cg.compileExpression(is, elem, scope, table)
		}
		is.Define("newarray", len(exp.Elements))
	case *ast.HashExpression:
		for key, value := range exp.Data {
			is.Define("putstring", fmt.Sprintf("\"%s\"", key))
			cg.compileExpression(is, value, scope, table)
		}
		is.Define("newhash", len(exp.Data)*2)
	case *ast.InfixExpression:
		cg.compileInfixExpression(is, exp, scope, table)
	case *ast.PrefixExpression:
		switch exp.Operator {
		case "!":
			cg.compileExpression(is, exp.Right, scope, table)
			is.Define("send", exp.Operator, 0)
		case "-":
			is.Define("putobject", 0)
			cg.compileExpression(is, exp.Right, scope, table)
			is.Define("send", exp.Operator, 1)
		}

	case *ast.IfExpression:
		cg.compileIfExpression(is, exp, scope, table)
	case *ast.SelfExpression:
		is.Define("putself")
	case *ast.YieldExpression:
		is.Define("putself")

		for _, arg := range exp.Arguments {
			cg.compileExpression(is, arg, scope, table)
		}

		is.Define("invokeblock", len(exp.Arguments))
	case *ast.CallExpression:
		cg.compileExpression(is, exp.Receiver, scope, table)

		for _, arg := range exp.Arguments {
			cg.compileExpression(is, arg, scope, table)
		}

		if exp.Block != nil {
			newTable := newLocalTable(table.depth+1)
			newTable.upper = table
			blockIndex := cg.blockCounter
			cg.blockCounter += 1
			cg.compileBlockArgExpression(blockIndex, exp, scope, newTable)
			is.Define("send", exp.Method, len(exp.Arguments), fmt.Sprintf("block:%d", blockIndex))
			return
		}
		is.Define("send", exp.Method, len(exp.Arguments))
	}
}

func (cg *CodeGenerator) compileBlockArgExpression(index int, exp *ast.CallExpression, scope *Scope, table *LocalTable) {
	is := &InstructionSet{}
	is.SetLabel(fmt.Sprintf("Block:%d", index))

	for i := 0; i < len(exp.BlockArguments); i++ {
		table.Set(exp.BlockArguments[i].Value)
	}

	cg.compileBlockStatement(is, exp.Block, scope, table)
	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileIfExpression(is *InstructionSet, exp *ast.IfExpression, scope *Scope, table *LocalTable) {
	cg.compileExpression(is, exp.Condition, scope, table)

	anchor1 := &Anchor{}
	is.Define("branchunless", anchor1)

	cg.compileBlockStatement(is, exp.Consequence, scope, table)

	anchor1.Line = is.Count + 1

	if exp.Alternative == nil {
		anchor1.Line -= 1
		is.Define("putnil")
		return
	}

	anchor2 := &Anchor{}
	is.Define("jump", anchor2)

	cg.compileBlockStatement(is, exp.Alternative, scope, table)

	anchor2.Line = is.Count
}

func (cg *CodeGenerator) compileInfixExpression(is *InstructionSet, node *ast.InfixExpression, scope *Scope, table *LocalTable) {
	cg.compileExpression(is, node.Left, scope, table)
	cg.compileExpression(is, node.Right, scope, table)
	is.Define("send", node.Operator, "1")
}

func (cg *CodeGenerator) compileBlockStatement(is *InstructionSet, stmt *ast.BlockStatement, scope *Scope, table *LocalTable) {
	for _, s := range stmt.Statements {
		cg.compileStatement(is, s, scope, table)
	}
}

func (cg *CodeGenerator) endInstructions(is *InstructionSet) {
	is.Define("leave")
}

func removeEmptyLine(s string) string {
	regex, err := regexp.Compile("\n+")
	if err != nil {
		panic(err)
	}
	s = regex.ReplaceAllString(s, "\n")

	return s
}

func newLocalTable(depth int) *LocalTable {
	s := make(map[string]int)
	return &LocalTable{store: s, depth: depth}
}

func newScope(scope *Scope, stmt ast.Statement) *Scope {
	return &Scope{Out: scope, LocalTable: newLocalTable(0), Self: stmt, Line: 0}
}
