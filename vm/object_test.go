package vm

import "testing"

func TestObjectClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Object.class.name`, "Class"},
		{`Object.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectTapMethod(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			a = 1
			a.tap do |int|
				int + 1
			end
`, 1},
		{
			`
			a = 1
			b = 2
			a.tap do |int|
				b = int + b
			end
			b
`, 3},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectTapMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`Object.new.tap`, "InternalError: Can't yield without a block", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

