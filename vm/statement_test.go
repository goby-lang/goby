package vm

import (
	"testing"
)

func TestComment(t *testing.T) {
	input := `
	# Comment
	class Foo
		# Comment
		def one
			# Comment
			1 # Comment
			# Comment
		end
		# Comment
	end
	# Comment
	Foo.new.one #=> Comment
	# Comment`

	evaluated := testEval(t, input)
	testIntegerObject(t, evaluated, 1)
}

func TestAssignStatementEvaluation(t *testing.T) {
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

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expectedValue)
	}
}

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

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
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

	evaluated := testEval(t, input)
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
		methodObj, ok := class.Methods.Get(expectedMethod.name)
		if !ok {
			t.Errorf("expect class %s to have method %s.", class.Name, expectedMethod.name)
		}

		method := methodObj.(*Method)
		if method.Name != expectedMethod.name {
			t.Errorf("expect method's name to be %s. got=%s", expectedMethod.name, method.Name)
		}

		if method.argc != len(expectedMethod.params) {
			t.Errorf("expect method %s to have %d parameters. got=%d", method.Name, len(expectedMethod.params), method.argc)
		}
	}
}

//func TestWhileStatement(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected int
//	}{
//		{
//			`
//		i = 10
//		while i > 0
//		  i--
//		end
//		i
//		`, 0},
//		{
//			`
//		a = [1, 2, 3, 4, 5]
//		i = 0
//		while i < a.length
//		  a[i]++
//		  i++
//		end
//		a[4]
//		`, 6},
//	}
//
//	for _, tt := range tests {
//		evaluated := testEval(t, tt.input)
//		testIntegerObject(t, evaluated, tt.expected)
//	}
//}
