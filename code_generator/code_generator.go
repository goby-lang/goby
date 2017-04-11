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
}

func (lt *LocalTable) Get(v string) int {
	i, ok := lt.store[v]

	if !ok {
		panic(fmt.Errorf("Can't find %s in local table.", v))
	}

	return i
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
	Self       ast.Statement
	Program    *ast.Program
	Out        *Scope
	LocalTable *LocalTable
	Line       int
}

type CodeGenerator struct {
	program         *ast.Program
	instructionSets []*InstructionSet
}

func New(program *ast.Program) *CodeGenerator {
	return &CodeGenerator{program: program}
}

func (cg *CodeGenerator) GenerateByteCode(program *ast.Program) string {
	scope := &Scope{Program: program, LocalTable: newLocalTable()}
	cg.compileProgram(program.Statements, scope)
	var out bytes.Buffer

	for _, is := range cg.instructionSets {
		out.WriteString(is.Compile())
	}
	return strings.TrimSpace(removeEmptyLine(out.String()))
}

func (cg *CodeGenerator) compileProgram(stmts []ast.Statement, scope *Scope) {
	label := &Label{Name: "ProgramStart"}
	is := &InstructionSet{Label: label}

	for _, statement := range stmts {
		cg.compileStatement(is, statement, scope)
	}

	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileStatement(is *InstructionSet, statement ast.Statement, scope *Scope) {
	scope.Line += 1
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		cg.compileExpression(is, stmt.Expression, scope)
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
		cg.compileAssignStmt(is, stmt, scope)
	case *ast.ClassStatement:
		is.Define("putself")

		if stmt.SuperClass != nil {
			is.Define("def_class", stmt.Name.Value, stmt.SuperClass.Value)
		} else {
			is.Define("def_class", stmt.Name.Value)
		}

		is.Define("pop")
		cg.compileClassStmt(stmt, scope)
	}
}

func (cg *CodeGenerator) compileClassStmt(stmt *ast.ClassStatement, scope *Scope) {
	scope = newScope(scope, stmt)
	is := &InstructionSet{}
	is.SetLabel(fmt.Sprintf("DefClass:%s", stmt.Name.Value))

	cg.compileBlockStatement(is, stmt.Body, scope)
	is.Define("leave")
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileAssignStmt(is *InstructionSet, stmt *ast.AssignStatement, scope *Scope) {
	switch name := stmt.Name.(type) {
	case *ast.Identifier:
		index := scope.LocalTable.Set(name.Value)
		cg.compileExpression(is, stmt.Value, scope)
		is.Define("setlocal", fmt.Sprint(index))
	case *ast.InstanceVariable:
		is.Define("setinstancevariable", name.Value)
	}
}

func (cg *CodeGenerator) compileDefStmt(stmt *ast.DefStatement, scope *Scope) {
	scope = newScope(scope, stmt)

	is := &InstructionSet{}
	is.SetLabel(fmt.Sprintf("Def:%s", stmt.Name.Value))

	for i := 0; i < len(stmt.Parameters); i++ {
		scope.LocalTable.Set(stmt.Parameters[i].Value)
	}

	cg.compileBlockStatement(is, stmt.BlockStatement, scope)
	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileExpression(is *InstructionSet, exp ast.Expression, scope *Scope) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		value := fmt.Sprintf("%d", scope.LocalTable.Get(exp.Value))
		is.Define("getlocal", value)
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
			cg.compileExpression(is, elem, scope)
		}
		is.Define("newarray", len(exp.Elements))
	case *ast.HashExpression:
		for key, value := range exp.Data {
			is.Define("putstring", fmt.Sprintf("\"%s\"", key))
			cg.compileExpression(is, value, scope)
		}
		is.Define("newhash", len(exp.Data)*2)
	case *ast.InfixExpression:
		cg.compileInfixExpression(is, exp, scope)
	case *ast.IfExpression:
		cg.compileIfExpression(is, exp, scope)
	case *ast.SelfExpression:
		is.Define("putself")
	case *ast.YieldExpression:
		is.Define("invokeblock")
		for i := len(exp.Arguments) - 1; i >= 0; i-- {
			cg.compileExpression(is, exp.Arguments[i], scope)
		}
	case *ast.CallExpression:
		cg.compileExpression(is, exp.Receiver, scope)

		for i := len(exp.Arguments) - 1; i >= 0; i-- {
			cg.compileExpression(is, exp.Arguments[i], scope)
		}

		if exp.Block != nil {
			cg.compileBlockArgExpression(exp, scope)
			is.Define("send", exp.Method, len(exp.Arguments), "block")
			return
		}
		is.Define("send", exp.Method, len(exp.Arguments))
	}
}

func (cg *CodeGenerator) compileBlockArgExpression(exp *ast.CallExpression, scope *Scope) {
	is := &InstructionSet{}
	is.SetLabel("Block")

	for i := 0; i < len(exp.BlockArguments); i++ {
		scope.LocalTable.Set(exp.BlockArguments[i].Value)
	}

	cg.compileBlockStatement(is, exp.Block, scope)
	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileIfExpression(is *InstructionSet, exp *ast.IfExpression, scope *Scope) {
	cg.compileExpression(is, exp.Condition, scope)

	anchor1 := &Anchor{}
	is.Define("branchunless", anchor1)

	cg.compileBlockStatement(is, exp.Consequence, scope)

	anchor1.Line = is.Count + 1

	if exp.Alternative == nil {
		return
	}

	anchor2 := &Anchor{}
	is.Define("jump", anchor2)

	cg.compileBlockStatement(is, exp.Alternative, scope)

	anchor2.Line = is.Count
}

func (cg *CodeGenerator) compileInfixExpression(is *InstructionSet, node *ast.InfixExpression, scope *Scope) {
	cg.compileExpression(is, node.Left, scope)
	cg.compileExpression(is, node.Right, scope)
	is.Define("send", node.Operator, "1")
}

func (cg *CodeGenerator) compileBlockStatement(is *InstructionSet, stmt *ast.BlockStatement, scope *Scope) {
	for _, s := range stmt.Statements {
		cg.compileStatement(is, s, scope)
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

func doubleQuoteString(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func newLocalTable() *LocalTable {
	s := make(map[string]int)
	return &LocalTable{store: s}
}

func newScope(scope *Scope, stmt ast.Statement) *Scope {
	return &Scope{Out: scope, LocalTable: newLocalTable(), Self: stmt, Line: 0}
}
