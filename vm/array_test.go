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

	result := m(nil, nil, nil).(*IntegerObject).Value

	if int(result) != expected {
		t.Fatalf("Expect length method returns array's length: %d. got=%d", expected, result)
	}
}

func TestPopMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "pop")
	last := m(nil, nil, nil).(*IntegerObject).Value

	if int(last) != 5 {
		t.Fatalf("Expect pop to return array's last  got=%d", last)
	}

	if array.length() != 4 {
		t.Fatalf("Expect pop remove last elements from array. got=%d", array.length())
	}
}

func TestPushMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "push")

	six := initilaizeInteger(6)
	seven := initilaizeInteger(7)
	m(nil, []Object{six, seven}, nil)

	if array.length() != 7 {
		t.Fatalf("Expect array's length to be 7(5 + 2). got=%d", array.length())
	}

	last := array.Elements[array.length()-1].(*IntegerObject).Value

	if int(last) != 7 {
		t.Fatalf("Expect last object to be 7. got=%d", last)
	}
}

func TestShiftMethod(t *testing.T) {
	array := initializeArray([]Object{initilaizeInteger(1), initilaizeInteger(2), initilaizeInteger(3), initilaizeInteger(4)})
	second := initializeArray([]Object{initilaizeInteger(2), initilaizeInteger(3), initilaizeInteger(4)})

	m := getBuiltInMethod(t, array, "shift")
	first := m(nil, nil, nil)

	testArrayObject(t, array, second)
	testIntegerObject(t, first, 1)
}

func TestShiftMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.shift(3, 3, 4, 5)
		`, newError("Expect 0 argument. got=4")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
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
			[1, 2, 3][100]
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
		{`
			[].at(1)
		`, nil},
		{`
			[1, 2, 10, 5].at(2)
		`, int64(10)},
		{`
			[1, "a", 10, 5].at(1)
		`, "a"},
		{`
			a = [1, "a", 10, 5]
			a.at(0)
		`, 1},
		{`
			a = [1, "a", 10, 5]
			a[2] = a.at(1)
			a[2]
		`, "a"},
		{`
			a = [1, 2, 3, 5, 10]
			a[0] = a.at(1) + a.at(2) + a.at(3) * a.at(4)
			a.at(0)
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
			_, ok := evaluated.(*NullObject)

			if !ok {

				t.Fatalf("expect input: \"%s\"'s result should be Null. got=%T(%s)", tt.input, evaluated, evaluated.Inspect())
			}
		}
	}
}

func TestEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		sum = 0
		[1, 2, 3, 4, 5].each do |i|
		  sum = sum + i
		end
		sum
		`, 15},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEachIndexMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`
		sum = 0
		[2, 3, 40, 5, 22].each_index do |i|
		  sum = sum + i
		end
		sum
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMapMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [1, 2, 7]
		a.map do |i|
			i + 3
		end
		`, initializeArray([]Object{initilaizeInteger(4), initilaizeInteger(5), initilaizeInteger(10)})},
		{`
		a = [true, false, true, false, true ]
		a.map do |i|
			!i
		end
		`, initializeArray([]Object{FALSE, TRUE, FALSE, TRUE, FALSE})},
		{`
		a = ["1", "sss", "qwe"]
		a.map do |i|
			i + "1"
		end
		`, initializeArray([]Object{initializeString("11"), initializeString("sss1"), initializeString("qwe1")})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestSelectMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [1, 2, 3, 4, 5]
		a.select do |i|
			i > 3
		end
		`, initializeArray([]Object{initilaizeInteger(4), initilaizeInteger(5)})},
		{`
		a = [true, false, true, false, true ]
		a.select do |i|
			i
		end
		`, initializeArray([]Object{TRUE, TRUE, TRUE})},
		{`
		a = ["test", "not2", "3", "test", "5"]
		a.select do |i|
			i == "test"
		end
		`, initializeArray([]Object{initializeString("test"), initializeString("test")})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestClearMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [1, 2, 3]
		a.clear
		`, initializeArray([]Object{})},
		{`
		a = []
		a.clear
		`, initializeArray([]Object{})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestConcatMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [1, 2]
		a.concat([3], [4])
		`, initializeArray([]Object{initilaizeInteger(1), initilaizeInteger(2), initilaizeInteger(3), initilaizeInteger(4)})},
		{`
		a = []
		a.concat([1], [2], ["a", "b"], [3], [4])
		`, initializeArray([]Object{initilaizeInteger(1), initilaizeInteger(2), initializeString("a"), initializeString("b"), initilaizeInteger(3), initilaizeInteger(4)})},
		{`
		a = [1, 2]
		a.concat()
		`, initializeArray([]Object{initilaizeInteger(1), initilaizeInteger(2)})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestConcatMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.concat(3)
		`, newError("Expect argument to be Array. got=*vm.IntegerObject")},
		{`
		a = []
		a.concat("a")
		`, newError("Expect argument to be Array. got=*vm.StringObject")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestCountMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *IntegerObject
	}{
		{`
		a = [1, 2]
		a.count
		`, initilaizeInteger(2)},
		{`
		a = [1, 2]
		a.count(1)
		`, initilaizeInteger(1)},
		{`
		a = ["a", "bb", "c", "db", "bb", 2]
		a.count("bb")
		`, initilaizeInteger(2)},
		{`
		a = [true, true, true, false, true]
		a.count(true)
		`, initilaizeInteger(4)},
		{`
		a = []
		a.count(true)
		`, initilaizeInteger(0)},
		{`
		a = [1, 2, 3, 4, 5, 6, 7, 8]
		a.count do |i|
			i > 3
		end
		`, initilaizeInteger(5)},
		{`
		a = ["a", "bb", "c", "db", "bb"]
		a.count do |i|
			i.size > 1
		end
		`, initilaizeInteger(3)},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected.Value)
	}
}

func TestCountMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.count(3, 3)
		`, newError("Expect one argument. got=2")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)

		err, ok := evaluated.(*Error)
		if !ok || err.Class.ReturnName() != ArgumentError {
			t.Errorf("Expect ArgumentError. got=%T (%+v)", err, err)
		}
	}
}

func TestRotateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [1, 2]
		a.rotate
		`, initializeArray([]Object{initilaizeInteger(2), initilaizeInteger(1)})},
		{`
		a = [1, 2, 3, 4]
		a.rotate(2)
		`, initializeArray([]Object{initilaizeInteger(3), initilaizeInteger(4), initilaizeInteger(1), initilaizeInteger(2)})},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestRotateMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.rotate("a")
		`, newError("Expect index argument to be Integer. got=*vm.StringObject")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestFirstMethod(t *testing.T) {
	testsInt := []struct {
		input    string
		expected *IntegerObject
	}{
		{`
		a = [1, 2]
		a.first
		`, initilaizeInteger(1)},
	}

	for _, tt := range testsInt {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected.Value)
	}

	testsArray := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [3, 4, 5, 1, 6]
		a.first(2)
		`, initializeArray([]Object{initilaizeInteger(3), initilaizeInteger(4)})},
		{`
		a = ["a", "b", "d", "q"]
		a.first(2)
		`, initializeArray([]Object{initializeString("a"), initializeString("b")})},
	}

	for _, tt := range testsArray {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestFirstMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.first("a")
		`, newError("Expect index argument to be Integer. got=*vm.StringObject")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", err.Message, tt.expected.Message)
		}
	}
}

func TestLastMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected *StringObject
	}{
		{`
		a = [1, 2, "a", 2, "b"]
		a.last
		`, initializeString("b")},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testStringObject(t, evaluated, tt.expected.Value)
	}

	testsArray := []struct {
		input    string
		expected *ArrayObject
	}{
		{`
		a = [3, 4, 5, 1, 6]
		a.last(3)
		`, initializeArray([]Object{initilaizeInteger(5), initilaizeInteger(1), initilaizeInteger(6)})},
		{`
		a = ["a", "b", "d", "q"]
		a.last(2)
		`, initializeArray([]Object{initializeString("d"), initializeString("q")})},
	}

	for _, tt := range testsArray {
		evaluated := testEval(t, tt.input)
		testArrayObject(t, evaluated, tt.expected)
	}
}

func TestLastMethodFail(t *testing.T) {
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`
		a = [1, 2]
		a.last("l")
		`, newError("Expect index argument to be Integer. got=*vm.StringObject")},
	}

	for _, tt := range testsFail {
		evaluated := testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", err.Message, tt.expected.Message)
		}
	}
}

func generateArray(length int) *ArrayObject {
	var elements []Object
	for i := 1; i <= length; i++ {
		int := initilaizeInteger(i)
		elements = append(elements, int)
	}
	return initializeArray(elements)
}
