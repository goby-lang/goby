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
	return lt.store[v]
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
		switch stmt.Expression.(type) {
		case *ast.IfExpression:
		default:
			is.Define("pop")
		}

	case *ast.DefStatement:
		cg.compileDefStmt(stmt, scope)
	case *ast.AssignStatement:
		cg.compileAssignStmt(is, stmt, scope)
	}
}

func (cg *CodeGenerator) compileAssignStmt(is *InstructionSet, stmt *ast.AssignStatement, scope *Scope) {
	n := stmt.Name.ReturnValue()
	index := scope.LocalTable.Set(n)
	cg.compileExpression(is, stmt.Value, scope)
	is.Define("setlocal", fmt.Sprint(index))
}

func (cg *CodeGenerator) compileDefStmt(stmt *ast.DefStatement, scope *Scope) {
	scope = newScope(scope, stmt)
	is := &InstructionSet{Label: &Label{Name: fmt.Sprintf("Def:%s", stmt.Name.Value)}}

	cg.compileBlockStatement(is, stmt.BlockStatement, scope)

	cg.endInstructions(is)
	cg.instructionSets = append(cg.instructionSets, is)
}

func (cg *CodeGenerator) compileExpression(is *InstructionSet, exp ast.Expression, scope *Scope) {
	switch exp := exp.(type) {
	case *ast.Identifier:
		value := fmt.Sprintf("%d", scope.LocalTable.Get(exp.Value))
		is.Define("getlocal", value)
	case *ast.IntegerLiteral:
		is.Define("putobject", fmt.Sprint(exp.Value))
	case *ast.StringLiteral:
		is.Define("putstring", exp.Value)
	case *ast.Boolean:
		is.Define("putobject", fmt.Sprint(exp.Value))
	case *ast.InfixExpression:
		cg.compileInfixExpression(is, exp, scope)
	case *ast.IfExpression:
		cg.compileIfExpression(is, exp, scope)
	}
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
	var operation string
	cg.compileExpression(is, node.Left, scope)
	cg.compileExpression(is, node.Right, scope)
	switch node.Operator {
	case "+":
		operation = "opt_plus"
	case "-":
		operation = "opt_minus"
	case "*":
		operation = "opt_mult"
	case "/":
		operation = "opt_div"
	case "==":
		operation = "opt_eq"
	case "<":
		operation = "opt_lt"
	case "<=":
		operation = "opt_le"
	case ">":
		operation = "opt_gl"
	case ">=":
		operation = "opt_ge"
	default:
		panic(fmt.Sprintf("Doesn't support %s operator", node.Operator))
	}
	is.Define(operation)
}

func (cg *CodeGenerator) compileBlockStatement(is *InstructionSet, stmt *ast.BlockStatement, scope *Scope) {
	for _, s := range stmt.Statements {
		cg.compileStatement(is, s, scope)
	}
}

func (cg *CodeGenerator) endInstructions(is *InstructionSet) {
	// if last instruction is pop, it should be replaced with leave
	if is.Instructions[len(is.Instructions)-1].Action == "pop" {
		is.Instructions = is.Instructions[:len(is.Instructions)-1]
		is.Count -= 1
	}
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

func newLocalTable() *LocalTable {
	s := make(map[string]int)
	return &LocalTable{store: s}
}

func newScope(scope *Scope, stmt ast.Statement) *Scope {
	return &Scope{Out: scope, LocalTable: newLocalTable(), Self: stmt, Line: 0}
}
