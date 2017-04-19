package vm

import (
	"testing"
)

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
