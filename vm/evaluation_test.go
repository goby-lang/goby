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

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	testIntegerObject(t, 0, evaluated, 310)
	vm.checkCFP(t, 0, 0)
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

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	testIntegerObject(t, 0, evaluated, 1)
	vm.checkCFP(t, 0, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestMethodCallWithoutParens(t *testing.T) {
	tests := []struct {
		input    string
		expected int
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

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

	vm.checkCFP(t, 0, 0)
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
		{`
		a = 1
		a+=1
		a
		`,
			2},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testIntegerObject(t, i, evaluated, tt.expectedValue)
	}
}

func TestIfExpressionEvaluation(t *testing.T) {
	tests := []struct {
		input      string
		expected   interface{}
		expectedSP int
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
			1,
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
			1,
		},
		{`
		if true
		   10
		end`,
			10,
			1,
		},
		{"if false; 10 end", nil, 1},
		{"if 1; 10; end", 10, 1},
		{"if 1 < 2; 10 end", 10, 1},
		{"if 1 > 2; 10 end", nil, 1},
		{"if 1 > 2; 10 else 20 end", 20, 1},
		{"if 1 < 2; 10 else 20 end", 10, 1},
		{"if nil; 10 else 20 end", 20, 1},
		{`
		if false
		  x = 1
		end # This pushes nil

		x # This pushes nil too
		`, nil, 2},
		{`
		def foo
		  if false
		    x = 1
	      end # This shouldn't push nil
		end

		foo # This should push nil
		`, nil, 1},
		{`
		def foo
		  x = 0
		  if true
		    x = 1
	      end # This shouldn't push nil
	      x
		end

		foo # This should push nil
		`, 1, 1},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, tt.expectedSP)
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
	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	testStringObject(t, 0, evaluated, "Bar")
	vm.checkCFP(t, 0, 0)
}
