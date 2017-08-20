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
	input := `
		class Foo
			def bar(x, y)
				x + y
			end

			def foo(y)
				y
			end

			def baz=(x)
			end
		end

		Foo
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	class := evaluated.(*RClass)

	expectedMethods := []struct {
		name   string
		params []string
	}{
		{name: "foo", params: []string{"y"}},
		{name: "bar", params: []string{"x", "y"}},
		{name: "baz=", params: []string{"x"}},
	}

	for _, expectedMethod := range expectedMethods {
		methodObj, ok := class.Methods.get(expectedMethod.name)
		if !ok {
			t.Errorf("expect class %s to have method %s.", class.Name, expectedMethod.name)
		}

		method := methodObj.(*MethodObject)
		if method.Name != expectedMethod.name {
			t.Errorf("expect method's name to be %s. got=%s", expectedMethod.name, method.Name)
		}

		if method.argc != len(expectedMethod.params) {
			t.Errorf("expect method %s to have %d parameters. got=%d", method.Name, len(expectedMethod.params), method.argc)
		}
	}

	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
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
