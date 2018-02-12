package parser

import (
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
	firstStmt.ShouldHasName("add")
	firstStmt.ShouldHasNormalParam("x")
	firstStmt.ShouldHasNormalParam("y")

	firstExpression := firstStmt.MethodBody().NthStmt(1).IsExpression(t)
	infixExp := firstExpression.IsInfixExpression(t)
	infixExp.ShouldHasOperator("+")
	infixExp.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	infixExp.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")

	secondStmt := program.NthStmt(2).IsDefStmt(t)
	secondStmt.ShouldHasName("foo")
	secondStmt.ShouldHasNoParam()

	secondExpression := secondStmt.MethodBody().NthStmt(1).IsExpression(t)
	secondExpression.IsIntegerLiteral(t).ShouldEqualTo(123)
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
	firstExp := stmt.MethodBody().NthStmt(1).IsExpression(t)
	firstYield := firstExp.IsYieldExpression(t)

	firstYield.NthArgument(1).IsIntegerLiteral(t).ShouldEqualTo(1)
	firstYield.NthArgument(2).IsIntegerLiteral(t).ShouldEqualTo(2)
	firstYield.NthArgument(3).IsIdentifier(t).ShouldHasName("bar")

	secondExp := stmt.MethodBody().NthStmt(2).IsExpression(t)
	secondExp.IsYieldExpression(t)
}
