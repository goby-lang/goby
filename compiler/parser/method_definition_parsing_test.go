package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"testing"
)

func TestDefStatement(t *testing.T) {
	input := `
	def add(x, y)
	  x + y
	end

	def foo
	  123;
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	firstStmt := program.FirstStmt().IsDefStmt(t)
	firstStmt.ShouldHasName(t, "add")
	firstStmt.ShouldHasNormalParam(t, "x")
	firstStmt.ShouldHasNormalParam(t, "y")

	firstExpressionStmt := firstStmt.MethodBody().NthStmt(1).IsExpressionStmt(t)
	testInfixExpression(t, firstExpressionStmt.Expression, "x", "+", "y")

	secondStmt := program.NthStmt(2).IsDefStmt(t)
	secondStmt.ShouldHasName(t, "foo")
	secondStmt.ShouldHasNoParam(t)

	secondExpressionStmt := secondStmt.MethodBody().NthStmt(1).IsExpressionStmt(t)
	testIntegerLiteral(t, secondExpressionStmt.Expression, 123)
}

func TestDefStatementWithYield(t *testing.T) {
	input := `
	def foo
	  yield(1, 2, bar)
	  yield
	end
	`
	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	stmt := program.FirstStmt().IsDefStmt(t)
	block := stmt.BlockStatement
	firstStmt, ok := block.Statements[0].(*ast.ExpressionStatement)
	firstYield := firstStmt.Expression.(*ast.YieldExpression)

	if !ok {
		t.Fatalf("Expect method's body is an YieldExpression. got=%T", block.Statements[0])
	}

	testIntegerLiteral(t, firstYield.Arguments[0], 1)
	testIntegerLiteral(t, firstYield.Arguments[1], 2)
	testIdentifier(t, firstYield.Arguments[2], "bar")

	secondStmt, ok := block.Statements[1].(*ast.ExpressionStatement)
	_, ok = secondStmt.Expression.(*ast.YieldExpression)

	if !ok {
		t.Fatalf("Expect method's body is an YieldExpression. got=%T", block.Statements[1])
	}
}

