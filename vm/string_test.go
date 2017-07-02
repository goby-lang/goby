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

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}

func TestEvalInfixStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"Stan " + "Lo"`, "Stan Lo"},
		{`"Dog" + "&" + "Cat"`, "Dog&Cat"},
		{`"Three " * 3`, "Three Three Three "},
		{`"Zero" * 0`, ""},
		{`"Minus" * 1`, "Minus"},
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
		{`"Hello"[1]`, "e"},
		{`"Hello"[5]`, NULL},
		{`"Ruby"[1] = "oo"`, "Rooby"},
		{`"Go"[2] = "by"`, "Goby"},
		{`"cat".capitalize`, "Cat"},
		{`"HELLO".capitalize`, "Hello"},
		{`"123word".capitalize`, "123word"},
		{`"Two Words".capitalize`, "Two words"},
		{`"first Lower".capitalize`, "First lower"},
		{`"all lower".capitalize`, "All lower"},
		{`"heLlo\nWoRLd".capitalize`, "Hello\nworld"},
		{`"hEllO".downcase`, "hello"},
		{`"MORE wOrds".downcase`, "more words"},
		{`"HeLlO\tWorLD".downcase`, "hello\tworld"},
		{`"hEllO".upcase`, "HELLO"},
		{`"MORE wOrds".upcase`, "MORE WORDS"},
		{`"Hello\nWorld".upcase`, "HELLO\nWORLD"},
		{`"Rooby".size`, 5},
		{`"New method".length`, 10},
		{`" ".length`, 1},
		{`"Reverse Rooby-lang".reverse`, "gnal-ybooR esreveR"},
		{`" ".reverse`, " "},
		{`"-123".reverse`, "321-"},
		{`"Hello\nWorld".reverse`, "dlroW\nolleH"},
		{`"".empty`, TRUE},
		{`"Hello".empty`, FALSE},
		{`"Hello".eql("Hello")`, TRUE},
		{`"Hello".eql("World")`, FALSE},
		{`"Hello".start_with("Hel")`, TRUE},
		{`"Hello".start_with("hel")`, FALSE},
		{`"Hello".end_with("llo")`, TRUE},
		{`"Hello".end_with("ell")`, FALSE},
		{`"Hello".insert(0, "X")`, "XHello"},
		{`"Hello".insert(2, "X")`, "HeXllo"},
		{`"Hello".insert(5, "X")`, "HelloX"},
		{`"Hello".chop`, "Hell"},
		{`"Hello".ljust(2)`, "Hello"},
		{`"Hello".ljust(7)`, "Hello  "},
		{`"Hello".rjust(2)`, "Hello"},
		{`"Hello".rjust(7)`, "  Hello"},
		{`"  Goby Lang   ".strip`, "Goby Lang"},
		{`"\nGoby Lang\r\t".strip`, "Goby Lang"},
		{`"Hello World".split("o")`, initArrayObject([]Object{initStringObject("Hell"), initStringObject(" W"), initStringObject("rld")})},
		{`"Hello".split("")`, initArrayObject([]Object{initStringObject("H"), initStringObject("e"), initStringObject("l"), initStringObject("l"), initStringObject("o")})},
		{`"Hello World".slice(1..6)`, "ello W"},
		{`"string".to_s`, "string"},
		{`"123".to_i`, 123},
		{`"string".to_i`, 0},
		{`"123string123".to_i`, 123},
		{`"string123".to_i`, 0},
		{`"Goby".to_a`, initArrayObject([]Object{initStringObject("G"), initStringObject("o"), initStringObject("b"), initStringObject("y")})},
		{`"More test".reverse.upcase`, "TSET EROM"},
		{`"Hello\nWorld".include("\n")`, true},
		{`"Hello\nWorld".include("\r")`, false},
		{`"Hello ".concat("World")`, "Hello World"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}
