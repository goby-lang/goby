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

	stmt := program.Statements[0].(*ast.ClassStatement)

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

func TestSimpleKeywordArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string] int
	}{
		{`
		def add(x: 111)
		end
		`, map[string]int {
			"x": 111,
		} },
		{`
		def add(x: 111, y: 222)
		end
		`, map[string]int {
			"x": 111,
			"y": 222,
		} },
		{`
		def add(x: 111, y: 222, z: 333)
		end
		`, map[string]int {
			"x": 111,
			"y": 222,
			"z": 333,
		} },
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatalf("At case %d " + err.Message, i)
		}

		firstStmt := program.Statements[0].(*ast.DefStatement)
		h, ok := firstStmt.Parameters[0].(*ast.HashExpression)

		if !ok {
			t.Fatalf("At case %d program.Statments[0] is not ast.HashExpression. got=%T", i, program.Statements[0])
		}

		for k, expected := range tt.expected {
			testIntegerLiteral(t, h.Data[k], expected)
		}
	}
}

func TestKeywordArgumentWithNoDefaultValue(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string] interface{}
	}{
		{`
		def add(x:)
		end
		`, map[string]interface{} {
			"x": nil,
		} },
		{`
		def add(x: 111, y:)
		end
		`, map[string]interface{} {
			"x": 111,
			"y": nil,
		} },
		{`
		def add(x: 111, y:, z:)
		end
		`, map[string]interface{} {
			"x": 111,
			"y": nil,
			"z": nil,
		} },
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatalf("At case %d " + err.Message, i)
		}

		firstStmt := program.Statements[0].(*ast.DefStatement)
		h, ok := firstStmt.Parameters[0].(*ast.HashExpression)

		if !ok {
			t.Fatalf("At case %d program.Statments[0] is not ast.HashExpression. got=%T", i, program.Statements[0])
		}

		for k, expected := range tt.expected {
			if expected == nil {
				if h.Data[k] != nil {
					t.Fatalf("At case %d h.Data[%T] is not nil. got=%T", i, k, h.Data[k])
				}
			} else {
				testIntegerLiteral(t, h.Data[k], expected.(int))
			}
		}
	}
}

func TestDefStatementArgumentDefinitionError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		def add(x, y, x)
			x + y
		end
		`, "Duplicate argument name: \"x\". Line: 1"},
		{`
		def add(a, b, b = [1, 2])
			a + b
		end
		`, "Duplicate argument name: \"b\". Line: 1"},
		{`
		def add(a, b = [1, 2], b = [3, 4])
			a + b
		end
		`, "Duplicate argument name: \"b\". Line: 1"},
		{`
		def add(a, b = 1, b = 2)
			a + b
		end
		`, "Duplicate argument name: \"b\". Line: 1"},
		{`
		def add(a = 1, b)
			a + b
		end
		`, "Normal argument \"b\" should be defined before optioned argument. Line: 1"},
		{`
		def add(a = 1, b, c)
			a + b
		end
		`, "Normal argument \"b\" should be defined before optioned argument. Line: 1"},
		{`
		def add(*a, b)
			a + b
		end
		`, "Normal argument \"b\" should be defined before splat argument. Line: 1"},
		{`
		def add(*a, b = 1, c)
			a + b
		end
		`, "Optioned argument \"b = 1\" should be defined before splat argument. Line: 1"},
		{`
		def add(*a, *b, c)
			a + b
		end
		`, "Can't define splat argument more than once. Line: 1"},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.ParseProgram()

		if err == nil {
			t.Fatalf("At case %d expect not to allow duplicate argument name", i)
		}

		if err.Message != tt.expected {
			t.Fatalf("At case %d expect error message to be:\n  %s. got: \n%s", i, tt.expected, err.Message)
		}
	}
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
	  i++
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
	secondCall := secondStmt.Expression.(*ast.CallExpression)
	testIdentifier(t, secondCall.Receiver, "i")
	testMethodName(t, secondCall, "++")
}
