package parser

import (
	"github.com/goby-lang/goby/compiler/ast"
	"github.com/goby-lang/goby/compiler/lexer"
	"testing"
)

func TestCallExpressionWithKeywordArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string] int
	}{
		{`
		add(x: 111)
		`, map[string]int {
			"x": 111,
		} },
		{`
		add(x: 111, y: 222)
		`, map[string]int {
			"x": 111,
			"y": 222,
		} },
		{`
		add(x: 111, y: 222, z: 333)
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

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		c := stmt.Expression.(*ast.CallExpression)

		h, ok := c.Arguments[0].(*ast.HashExpression)
		if !ok {
			t.Fatalf("At case %d c.Statments[0] is not ast.HashExpression. got=%T", i, c.Arguments[0])
		}

		for k, expected := range tt.expected {
			testIntegerLiteral(t, h.Data[k], expected)
		}
	}
}
