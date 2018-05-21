package vm

import "testing"

func TestChannelClassSuperclass(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Channel.class.name`, "Class"},
		{`Channel.superclass.name`, "Object"},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestChannelCloseFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`c = Channel.new; c.close(1)`, "ArgumentError: Expect 0 arguments. got: 1", 1},
		{`c = Channel.new; c.close;c.close`, "ChannelCloseError: The channel is already closed.", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestChannelReceiveFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`c = Channel.new; c.receive(1)`, "ArgumentError: Expect 0 arguments. got: 1", 1},
		{`c = Channel.new; c.close; c.receive`, "ChannelCloseError: The channel is already closed.", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestChannelDeliverFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`c = Channel.new; c.deliver`, "ArgumentError: Expect 1 arguments. got: 0", 1},
		{`c = Channel.new; c.deliver 1, 2`, "ArgumentError: Expect 1 arguments. got: 2", 1},
		{`c = Channel.new; c.close; c.deliver 1`, "ChannelCloseError: The channel is already closed.", 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkErrorMsg(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}
