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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestStringToArrayMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Dog" == "Dog"`, true},
		{`"1234" > "123"`, true},
		{`"123" > "1235"`, false},
		{`"1234" < "123"`, false},
		{`"1234" < "12jdkfj3"`, true},
		{`"1234" != "123"`, true},
		{`"123" != "123"`, false},
		{`"1234" <=> "1234"`, 0},
		{`"1234" <=> "4"`, -1},
		{`"abcdef" <=> "abcde"`, 1},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestStringConparisonFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"a" < 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"a" > 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"a" == 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"a" <=> 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"a" != 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
	}
	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
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
		{`"Hello\nWorld"[5]`, "\n"},
		{`"Ruby"[1] = "oo"`, "Rooby"},
		{`"Go"[2] = "by"`, "Goby"},
		{`"Ruby"[-3] = "oo"`, "Rooby"},
		{`"Hello"[-5] = "Tr"`, "Trello"},
		{`"Hello\nWorld"[5] = " "`, "Hello World"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStringOperationFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Taipei" + 101`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(101))},
		{`"Taipei" * "101"`, initErrorObject(TypeErrorClass, "Expect argument to be Integer. got=%T", vm.initStringObject("101"))},
		{`"Taipei" * (-101)`, initErrorObject(ArgumentErrorClass, "Second argument must be greater than or equal to 0. got=%v", -101)},
		{`"Taipei"[1] = 1`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"Taipei"[1] = true`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		// TODO: Implement test for empty index or wrong index type
		//{`"Taipei"[]`, initErrorObject(ArgumentErrorClass, "Expect 1 argument. got=%v", "0")},
		// {`"Taipei"[true] = 101`, newError("expect argument to be Integer type")},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestCountingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"abcde".count`, 5},
		{`"哈囉！世界！".count`, 6},
		{`"Hello\nWorld".count`, 11},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestCapitalizingString(t *testing.T) {
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestConcatenatingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello ".concat("World")`, "Hello World"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestConcatenatingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"a".concat(1)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"a".concat(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		{`"a".concat(nil)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", NULL)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestDeletingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello hello HeLlo".delete("el")`, "Hlo hlo HeLlo"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestDeleteStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Hello hello HeLlo".delete(1)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(1))},
		{`"Hello hello HeLlo".delete(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		{`"Hello hello HeLlo".delete(nil)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", NULL)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestStringEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"".empty`, true},
		{`"Hello".empty`, false},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStringEqual(t *testing.T) {
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStringStartWith(t *testing.T) {
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStringStartWithFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Taipei".start_with(101)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(101))},
		{`"Hello".start_with(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		{`"Hello".start_with(1..5)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initRangeObject(1, 5))},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestStringEndWith(t *testing.T) {
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestStringEndWithFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Taipei".end_with(101)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(101))},
		{`"Hello".end_with(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		{`"Hello".end_with(1..5)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initRangeObject(1, 5))},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestInsertingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".insert(0, "X")`, "XHello"},
		{`"Hello".insert(2, "X")`, "HeXllo"},
		{`"Hello".insert(5, "X")`, "HelloX"},
		{`"Hello".insert(-2, "X")`, "HelXlo"},
		{`"Hello".insert(-6, "X")`, "XHello"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestInsertingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Taipei".insert("6", " 101")`, initErrorObject(TypeErrorClass, "Expect index to be Integer. got=%T", vm.initStringObject("6"))},
		{`"Taipei".insert(6, 101)`, initErrorObject(TypeErrorClass, "Expect insert string to be String. got=%T", vm.initIntegerObject(101))},
		{`"Taipei".insert(-8, "101")`, initErrorObject(ArgumentErrorClass, "Index value out of range. got=%v", "-8")},
		{`"Taipei".insert(7, "101")`, initErrorObject(ArgumentErrorClass, "Index value out of range. got=%v", "7")},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestChoppingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".chop`, "Hell"},
		{`"Hello\n".chop`, "Hello"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestLeftJustifyingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".ljust(2)`, "Hello"},
		{`"Hello".ljust(7)`, "Hello  "},
		{`"Hello".ljust(10, "xo")`, "Helloxoxox"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestLeftJustifyStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Hello".ljust(true)`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", TRUE)},
		{`"Hello".ljust("World")`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", vm.initStringObject("World"))},
		{`"Hello".ljust(2..5)`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", vm.initRangeObject(2, 5))},
		{`"Hello".ljust(10, 10)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", vm.initIntegerObject(10))},
		{`"Hello".ljust(10, 2..5)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", vm.initRangeObject(2, 5))},
		{`"Hello".ljust(10, true)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", TRUE)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestRightJustifyingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".rjust(2)`, "Hello"},
		{`"Hello".rjust(7)`, "  Hello"},
		{`"Hello".rjust(10, "xo")`, "xoxoxHello"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestRightJustifyStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Hello".rjust(true)`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", TRUE)},
		{`"Hello".rjust("World")`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", vm.initStringObject("World"))},
		{`"Hello".rjust(2..5)`, initErrorObject(TypeErrorClass, "Expect justify width to be Integer. got=%T", vm.initRangeObject(2, 5))},
		{`"Hello".rjust(10, 10)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", vm.initIntegerObject(10))},
		{`"Hello".rjust(10, 2..5)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", vm.initRangeObject(2, 5))},
		{`"Hello".rjust(10, true)`, initErrorObject(TypeErrorClass, "Expect padding string to be String. got=%T", TRUE)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestStrippingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"  Goby Lang   ".strip`, "Goby Lang"},
		{`"\nGoby Lang\r\t".strip`, "Goby Lang"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSplittingString(t *testing.T) {
	vm := initTestVM()
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
	}

	for i, tt := range tests {
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSplittingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Hello World".split(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
		{`"Hello World".split(123)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(123))},
		{`"Hello World".split(1..2)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initRangeObject(1, 2))},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestSlicingString(t *testing.T) {
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
		{`"Hello World".slice(4)`, "o"},
		{`"Hello\nWorld".slice(5)`, "\n"},
		{`"Hello World".slice(-3)`, "r"},
		{`"Hello World".slice(-11)`, "H"},
		{`"Hello World".slice(-12)`, nil},
		{`"Hello World".slice(11)`, nil},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestSlicingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Goby Lang".slice("Hello")`, initErrorObject(ArgumentErrorClass, "Expect slice range is Range or Integer type. got=%T", vm.initStringObject("Hello"))},
		{`"Goby Lang".slice(true)`, initErrorObject(ArgumentErrorClass, "Expect slice range is Range or Integer type. got=%T", TRUE)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestReplacingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Hello".replace("World")`, "World"},
		{`"您好".replace("再見")`, "再見"},
		{`"Ruby\nLang".replace("Goby\nLang")`, "Goby\nLang"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestReplacingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Taipei".replace(101)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", vm.initIntegerObject(101))},
		{`"Taipei".replace(true)`, initErrorObject(TypeErrorClass, "Expect argument to be String. got=%T", TRUE)},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
	}
}

func TestUpcasingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hEllO".upcase`, "HELLO"},
		{`"MORE wOrds".upcase`, "MORE WORDS"},
		{`"Hello\nWorld".upcase`, "HELLO\nWORLD"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestDowncasingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hEllO".downcase`, "hello"},
		{`"MORE wOrds".downcase`, "more words"},
		{`"HeLlO\tWorLD".downcase`, "hello\tworld"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestSizeAndLengthOfString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Rooby".size`, 5},
		{`"New method".length`, 10},
		{`" ".length`, 1},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestReversingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Reverse Rooby-lang".reverse`, "gnal-ybooR esreveR"},
		{`" ".reverse`, " "},
		{`"-123".reverse`, "321-"},
		{`"Hello\nWorld".reverse`, "dlroW\nolleH"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestIncludingString(t *testing.T) {
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
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

func TestGlobalSubstitutingString(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Ruby".gsub("Ru", "Go")`, "Goby"},
		{`"Hello World".gsub(" ", "\n")`, "Hello\nWorld"},
		{`"Hello World".gsub("Hello", "Goby")`, "Goby World"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
	}
}

func TestGlobalSubstitutingStringFail(t *testing.T) {
	vm := initTestVM()
	testsFail := []struct {
		input    string
		expected *Error
	}{
		{`"Ruby".gsub()`, initErrorObject(ArgumentErrorClass, "Expect to have 2 arguments. got=%v", 0)},
		{`"Ruby".gsub("Ru")`, initErrorObject(ArgumentErrorClass, "Expect to have 2 arguments. got=%v", 1)},
		{`"Ruby".gsub(123, "Go")`, initErrorObject(TypeErrorClass, "Expect pattern to be String. got=%T", vm.initIntegerObject(123))},
		{`"Ruby".gsub("Ru", 456)`, initErrorObject(TypeErrorClass, "Expect replacement to be String. got=%T", vm.initIntegerObject(456))},
	}

	for _, tt := range testsFail {
		evaluated := vm.testEval(t, tt.input)
		err, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("Expect error. got=%T (%+v)", err, err)
		}
		if err.Message != tt.expected.Message {
			t.Errorf("Expect error message \"%s\". got=\"%s\"", tt.expected.Message, err.Message)
		}
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
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}
