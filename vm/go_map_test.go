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

func TestGoMapGetMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = { foo: "bar" }
		m = GoMap.new(h)
		m.get("foo")
		`, "bar"},
		{`
		m = GoMap.new
		m.get("foo")
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGoMapSetMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		m = GoMap.new
		m.set("foo", "bar")
		m.get("foo")
		`, "bar"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGoMapToHashMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = { foo: "bar" }
		m = GoMap.new(h)
		h2 = m.to_hash
		h2[:foo]
		`, "bar"},
		{`
		m = GoMap.new
		h = m.to_hash
		h[:foo]
		`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
