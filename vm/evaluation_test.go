package vm

import (
	"testing"
)

func TestComplexEvaluation(t *testing.T) {
	input := `
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
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	testIntegerObject(t, 0, evaluated, 310)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestComment(t *testing.T) {
	input := `
	# Comment
	class Foo
		# Comment
		def one # Comment
			# Comment
			1 # Comment
			# Comment
		end
		# Comment

		def bar(x) # Comment
		  123
		end  # Comment
	end
	# Comment
	Foo.new.one #=> Comment
	# Comment`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	testIntegerObject(t, 0, evaluated, 1)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			class Foo
			  def bar
			    10
			  end

			  def baz(x)
			    x + 100
			  end

			  def foo
			    x = baz(bar)
			    x
			  end
			end

			Foo.new.foo
			`, 110,
		},
		{
			`
			class Foo
			  def baz
			    @foo = 100
			  end

			  def foo
			    @foo
			  end

			  def bar
			    baz
			    foo
			  end
			end

			Foo.new.bar
			`, 100,
		},
		{
			`
			class Foo
			  def baz
			    @foo = 100
			  end

			  def foo
			    @foo
			  end

			  def baz2
			    @foo = @foo + 100
			  end

			  def bar
			    baz
			    baz2
			    foo
			  end
			end

			Foo.new.bar
			`, 200,
		},
		{
			`
			class Foo
			  def set_x(x) # Set x
			    @x = x
			  end

			  def foo # Set x and plus a
			    set_x(10)
			    a = 10
			    @x + a
			  end
			end

			f = Foo.new
			f.foo
			`,
			20,
		},
		{
			`
			class Foo
			  def bar=(x)
			    @bar = x
			  end

			  def bar
			    @bar
			  end
			end

			f = Foo.new
			f.bar = 10
			f.bar
			`,
			10,
		},
		{
			`
			class Foo
			  def set_x(x)
			    @x = x
			  end

			  def foo
			    set_x(10 + 10 * 100)
			    a = 10
			    @x + a
			  end
			end

			f = Foo.new
			f.foo
			`,
			1020,
		},
		{
			`class Foo
				def bar
					10
				end

				def foo
					bar = 100
					10 + bar
				end
			end

			f = Foo.new
			f.foo
			`,
			110,
		},
		{
			`class Foo
				def bar
					10
				end

				def foo
					a = 10
					bar + a
				end
			end

			Foo.new.foo
			`,
			20,
		},
		{
			`class Foo
				def self.bar
					10
				end

				def self.foo
					a = 10
					bar + a
				end
			end

			Foo.foo
			`,
			20,
		},
		{
			`class Foo
				def bar
					100
				end

				def self.bar
					10
				end

				def foo
					a = 10
					bar + a
				end
			end

			Foo.new.foo
			`,
			110,
		},
		{`
		class Foo
		  def bar
		  end
		end

		Foo.new.bar
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMethodCallWithDefaultArgValue(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		def foo(x = 10)
		  x
		end

		foo
		`, 10},
		{`
		def foo(x = 10, y)
		  x + y
		end

		foo(100, 10)
		`, 110},
		{`
		def foo(x, y = 10)
		  x + y
		end

		foo(100)
		`, 110},
		{`
		def foo(x = 100, y = 10)
		  x + y
		end

		foo
		`, 110},
		{`
		def foo(x = 100, y = 10)
		  x + y
		end

		foo(200)
		`, 210},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMethodCallWithBlockArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
				class Foo
				  def bar
				    yield(1, 3, 5)
				  end
				end

				Foo.new.bar do |first, second, third|
				  first + second * third
				end

				`, 16},
		{`
				class Foo
				  def bar
				    yield
				  end
				end

				Foo.new.bar do
				  3
				end

				`, 3},
		{`
				class Bar
				  def foo
				    yield(10)
				  end
				end

				class Foo
				  def bar
				    yield
				  end
				end

				Bar.new.foo do |num|
				  Foo.new.bar do
				    3 * num
				  end
				end

				`, 30},
		{`
				class Foo
				  def bar
				    0
				  end
				end

				Foo.new.bar do
				  3
				end

				`, 0},
		{`
				class Foo
				  def bar
				    yield
				  end
				end

				i = 10
				Foo.new.bar do
				  i = 3 + i
				end
				i

				`, 13},
		{`
		class Car
		  def initialize
		    yield(self)
		  end

		  def doors=(ds)
		    @doors = ds
		  end

		  def doors
		    @doors
		  end
		end

		car = Car.new do |c|
		  c.doors = 4
		end

		car.doors
				`,
			4},
		{`
		class Foo
		  def bar(x)
		    yield(x)
		  end
		end

		f = Foo.new
		x = 100
		y = 10

		f.bar(10) do |x|
                  y = x + y
		end

		y
		`, 20},
		{`
		class Foo
		  def self.bar
		    yield(10)
		  end
		end

		i = 1
		r = 0
		Foo.bar do |ten|
		  r = ten + i
		end

		r
		`,
			11},
		{
			`
			def foo
			  yield(10)
			end

			foo do |ten|
			  ten + 20
			end
			`, 30},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMethodCallWithNestedBlock(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		class Foo
		  def bar
		    yield
		  end
		end

		a = 100
		i = 10
		b = 1000

		f = Foo.new

		f.bar do
		  i = 3 * a
		  f.bar do
		    i = 3 + i
		  end
		end
		i

		`, 303},
		{`
		class Foo
		  def bar
		    yield
		  end
		end

		i = 10
		a = 100
		b = 1000

		f = Foo.new

		f.bar do
		  a = 20
		  f.bar do
		    b = (3 + i) * a
		  end
		end
		b

		`, 260},
		{
			`
			def foo(x)
			  yield(x + 10)
			end

			def bar(y)
			  foo(y) do |f|
			    yield(f)
			  end
			end

			a = 0
			bar(100) do |b|
			  a = b
			end

			a
			`,
			110},
		{
			`
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
			`,
			210},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMethodCallWithoutParens(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		def foo(x)
		  x
		end

		a = foo 10
		a
		`, 10},
		{
			`
			class Foo
			  def bar
			    10
			  end

			  def baz(x)
			    x + 100
			  end

			  def foo
			    x = baz bar
			    x
			  end
			end

			Foo.new.foo
			`, 110,
		},
		{
			`
			class Foo
			  def set_x(x0)
			    @x = x0
			  end

			  def foo
			    set_x 10
			    a = 10
			    @x + a
			  end
			end

			f = Foo.new
			f.foo
			`,
			20,
		},
		{
			`
			class Foo
			  def set_x(x1, x2)
			    @x = x1
			    @y = x2
			  end

			  def foo
			    set_x 10,11
			    a = 10
			    @x + a +@y
			  end
			end

			f = Foo.new
			f.foo
			`,
			31,
		},
		{
			`
			class Foo
			  def set_x(x1, x2)
			    @x1 = x1
			    @x2 = x2
			  end

			  def set_y(y1, y2, y3)
			    @y3 = y3
			    @y1 = y1
			  end

			  def foo
			    set_x 15,17
			    set_y 3,4,5
			    set_x 10,11
			    @x1 + @x2 + @y3
			  end
			end

			f = Foo.new
			f.foo
			`,
			26,
		},
		{
			`
			class Foo
			  attr_reader :x, :y

			  def set_x(x1)
			    @x1 = x1
			  end

			  def set_y(y1, y2, y3)
			    @y = y1 + y2 + y3
			  end

			  def foo
			    set_x 10
			    set_y 1,2,3
			    set_y 3,4,5
			    @x1 + @y
			  end
			end

			f = Foo.new
			f.foo
			`,
			22,
		},
		{
			`

			class Foo
   			  attr_reader :y

			  def set_y(y1, y2, y3)
			    @y = y1 + y2 + y3
			  end

			end

			f = Foo.new
			f.set_y 3,4,5

			f.y
			`,
			12,
		},
		{
			`
			class Foo
			  attr_reader :x, :y

			  def set_x(x1)
			    @x = x1
			  end

			  def set_y(y1, y2, y3)
			    @y = y1 + y2 + y3
			  end
			end

			f = Foo.new
			f.set_x 1
			f.set_y 4,5,6
			f.x + f.y
			`,
			16,
		},

		{
			`
			class Foo
			  attr_reader :x, :y

			  def set_x(x1, x2, x3)
			    @x = x1 + x2 + x3
			  end

			  def set_y(y1, y2, y3)
			    @y = y1 + y2 + y3
			  end
			end

			f = Foo.new
			f.set_x 1,2,3
			f.set_y 4,5,6
			f.x + f.y
			`,
			21,
		},
		{
			`
			class Foo
			  def bar
			    yield 10
			  end
			end

			i = 0
			Foo.new.bar do |ten|
			  i = ten
			end
			i
			`, 10,
		},
		{
			`

		    def bar
			  yield(10)
		    end

			i = 0
			bar do |ten|
			  i = ten
			end
			i
			`, 10,
		},
		{`
		def foo
		  10
		end

		def double(x)
		  x * 2
		end

		double foo
		`, 20},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			class Bar
				def self.foo
					10
				end
			end
			Bar.foo;
			`,
			10,
		},
		{
			`
			class Bar
				def self.foo
					10
				end
			end
			class Foo < Bar; end
			class FooBar < Foo; end
			FooBar.foo
			`,
			10,
		},
		{
			`
			class Foo
				def self.foo
					10
				end
			end

			class Bar < Foo; end
			Bar.foo
			`,
			10,
		},
		{
			`
			class Foo
				def self.foo
					10
				end
			end

			class Bar < Foo
				def self.foo
					100
				end
			end
			Bar.foo
			`,
			100,
		},
		{
			`
			class Bar
				def self.foo
					bar
				end

				def self.bar
					100
				end

				def bar
					1000
				end
			end
			Bar.foo
			`,
			100,
		},
		{
			`
			# Test class method call inside class method.
			class JobPosition
				def initialize(name)
					@name = name
				end

				def self.engineer
					new("Engineer")
				end

				def name
					@name
				end
			end
			job = JobPosition.engineer
			job.name
			`,
			"Engineer",
		},
		{
			`
			class Foo; end
			Foo.new.class.name
			`,
			"Foo",
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestInstanceMethodCall(t *testing.T) {
	input := `

		class Bar
			def set(x)
				@x = x
			end
		end

		class Foo < Bar
			def add(x, y)
				x + y
			end
		end

		class FooBar < Foo
			def get
				@x
			end
		end

		fb = FooBar.new
		fb.set(100)
		fb.add(10, fb.get)
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*Error).Message)
	}

	result, ok := evaluated.(*IntegerObject)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 110 {
		t.Errorf("expect result to be 110. got=%d", result.Value)
	}

	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestPostfixMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		a = 1
		a++
		a
		`, 2},
		{`
		a = 10
		a--
		a
		`,
			9},
		{`
		a = 0
		a--
		a
		`,
			-1},
		{`
		a = -5
		a++
		a
		`,
			-4},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBangPrefixMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!5", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMinusPrefixMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"-5", -5},
		{"-10", -10},
		{"-(-10)", 10},
		{"-(-5)", 5},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestSelfExpressionEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`self.class.name`, "Object"},
		{
			`
			class Bar
				def whoami
					"Instance of " + self.class.name
				end
			end

			Bar.new.whoami
		`, "Instance of Bar"},
		{
			`
			class Foo
				Self = self

				def get_self
					Self
				end
			end

			Foo.new.get_self.name
			`,
			"Foo"},
		{
			`
			class Foo
				def class
					Foo
				end
			end

			Foo.new.class.name
			`,
			"Foo"},
		{
			`
			class Foo
				def class_name
					self.class.name
				end
			end

			Foo.new.class_name
			`,
			"Foo"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestInstanceVariableEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		class Foo
			def set(x)
				@x = x;
			end

			def get
				@x
			end

			def double_get
				self.get() * 2;
			end
		end

		class Bar
			def set(x)
				@x = x;
			end

			def get
				@x
			end
		end

		f1 = Foo.new
		f1.set(10)

		f2 = Foo.new
		f2.set(20)

		b = Bar.new
		b.set(10)

		f2.double_get() + f1.get() + b.get()
	`, 60},
		{`
		class Foo
		  attr_reader("bar")
		end

		Foo.new.bar
		`, nil},
		{`
		class Foo
		  def bar
		    @x
		  end
		end

		Foo.new.bar
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestAssignmentEvaluation(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5; b = a; b;", 5},
		{"a = 5; b = a; c = a + b + 5; c;", 15},
		{"a = 5; b = 10; c = if a > b; 100 else 50 end; c", 50},
		{"Bar = 100; Bar", 100},
		{`
		a = 100
		b = a
		b = 1000
		a
		`, 100},
		{`
		a = 100
		b = a
		a = 1000
		b
		`, 100},
		{`
		i = 0

		if a = 10
		  i = 100
		end

		i + a
		`, 110},
		{`
		i = 0

		if @a = 10
		  i = 100
		end

		i + @a
		`, 110},
		{`a = b = 10; a`, 10},
		{`a = b = c = 10; a`, 10},
		{`
		i = 100
		a = b = i + 10
		a + b
		`, 220},
		{`
		def foo(x)
		  x
		end

		foo(a = b = c = d = 10)
		`, 10},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		testIntegerObject(t, i, evaluated, tt.expectedValue)
		v.checkCFP(t, 0, 0)
		v.checkSP(t, i, 1)
	}
}

func TestAssignmentByOperationEvaluation(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"a = 5; a += 2; a;", 7},
		{"a = 5; a -= 10; a;", -5},
		{"a = 5; a += 2 * 3 + 5; a;", 16},
		{"a = 5; a -= 2 * 3 + 5; a;", -6},
		{"a = false; a ||= true; a;", true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expectedValue)
		v.checkCFP(t, 0, 0)
		v.checkSP(t, i, 1)
	}
}

func TestIfExpressionEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			if 10 > 5
				100
			else
				-10
			end
			`,
			100,
		},
		{
			`
			if 5 != 5
				false
			else
				true
			end
			`,
			true,
		},
		{`
		if true
		   10
		end`,
			10,
		},
		{"if false; 10 end", nil},
		{"if 1; 10; end", 10},
		{"if 1 < 2; 10 end", 10},
		{"if 1 > 2; 10 end", nil},
		{"if 1 > 2; 10 else 20 end", 20},
		{"if 1 < 2; 10 else 20 end", 10},
		{"if nil; 10 else 20 end", 20},
		{`
		if false
		  x = 1
		end # This pushes nil

		x # This pushes nil too
		`, nil},
		{`
		def foo
		  if false
		    x = 1
	      end # This shouldn't push nil
		end

		foo # This should push nil
		`, nil},
		{`
		def foo
		  x = 0
		  if true
		    x = 1
	      end # This shouldn't push nil
	      x
		end

		foo # This should push nil
		`, 1},
		{`
			a = 10
			b = 5
			if a > b
			  puts(123)
			  c = 10
			else
			  c = 5
			end

			c + 1
		`, 11},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassInheritance(t *testing.T) {
	input := `
		class Bar
		end

		class Foo < Bar
		  def self.add
		    10
		  end
		end

		Foo.superclass.name
	`
	v := initTestVM()
	evaluated := v.testEval(t, input)

	testStringObject(t, 0, evaluated, "Bar")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestMultiVarAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		a, b = [1, 2]
		a
		`,
			1},
		{`
		a, b = [1, 2]
		b
		`,
			2},

		{`
		a, b, c = [1, 2, 3]
		c
		`,
			3},
		{`
		a, b, c = [1]
		b
		`,
			nil},
		{`
		a, b, c = [1]
		c
		`,
			nil},
		{`
		arr = [1, 2, 3]
		a, b, c = arr
		b
		`, 2},
		{`
		arr = [1, 2, 3]
		a, b, c = arr
		c
		`, 3},
		{`
		arr = [1]
		a, b, c = arr
		a
		`, 1},
		{`
		arr = [1]
		a, b, c = arr
		b
		`, nil},
		{`
		arr = [1]
		a, b, c = arr
		c
		`, nil},
		{`
		arr = [1]
		a, b, c, d = arr
		d
		`, nil},
		{`
		arr = [1, 2, 3]
		@a, @b, c = arr
		@a
		`, 1},
		{`
		arr = [1, 2, 3]
		@a, @b, c = arr
		@b
		`, 2},
		{`
		arr = [1, 2, 3]
		@a, @b, c = arr
		c
		`, 3},
		{`
		class Foo
		  attr_reader :a, :b, :c

		  def bar(arr)
		    @a, @b, c = arr
		    @c = @a + @b + c
		  end
		end

		f = Foo.new

		f.bar([10, 100, 200])
		f.a
		`, 10},
		{`
		class Foo
		  attr_reader :a, :b, :c

		  def bar(arr)
		    @a, @b, c = arr
		    @c = @a + @b + c
		  end
		end

		f = Foo.new

		f.bar([10, 100, 200])
		f.b
		`, 100},
		{`
		class Foo
		  attr_reader :a, :b, :c

		  def bar(arr)
		    @a, @b, c = arr
		    @c = @a + @b + c
		  end
		end

		f = Foo.new

		f.bar([10, 100, 200])
		f.c
		`, 310},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
