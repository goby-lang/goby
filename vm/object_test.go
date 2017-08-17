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
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
