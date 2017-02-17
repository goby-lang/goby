package parser

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/lexer"
	"github.com/st0012/rooby/token"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}

		if !testLetStatement(t, program.Statements[0], tt.expectedIdentifier) {
			return
		}
	}
}

func TestConstantAssignment(t *testing.T) {
	input := `
	let Foo = 5;
	Foo;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// First statement
	letStmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expect first statement to be LetStatement. got=%T", program.Statements[0])
	}

	variableName, ok := letStmt.Name.(*ast.Constant)
	if !ok {
		t.Fatalf("expect statement's name to be a constant. got=%T (%s)", letStmt.Name, letStmt.Name.String())
	}

	if variableName.Value != "Foo" {
		t.Fatalf("expect variable's name to be %s. got=%s", "Foo", variableName.Value)
	}

	// Second statement

	expStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expect second statement to be ExpressionStatement. got=%T", program.Statements[1])
	}

	variable, ok := expStmt.Expression.(*ast.Constant)
	if !ok {
		t.Fatalf("expect expression to be a constant. got=%T", expStmt.Expression)
	}

	if variable.Value != "Foo" {
		t.Fatalf("expect variable's name to be %s. got=%s", "Foo", variable.Value)
	}
}

func TestInstanceVariableAssignment(t *testing.T) {
	input := `
	let @foo = 5;
	@foo;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// First statement
	letStmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expect first statement to be LetStatement. got=%T", program.Statements[0])
	}

	variableName, ok := letStmt.Name.(*ast.InstanceVariable)
	if !ok {
		t.Fatalf("expect statement's name to be an instance variable. got=%T (%s)", letStmt.Name, letStmt.Name.String())
	}

	if variableName.Value != "@foo" {
		t.Fatalf("expect variable's name to be %s. got=%s", "@foo", variableName.Value)
	}

	// Second statement

	expStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expect second statement to be ExpressionStatement. got=%T", program.Statements[1])
	}

	variable, ok := expStmt.Expression.(*ast.InstanceVariable)
	if !ok {
		t.Fatalf("expect expression to be an instance variable. got=%T", expStmt.Expression)
	}

	if variable.Value != "@foo" {
		t.Fatalf("expect variable's name to be %s. got=%s", "@foo", variable.Value)
	}
}

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

		program := p.ParseProgram()
		checkParserErrors(t, p)

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
	class Foo {
		def bar(x, y) {
			x + y;
		}
	}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ClassStatement)

	if stmt.Token.Type != token.CLASS {
		t.Fatalf("expect token to be CLASS. got=%T", stmt.Token)
	}

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

func TestDefStatement(t *testing.T) {
	input := `
		def add(x, y) {
			x + y
		}
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.DefStatement)

	if stmt.Token.Type != token.DEF {
		t.Fatalf("expect DefStatement's token to be 'DEF'. got=%T", stmt.Token.Type)
	}

	testLiteralExpression(t, stmt.Parameters[0], "x")
	testLiteralExpression(t, stmt.Parameters[1], "y")

	expressionStmt := stmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)

	testInfixExpression(t, expressionStmt.Expression, "x", "+", "y")
}
