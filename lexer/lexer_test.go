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
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
	}{
		{token.IDENT, "five", 1},
		{token.ASSIGN, "=", 1},
		{token.INT, "5", 1},
		{token.SEMICOLON, ";", 1},
		{token.IDENT, "ten", 2},
		{token.ASSIGN, "=", 2},
		{token.INT, "10", 2},
		{token.SEMICOLON, ";", 2},

		// class is default to be ident
		{token.IDENT, "class", 4},
		{token.CONSTANT, "Person", 4},

		{token.DEF, "def", 5},
		{token.IDENT, "initialize", 5},
		{token.LPAREN, "(", 5},
		{token.IDENT, "a", 5},
		{token.RPAREN, ")", 5},
		{token.INSTANCE_VARIABLE, "@a", 6},
		{token.ASSIGN, "=", 6},
		{token.IDENT, "a", 6},
		{token.SEMICOLON, ";", 6},
		{token.END, "end", 7},

		{token.DEF, "def", 9},
		{token.IDENT, "add", 9},
		{token.LPAREN, "(", 9},
		{token.IDENT, "x", 9},
		{token.COMMA, ",", 9},
		{token.IDENT, "y", 9},
		{token.RPAREN, ")", 9},
		{token.IDENT, "x", 10},
		{token.PLUS, "+", 10},
		{token.IDENT, "y", 10},
		{token.SEMICOLON, ";", 10},
		{token.END, "end", 11},

		{token.DEF, "def", 13},
		{token.IDENT, "ten", 13},
		{token.LPAREN, "(", 13},
		{token.RPAREN, ")", 13},
		{token.SELF, "self", 14},
		{token.DOT, ".", 14},
		{token.IDENT, "add", 14},
		{token.LPAREN, "(", 14},
		{token.INT, "1", 14},
		{token.COMMA, ",", 14},
		{token.INT, "9", 14},
		{token.RPAREN, ")", 14},
		{token.SEMICOLON, ";", 14},
		{token.END, "end", 15},
		{token.END, "end", 16},

		{token.IDENT, "p", 18},
		{token.ASSIGN, "=", 18},
		{token.CONSTANT, "Person", 18},
		{token.DOT, ".", 18},
		{token.IDENT, "new", 18},
		{token.SEMICOLON, ";", 18},
		{token.IDENT, "result", 19},
		{token.ASSIGN, "=", 19},
		{token.IDENT, "p", 19},
		{token.DOT, ".", 19},
		{token.IDENT, "add", 19},
		{token.LPAREN, "(", 19},
		{token.IDENT, "five", 19},
		{token.COMMA, ",", 19},
		{token.IDENT, "ten", 19},
		{token.RPAREN, ")", 19},
		{token.SEMICOLON, ";", 19},

		{token.BANG, "!", 21},
		{token.MINUS, "-", 21},
		{token.SLASH, "/", 21},
		{token.ASTERISK, "*", 21},
		{token.INT, "5", 21},
		{token.SEMICOLON, ";", 21},

		{token.INT, "5", 22},
		{token.LT, "<", 22},
		{token.INT, "10", 22},
		{token.GT, ">", 22},
		{token.INT, "5", 22},
		{token.SEMICOLON, ";", 22},

		{token.IF, "if", 24},
		{token.INT, "5", 24},
		{token.LT, "<", 24},
		{token.INT, "10", 24},
		{token.RETURN, "return", 25},
		{token.TRUE, "true", 25},
		{token.SEMICOLON, ";", 25},
		{token.ELSE, "else", 26},
		{token.RETURN, "return", 27},
		{token.FALSE, "false", 27},
		{token.SEMICOLON, ";", 27},
		{token.END, "end", 28},

		{token.STRING, "string1", 30},
		{token.SEMICOLON, ";", 30},
		{token.STRING, "string2", 31},
		{token.SEMICOLON, ";", 31},

		{token.INT, "10", 33},
		{token.EQ, "==", 33},
		{token.INT, "10", 33},
		{token.SEMICOLON, ";", 33},

		{token.INT, "10", 35},
		{token.NOT_EQ, "!=", 35},
		{token.INT, "9", 35},
		{token.SEMICOLON, ";", 35},

		{token.COMMENT, "# This is comment.", 37},
		{token.COMMENT, "# And I should be ignored.", 38},

		{token.LBRACKET, "[", 40},
		{token.INT, "1", 40},
		{token.COMMA, ",", 40},
		{token.INT, "2", 40},
		{token.COMMA, ",", 40},
		{token.INT, "3", 40},
		{token.COMMA, ",", 40},
		{token.INT, "4", 40},
		{token.COMMA, ",", 40},
		{token.INT, "5", 40},
		{token.RBRACKET, "]", 40},

		{token.LBRACKET, "[", 41},
		{token.STRING, "test", 41},
		{token.COMMA, ",", 41},
		{token.STRING, "test", 41},
		{token.RBRACKET, "]", 41},

		{token.LBRACE, "{", 43},
		{token.IDENT, "test", 43},
		{token.COLON, ":", 43},
		{token.STRING, "123", 43},
		{token.RBRACE, "}", 43},

		{token.IDENT, "a", 45},
		{token.INCR, "++", 45},
		{token.IDENT, "b", 46},
		{token.DECR, "--", 46},

		{token.WHILE, "while", 48},
		{token.IDENT, "i", 48},
		{token.LT, "<", 48},
		{token.INT, "10", 48},
		{token.DO, "do", 48},
		{token.IDENT, "puts", 49},
		{token.LPAREN, "(", 49},
		{token.IDENT, "i", 49},
		{token.RPAREN, ")", 49},
		{token.IDENT, "i", 50},
		{token.INCR, "++", 50},
		{token.END, "end", 51},

		{token.EOF, "", 52},
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
