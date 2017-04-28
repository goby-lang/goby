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
		testStringObject(t, evaluated, tt.expected)
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
		{`"cat".capitalize`, "Cat"},
		{`"HELLO".capitalize`, "Hello"},
		{`"123word".capitalize`, "123word"},
		{`"Two Words".capitalize`, "Two words"},
		{`"first Lower".capitalize`, "First lower"},
		{`"all lower".capitalize`, "All lower"},
		{`"hEllO".downcase`, "hello"},
		{`"MORE wOrds".downcase`, "more words"},
		{`"hEllO".upcase`, "HELLO"},
		{`"MORE wOrds".upcase`, "MORE WORDS"},
		{`"Rooby".size`, 5},
		{`"New method".length`, 10},
		{`" ".length`, 1},
		{`"Reverse Rooby-lang".reverse`, "gnal-ybooR esreveR"},
		{`" ".reverse`, " "},
		{`"-123".reverse`, "321-"},
		{`"string".to_s`, "string"},
		{`"More test".reverse.upcase`, "TSET EROM"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		switch tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, tt.expected.(bool))
		case string:
			testStringObject(t, evaluated, tt.expected.(string))
		case int:
			testIntegerObject(t, evaluated, tt.expected.(int))
		}
	}
}
