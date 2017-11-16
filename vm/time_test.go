package vm

import (
	"testing"
)

func TestTimeClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Time.class.name`, "Class"},
		{`Time.superclass.name`, "Object"},
		{`Time.ancestors.to_s`, "[Time, Object]"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestTimeNew(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Time.new('2017-05-30').to_s`,
			"2017-05-30 00:00:00 +0000 UTC"},
		{`Time.new('2017-05-30 18:00').to_s`,
			"2017-05-30 18:00:00 +0000 UTC"},
		{`Time.new('2017-05-30 9:00').to_s`,
			"2017-05-30 09:00:00 +0000 UTC"},
		{`Time.new('2017-05-30 23:00 JST').to_s`,
			"2017-05-30 23:00:00 +0900 JST"},
		{`Time.new('2017-05-30 23:59 JST').to_s`,
			"2017-05-30 23:59:00 +0900 JST"},
		{`Time.new('2017-05-30 23:59:59 JST').to_s`,
			"2017-05-30 23:59:59 +0900 JST"},
		{`Time.new('2017-05-30  23:59:59  JST').to_s`,
			"2017-05-30 23:59:59 +0900 JST"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestTimeNewFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Time.new('2017-5-30').to_s`,
			"ArgumentError: Invalid time format. got=2017-5-30", 1, 1},
		{`Time.new('2017-05-30 00:00:00:00 JST')`,
			"ArgumentError: Invalid time format. got=2017-05-30 00:00:00:00 JST", 1, 1},
		{`Time.new('2017-Jan-30 00:0 JST')`,
			"ArgumentError: Invalid time format. got=2017-Jan-30 00:0 JST", 1, 1},
		{`Time.new('09:00 JST')`,
			"ArgumentError: Invalid time format. got=09:00 JST", 1, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
