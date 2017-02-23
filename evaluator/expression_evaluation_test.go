package evaluator

import (
	"github.com/st0012/rooby/object"
	"testing"
)

func TestClassMethodEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			class Bar {
				def self.foo {
					10
				}
			}
			Bar.foo;
			`,
			10,
		},
		{
			`
			class Foo {
				def self.foo {
					10
				}
			}

			class Bar < Foo {}
			Bar.foo;
			`,
			10,
		},
		{
			`
			class Foo {
				def self.foo {
					10
				}
			}

			class Bar < Foo {
				def self.foo {
					100
				}
			}
			Bar.foo;
			`,
			100,
		},
		{
			`
			class Bar {
				def self.foo {
					self.bar();
				}

				def self.bar {
					100
				}

				def bar {
					1000
				}
			}
			Bar.foo();
			`,
			100,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestSelfExpressionEvaluation(t *testing.T) {
	tests := []struct {
		input        string
		expected_obj string
	}{
		{`self`, object.BASE_OBJECT_OBJ},
		{
			`
			class Bar {
				def whoami {
					self
				}
			}

			Bar.new.whoami;
		`, object.BASE_OBJECT_OBJ},
		{
			`
			class Foo {
				Self = self;

				def get_self {
					Self
				}
			}

			Foo.new.get_self;
			`,
			object.CLASS_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		if isError(evaluated) {
			t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
		}

		if string(evaluated.Type()) != tt.expected_obj {
			t.Fatalf("expect self to return %s. got=%s", string(tt.expected_obj), evaluated.Type())
		}
	}
}

func TestEvalInstanceVariable(t *testing.T) {
	input := `
		class Foo {
			def set(x) {
				@x = x;
			}

			def get {
				@x
			}

			def double_get {
				self.get() * 2;
			}
		}

		class Bar {
			def set(x) {
				@x = x;
			}

			def get {
				@x
			}
		}

		f1 = Foo.new;
		f1.set(10);

		f2 = Foo.new;
		f2.set(20);

		b = Bar.new;
		b.set(10)

		f2.double_get() + f1.get() + b.get();
	`

	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
	}

	result, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 60 {
		t.Fatalf("expect result to be 60. got=%d", result.Value)
	}
}

func TestEvalInstanceMethodCall(t *testing.T) {
	input := `

		class Bar {
			def set(x) {
				@x = x;
			}
		}

		class Foo < Bar {
			def add(x, y) {
				x + y
			}
		}

		class FooBar < Foo {
			def get {
				@x
			}
		}

		fb = FooBar.new;
		fb.set(100);
		fb.add(10, fb.get());
	`

	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
	}

	result, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 110 {
		t.Errorf("expect result to be 110. got=%d", result.Value)
	}
}

func TestEvalCustomInitializeMethod(t *testing.T) {
	input := `
		class Foo {
			def initialize(x, y) {
				@x = x;
				@y = y;
			}

			def bar {
				@x + @y;
			}
		}

		f = Foo.new(10, 20);
		f.bar;
	`

	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
	}

	result, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 30 {
		t.Errorf("expect result to be 30. got=%d", result.Value)
	}
}

func TestEvalClassInheritance(t *testing.T) {
	input := `
		class Foo {
			def add(x, y) {
				x + y
			}
		}
		Foo.new.add(10, 11)
	`

	evaluated := testEval(t, input)

	if isError(evaluated) {
		t.Fatalf("got Error: %s", evaluated.(*object.Error).Message)
	}

	result, ok := evaluated.(*object.Integer)

	if !ok {
		t.Errorf("expect result to be an integer. got=%T", evaluated)
	}

	if result.Value != 21 {
		t.Errorf("expect result to be 21. got=%d", result.Value)
	}
}

func TestEvalInfixIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`if (10 > 5) {
				100
			} else {
				-10
			}
			`,
			100,
		},
		{
			`if (5 != 5) {
				false
			} else {
				true
			}
			`,
			true,
		},
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, tt.expected.(int64))
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case nil:
			testNullObject(t, evaluated)
		}

	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Stan " + "Lo"`, "Stan Lo"},
		{`"Dog" + "&" + "Cat"`, "Dog&Cat"},
		{`"Dog" == "Dog"`, true},
		{`"1234" > "123"`, true},
		{`"1234" < "123"`, false},
		{`"1234" != "123"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case string:
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestEvalInfixBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBangPrefixExpression(t *testing.T) {
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

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalMinusPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"-10", -10},
		{"--10", 10},
		{"--5", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"st0012"`, "st0012"},
		{`'Monkey'`, "Monkey"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
