package code_generator

import (
	"fmt"
	"github.com/st0012/Rooby/ast"
	"regexp"
	"strings"
)

type LocalTable struct {
	store map[int]string
}

func (lt *LocalTable) Get(index int) (string, bool) {
	obj, ok := lt.store[index]
	return obj, ok
}

func (lt *LocalTable) Set(index int, val string) string {
	lt.store[index] = val
	return val
}

type Scope struct {
	Self    *ast.Statement
	Program *ast.Program
	Out     *Scope
}

var instructionSet []string

func GenerateByteCode(program *ast.Program) string {
	scope := &Scope{Program: program}
	compileProgram(program.Statements, scope)
	instructions := strings.Join(instructionSet, "\n")
	return strings.TrimSpace(removeEmptyLine(instructions))
}

func compileProgram(stmts []ast.Statement, scope *Scope) {
	var result []string

	result = append(result, "<ProgramStart>")
	for _, statement := range stmts {
		s := &Scope{Self: &statement, Out: scope}
		result = append(result, compileStatement(statement, s))
	}

	result = append(result, "leave")
	instructionSet = append(instructionSet, strings.Join(result, "\n"))
}

func compileStatement(statement ast.Statement, scope *Scope) string {
	switch stmt := statement.(type) {
	case *ast.ExpressionStatement:
		return compileExpression(stmt.Expression, scope)
	case *ast.DefStatement:
		compileDefStmt(stmt, scope)
		return ""
	}

	return ""
}

func compileDefStmt(stmt *ast.DefStatement, scope *Scope) {
	var result []string

	result = append(result, fmt.Sprintf("<Def:%s>", stmt.Name.Value))

	for _, s := range stmt.BlockStatement.Statements {
		result = append(result, compileStatement(s, scope))
	}
	result = append(result, "leave")
	instructionSet = append(instructionSet, strings.Join(result, "\n"))
}

func compileExpression(exp ast.Expression, scope *Scope) string {
	switch exp := exp.(type) {
	case *ast.IntegerLiteral:
		return fmt.Sprintf(`putobject %d`, exp.Value)
	case *ast.InfixExpression:
		return compileInfixExpression(exp, scope)
	}

	return ""
}

func compileInfixExpression(node *ast.InfixExpression, scope *Scope) string {
	var operation string
	left := compileExpression(node.Left, scope)
	right := compileExpression(node.Right, scope)
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
