package parser

import (
	"github.com/st0012/rooby/ast"
	"github.com/st0012/rooby/lexer"
	"testing"
)

func TestMethodChainExpression(t *testing.T) {
	input := `
		Person.new(a, b).bar(c).add(d);
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)

	firstCall := stmt.Expression.(*ast.CallExpression)

	testIdentifier(t, firstCall.Method, "add")
	testIdentifier(t, firstCall.Arguments[0], "d")

	secondCall := firstCall.Receiver.(*ast.CallExpression)

	testIdentifier(t, secondCall.Method, "bar")
	testIdentifier(t, secondCall.Arguments[0], "c")

	thirdCall := secondCall.Receiver.(*ast.CallExpression)

	testIdentifier(t, thirdCall.Method, "new")
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
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
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
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expcted=%q, got=%q", tt.expected, actual)
		}
	}
}
