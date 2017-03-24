package code_generator

import (
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
}

type CodeGenerator struct {
	program        *ast.Program
	instructionSet []string
}

func New(program *ast.Program) *CodeGenerator {
	return &CodeGenerator{program: program}
}

func (cg *CodeGenerator) GenerateByteCode(program *ast.Program) string {
	scope := &Scope{Program: program, LocalTable: newLocalTable()}
	cg.compileProgram(program.Statements, scope)
	instructions := strings.Join(cg.instructionSet, "\n")
	return strings.TrimSpace(removeEmptyLine(instructions))
}

func (cg *CodeGenerator) compileProgram(stmts []ast.Statement, scope *Scope) {
	var result []string

	result = append(result, "<ProgramStart>")
	for _, statement := range stmts {
		result = append(result, cg.compileStatement(statement, scope))
	}

	result = append(result, "leave")
	cg.instructionSet = append(cg.instructionSet, strings.Join(result, "\n"))
}

func (cg *CodeGenerator) compileStatement(statement ast.Statement, scope *Scope) string {
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		return cg.compileExpression(stmt.Expression, scope)
	case *ast.DefStatement:
		s := newScope(scope, stmt)
		cg.compileDefStmt(stmt, s)
		return ""
	case *ast.AssignStatement:
		return cg.compileAssignStmt(stmt, scope)
	}

	return ""
}

func (cg *CodeGenerator) compileAssignStmt(stmt *ast.AssignStatement, scope *Scope) string {
	n := stmt.Name.ReturnValue()
	index := scope.LocalTable.Set(n)
	result := fmt.Sprintf(`
%s
setlocal %d
`, cg.compileExpression(stmt.Value, scope), index)
	return result
}

func (cg *CodeGenerator) compileDefStmt(stmt *ast.DefStatement, scope *Scope) {
	var result []string

	result = append(result, fmt.Sprintf("<Def:%s>", stmt.Name.Value))

	for _, s := range stmt.BlockStatement.Statements {
		result = append(result, cg.compileStatement(s, scope))
	}
	result = append(result, "leave")
	cg.instructionSet = append(cg.instructionSet, strings.Join(result, "\n"))
}

func (cg *CodeGenerator) compileExpression(exp ast.Expression, scope *Scope) string {
	switch exp := exp.(type) {
	case *ast.Identifier:
		return fmt.Sprintf(`getlocal %d`, scope.LocalTable.Get(exp.Value))
	case *ast.IntegerLiteral:
		return fmt.Sprintf(`putobject %d`, exp.Value)
	case *ast.InfixExpression:
		return cg.compileInfixExpression(exp, scope)
	}

	return ""
}

func (cg *CodeGenerator) compileInfixExpression(node *ast.InfixExpression, scope *Scope) string {
	var operation string
	left := cg.compileExpression(node.Left, scope)
	right := cg.compileExpression(node.Right, scope)
	switch node.Operator {
	case "+":
		operation = "opt_plus"
	case "-":
		operation = "opt_minus"
	case "*":
		operation = "opt_mult"
	case "/":
		operation = "opt_div"
	default:
		panic(fmt.Sprintf("Doesn't support %s operator", node.Operator))
	}
	return fmt.Sprintf(`
%s
%s
%s
`, left, right, operation)
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
	return &Scope{Out: scope, LocalTable: newLocalTable(), Self: stmt}
}
