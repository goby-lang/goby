package vm

import (
	"testing"
)

func TestConcurrentHashClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.class.name
		`, "Class"},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.superclass.name
		`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashClassNew(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new
		`, map[string]interface{}{}},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({a: 1, b: 2})
		`, map[string]interface{}{"a": 1, "b": 2}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testConcurrentHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashClassNewFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new(true)
		`, "TypeError: Expect argument to be Hash. got: Boolean", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new(1, 2)
		`, "ArgumentError: Expect 0 or 1 arguments, got 2", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestEvalConcurrentHashExpression(t *testing.T) {
	input := `
	require 'concurrent/hash'
	Concurrent::Hash.new({ foo: 123, bar: "test", Baz: true })
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	h, ok := evaluated.(*ConcurrentHashObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be a concurrent hash. got: %T", evaluated)
	}

	iterator := func(key, value interface{}) bool {
		switch key {
		case "foo":
			testIntegerObject(t, 0, value.(Object), 123)
		case "bar":
			testStringObject(t, 0, value.(Object), "test")
		case "Baz":
			testBooleanObject(t, 0, value.(Object), true)
		}

		return true
	}

	h.internalMap.Range(iterator)

	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestConcurrentHashAccessOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({})[:foo]
		`, nil},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({})[:foo123]
		`, nil},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo123: 100 })[:foo123]
		`, 100},
		{`
		require 'concurrent/hash'
		{}["foo"]
		`, nil},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ bar: "foo" })[:bar]
		`, "foo"},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ bar: "foo" })["bar"]
		`, "foo"},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo: 2, bar: "foo" })[:foo]
		`, 2},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo: 2, bar: "foo" })["foo"]
		`, 2},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ bar: "Foo" })
		h["bar"]
		`, "Foo"},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ bar: 1, foo: 2 })
		h["foo"] = h["bar"]
		h["foo"]

		`, 1},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({})
		h["foo"] = 100
		h["foo"]
		`, 100},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({})
		h["foo"] = Concurrent::Hash.new({ bar: 100 })
		h["foo"]["bar"]
		`, 100},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ foo: { bar: [1, 2, 3] }})
		h["foo"]["bar"][0] + h["foo"]["bar"][1]
		`, 3},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({})
		h["foo"] = 100
		h["bar"]
		`, nil},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ foo: 1, bar: 5, baz: 10 })
		h["foo"] = h["bar"] * h["baz"]
		h["foo"]
		`, 50},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashAccessOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 })[]`, "ArgumentError: Expect 1 argument. got: 0", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 })[true]`, "TypeError: Expect argument to be String. got: Boolean", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 })[true] = 1`, "TypeError: Expect argument to be String. got: Boolean", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashDeleteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("a")
		h["a"]
		`, nil},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("a")
		h["b"]
		`, "Hello"},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("a")
		h["c"]
		`, true},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("b")
		h["a"]
		`, 1},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("b")
		h["b"]
		`, nil},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("b")
		h["c"]
		`, true},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("c")
		h["a"]
		`, 1},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("c")
		h["b"]
		`, "Hello"},
		{`
		require 'concurrent/hash'
		h = Concurrent::Hash.new({ a: 1, b: "Hello", c: true })
		h.delete("c")
		h["c"]
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

func TestConcurrentHashDeleteMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: "Hello", c: true }).delete`, "ArgumentError: Expect 1 argument. got: 0", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: "Hello", c: true }).delete("a", "b")`, "ArgumentError: Expect 1 argument. got: 2", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: "Hello", c: true }).delete(123)`, "TypeError: Expect argument to be String. got: Integer", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: "Hello", c: true }).delete(true)`, "TypeError: Expect argument to be String. got: Boolean", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		// return value
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ b: "2" }).each do end
		`, map[string]interface{}{"b": "2"}},
		// empty hash
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ }).each do end
		`, map[string]interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testConcurrentHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	tests2 := []struct {
		input    string
		expected [][]interface{}
	}{
		// block yielding
		{`
		require 'concurrent/hash'
		output = []
		h = Concurrent::Hash.new({ b: "2" })
		h.each do |k, v|
			output.push([k, v])
		end
		output
		`, [][]interface{}{{"b", "2"}}},
	}

	for i, tt := range tests2 {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testBidimensionalArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashEachMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2}).each("Hello") do end`, "ArgumentError: Expect 0 arguments. got: 1", 3, 1},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2}).each`, "InternalError: Can't yield without a block", 3, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashHasKeyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: "Hello", b: 123, c: true }).has_key?("a")`, true},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: "Hello", b: 123, c: true }).has_key?("d")`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashHasKeyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).has_key?`, "ArgumentError: Expect 1 argument. got: 0", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).has_key?(true, { hello: "World" })`, "ArgumentError: Expect 1 argument. got: 2", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).has_key?(true)`, "TypeError: Expect argument to be String. got: Boolean", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).has_key?(123)`, "TypeError: Expect argument to be String. got: Integer", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToJSONMethodWithArray(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: [1, "2", true]}).to_json
		`, struct {
			A int           `json:"a"`
			B []interface{} `json:"b"`
		}{
			A: 1,
			B: []interface{}{1, "2", true},
		}},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: [1, "2", [4, 5, nil], { foo: "bar" }]}).to_json
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToJSONMethodWithNestedHash(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: { c: 2 }}).to_json
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
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: { c: 2, d: { e: "foo" }}}).to_json
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToJSONMethodWithBasicTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new.to_json
		`, struct{}{}},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).to_json
		`, struct {
			A int `json:"a"`
			B int `json:"b"`
		}{
			1,
			2,
		}},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo: "bar", b: 2 }).to_json
		`, struct {
			Foo string `json:"foo"`
			B   int    `json:"b"`
		}{
			"bar",
			2,
		}},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo: "bar", b: 2, boolean: true }).to_json
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
		require 'concurrent/hash'
		Concurrent::Hash.new({ foo: "bar", b: 2, boolean: true, nothing: nil }).to_json
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToJSONMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).to_json(123)`, "ArgumentError: Expect 0 argument. got: 1", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).to_json(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToStringMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1 }).to_s`, "{ a: 1 }"},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ b: "Hello" }).to_s`, "{ b: \"Hello\" }"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestConcurrentHashToStringMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).to_s(123)`, "ArgumentError: Expect 0 argument. got: 1", 3, 3},
		{`
		require 'concurrent/hash'
		Concurrent::Hash.new({ a: 1, b: 2 }).to_s(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 3, 3},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}
