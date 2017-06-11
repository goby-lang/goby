package vm

import (
	"encoding/json"
	"testing"
	"reflect"
)

func TestHashToJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		{ a: 1, b: 2 }.to_json
		`, struct {
			A int `json:"a"`
			B int `json:"b"`
		}{
			A: 1,
			B: 2,
		}},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		compareJSONResult(t, evaluated, tt.expected)
	}
}

func TestHashLength(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		{ a: 1, b: 2 }.length
		`, 2},
		{`
		{}.length
		`, 0},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func TestEvalHashExpression(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }
	`

	evaluated := testEval(t, input)

	h, ok := evaluated.(*HashObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be a hash. got=%T", evaluated)
	}

	for key, value := range h.Pairs {
		switch key {
		case "foo":
			testIntegerObject(t, value, 123)
		case "bar":
			testStringObject(t, value, "test")
		case "baz":
			testBooleanObject(t, value, true)
		}
	}
}

func TestEvalHashAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			{}["foo"]
		`, nil},
		{`
			{ bar: "foo" }["bar"]
		`, "foo"},
		{`
			{ foo: 2, bar: "foo" }["foo"]
		`, 2},
		{`
			h = { bar: "Foo" }
			h["bar"]
		`, "Foo"},
		{`
			h = { bar: 1, foo: 2 }
			h["foo"] = h["bar"]
			h["foo"]

		`, 1},
		{`
			h = {}
			h["foo"] = 100
			h["foo"]
		`, 100},
		{`
			h = {}
			h["foo"] = 100
			h["bar"]
		`, nil},
		{`
			h = { foo: 1, bar: 5, baz: 10 }
			h["foo"] = h["bar"] * h["baz"]
			h["foo"]
		`, 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func compareJSONResult(t *testing.T, evaluated Object, exp interface{}) {
	expected, err := json.Marshal(exp)

	if err != nil {
		t.Fatal(err.Error())
	}

	s := evaluated.(*StringObject).Value

	r, err := JSONBytesEqual([]byte(s), expected)

	if err != nil {
		t.Fatal(err.Error())
	}

	if !r {
		t.Fatalf("Expect json:\n%s \n\n got: %s", string(expected), s)
	}
}