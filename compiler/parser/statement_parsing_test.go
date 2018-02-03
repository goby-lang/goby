package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"github.com/goby-lang/goby/compiler/token"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return x;", "x"},
		{"return true;", true},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", returnStmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
		testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue)
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
	defStmt.HasNormalParam(t, "x")
	defStmt.HasNormalParam(t, "c")

	body, ok := defStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expect body should be an expression statement. got=%T", body)
	}

	if !testInfixExpression(t, body.Expression, "x", "+", "y") {
		return
	}
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

	stmt := program.Statements[0].(*ast.ModuleStatement)

	testConstant(t, stmt.Name, "Foo")

	defStmt := stmt.Body.Statements[0].(*ast.DefStatement)

	testIdentifier(t, defStmt.Name, "bar")
	testIdentifier(t, defStmt.Parameters[0], "x")
	testIdentifier(t, defStmt.Parameters[1], "y")

	body, ok := defStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expect body should be an expression statement. got=%T", body)
	}

	if !testInfixExpression(t, body.Expression, "x", "+", "y") {
		return
	}
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

	stmt := program.Statements[0].(*ast.ClassStatement)

	testConstant(t, stmt.Name, "Foo")
	testConstant(t, stmt.SuperClass, "Bar")

	defStmt := stmt.Body.Statements[0].(*ast.DefStatement)

	testIdentifier(t, defStmt.Name, "bar")
	testIdentifier(t, defStmt.Parameters[0], "x")
	testIdentifier(t, defStmt.Parameters[1], "y")

	body, ok := defStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expect body should be an expression statement. got=%T", body)
	}

	if !testInfixExpression(t, body.Expression, "x", "+", "y") {
		return
	}
}

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

	firstStmt := program.Statements[0].(*ast.DefStatement)

	testLiteralExpression(t, firstStmt.Parameters[0], "x")
	testLiteralExpression(t, firstStmt.Parameters[1], "y")

	firstExpressionStmt := firstStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)

	testInfixExpression(t, firstExpressionStmt.Expression, "x", "+", "y")

	secondStmt := program.Statements[1].(*ast.DefStatement)

	if secondStmt.Token.Type != token.Def {
		t.Fatalf("expect DefStatement's token to be 'DEF'. got=%T", secondStmt.Token.Type)
	}

	if len(secondStmt.Parameters) != 0 {
		t.Fatalf("expect second method definition not having any parameters")
	}

	secondExpressionStmt := secondStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)
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

	stmt := program.Statements[0].(*ast.DefStatement)
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
