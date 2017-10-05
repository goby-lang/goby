package vm

import (
	"testing"
)

func TestReturnStatementEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`
	class Foo
	  def self.bar
	    return 100
	    10
	  end
	end

	Foo.bar
			`,
			100,
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)

	}
}

func TestDefStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		def foo
		  10
		end

		foo
		`, 10},
		{`
		def foo(x)
		  x + 10
		end

		foo(10)
		`, 20},
		{`
		def foo(x = 100, y=10)
		  x + y
		end

		foo
		`, 110},
		{`
		def foo(x = 100, y=10)
		  x + y
		end

		foo(20)
		`, 30},
		{`
		def foo(x, y=10)
		  x + y
		end

		foo(100)
		`, 110},
		{`
		def foo(x=10, y=11, z=12)
		  x + y + z
		end

		foo(10, 20)
		`, 42},
		{`
		class Foo; end

		def Foo.foo
		  10
		end

		Foo.foo
		`, 10},
		{`
		class Foo
		  def ten
		    10
		  end
		end

		f1 = Foo.new
		f2 = Foo.new

		def f1.ten
		  20
		end

		f2.ten + f1.ten
		`, 30},
		{`
		a = 1

		def a.foo
		  10
		end

		a.foo
		`, 10},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDefStatementWithSplatArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{

		{`
		def foo(a, *b)
		  b.each do |i|
		    a = a + i
		  end
		  a
		end

		foo(10, 15, 20)
		`, 45},
		{`
		def foo(a, b, *c)
		  c.each do |i|
		    a = a + i
		  end
		  a + b
		end

		foo(10, 15, 20, 25)
		`, 70},
		{`
		def foo(a, b = 15, *c)
		  c.each do |i|
		    a = a + i
		  end
		  a + b
		end

		foo(10, 20, 25)
		`, 55},
		{`
		def foo(*a)
		  r = 0
		  a.each do |i|
		    r += i
		  end
		  r
		end

		foo(10, 20, 30)
		`, 60},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestDefStatementWithKeywordArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		def foo(a:)
		  a
		end

		foo(a:10)
		`, 10},
		{`
		def foo(a: 20)
		  a
		end

		foo
		`, 20},
		{`
		def foo(a: 20)
		  a
		end

		foo(a: 10)
		`, 10},
		{`
		def foo(a:, b:)
		  a - b
		end

		foo(a:10, b: 20)
		`, -10},
		{`
		def foo(a: 10, b:)
		  a - b
		end

		foo(b: 20)
		`, -10},
		{`
		def foo(a:, b: 20)
		  a - b
		end

		foo(a: 10)
		`, -10},
		{`
		def foo(a:, b:)
		  a - b
		end

		foo(b:10, a: 20)
		`, 10},
		{`
		def foo(foo, a:, b:)
		  a - b + foo
		end

		foo(100, a:10, b: 20)
		`, 90},
		{`
		def foo(foo, a: 10, b:)
		  a - b + foo
		end

		foo(100, b: 20)
		`, 90},
		// Two normal arguments plus two keyword arguments
		{`
		def foo(bar, foo, a:, b:)
		  a - b + foo - bar
		end

		foo(40, 100, a:10, b: 20)
		`, 50},
		{`
		def foo(bar, foo, a:, b:)
		  a - b + foo - bar
		end

		foo(40, 100, b: 20, a: 10)
		`, 50},
		{`
		def foo(bar, foo, a:, b:)
		  a - b + foo - bar
		end

		foo(b: 20, a: 10, 40, 100)
		`, 50},
		//{`
		//def foo(bar, foo = 100, a:, b:)
		//  a - b + foo - bar
		//end
		//
		//foo(b: 20, a: 10, 40)
		//`, 50},

		// Add splat arguments
		{`
		def foo(bar, foo, a:, b:, *args)
		  a - b + foo - bar + *args[1]
		end

		foo(b: 20, a: 10, 40, 100, "foo", 50)
		`, 100},
		{`
		def foo(bar, foo, a:, b:, *args)
		  a - b + foo - bar + *args[1]
		end

		foo(b: 20, a: 10, 40, 100, "foo", 50)
		`, 100},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestModuleStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		module Bar
		  def bar
		    10
		  end
		end

		class Foo
		  include Bar
		end

		Foo.new.bar
		`,
			10},
		{`
		module Bar
		  def bar
		    10
		  end
		end

		include Bar

		bar
		`,
			10},
		{`module Foo
			  def ten
			    10
			  end
			end

			class Baz
			  def ten
			    1
			  end

			  def five
			    5
			  end
			end

			class Bar < Baz
			  include Foo
			end

			b = Bar.new
			b.ten * b.five
`, 50},
		{`
		module Foo
		  def ten
			10
		  end
		  def twenty
			20
		  end
		end

		module Bar
		  def twenty
			"20th"
		  end
		end

		class Baz
		  include(Foo)
		  include(Bar)
		  include(Foo)
		end

		a = Baz.new
		a.twenty
		`, "20th"},
		{`
		module Foo
		  def ten
			10
		  end
		end

		module Bar
		  def twenty
			20
		  end
		end

		class Baz
		  include(Bar)
		  include(Foo)
		end

		class Baz
		  include(Bar)
		  include(Foo)
		end

		b = Baz.new
		b.ten + b.twenty
		`, 30},
		{`
		module Foo
		  def ten
		    10
		  end
		end

		class Bar
		  extend Foo
		end

		Bar.ten
		`, 10},
		{`
		module Foo; end

		class Bar
		  extend Foo
		end

		module Foo
		  def ten
		    10
		  end
		end

		Bar.ten
		`, 10},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			`
			i = 10
			while i < 0 do
			  i = i + 1
			end

			i
			`, 10},
		{
			`
			i = 10

			while i < 0 do
			  10
			  i = i + 1
			end

			i
			`, 10},
		{
			`
			i = 10

			while i < 0 do
			  i = i + 1
			  10
			end

			i
			`, 10},
		{
			`
		i = 10
		while i > 0 do
		  i -= 1
		end
		i
		`, 0},
		{
			`
		a = [1, 2, 3, 4, 5]
		i = 0
		while i < a.length do
			a[i] += 1
			i += 1
		end
		a[4]
		`, 6},
		// These are regression tests for #396
		// Which should prevent parser from peeking the do keyword and consider identifier as method call
		{`
		i = 0
		l = 10
		while i < l do
		  i += 1
		end

		i
		`, 10},
		{`
		f = false
		i = 0
		while f do
		  i += 1
		end

		i
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestNextStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
x = 0
y = 0

while x < 10 do
  x = x + 1
  if x == 5
	next
  end
  y = y + 1
end

x + y
		`, 19},
		{`
x = 0
y = 0
i = 0

while x < 10 do
  x = x + 1
  while y < 5 do
	y = y + 1

	if y == 3
	  next
	end

	i = i + x * y
  end
end

i
		`, 12},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBreakStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
x = 0
y = 0

while x < 10 do
  x = x + 1
  if x == 5
	break
  end
  y = y + 1
end

x + y
		`, 9},
		{`
x = 0
y = 0
i = 0

while x < 10 do
  x = x + 1
  while y < 5 do
	y = y + 1

	if y == 3
	  break
	end

	i = i + x * y
  end
end

a = i * 10
a + 100
		`, 310},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
