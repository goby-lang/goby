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

func TestObjectMutationInThread(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		c = Channel.new

		i = 0
		thread do
		  i += 1
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
		  i += 1
		  c.deliver(i)
		end

		# if we put i += 1 here then it will execute along with other thread,
		# which will cause a race condition.
		# The following "receive" is needed to block the main process until thread is finished
		c.receive
		i += 1
		i
		`, 2},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}

func TestObjectDeliveryBetweenThread(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
		c = Channel.new

		thread do
		  s = "123"
		  c.deliver(s)
		end

		c.receive
		`, "123"},
		{`
		c = Channel.new

		thread do
		  h = "Hello"
		  w = "World"
		  c.deliver(h)
		  c.deliver(w)
		end

		h = c.receive
		w = c.receive

		h + " " + w
		`, "Hello World"},
		{`
		class Foo
		  def bar
		    100
		  end
		end

		c = Channel.new

		thread do
		  f = Foo.new
		  c.deliver(f)
		end

		c.receive.bar
		`, 100},
		{`
		c = Channel.new
		c2 = Channel.new

		thread do
		  1001.times do |i| # i start from 0 to 1000
		  	c.deliver(i)
		  end

		  c2.receive
		  c.deliver(100)
		end

		r = 0

		1001.times do
		  r = r + c.receive
		end

		c2.deliver(true) # block thread until it finishes the loop
		r + c.receive
		`, 500600},
		{`
		c = Channel.new

		1001.times do |i| # i start from 0 to 1000
		  thread do
		  	c.deliver(i)
		  end
		end

		r = 0
		1001.times do
		  r = r + c.receive
		end

		r
		`, 500500},
	}

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		VerifyExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
