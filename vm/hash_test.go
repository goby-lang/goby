package vm

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func TestHashClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Hash.class.name`, "Class"},
		{`Hash.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestEvalHashExpression(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", Baz: true }
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	h, ok := evaluated.(*HashObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be a hash. got: %T", evaluated)
	}

	for key, value := range h.Pairs {
		switch key {
		case "foo":
			testIntegerObject(t, 0, value, 123)
		case "bar":
			testStringObject(t, 0, value, "test")
		case "Baz":
			testBooleanObject(t, 0, value, true)
		}
	}

	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestHashAccessOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			{}[:foo]
		`, nil},
		{`
			{}[:foo123]
		`, nil},
		{`
			{ foo123: 100 }[:foo123]
		`, 100},
		{`
			{}["foo"]
		`, nil},
		{`
			{ bar: "foo" }[:bar]
		`, "foo"},
		{`
			{ bar: "foo" }["bar"]
		`, "foo"},
		{`
			{ foo: 2, bar: "foo" }[:foo]
		`, 2},
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashAccessWithDefaultOperation(t *testing.T) {
	valueTests := []struct {
		input    string
		expected interface{}
	}{
		{`
			h = {}
			h.default = 0
			h['c']
		`, 0},
		{`
			h = {}
			h.default = 0
			h['d'] += 2
			h['d']
		`, 2},
	}

	for i, tt := range valueTests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	hashTests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{`
			h = {}
			h.default = 0
			h
		`, map[string]interface{}{}},
		{`
			h = {}
			h.default = 0
			h['d'] += 2
			h
		`, map[string]interface{}{"d": 2}},
	}

	for i, tt := range hashTests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashAccessOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }[]`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`{ a: 1, b: 2 }[true]`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`{ a: 1, b: 2 }[true] = 1`, "TypeError: Expect argument to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashAnyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
      { a: 1, b: 2 }.any? do |k, v|
        v == 2
      end
		`, true},
		{`
      { a: 1, b: 2 }.any? do |k, v|
        v
      end
		`, true},
		{`
      { a: 1, b: 2 }.any? do |k, v|
        v == 5
      end
		`, false},
		{`
      { a: 1, b: 2 }.any? do |k, v|
        nil
      end
		`, false},
		{`
      { }.any? do |k, v|
        true
      end
		`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashAnyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{  }.any?(123) do end`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{  }.any?`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashClearMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		// object modification
		{`
			hash = { foo: 123, bar: "test" }
			hash.clear
			hash
		`, map[string]interface{}{}},

		// return value
		{`
			{ foo: 123, bar: "test" }.clear
		`, map[string]interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashClearMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.clear(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.clear(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashDefaultOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			h = {}
			h.default
		`, nil},
		{`
			h = {}
			h.default = 0
			h.default
		`, 0},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashDefaultSetOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ }.default = *[1, 2]`, "ArgumentError: Expected 1 argument, got 2", 1},
		{`{ }.default = []`, "ArgumentError: Arrays and Hashes are not accepted as default values", 1},
		{`{ }.default = {}`, "ArgumentError: Arrays and Hashes are not accepted as default values", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashDigMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
			{ a: 1, b: 2 }.dig(:a)
		`, 1},
		{`
			{ a: {}, b: 2 }.dig(:a, :b)
		`, nil},
		{`
			{ a: {}, b: 2 }.dig(:a, :b, :c)
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

func TestHashDigMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: [], b: 2 }.dig`, "ArgumentError: Expected 1+ arguments, got 0", 1},
		{`{ a: 1, b: 2 }.dig(:a, :b)`, "TypeError: Expect target to be Diggable, got Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		// return value
		{`
			{ b: "2", a: 1 }.each do end
		`, map[string]interface{}{"a": 1, "b": "2"}},
		// empty hash
		{`
			{ }.each do end
		`, map[string]interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	tests2 := []struct {
		input    string
		expected [][]interface{}
	}{
		// block yielding
		{`
			output = []
			h = { b: "2", a: 1 }
			h.each do |k, v|
				output.push([k, v])
			end
			output
		`, [][]interface{}{{"a", 1}, {"b", "2"}}},
	}

	for i, tt := range tests2 {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testBidimensionalArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2}.each("Hello") do end
		`, "ArgumentError: Expect 0 arguments. got: 1", 1},
		{`{ a: 1, b: 2}.each`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachKeyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
			{ a: "Hello", b: "World", c: "Goby" }.each_key do |key|
			  # Empty Block
			end
		`, []interface{}{"a", "b", "c"}},
		{`
			{ b: "Hello", c: "World", a: "Goby" }.each_key do |key|
			  # Empty Block
			end
		`, []interface{}{"a", "b", "c"}},
		{`
			{ b: "Hello", c: "World", b: "Goby" }.each_key do |key|
			  # Empty Block
			end
		`, []interface{}{"b", "c"}},
		{`
			arr = []
			{ a: "Hello", b: "World", c: "Goby" }.each_key do |key|
			  arr.push(key)
			end
			arr
		`, []interface{}{"a", "b", "c"}},
		{`
			arr = []
			{}.each_key do |key|
			  arr.push(key)
			end
			arr
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachKeyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2, c: 3 }.each_key("Hello") do |key|
		  puts key
		end
		`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2, c: 3 }.each_key`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachValueMethod(t *testing.T) {
	hashTests := []struct {
		input    string
		expected []interface{}
	}{
		{`
			{ a: "Hello", b: 123, c: true }.each_value do |v|
			  # Empty Block
			end
		`, []interface{}{"Hello", 123, true}},
		{`
			{ b: "Hello", c: 123, a: true }.each_value do |v|
			  # Empty Block
			end
		`, []interface{}{true, "Hello", 123}},
		{`
			{ a: "Hello", b: 123, a: true }.each_value do |v|
			  # Empty Block
			end
		`, []interface{}{true, 123}},
		{`
			{}.each_value do |v|
			  # Empty Block
			end
		`, []interface{}{}},
	}

	for i, tt := range hashTests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	normalTests := []struct {
		input    string
		expected interface{}
	}{
		{`
			sum = 0
			{ a: 1, b: 2, c: 3, d: 4, e: 5 }.each_value do |v|
			  sum = sum + v
			end
			sum
			`, 15},
		{`
			sum = 0
			{ a: 1, b: 2, a: 3, b: 4, a: 5 }.each_value do |v|
			  sum = sum + v
			end
			sum
			`, 9},
		{`
			string = ""
			{ a: "Hello", b: "World", c: "Goby", d: "Lang" }.each_value do |v|
			  string = string + v + " "
			end
			string
			`, "Hello World Goby Lang "},
		{`
			string = ""
			{}.each_value do |v|
			  string = string + v + " "
			end
			string
			`, ""},
	}

	for i, tt := range normalTests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashEachValueMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2, c: 3 }.each_value("Hello") do |value|
		  puts value
		end
		`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2, c: 3 }.each_value`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{}.empty?`, true},
		{`{ a: "Hello" }.empty?`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashEmptyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.empty?(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.empty?(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: 1 }.eql?({ a: 1 })`, true},
		{`{ a: 1 }.eql?({ a: 1, b: 2 })`, false},
		{`{ a: 1, b: 2 }.eql?({ a: 1, b: 2 })`, true},
		{`{ a: 1, b: 2 }.eql?({ b: 2, a: 1 })`, true},
		{`{ a: 1, b: 2 }.eql?({ a: 2, b: 1 })`, false},
		{`{ a: 1, b: 2 }.eql?({ a: 2, b: 2, a: 1 })`, true},
		{`{ a: 1, b: 2 }.eql?({ a: 1, b: 2, a: 2 })`, false},
		{`{ a: [1, 2, 3], b: { hello: "World" } }.eql?({ a: [1, 2, 3], b: { hello: "World"} })`, true},
		{`{ a: [1, 2, 3], b: { hello: "World" } }.eql?({ a: [3, 2, 1], b: { hello: "World"} })`, false},
		{`{ b: { hello: "World", lang: "Goby" } }.eql?({ b: { lang: "Goby", hello: "World"} })`, true},
		{`{ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { hello: "World", lang: "Goby" }, range: 2..5, null: nil }.eql?({ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { hello: "World", lang: "Goby" }, range: 2..5, null: nil })`, true},
		{`{ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { lang: "Goby", hello: "World" }, range: 2..5, null: nil }.eql?({ range: 2..5, null: nil, string: "Goby", number: 1, array: [1, "2", true], boolean: true, hash: { hello: "World", lang: "Goby" } })`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashEqualMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.eql?`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`{ a: 1, b: 2 }.eql?(true, { hello: "World" })`, "ArgumentError: Expect 1 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashDeleteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = { a: 1, b: "Hello", c: true }.delete("a")
		h["a"]
		`, nil},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("a")
		h["b"]
		`, "Hello"},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("a")
		h["c"]
		`, true},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("b")
		h["a"]
		`, 1},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("b")
		h["b"]
		`, nil},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("b")
		h["c"]
		`, true},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("c")
		h["a"]
		`, 1},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("c")
		h["b"]
		`, "Hello"},
		{`
		h = { a: 1, b: "Hello", c: true }.delete("c")
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

func TestHashDeleteMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: "Hello", c: true }.delete`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`{ a: 1, b: "Hello", c: true }.delete("a", "b")`, "ArgumentError: Expect 1 argument. got: 2", 1},
		{`{ a: 1, b: "Hello", c: true }.delete(123)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`{ a: 1, b: "Hello", c: true }.delete(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashDeleteIfMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		// Since the method returns the hash itself, for compactness we perform the
		// tests on the return value, but we still make sure, with the first test,
		// that the hash itself is modified.
		{`
			hash = { a: 1, b: 2 }
			hash.delete_if do |k, v| v == 1 end
			hash
		`, map[string]interface{}{"b": 2}},
		{`
			{ a: 1, b: 2 }.delete_if do |k, v| v == 1 end
		`, map[string]interface{}{"b": 2}},
		{`
			{ a: 1, b: 2 }.delete_if do |k, v| 5 end
		`, map[string]interface{}{}},
		{`
			{ a: 1, b: 2 }.delete_if do |k, v| false end
		`, map[string]interface{}{"a": 1, "b": 2}},
		{`
			{ a: 1, b: 2 }.delete_if do |k, v| nil end
		`, map[string]interface{}{"a": 1, "b": 2}},
		{`
			{ }.delete_if do |k, v| true end
		`, map[string]interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashDeleteIfMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ }.delete_if(123) do end`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ }.delete_if`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashHasKeyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: "Hello", b: 123, c: true }.has_key?("a")`, true},
		{`{ a: "Hello", b: 123, c: true }.has_key?("d")`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashHasKeyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.has_key?`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`{ a: 1, b: 2 }.has_key?(true, { hello: "World" })`, "ArgumentError: Expect 1 argument. got: 2", 1},
		{`{ a: 1, b: 2 }.has_key?(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`{ a: 1, b: 2 }.has_key?(123)`, "TypeError: Expect argument to be String. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashHasValueMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: "Hello", b: 123, c: true }.has_value?("Hello")`, true},
		{`{ a: "Hello", b: 123, c: true }.has_value?("World")`, false},
		{`{ a: "Hello", b: 123, c: true }.has_value?(123)`, true},
		{`{ a: "Hello", b: 123, c: true }.has_value?(false)`, false},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value?({ lang: "Goby", arr: [3, 1, 2] })`, true},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value?({ lang: "Goby", arr: [1, 2, 3] })`, false},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value?({ arr: [3, 1, 2], lang: "Goby" })`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashHasValueMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.has_value?`, "ArgumentError: Expect 1 argument. got: 0", 1},
		{`{ a: 1, b: 2 }.has_value?(true, { hello: "World" })`, "ArgumentError: Expect 1 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashKeysMethod(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }.keys
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be Array. got: %T", evaluated)
	} else if arr.length() != 3 {
		t.Fatalf("Expect evaluated array length to be 3. got: %d", arr.length())
	}

	var evaluatedArr []string
	for _, k := range arr.Elements {
		evaluatedArr = append(evaluatedArr, k.(*StringObject).value)
	}
	sort.Strings(evaluatedArr)
	if !reflect.DeepEqual(evaluatedArr, []string{"bar", "baz", "foo"}) {
		t.Fatalf("Expect evaluated array to be [\"bar\", \"baz\", \"foo\". got: %v", evaluatedArr)
	}

	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestHashKeysMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.keys(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.keys(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashLengthMethod(t *testing.T) {
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashLengthMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.length(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.length(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashMapValuesMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		h["a"]
		`, 3},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		h["b"]
		`, 6},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		h["c"]
		`, 9},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		result["a"]
		`, 3},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		result["b"]
		`, 6},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.map_values do |v|
		  v * 3
		end
		result["c"]
		`, 9},
		{`
		h = {}
		result = h.map_values do |v|
		  v * 3
		end
		result["c"]
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

func TestHashMapValuesMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2, c: 3 }.map_values("Hello") do |value|
		  value * 3
		end
		`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2, c: 3 }.map_values`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashMergeMethod(t *testing.T) {
	input := []string{
		`{ a: "Hello", b: 2..5 }.merge({ b: true, c: 123, d: ["World", 456, false] })`,
		`{ b: 123, d: false }.merge({ a: "Hello", c: 123 }, { b: true, d: ["World"] }, { d: ["World", 456, false] })`,
	}

	for i, value := range input {
		v := initTestVM()
		evaluated := v.testEval(t, value, getFilename())

		h, ok := evaluated.(*HashObject)
		if !ok {
			t.Fatalf("Expect evaluated value to be a hash. got: %T", evaluated)
		}

		for key, value := range h.Pairs {
			switch key {
			case "a":
				testStringObject(t, i, value, "Hello")
			case "b":
				testBooleanObject(t, i, value, true)
			case "c":
				testIntegerObject(t, i, value, 123)
			case "d":
				testArrayObject(t, i, value, []interface{}{"World", 456, false})
			}
		}

		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashMergeMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.merge`, "ArgumentError: Expect at least 1 argument. got: 0", 1},
		{`{ a: 1, b: 2 }.merge(true, { hello: "World" })`, "TypeError: Expect argument to be Hash. got: Boolean", 1},
		{`{ a: 1, b: 2 }.merge({ hello: "World" }, 123, "Hello")`, "TypeError: Expect argument to be Hash. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashSortedKeysMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`{ a: 1, b: 2, c: 3 }.sorted_keys`, []interface{}{"a", "b", "c"}},
		{`{ c: 1, b: 2, a: 3 }.sorted_keys`, []interface{}{"a", "b", "c"}},
		{`{ b: 1, a: 2, c: 3 }.sorted_keys`, []interface{}{"a", "b", "c"}},
		{`{ b: 1, a: 2, b: 3 }.sorted_keys`, []interface{}{"a", "b"}},
		{`{ c: 1, a: 2, a: 3 }.sorted_keys`, []interface{}{"a", "c"}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashSelectMethod(t *testing.T) {
	testsSortedArray := []struct {
		input    string
		expected map[string]interface{}
	}{
		{`
			{ a: 1, b: 2 }.select do |k, v|
			  v == 2
			end
		`, map[string]interface{}{"b": 2}},
		{`
			{ a: 1, b: 2 }.select do |k, v|
			  5
			end
		`, map[string]interface{}{ "a": 1, "b": 2 }},
		{`
			{ a: 1, b: 2 }.select do |k, v|
			  nil
			end
		`, map[string]interface{}{}},
		{`
			{ a: 1, b: 2 }.select do |k, v|
			  false
			end
		`, map[string]interface{}{}},
		{`
			{ }.select do end
		`, map[string]interface{}{}},
		// non-destructivity specification
		{`
			source = { a: 1, b: 2 }
			source.select do |k, v| true end
			source
		`, map[string]interface{}{ "a": 1, "b": 2 }},
	}

	for i, tt := range testsSortedArray {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testHashObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashSelectMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ }.select(123) do end`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ }.select`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashSortedKeysMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.sorted_keys(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.sorted_keys(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashToArrayMethod(t *testing.T) {
	testsSortedArray := []struct {
		input    string
		expected []interface{}
	}{
		{`{ a: 1, b: 2, c: 3 }.to_a(true)[0]`, []interface{}{"a", 1}},
		{`{ a: 1, b: 2, c: 3 }.to_a(true)[1]`, []interface{}{"b", 2}},
		{`{ a: 1, b: 2, c: 3 }.to_a(true)[2]`, []interface{}{"c", 3}},
		{`{ b: 1, c: 2, a: 3 }.to_a(true)[0]`, []interface{}{"a", 3}},
		{`{ b: 1, c: 2, a: 3 }.to_a(true)[1]`, []interface{}{"b", 1}},
		{`{ b: 1, c: 2, a: 3 }.to_a(true)[2]`, []interface{}{"c", 2}},
	}

	for i, tt := range testsSortedArray {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}

	input := `
	{ a: 123, b: "test", c: true, d: [1, "Goby", false] }.to_a
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be Array. got: %T", evaluated)
	} else if arr.length() != 4 {
		t.Fatalf("Expect evaluated array length to be 4. got: %d", arr.length())
	}

	evaluatedArr := make(map[string]Object)
	for _, p := range arr.Elements {
		pair := p.(*ArrayObject)
		evaluatedArr[pair.Elements[0].(*StringObject).value] = pair.Elements[1]
	}

	for k, v := range evaluatedArr {
		switch k {
		case "a":
			testIntegerObject(t, 0, v, 123)
		case "b":
			testStringObject(t, 0, v, "test")
		case "c":
			testBooleanObject(t, 0, v, true)
		case "d":
			testArrayObject(t, 0, v, []interface{}{1, "Goby", false})
		}
	}
	v.checkCFP(t, 0, 0)
}

func TestHashToArrayMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.to_a(true, { hello: "World" })`, "ArgumentError: Expect 0..1 argument. got: 2", 1},
		{`{ a: 1, b: 2 }.to_a(123)`, "TypeError: Expect argument to be Boolean. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
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
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		compareJSONResult(t, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashToJSONMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.to_json(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.to_json(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashToStringMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: 1 }.to_s`, "{ a: 1 }"},
		{`{ a: 1, b: "Hello" }.to_s`, "{ a: 1, b: \"Hello\" }"},
		{`{ a: 1, b: [1, true, "Hello", 1..2], c: { lang: "Goby" } }.to_s`, "{ a: 1, b: [1, true, \"Hello\", (1..2)], c: { lang: \"Goby\" } }"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashToStringMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.to_s(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.to_s(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashTransformValuesMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		h["a"]
		`, 1},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		h["b"]
		`, 2},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		h["c"]
		`, 3},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		result["a"]
		`, 3},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		result["b"]
		`, 6},
		{`
		h = { a: 1, b: 2, c: 3 }
		result = h.transform_values do |v|
		  v * 3
		end
		result["c"]
		`, 9},
		{`
		h = {}
		result = h.transform_values do |v|
		  v * 3
		end
		result["c"]
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

func TestHashTransformValuesMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2, c: 3 }.transform_values("Hello") do |value|
		  value * 3
		end
		`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2, c: 3 }.transform_values`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashValuesMethod(t *testing.T) {
	input := `
	{ a: 123, b: "test", c: true, d: [1, "Goby", false] }.values
	`

	v := initTestVM()
	evaluated := v.testEval(t, input, getFilename())

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be Array. got: %T", evaluated)
	} else if arr.length() != 4 {
		t.Fatalf("Expect evaluated array length to be 4. got: %d", arr.length())
	}

	for _, v := range arr.Elements {
		switch value := v.(type) {
		case *IntegerObject:
			testIntegerObject(t, 0, value, 123)
		case *StringObject:
			testStringObject(t, 0, v, "test")
		case *BooleanObject:
			testBooleanObject(t, 0, v, true)
		case *ArrayObject:
			testArrayObject(t, 0, v, []interface{}{1, "Goby", false})
		}
	}
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestHashValuesMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.values(123)`, "ArgumentError: Expect 0 argument. got: 1", 1},
		{`{ a: 1, b: 2 }.values(true, { hello: "World" })`, "ArgumentError: Expect 0 argument. got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestHashValuesAtMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		{ a: 1, b: "2" }.values_at("a", "c")
		`, []interface{}{1, nil}},
		{`
		{ a: 1, b: "2" }.values_at()
		`, []interface{}{}},
		{`
		{}.values_at("a")
		`, []interface{}{nil}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestHashValuesAtMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`{ a: 1, b: 2 }.values_at(123)`, "TypeError: Expect argument to be String. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
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

	s := evaluated.(*StringObject).value

	r, err := JSONBytesEqual([]byte(s), expected)

	if err != nil {
		t.Fatal(err.Error())
	}

	if !r {
		t.Fatalf("Expect json:\n%s \n\n got: %s", string(expected), s)
	}
}
