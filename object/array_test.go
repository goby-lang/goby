package object

import (
	"testing"
)

func init() {

}

func TestLengthMethod(t *testing.T) {
	expected := 5
	array := generateArray(expected)
	m := getBuiltInMethod(t, array, "length")

	result := m().(*IntegerObject).Value

	if int(result) != expected {
		t.Fatalf("Expect length method returns array's length: %d. got=%d", expected, result)
	}
}

func TestPopMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "pop")
	last := m().(*IntegerObject).Value

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
	m(six, seven)

	if array.Length() != 7 {
		t.Fatalf("Expect array's length to be 7(5 + 2). got=%d", array.Length())
	}

	last := array.Elements[array.Length()-1].(*IntegerObject).Value

	if int(last) != 7 {
		t.Fatalf("Expect last object to be 7. got=%d", last)
	}
}

func generateArray(length int) *ArrayObject {
	var elements []Object
	for i := 1; i <= length; i++ {
		int := InitilaizeInteger(i)
		elements = append(elements, int)
	}
	return InitializeArray(elements)
}
