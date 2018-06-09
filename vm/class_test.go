package vm

import "testing"

func TestClassClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Class.class.name`, "Class"},
		{`Class.superclass.name`, "Module"},
		{`Module.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestAttrReaderAndWriter(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassInheritModuleError(t *testing.T) {
	input := `module Foo
end

class Bar < Foo
end

a = Bar.new()
	`
	expected := `InternalError: Module inheritance is not supported: Foo`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	checkErrorMsg(t, i, evaluated, expected)
	v.checkCFP(t, 0, 1)
	v.checkSP(t, 0, 3)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassInstanceVariableFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_get
		`, "ArgumentError: Expect 1 arguments. got: 0", 1},
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_get("@foo", 2)
		`, "ArgumentError: Expect 1 arguments. got: 2", 1},
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_set
				`, "ArgumentError: Expect 2 arguments. got: 0", 1},
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_set("@bar")
				`, "ArgumentError: Expect 2 arguments. got: 1", 1},
		{`
		class Bar
		  @foo = 1
		end

		Bar.instance_variable_set("@bar", 2, 3)
				`, "ArgumentError: Expect 2 arguments. got: 3", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestCustomClassConstructor(t *testing.T) {
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

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, 30)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestClassNamespace(t *testing.T) {
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		{
			`
			Object.class.name
			`, "Class",
		},
		{
			`
			Class.class.name
			`, "Class",
		},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestForbiddenInclusionWithClass(t *testing.T) {
	input := `class Foo
end

class Bar
  include Foo
end
	`
	expected := `TypeError: Expect argument to be a module. got=Class`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	checkErrorMsg(t, i, evaluated, expected)
	v.checkCFP(t, 0, 2)
	v.checkSP(t, 0, 1)
}

func TestForbiddenExtensionWithClass(t *testing.T) {
	input := `class Foo
end

class Bar
  extend Foo
end
	`
	expected := `TypeError: Expect argument to be a module. got=Class`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	checkErrorMsg(t, i, evaluated, expected)
	v.checkCFP(t, 0, 2)
	v.checkSP(t, 0, 1)
}

// Method tests

func TestMethodsMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		class C
		  def hola
		  end

		  def hi
		  end
		end
		C.new.methods.first(2) == ["hi", "hola"]
		`, true},
		{`
		class C
		  def hi
		  end
		end
		class D < C
		  def hola
		  end
		end
		D.new.methods.first(2) == ["hola", "hi"]
		`, true},
		{`
		class C
		  def hi
		  end
		end
		c = C.new
		def c.hola
		end
		c.methods.first(2) == ["hola", "hi"]
		`, true},
		{`
		class C
		end
		C.new.methods.include?("to_s")
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestAncestorsMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		class C
		end
		C.ancestors == [C, Object]
		`, true},
		{`
		module M
		end
		class C
		end
		class C2 < C
		  include M
		end
		class C3 < C2
		end
		C3.ancestors == [C3, C2, M, C, Object]
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestBuiltinClassMonkeyPatching(t *testing.T) {
	input := `
	class String
	  def buz
	    "buz"
	  end
	end

	"123".buz
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, "buz")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestClassGreaterThanMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		Object > Array
		`, true},
		{`
		Array > Object
		`, false},
		{`
		Object > Object
		`, false},
		{`
		(Array > Hash).nil?
		`, true},
		{`
		module M
		end
		class C
		end
		class C2 < C
		  include M
		end
		class C3 < C2
		end
		M > C3
		`, true},
		{`
		module M
		end
		class C
		end
		(M > C).nil?
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassGreaterThanMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Array > 1`, "TypeError: Expect argument to be a module. got=Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestClassGreaterThanOrEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		Object >= Array
		`, true},
		{`
		Array >= Object
		`, false},
		{`
		Object >= Object
		`, true},
		{`
		(Array >= Hash).nil?
		`, true},
		{`
		module M
		end
		class C
		end
		class C2 < C
		  include M
		end
		class C3 < C2
		end
		M >= C3
		`, true},
		{`
		module M
		end
		class C
		end
		(M >= C).nil?
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassGreaterThanOrEqualMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Array >= 1`, "TypeError: Expect argument to be a module. got=Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestClassLesserThanMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		Object < Array
		`, false},
		{`
		Array < Object
		`, true},
		{`
		Object < Object
		`, false},
		{`
		(Array < Hash).nil?
		`, true},
		{`
		module M
		end
		class C
		end
		class C2 < C
		  include M
		end
		class C3 < C2
		end
		C3 < M
		`, true},
		{`
		module M
		end
		class C
		end
		(M < C).nil?
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassLesserThanMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Array < 1`, "TypeError: Expect argument to be a module. got=Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestClassLesserThanOrEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		Object <= Array
		`, false},
		{`
		Array <= Object
		`, true},
		{`
		Object <= Object
		`, true},
		{`
		(Array <= Hash).nil?
		`, true},
		{`
		module M
		end
		class C
		end
		class C2 < C
		  include M
		end
		class C3 < C2
		end
		C3 <= M
		`, true},
		{`
		module M
		end
		class C
		end
		(M <= C).nil?
		`, true},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassLesserThanOrEqualMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Array <= 1`, "TypeError: Expect argument to be a module. got=Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestInheritsMethodMissingMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		class Bar
		end
		
		Bar.new.inherits_method_missing?

`, false},
		{`
		class Bar
		  inherits_method_missing
		end
		
		Bar.new.inherits_method_missing?

`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*Error).message)
		}

		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestSendMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		def foo
		  10
		end

		send(:foo)
		`, 10},
		{`
		class Foo
		  def bar
		    10
		  end
		end

		Foo.new.send(:bar)
		`, 10},
		{`
		class Foo
		  def self.bar
		    10
		  end
		end

		Foo.send(:bar)
		`, 10},
		{`
		def foo(x)
		  10 + x
		end

		send(:foo, 5)
		`, 15},
		{`
		class Foo
		  def bar(x)
		    10 + x
		  end
		end

		Foo.new.send(:bar, 5)
		`, 15},
		{`
		class Foo
		  def self.bar(x)
		    10 + x
		  end
		end

		Foo.send(:bar, 5)
		`, 15},
		{`
		class Math
		  def self.add(x, y)
		    x + y
		  end
		end

		Math.send(:add, 10, 15)
		`, 25},
		{`
		class Foo
		  def bar(x, y)
		    yield x, y
		  end
		end
		a = Foo.new
		a.send(:bar, 7, 8) do |i, j| i * j; end
		`, 56},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestRequireRelativeMethod(t *testing.T) {
	input := `
	require_relative("../test_fixtures/require_test/foo")

	fifty = Foo.bar(5)

	Foo.baz do |hundred|
	  hundred + fifty + Bar.baz
	end
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, 160)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestRequireMethod(t *testing.T) {
	input := `
	require "uri"

	u = URI.parse("http://example.com")
	u.scheme
	`
	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())
	VerifyExpected(t, 0, evaluated, "http")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestRequireMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`require "bar"`, `InternalError: Can't require "bar"`, 1},
		{`require "db", "json"`, `ArgumentError: Expect 1 argument. got: 2`, 1},
		{`require_relative "db", "json"`, `ArgumentError: Expect 1 argument. got: 2`, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestRaiseMethod(t *testing.T) {
	testsFail := []struct {
		input       string
		expected    string
		expectedCFP int
		expectedSP  int
	}{
		{`raise`, "InternalError: ", 1, 1},
		{`raise "Foo"`, "InternalError: 'Foo'", 1, 1},
		{`
		class BarError; end
		raise BarError, "Foo"`, "BarError: 'Foo'", 1, 1},
		{`
		class FooError; end

		def raise_foo
		  raise FooError, "Foo"
		end

		raise_foo
		`,
			// Expect CFP to be 2 is because the `raise_foo`'s frame is not popped
			// Expect SP to be 2 cause the program got stopped before it replaces receiver with the return value (error)
			// TODO: This means we need to pop error object when implementing `rescue`
			"FooError: 'Foo'", 2, 2},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, tt.expectedSP)
	}
}

func TestRaiseMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`raise "Foo", "Bar"`, "ArgumentError: Expect error class, got: String", 1},
		{`
		class BarError; end
		raise BarError, "Foo", "Bar"`, "ArgumentError: Expect at most 2 arguments. got: 3", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestResponseToMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`
		1.respond_to? :to_i
		`, true},
		{`
		"string".respond_to? "+"
		`, true},
		{`
		1.respond_to? :numerator
		`, false},
		{`
		Class.respond_to? "respond_to?"
		`, true},
		{`
		Class.respond_to? :phantom
		`, false},
	}
	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// With the current framework, only exit() failures can be tested.
func TestExitMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`exit("abc")`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`exit(1, 2)`, "ArgumentError: Expected at most 1 argument; got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestGeneralIsAMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.is_a?(Integer)`, true},
		{`123.is_a?(Object)`, true},
		{`123.is_a?(String)`, false},
		{`123.is_a?(Range)`, false},
		{`"Hello World".is_a?(String)`, true},
		{`"Hello World".is_a?(Object)`, true},
		{`"Hello World".is_a?(Boolean)`, false},
		{`"Hello World".is_a?(Array)`, false},
		{`(2..10).is_a?(Range)`, true},
		{`(2..10).is_a?(Object)`, true},
		{`(2..10).is_a?(Null)`, false},
		{`(2..10).is_a?(Hash)`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a?(Hash)`, true},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a?(Object)`, true},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a?(Class)`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.is_a?(Array)`, false},
		{`[1, 2, 3, 4, 5].is_a?(Array)`, true},
		{`[1, 2, 3, 4, 5].is_a?(Object)`, true},
		{`[1, 2, 3, 4, 5].is_a?(Null)`, false},
		{`[1, 2, 3, 4, 5].is_a?(String)`, false},
		{`true.is_a?(Boolean)`, true},
		{`true.is_a?(Object)`, true},
		{`true.is_a?(Array)`, false},
		{`true.is_a?(Integer)`, false},
		{`String.is_a?(Class)`, true},
		{`String.is_a?(String)`, false},
		{`String.is_a?(Array)`, false},
		{`nil.is_a?(Null)`, true},
		{`nil.is_a?(Object)`, true},
		{`nil.is_a?(String)`, false},
		{`nil.is_a?(Range)`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGeneralIsAMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`123.is_a?`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`Class.is_a?`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`123.is_a?(123, 456)`, "ArgumentError: Expect 1 argument. got: 2", 1},
		{`123.is_a?(Integer, String)`, "ArgumentError: Expect 1 argument. got: 2", 1},
		{`123.is_a?(true)`, "TypeError: Expect argument to be Class. got: Boolean", 1},
		{`Class.is_a?(true)`, "TypeError: Expect argument to be Class. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestGeneralIsNilMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`123.nil?`, false},
		{`"Hello World".nil?`, false},
		{`(2..10).nil?`, false},
		{`{ a: 1, b: "2", c: ["Goby", 123] }.nil?`, false},
		{`[1, 2, 3, 4, 5].nil?`, false},
		{`true.nil?`, false},
		{`String.nil?`, false},
		{`nil.nil?`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGeneralIsNilMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`123.nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`"Fail".nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`[1, 2, 3].nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2, c: 3 }.nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`(1..10).nil?("Hello")`, "ArgumentError: Expect 0 argument. got: 1", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassNameClassMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Taipei".name`, "UndefinedMethodError: Undefined Method 'name' for Taipei", 1},
		{`123.name`, "UndefinedMethodError: Undefined Method 'name' for 123", 1},
		{`true.name`, "UndefinedMethodError: Undefined Method 'name' for true", 1},
		{`Integer.name(Integer)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`String.name(Hash, Array)`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
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
		{`Module.superclass.name`, "Object"},
		{`
		module Bar; end
		class Foo
		  include Bar
		end
		Foo.superclass.name
		`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestClassSuperclassClassMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Taipei".superclass`, "UndefinedMethodError: Undefined Method 'superclass' for Taipei", 1},
		{`123.superclass`, "UndefinedMethodError: Undefined Method 'superclass' for 123", 1},
		{`true.superclass`, "UndefinedMethodError: Undefined Method 'superclass' for true", 1},
		{`Integer.superclass(Integer)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`String.superclass(Hash, Array)`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestClassSingletonClassMethod(t *testing.T) {
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
		// Check if this works on non-class objects
		{`'a'.singleton_class.to_s.slice(1..16).to_s`, "<Class:#<String:"},
		{`1.singleton_class.to_s.slice(1..17).to_s`, "<Class:#<Integer:"},
		{`nil.singleton_class.to_s.slice(1..14).to_s`, "<Class:#<Null:"},
		{`[1,2].singleton_class.to_s.slice(1..15).to_s`, "<Class:#<Array:"},
		{`{key: "value"}.singleton_class.to_s.slice(1..14).to_s`, "<Class:#<Hash:"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestInstanceEvalMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		class Foo
		  def initialize
			@secret = 99
		  end
		end

		f = Foo.new
		f.instance_eval do
		  @secret
		end
`, 99},
		{`
		string = "String"
		string.instance_eval do
		  def new_method
			self.reverse
		  end
		end
		string.new_method
`, "gnirtS"},
		{`"a".instance_eval`, "a"},
		{`"a".instance_eval do end`, "a"},
		{`
		class Foo
		  def bar
			10
		  end
		end

		block = Block.new do
		  self.bar
		end

		Foo.new.instance_eval block
		`, 10},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestInstanceEvalMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"s".instance_eval(1, 1)`, "ArgumentError: Expect at most 1 arguments. got: 2", 1},
		{`"s".instance_eval(true)`, "TypeError: Expect argument to be Block. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestObjectIdMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// immutable objects
		{`Object.new.object_id.is_a?(Integer)`, true},
		{`nil.object_id == nil.object_id`, true},
		{`true.object_id == true.object_id`, true},
		{`false.object_id == false.object_id`, true},
		{`CONSTANT=1; CONSTANT.object_id == CONSTANT.object_id`, true},
		{`ARGV.object_id == ARGV.object_id`, true},
		{`STDIN.object_id == STDIN.object_id`, true},
		{`STDOUT.object_id == STDOUT.object_id`, true},
		{`STDERR.object_id == STDERR.object_id`, true},
		{`ENV.object_id == ENV.object_id`, true},
		{`Class.object_id == Class.object_id`, true},
		{`Object.object_id == Object.object_id`, true},
		{`Integer.object_id == Integer.object_id`, true},
		// other objects
		{`a = 1.object_id; b = 1.object_id; a == b`, false},
		{`a = "a".object_id; b = "a".object_id; a == b`, false},
		{`a = 1.object_id; b = a; a.object_id == b.object_id`, true},
		{`a = "a".object_id; b = a; a.object_id == b.object_id`, true},
		{`a = 'a'.object_id; b = a; a.object_id == b.object_id`, true},
		{`a = 1..100; b = a; a.object_id == b.object_id`, true},
		{`a = "3.14".to_f; b = a; a.object_id == b.object_id`, true},
		{`a = [1, 2, 3];b = a; a.object_id == b.object_id`, true},
		{`a = {key: 2} ;b = a; a.object_id == b.object_id`, true},
		{`a = :symbol ;b = a; a.object_id == b.object_id`, true},
		{`a = Regexp.new("aa") ;b = a; a.object_id == b.object_id`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
