package parser

import (
	"github.com/gooby-lang/gooby/compiler/ast"
	"github.com/gooby-lang/gooby/compiler/lexer"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5", 5},
		{"return 'x'", "x"},
		{"return true", true},
		{"return foo", ast.TestableIdentifierValue("foo")},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		returnStmt := program.FirstStmt().IsReturnStmt(t)
		returnStmt.ShouldHaveValue(tt.expectedValue)
	}
}

func TestClassStatement(t *testing.T) {
	input := `
	class Foo
	  def bar(x, y)
	    x + y
	  end
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	stmt := program.FirstStmt().IsClassStmt(t)
	stmt.ShouldHaveName("Foo")
	defStmt := stmt.HasMethod("bar")
	defStmt.ShouldHaveNormalParam("x")
	defStmt.ShouldHaveNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHaveOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHaveName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHaveName("y")
}

func TestModuleStatement(t *testing.T) {
	input := `
	module Foo
	  def bar(x, y)
	    x + y
	  end
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	stmt := program.FirstStmt().IsModuleStmt(t)
	stmt.ShouldHaveName("Foo")
	defStmt := stmt.HasMethod(t, "bar")
	defStmt.ShouldHaveNormalParam("x")
	defStmt.ShouldHaveNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHaveOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHaveName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHaveName("y")
}

func TestClassStatementWithInheritance(t *testing.T) {
	input := `
	class Foo < Bar
	  def bar(x, y)
	    x + y;
	  end
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	classStmt := program.FirstStmt().IsClassStmt(t)
	classStmt.ShouldHaveName("Foo")
	classStmt.ShouldInherit("Bar")

	defStmt := classStmt.HasMethod("bar")
	defStmt.ShouldHaveNormalParam("x")
	defStmt.ShouldHaveNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHaveOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHaveName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHaveName("y")
}

func TestWhileStatement(t *testing.T) {
	input := `
	while i < a.length do
	  puts(i)
	  i += 1
	end
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	whileStatement := program.FirstStmt().IsWhileStmt(t)

	infix := whileStatement.ConditionExpression().IsInfixExpression(t)
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHaveName("i")
	infix.ShouldHaveOperator("<")
	callExp := infix.TestableRightExpression().IsCallExpression(t)
	callExp.ShouldHaveMethodName("length")

	if callExp.Block != nil {
		t.Fatalf("Condition expression shouldn't have block")
	}

	// Test block
	block := whileStatement.CodeBlock()
	firstExp := block.NthStmt(1).IsExpression(t)
	firstCall := firstExp.IsCallExpression(t)
	firstCall.ShouldHaveMethodName("puts")
	firstCall.NthArgument(1).IsIdentifier(t).ShouldHaveName("i")

	secondExp := block.NthStmt(2).IsExpression(t)
	secondCall := secondExp.IsAssignExpression(t)
	secondCall.NthVariable(1).IsIdentifier(t).ShouldHaveName("i")
}

func TestWhileStatementWithoutDoKeywordFail(t *testing.T) {
	input := `
	while i < a.length
	  puts(i)
	  i += 1
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "expected next token to be DO, got IDENT(puts) instead. Line: 2" {
		t.Fatal("Condition expression should be followed by a do keyword")
	}

}

func TestInvalidMethodNameFail(t *testing.T) {
	input := `
	def ()
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "Invalid method name: (. Line: 1" {
		t.Fatal(err.Message)
	}
}

func TestInvalidParameter(t *testing.T) {
	input := `
	def foo(@a)
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "Invalid parameters: @a. Line: 1" {
		t.Fatal(err.Message)
	}
}

func TestInvalidMultipleParameter(t *testing.T) {
	input := `
	def foo(a, @b, c)
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "Invalid parameters: @b. Line: 1" {
		t.Fatal(err.Message)
	}
}

func TestParenthesisInNextLineIsNotParameter(t *testing.T) {
	input := `
	def test
	  (@instance = 1)
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

}

func TestInvalidIfStatement(t *testing.T) {
	input := `
	if
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "syntax error, unexpected end Line: 2" {
		t.Fatal(err.Message)
	}
}

func TestInvalidIfStatementWithParentheses(t *testing.T) {
	input := `
	if ()
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "expected next token to be ), got END(end) instead. Line: 2" {
		t.Fatal(err.Message)
	}
}

func TestInvalidIfStatementWithReverseParentheses(t *testing.T) {
	input := `
	if )(
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "expected next token to be ), got EOF() instead. Line: 2" {
		t.Fatal(err.Message)
	}
}

func TestInvalidIfStatementWithRightParentheses(t *testing.T) {
	input := `
	if )
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "unexpected ) Line: 1" {
		t.Fatal(err.Message)
	}
}

func TestInvalidIfStatementWithLeftParentheses(t *testing.T) {
	input := `
	if (
	end`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err.Message != "syntax error, unexpected end Line: 2" {
		t.Fatal(err.Message)
	}
}
