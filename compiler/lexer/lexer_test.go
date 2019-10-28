package lexer

import (
	"github.com/goby-lang/goby/compiler/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input   string
		expects []struct {
			expectedType    token.Type
			expectedLiteral string
			expectedLine    int
		}
	}{
		{
			`
				five = 5;
				ten = 10;
				_ = 15
			`,
			[]struct {
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
			},
		},
		{
			`
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
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				// class is default to be ident
				{token.Class, "class", 1},
				{token.Constant, "Person", 1},

				{token.Def, "def", 2},
				{token.Ident, "initialize", 2},
				{token.LParen, "(", 2},
				{token.Ident, "a", 2},
				{token.RParen, ")", 2},
				{token.InstanceVariable, "@a", 3},
				{token.Assign, "=", 3},
				{token.Ident, "a", 3},
				{token.Semicolon, ";", 3},
				{token.End, "end", 4},

				{token.Def, "def", 6},
				{token.Ident, "add", 6},
				{token.LParen, "(", 6},
				{token.Ident, "x", 6},
				{token.Comma, ",", 6},
				{token.Ident, "y", 6},
				{token.RParen, ")", 6},
				{token.Ident, "x", 7},
				{token.Plus, "+", 7},
				{token.Ident, "y", 7},
				{token.Semicolon, ";", 7},
				{token.End, "end", 8},

				{token.Def, "def", 10},
				{token.Ident, "ten", 10},
				{token.LParen, "(", 10},
				{token.RParen, ")", 10},
				{token.Self, "self", 11},
				{token.Dot, ".", 11},
				{token.Ident, "add", 11},
				{token.LParen, "(", 11},
				{token.Int, "1", 11},
				{token.Comma, ",", 11},
				{token.Int, "9", 11},
				{token.RParen, ")", 11},
				{token.Semicolon, ";", 11},
				{token.End, "end", 12},
				{token.End, "end", 13},

				{token.Ident, "p", 15},
				{token.Assign, "=", 15},
				{token.Constant, "Person", 15},
				{token.Dot, ".", 15},
				{token.Ident, "new", 15},
				{token.Semicolon, ";", 15},
				{token.Ident, "result", 16},
				{token.Assign, "=", 16},
				{token.Ident, "p", 16},
				{token.Dot, ".", 16},
				{token.Ident, "add", 16},
				{token.LParen, "(", 16},
				{token.Ident, "five", 16},
				{token.Comma, ",", 16},
				{token.Ident, "ten", 16},
				{token.RParen, ")", 16},
				{token.Semicolon, ";", 16},
			},
		},
		{
			`
	!-/*5;
	5 < 10 >5;
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Bang, "!", 1},
				{token.Minus, "-", 1},
				{token.Slash, "/", 1},
				{token.Asterisk, "*", 1},
				{token.Int, "5", 1},
				{token.Semicolon, ";", 1},

				{token.Int, "5", 2},
				{token.LT, "<", 2},
				{token.Int, "10", 2},
				{token.GT, ">", 2},
				{token.Int, "5", 2},
				{token.Semicolon, ";", 2},
			},
		},
		{
			`
	if 5 < 10
	 return true;
	elsif 10 == 11
	 return true;
	else
	 return false;
	end
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.If, "if", 1},
				{token.Int, "5", 1},
				{token.LT, "<", 1},
				{token.Int, "10", 1},
				{token.Return, "return", 2},
				{token.True, "true", 2},
				{token.Semicolon, ";", 2},
				{token.ElsIf, "elsif", 3},
				{token.Int, "10", 3},
				{token.Eq, "==", 3},
				{token.Int, "11", 3},
				{token.Return, "return", 4},
				{token.True, "true", 4},
				{token.Semicolon, ";", 4},
				{token.Else, "else", 5},
				{token.Return, "return", 6},
				{token.False, "false", 6},
				{token.Semicolon, ";", 6},
				{token.End, "end", 7},
			},
		}, {
			`
	"string1";
	'string2';			
`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.String, "string1", 1},
				{token.Semicolon, ";", 1},
				{token.String, "string2", 2},
				{token.Semicolon, ";", 2},
			},
		}, {
			`
	10 == 10;
	10 != 9;
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{

				{token.Int, "10", 1},
				{token.Eq, "==", 1},
				{token.Int, "10", 1},
				{token.Semicolon, ";", 1},

				{token.Int, "10", 2},
				{token.NotEq, "!=", 2},
				{token.Int, "9", 2},
				{token.Semicolon, ";", 2},
			},
		}, {
			`
	# This is comment.
	# And I should be ignored.
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Comment, "# This is comment.", 1},
				{token.Comment, "# And I should be ignored.", 2},
			},
		}, {
			`
	[1, 2, 3, 4, 5]
	["test", "test"]
	{ test: "12" }
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.LBracket, "[", 1},
				{token.Int, "1", 1},
				{token.Comma, ",", 1},
				{token.Int, "2", 1},
				{token.Comma, ",", 1},
				{token.Int, "3", 1},
				{token.Comma, ",", 1},
				{token.Int, "4", 1},
				{token.Comma, ",", 1},
				{token.Int, "5", 1},
				{token.RBracket, "]", 1},

				{token.LBracket, "[", 2},
				{token.String, "test", 2},
				{token.Comma, ",", 2},
				{token.String, "test", 2},
				{token.RBracket, "]", 2},

				{token.LBrace, "{", 3},
				{token.Ident, "test", 3},
				{token.Colon, ":", 3},
				{token.String, "12", 3},
				{token.RBrace, "}", 3},
			},
		}, {
			`
	require_relative "foo"
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Ident, "require_relative", 1},
				{token.String, "foo", 1},
			},
		}, {
			`
	'a' =~ 'a';
	10 <= 10;
	10 >= 10;
	a = 1 <=> 2;
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.String, "a", 1},
				{token.Match, "=~", 1},
				{token.String, "a", 1},
				{token.Semicolon, ";", 1},

				{token.Int, "10", 2},
				{token.LTE, "<=", 2},
				{token.Int, "10", 2},
				{token.Semicolon, ";", 2},

				{token.Int, "10", 3},
				{token.GTE, ">=", 3},
				{token.Int, "10", 3},
				{token.Semicolon, ";", 3},

				{token.Ident, "a", 4},
				{token.Assign, "=", 4},
				{token.Int, "1", 4},
				{token.COMP, "<=>", 4},
				{token.Int, "2", 4},
				{token.Semicolon, ";", 4},
			},
		}, {
			`
	8 ** 10;
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Int, "8", 1},
				{token.Pow, "**", 1},
				{token.Int, "10", 1},
				{token.Semicolon, ";", 1},
			},
		}, {
			`
	true && false;
	false || true;			
`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.True, "true", 1},
				{token.And, "&&", 1},
				{token.False, "false", 1},
				{token.Semicolon, ";", 1},

				{token.False, "false", 2},
				{token.Or, "||", 2},
				{token.True, "true", 2},
				{token.Semicolon, ";", 2},
			},
		}, {
			`
	nil
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Null, "nil", 1},
			},
		}, {
			`
	module Foo
	end
	
	foo.module
	
	require "foo"
	
	Foo::Bar
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Module, "module", 1},
				{token.Constant, "Foo", 1},
				{token.End, "end", 2},

				{token.Ident, "foo", 4},
				{token.Dot, ".", 4},
				{token.Ident, "module", 4},

				{token.Ident, "require", 6},
				{token.String, "foo", 6},

				{token.Constant, "Foo", 8},
				{token.ResolutionOperator, "::", 8},
				{token.Constant, "Bar", 8},
			},
		}, {
			`
	Person.class
	class Foo
	  def class
		Foo
	  end
	  def self.bar
		10
	  end
	end
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Constant, "Person", 1},
				{token.Dot, ".", 1},
				{token.Ident, "class", 1},

				{token.Class, "class", 2},
				{token.Constant, "Foo", 2},
				{token.Def, "def", 3},
				{token.Ident, "class", 3},
				{token.Constant, "Foo", 4},
				{token.End, "end", 5},
				{token.Def, "def", 6},
				{token.Self, "self", 6},
				{token.Dot, ".", 6},
				{token.Ident, "bar", 6},
				{token.Int, "10", 7},
				{token.End, "end", 8},
				{token.End, "end", 9},
			},
		}, {
			`
	10 % 5
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Int, "10", 1},
				{token.Modulo, "%", 1},
				{token.Int, "5", 1},
			},
		}, {
			`
	""
	''
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.String, "", 1},
				{token.String, "", 2},
			},
		}, {
			`
	next
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Next, "next", 1},
			},
		},
		{
			`
	:apple

			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.String, "apple", 1},
			},
		}, {
			`
	{ test:"abc" }
	{ test: :abc }
	{ test:1 }
	{ test: abc }
	{ test:abc }
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.LBrace, "{", 1},
				{token.Ident, "test", 1},
				{token.Colon, ":", 1},
				{token.String, "abc", 1},
				{token.RBrace, "}", 1},

				{token.LBrace, "{", 2},
				{token.Ident, "test", 2},
				{token.Colon, ":", 2},
				{token.String, "abc", 2},
				{token.RBrace, "}", 2},

				{token.LBrace, "{", 3},
				{token.Ident, "test", 3},
				{token.Colon, ":", 3},
				{token.Int, "1", 3},
				{token.RBrace, "}", 3},

				{token.LBrace, "{", 4},
				{token.Ident, "test", 4},
				{token.Colon, ":", 4},
				{token.Ident, "abc", 4},
				{token.RBrace, "}", 4},

				{token.LBrace, "{", 5},
				{token.Ident, "test", 5},
				{token.Colon, ":", 5},
				{token.Ident, "abc", 5},
				{token.RBrace, "}", 5},
			},
		}, {
			`
	(1..5)
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{

				{token.LParen, "(", 1},
				{token.Int, "1", 1},
				{token.Range, "..", 1},
				{token.Int, "5", 1},
				{token.RParen, ")", 1},
			},
		}, {
			`
	while i < 10 do
	 break
	end
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.While, "while", 1},
				{token.Ident, "i", 1},
				{token.LT, "<", 1},
				{token.Int, "10", 1},
				{token.Do, "do", 1},
				{token.Break, "break", 2},
				{token.End, "end", 3},
			},
		}, {
			`
	a += 1
	b -= 2
	c ||= true
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				{token.Ident, "a", 1},
				{token.PlusEq, "+=", 1},
				{token.Int, "1", 1},

				{token.Ident, "b", 2},
				{token.MinusEq, "-=", 2},
				{token.Int, "2", 2},

				{token.Ident, "c", 3},
				{token.OrEq, "||=", 3},
				{token.True, "true", 3},
			},
		}, {
			`
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
			`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
				expectedLine    int
			}{
				// Escaped character tests
				{token.String, "\nstring\n", 1},
				{token.String, "\\nstring\\n", 2},
				{token.String, "\tstring\t", 3},
				{token.String, "\\tstring\\t", 4},
				{token.String, "\vstring\v", 5},
				{token.String, "\\vstring\\v", 6},
				{token.String, "\fstring\f", 7},
				{token.String, "\\fstring\\f", 8},
				{token.String, "\rstring\r", 9},
				{token.String, "\\rstring\\r", 10},
				{token.String, "\"string\"", 11},
				{token.String, "\\\"string\\\"", 12},
				{token.String, "'string'", 13},
				{token.String, "'string'", 14},

				{token.EOF, "", 15},
			},
		},
	}

	for i, tt := range tests {
		l := New(tt.input)

		for _, expect := range tt.expects {
			tok := l.NextToken()

			if tok.Type != expect.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, expect.expectedType, tok.Type)
			}
			if tok.Literal != expect.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, expect.expectedLiteral, tok.Literal)

			}
			if tok.Line != expect.expectedLine {
				t.Fatalf("tests[%d] - line number wrong. expected=%d, got=%d", i, expect.expectedLine, tok.Line)
			}
		}
	}
}
