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
		{"return foo", ast.TestingIdentifier("foo")},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		returnStmt := program.FirstStmt().IsReturnStmt(t)
		returnStmt.ShouldHasValue(t, tt.expectedValue)
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

	stmt := program.FirstStmt().IsClassStmt(t, "Foo")
	defStmt := stmt.HasMethod(t, "bar")
	defStmt.ShouldHasNormalParam(t, "x")
	defStmt.ShouldHasNormalParam(t, "y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)

	testInfixExpression(t, methodBodyExp, "x", "+", "y")
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

	stmt := program.FirstStmt().IsModuleStmt(t, "Foo")
	defStmt := stmt.HasMethod(t, "bar")
	defStmt.ShouldHasNormalParam(t, "x")
	defStmt.ShouldHasNormalParam(t, "y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)

	testInfixExpression(t, methodBodyExp, "x", "+", "y")
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

	classStmt := program.FirstStmt().IsClassStmt(t, "Foo")
	classStmt.ShouldInherits(t, "Bar")

	defStmt := classStmt.HasMethod(t, "bar")
	defStmt.ShouldHasNormalParam(t, "x")
	defStmt.ShouldHasNormalParam(t, "y")

	methodBodyExp := defStmt.MethodBody().NthStmt(1).IsExpression(t)

	testInfixExpression(t, methodBodyExp, "x", "+", "y")
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

	whileStatement := program.Statements[0].(*ast.WhileStatement)

	infix := whileStatement.Condition.(*ast.InfixExpression)

	testIdentifier(t, infix.Left, "i")

	if infix.Operator != "<" {
		t.Fatalf("Expect condition's infix operator to be '<'. got=%s", infix.Operator)
	}

	callExp, ok := infix.Right.(*ast.CallExpression)

	if !ok {
		t.Fatalf("Expect infix's right to be a CallExpression. got=%T", infix.Right)
	}

	testMethodName(t, callExp, "length")

	if callExp.Block != nil {
		t.Fatalf("Condition expression shouldn't have block")
	}

	// Test block
	block := whileStatement.Body
	firstStmt := block.Statements[0].(*ast.ExpressionStatement)
	firstCall := firstStmt.Expression.(*ast.CallExpression)
	testMethodName(t, firstCall, "puts")
	testIdentifier(t, firstCall.Arguments[0], "i")

	secondStmt := block.Statements[1].(*ast.ExpressionStatement)
	secondCall := secondStmt.Expression.(*ast.AssignExpression)
	testIdentifier(t, secondCall.Variables[0], "i")
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
