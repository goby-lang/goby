package vm

import (
	"testing"
)

func init() {

}

func TestLengthMethod(t *testing.T) {
	expected := 5
	array := generateArray(expected)
	m := getBuiltInMethod(t, array, "length")

	result := m(nil, nil).(*IntegerObject).Value

	if int(result) != expected {
		t.Fatalf("Expect length method returns array's length: %d. got=%d", expected, result)
	}
}

func TestPopMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "pop")
	last := m(nil, nil).(*IntegerObject).Value

	if int(last) != 5 {
		t.Fatalf("Expect pop to return array's last  got=%d", last)
	}

	if array.Length() != 4 {
		t.Fatalf("Expect pop remove last elements from array. got=%d", array.Length())
	}
}

func TestPushMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "push")

	six := InitilaizeInteger(6)
	seven := InitilaizeInteger(7)
	m([]Object{six, seven}, nil)

	if array.Length() != 7 {
		t.Fatalf("Expect array's length to be 7(5 + 2). got=%d", array.Length())
	}

	last := array.Elements[array.Length()-1].(*IntegerObject).Value

	if int(last) != 7 {
		t.Fatalf("Expect last object to be 7. got=%d", last)
	}
}

func TestEvalArrayExpression(t *testing.T) {
	input := `
	[1, "234", true]
	`

	evaluated := testEval(t, input)

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be an array. got=%T", evaluated)
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testStringObject(t, arr.Elements[1], "234")
	testBooleanObject(t, arr.Elements[2], true)
}

func TestEvalArrayIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			[][1]
		`, nil},
		{`
			[1, 2, 10, 5][2]
		`, int64(10)},
		{`
			[1, "a", 10, 5][1]
		`, "a"},
		{`
			a = [1, "a", 10, 5]
			a[0]
		`, 1},
		{`
			a = [1, "a", 10, 5]
			a[2] = a[1]
			a[2]

		`, "a"},
		{`
			a = []
			a[10] = 100
			a[10]
		`, 100},
		{`
			a = []
			a[10] = 100
			a[0]
		`, nil},
		{`
			a = [1, 2 ,3 ,5 , 10]
			a[0] = a[1] + a[2] + a[3] * a[4]
			a[0]
		`, 55},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case nil:
			_, ok := evaluated.(*Null)

			if !ok {

				t.Fatalf("expect input: \"%s\"'s result should be Null. got=%T(%s)", tt.input, evaluated, evaluated.Inspect())
			}
		}
	}
}

//func TestEachMethod(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected int
//	}{
//		{`
//		sum = 0
//		puts(self)
//		[1, 2, 3, 4, 5].each do |i|
//		  puts(self)
//		  puts(sum)
//		  sum = sum + i
//		end
//		sum
//		`, 15},
//	}
//
//	for _, tt := range tests {
//		evaluated := testEval(t, tt.input)
//		testIntegerObject(t, evaluated, tt.expected)
//	}
//}

func generateArray(length int) *ArrayObject {
	var elements []Object
	for i := 1; i <= length; i++ {
		int := InitilaizeInteger(i)
		elements = append(elements, int)
	}
	return InitializeArray(elements)
}
