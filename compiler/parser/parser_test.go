package parser

import (
	"fmt"
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"testing"
)

func TestMethodChainExpression(t *testing.T) {
	input := `
		Person.new(a, b).bar(c).add(d);
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	firstCall := stmt.Expression.(*ast.CallExpression)

	testMethodName(t, firstCall, "add")
	testIdentifier(t, firstCall.Arguments[0], "d")

	secondCall := firstCall.Receiver.(*ast.CallExpression)

	testMethodName(t, secondCall, "bar")
	testIdentifier(t, secondCall.Arguments[0], "c")

	thirdCall := secondCall.Receiver.(*ast.CallExpression)

	testMethodName(t, thirdCall, "new")
	testIdentifier(t, thirdCall.Arguments[0], "a")
	testIdentifier(t, thirdCall.Arguments[1], "b")

	originalReceiver := thirdCall.Receiver.(*ast.Constant)
	testConstant(t, originalReceiver, "Person")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"-5 * -5",
			"((-5) * (-5))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + n.add(b * c) + d",
			"((a + n.add((b * c))) + d)",
		},
		{
			"n.add(a, b, 1, 2 * 3, 4 + 5, m.add(6, 7 * 8))",
			"n.add(a, b, 1, (2 * 3), (4 + 5), m.add(6, (7 * 8)))",
		},
		{
			"n.add(a + b + c * d / f + g)",
			"n.add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatal(err.Message)
		}

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expcted=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIgnoreComments(t *testing.T) {
	input := `
		# This is comment.
		# Ignore me!
		p.add(1, 2 * 3, 4 + 5);
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatal(err.Message)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expect parser to ignore comment")
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("expect parser to ignore comment and return only call expression")
	}

}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expect exp to be IntegerLiteral. got=%T", exp)
	}
	if il.Value != value {
		t.Errorf("il.Value is not %d. got=%d", value, il.Value)
		return false
	}
	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral not %d. got=%s", value, il.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	sl, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expect exp to be StringLiteral. got=%T", exp)
	}
	if sl.Value != value {
		t.Errorf("il.Value is not %s. got=%s", value, sl.Value)
		return false
	}
	if sl.TokenLiteral() != value {
		t.Errorf("il.TokenLiteral not %s. got=%s", value, sl.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testConstant(t *testing.T, exp ast.Expression, value string) bool {
	constant, ok := exp.(*ast.Constant)
	if !ok {
		t.Errorf("exp not *ast.Constant. got=%T", exp)
		return false
	}
	if constant.Value != value {
		t.Errorf("constant.Value not %s. got=%s", value, constant.Value)
		return false
	}

	if constant.TokenLiteral() != value {
		t.Errorf("constant.TokenLiteral not %s. got=%s", value, constant.TokenLiteral())
		return false
	}

	return true
}

func testInstanceVariable(t *testing.T, exp ast.Expression, value string) bool {
	instVar, ok := exp.(*ast.InstanceVariable)
	if !ok {
		t.Errorf("exp not *ast.InstanceVariable. got=%T", exp)
		return false
	}
	if instVar.Value != value {
		t.Errorf("instVar.Value not %s. got=%s", value, instVar.Value)
		return false
	}

	if instVar.TokenLiteral() != value {
		t.Errorf("instVar.TokenLiteral not %s. got=%s", value, instVar.TokenLiteral())
		return false
	}

	return true
}

func testMethodName(t *testing.T, exp ast.Expression, value string) {
	callExp, ok := exp.(*ast.CallExpression)

	if !ok {
		t.Errorf("expect exp to be a CallExpression. got=%T", exp)
	}

	if callExp.Method != value {
		t.Errorf("expect method name to be %s. got=%s", value, callExp.Method)
	}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not %T. got=%T", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("opExp's operator is not %s. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func testBoolLiteral(t *testing.T, exp ast.Expression, v bool) bool {
	bo, ok := exp.(*ast.BooleanExpression)
	if !ok {
		t.Errorf("exp is not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != v {
		t.Errorf("bo.Value is not %t. got=%t", v, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", v) {
		t.Errorf("bo.TokenLiteral is not %t. got=%t", v, exp.TokenLiteral())
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expcted interface{},
) bool {
	switch v := expcted.(type) {
	case int:
		return testIntegerLiteral(t, exp, v)
	case int64:
		return testIntegerLiteral(t, exp, int(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
