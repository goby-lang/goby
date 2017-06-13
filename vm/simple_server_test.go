package vm

import "testing"

func TestServerInitialization(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		require "net/simple_server"

		s = Net::SimpleServer.new(4000)
		s.port
		`, 4000},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}
