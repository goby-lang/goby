package parser

import (
	"github.com/gooby-lang/gooby/compiler/lexer"
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

	def bar(x = 10, y: ); end

	def baz(z: 100, *s); end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	firstStmt := program.FirstStmt().IsDefStmt(t)
	firstStmt.ShouldHaveName("add")
	firstStmt.ShouldHaveNormalParam("x")
	firstStmt.ShouldHaveNormalParam("y")

	firstExpression := firstStmt.MethodBody().NthStmt(1).IsExpression(t)
	infixExp := firstExpression.IsInfixExpression(t)
	infixExp.ShouldHaveOperator("+")
	infixExp.TestableLeftExpression().IsIdentifier(t).ShouldHaveName("x")
	infixExp.TestableRightExpression().IsIdentifier(t).ShouldHaveName("y")

	secondStmt := program.NthStmt(2).IsDefStmt(t)
	secondStmt.ShouldHaveName("foo")
	secondStmt.ShouldHaveNoParam()

	secondExpression := secondStmt.MethodBody().NthStmt(1).IsExpression(t)
	secondExpression.IsIntegerLiteral(t).ShouldEqualTo(123)

	thirdStmt := program.NthStmt(3).IsDefStmt(t)
	thirdStmt.ShouldHaveName("bar")
	thirdStmt.ShouldHaveOptionalParam("x")
	thirdStmt.ShouldHaveRequiredKeywordParam("y")

	fourthStmt := program.NthStmt(4).IsDefStmt(t)
	fourthStmt.ShouldHaveName("baz")
	fourthStmt.ShouldHaveOptionalKeywordParam("z")
	fourthStmt.ShouldHaveSplatParam("s")
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
	firstYield.NthArgument(3).IsIdentifier(t).ShouldHaveName("bar")

	secondExp := stmt.MethodBody().NthStmt(2).IsExpression(t)
	secondExp.IsYieldExpression(t)
}
