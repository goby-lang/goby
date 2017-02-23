package parser

import (
	"fmt"
	"github.com/st0012/rooby/ast"
	"testing"
)

func testAssignStatement(t *testing.T, s ast.Statement, name string, value interface{}) bool {
	as, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if as.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, as.Name)
		return false
	}

	switch v := value.(type) {
	case int64:
		testIntegerLiteral(t, as.Value, v)
	case string:
		testIdentifier(t, as.Value, v)
	case bool:
		testBoolLiteral(t, as.Value, v)
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
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
	bo, ok := exp.(*ast.Boolean)
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
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
