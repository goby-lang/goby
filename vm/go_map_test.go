package vm

import (
	"testing"
)

func TestGoMapInitWithoutArg(t *testing.T) {
	input := `
	GoMap.new
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	_, ok := evaluated.(*GoMap)

	if !ok {
		t.Errorf("Expect object to be an instance of GoMap. got: %s", evaluated.toString())
	}
}

func TestGoMapInitWithHash(t *testing.T) {
	input := `
	h = { foo: "bar" }
	GoMap.new(h)
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	m, ok := evaluated.(*GoMap)

	if !ok {
		t.Fatalf("Expect object to be an instance of GoMap. got: %s", evaluated.toString())
	}

	bar, ok := m.data["foo"]

	if !ok {
		t.Fatal("Expect object's data to contains \"foo\" key")
	}

	b := bar.(*StringObject)

	testStringObject(t, 0, b, "bar")
}
