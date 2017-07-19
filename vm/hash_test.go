package vm

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"
)

func TestEvalHashExpression(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

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
		case "baz":
			testBooleanObject(t, 0, value, true)
		}
	}

	vm.checkCFP(t, 0, 0)
}

func TestHashAccessOperation(t *testing.T) {
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

func TestHashAccessOperationFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }[]`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`{ a: 1, b: 2 }[true]`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`{ a: 1, b: 2 }[true] = 1`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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

func TestHashClearMethod(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }.clear
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	h, ok := evaluated.(*HashObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be a hash. got: %T", evaluated)
	} else if h.length() != 0 {
		t.Fatalf("Expect length of pairs of hash to be 0. got: %v", h.length())
	}

	vm.checkCFP(t, 0, 0)
}

func TestHashClearMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.clear(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.clear(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashEachKeyMethodFail(t *testing.T) {
	vm := initTestVM()

	testArgumentError := `
	{ a: 1, b: 2, c: 3 }.each_key("Hello") do |key|
	  puts key
	end
	`
	evaluated := vm.testEval(t, testArgumentError)
	checkError(t, 0, evaluated, ArgumentError, "ArgumentError: Expect 0 argument. got: 1")
	vm.checkCFP(t, 0, 2)

	testInternalError := `{ a: 1, b: 2, c: 3 }.each_key`
	evaluated = vm.testEval(t, testInternalError)
	checkError(t, 1, evaluated, InternalError, "InternalError: Can't yield without a block")
	vm.checkCFP(t, 1, 3)
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
	}

	for i, tt := range hashTests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
	}

	for i, tt := range normalTests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashEachValueMethodFail(t *testing.T) {
	vm := initTestVM()

	testArgumentError := `
	{ a: 1, b: 2, c: 3 }.each_value("Hello") do |value|
	  puts value
	end
	`
	evaluated := vm.testEval(t, testArgumentError)
	checkError(t, 0, evaluated, ArgumentError, "ArgumentError: Expect 0 argument. got: 1")
	vm.checkCFP(t, 0, 2)

	testInternalError := `{ a: 1, b: 2, c: 3 }.each_value`
	evaluated = vm.testEval(t, testInternalError)
	checkError(t, 1, evaluated, InternalError, "InternalError: Can't yield without a block")
	vm.checkCFP(t, 1, 3)
}

func TestHashEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{}.empty`, true},
		{`{ a: "Hello" }.empty`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashEmptyMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.empty(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.empty(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestHashEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: 1 }.eql({ a: 1 })`, true},
		{`{ a: 1 }.eql({ a: 1, b: 2 })`, false},
		{`{ a: 1, b: 2 }.eql({ a: 1, b: 2 })`, true},
		{`{ a: 1, b: 2 }.eql({ b: 2, a: 1 })`, true},
		{`{ a: 1, b: 2 }.eql({ a: 2, b: 1 })`, false},
		{`{ a: 1, b: 2 }.eql({ a: 2, b: 2, a: 1 })`, true},
		{`{ a: 1, b: 2 }.eql({ a: 1, b: 2, a: 2 })`, false},
		{`{ a: [1, 2, 3], b: { hello: "World" } }.eql({ a: [1, 2, 3], b: { hello: "World"} })`, true},
		{`{ a: [1, 2, 3], b: { hello: "World" } }.eql({ a: [3, 2, 1], b: { hello: "World"} })`, false},
		{`{ b: { hello: "World", lang: "Goby" } }.eql({ b: { lang: "Goby", hello: "World"} })`, true},
		{`{ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { hello: "World", lang: "Goby" }, range: 2..5, null: nil }.eql({ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { hello: "World", lang: "Goby" }, range: 2..5, null: nil })`, true},
		{`{ number: 1, boolean: true, string: "Goby", array: [1, "2", true], hash: { lang: "Goby", hello: "World" }, range: 2..5, null: nil }.eql({ range: 2..5, null: nil, string: "Goby", number: 1, array: [1, "2", true], boolean: true, hash: { hello: "World", lang: "Goby" } })`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashEqualMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.eql`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`{ a: 1, b: 2 }.eql(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 1 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestHashHasKeyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: "Hello", b: 123, c: true }.has_key("a")`, true},
		{`{ a: "Hello", b: 123, c: true }.has_key("d")`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashHasKeyMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.has_key`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`{ a: 1, b: 2 }.has_key(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 1 argument. got: 2"},
		{`{ a: 1, b: 2 }.has_key(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`{ a: 1, b: 2 }.has_key(123)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestHashHasValueMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{ a: "Hello", b: 123, c: true }.has_value("Hello")`, true},
		{`{ a: "Hello", b: 123, c: true }.has_value("World")`, false},
		{`{ a: "Hello", b: 123, c: true }.has_value(123)`, true},
		{`{ a: "Hello", b: 123, c: true }.has_value(false)`, false},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value({ lang: "Goby", arr: [3, 1, 2] })`, true},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value({ lang: "Goby", arr: [1, 2, 3] })`, false},
		{`{ a: "Hello", b: { lang: "Goby", arr: [3, 1, 2] }, c: true }.has_value({ arr: [3, 1, 2], lang: "Goby" })`, true},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashHasValueMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.has_value`, ArgumentError, "ArgumentError: Expect 1 argument. got: 0"},
		{`{ a: 1, b: 2 }.has_value(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 1 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}

func TestHashKeysMethod(t *testing.T) {
	input := `
	{ foo: 123, bar: "test", baz: true }.keys
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be Array. got: %T", evaluated)
	} else if arr.length() != 3 {
		t.Fatalf("Expect evaluated array length to be 3. got: %d", arr.length())
	}

	var evaluatedArr []string
	for _, k := range arr.Elements {
		evaluatedArr = append(evaluatedArr, k.(*StringObject).Value)
	}
	sort.Strings(evaluatedArr)
	if !reflect.DeepEqual(evaluatedArr, []string{"bar", "baz", "foo"}) {
		t.Fatalf("Expect evaluated array to be [\"bar\", \"baz\", \"foo\". got: %v", evaluatedArr)
	}

	vm.checkCFP(t, 0, 0)
}

func TestHashKeysMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.keys(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.keys(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashLengthMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.length(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.length(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashMapValuesMethodFail(t *testing.T) {
	vm := initTestVM()

	testArgumentError := `
	{ a: 1, b: 2, c: 3 }.map_values("Hello") do |value|
	  value * 3
	end
	`
	evaluated := vm.testEval(t, testArgumentError)
	checkError(t, 0, evaluated, ArgumentError, "ArgumentError: Expect 0 argument. got: 1")
	vm.checkCFP(t, 0, 2)

	testInternalError := `{ a: 1, b: 2, c: 3 }.map_values`
	evaluated = vm.testEval(t, testInternalError)
	checkError(t, 1, evaluated, InternalError, "InternalError: Can't yield without a block")
	vm.checkCFP(t, 1, 3)
}

func TestHashMergeMethod(t *testing.T) {
	input := []string{
		`{ a: "Hello", b: 2..5 }.merge({ b: true, c: 123, d: ["World", 456, false] })`,
		`{ b: 123, d: false }.merge({ a: "Hello", c: 123 }, { b: true, d: ["World"] }, { d: ["World", 456, false] })`,
	}

	for i, v := range input {
		vm := initTestVM()
		evaluated := vm.testEval(t, v)

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

		vm.checkCFP(t, i, 0)
	}
}

func TestHashMergeMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.merge`, ArgumentError, "ArgumentError: Expect at least 1 argument. got: 0"},
		{`{ a: 1, b: 2 }.merge(true, { hello: "World" })`, TypeError, "TypeError: Expect argument to be Hash. got: Boolean"},
		{`{ a: 1, b: 2 }.merge({ hello: "World" }, 123, "Hello")`, TypeError, "TypeError: Expect argument to be Hash. got: Integer"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashSortedKeysMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.sorted_keys(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.sorted_keys(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}

	input := `
	{ a: 123, b: "test", c: true, d: [1, "Goby", false] }.to_a
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

	arr, ok := evaluated.(*ArrayObject)
	if !ok {
		t.Fatalf("Expect evaluated value to be Array. got: %T", evaluated)
	} else if arr.length() != 4 {
		t.Fatalf("Expect evaluated array length to be 4. got: %d", arr.length())
	}

	evaluatedArr := make(map[string]Object)
	for _, p := range arr.Elements {
		pair := p.(*ArrayObject)
		evaluatedArr[pair.Elements[0].(*StringObject).Value] = pair.Elements[1]
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
	vm.checkCFP(t, 0, 0)
}

func TestHashToArrayMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.to_a(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0..1 argument. got: 2"},
		{`{ a: 1, b: 2 }.to_a(123)`, TypeError, "TypeError: Expect argument to be Boolean. got: Integer"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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

func TestHashToJSONMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.to_json(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.to_json(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashToStringMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.to_s(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.to_s(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestHashTransformValuesMethodFail(t *testing.T) {
	vm := initTestVM()

	testArgumentError := `
	{ a: 1, b: 2, c: 3 }.transform_values("Hello") do |value|
	  value * 3
	end
	`
	evaluated := vm.testEval(t, testArgumentError)
	checkError(t, 0, evaluated, ArgumentError, "ArgumentError: Expect 0 argument. got: 1")
	vm.checkCFP(t, 0, 2)

	testInternalError := `{ a: 1, b: 2, c: 3 }.transform_values`
	evaluated = vm.testEval(t, testInternalError)
	checkError(t, 1, evaluated, InternalError, "InternalError: Can't yield without a block")
	vm.checkCFP(t, 1, 3)
}

func TestHashValuesMethod(t *testing.T) {
	input := `
	{ a: 123, b: "test", c: true, d: [1, "Goby", false] }.values
	`

	vm := initTestVM()
	evaluated := vm.testEval(t, input)

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
	vm.checkCFP(t, 0, 0)
}

func TestHashValuesMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`{ a: 1, b: 2 }.values(123)`, ArgumentError, "ArgumentError: Expect 0 argument. got: 1"},
		{`{ a: 1, b: 2 }.values(true, { hello: "World" })`, ArgumentError, "ArgumentError: Expect 0 argument. got: 2"},
	}

	for i, tt := range testsFail {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		vm.checkCFP(t, i, 1)
	}
}
