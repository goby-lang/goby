package vm

import (
	"testing"
)

func TestRegexpClassCreation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//{`Regexp.new()`, ""},
		{`"Hello ".concat("World")`, "Hello World"},
		//{`Regexp.new('🍣Goby🍺').class`, "Regexp"},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input)
		checkExpected(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
	}
}

//func TestRegexpMatch(t *testing.T) {
//	tests := []struct {
//		input    string
//		expected interface{}
//	}{
//		{`
//		  re = Regexp.new("🍣Goby🍺"
//		  re.match?("Hello, 🍣Goby🍺!")
//		`, true},
//	}
//
//	for i, tt := range tests {
//		vm := initTestVM()
//		evaluated := vm.testEval(t, tt.input)
//		checkExpected(t, i, evaluated, tt.expected)
//		vm.checkCFP(t, i, 0)
//	}
//}
