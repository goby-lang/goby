package vm

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHashComparisonOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: 1, b: 2 } == 123`, false},
		{`{ a: 1, b: 2 } == "123"`, false},
		{`{ a: 1, b: 2 } == "124"`, false},
		{`{ a: 1, b: 2 } == (1..3)`, false},
		{`{ a: 1, b: 2 } == { a: 1, b: 2 }`, true},
		{`{ b: 2, a: 1 } == { a: 1, b: 2 }`, true}, // Hash has no order issue
		{`{ a: 1, b: 2 } == { a: 2, b: 1 }`, false},
		{`{ a: 1, b: 2 } == { b: 1, a: 2 }`, false},
		{`{ a: 1, b: 2 } == { a: 1, b: 2, c: 3 }`, false},
		{`{ a: 1, b: 2 } == { a: 2, b: 2, a: 1 }`, true}, // Hash front key will be overwritten if duplicated
		{`{ a: [1, 2, 3], b: 2 } == { a: [1, 2, 3], b: 2 }`, true},
		{`{ a: [1, 2, 3], b: 2 } == { a: [3, 2, 1], b: 2 }`, false}, // Hash of array has order issue
		{`{ a: 1, b: 2 } == [1, "String", true, 2..5]`, false},
		{`{ a: 1, b: 2 } == Integer`, false},
		{`{ a: 1, b: 2 } != 123`, true},
		{`{ a: 1, b: 2 } != "123"`, true},
		{`{ a: 1, b: 2 } != "124"`, true},
		{`{ a: 1, b: 2 } != (1..3)`, true},
		{`{ a: 1, b: 2 } != { a: 1, b: 2 }`, false},
		{`{ b: 2, a: 1 } != { a: 1, b: 2 }`, false}, // Hash has no order issue
		{`{ a: 1, b: 2 } != { a: 2, b: 1 }`, true},
		{`{ a: 1, b: 2 } != { b: 1, a: 2 }`, true},
		{`{ a: 1, b: 2 } != { a: 1, b: 2, c: 3 }`, true},
		{`{ a: 1, b: 2 } != { a: 2, b: 2, a: 1 }`, false}, // Hash front key will be overwritten if duplicated
		{`{ a: [1, 2, 3], b: 2 } != { a: [1, 2, 3], b: 2 }`, false},
		{`{ a: [1, 2, 3], b: 2 } != { a: [3, 2, 1], b: 2 }`, true}, // Hash of array has order issue
		{`{ a: 1, b: 2 } != [1, "String", true, 2..5]`, true},
		{`{ a: 1, b: 2 } != Integer`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashToJSONMethodWithArray(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		{ a: 1, b: [1, "2", true]}.to_json
		`, struct {
			A int           `json:"a"`
			B []interface{} `json:"b"`
		}{
			A: 1,
			B: []interface{}{1, "2", true},
		}},
		{`
		{ a: 1, b: [1, "2", [4, 5, nil], { foo: "bar" }]}.to_json
		`, struct {
			A int           `json:"a"`
			B []interface{} `json:"b"`
		}{
			A: 1,
			B: []interface{}{
				1, "2", []interface{}{4, 5, nil}, struct {
					Foo string `json:"foo"`
				}{
					"bar",
				},
			},
		}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		compareJSONResult(t, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashToJSONMethodWithNestedHash(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		{ a: 1, b: { c: 2 }}.to_json
		`, struct {
			A int `json:"a"`
			B struct {
				C int `json:"c"`
			} `json:"b"`
		}{
			1,
			struct {
				C int `json:"c"`
			}{C: 2},
		}},
		{`
		{ a: 1, b: { c: 2, d: { e: "foo" }}}.to_json
		`, struct {
			A int `json:"a"`
			B struct {
				C int `json:"c"`
				D struct {
					E string `json:"e"`
				} `json:"d"`
			} `json:"b"`
		}{
			1,
			struct {
				C int `json:"c"`
				D struct {
					E string `json:"e"`
				} `json:"d"`
			}{C: 2, D: struct {
				E string `json:"e"`
			}{E: "foo"}},
		}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		compareJSONResult(t, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashToJSONMethodWithBasicTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		{}.to_json
		`, struct{}{}},
		{`
		{ a: 1, b: 2 }.to_json
		`, struct {
			A int `json:"a"`
			B int `json:"b"`
		}{
			1,
			2,
		}},
		{`
		{ foo: "bar", b: 2 }.to_json
		`, struct {
			Foo string `json:"foo"`
			B   int    `json:"b"`
		}{
			"bar",
			2,
		}},
		{`
		{ foo: "bar", b: 2, boolean: true }.to_json
		`, struct {
			Foo     string `json:"foo"`
			B       int    `json:"b"`
			Boolean bool   `json:"boolean"`
		}{
			"bar",
			2,
			true,
		}},
		{`
		{ foo: "bar", b: 2, boolean: true, nothing: nil }.to_json
		`, struct {
			Foo     string      `json:"foo"`
			B       int         `json:"b"`
			Boolean bool        `json:"boolean"`
			Nothing interface{} `json:"nothing"`
		}{
			"bar",
			2,
			true,
			nil,
		}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		compareJSONResult(t, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestEvalHashExpression(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	h, ok := evaluated.(*HashObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be a hash. got=%T", evaluated)
	}

	for key, value := range h.Pairs {
		switch key {
		case "foo":
			testIntegerObject(t, 0, value, 123)
		case "bar":
			testStringObject(t, 0, value, "test")
		case "baz":
			testBooleanObject(t, 0, value, true)
		}
	}

	vm.checkCFP(t, 0, 0)
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
			h["foo"] = { bar: 100 }
			h["foo"]["bar"]
		`, 100},
		{`
			h = { foo: { bar: [1, 2, 3] }}
			h["foo"]["bar"][0] + h["foo"]["bar"][1]
		`, 3},
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

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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

// We can't compare string directly because the key/value's order might change and we can't control it.
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
