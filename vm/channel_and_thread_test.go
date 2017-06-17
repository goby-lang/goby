package vm

import "testing"

func TestObjectMutationInThread(t *testing.T) {
	tests := []struct{
		input string
		expected interface{}
	}{
		{`
		c = Channel.new

		i = 0
		thread do
		  i++
		  c.deliver(i)
		end

		# Used to block main process until thread is finished
		c.receive
		i
		`, 1},
		{`
		c = Channel.new

		i = 0
		thread do
		  i++
		  c.deliver(i)
		end

		i++
		# Used to block main process until thread is finished
		c.receive
		i
		`, 2},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		checkExpected(t, evaluated, tt.expected)
	}
}
