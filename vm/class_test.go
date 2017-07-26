package vm

import "testing"

func TestAttrReaderAndWriter(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		class Foo
		  attr_reader :bar

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
		  attr_writer :bar

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
		  attr_writer :bar
		  attr_reader :bar
		end

		f = Foo.new
		f.bar = 10
		f.bar

		`, 10},
		{`
		class Foo
		  attr_accessor :bar
		end

		f = Foo.new
		f.bar = 10
		f.bar

		`, 10},
		{`
		class Foo
		  attr_accessor :foo, :bar
		end

		f = Foo.new
		f.bar = 10
		f.foo = 100
		f.bar + f.foo

		`, 110},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
	expected := `InternalError: Module inheritance is not supported: Foo`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	if !isError(evaluated) {
		t.Fatalf("Should return an error when a class inherits a module")
	}

	err := evaluated.(*Error)

	if err.Message != expected {
		t.Fatalf("Error message should be '%s'. got: %s", expected, err.Message)
	}
	vm.checkCFP(t, 0, 1)
}

func TestClassInstanceVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_get("@foo")
		`, 1},
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_set("@bar", 100)
		Bar.instance_variable_set("@foo", 20)
		Bar.instance_variable_get("@foo") + Bar.instance_variable_get("@bar")
		`, 120},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, 30)
	vm.checkCFP(t, 0, 0)
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

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestMonkeyPatchBuiltInClass(t *testing.T) {
	input := `
	class String
	  def buz
	    "buz"
	  end
	end

	"123".buz
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, "buz")
	vm.checkCFP(t, 0, 0)
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
		{`
		class X
		  Bar = 100
		end

		module Foo
		  Bar = 10

		  class Baz < X
			def self.result
			  Bar
			end
		  end
		end

		Foo::Baz.result`, 10},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestRequireRelative(t *testing.T) {
	input := `
	require_relative("../test_fixtures/require_test/foo")

	fifty = Foo.bar(5)

	Foo.baz do |hundred|
	  hundred + fifty + Bar.baz
	end
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, 160)
	vm.checkCFP(t, 0, 0)
}

func TestRequireSuccess(t *testing.T) {
	input := `
	require "file"

	File.extname("foo.rb")
	`
	vm := initTestVM()
	evaluated := vm.testEval(t, input)
	checkExpected(t, 0, evaluated, ".rb")
	vm.checkCFP(t, 0, 0)
}

func TestRequireFail(t *testing.T) {
	input := `
	require "bar"
	`
	expected := `InternalError: Can't require "bar"`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	if !isError(evaluated) {
		t.Fatalf("Should return an error")
	}

	err := evaluated.(*Error)
	if err.Message != expected {
		t.Fatalf("Error message should be '%s'. got: %s", expected, err.Message)
	}
	vm.checkCFP(t, 0, 1)
}

func TestGeneralIsNilMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.is_nil`, false},
		{`"Hello World".is_nil`, false},
		{`(2..10).is_nil`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_nil`, false},
		{`[1, 2, 3, 4, 5].is_nil`, false},
		{`true.is_nil`, false},
		{`String.is_nil`, false},
		{`nil.is_nil`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestGeneralIsNilMethodFail(t *testing.T) {
	testsFail := []struct {
		input  string
		errMsg string
	}{
		{`123.is_nil("Hello")`, "ArgumentError: Expect 0 argument. got: 1"},
		{`"Fail".is_nil("Hello")`, "ArgumentError: Expect 0 argument. got: 1"},
		{`[1, 2, 3].is_nil("Hello")`, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2, c: 3 }.is_nil("Hello")`, "ArgumentError: Expect 0 argument. got: 1"},
		{`(1..10).is_nil("Hello")`, "ArgumentError: Expect 0 argument. got: 1"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, ArgumentError, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestClassGeneralComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Integer == 123`, false},
		{`Integer == "123"`, false},
		{`Integer == "124"`, false},
		{`Integer == (1..3)`, false},
		{`Integer == { a: 1, b: 2 }`, false},
		{`Integer == [1, "String", true, 2..5]`, false},
		{`Integer == Integer`, true},
		{`Integer == String`, false},
		{`123.class == Integer`, true},
		{`Integer == Object`, false},
		{`Integer.superclass == Object`, true},
		{`123.class.superclass == Object`, true},
		{`Integer != 123`, true},
		{`Integer != "123"`, true},
		{`Integer != "124"`, true},
		{`Integer != (1..3)`, true},
		{`Integer != { a: 1, b: 2 }`, true},
		{`Integer != [1, "String", true, 2..5]`, true},
		{`Integer != Integer`, false},
		{`Integer != String`, true},
		{`123.class != Integer`, false},
		{`Integer != Object`, true},
		{`Integer.superclass != Object`, false},
		{`123.class.superclass != Object`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestGeneralAssignmentByOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`a = 123;    a ||= 456;                  a;`, 123},
		{`a = 123;    a ||= true;                 a;`, 123},
		{`a = "Goby"; a ||= "Fish";               a;`, "Goby"},
		{`a = (1..3); a ||= [1, 2, 3];          a.to_s;`, "(1..3)"},
		{`a = false;  a ||= 123;                  a;`, 123},
		{`a = nil;    a ||= { b: 1 };             a["b"];`, 1},
		{`a = false;  a ||= false;                a;`, false},
		{`a = nil;    a ||= false;                a;`, false},
		{`a = false;  a ||= nil;                  a;`, nil},
		{`a = nil;    a ||= nil;                  a;`, nil},
		{`a = false;  a ||= nil || false;         a;`, false},
		{`a = false;  a ||= false || nil;         a;`, nil},
		{`a = false;  a ||= true && false || nil; a;`, nil},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestGeneralIsAMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.is_a(Integer)`, true},
		{`123.is_a(Object)`, true},
		{`123.is_a(String)`, false},
		{`123.is_a(Range)`, false},
		{`"Hello World".is_a(String)`, true},
		{`"Hello World".is_a(Object)`, true},
		{`"Hello World".is_a(Boolean)`, false},
		{`"Hello World".is_a(Array)`, false},
		{`(2..10).is_a(Range)`, true},
		{`(2..10).is_a(Object)`, true},
		{`(2..10).is_a(Null)`, false},
		{`(2..10).is_a(Hash)`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a(Hash)`, true},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a(Object)`, true},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a(Class)`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a(Array)`, false},
		{`[1, 2, 3, 4, 5].is_a(Array)`, true},
		{`[1, 2, 3, 4, 5].is_a(Object)`, true},
		{`[1, 2, 3, 4, 5].is_a(Null)`, false},
		{`[1, 2, 3, 4, 5].is_a(String)`, false},
		{`true.is_a(Boolean)`, true},
		{`true.is_a(Object)`, true},
		{`true.is_a(Array)`, false},
		{`true.is_a(Integer)`, false},
		{`String.is_a(Class)`, true},
		{`String.is_a(String)`, false},
		{`String.is_a(Array)`, false},
		{`nil.is_a(Null)`, true},
		{`nil.is_a(Object)`, true},
		{`nil.is_a(String)`, false},
		{`nil.is_a(Range)`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestGeneralIsAMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`123.is_a`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`Class.is_a`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`123.is_a(123, 456)`, ArgumentError, "ArgumentError: Expect 1 argument. got: 2"},
		{`123.is_a(Integer, String)`, ArgumentError, "ArgumentError: Expect 1 argument. got: 2"},
		{`123.is_a(true)`, TypeError, "TypeError: Expect argument to be Class. got: Boolean"},
		{`Class.is_a(true)`, TypeError, "TypeError: Expect argument to be Class. got: Boolean"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestClassNameClassMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Integer.name`, "Integer"},
		{`String.name`, "String"},
		{`Boolean.name`, "Boolean"},
		{`Range.name`, "Range"},
		{`Hash.name`, "Hash"},
		{`Array.name`, "Array"},
		{`Class.name`, "Class"},
		{`Object.name`, "Object"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestClassNameClassMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei".name`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'name' for Taipei"},
		{`123.name`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'name' for 123"},
		{`true.name`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'name' for true"},
		{`Integer.name(Integer)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`String.name(Hash, Array)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestClassSuperclassClassMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Integer.superclass.name`, "Object"},
		{`String.superclass.name`, "Object"},
		{`Boolean.superclass.name`, "Object"},
		{`Range.superclass.name`, "Object"},
		{`Hash.superclass.name`, "Object"},
		{`Array.superclass.name`, "Object"},
		{`Object.superclass.name`, "Object"},
		{`Class.superclass.name`, "Object"},
		{`
		module Bar; end
		class Foo
		  include Bar
		end
		Foo.superclass.name
		`, "Object"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestClassSuperclassClassMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei".superclass`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'superclass' for Taipei"},
		{`123.superclass`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'superclass' for 123"},
		{`true.superclass`, UndefinedMethodError, "UndefinedMethodError: Undefined Method 'superclass' for true"},
		{`Integer.superclass(Integer)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`String.superclass(Hash, Array)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestClassSingletonClassClassMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`Integer.singleton_class.name`, "#<Class:Integer>"},
		{`Integer.singleton_class.superclass.name`, "#<Class:Object>"},
		{`
		class Bar; end
		Bar.singleton_class.name
		`, "#<Class:Bar>"},
		{`
		class Bar; end
		class Foo < Bar; end
		Foo.singleton_class.superclass.name
		`, "#<Class:Bar>"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}
