package lexer

import (
	"github.com/goby-lang/goby/compiler/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
	five = 5;
	ten = 10;
	_ = 15

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

	require_relative "foo"

	'a' =~ 'a';
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
     1_23
     12_3.45_6
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
		{token.Ident, "_", 3},
		{token.Assign, "=", 3},
		{token.Int, "15", 3},

		// class is default to be ident
		{token.Class, "class", 5},
		{token.Constant, "Person", 5},

		{token.Def, "def", 6},
		{token.Ident, "initialize", 6},
		{token.LParen, "(", 6},
		{token.Ident, "a", 6},
		{token.RParen, ")", 6},
		{token.InstanceVariable, "@a", 7},
		{token.Assign, "=", 7},
		{token.Ident, "a", 7},
		{token.Semicolon, ";", 7},
		{token.End, "end", 8},

		{token.Def, "def", 10},
		{token.Ident, "add", 10},
		{token.LParen, "(", 10},
		{token.Ident, "x", 10},
		{token.Comma, ",", 10},
		{token.Ident, "y", 10},
		{token.RParen, ")", 10},
		{token.Ident, "x", 11},
		{token.Plus, "+", 11},
		{token.Ident, "y", 11},
		{token.Semicolon, ";", 11},
		{token.End, "end", 12},

		{token.Def, "def", 14},
		{token.Ident, "ten", 14},
		{token.LParen, "(", 14},
		{token.RParen, ")", 14},
		{token.Self, "self", 15},
		{token.Dot, ".", 15},
		{token.Ident, "add", 15},
		{token.LParen, "(", 15},
		{token.Int, "1", 15},
		{token.Comma, ",", 15},
		{token.Int, "9", 15},
		{token.RParen, ")", 15},
		{token.Semicolon, ";", 15},
		{token.End, "end", 16},
		{token.End, "end", 17},

		{token.Ident, "p", 19},
		{token.Assign, "=", 19},
		{token.Constant, "Person", 19},
		{token.Dot, ".", 19},
		{token.Ident, "new", 19},
		{token.Semicolon, ";", 19},
		{token.Ident, "result", 20},
		{token.Assign, "=", 20},
		{token.Ident, "p", 20},
		{token.Dot, ".", 20},
		{token.Ident, "add", 20},
		{token.LParen, "(", 20},
		{token.Ident, "five", 20},
		{token.Comma, ",", 20},
		{token.Ident, "ten", 20},
		{token.RParen, ")", 20},
		{token.Semicolon, ";", 20},

		{token.Bang, "!", 22},
		{token.Minus, "-", 22},
		{token.Slash, "/", 22},
		{token.Asterisk, "*", 22},
		{token.Int, "5", 22},
		{token.Semicolon, ";", 22},

		{token.Int, "5", 23},
		{token.LT, "<", 23},
		{token.Int, "10", 23},
		{token.GT, ">", 23},
		{token.Int, "5", 23},
		{token.Semicolon, ";", 23},

		{token.If, "if", 25},
		{token.Int, "5", 25},
		{token.LT, "<", 25},
		{token.Int, "10", 25},
		{token.Return, "return", 26},
		{token.True, "true", 26},
		{token.Semicolon, ";", 26},
		{token.ElsIf, "elsif", 27},
		{token.Int, "10", 27},
		{token.Eq, "==", 27},
		{token.Int, "11", 27},
		{token.Return, "return", 28},
		{token.True, "true", 28},
		{token.Semicolon, ";", 28},
		{token.Else, "else", 29},
		{token.Return, "return", 30},
		{token.False, "false", 30},
		{token.Semicolon, ";", 30},
		{token.End, "end", 31},

		{token.String, "string1", 33},
		{token.Semicolon, ";", 33},
		{token.String, "string2", 34},
		{token.Semicolon, ";", 34},

		{token.Int, "10", 36},
		{token.Eq, "==", 36},
		{token.Int, "10", 36},
		{token.Semicolon, ";", 36},

		{token.Int, "10", 38},
		{token.NotEq, "!=", 38},
		{token.Int, "9", 38},
		{token.Semicolon, ";", 38},

		{token.Comment, "# This is comment.", 40},
		{token.Comment, "# And I should be ignored.", 41},

		{token.LBracket, "[", 43},
		{token.Int, "1", 43},
		{token.Comma, ",", 43},
		{token.Int, "2", 43},
		{token.Comma, ",", 43},
		{token.Int, "3", 43},
		{token.Comma, ",", 43},
		{token.Int, "4", 43},
		{token.Comma, ",", 43},
		{token.Int, "5", 43},
		{token.RBracket, "]", 43},

		{token.LBracket, "[", 44},
		{token.String, "test", 44},
		{token.Comma, ",", 44},
		{token.String, "test", 44},
		{token.RBracket, "]", 44},

		{token.LBrace, "{", 46},
		{token.Ident, "test", 46},
		{token.Colon, ":", 46},
		{token.String, "123", 46},
		{token.RBrace, "}", 46},

		{token.Ident, "require_relative", 48},
		{token.String, "foo", 48},

		{token.String, "a", 50},
		{token.Match, "=~", 50},
		{token.String, "a", 50},
		{token.Semicolon, ";", 50},

		{token.Int, "10", 51},
		{token.LTE, "<=", 51},
		{token.Int, "10", 51},
		{token.Semicolon, ";", 51},

		{token.Int, "10", 52},
		{token.GTE, ">=", 52},
		{token.Int, "10", 52},
		{token.Semicolon, ";", 52},

		{token.Ident, "a", 53},
		{token.Assign, "=", 53},
		{token.Int, "1", 53},
		{token.COMP, "<=>", 53},
		{token.Int, "2", 53},
		{token.Semicolon, ";", 53},

		{token.Int, "8", 55},
		{token.Pow, "**", 55},
		{token.Int, "10", 55},
		{token.Semicolon, ";", 55},

		{token.True, "true", 57},
		{token.And, "&&", 57},
		{token.False, "false", 57},
		{token.Semicolon, ";", 57},

		{token.False, "false", 58},
		{token.Or, "||", 58},
		{token.True, "true", 58},
		{token.Semicolon, ";", 58},

		{token.Null, "nil", 60},

		{token.Module, "module", 62},
		{token.Constant, "Foo", 62},
		{token.End, "end", 63},

		{token.Ident, "foo", 65},
		{token.Dot, ".", 65},
		{token.Ident, "module", 65},

		{token.Ident, "require", 67},
		{token.String, "foo", 67},

		{token.Constant, "Foo", 69},
		{token.ResolutionOperator, "::", 69},
		{token.Constant, "Bar", 69},

		{token.Constant, "Person", 71},
		{token.Dot, ".", 71},
		{token.Ident, "class", 71},

		{token.Class, "class", 73},
		{token.Constant, "Foo", 73},
		{token.Def, "def", 74},
		{token.Ident, "class", 74},
		{token.Constant, "Foo", 75},
		{token.End, "end", 76},
		{token.Def, "def", 77},
		{token.Self, "self", 77},
		{token.Dot, ".", 77},
		{token.Ident, "bar", 77},
		{token.Int, "10", 78},
		{token.End, "end", 79},
		{token.End, "end", 80},

		{token.Int, "10", 82},
		{token.Modulo, "%", 82},
		{token.Int, "5", 82},

		{token.String, "", 84},
		{token.String, "", 85},

		{token.Next, "next", 87},
		{token.String, "apple", 88},

		{token.LBrace, "{", 89},
		{token.Ident, "test", 89},
		{token.Colon, ":", 89},
		{token.String, "abc", 89},
		{token.RBrace, "}", 89},

		{token.LBrace, "{", 90},
		{token.Ident, "test", 90},
		{token.Colon, ":", 90},
		{token.String, "abc", 90},
		{token.RBrace, "}", 90},

		{token.LBrace, "{", 91},
		{token.Ident, "test", 91},
		{token.Colon, ":", 91},
		{token.Int, "50", 91},
		{token.RBrace, "}", 91},

		{token.LBrace, "{", 92},
		{token.Ident, "test", 92},
		{token.Colon, ":", 92},
		{token.Ident, "abc", 92},
		{token.RBrace, "}", 92},

		{token.LBrace, "{", 93},
		{token.Ident, "test", 93},
		{token.Colon, ":", 93},
		{token.Ident, "abc", 93},
		{token.RBrace, "}", 93},

		{token.LParen, "(", 95},
		{token.Int, "1", 95},
		{token.Range, "..", 95},
		{token.Int, "5", 95},
		{token.RParen, ")", 95},

		{token.While, "while", 97},
		{token.Ident, "i", 97},
		{token.LT, "<", 97},
		{token.Int, "10", 97},
		{token.Do, "do", 97},
		{token.Break, "break", 98},
		{token.End, "end", 99},

		{token.Ident, "a", 101},
		{token.PlusEq, "+=", 101},
		{token.Int, "1", 101},

		{token.Ident, "b", 102},
		{token.MinusEq, "-=", 102},
		{token.Int, "2", 102},

		{token.Ident, "c", 103},
		{token.OrEq, "||=", 103},
		{token.True, "true", 103},

		// Escaped character tests
		{token.String, "\nstring\n", 105},
		{token.String, "\\nstring\\n", 106},
		{token.String, "\tstring\t", 107},
		{token.String, "\\tstring\\t", 108},
		{token.String, "\vstring\v", 109},
		{token.String, "\\vstring\\v", 110},
		{token.String, "\fstring\f", 111},
		{token.String, "\\fstring\\f", 112},
		{token.String, "\rstring\r", 113},
		{token.String, "\\rstring\\r", 114},
		{token.String, "\"string\"", 115},
		{token.String, "\\\"string\\\"", 116},
		{token.String, "'string'", 117},
		{token.String, "'string'", 118},
		{token.Int, "123", 119},
		{token.Int, "123", 120},
		{token.Dot, ".", 120},
		{token.Int, "456", 120},
		{token.EOF, "", 121},
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
