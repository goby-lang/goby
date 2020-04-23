package ripper

import (
	"strings"
	"testing"

	"github.com/goby-lang/goby/vm"
	"regexp"
)

type errorTestCase struct {
	input       string
	expected    string
	expectedCFP int
}

func TestRipperClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`require 'ripper'; Ripper.class.name`, "Class"},
		{`require 'ripper'; Ripper.superclass.name`, "Object"},
		{`require 'ripper'; Ripper.ancestors.to_s`, "[Ripper, Object]"},
	}

	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestRipperClassCreationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require 'ripper'; Ripper.new`, "NoMethodError: Undefined Method 'new' for Ripper", 1},
	}

	for i, tt := range testsFail {
		evaluated := vm.ExecAndReturn(t, tt.input)
		checkErrorMsg(t, i, evaluated, tt.expected)
	}
}

func TestRipperParse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`require 'ripper'; Ripper.parse "
	class Bar
		def self.foo
			10
		end
	end
	class Foo < Bar; end
	class FooBar < Foo; end
	FooBar.foo
"`, "class Bar {\ndef foo() {\n10\n}\n}class Foo {\n\n}class FooBar {\n\n}FooBar.foo()"},
		{`require 'ripper'; Ripper.parse "
	def foo(x)
	  yield(x + 10)
	end
	def bar(y)
	  foo(y) do |f|
		yield(f)
	  end
	end
	def baz(z)
	  bar(z + 100) do |b|
		yield(b)
	  end
	end
	a = 0
	baz(100) do |b|
	  a = b
	end
	a

	class Foo
	  def bar
		100
	  end
	end
	module Baz
	  class Bar
		def bar
		  Foo.new.bar
		end
	  end
	end
	Baz::Bar.new.bar + a
"`, "def foo(x) {\nyield((x + 10))\n}def bar(y) {\nself.foo(y) do |f|\nyield(f)\nend\n}def baz(z) {\nself.bar((z + 100)) do |b|\nyield(b)\nend\n}a = 0self.baz(100) do |b|\na = b\nendaclass Foo {\ndef bar() {\n100\n}\n}module Baz {\nclass Bar {\ndef bar() {\nFoo.new().bar()\n}\n}\n}((Baz :: Bar).new().bar() + a)"},
		{`require 'ripper'; Ripper.parse "
	def bar(block)
	block.call + get_block.call
	end
	
	def foo
		bar(get_block) do
  		20
		end
	end
	
	foo do
		10
	end
"`, "def bar(block) {\n(block.call() + get_block.call())\n}def foo() {\nself.bar(get_block) do\n20\nend\n}self.foo() do\n10\nend"},
	}

	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestRipperParseFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require 'ripper'; Ripper.parse`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`require 'ripper'; Ripper.parse(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`require 'ripper'; Ripper.parse(1.2)`, "TypeError: Expect argument to be String. got: Float", 1},
		{`require 'ripper'; Ripper.parse(["puts", "123"])`, "TypeError: Expect argument to be String. got: Array", 1},
		{`require 'ripper'; Ripper.parse({key: 1})`, "TypeError: Expect argument to be String. got: Hash", 1},
	}

	for i, tt := range testsFail {
		evaluated := vm.ExecAndReturn(t, tt.input)
		checkErrorMsg(t, i, evaluated, tt.expected)
	}
}

func TestRipperTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`require 'ripper'; Ripper.tokenize("
	class Bar
		def self.foo
			10
		end
	end
	class Foo < Bar; end
	class FooBar < Foo; end
	FooBar.foo
").to_s`, `["class", "Bar", "def", "self", ".", "foo", "10", "end", "end", "class", "Foo", "<", "Bar", ";", "end", "class", "FooBar", "<", "Foo", ";", "end", "FooBar", ".", "foo", "EOF"]`},
		{`require 'ripper'; Ripper.tokenize("
	def foo(x)
	  yield(x + 10)
	end
	def bar(y)
	  foo(y) do |f|
		yield(f)
	  end
	end
	def baz(z)
	  bar(z + 100) do |b|
		yield(b)
	  end
	end
	a = 0
	baz(100) do |b|
	  a = b
	end
	a

	class Foo
	  def bar
		100
	  end
	end
	module Baz
	  class Bar
		def bar
		  Foo.new.bar
		end
	  end
	end
	Baz::Bar.new.bar + a
").to_s`, `["def", "foo", "(", "x", ")", "yield", "(", "x", "+", "10", ")", "end", "def", "bar", "(", "y", ")", "foo", "(", "y", ")", "do", "|", "f", "|", "yield", "(", "f", ")", "end", "end", "def", "baz", "(", "z", ")", "bar", "(", "z", "+", "100", ")", "do", "|", "b", "|", "yield", "(", "b", ")", "end", "end", "a", "=", "0", "baz", "(", "100", ")", "do", "|", "b", "|", "a", "=", "b", "end", "a", "class", "Foo", "def", "bar", "100", "end", "end", "module", "Baz", "class", "Bar", "def", "bar", "Foo", ".", "new", ".", "bar", "end", "end", "end", "Baz", "::", "Bar", ".", "new", ".", "bar", "+", "a", "EOF"]`},
		{`require 'ripper'; Ripper.tokenize("
	def bar(block)
	block.call + get_block.call
	end
	
	def foo
		bar(get_block) do
  		20
		end
	end
	
	foo do
		10
	end
").to_s`, `["def", "bar", "(", "block", ")", "block", ".", "call", "+", "get_block", ".", "call", "end", "def", "foo", "bar", "(", "get_block", ")", "do", "20", "end", "end", "foo", "do", "10", "end", "EOF"]`},
	}

	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestRipperTokenFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require 'ripper'; Ripper.tokenize`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`require 'ripper'; Ripper.tokenize(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`require 'ripper'; Ripper.tokenize(1.2)`, "TypeError: Expect argument to be String. got: Float", 1},
		{`require 'ripper'; Ripper.tokenize(["puts", "123"])`, "TypeError: Expect argument to be String. got: Array", 1},
		{`require 'ripper'; Ripper.tokenize({key: 1})`, "TypeError: Expect argument to be String. got: Hash", 1},
	}

	for i, tt := range testsFail {
		evaluated := vm.ExecAndReturn(t, tt.input)
		checkErrorMsg(t, i, evaluated, tt.expected)
	}
}

func TestRipperLex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`require 'ripper'; Ripper.lex("
	class Bar
		def self.foo
			10
		end
	end
	class Foo < Bar; end
	class FooBar < Foo; end
	FooBar.foo
").to_s`, `[[1, "on_class", "class"], [1, "on_constant", "Bar"], [2, "on_def", "def"], [2, "on_self", "self"], [2, "on_dot", "."], [2, "on_ident", "foo"], [3, "on_int", "10"], [4, "on_end", "end"], [5, "on_end", "end"], [6, "on_class", "class"], [6, "on_constant", "Foo"], [6, "on_lt", "<"], [6, "on_constant", "Bar"], [6, "on_semicolon", ";"], [6, "on_end", "end"], [7, "on_class", "class"], [7, "on_constant", "FooBar"], [7, "on_lt", "<"], [7, "on_constant", "Foo"], [7, "on_semicolon", ";"], [7, "on_end", "end"], [8, "on_constant", "FooBar"], [8, "on_dot", "."], [8, "on_ident", "foo"], [9, "on_eof", ""]]`},
		{`require 'ripper'; Ripper.lex("
	def foo(x)
	  yield(x + 10)
	end
	def bar(y)
	  foo(y) do |f|
		yield(f)
	  end
	end
	def baz(z)
	  bar(z + 100) do |b|
		yield(b)
	  end
	end
	a = 0
	baz(100) do |b|
	  a = b
	end
	a

	class Foo
	  def bar
		100
	  end
	end
	module Baz
	  class Bar
		def bar
		  Foo.new.bar
		end
	  end
	end
	Baz::Bar.new.bar + a
").to_s`, `[[1, "on_def", "def"], [1, "on_ident", "foo"], [1, "on_lparen", "("], [1, "on_ident", "x"], [1, "on_rparen", ")"], [2, "on_yield", "yield"], [2, "on_lparen", "("], [2, "on_ident", "x"], [2, "on_plus", "+"], [2, "on_int", "10"], [2, "on_rparen", ")"], [3, "on_end", "end"], [4, "on_def", "def"], [4, "on_ident", "bar"], [4, "on_lparen", "("], [4, "on_ident", "y"], [4, "on_rparen", ")"], [5, "on_ident", "foo"], [5, "on_lparen", "("], [5, "on_ident", "y"], [5, "on_rparen", ")"], [5, "on_do", "do"], [5, "on_bar", "|"], [5, "on_ident", "f"], [5, "on_bar", "|"], [6, "on_yield", "yield"], [6, "on_lparen", "("], [6, "on_ident", "f"], [6, "on_rparen", ")"], [7, "on_end", "end"], [8, "on_end", "end"], [9, "on_def", "def"], [9, "on_ident", "baz"], [9, "on_lparen", "("], [9, "on_ident", "z"], [9, "on_rparen", ")"], [10, "on_ident", "bar"], [10, "on_lparen", "("], [10, "on_ident", "z"], [10, "on_plus", "+"], [10, "on_int", "100"], [10, "on_rparen", ")"], [10, "on_do", "do"], [10, "on_bar", "|"], [10, "on_ident", "b"], [10, "on_bar", "|"], [11, "on_yield", "yield"], [11, "on_lparen", "("], [11, "on_ident", "b"], [11, "on_rparen", ")"], [12, "on_end", "end"], [13, "on_end", "end"], [14, "on_ident", "a"], [14, "on_assign", "="], [14, "on_int", "0"], [15, "on_ident", "baz"], [15, "on_lparen", "("], [15, "on_int", "100"], [15, "on_rparen", ")"], [15, "on_do", "do"], [15, "on_bar", "|"], [15, "on_ident", "b"], [15, "on_bar", "|"], [16, "on_ident", "a"], [16, "on_assign", "="], [16, "on_ident", "b"], [17, "on_end", "end"], [18, "on_ident", "a"], [20, "on_class", "class"], [20, "on_constant", "Foo"], [21, "on_def", "def"], [21, "on_ident", "bar"], [22, "on_int", "100"], [23, "on_end", "end"], [24, "on_end", "end"], [25, "on_module", "module"], [25, "on_constant", "Baz"], [26, "on_class", "class"], [26, "on_constant", "Bar"], [27, "on_def", "def"], [27, "on_ident", "bar"], [28, "on_constant", "Foo"], [28, "on_dot", "."], [28, "on_ident", "new"], [28, "on_dot", "."], [28, "on_ident", "bar"], [29, "on_end", "end"], [30, "on_end", "end"], [31, "on_end", "end"], [32, "on_constant", "Baz"], [32, "on_resolutionoperator", "::"], [32, "on_constant", "Bar"], [32, "on_dot", "."], [32, "on_ident", "new"], [32, "on_dot", "."], [32, "on_ident", "bar"], [32, "on_plus", "+"], [32, "on_ident", "a"], [33, "on_eof", ""]]`},
		{`require 'ripper'; Ripper.lex("
	def bar(block)
	block.call + get_block.call
	end
	
	def foo
		bar(get_block) do
  		20
		end
	end
	
	foo do
		10
	end
").to_s`, `[[1, "on_def", "def"], [1, "on_ident", "bar"], [1, "on_lparen", "("], [1, "on_ident", "block"], [1, "on_rparen", ")"], [2, "on_ident", "block"], [2, "on_dot", "."], [2, "on_ident", "call"], [2, "on_plus", "+"], [2, "on_get_block", "get_block"], [2, "on_dot", "."], [2, "on_ident", "call"], [3, "on_end", "end"], [5, "on_def", "def"], [5, "on_ident", "foo"], [6, "on_ident", "bar"], [6, "on_lparen", "("], [6, "on_get_block", "get_block"], [6, "on_rparen", ")"], [6, "on_do", "do"], [7, "on_int", "20"], [8, "on_end", "end"], [9, "on_end", "end"], [11, "on_ident", "foo"], [11, "on_do", "do"], [12, "on_int", "10"], [13, "on_end", "end"], [14, "on_eof", ""]]`},
	}

	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)
		vm.VerifyExpected(t, i, evaluated, tt.expected)
	}
}

func TestRipperLexFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require 'ripper'; Ripper.lex`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`require 'ripper'; Ripper.lex(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`require 'ripper'; Ripper.lex(1.2)`, "TypeError: Expect argument to be String. got: Float", 1},
		{`require 'ripper'; Ripper.lex(["puts", "123"])`, "TypeError: Expect argument to be String. got: Array", 1},
		{`require 'ripper'; Ripper.lex({key: 1})`, "TypeError: Expect argument to be String. got: Hash", 1},
	}

	for i, tt := range testsFail {
		evaluated := vm.ExecAndReturn(t, tt.input)
		checkErrorMsg(t, i, evaluated, tt.expected)
	}
}

func TestRipperInstruction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Single method definition
		{`require 'ripper'; Ripper.instruction("
def foo
  10
end
  ").to_s`, `
[
  { 
    arg_types: { names: [], types: [] }, 
    instructions: [
      { action: "putobject", line: 0, params: ["10"], source_line: 3 }, 
      { action: "leave", line: 1, params: [], source_line: 2 }
    ], 
    name: "foo", 
    type: "Def" 
  }, 
  {
     arg_types: { names: [], types: [] }, 
     instructions: [
       { action: "putobject", line: 0, params: ["10"], source_line: 3 }, 
       { action: "leave", line: 1, params: [], source_line: 2 }
     ], 
     name: "foo", 
     type: "Def" 
  }, 
  { 
     instructions: [
       { action: "putself", line: 0, params: [], source_line: 2 }, 
       { action: "putstring", line: 1, params: ["foo"], source_line: 2 }, 
       { action: "def_method", line: 2, params: ["0"], source_line: 2 }, 
       { action: "leave", line: 3, params: [], source_line: 2 }
     ], 
     name: "ProgramStart", 
     type: "ProgramStart" 
  }, 
  { 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 2 }, 
      { action: "putstring", line: 1, params: ["foo"], source_line: 2 }, 
      { action: "def_method", line: 2, params: ["0"], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  }
]
`,
	},
		// Single class definition
		{`require 'ripper'; Ripper.instruction("
    class Foo
		end
		 ").to_s`, `
[
  { 
     instructions: [
       { action: "leave", line: 0, params: [], source_line: 2 }
     ], 
     name: "Foo", 
     type: "DefClass" 
  }, 
  { 
     instructions: [
       { action: "leave", line: 0, params: [], source_line: 2 }
     ], 
     name: "Foo", 
     type: "DefClass" 
  }, 
  { 
     instructions: [
       { action: "putself", line: 0, params: [], source_line: 2 }, 
       { action: "def_class", line: 1, params: ["class:Foo"], source_line: 2 }, 
       { action: "pop", line: 2, params: [], source_line: 2 }, 
       { action: "leave", line: 3, params: [], source_line: 2 }
     ], name: "ProgramStart", type: "ProgramStart" 
  },
  { 
     instructions: [
       { action: "putself", line: 0, params: [], source_line: 2 }, 
       { action: "def_class", line: 1, params: ["class:Foo"], source_line: 2 }, 
       { action: "pop", line: 2, params: [], source_line: 2 }, 
       { action: "leave", line: 3, params: [], source_line: 2 }], name: "ProgramStart", type: "ProgramStart" 
  }
]
`},
		// Single method call
		{`require 'ripper'; Ripper.instruction("
    Array.methods
		 ").to_s`, `
[
  { 
    arg_set: { names: [], types: [] }, 
    instructions: [
      { action: "getconstant", line: 0, params: ["Array", "false"], source_line: 2 }, 
      { action: "send", line: 1, params: ["methods", "0", "", "&{[] []}"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  }, 
  { 
    arg_set: { names: [], types: [] }, 
    instructions: [
      { action: "getconstant", line: 0, params: ["Array", "false"], source_line: 2 }, 
      { action: "send", line: 1, params: ["methods", "0", "", "&{[] []}"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ],
    name: "ProgramStart", 
    type: "ProgramStart" 
  }
]`},
		// Single method call with a block
		{`require 'ripper'; Ripper.instruction("
    10.times do |i| end
		 ").to_s`, `
[
  { 
		arg_types:{names:["i"],types:[0]},
    instructions: [
      { action: "leave", line: 0, params: [], source_line: 2 }
    ], 
    name: "0", 
    type: "Block" 
  }, 
  {
		arg_types:{names:["i"],types:[0]},
		instructions: [
		 { action: "leave", line: 0, params: [], source_line: 2 }
		], 
		name: "0", 
		type: "Block" 
  }, 
  { 
     arg_set: { names: [], types: [] }, 
     instructions: [
       { action: "putobject", line: 0, params: ["10"], source_line: 2 }, 
       { action: "send", line: 1, params: ["times", "0", "block:0", "&{[] []}"], source_line: 2 }, 
       { action: "pop", line: 2, params: [], source_line: 2 }, { action: "leave", line: 3, params: [], source_line: 2 }
     ],
     name: "ProgramStart", 
     type: "ProgramStart" 
  }, 
  { 
    arg_set: { names: [], types: [] }, 
    instructions: [
      { action: "putobject", line: 0, params: ["10"], source_line: 2 }, 
      { action: "send", line: 1, params: ["times", "0", "block:0", "&{[] []}"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  }
]`},
		// Single module definition
		{`require 'ripper'; Ripper.instruction("
		module Bar
    end
		 ").to_s`, `
[
  { instructions: [{ action: "leave", line: 0, params: [], source_line: 2 }], name: "Bar", type: "DefClass" }, 
  { instructions: [{ action: "leave", line: 0, params: [], source_line: 2 }], name: "Bar", type: "DefClass" }, 
  { 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 2 }, 
      { action: "def_class", line: 1, params: ["module:Bar"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  },
  { 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 2 }, 
      { action: "def_class", line: 1, params: ["module:Bar"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  }
]`},
		// Single assignment
		{`require 'ripper'; Ripper.instruction("
		 a = 1
		 ").to_s`, `
[
  { instructions: [{ action: "putobject", line: 0, params: ["1"], source_line: 2 }, { action: "setlocal", line: 1, params: ["0", "0"], source_line: 2 }, { action: "pop", line: 2, params: [], source_line: 2 }, { action: "leave", line: 3, params: [], source_line: 2 }], name: "ProgramStart", type: "ProgramStart" }, 
  { instructions: [{ action: "putobject", line: 0, params: ["1"], source_line: 2 }, { action: "setlocal", line: 1, params: ["0", "0"], source_line: 2 }, { action: "pop", line: 2, params: [], source_line: 2 }, { action: "leave", line: 3, params: [], source_line: 2 }], name: "ProgramStart", type: "ProgramStart" }
]`},
		// More complicated cases
		{`require 'ripper'
	Ripper.instruction("
	class Bar
		def self.foo
			10
		end
	end
	class Foo < Bar; end
	class FooBar < Foo; end
	FooBar.foo
  ").to_s`, `
[
  { 
    arg_types: { names: [], types: [] }, 
    instructions: [
      { action: "putobject", line: 0, params: ["10"], source_line: 4 }, 
      { action: "leave", line: 1, params: [], source_line: 3 }
    ], 
    name: "foo", 
    type: "Def" 
  }, 
  { 
    arg_types: { names: [], types: [] }, 
    instructions: [
      { action: "putobject", line: 0, params: ["10"], source_line: 4 }, 
      { action: "leave", line: 1, params: [], source_line: 3 }
    ], 
    name: "foo", 
    type: "Def" 
  }, 
  { 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 3 }, 
      { action: "putstring", line: 1, params: ["foo"], source_line: 3 }, 
      { action: "def_singleton_method", line: 2, params: ["0"], source_line: 3 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "Bar", 
    type: "DefClass" 
  }, 
  { 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 3 }, 
      { action: "putstring", line: 1, params: ["foo"], source_line: 3 }, 
      { action: "def_singleton_method", line: 2, params: ["0"], source_line: 3 }, 
      { action: "leave", line: 3, params: [], source_line: 2 }
    ], 
    name: "Bar", 
    type: "DefClass" 
  }, 
  { 
    instructions: [{ action: "leave", line: 0, params: [], source_line: 7 }], name: "Foo", type: "DefClass" 
  }, 
  { 
    instructions: [{ action: "leave", line: 0, params: [], source_line: 7 }], name: "Foo", type: "DefClass" 
  }, 
  { 
    instructions: [{ action: "leave", line: 0, params: [], source_line: 8 }], name: "FooBar", type: "DefClass" 
  }, 
  { 
    instructions: [{ action: "leave", line: 0, params: [], source_line: 8 }], name: "FooBar", type: "DefClass" 
  }, 
  { 
    arg_set: { names: [], types: [] }, 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 2 }, 
      { action: "def_class", line: 1, params: ["class:Bar"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "putself", line: 3, params: [], source_line: 7 }, 
      { action: "getconstant", line: 4, params: ["Bar", "false"], source_line: 7 }, 
      { action: "def_class", line: 5, params: ["class:Foo", "Bar"], source_line: 7 }, 
      { action: "pop", line: 6, params: [], source_line: 7 }, 
      { action: "pop", line: 7, params: [], source_line: 7 }, 
      { action: "putself", line: 8, params: [], source_line: 8 }, 
      { action: "getconstant", line: 9, params: ["Foo", "false"], source_line: 8 }, 
      { action: "def_class", line: 10, params: ["class:FooBar", "Foo"], source_line: 8 }, 
      { action: "pop", line: 11, params: [], source_line: 8 }, { action: "pop", line: 12, params: [], source_line: 8 }, 
      { action: "getconstant", line: 13, params: ["FooBar", "false"], source_line: 9 }, 
      { action: "send", line: 14, params: ["foo", "0", "", "&{[] []}"], source_line: 9 }, 
      { action: "pop", line: 15, params: [], source_line: 9 }, 
      { action: "leave", line: 16, params: [], source_line: 9 }
    ], 
    name: "ProgramStart", 
    type: "ProgramStart" 
  }, 
  { 
    arg_set: { names: [], types: [] }, 
    instructions: [
      { action: "putself", line: 0, params: [], source_line: 2 }, 
      { action: "def_class", line: 1, params: ["class:Bar"], source_line: 2 }, 
      { action: "pop", line: 2, params: [], source_line: 2 }, 
      { action: "putself", line: 3, params: [], source_line: 7 }, 
      { action: "getconstant", line: 4, params: ["Bar", "false"], source_line: 7 }, 
      { action: "def_class", line: 5, params: ["class:Foo", "Bar"], source_line: 7 }, 
      { action: "pop", line: 6, params: [], source_line: 7 }, 
      { action: "pop", line: 7, params: [], source_line: 7 }, 
      { action: "putself", line: 8, params: [], source_line: 8 }, 
      { action: "getconstant", line: 9, params: ["Foo", "false"], source_line: 8 }, 
      { action: "def_class", line: 10, params: ["class:FooBar", "Foo"], source_line: 8 }, 
      { action: "pop", line: 11, params: [], source_line: 8 }, 
      { action: "pop", line: 12, params: [], source_line: 8 }, 
      { action: "getconstant", line: 13, params: ["FooBar", "false"], source_line: 9 }, 
      { action: "send", line: 14, params: ["foo", "0", "", "&{[] []}"], source_line: 9 }, 
      { action: "pop", line: 15, params: [], source_line: 9 }, 
	  { action: "leave", line: 16, params: [], source_line: 9 }
    ],
	name: "ProgramStart", 
	type: "ProgramStart" 
}
]`},
		{`require 'ripper'
Ripper.instruction("
	def foo(x)
	  yield(x + 10)
	end
	def bar(y)
	  foo(y) do |f|
		yield(f)
	  end
	end
	def baz(z)
	  bar(z + 100) do |b|
		yield(b)
	  end
	end
	a = 0
	baz(100) do |b|
	  a = b
	end
	a

	class Foo
	  def bar
		100
	  end
	end
	module Baz
	  class Bar
		def bar
		  Foo.new.bar
		end
	  end
	end
	Baz::Bar.new.bar + a
").to_s`, `
[{
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: ["x"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 3
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 3
  }, {
    action: "putobject",
    line: 2,
    params: ["10"],
    source_line: 3
  }, {
    action: "send",
    line: 3,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "invokeblock",
    line: 4,
    params: ["1"],
    source_line: 3
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 2
  }],
  name: "foo",
  type: "Def"
}, {
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: ["x"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 3
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 3
  }, {
    action: "putobject",
    line: 2,
    params: ["10"],
    source_line: 3
  }, {
    action: "send",
    line: 3,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "invokeblock",
    line: 4,
    params: ["1"],
    source_line: 3
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 2
  }],
  name: "foo",
  type: "Def"
}, {
  arg_types: {
    names: ["f"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 7
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 7
  }, {
    action: "invokeblock",
    line: 2,
    params: ["1"],
    source_line: 7
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 6
  }],
  name: "0",
  type: "Block"
}, {
  arg_types: {
    names: ["f"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 7
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 7
  }, {
    action: "invokeblock",
    line: 2,
    params: ["1"],
    source_line: 7
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 6
  }],
  name: "0",
  type: "Block"
}, {
  arg_set: {
    names: ["y"],
    types: [0]
  },
  arg_types: {
    names: ["y"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 6
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 6
  }, {
    action: "send",
    line: 2,
    params: ["foo", "1", "block:0", "&{[y][0]}"],
    source_line: 6
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 5
  }],
  name: "bar",
  type: "Def"
}, {
  arg_set: {
    names: ["y"],
    types: [0]
  },
  arg_types: {
    names: ["y"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 6
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 6
  }, {
    action: "send",
    line: 2,
    params: ["foo", "1", "block:0", "&{[y][0]}"],
    source_line: 6
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 5
  }],
  name: "bar",
  type: "Def"
}, {
  arg_types: {
    names: ["b"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 12
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 12
  }, {
    action: "invokeblock",
    line: 2,
    params: ["1"],
    source_line: 12
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 11
  }],
  name: "1",
  type: "Block"
}, {
  arg_types: {
    names: ["b"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 12
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 12
  }, {
    action: "invokeblock",
    line: 2,
    params: ["1"],
    source_line: 12
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 11
  }],
  name: "1",
  type: "Block"
}, {
  arg_set: {
    names: [""],
    types: [0]
  },
  arg_types: {
    names: ["z"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 11
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 11
  }, {
    action: "putobject",
    line: 2,
    params: ["100"],
    source_line: 11
  }, {
    action: "send",
    line: 3,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 11
  }, {
    action: "send",
    line: 4,
    params: ["bar", "1", "block:1", "&{[][0]}"],
    source_line: 11
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 10
  }],
  name: "baz",
  type: "Def"
}, {
  arg_set: {
    names: [""],
    types: [0]
  },
  arg_types: {
    names: ["z"],
    types: [0]
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 11
  }, {
    action: "getlocal",
    line: 1,
    params: ["0", "0"],
    source_line: 11
  }, {
    action: "putobject",
    line: 2,
    params: ["100"],
    source_line: 11
  }, {
    action: "send",
    line: 3,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 11
  }, {
    action: "send",
    line: 4,
    params: ["bar", "1", "block:1", "&{[][0]}"],
    source_line: 11
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 10
  }],
  name: "baz",
  type: "Def"
}, {
  arg_types: {
    names: ["b"],
    types: [0]
  },
  instructions: [{
    action: "getlocal",
    line: 0,
    params: ["0", "0"],
    source_line: 17
  }, {
    action: "setlocal",
    line: 1,
    params: ["1", "0"],
    source_line: 17
  }, {
    action: "leave",
    line: 2,
    params: [],
    source_line: 16
  }],
  name: "2",
  type: "Block"
}, {
  arg_types: {
    names: ["b"],
    types: [0]
  },
  instructions: [{
    action: "getlocal",
    line: 0,
    params: ["0", "0"],
    source_line: 17
  }, {
    action: "setlocal",
    line: 1,
    params: ["1", "0"],
    source_line: 17
  }, {
    action: "leave",
    line: 2,
    params: [],
    source_line: 16
  }],
  name: "2",
  type: "Block"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["100"],
    source_line: 23
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 22
  }],
  name: "bar",
  type: "Def"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["100"],
    source_line: 23
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 22
  }],
  name: "bar",
  type: "Def"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 22
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 22
  }, {
    action: "def_method",
    line: 2,
    params: ["0"],
    source_line: 22
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 21
  }],
  name: "Foo",
  type: "DefClass"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 22
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 22
  }, {
    action: "def_method",
    line: 2,
    params: ["0"],
    source_line: 22
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 21
  }],
  name: "Foo",
  type: "DefClass"
}, {
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "getconstant",
    line: 0,
    params: ["Foo", "false"],
    source_line: 29
  }, {
    action: "send",
    line: 1,
    params: ["new", "0", "", "&{[][]}"],
    source_line: 29
  }, {
    action: "send",
    line: 2,
    params: ["bar", "0", "", "&{[][]}"],
    source_line: 29
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 28
  }],
  name: "bar",
  type: "Def"
}, {
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "getconstant",
    line: 0,
    params: ["Foo", "false"],
    source_line: 29
  }, {
    action: "send",
    line: 1,
    params: ["new", "0", "", "&{[][]}"],
    source_line: 29
  }, {
    action: "send",
    line: 2,
    params: ["bar", "0", "", "&{[][]}"],
    source_line: 29
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 28
  }],
  name: "bar",
  type: "Def"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 28
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 28
  }, {
    action: "def_method",
    line: 2,
    params: ["0"],
    source_line: 28
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 27
  }],
  name: "Bar",
  type: "DefClass"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 28
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 28
  }, {
    action: "def_method",
    line: 2,
    params: ["0"],
    source_line: 28
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 27
  }],
  name: "Bar",
  type: "DefClass"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 27
  }, {
    action: "def_class",
    line: 1,
    params: ["class:Bar"],
    source_line: 27
  }, {
    action: "pop",
    line: 2,
    params: [],
    source_line: 27
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 26
  }],
  name: "Baz",
  type: "DefClass"
}, {
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 27
  }, {
    action: "def_class",
    line: 1,
    params: ["class:Bar"],
    source_line: 27
  }, {
    action: "pop",
    line: 2,
    params: [],
    source_line: 27
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 26
  }],
  name: "Baz",
  type: "DefClass"
}, {
  arg_set: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 2
  }, {
    action: "putstring",
    line: 1,
    params: ["foo"],
    source_line: 2
  }, {
    action: "def_method",
    line: 2,
    params: ["1"],
    source_line: 2
  }, {
    action: "putself",
    line: 3,
    params: [],
    source_line: 5
  }, {
    action: "putstring",
    line: 4,
    params: ["bar"],
    source_line: 5
  }, {
    action: "def_method",
    line: 5,
    params: ["1"],
    source_line: 5
  }, {
    action: "putself",
    line: 6,
    params: [],
    source_line: 10
  }, {
    action: "putstring",
    line: 7,
    params: ["baz"],
    source_line: 10
  }, {
    action: "def_method",
    line: 8,
    params: ["1"],
    source_line: 10
  }, {
    action: "putobject",
    line: 9,
    params: ["0"],
    source_line: 15
  }, {
    action: "setlocal",
    line: 10,
    params: ["0", "0"],
    source_line: 15
  }, {
    action: "pop",
    line: 11,
    params: [],
    source_line: 15
  }, {
    action: "putself",
    line: 12,
    params: [],
    source_line: 16
  }, {
    action: "putobject",
    line: 13,
    params: ["100"],
    source_line: 16
  }, {
    action: "send",
    line: 14,
    params: ["baz", "1", "block:2", "&{[][0]}"],
    source_line: 16
  }, {
    action: "pop",
    line: 15,
    params: [],
    source_line: 16
  }, {
    action: "getlocal",
    line: 16,
    params: ["0", "0"],
    source_line: 19
  }, {
    action: "pop",
    line: 17,
    params: [],
    source_line: 19
  }, {
    action: "putself",
    line: 18,
    params: [],
    source_line: 21
  }, {
    action: "def_class",
    line: 19,
    params: ["class:Foo"],
    source_line: 21
  }, {
    action: "pop",
    line: 20,
    params: [],
    source_line: 21
  }, {
    action: "putself",
    line: 21,
    params: [],
    source_line: 26
  }, {
    action: "def_class",
    line: 22,
    params: ["module:Baz"],
    source_line: 26
  }, {
    action: "pop",
    line: 23,
    params: [],
    source_line: 26
  }, {
    action: "getconstant",
    line: 24,
    params: ["Baz", "true"],
    source_line: 33
  }, {
    action: "getconstant",
    line: 25,
    params: ["Bar", "false"],
    source_line: 33
  }, {
    action: "send",
    line: 26,
    params: ["new", "0", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "send",
    line: 27,
    params: ["bar", "0", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "getlocal",
    line: 28,
    params: ["0", "0"],
    source_line: 33
  }, {
    action: "send",
    line: 29,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "pop",
    line: 30,
    params: [],
    source_line: 33
  }, {
    action: "leave",
    line: 31,
    params: [],
    source_line: 33
  }],
  name: "ProgramStart",
  type: "ProgramStart"
}, {
  arg_set: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 2
  }, {
    action: "putstring",
    line: 1,
    params: ["foo"],
    source_line: 2
  }, {
    action: "def_method",
    line: 2,
    params: ["1"],
    source_line: 2
  }, {
    action: "putself",
    line: 3,
    params: [],
    source_line: 5
  }, {
    action: "putstring",
    line: 4,
    params: ["bar"],
    source_line: 5
  }, {
    action: "def_method",
    line: 5,
    params: ["1"],
    source_line: 5
  }, {
    action: "putself",
    line: 6,
    params: [],
    source_line: 10
  }, {
    action: "putstring",
    line: 7,
    params: ["baz"],
    source_line: 10
  }, {
    action: "def_method",
    line: 8,
    params: ["1"],
    source_line: 10
  }, {
    action: "putobject",
    line: 9,
    params: ["0"],
    source_line: 15
  }, {
    action: "setlocal",
    line: 10,
    params: ["0", "0"],
    source_line: 15
  }, {
    action: "pop",
    line: 11,
    params: [],
    source_line: 15
  }, {
    action: "putself",
    line: 12,
    params: [],
    source_line: 16
  }, {
    action: "putobject",
    line: 13,
    params: ["100"],
    source_line: 16
  }, {
    action: "send",
    line: 14,
    params: ["baz", "1", "block:2", "&{[][0]}"],
    source_line: 16
  }, {
    action: "pop",
    line: 15,
    params: [],
    source_line: 16
  }, {
    action: "getlocal",
    line: 16,
    params: ["0", "0"],
    source_line: 19
  }, {
    action: "pop",
    line: 17,
    params: [],
    source_line: 19
  }, {
    action: "putself",
    line: 18,
    params: [],
    source_line: 21
  }, {
    action: "def_class",
    line: 19,
    params: ["class:Foo"],
    source_line: 21
  }, {
    action: "pop",
    line: 20,
    params: [],
    source_line: 21
  }, {
    action: "putself",
    line: 21,
    params: [],
    source_line: 26
  }, {
    action: "def_class",
    line: 22,
    params: ["module:Baz"],
    source_line: 26
  }, {
    action: "pop",
    line: 23,
    params: [],
    source_line: 26
  }, {
    action: "getconstant",
    line: 24,
    params: ["Baz", "true"],
    source_line: 33
  }, {
    action: "getconstant",
    line: 25,
    params: ["Bar", "false"],
    source_line: 33
  }, {
    action: "send",
    line: 26,
    params: ["new", "0", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "send",
    line: 27,
    params: ["bar", "0", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "getlocal",
    line: 28,
    params: ["0", "0"],
    source_line: 33
  }, {
    action: "send",
    line: 29,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 33
  }, {
    action: "pop",
    line: 30,
    params: [],
    source_line: 33
  }, {
    action: "leave",
    line: 31,
    params: [],
    source_line: 33
  }],
  name: "ProgramStart",
  type: "ProgramStart"
}]`},

		{`require 'ripper'
Ripper.instruction("
	def bar(block)
	block.call + get_block.call
	end

	def foo
		bar(get_block) do
 		20
		end
	end

	foo do
		10
	end
").to_s`, `
[{
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: ["block"],
    types: [0]
  },
  instructions: [{
    action: "getlocal",
    line: 0,
    params: ["0", "0"],
    source_line: 3
  }, {
    action: "send",
    line: 1,
    params: ["call", "0", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "getblock",
    line: 2,
    params: [],
    source_line: 3
  }, {
    action: "send",
    line: 3,
    params: ["call", "0", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "send",
    line: 4,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 2
  }],
  name: "bar",
  type: "Def"
}, {
  arg_set: {
    names: [],
    types: []
  },
  arg_types: {
    names: ["block"],
    types: [0]
  },
  instructions: [{
    action: "getlocal",
    line: 0,
    params: ["0", "0"],
    source_line: 3
  }, {
    action: "send",
    line: 1,
    params: ["call", "0", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "getblock",
    line: 2,
    params: [],
    source_line: 3
  }, {
    action: "send",
    line: 3,
    params: ["call", "0", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "send",
    line: 4,
    params: ["+", "1", "", "&{[][]}"],
    source_line: 3
  }, {
    action: "leave",
    line: 5,
    params: [],
    source_line: 2
  }],
  name: "bar",
  type: "Def"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["20"],
    source_line: 8
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 7
  }],
  name: "0",
  type: "Block"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["20"],
    source_line: 8
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 7
  }],
  name: "0",
  type: "Block"
}, {
  arg_set: {
    names: [""],
    types: [0]
  },
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 7
  }, {
    action: "getblock",
    line: 1,
    params: [],
    source_line: 7
  }, {
    action: "send",
    line: 2,
    params: ["bar", "1", "block:0", "&{[][0]}"],
    source_line: 7
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 6
  }],
  name: "foo",
  type: "Def"
}, {
  arg_set: {
    names: [""],
    types: [0]
  },
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 7
  }, {
    action: "getblock",
    line: 1,
    params: [],
    source_line: 7
  }, {
    action: "send",
    line: 2,
    params: ["bar", "1", "block:0", "&{[][0]}"],
    source_line: 7
  }, {
    action: "leave",
    line: 3,
    params: [],
    source_line: 6
  }],
  name: "foo",
  type: "Def"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["10"],
    source_line: 13
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 12
  }],
  name: "1",
  type: "Block"
}, {
  arg_types: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putobject",
    line: 0,
    params: ["10"],
    source_line: 13
  }, {
    action: "leave",
    line: 1,
    params: [],
    source_line: 12
  }],
  name: "1",
  type: "Block"
}, {
  arg_set: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 2
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 2
  }, {
    action: "def_method",
    line: 2,
    params: ["1"],
    source_line: 2
  }, {
    action: "putself",
    line: 3,
    params: [],
    source_line: 6
  }, {
    action: "putstring",
    line: 4,
    params: ["foo"],
    source_line: 6
  }, {
    action: "def_method",
    line: 5,
    params: ["0"],
    source_line: 6
  }, {
    action: "putself",
    line: 6,
    params: [],
    source_line: 12
  }, {
    action: "send",
    line: 7,
    params: ["foo", "0", "block:1", "&{[][]}"],
    source_line: 12
  }, {
    action: "pop",
    line: 8,
    params: [],
    source_line: 12
  }, {
    action: "leave",
    line: 9,
    params: [],
    source_line: 12
  }],
  name: "ProgramStart",
  type: "ProgramStart"
}, {
  arg_set: {
    names: [],
    types: []
  },
  instructions: [{
    action: "putself",
    line: 0,
    params: [],
    source_line: 2
  }, {
    action: "putstring",
    line: 1,
    params: ["bar"],
    source_line: 2
  }, {
    action: "def_method",
    line: 2,
    params: ["1"],
    source_line: 2
  }, {
    action: "putself",
    line: 3,
    params: [],
    source_line: 6
  }, {
    action: "putstring",
    line: 4,
    params: ["foo"],
    source_line: 6
  }, {
    action: "def_method",
    line: 5,
    params: ["0"],
    source_line: 6
  }, {
    action: "putself",
    line: 6,
    params: [],
    source_line: 12
  }, {
    action: "send",
    line: 7,
    params: ["foo", "0", "block:1", "&{[][]}"],
    source_line: 12
  }, {
    action: "pop",
    line: 8,
    params: [],
    source_line: 12
  }, {
    action: "leave",
    line: 9,
    params: [],
    source_line: 12
  }],
  name: "ProgramStart",
  type: "ProgramStart"
}]`},
	}
	for i, tt := range tests {
		evaluated := vm.ExecAndReturn(t, tt.input)

		result := removeRedundantSpaces(evaluated.Value().(string))
		expected := removeRedundantSpaces(tt.expected)

		if result != expected {
			t.Fatalf("Expect:\n%s. Got:\n%s\nAt case %d", expected, result, i)
		}
	}
}

func TestRipperInstructionFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require 'ripper'; Ripper.instruction`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`require 'ripper'; Ripper.instruction("a = 1", "a = 2")`, "ArgumentError: Expect 1 argument. got=2", 1},
		{`require 'ripper'; Ripper.instruction(Array)`, "TypeError: Expect argument to be String. got: Class", 1},
		{`require 'ripper'; Ripper.instruction("10.times do")`, "InternalError: invalid code: 10.times do", 1},
	}

	for i, tt := range testsFail {
		evaluated := vm.ExecAndReturn(t, tt.input)
		checkErrorMsg(t, i, evaluated, tt.expected)
	}
}

// Error test helper methods

type Error struct {
	*vm.BaseObj
	message      string
	stackTraces  []string
	storedTraces bool
	Type         string
}

func checkErrorMsg(t *testing.T, index int, evaluated Object, expectedErrMsg string) {
	t.Helper()
	err, ok := evaluated.(*vm.Error)
	if !ok {
		t.Fatalf("At test case %d: Expect Error. got=%T (%+v)", index, evaluated, evaluated)
	}

	message := strings.Split(err.Message(), "\n")
	if message[0] != expectedErrMsg {
		t.Fatalf("At test case %d: Expect error message to be:\n  %s. got: \n%s", index, expectedErrMsg, err.Message())
	}
}

func removeRedundantSpaces(str string) string {
	space := regexp.MustCompile(`\s`)
	newline := regexp.MustCompile(`\n`)
	str = space.ReplaceAllString(str, "")
	return newline.ReplaceAllString(str, "")
}
