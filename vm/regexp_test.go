package vm

import (
	"testing"
)

func TestRegexpClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Regexp.class.name`, "Class"},
		{`Regexp.superclass.name`, "Object"},
		{`Regexp.ancestors.to_s`, "[Regexp, Object]"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

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
		evaluated := vm.testEval(t, tt.input, getFilename())
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
//		evaluated := vm.testEval(t, tt.input, getFilename())
//		checkExpected(t, i, evaluated, tt.expected)
//		vm.checkCFP(t, i, 0)
//	}
//}
