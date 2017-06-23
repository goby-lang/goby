package vm

import "testing"

func TestEachThroughRange(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		r = 0
		(1..(1+4)).each do |i|
		  r = r + i
		end
		r
		`, 15},
		{`
		r = 0
		a = 1
		b = 5
		(a..b).each do |i|
		  r = r + i
		end
		r
		`, 15},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}
