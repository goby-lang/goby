package parser_test

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/lexer"
	"github.com/st0012/rooby/token"
	"testing"
	"github.com/st0012/rooby/parser"
)

func TestAssignStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
		{"@foobar = y;", "@foobar", "y"},
		{"Foo = '123';", "Foo", 10},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}

		testAssignStatement(t, program.Statements[0], tt.expectedIdentifier, tt.expectedValue)
	}
}

func TestConstantAssignment(t *testing.T) {
	input := `
	Foo = 5;
	Foo;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// First statement
	testAssignStatement(t, program.Statements[0], "Foo", 5)

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
	@foo = 5;
	@foo;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	// First statement
	testAssignStatement(t, program.Statements[0], "@foo", 5)

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
		p := parser.New(l)

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
	p := parser.New(l)
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

func TestClassStatementWithInheritance(t *testing.T) {
	input := `
	class Foo < Bar {
		def bar(x, y) {
			x + y;
		}
	}
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ClassStatement)

	if stmt.Token.Type != token.CLASS {
		t.Fatalf("expect token to be CLASS. got=%T", stmt.Token)
	}

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
		def add(x, y) {
			x + y
		}

		def foo {
			123;
		}
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	firstStmt := program.Statements[0].(*ast.DefStatement)

	if firstStmt.Token.Type != token.DEF {
		t.Fatalf("expect DefStatement's token to be 'DEF'. got=%T", firstStmt.Token.Type)
	}

	testLiteralExpression(t, firstStmt.Parameters[0], "x")
	testLiteralExpression(t, firstStmt.Parameters[1], "y")

	firstExpressionStmt := firstStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)

	testInfixExpression(t, firstExpressionStmt.Expression, "x", "+", "y")

	secondStmt := program.Statements[1].(*ast.DefStatement)

	if secondStmt.Token.Type != token.DEF {
		t.Fatalf("expect DefStatement's token to be 'DEF'. got=%T", secondStmt.Token.Type)
	}

	if len(secondStmt.Parameters) != 0 {
		t.Fatalf("expect second method definition not having any parameters")
	}

	secondExpressionStmt := secondStmt.BlockStatement.Statements[0].(*ast.ExpressionStatement)
	testIntegerLiteral(t, secondExpressionStmt.Expression, 123)
}
