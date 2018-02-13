package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
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
		returnStmt.ShouldHasValue(tt.expectedValue)
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
	stmt.ShouldHasName("Foo")
	defStmt := stmt.HasMethod("bar")
	defStmt.ShouldHasNormalParam("x")
	defStmt.ShouldHasNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHasOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
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
	stmt.ShouldHasName("Foo")
	defStmt := stmt.HasMethod(t, "bar")
	defStmt.ShouldHasNormalParam("x")
	defStmt.ShouldHasNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHasOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
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
	classStmt.ShouldHasName("Foo")
	classStmt.ShouldInherits("Bar")

	defStmt := classStmt.HasMethod("bar")
	defStmt.ShouldHasNormalParam("x")
	defStmt.ShouldHasNormalParam("y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)
	infix := methodBodyExp.IsInfixExpression(t)
	infix.ShouldHasOperator("+")
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHasName("x")
	infix.TestableRightExpression().IsIdentifier(t).ShouldHasName("y")
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
	infix.TestableLeftExpression().IsIdentifier(t).ShouldHasName("i")
	infix.ShouldHasOperator("<")
	callExp := infix.TestableRightExpression().IsCallExpression(t)
	callExp.ShouldHasMethodName("length")

	if callExp.Block != nil {
		t.Fatalf("Condition expression shouldn't have block")
	}

	// Test block
	block := whileStatement.CodeBlock()
	firstExp := block.NthStmt(1).IsExpression(t)
	firstCall := firstExp.IsCallExpression(t)
	firstCall.ShouldHasMethodName("puts")
	firstCall.NthArgument(1).IsIdentifier(t).ShouldHasName("i")

	secondExp := block.NthStmt(2).IsExpression(t)
	secondCall := secondExp.IsAssignExpression(t)
	secondCall.NthVariable(1).IsIdentifier(t).ShouldHasName("i")
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
