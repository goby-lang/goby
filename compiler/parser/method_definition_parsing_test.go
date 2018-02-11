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
	firstStmt.ShouldHasName(t, "add")
	firstStmt.ShouldHasNormalParam(t, "x")
	firstStmt.ShouldHasNormalParam(t, "y")

	firstExpression := firstStmt.MethodBody().NthStmt(1).IsExpression(t)
	testInfixExpression(t, firstExpression, "x", "+", "y")

	secondStmt := program.NthStmt(2).IsDefStmt(t)
	secondStmt.ShouldHasName(t, "foo")
	secondStmt.ShouldHasNoParam(t)

	secondExpression := secondStmt.MethodBody().NthStmt(1).IsExpression(t)
	testIntegerLiteral(t, secondExpression, 123)
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

	testIntegerLiteral(t, firstYield.Arguments[0], 1)
	testIntegerLiteral(t, firstYield.Arguments[1], 2)
	testIdentifier(t, firstYield.Arguments[2], "bar")

	secondExp := stmt.MethodBody().NthStmt(2).IsExpression(t)
	secondExp.IsYieldExpression(t)
}

