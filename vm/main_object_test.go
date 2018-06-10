package vm

import "testing"

func TestMainToS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`to_s`, "main"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
