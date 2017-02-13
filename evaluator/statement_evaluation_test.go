package evaluator

import (
	"github.com/st0012/rooby/object"
	"testing"
)

//func TestClosures(t *testing.T) {
//	input := `
//let newAdder = fn(x) {
//  fn(y) { x + y };
//};
//
//let addTwo = newAdder(2);
//addTwo(2);`
//	testIntegerObject(t, testEval(input), 4)
//}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 5; let b = 10; let c = if (a > b) { 100; } else { 50; }", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
    if (10 > 1) {
      if (10 > 1) {
	return 10;
      }

      return 1;
    }
    `,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClassStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`class Foo {}`, "Foo"},
		{
			`class Foo {
				def bar() {
					x;
				}
			}`, "Foo"},
		{
			`class Bar {}
			class Foo {}
			Bar
			`, "Bar"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testClassObject(t, evaluated, tt.expected)
	}
}

func TestDefStatement(t *testing.T) {
	input := `
		class Foo {
			def bar(x, y) {
				x + y;
			}

			def foo(y) {
				y;
			}
		}
	`

	evaluated := testEval(t, input)
	class := evaluated.(*object.Class)

	expectedMethods := []struct {
		name   string
		params []string
	}{
		{name: "foo", params: []string{"y"}},
		{name: "bar", params: []string{"x", "y"}},
	}

	for _, expectedMethod := range expectedMethods {
		methodObj, ok := class.Body.Get("_method_" + expectedMethod.name)
		if !ok {
			t.Errorf("expect class %s to have method %s.", class.Name, expectedMethod.name)
		}

		method := methodObj.(*object.Method)
		if method.Name != expectedMethod.name {
			t.Errorf("expect method's name to be %s. got=%s", expectedMethod.name, method.Name)
		}
		for i, expectedParam := range expectedMethod.params {
			if method.Parameters[i].Value != expectedParam {
				t.Errorf("expect method %s's parameters to have %s. got=%s", expectedMethod.name, expectedParam, method.Parameters[i].Value)
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
	    if (10 > 1) {
	      if (10 > 1) {
		return true + false;
	      }

	      return 1;
	    }
	    `,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		//{
		//	"let add = fn(x, y) { x + y; }; add(1, 2, 3);",
		//	"wrong arguments: expect=2, got=3",
		//},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}