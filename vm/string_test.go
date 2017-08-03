package vm

import (
	"testing"
)

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"st0012"`, "st0012"},
		{`'Monkey'`, "Monkey"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"string".to_s`, "string"},
		{`"123".to_i`, 123},
		{`"string".to_i`, 0},
		{`"123string123".to_i`, 123},
		{`"string123".to_i`, 0},
		{`
		  arr = "Goby".to_a
		  arr[0]
		`, "G"},
		{`
		  arr = "Goby".to_a
		  arr[1]
		`, "o"},
		{`
		  arr = "Goby".to_a
		  arr[2]
		`, "b"},
		{`
		  arr = "Goby".to_a
		  arr[3]
		`, "y"},
		{`
		  arr = "🍣Goby🍺".to_a
		  arr[0]
		`, "🍣"},
		{`
		  arr = "🍣Goby🍺".to_a
		  arr[5]
		`, "🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"123" == "123"`, true},
		{`"123" == "124"`, false},
		{`"123" == 123`, false},
		{`"123" == (1..3)`, false},
		{`"123" == { a: 1, b: 2 }`, false},
		{`"123" == [1, "String", true, 2..5]`, false},
		{`"123" != "123"`, false},
		{`"123" != "124"`, true},
		{`"123" != 123`, true},
		{`"123" != (1..3)`, true},
		{`"123" != { a: 1, b: 2 }`, true},
		{`"123" != [1, "String", true, 2..5]`, true},
		{`"123" != String`, true},
		{`"1234" > "123"`, true},
		{`"123" > "1235"`, false},
		{`"1234" < "123"`, false},
		{`"1234" < "12jdkfj3"`, true},
		{`"1234" <=> "1234"`, 0},
		{`"1234" <=> "4"`, -1},
		{`"abcdef" <=> "abcde"`, 1},
		{`"一" <=> "一"`, 0},
		{`"二" <=> "一"`, 1},
		{`"一" <=> "二"`, -1},
		{`"🍣" <=> "🍣"`, 0},
		{`"🍣" <=> "一"`, 1},
		{`"一" <=> "🍣"`, -1},
		{`"🍺" <=> "🍣"`, 1},
		{`"🍣" <=> "🍺"`, -1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringConparisonFail(t *testing.T) {
	testsFail := []struct {
		input  string
		errMsg string
	}{
		{`"a" < 1`, "TypeError: Expect argument to be String. got: Integer"},
		{`"a" > 1`, "TypeError: Expect argument to be String. got: Integer"},
		{`"a" <=> 1`, "TypeError: Expect argument to be String. got: Integer"},
	}
	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, TypeError, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Stan " + "Lo"`, "Stan Lo"},
		{`"Dog" + "&" + "Cat"`, "Dog&Cat"},
		{`"Three " * 3`, "Three Three Three "},
		{`"Zero" * 0`, ""},
		{`"Minus" * 1`, "Minus"},
		{`"Hello"[1]`, "e"},
		{`"Hello"[5]`, nil},
		{`"Hello"[-1]`, "o"},
		{`"Hello"[-6]`, nil},
		{`"Hello🍣"[5]`, "🍣"},
		{`"Hello🍣"[-1]`, "🍣"},
		{`"Hello\nWorld"[5]`, "\n"},
		{`"Ruby"[1] = "oo"`, "Rooby"},
		{`"Go"[2] = "by"`, "Goby"},
		{`"Ruby"[-3] = "oo"`, "Rooby"},
		{`"Hello"[-5] = "Tr"`, "Trello"},
		{`"Hello\nWorld"[5] = " "`, "Hello World"},
		{`"Hello🍣"[5] = "🍺"`, "Hello🍺"},
		{`"Hello🍣"[1] = "🍺"`, "H🍺llo🍣"},
		{`"Hello🍣"[-1] = "🍺"`, "Hello🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringOperationFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei" + 101`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Taipei" * "101"`, TypeError, "TypeError: Expect argument to be Integer. got: String"},
		{`"Taipei" * (-101)`, ArgumentError, "ArgumentError: Second argument must be greater than or equal to 0. got=-101"},
		{`"Taipei"[1] = 1`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Taipei"[1] = true`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Taipei"[]`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Taipei"[true] = 101`, TypeError, "TypeError: Expect argument to be Integer. got: Boolean"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringCapitalizeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"cat".capitalize`, "Cat"},
		{`"HELLO".capitalize`, "Hello"},
		{`"123word".capitalize`, "123word"},
		{`"Two Words".capitalize`, "Two words"},
		{`"first Lower".capitalize`, "First lower"},
		{`"all lower".capitalize`, "All lower"},
		{`"heLlo\nWoRLd".capitalize`, "Hello\nworld"},
		{`"🍣HeLlO🍺".capitalize`, "🍣hello🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringChopMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".chop`, "Hell"},
		{`"Hello\n".chop`, "Hello"},
		{`"Hello🍣".chop`, "Hello"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringConcatenateMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello ".concat("World")`, "Hello World"},
		{`"Hello World".concat("🍣")`, "Hello World🍣"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringConcatenateMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"a".concat`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"a".concat("Hello", "World")`, ArgumentError, "ArgumentError: Expect 1 argument. got=2"},
		{`"a".concat(1)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"a".concat(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"a".concat(nil)`, TypeError, "TypeError: Expect argument to be String. got: Null"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringCountMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"abcde".count`, 5},
		{`"哈囉！世界！".count`, 6},
		{`"Hello\nWorld".count`, 11},
		{`"Hello\nWorld🍣".count`, 12},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringDeleteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello hello HeLlo".delete("el")`, "Hlo hlo HeLlo"},
		{`"Hello 🍣 Hello 🍣 Hello".delete("🍣")`, "Hello  Hello  Hello"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringDeleteMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Hello hello HeLlo".delete`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Hello hello HeLlo".delete(1)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Hello hello HeLlo".delete(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Hello hello HeLlo".delete(nil)`, TypeError, "TypeError: Expect argument to be String. got: Null"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringDowncaseMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hEllO".downcase`, "hello"},
		{`"MORE wOrds".downcase`, "more words"},
		{`"HeLlO\tWorLD".downcase`, "hello\tworld"},
		{`"🍣HeLlO🍺".downcase`, "🍣hello🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEndWithMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".end_with("llo")`, true},
		{`"Hello".end_with("Hello")`, true},
		{`"Hello".end_with("Hello ")`, false},
		{`"哈囉！世界！".end_with("世界！")`, true},
		{`"Hello".end_with("ell")`, false},
		{`"哈囉！世界！".end_with("哈囉！")`, false},
		{`"🍣Hello🍺".end_with("🍺")`, true},
		{`"🍣Hello🍺".end_with("🍣")`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEndWithMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei".end_with("1", "0", "1")`, ArgumentError, "ArgumentError: Expect 1 argument. got=3"},
		{`"Taipei".end_with(101)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Hello".end_with(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Hello".end_with(1..5)`, TypeError, "TypeError: Expect argument to be String. got: Range"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"".empty`, true},
		{`"Hello".empty`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".eql("Hello")`, true},
		{`"Hello\nWorld".eql("Hello\nWorld")`, true},
		{`"Hello".eql("World")`, false},
		{`"Hello".eql(1)`, false},
		{`"Hello".eql(true)`, false},
		{`"Hello".eql(2..4)`, false},
		{`"Hello🍣".eql("Hello🍣")`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEqualMethodFail(t *testing.T) {
	testsFail := []struct {
		input  string
		errMsg string
	}{
		{`"Hello".eql`, "ArgumentError: Expect 1 argument. got=0"},
		{`"Hello".eql("Hello", "World")`, "ArgumentError: Expect 1 argument. got=2"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, ArgumentError, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringGlobalSubstituteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Ruby".gsub("Ru", "Go")`, "Goby"},
		{`"Hello World".gsub(" ", "\n")`, "Hello\nWorld"},
		{`"Hello World".gsub("Hello", "Goby")`, "Goby World"},
		{`"Hello 🍣 Hello 🍣 Hello".gsub("🍣", "🍺")`, "Hello 🍺 Hello 🍺 Hello"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringGlobalSubstituteMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Ruby".gsub()`, ArgumentError, "ArgumentError: Expect 2 arguments. got=0"},
		{`"Ruby".gsub("Ru")`, ArgumentError, "ArgumentError: Expect 2 arguments. got=1"},
		{`"Ruby".gsub(123, "Go")`, TypeError, "TypeError: Expect pattern to be String. got: Integer"},
		{`"Ruby".gsub("Ru", 456)`, TypeError, "TypeError: Expect replacement to be String. got: Integer"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringIncludeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello\nWorld".include("Hello")`, true},
		{`"Hello\nWorld".include("Hello\nWorld")`, true},
		{`"Hello\nWorld".include("Hello World")`, false},
		{`"Hello\nWorld".include("Hello\nWorld!")`, false},
		{`"Hello\nWorld".include("\n")`, true},
		{`"Hello\nWorld".include("\r")`, false},
		{`"Hello🍣".include("🍣")`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringIncludeMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Goby".include`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Goby".include("Ruby", "Lang")`, ArgumentError, "ArgumentError: Expect 1 argument. got=2"},
		{`"Goby".include(2)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Goby".include(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Goby".include(nil)`, TypeError, "TypeError: Expect argument to be String. got: Null"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringInsertMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".insert(0, "X")`, "XHello"},
		{`"Hello".insert(2, "X")`, "HeXllo"},
		{`"Hello".insert(5, "X")`, "HelloX"},
		{`"Hello".insert(-2, "X")`, "HelXlo"},
		{`"Hello".insert(-6, "X")`, "XHello"},
		{`"Hello".insert(0, "🍣")`, "🍣Hello"},
		{`"Hello".insert(2, "🍣")`, "He🍣llo"},
		{`"Hello".insert(5, "🍣")`, "Hello🍣"},
		{`"Hello".insert(-2, "🍣")`, "Hel🍣lo"},
		{`"Hello".insert(-6, "🍣")`, "🍣Hello"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringInsertMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Goby Lang".insert`, ArgumentError, "ArgumentError: Expect 2 arguments. got=0"},
		{`"Taipei".insert(6, " ", "101")`, ArgumentError, "ArgumentError: Expect 2 arguments. got=3"},
		{`"Taipei".insert("6", " 101")`, TypeError, "TypeError: Expect argument to be Integer. got: String"},
		{`"Taipei".insert(6, 101)`, TypeError, "TypeError: Expect insert string to be String. got: Integer"},
		{`"Taipei".insert(-8, "101")`, ArgumentError, "ArgumentError: Index value out of range. got=-8"},
		{`"Taipei".insert(7, "101")`, ArgumentError, "ArgumentError: Index value out of range. got=7"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringLeftJustifyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".ljust(2)`, "Hello"},
		{`"Hello".ljust(7)`, "Hello  "},
		{`"Hello".ljust(10, "xo")`, "Helloxoxox"},
		{`"Hello".ljust(10, "🍣🍺")`, "Hello🍣🍺🍣🍺🍣"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringLeftJustifyMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Hello".ljust`, ArgumentError, "ArgumentError: Expect 1..2 arguments. got=0"},
		{`"Hello".ljust(1, 2, 3, 4, 5)`, ArgumentError, "ArgumentError: Expect 1..2 arguments. got=5"},
		{`"Hello".ljust(true)`, TypeError, "TypeError: Expect justify width to be Integer. got: Boolean"},
		{`"Hello".ljust("World")`, TypeError, "TypeError: Expect justify width to be Integer. got: String"},
		{`"Hello".ljust(2..5)`, TypeError, "TypeError: Expect justify width to be Integer. got: Range"},
		{`"Hello".ljust(10, 10)`, TypeError, "TypeError: Expect padding string to be String. got: Integer"},
		{`"Hello".ljust(10, 2..5)`, TypeError, "TypeError: Expect padding string to be String. got: Range"},
		{`"Hello".ljust(10, true)`, TypeError, "TypeError: Expect padding string to be String. got: Boolean"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringLengthMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"New method".length`, 10},
		{`" ".length`, 1},
		{`"🍣🍣🍣".length`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".replace("World")`, "World"},
		{`"您好".replace("再見")`, "再見"},
		{`"Ruby\nLang".replace("Goby\nLang")`, "Goby\nLang"},
		{`"Hello🍣".replace("World🍺")`, "World🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei".replace`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Taipei".replace(101)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Taipei".replace(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringReverseMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Reverse Rooby-lang".reverse`, "gnal-ybooR esreveR"},
		{`" ".reverse`, " "},
		{`"-123".reverse`, "321-"},
		{`"Hello\nWorld".reverse`, "dlroW\nolleH"},
		{`"Hello 🍣🍺 World".reverse`, "dlroW 🍺🍣 olleH"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringRightJustifyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".rjust(2)`, "Hello"},
		{`"Hello".rjust(7)`, "  Hello"},
		{`"Hello".rjust(10, "xo")`, "xoxoxHello"},
		{`"Hello".rjust(10, "🍣🍺")`, "🍣🍺🍣🍺🍣Hello"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringRightJustifyFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Hello".rjust`, ArgumentError, "ArgumentError: Expect 1..2 arguments. got=0"},
		{`"Hello".rjust(1, 2, 3, 4, 5)`, ArgumentError, "ArgumentError: Expect 1..2 arguments. got=5"},
		{`"Hello".rjust(true)`, TypeError, "TypeError: Expect justify width to be Integer. got: Boolean"},
		{`"Hello".rjust("World")`, TypeError, "TypeError: Expect justify width to be Integer. got: String"},
		{`"Hello".rjust(2..5)`, TypeError, "TypeError: Expect justify width to be Integer. got: Range"},
		{`"Hello".rjust(10, 10)`, TypeError, "TypeError: Expect padding string to be String. got: Integer"},
		{`"Hello".rjust(10, 2..5)`, TypeError, "TypeError: Expect padding string to be String. got: Range"},
		{`"Hello".rjust(10, true)`, TypeError, "TypeError: Expect padding string to be String. got: Boolean"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringSizeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Rooby".size`, 5},
		{`" ".size`, 1},
		{`"🍣🍺🍺🍣".size`, 4},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringSliceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello World".slice(1..6)`, "ello W"},
		{`"1234567890".slice(6..1)`, ""},
		{`"1234567890".slice(11..1)`, nil},
		{`"1234567890".slice(11..-1)`, nil},
		{`"1234567890".slice(-10..1)`, "12"},
		{`"1234567890".slice(-5..1)`, ""},
		{`"1234567890".slice(-10..-1)`, "1234567890"},
		{`"1234567890".slice(-10..-11)`, ""},
		{`"1234567890".slice(1..-1)`, "234567890"},
		{`"1234567890".slice(1..-1234)`, ""},
		{`"1234567890".slice(-10..-5)`, "123456"},
		{`"1234567890".slice(-5..-10)`, ""},
		{`"1234567890".slice(-11..5)`, nil},
		{`"1234567890".slice(-10..-12)`, ""},
		{`"1234567890".slice(-11..-12)`, nil},
		{`"1234567890".slice(-11..-5)`, nil},
		{`"Hello 🍣🍺 World".slice(1..6)`, "ello 🍣"},
		{`"Hello 🍣🍺 World".slice(-10..7)`, "o 🍣🍺"},
		{`"Hello 🍣🍺 World".slice(1..-1)`, "ello 🍣🍺 World"},
		{`"Hello 🍣🍺 World".slice(-12..-5)`, "llo 🍣🍺 W"},
		{`"Hello World".slice(4)`, "o"},
		{`"Hello\nWorld".slice(5)`, "\n"},
		{`"Hello World".slice(-3)`, "r"},
		{`"Hello World".slice(-11)`, "H"},
		{`"Hello World".slice(-12)`, nil},
		{`"Hello World".slice(11)`, nil},
		{`"Hello 🍣🍺 World".slice(6)`, "🍣"},
		{`"Hello 🍣🍺 World".slice(-7)`, "🍺"},
		{`"Hello 🍣🍺 World".slice(-10)`, "o"},
		{`"Hello 🍣🍺 World".slice(-15)`, nil},
		{`"Hello 🍣🍺 World".slice(14)`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringSliceMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Goby Lang".slice`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Goby Lang".slice("Hello")`, TypeError, "TypeError: Expect slice range to be Range or Integer. got: String"},
		{`"Goby Lang".slice(true)`, TypeError, "TypeError: Expect slice range to be Range or Integer. got: Boolean"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringSplitMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		arr = "Hello World".split("o")
		arr[0]
		`, "Hell"},
		{`
		arr = "Hello World".split("o")
		arr[1]
		`, " W"},
		{`
		arr = "Hello World".split("o")
		arr[2]
		`, "rld"},
		{`
		arr = "Hello".split("")
		arr[0]
		`, "H"},
		{`
		arr = "Hello".split("")
		arr[1]
		`, "e"},
		{`
		arr = "Hello".split("")
		arr[2]
		`, "l"},
		{`
		arr = "Hello".split("")
		arr[3]
		`, "l"},
		{`
		arr = "Hello".split("")
		arr[4]
		`, "o"},
		{`
		arr = "Hello\nWorld\nGoby".split("\n")
		arr[0]
		`, "Hello"},
		{`
		arr = "Hello\nWorld\nGoby".split("\n")
		arr[1]
		`, "World"},
		{`
		arr = "Hello\nWorld\nGoby".split("\n")
		arr[2]
		`, "Goby"},
		{`
		arr = "Hello🍺World🍺Goby".split("🍺")
		arr[0]
		`, "Hello"},
		{`
		arr = "Hello🍺World🍺Goby".split("🍺")
		arr[1]
		`, "World"},
		{`
		arr = "Hello🍺World🍺Goby".split("🍺")
		arr[2]
		`, "Goby"},
		{`
		arr = "Hello🍺World🍣Goby".split("🍺")
		arr[0]
		`, "Hello"},
		{`
		arr = "Hello🍺World🍣Goby".split("🍺")
		arr[1]
		`, "World🍣Goby"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringSplitMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Hello World".split`, ArgumentError, "ArgumentError: Expect 1 argument. got=0"},
		{`"Hello World".split(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Hello World".split(123)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Hello World".split(1..2)`, TypeError, "TypeError: Expect argument to be String. got: Range"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringStartWithMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".start_with("Hel")`, true},
		{`"Hello".start_with("Hello")`, true},
		{`"Hello".start_with("Hello ")`, false},
		{`"哈囉！世界！".start_with("哈囉！")`, true},
		{`"Hello".start_with("hel")`, false},
		{`"哈囉！世界".start_with("世界！")`, false},
		{`"🍣Hello🍺".start_with("🍣")`, true},
		{`"🍣Hello🍺".start_with("🍺")`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringStartWithMethodFail(t *testing.T) {
	testsFail := []struct {
		input   string
		errType string
		errMsg  string
	}{
		{`"Taipei".start_with("1", "0", "1")`, ArgumentError, "ArgumentError: Expect 1 argument. got=3"},
		{`"Taipei".start_with(101)`, TypeError, "TypeError: Expect argument to be String. got: Integer"},
		{`"Hello".start_with(true)`, TypeError, "TypeError: Expect argument to be String. got: Boolean"},
		{`"Hello".start_with(1..5)`, TypeError, "TypeError: Expect argument to be String. got: Range"},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkError(t, i, evaluated, tt.errType, tt.errMsg)
		v.checkCFP(t, i, 1)
		v.checkSP(t, i, 1)
	}
}

func TestStringStripMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"  Goby Lang   ".strip`, "Goby Lang"},
		{`"\nGoby Lang\r\t".strip`, "Goby Lang"},
		{`" \t 🍣 Goby Lang 🍺 \r\n ".strip`, "🍣 Goby Lang 🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringUpcaseMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hEllO".upcase`, "HELLO"},
		{`"MORE wOrds".upcase`, "MORE WORDS"},
		{`"Hello\nWorld".upcase`, "HELLO\nWORLD"},
		{`"🍣Hello🍺".upcase`, "🍣HELLO🍺"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestChainingStringMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"More test".reverse.upcase`, "TSET EROM"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
