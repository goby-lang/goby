package lexer

import (
	"github.com/goby-lang/goby/compiler/token"
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
	elsif 10 == 11
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

	nil

	module Foo
	end

	foo.module

	require "foo"

	Foo::Bar

	Person.class

	class Foo
	  def class
	    Foo
	  end
	  def self.bar
       	    10
   	  end
	end

	10 % 5

	""
	''

	next
	:apple
	{ test:"abc" }
	{ test: :abc }
	{ test:50 }
	{ test: abc }
	{ test:abc }

	(1..5)

	while i < 10 do
	  break
	end

	a += 1
	b -= 2
	c ||= true

	"\nstring\n"
	'\nstring\n'
	"\tstring\t"
	'\tstring\t'
	"\vstring\v"
	'\vstring\v'
	"\fstring\f"
	'\fstring\f'
	"\rstring\r"
	'\rstring\r'
	"\"string\""
	'\"string\"'
	"\'string\'"
	'\'string\''
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
		{token.Class, "class", 4},
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
		{token.ElsIf, "elsif", 26},
		{token.Int, "10", 26},
		{token.Eq, "==", 26},
		{token.Int, "11", 26},
		{token.Return, "return", 27},
		{token.True, "true", 27},
		{token.Semicolon, ";", 27},
		{token.Else, "else", 28},
		{token.Return, "return", 29},
		{token.False, "false", 29},
		{token.Semicolon, ";", 29},
		{token.End, "end", 30},

		{token.String, "string1", 32},
		{token.Semicolon, ";", 32},
		{token.String, "string2", 33},
		{token.Semicolon, ";", 33},

		{token.Int, "10", 35},
		{token.Eq, "==", 35},
		{token.Int, "10", 35},
		{token.Semicolon, ";", 35},

		{token.Int, "10", 37},
		{token.NotEq, "!=", 37},
		{token.Int, "9", 37},
		{token.Semicolon, ";", 37},

		{token.Comment, "# This is comment.", 39},
		{token.Comment, "# And I should be ignored.", 40},

		{token.LBracket, "[", 42},
		{token.Int, "1", 42},
		{token.Comma, ",", 42},
		{token.Int, "2", 42},
		{token.Comma, ",", 42},
		{token.Int, "3", 42},
		{token.Comma, ",", 42},
		{token.Int, "4", 42},
		{token.Comma, ",", 42},
		{token.Int, "5", 42},
		{token.RBracket, "]", 42},

		{token.LBracket, "[", 43},
		{token.String, "test", 43},
		{token.Comma, ",", 43},
		{token.String, "test", 43},
		{token.RBracket, "]", 43},

		{token.LBrace, "{", 45},
		{token.Ident, "test", 45},
		{token.Colon, ":", 45},
		{token.String, "123", 45},
		{token.RBrace, "}", 45},

		{token.Ident, "a", 47},
		{token.Incr, "++", 47},
		{token.Ident, "b", 48},
		{token.Decr, "--", 48},

		{token.While, "while", 50},
		{token.Ident, "i", 50},
		{token.LT, "<", 50},
		{token.Int, "10", 50},
		{token.Do, "do", 50},
		{token.Ident, "puts", 51},
		{token.LParen, "(", 51},
		{token.Ident, "i", 51},
		{token.RParen, ")", 51},
		{token.Ident, "i", 52},
		{token.Incr, "++", 52},
		{token.End, "end", 53},

		{token.Ident, "require_relative", 55},
		{token.String, "foo", 55},

		{token.Int, "10", 57},
		{token.LTE, "<=", 57},
		{token.Int, "10", 57},
		{token.Semicolon, ";", 57},

		{token.Int, "10", 58},
		{token.GTE, ">=", 58},
		{token.Int, "10", 58},
		{token.Semicolon, ";", 58},

		{token.Ident, "a", 59},
		{token.Assign, "=", 59},
		{token.Int, "1", 59},
		{token.COMP, "<=>", 59},
		{token.Int, "2", 59},
		{token.Semicolon, ";", 59},

		{token.Int, "8", 61},
		{token.Pow, "**", 61},
		{token.Int, "10", 61},
		{token.Semicolon, ";", 61},

		{token.True, "true", 63},
		{token.And, "&&", 63},
		{token.False, "false", 63},
		{token.Semicolon, ";", 63},

		{token.False, "false", 64},
		{token.Or, "||", 64},
		{token.True, "true", 64},
		{token.Semicolon, ";", 64},

		{token.Null, "nil", 66},

		{token.Module, "module", 68},
		{token.Constant, "Foo", 68},
		{token.End, "end", 69},

		{token.Ident, "foo", 71},
		{token.Dot, ".", 71},
		{token.Ident, "module", 71},

		{token.Ident, "require", 73},
		{token.String, "foo", 73},

		{token.Constant, "Foo", 75},
		{token.ResolutionOperator, "::", 75},
		{token.Constant, "Bar", 75},

		{token.Constant, "Person", 77},
		{token.Dot, ".", 77},
		{token.Ident, "class", 77},

		{token.Class, "class", 79},
		{token.Constant, "Foo", 79},
		{token.Def, "def", 80},
		{token.Ident, "class", 80},
		{token.Constant, "Foo", 81},
		{token.End, "end", 82},
		{token.Def, "def", 83},
		{token.Self, "self", 83},
		{token.Dot, ".", 83},
		{token.Ident, "bar", 83},
		{token.Int, "10", 84},
		{token.End, "end", 85},
		{token.End, "end", 86},

		{token.Int, "10", 88},
		{token.Modulo, "%", 88},
		{token.Int, "5", 88},

		{token.String, "", 90},
		{token.String, "", 91},

		{token.Next, "next", 93},
		{token.String, "apple", 94},

		{token.LBrace, "{", 95},
		{token.Ident, "test", 95},
		{token.Colon, ":", 95},
		{token.String, "abc", 95},
		{token.RBrace, "}", 95},

		{token.LBrace, "{", 96},
		{token.Ident, "test", 96},
		{token.Colon, ":", 96},
		{token.String, "abc", 96},
		{token.RBrace, "}", 96},

		{token.LBrace, "{", 97},
		{token.Ident, "test", 97},
		{token.Colon, ":", 97},
		{token.Int, "50", 97},
		{token.RBrace, "}", 97},

		{token.LBrace, "{", 98},
		{token.Ident, "test", 98},
		{token.Colon, ":", 98},
		{token.Ident, "abc", 98},
		{token.RBrace, "}", 98},

		{token.LBrace, "{", 99},
		{token.Ident, "test", 99},
		{token.Colon, ":", 99},
		{token.Ident, "abc", 99},
		{token.RBrace, "}", 99},

		{token.LParen, "(", 101},
		{token.Int, "1", 101},
		{token.Range, "..", 101},
		{token.Int, "5", 101},
		{token.RParen, ")", 101},

		{token.While, "while", 103},
		{token.Ident, "i", 103},
		{token.LT, "<", 103},
		{token.Int, "10", 103},
		{token.Do, "do", 103},
		{token.Break, "break", 104},
		{token.End, "end", 105},

		{token.Ident, "a", 107},
		{token.PlusEq, "+=", 107},
		{token.Int, "1", 107},

		{token.Ident, "b", 108},
		{token.MinusEq, "-=", 108},
		{token.Int, "2", 108},

		{token.Ident, "c", 109},
		{token.OrEq, "||=", 109},
		{token.True, "true", 109},

		// Escaped character tests
		{token.String, "\nstring\n", 111},
		{token.String, "\\nstring\\n", 112},
		{token.String, "\tstring\t", 113},
		{token.String, "\\tstring\\t", 114},
		{token.String, "\vstring\v", 115},
		{token.String, "\\vstring\\v", 116},
		{token.String, "\fstring\f", 117},
		{token.String, "\\fstring\\f", 118},
		{token.String, "\rstring\r", 119},
		{token.String, "\\rstring\\r", 120},
		{token.String, "\"string\"", 121},
		{token.String, "\\\"string\\\"", 122},
		{token.String, "'string'", 123},
		{token.String, "'string'", 124},

		{token.EOF, "", 125},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)

		}
		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line number wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}
	}
}
