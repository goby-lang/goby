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
		t.Errorf("Expect object to be an instance of GoMap. got: %s", evaluated.ToString())
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
		t.Fatalf("Expect object to be an instance of GoMap. got: %s", evaluated.ToString())
	}

	bar, ok := m.data["foo"]

	if !ok {
		t.Fatal("Expect object's data to contains \"foo\" key")
	}

	if bar.(string) != "bar" {
		t.Fatal("Expect \"foo\" key has Go's \"bar\" string")
	}
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGoMapGetMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`m = GoMap.new;m.get("foo", 1)`, "ArgumentError: Expect 1 argument(s). got: 2", 1},
		{`m = GoMap.new;m.get`, "ArgumentError: Expect 1 argument(s). got: 0", 1},
		{`m = GoMap.new;m.get(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGoMapSetMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`m = GoMap.new;m.set("foo")`, "ArgumentError: Expect 2 argument(s). got: 1", 1},
		{`m = GoMap.new;m.set("foo", "bar", "baz")`, "ArgumentError: Expect 2 argument(s). got: 3", 1},
		{`m = GoMap.new;m.set(1, "foo")`, "TypeError: Expect argument #1 to be String. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestGoMapToHashMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`m = GoMap.new;m.to_hash(1)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`m = GoMap.new;m.to_hash(1, 2)`, "ArgumentError: Expect 0 argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
