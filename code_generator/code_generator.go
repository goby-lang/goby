package code_generator

import (
	"fmt"
	"github.com/st0012/Rooby/ast"
	"regexp"
	"strings"
)

type LocalTable struct {
	store []string
	count int
}

func (lt *LocalTable) Get(index int) string {
	return lt.store[index]
}

func (lt *LocalTable) Set(val string) int {
	if len(lt.store) >= lt.count {
		lt.store[lt.count] = val
	} else {
		lt.store = append(lt.store, val)
	}

	lt.count += 1
	return lt.count
}

type Scope struct {
	Self       ast.Statement
	Program    *ast.Program
	Out        *Scope
	LocalTable *LocalTable
}

var instructionSet []string

func newLocalTable() *LocalTable {
	s := make([]string, 1)
	return &LocalTable{store: s}
}

func newScope(scope *Scope, stmt ast.Statement) *Scope {
	return &Scope{Out: scope, LocalTable: newLocalTable(), Self: stmt}
}

func GenerateByteCode(program *ast.Program) string {
	scope := &Scope{Program: program, LocalTable: newLocalTable()}
	compileProgram(program.Statements, scope)
	instructions := strings.Join(instructionSet, "\n")
	return strings.TrimSpace(removeEmptyLine(instructions))
}

func compileProgram(stmts []ast.Statement, scope *Scope) {
	var result []string

	result = append(result, "<ProgramStart>")
	for _, statement := range stmts {
		s := newScope(scope, statement)
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
		s := newScope(scope, stmt)
		compileDefStmt(stmt, s)
		return ""
	case *ast.AssignStatement:
		s := newScope(scope, stmt)
		return compileAssignStmt(stmt, s)
	}

	return ""
}

func compileAssignStmt(stmt *ast.AssignStatement, scope *Scope) string {
	n := stmt.Name.ReturnValue()
	index := scope.LocalTable.Set(n)
	result := fmt.Sprintf(`
%s
setlocal %d
`, compileExpression(stmt.Value, scope), index)
	return result
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
