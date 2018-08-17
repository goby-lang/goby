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
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataCaptures(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`'a1bca2'.match(Regexp.new('a.')).to_s`, "#<MatchData 0:\"a1\">"},
		{`'a1bca2'.match(Regexp.new('(a.)')).to_s`, "#<MatchData 0:\"a1\" 1:\"a1\">"},
		{`'a1bca2'.match(Regexp.new('(a.)(b.)')).to_s`, "#<MatchData 0:\"a1bc\" 1:\"a1\" 2:\"bc\">"},
		{`'abcd'.match(Regexp.new('a(?<first>b)(?<second>c)')).to_s`, "#<MatchData 0:\"abc\" first:\"b\" second:\"c\">"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataCapturesFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'a1bca2'.match(1, 2)`, "ArgumentError: Expect 1 argument(s). got: 2", 1},
		{`'a1bca2'.match('a.')`, "TypeError: Expect argument to be Regexp. got: String", 1},
		{`'a1bca2'.match(Regexp.new('a')).captures('a')`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`'a1bca2'.match(Regexp.new('a')).captures('a', 'b')`, "ArgumentError: Expect 0 argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataToAMethod(t *testing.T) {
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
		verifyArrayObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestMatchDataToAMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'a1bca2'.match(Regexp.new('a.')).to_a(1)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`'a1bca2'.match(Regexp.new('a.')).to_a(1, 2)`, "ArgumentError: Expect 0 argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataLengthMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`'abc'.match(Regexp.new('(a)(b)c')).length`, 3},
		{`'abc'.match(Regexp.new('(a)(b)(c)')).length`, 4},
		{`'abc'.match(Regexp.new('abc')).length`, 1},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataLengthMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'abc'.match(Regexp.new('(a)(b)c')).length(1)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`'abc'.match(Regexp.new('(a)(b)c')).length(1, 2)`, "ArgumentError: Expect 0 argument(s). got: 2", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestMatchDataToHMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{`'abcd'.match(Regexp.new('a.')).to_h`, map[string]interface{}{"0": "ab"}},
		{`'abcd'.match(Regexp.new('a(b)(c)')).to_h`, map[string]interface{}{"0": "abc", "1": "b", "2": "c"}},
		{`'abcd'.match(Regexp.new('a(?<first>b)(?<second>c)')).to_h`, map[string]interface{}{"0": "abc", "first": "b", "second": "c"}},
	}

	for i, tt := range tests {
		vm := initTestVM()
		evaluated := vm.testEval(t, tt.input, getFilename())
		verifyHashObject(t, i, evaluated, tt.expected)
		vm.checkCFP(t, i, 0)
		vm.checkSP(t, i, 1)
	}
}

func TestMatchDataToHMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`'abcd'.match(Regexp.new('a.')).to_h(1)`, "ArgumentError: Expect 0 argument(s). got: 1", 1},
		{`'abcd'.match(Regexp.new('a.')).to_h(1, 2)`, "ArgumentError: Expect 0 argument(s). got: 2", 1},
	}
	
	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
