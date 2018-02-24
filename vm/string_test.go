package vm

import (
	"testing"
)

func TestStringClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`String.class.name`, "Class"},
		{`String.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"st0012"`, "st0012"},
		{`'Monkey'`, "Monkey"},
		{`"\"Maxwell\""`, "\"Maxwell\""},
		{`"'Alexius'"`, "'Alexius'"},
		{`"\'Maxwell\'"`, "'Maxwell'"},
		{`'\'Alexius\''`, "'Alexius'"},
		{`"Maxwell\nAlexius"`, "Maxwell\nAlexius"},
		{`'Maxwell\nAlexius'`, "Maxwell\\nAlexius"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		{`"Maxwell\nAlexius".to_s`, "Maxwell\nAlexius"},
		{`'Maxwell\nAlexius'.to_s`, "Maxwell\\nAlexius"},
		{`"\"Maxwell\"".to_s`, "\"Maxwell\""},
		{`'\"Maxwell\"'.to_s`, "\\\"Maxwell\\\""},
		{`"\'Maxwell\'".to_s`, "'Maxwell'"},
		{`'\'Maxwell\''.to_s`, "'Maxwell'"},
		{`"123".to_i`, 123},
		{`"string".to_i`, 0},
		{`" \t123".to_i`, 123},
		{`"123string123".to_i`, 123},
		{`"string123".to_i`, 0},
		{`"123.5".to_f`, 123.5},
		{`".5".to_f`, 0.5},
		{`"  123.5".to_f`, 123.5},
		{`"3.5e2".to_f`, 350.0},
		{`"3.5ef".to_f`, 0.0},
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
		{`
		  arr = "Maxwell\nAlexius".to_a
		  arr[7]
		 `, "\n"},
		{`
		  arr = 'Maxwell\nAlexius'.to_a
		  arr[7]
		 `, "\\"},
		{`
		  arr = 'Maxwell\nAlexius'.to_a
		  arr[8]
		 `, "n"},
		{`
		  arr = "\"Maxwell\"".to_a
		  arr[0]
		 `, "\""},
		{`
		  arr = "\"Maxwell\"".to_a
		  arr[-1]
		 `, "\""},
		{`
		  arr = "\'Maxwell\'".to_a
		  arr[0]
		 `, "'"},
		{`
		  arr = "\'Maxwell\'".to_a
		  arr[-1]
		 `, "'"},
		{`
		  arr = '\"Maxwell\"'.to_a
		  arr[0]
		 `, "\\"},
		{`
		  arr = '\"Maxwell\"'.to_a
		  arr[1]
		 `, "\""},
		{`
		  arr = '\"Maxwell\"'.to_a
		  arr[-1]
		 `, "\""},
		{`
		  arr = '\"Maxwell\"'.to_a
		  arr[-2]
		 `, "\\"},
		{`
		  arr = '\'Maxwell\''.to_a
		  arr[0]
		 `, "'"},
		{`
		  arr = '\'Maxwell\''.to_a
		  arr[-1]
		 `, "'"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		{`"\"Maxwell\"" == '"Maxwell"'`, true},
		{`"\'Maxwell\'" == '\'Maxwell\''`, true},
		{`"\"Maxwell\"" == '\"Maxwell\"'`, false},
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringComparisonFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"a" < 1`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"a" > 1`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"a" <=> 1`, "TypeError: Expect argument to be String. got: Integer", 1},
	}
	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringMatchOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"abc" =~ Regexp.new("bc")`, 1},
		{`"abc" =~ Regexp.new("d")`, nil},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringMatchOperatorFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"abc" =~ *[1, 2]`, "ArgumentError: Expect 1 argument. got=2", 1},
		{`"abc" =~ 'a'`, "TypeError: Expect argument to be Regexp. got: String", 1},
	}
	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		{`"Hello"[1..3]`, "ell"},
		{`"Hello"[1..5]`, "ello"},
		{`"Hello"[-3..-1]`, "llo"},
		{`"Hello"[-6..-1]`, nil},
		{`"Hello🍣"[-3..-3]`, "l"},
		{`"Hello🍣"[1..-1]`, "ello🍣"},
		{`"Hello\nWorld"[5]`, "\n"},
		{`"\"Maxwell\""[0]`, "\""},
		{`"\"Maxwell\""[-1]`, "\""},
		{`"\'Maxwell\'"[0]`, "'"},
		{`"\'Maxwell\'"[-1]`, "'"},
		{`'\"Maxwell\"'[0]`, "\\"},
		{`'\"Maxwell\"'[1]`, "\""},
		{`'\"Maxwell\"'[-1]`, "\""},
		{`'\"Maxwell\"'[-2]`, "\\"},
		{`'\'Maxwell\''[0]`, "'"},
		{`'\'Maxwell\''[-1]`, "'"},
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringOperationFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Taipei" + 101`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Taipei" * "101"`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`"Taipei" * (-101)`, "ArgumentError: Second argument must be greater than or equal to 0. got=-101", 1},
		{`"Taipei"[1] = 1`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Taipei"[1] = true`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Taipei"[]`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Taipei"[true] = 101`, "TypeError: Expect argument to be Integer. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Method test

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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringConcatenateMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"a".concat`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"a".concat("Hello", "World")`, "ArgumentError: Expect 1 argument. got=2", 1},
		{`"a".concat(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"a".concat(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"a".concat(nil)`, "TypeError: Expect argument to be String. got: Null", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringDeleteMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Hello hello HeLlo".delete`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Hello hello HeLlo".delete(1)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Hello hello HeLlo".delete(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Hello hello HeLlo".delete(nil)`, "TypeError: Expect argument to be String. got: Null", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachByteMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		arr = []
		"Hello\nWorld".each_byte do |byte|
		  arr.push(byte)
		end
		arr
		`, []interface{}{72, 101, 108, 108, 111, 10, 87, 111, 114, 108, 100}},
		{`
		arr = []
		"Sushi 🍣".each_byte do |byte|
		  arr.push(byte)
		end
		arr
		`, []interface{}{83, 117, 115, 104, 105, 32, 240, 159, 141, 163}},
		// cases for providing an empty block
		{`
		a = "Sushi 🍣".each_byte do; end
		a.to_a
		`, []interface{}{"S", "u", "s", "h", "i", " ", "🍣"}},
		{`
		a = "Sushi 🍣".each_byte do |i|; end
		a.to_a
		`, []interface{}{"S", "u", "s", "h", "i", " ", "🍣"}},
		{`
		a = "".each_byte do; end
		a.to_a
		`, []interface{}{}},
		{`
		a = "".each_byte do |i|; end
		a.to_a
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachByteMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		"Taipei".each_byte(101) do |byte|
		  puts byte
		end
		`, "ArgumentError: Expect 0 argument. got=1", 1},
		{`"Taipei".each_byte`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachCharMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		arr = []
		"Hello\nWorld".each_char do |char|
		  arr.push(char)
		end
		arr
		`, []interface{}{"H", "e", "l", "l", "o", "\n", "W", "o", "r", "l", "d"}},
		// cases for providing an empty block
		{`
		a = "Sushi 🍣".each_char do; end; a.to_a
		`, []interface{}{"S", "u", "s", "h", "i", " ", "🍣"}},
		{`
		a = "Sushi 🍣".each_char do |i|; end; a.to_a
		`, []interface{}{"S", "u", "s", "h", "i", " ", "🍣"}},
		{`
		a = "".each_char do; end
		a.to_a
		`, []interface{}{}},
		{`
		a = "".each_char do |i|; end
		a.to_a
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachCharMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		"Taipei".each_char(101) do |char|
		  puts char
		end
		`, "ArgumentError: Expect 0 argument. got=1", 1},
		{`"Taipei".each_char`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachLineMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`
		arr = []
		"Hello\nWorld\nGoby".each_line do |line|
		  arr.push(line)
		end
		arr
		`, []interface{}{"Hello", "World", "Goby"}},
		{`
		arr = []
		"Max\vwell\nAlex\fius".each_line do |line|
		  arr.push(line)
		end
		arr
		`, []interface{}{"Max\vwell", "Alex\fius"}},
		// cases for providing an empty block
		{`
		a = "Max\vwell\nAlex\fius".each_line do; end; a.to_a
		`, []interface{}{"M", "a", "x", "\v", "w", "e", "l", "l", "\n", "A", "l", "e", "x", "\f", "i", "u", "s"}},
		{`
		a = "Max\vwell\nAlex\fius".each_line do |i|; end; a.to_a
		`, []interface{}{"M", "a", "x", "\v", "w", "e", "l", "l", "\n", "A", "l", "e", "x", "\f", "i", "u", "s"}},
		{`
		a = "".each_line do; end; a.to_a
		`, []interface{}{}},
		{`
		a = "".each_line do |i|; end; a.to_a
		`, []interface{}{}},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyArrayObject(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEachLineMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		"Taipei".each_line(101) do |line|
		  puts line
		end
		`, "ArgumentError: Expect 0 argument. got=1", 1},
		{`"Taipei".each_line`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringEndWithMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".end_with?("llo")`, true},
		{`"Hello".end_with?("Hello")`, true},
		{`"Hello".end_with?("Hello ")`, false},
		{`"哈囉！世界！".end_with?("世界！")`, true},
		{`"Hello".end_with?("ell")`, false},
		{`"哈囉！世界！".end_with?("哈囉！")`, false},
		{`"🍣Hello🍺".end_with?("🍺")`, true},
		{`"🍣Hello🍺".end_with?("🍣")`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEndWithMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Taipei".end_with?("1", "0", "1")`, "ArgumentError: Expect 1 argument. got=3", 1},
		{`"Taipei".end_with?(101)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Hello".end_with?(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Hello".end_with?(1..5)`, "TypeError: Expect argument to be String. got: Range", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringEmptyMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"".empty?`, true},
		{`"Hello".empty?`, false},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEqualMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".eql?("Hello")`, true},
		{`"Hello\nWorld".eql?("Hello\nWorld")`, true},
		{`"Hello".eql?("World")`, false},
		{`"Hello".eql?(1)`, false},
		{`"Hello".eql?(true)`, false},
		{`"Hello".eql?(2..4)`, false},
		{`"Hello🍣".eql?("Hello🍣")`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringEqualMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Hello".eql?`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Hello".eql?("Hello", "World")`, "ArgumentError: Expect 1 argument. got=2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringIncludeMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello\nWorld".include?("Hello")`, true},
		{`"Hello\nWorld".include?("Hello\nWorld")`, true},
		{`"Hello\nWorld".include?("Hello World")`, false},
		{`"Hello\nWorld".include?("Hello\nWorld!")`, false},
		{`"Hello\nWorld".include?("\n")`, true},
		{`"Hello\nWorld".include?("\r")`, false},
		{`"Hello🍣".include?("🍣")`, true},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringIncludeMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Goby".include?`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Goby".include?("Ruby", "Lang")`, "ArgumentError: Expect 1 argument. got=2", 1},
		{`"Goby".include?(2)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Goby".include?(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Goby".include?(nil)`, "TypeError: Expect argument to be String. got: Null", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringInsertMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Goby Lang".insert`, "ArgumentError: Expect 2 arguments. got=0", 1},
		{`"Taipei".insert(6, " ", "101")`, "ArgumentError: Expect 2 arguments. got=3", 1},
		{`"Taipei".insert("6", " 101")`, "TypeError: Expect argument to be Integer. got: String", 1},
		{`"Taipei".insert(6, 101)`, "TypeError: Expect insert string to be String. got: Integer", 1},
		{`"Taipei".insert(-8, "101")`, "ArgumentError: Index value out of range. got=-8", 1},
		{`"Taipei".insert(7, "101")`, "ArgumentError: Index value out of range. got=7", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringLeftJustifyMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Hello".ljust`, "ArgumentError: Expect 1..2 arguments. got=0", 1},
		{`"Hello".ljust(1, 2, 3, 4, 5)`, "ArgumentError: Expect 1..2 arguments. got=5", 1},
		{`"Hello".ljust(true)`, "TypeError: Expect justify width to be Integer. got: Boolean", 1},
		{`"Hello".ljust("World")`, "TypeError: Expect justify width to be Integer. got: String", 1},
		{`"Hello".ljust(2..5)`, "TypeError: Expect justify width to be Integer. got: Range", 1},
		{`"Hello".ljust(10, 10)`, "TypeError: Expect padding string to be String. got: Integer", 1},
		{`"Hello".ljust(10, 2..5)`, "TypeError: Expect padding string to be String. got: Range", 1},
		{`"Hello".ljust(10, true)`, "TypeError: Expect padding string to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Goby!!".match(Regexp.new("G(o)b(y)")).to_s`, "#<MatchData 0:\"Goby\" 1:\"o\" 2:\"y\">"},
		{`"Ruby".match(Regexp.new("G(o)b(y)"))`, nil},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestStringMatchFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'a'.match(Regexp.new("abc"), 1)`, "ArgumentError: Expect 1 argument. got=2", 1},
		{`'a'.match(1)`, "TypeError: Expect argument to be Regexp. got: Integer", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Ruby Lang Ruby Ruby".replace("Ru", "Go")`, "Goby Lang Goby Goby"},
		{`"🍣Ruby🍣Lang".replace("Ru", "Go")`, "🍣Goby🍣Lang"},
		{`re = Regexp.new("(Ru|ru)");"Ruby Lang ruby lang".replace(re, "Go")`, "Goby Lang Goby lang"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Invalid".replace`, "ArgumentError: Expect 2 arguments. got=0", 1},
		{`"Invalid".replace("string")`, "ArgumentError: Expect 2 arguments. got=1", 1},
		{`"Invalid".replace("string", "replace", true)`, "ArgumentError: Expect 2 arguments. got=3", 1},
		{`"Invalid".replace(true, "replacement")`, "TypeError: Expect pattern to be String or Regexp. got: Boolean", 1},
		{`"Invalid".replace("pattern", true)`, "TypeError: Expect replacement to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceOnceMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Ruby Lang Ruby Ruby".replace_once("Ru", "Go")`, "Goby Lang Ruby Ruby"},
		{`"🍣Ruby🍣Lang Ruby".replace_once("Ru", "Go")`, "🍣Goby🍣Lang Ruby"},
		{`re = Regexp.new("(Ru|ru)");"Ruby Lang ruby lang".replace_once(re, "Go")`, "Goby Lang ruby lang"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringReplaceOnceMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Invalid".replace_once`, "ArgumentError: Expect 2 arguments. got=0", 1},
		{`"Invalid".replace_once("string")`, "ArgumentError: Expect 2 arguments. got=1", 1},
		{`"Invalid".replace_once("string", "replace", true)`, "ArgumentError: Expect 2 arguments. got=3", 1},
		{`"Invalid".replace_once(true, "replacement")`, "TypeError: Expect pattern to be String or Regexp. got: Boolean", 1},
		{`"Invalid".replace_once("pattern", true)`, "TypeError: Expect replacement to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringRightJustifyFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Hello".rjust`, "ArgumentError: Expect 1..2 arguments. got=0", 1},
		{`"Hello".rjust(1, 2, 3, 4, 5)`, "ArgumentError: Expect 1..2 arguments. got=5", 1},
		{`"Hello".rjust(true)`, "TypeError: Expect justify width to be Integer. got: Boolean", 1},
		{`"Hello".rjust("World")`, "TypeError: Expect justify width to be Integer. got: String", 1},
		{`"Hello".rjust(2..5)`, "TypeError: Expect justify width to be Integer. got: Range", 1},
		{`"Hello".rjust(10, 10)`, "TypeError: Expect padding string to be String. got: Integer", 1},
		{`"Hello".rjust(10, 2..5)`, "TypeError: Expect padding string to be String. got: Range", 1},
		{`"Hello".rjust(10, true)`, "TypeError: Expect padding string to be String. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringSliceMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Goby Lang".slice`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Goby Lang".slice("Hello")`, "TypeError: Expect slice range to be Range or Integer. got: String", 1},
		{`"Goby Lang".slice(true)`, "TypeError: Expect slice range to be Range or Integer. got: Boolean", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringSplitMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Hello World".split`, "ArgumentError: Expect 1 argument. got=0", 1},
		{`"Hello World".split(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Hello World".split(123)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Hello World".split(1..2)`, "TypeError: Expect argument to be String. got: Range", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestStringStartWithMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`"Taipei".start_with("1", "0", "1")`, "ArgumentError: Expect 1 argument. got=3", 1},
		{`"Taipei".start_with(101)`, "TypeError: Expect argument to be String. got: Integer", 1},
		{`"Hello".start_with(true)`, "TypeError: Expect argument to be String. got: Boolean", 1},
		{`"Hello".start_with(1..5)`, "TypeError: Expect argument to be String. got: Range", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
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
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

// Other test

func TestStringMethodChaining(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"More test".reverse.upcase`, "TSET EROM"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestFormattedString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`String.fmt("This is %s", "goby")`, "This is goby"},
		{`String.fmt("This is %slang", "goby")`, "This is gobylang"},
		{`String.fmt("This is %s %s", "goby", "ruby")`, "This is goby ruby"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		verifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
