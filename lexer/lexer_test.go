package lexer

import (
	"github.com/rooby-lang/rooby/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	five = 5;
	ten = 10;

	class Person
	  def initialize(a)
	    @a = a;
	  end

	  def add(x, y)
	    x + y;
	  end

	  def ten()
	    self.add(1, 9);
	  end
	end

	p = Person.new;
	result = p.add(five, ten);

	!-/*5;
	5 < 10 >5;

	if 5 < 10
	  return true;
	else
	  return false;
	end

	"string1";
	'string2';

	10 == 10;

	10 != 9;

	# This is comment.
	# And I should be ignored.

	[1, 2, 3, 4, 5]
	["test", "test"]

	{ test: "123" }

	a++
	b--

	while i < 10 do
	  puts(i)
	  i++
	end

	require_relative "foo"

	10 <= 10;
	10 >= 10;
	a = 1 <=> 2;

	8 ** 10;

	true && false;
	false || true;
	`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLine    int
	}{
		{token.Ident, "five", 1},
		{token.Assign, "=", 1},
		{token.Int, "5", 1},
		{token.Semicolon, ";", 1},
		{token.Ident, "ten", 2},
		{token.Assign, "=", 2},
		{token.Int, "10", 2},
		{token.Semicolon, ";", 2},

		// class is default to be ident
		{token.Ident, "class", 4},
		{token.Constant, "Person", 4},

		{token.Def, "def", 5},
		{token.Ident, "initialize", 5},
		{token.LParen, "(", 5},
		{token.Ident, "a", 5},
		{token.RParen, ")", 5},
		{token.InstanceVariable, "@a", 6},
		{token.Assign, "=", 6},
		{token.Ident, "a", 6},
		{token.Semicolon, ";", 6},
		{token.End, "end", 7},

		{token.Def, "def", 9},
		{token.Ident, "add", 9},
		{token.LParen, "(", 9},
		{token.Ident, "x", 9},
		{token.Comma, ",", 9},
		{token.Ident, "y", 9},
		{token.RParen, ")", 9},
		{token.Ident, "x", 10},
		{token.Plus, "+", 10},
		{token.Ident, "y", 10},
		{token.Semicolon, ";", 10},
		{token.End, "end", 11},

		{token.Def, "def", 13},
		{token.Ident, "ten", 13},
		{token.LParen, "(", 13},
		{token.RParen, ")", 13},
		{token.Self, "self", 14},
		{token.Dot, ".", 14},
		{token.Ident, "add", 14},
		{token.LParen, "(", 14},
		{token.Int, "1", 14},
		{token.Comma, ",", 14},
		{token.Int, "9", 14},
		{token.RParen, ")", 14},
		{token.Semicolon, ";", 14},
		{token.End, "end", 15},
		{token.End, "end", 16},

		{token.Ident, "p", 18},
		{token.Assign, "=", 18},
		{token.Constant, "Person", 18},
		{token.Dot, ".", 18},
		{token.Ident, "new", 18},
		{token.Semicolon, ";", 18},
		{token.Ident, "result", 19},
		{token.Assign, "=", 19},
		{token.Ident, "p", 19},
		{token.Dot, ".", 19},
		{token.Ident, "add", 19},
		{token.LParen, "(", 19},
		{token.Ident, "five", 19},
		{token.Comma, ",", 19},
		{token.Ident, "ten", 19},
		{token.RParen, ")", 19},
		{token.Semicolon, ";", 19},

		{token.Bang, "!", 21},
		{token.Minus, "-", 21},
		{token.Slash, "/", 21},
		{token.Asterisk, "*", 21},
		{token.Int, "5", 21},
		{token.Semicolon, ";", 21},

		{token.Int, "5", 22},
		{token.LT, "<", 22},
		{token.Int, "10", 22},
		{token.GT, ">", 22},
		{token.Int, "5", 22},
		{token.Semicolon, ";", 22},

		{token.If, "if", 24},
		{token.Int, "5", 24},
		{token.LT, "<", 24},
		{token.Int, "10", 24},
		{token.Return, "return", 25},
		{token.True, "true", 25},
		{token.Semicolon, ";", 25},
		{token.Else, "else", 26},
		{token.Return, "return", 27},
		{token.False, "false", 27},
		{token.Semicolon, ";", 27},
		{token.End, "end", 28},

		{token.String, "string1", 30},
		{token.Semicolon, ";", 30},
		{token.String, "string2", 31},
		{token.Semicolon, ";", 31},

		{token.Int, "10", 33},
		{token.Eq, "==", 33},
		{token.Int, "10", 33},
		{token.Semicolon, ";", 33},

		{token.Int, "10", 35},
		{token.NotEq, "!=", 35},
		{token.Int, "9", 35},
		{token.Semicolon, ";", 35},

		{token.Comment, "# This is comment.", 37},
		{token.Comment, "# And I should be ignored.", 38},

		{token.LBracket, "[", 40},
		{token.Int, "1", 40},
		{token.Comma, ",", 40},
		{token.Int, "2", 40},
		{token.Comma, ",", 40},
		{token.Int, "3", 40},
		{token.Comma, ",", 40},
		{token.Int, "4", 40},
		{token.Comma, ",", 40},
		{token.Int, "5", 40},
		{token.RBracket, "]", 40},

		{token.LBracket, "[", 41},
		{token.String, "test", 41},
		{token.Comma, ",", 41},
		{token.String, "test", 41},
		{token.RBracket, "]", 41},

		{token.LBrace, "{", 43},
		{token.Ident, "test", 43},
		{token.Colon, ":", 43},
		{token.String, "123", 43},
		{token.RBrace, "}", 43},

		{token.Ident, "a", 45},
		{token.Incr, "++", 45},
		{token.Ident, "b", 46},
		{token.Decr, "--", 46},

		{token.While, "while", 48},
		{token.Ident, "i", 48},
		{token.LT, "<", 48},
		{token.Int, "10", 48},
		{token.Do, "do", 48},
		{token.Ident, "puts", 49},
		{token.LParen, "(", 49},
		{token.Ident, "i", 49},
		{token.RParen, ")", 49},
		{token.Ident, "i", 50},
		{token.Incr, "++", 50},
		{token.End, "end", 51},

		{token.RequireRelative, "require_relative", 53},
		{token.String, "foo", 53},

		{token.Int, "10", 55},
		{token.LTE, "<=", 55},
		{token.Int, "10", 55},
		{token.Semicolon, ";", 55},

		{token.Int, "10", 56},
		{token.GTE, ">=", 56},
		{token.Int, "10", 56},
		{token.Semicolon, ";", 56},

		{token.Ident, "a", 57},
		{token.Assign, "=", 57},
		{token.Int, "1", 57},
		{token.COMP, "<=>", 57},
		{token.Int, "2", 57},
		{token.Semicolon, ";", 57},

		{token.Int, "8", 59},
		{token.Pow, "**", 59},
		{token.Int, "10", 59},
		{token.Semicolon, ";", 59},

		{token.True, "true", 61},
		{token.And, "&&", 61},
		{token.False, "false", 61},
		{token.Semicolon, ";", 61},

		{token.False, "false", 62},
		{token.Or, "||", 62},
		{token.True, "true", 62},
		{token.Semicolon, ";", 62},

		{token.EOF, "", 63},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. exprected=%q, got=%q", i, tt.expectedType, tok.Type)

		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. exprected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)

		}
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line number wrong. exprected=%d, got=%d", i, tt.expectedLine, tok.Line)

		}
	}
}
