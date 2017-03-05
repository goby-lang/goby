package initializer_test

import (
	"github.com/st0012/Rooby/initializer"
	"github.com/st0012/Rooby/object"
	"testing"
)

func TestLengthMethod(t *testing.T) {
	expected := 5
	array := generateArray(expected)
	m := getBuiltInMethod(t, array, "length")

	result := m().(*object.IntegerObject).Value

	if int(result) != expected {
		t.Fatalf("Expect length method returns array's length: %d. got=%d", expected, result)
	}
}

func TestPopMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "pop")
	last := m().(*object.IntegerObject).Value

	if int(last) != 5 {
		t.Fatalf("Expect pop to return array's last object. got=%d", last)
	}

	if array.Length() != 4 {
		t.Fatalf("Expect pop remove last elements from array. got=%d", array.Length())
	}
}

func TestPushMethod(t *testing.T) {
	array := generateArray(5)
	m := getBuiltInMethod(t, array, "push")

	six := initializer.InitilaizeInteger(6)
	seven := initializer.InitilaizeInteger(7)
	m(six, seven)

	if array.Length() != 7 {
		t.Fatalf("Expect array's length to be 7(5 + 2). got=%d", array.Length())
	}

	last := array.Elements[array.Length() - 1].(*object.IntegerObject).Value

	if int(last) != 7 {
		t.Fatalf("Expect last object to be 7. got=%d", last)
	}
}

func generateArray(length int) *object.ArrayObject {
	var elements []object.Object
	for i := 1; i <= length; i++ {
		int := initializer.InitilaizeInteger(i)
		elements = append(elements, int)
	}
	return initializer.InitializeArray(elements)
}
