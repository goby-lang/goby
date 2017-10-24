package vm

import (
	"testing"
)

func TestMatchDataClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`MatchData.class.name`, "Class"},
		{`String.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataCaptures(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'a1bca2'.match(Regexp.new('a.')).to_s`, "#<MatchData \"a1\">"},
		{`'a1bca2'.match(Regexp.new('(a.)')).to_s`, "#<MatchData \"a1\" 1:\"a1\">"},
		{`'a1bca2'.match(Regexp.new('(a.)(b.)')).to_s`, "#<MatchData \"a1bc\" 1:\"a1\" 2:\"bc\">"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataCapturesFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'a1bca2'.match(1, 2)`, "ArgumentError: Expect 1 argument. got=2", 1, 1},
		{`'a1bca2'.match('a.')`, "TypeError: Expect argument to be Regexp. got: String", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataToA(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`'a1bca2'.match(Regexp.new('a.')).to_a`, []interface{}{"a1"}},
		{`'a1bca2'.match(Regexp.new('(a.)')).to_a`, []interface{}{"a1", "a1"}},
		{`'a1bca2'.match(Regexp.new('(a.)(b.)')).to_a`, []interface{}{"a1bc", "a1", "bc"}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		testArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestMatchDataToAFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'a1bca2'.match(Regexp.new('a.')).to_a(1)`, "ArgumentError: Expect 0 argument. got=1", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
