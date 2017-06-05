package vm

import "testing"

func TestMonkeyPatchBuiltInClass(t *testing.T) {
	input := `
	class String
	  def buz
	    "buz"
	  end
	end

	"123".buz
	`

	evaluated := testEval(t, input)
	checkExpected(t, evaluated, "buz")
}

func TestRequireRelative(t *testing.T) {
	input := `
	require_relative("../test_fixtures/require_test/foo")

	fifty = Foo.bar(5)

	Foo.baz do |hundred|
	  hundred + fifty + Bar.baz
	end
	`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 160)
}

func TestDefSingletonMethtod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		class Foo
		  def self.bar
		    10
		  end
		end

		Foo.bar
		`, 10},
		{`
		module Foo
		  def self.bar
		    10
		  end
		end

		Foo.bar
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s.\n Input %s", evaluated.(*Error).Message, tt.input)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestAttrReaderAndWriter(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		class Foo
		  attr_reader("bar")

		  def set_bar(bar)
		    @bar = bar
		  end
		end

		f = Foo.new
		f.set_bar(10)
		f.bar

		`, 10},
		{`
		class Foo
		  attr_writer("bar")

		  def bar
		    @bar
		  end
		end

		f = Foo.new
		f.bar = 10
		f.bar

		`, 10},
		{`
		class Foo
		  attr_writer("bar")
		  attr_reader("bar")
		end

		f = Foo.new
		f.bar = 10
		f.bar

		`, 10},
		{`
		class Foo
		  attr_accessor("bar")
		end

		f = Foo.new
		f.bar = 10
		f.bar

		`, 10},
		{`
		class Foo
		  attr_accessor("foo", "bar")
		end

		f = Foo.new
		f.bar = 10
		f.foo = 100
		f.bar + f.foo

		`, 110},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s.\n Input %s", evaluated.(*Error).Message, tt.input)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestNamespace(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		module Foo
		  class Bar
		    def bar
		      10
		    end
		  end
		end

		Foo::Bar.new.bar
		`, 10},
		{`
		class Foo
		  class Bar
		    def bar
		      10
		    end
		  end
		end

		Foo::Bar.new.bar
		`, 10},
		{`
		class Foo
		  def bar
		    100
		  end

		  class Bar
		    def bar
		      10
		    end
		  end
		end

		Foo.new.bar + Foo::Bar.new.bar
		`, 110},
		{`
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

		Baz::Bar.new.bar
		`, 100},
		{`
		module Baz
		  class Bar
		    class Foo
		      def bar
			100
		      end
		    end
		  end
		end

		Baz::Bar::Foo.new.bar
		`, 100},
		{`
		module Baz
		  class Foo
		    def bar
		      100
		    end
		  end

		  class Bar
		    def bar
		      Foo.new.bar
		    end
		  end
		end

		Baz::Bar.new.bar
		`, 100},
		{`
		module Baz
		  class Bar
		    def bar
		      Foo.new.bar
		    end

		    class Foo
		      def bar
			100
		      end
		    end
		  end
		end

		Baz::Bar.new.bar
		`, 100},
		{`
		module Foo
		  class Bar
		    def bar
		      10
		    end
		  end
		end

		module Baz
		  class Bar < Foo::Bar
		    def foo
		      100
		    end
		  end
		end

		b = Baz::Bar.new
		b.foo + b.bar
		`, 110},
		{`
		module A
		  class B
		    class C
		      class D
		        def e
		          10
		        end
		      end
		    end
		  end
		end

		A::B::C::D.new.e
		`, 10},
		{`
		class Foo
		  def self.bar
		    10
		  end
		end

		Object::Foo.bar
		`, 10},

		{`
		Foo = 10

		Object::Foo
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s.\n Input %s", evaluated.(*Error).Message, tt.input)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestRequireSuccess(t *testing.T) {
	input := `
	require("file")

	File.extname("foo.rb")
	`
	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*Error).Message)
	}

	testStringObject(t, evaluated, ".rb")

}

func TestRequireFail(t *testing.T) {
	input := `
	require("bar")
	`
	expected := `Can't require "bar"`

	evaluated := testEval(t, input)

	if !isError(evaluated) {
		t.Fatalf("Should return an error")
	}

	if evaluated.(*Error).Message != expected {
		t.Fatalf("Error message should be '%s'", expected)
	}
}

func TestPrimitiveType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`100.class.name
			`,
			"Integer",
		},
		{
			`Integer.name
			`,
			"Integer",
		},
		{
			`"123".class.name
			`,
			"String",
		},
		{
			`String.name
			`,
			"String",
		},
		{
			`true.class.name
			`,
			"Boolean",
		},
		{
			`Boolean.name
			`,
			"Boolean",
		},
		{
			`
			nil.class.name
			`,
			"Null",
		},
		{
			`
			Integer.name
			`,
			"Integer",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).Message)
		}

		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalCustomConstructor(t *testing.T) {
	input := `
		class Foo
			def initialize(x, y)
				@x = x
				@y = y
			end

			def bar
				@x + @y
			end
		end

		f = Foo.new(10, 20)
		f.bar
	`

	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*Error).Message)
	}

	result, ok := evaluated.(*IntegerObject)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 30 {
		t.Errorf("expect result to be 30. got=%d", result.Value)
	}
}

func TestClassInheritModule(t *testing.T) {
	input := `
module Foo
end

class Bar < Foo
end

a = Bar.new()
	`
	expected := `Module inheritance is not supported: Foo`

	evaluated := testEval(t, input)

	if !isError(evaluated) {
		t.Fatalf("Should return an error when a class inherits a module")
	}

	if evaluated.(*Error).Message != expected {
		t.Fatalf("Error message should be '%s'\n result %s", expected, evaluated)
	}
}
