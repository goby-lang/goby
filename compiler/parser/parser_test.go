package parser

import (
	"github.com/gooby-lang/gooby/compiler/ast"
	"github.com/gooby-lang/gooby/compiler/lexer"
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

	exp := program.FirstStmt().IsExpression(t)
	firstCall := exp.IsCallExpression(t)
	firstCall.ShouldHaveMethodName("add")
	argName := firstCall.NthArgument(1)
	argName.IsIdentifier(t).ShouldHaveName("d")

	secondCall := firstCall.TestableReceiver().IsCallExpression(t)
	secondCall.ShouldHaveMethodName("bar")
	argName = secondCall.NthArgument(1)
	argName.IsIdentifier(t).ShouldHaveName("c")

	thirdCall := secondCall.TestableReceiver().IsCallExpression(t)
	thirdCall.ShouldHaveMethodName("new")
	argName1 := thirdCall.NthArgument(1)
	argName1.IsIdentifier(t).ShouldHaveName("a")
	argName2 := thirdCall.NthArgument(2)
	argName2.IsIdentifier(t).ShouldHaveName("b")

	originalReceiver := thirdCall.TestableReceiver().IsConstant(t)
	originalReceiver.ShouldHaveName("Person")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"(-a * b)",
		},
		{
			"!-a",
			"!-a",
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
			"(-5 * -5)",
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
			"-(5 + 5)",
		},
		{
			"!(true == true)",
			"!(true == true)",
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

// Makes sure to prohibit calling a capitalized method on toplevel
func TestProhibitingCallingCapitalizedMethod(t *testing.T) {
	input := `
	Const()
	`

	l := lexer.New(input)
	p := New(l)
	_, err := p.ParseProgram()

	if err == nil {
		t.Fatal("Calling a capitalized method on toplevel should be prohibited")
	} else {
		if err.Message != "cannot call CONSTANT with (. Line: 1" {
			t.Fatal("Error should be: 'cannot call CONSTANT with (. Line: 1': ", err.Message)
		}
	}
}

// If parser doesn't crash then we covered panic successfully
func TestRecoverMechanism(t *testing.T) {
	input := `
	if )(
	`
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()
}
