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

	a++
	b--

	while i < 10 do
	  puts(i)
	  i++
	end

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

		{token.Ident, "a", 48},
		{token.Incr, "++", 48},
		{token.Ident, "b", 49},
		{token.Decr, "--", 49},

		{token.While, "while", 51},
		{token.Ident, "i", 51},
		{token.LT, "<", 51},
		{token.Int, "10", 51},
		{token.Do, "do", 51},
		{token.Ident, "puts", 52},
		{token.LParen, "(", 52},
		{token.Ident, "i", 52},
		{token.RParen, ")", 52},
		{token.Ident, "i", 53},
		{token.Incr, "++", 53},
		{token.End, "end", 54},

		{token.Ident, "require_relative", 56},
		{token.String, "foo", 56},

		{token.String, "a", 58},
		{token.Match, "=~", 58},
		{token.String, "a", 58},
		{token.Semicolon, ";", 58},

		{token.Int, "10", 59},
		{token.LTE, "<=", 59},
		{token.Int, "10", 59},
		{token.Semicolon, ";", 59},

		{token.Int, "10", 60},
		{token.GTE, ">=", 60},
		{token.Int, "10", 60},
		{token.Semicolon, ";", 60},

		{token.Ident, "a", 61},
		{token.Assign, "=", 61},
		{token.Int, "1", 61},
		{token.COMP, "<=>", 61},
		{token.Int, "2", 61},
		{token.Semicolon, ";", 61},

		{token.Int, "8", 63},
		{token.Pow, "**", 63},
		{token.Int, "10", 63},
		{token.Semicolon, ";", 63},

		{token.True, "true", 65},
		{token.And, "&&", 65},
		{token.False, "false", 65},
		{token.Semicolon, ";", 65},

		{token.False, "false", 66},
		{token.Or, "||", 66},
		{token.True, "true", 66},
		{token.Semicolon, ";", 66},

		{token.Null, "nil", 68},

		{token.Module, "module", 70},
		{token.Constant, "Foo", 70},
		{token.End, "end", 71},

		{token.Ident, "foo", 73},
		{token.Dot, ".", 73},
		{token.Ident, "module", 73},

		{token.Ident, "require", 75},
		{token.String, "foo", 75},

		{token.Constant, "Foo", 77},
		{token.ResolutionOperator, "::", 77},
		{token.Constant, "Bar", 77},

		{token.Constant, "Person", 79},
		{token.Dot, ".", 79},
		{token.Ident, "class", 79},

		{token.Class, "class", 81},
		{token.Constant, "Foo", 81},
		{token.Def, "def", 82},
		{token.Ident, "class", 82},
		{token.Constant, "Foo", 83},
		{token.End, "end", 84},
		{token.Def, "def", 85},
		{token.Self, "self", 85},
		{token.Dot, ".", 85},
		{token.Ident, "bar", 85},
		{token.Int, "10", 86},
		{token.End, "end", 87},
		{token.End, "end", 88},

		{token.Int, "10", 90},
		{token.Modulo, "%", 90},
		{token.Int, "5", 90},

		{token.String, "", 92},
		{token.String, "", 93},

		{token.Next, "next", 95},
		{token.String, "apple", 96},

		{token.LBrace, "{", 97},
		{token.Ident, "test", 97},
		{token.Colon, ":", 97},
		{token.String, "abc", 97},
		{token.RBrace, "}", 97},

		{token.LBrace, "{", 98},
		{token.Ident, "test", 98},
		{token.Colon, ":", 98},
		{token.String, "abc", 98},
		{token.RBrace, "}", 98},

		{token.LBrace, "{", 99},
		{token.Ident, "test", 99},
		{token.Colon, ":", 99},
		{token.Int, "50", 99},
		{token.RBrace, "}", 99},

		{token.LBrace, "{", 100},
		{token.Ident, "test", 100},
		{token.Colon, ":", 100},
		{token.Ident, "abc", 100},
		{token.RBrace, "}", 100},

		{token.LBrace, "{", 101},
		{token.Ident, "test", 101},
		{token.Colon, ":", 101},
		{token.Ident, "abc", 101},
		{token.RBrace, "}", 101},

		{token.LParen, "(", 103},
		{token.Int, "1", 103},
		{token.Range, "..", 103},
		{token.Int, "5", 103},
		{token.RParen, ")", 103},

		{token.While, "while", 105},
		{token.Ident, "i", 105},
		{token.LT, "<", 105},
		{token.Int, "10", 105},
		{token.Do, "do", 105},
		{token.Break, "break", 106},
		{token.End, "end", 107},

		{token.Ident, "a", 109},
		{token.PlusEq, "+=", 109},
		{token.Int, "1", 109},

		{token.Ident, "b", 110},
		{token.MinusEq, "-=", 110},
		{token.Int, "2", 110},

		{token.Ident, "c", 111},
		{token.OrEq, "||=", 111},
		{token.True, "true", 111},

		// Escaped character tests
		{token.String, "\nstring\n", 113},
		{token.String, "\\nstring\\n", 114},
		{token.String, "\tstring\t", 115},
		{token.String, "\\tstring\\t", 116},
		{token.String, "\vstring\v", 117},
		{token.String, "\\vstring\\v", 118},
		{token.String, "\fstring\f", 119},
		{token.String, "\\fstring\\f", 120},
		{token.String, "\rstring\r", 121},
		{token.String, "\\rstring\\r", 122},
		{token.String, "\"string\"", 123},
		{token.String, "\\\"string\\\"", 124},
		{token.String, "'string'", 125},
		{token.String, "'string'", 126},

		{token.EOF, "", 127},
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
